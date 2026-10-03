package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
)

// ProxyEngine coordinates primary upstream, fallback failovers, caching, and FinOps.
type ProxyEngine struct {
	circuitBreaker domain.CircuitBreaker
	fallbackRouter *FallbackRouter
	cacheRepo      *repository.CacheRepository
	telemetryHub   *repository.TelemetryHub
	finopsService  *FinOpsService
	guardrail      *GuardrailService
	httpClient     *http.Client
	openAIKey      string
	primaryURL     string
}

// NewProxyEngine creates a new proxy coordinator.
func NewProxyEngine(cb domain.CircuitBreaker, fb *FallbackRouter, cr *repository.CacheRepository, th *repository.TelemetryHub, fs *FinOpsService, gr *GuardrailService, primaryURL, openAIKey string, timeout time.Duration) *ProxyEngine {
	return &ProxyEngine{
		circuitBreaker: cb, fallbackRouter: fb, cacheRepo: cr, telemetryHub: th,
		finopsService: fs, guardrail: gr, httpClient: &http.Client{Timeout: timeout},
		openAIKey: openAIKey, primaryURL: primaryURL,
	}
}

// ProcessChat handles chat completions with cache, resilience, and telemetry.
func (p *ProxyEngine) ProcessChat(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	if p.finopsService != nil {
		if err := p.finopsService.CheckBudget(ctx); err != nil {
			return nil, err
		}
	}
	var redactions map[string]string
	if p.guardrail != nil {
		gRes := p.guardrail.InspectAndSanitize(req.Messages)
		if gRes.IsBlocked {
			return nil, domain.NewAppException(constants.ErrCodePolicyViolation, gRes.BlockReason, http.StatusUnprocessableEntity, nil)
		}
		req.Messages, redactions = gRes.SanitizedMessages, gRes.Redactions
	}
	start := time.Now()
	if p.cacheRepo != nil {
		if cached, hit, _ := p.cacheRepo.Get(ctx, req); hit && cached != nil {
			p.emitTelemetry(ctx, req.Model, cached.Usage.PromptTokens, cached.Usage.CompletionTokens, true, "redis-cache", start)
			return p.unmaskResult(cached, redactions), nil
		}
	}
	resp, provider, err := p.executeRouting(ctx, req)
	if err != nil {
		return nil, err
	}
	if p.cacheRepo != nil {
		_ = p.cacheRepo.Set(ctx, req, resp, 24*time.Hour)
	}
	p.emitTelemetry(ctx, req.Model, resp.Usage.PromptTokens, resp.Usage.CompletionTokens, false, provider, start)
	return p.unmaskResult(resp, redactions), nil
}

func (p *ProxyEngine) unmaskResult(resp *domain.ChatCompletionResponse, redactions map[string]string) *domain.ChatCompletionResponse {
	if p.guardrail == nil || len(redactions) == 0 || resp == nil {
		return resp
	}
	for i := range resp.Choices {
		resp.Choices[i].Message.Content = p.guardrail.UnmaskResponse(resp.Choices[i].Message.Content, redactions)
	}
	return resp
}

func (p *ProxyEngine) executeRouting(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, string, error) {
	if !p.circuitBreaker.CanExecute() {
		resp, err := p.fallbackRouter.ExecuteFallback(ctx, req)
		return resp, "gemini-fallback", err
	}
	if p.openAIKey == "" {
		p.circuitBreaker.RecordSuccess()
		return p.executeMockPrimary(req), "openai-mock", nil
	}
	resp, err := p.callPrimaryOpenAI(ctx, req)
	if err != nil {
		p.circuitBreaker.RecordFailure()
		fallbackResp, fbErr := p.fallbackRouter.ExecuteFallback(ctx, req)
		return fallbackResp, "gemini-fallback", fbErr
	}
	p.circuitBreaker.RecordSuccess()
	return resp, "openai-primary", nil
}

func (p *ProxyEngine) emitTelemetry(ctx context.Context, model string, pTokens, cTokens int, isHit bool, provider string, start time.Time) {
	cost, saved := domain.CalculateCostEstimates(model, pTokens, cTokens, isHit, provider)
	if p.finopsService != nil && !isHit {
		_ = p.finopsService.RecordSpend(ctx, cost)
	}
	if p.telemetryHub == nil {
		return
	}
	_ = p.telemetryHub.Publish(ctx, domain.TelemetryEvent{
		ID: fmt.Sprintf("evt-%d", time.Now().UnixNano()), Timestamp: time.Now().UTC(),
		Model: model, PromptTokens: pTokens, CompletionTokens: cTokens,
		TotalTokens: pTokens + cTokens, EstimatedCostUSD: cost, DollarsSavedUSD: saved,
		LatencyMs: time.Since(start).Milliseconds(), IsCacheHit: isHit,
		Provider: provider, CircuitState: string(p.circuitBreaker.GetStatus().State),
	})
}

func (p *ProxyEngine) executeMockPrimary(req domain.ChatCompletionRequest) *domain.ChatCompletionResponse {
	return &domain.ChatCompletionResponse{
		ID: fmt.Sprintf("chatcmpl-%d", time.Now().Unix()), Object: "chat.completion",
		Created: time.Now().Unix(), Model: req.Model,
		Choices: []domain.ChatChoice{{
			Index: 0, Message: domain.ChatMessage{Role: "assistant", Content: "[AegisAI Primary] Gateway processed successfully."},
			FinishReason: "stop",
		}},
		Usage: domain.UsageMetrics{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30},
	}
}

func (p *ProxyEngine) callPrimaryOpenAI(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	bodyBytes, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/chat/completions", p.primaryURL), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("proxyEngine.NewRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.openAIKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream status error: %v", err)
	}
	defer resp.Body.Close()

	var result domain.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("proxyEngine.Decode: %w", err)
	}
	return &result, nil
}
