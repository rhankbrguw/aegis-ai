package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
)

// StreamingService handles token-by-token SSE streaming and TTFT tracking.
type StreamingService struct {
	guardrail    *GuardrailService
	telemetryHub *repository.TelemetryHub
}

// NewStreamingService constructs a StreamingService.
func NewStreamingService(gr *GuardrailService, th *repository.TelemetryHub) *StreamingService {
	return &StreamingService{guardrail: gr, telemetryHub: th}
}

// StreamChat executes streaming token generation, emitting chunks to chunkSink.
func (s *StreamingService) StreamChat(ctx context.Context, req domain.ChatCompletionRequest, sink func(domain.ChatCompletionChunk) error) error {
	if s.guardrail != nil {
		gRes := s.guardrail.InspectAndSanitize(req.Messages)
		if gRes.IsBlocked {
			return domain.NewAppException(constants.ErrCodePolicyViolation, gRes.BlockReason, 422, nil)
		}
		req.Messages = gRes.SanitizedMessages
	}

	start := time.Now()
	content := "[AegisAI Stream] Gateway token streaming active."
	tokens := strings.Fields(content)
	streamID := fmt.Sprintf("chatcmpl-stream-%d", time.Now().UnixNano())

	for i, token := range tokens {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			chunk := s.buildChunk(streamID, req.Model, token+" ", i == len(tokens)-1)
			if err := sink(chunk); err != nil {
				return err
			}
			if i == 0 {
				s.recordTTFT(ctx, req.Model, time.Since(start).Milliseconds())
			}
		}
	}
	return nil
}

func (s *StreamingService) buildChunk(id, model, text string, isLast bool) domain.ChatCompletionChunk {
	var finishReason *string
	if isLast {
		stop := "stop"
		finishReason = &stop
	}
	return domain.ChatCompletionChunk{
		ID: id, Object: "chat.completion.chunk", Created: time.Now().Unix(), Model: model,
		Choices: []domain.ChatChunkChoice{{
			Index: 0, Delta: domain.ChatMessageDelta{Content: text}, FinishReason: finishReason,
		}},
	}
}

func (s *StreamingService) recordTTFT(ctx context.Context, model string, ttftMs int64) {
	if s.telemetryHub == nil {
		return
	}
	_ = s.telemetryHub.Publish(ctx, domain.TelemetryEvent{
		ID: fmt.Sprintf("evt-ttft-%d", time.Now().UnixNano()), Timestamp: time.Now().UTC(),
		Model: model, LatencyMs: ttftMs, Provider: "stream-multiplexer",
	})
}
