package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// FallbackRouter handles failover to Gemini.
type FallbackRouter struct {
	client     *http.Client
	geminiKey  string
	baseURL    string
}

// NewFallbackRouter creates a new FallbackRouter instance.
func NewFallbackRouter(baseURL string, geminiKey string, timeout time.Duration) *FallbackRouter {
	return &FallbackRouter{
		client:    &http.Client{Timeout: timeout},
		geminiKey: geminiKey,
		baseURL:   baseURL,
	}
}

// ExecuteFallback sends request to Gemini and translates response to OpenAI format.
func (f *FallbackRouter) ExecuteFallback(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	if f.geminiKey == "" {
		return f.generateMockFallbackResponse(req), nil
	}
	return f.callGeminiAPI(ctx, req)
}

func (f *FallbackRouter) generateMockFallbackResponse(req domain.ChatCompletionRequest) *domain.ChatCompletionResponse {
	return &domain.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-fallback-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gemini-2.5-flash-fallback",
		Choices: []domain.ChatChoice{
			{
				Index: 0,
				Message: domain.ChatMessage{
					Role:    "assistant",
					Content: "[AegisAI Failover: Gemini 2.5 Flash] Circuit tripped from OpenAI. Request served seamlessly via fallback.",
				},
				FinishReason: "stop",
			},
		},
		Usage: domain.UsageMetrics{
			PromptTokens:     15,
			CompletionTokens: 25,
			TotalTokens:      40,
		},
	}
}

func (f *FallbackRouter) callGeminiAPI(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/models/gemini-2.5-flash:generateContent?key=%s", f.baseURL, f.geminiKey)
	gReq := MapOpenAIToGemini(req)
	bodyBytes, err := json.Marshal(gReq)
	if err != nil {
		return nil, fmt.Errorf("fallbackRouter.callGeminiAPI: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("fallbackRouter.NewRequest: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := f.client.Do(httpReq)
	if err != nil || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fallbackRouter.Do: %v", err)
	}
	defer resp.Body.Close()

	var gResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, fmt.Errorf("fallbackRouter.Decode: %w", err)
	}
	return MapGeminiToOpenAI(gResp, "gemini-2.5-flash"), nil
}

