package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/handler"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func TestChatHandler_StreamingSSE(t *testing.T) {
	cb := service.NewMemoryCircuitBreaker(3, time.Minute)
	fb := service.NewFallbackRouter("https://generativelanguage.googleapis.com", "", time.Second)
	guardrail := service.NewGuardrailService()
	proxy := service.NewProxyEngine(cb, fb, nil, nil, nil, guardrail, "https://api.openai.com/v1", "", time.Second)
	streamingSvc := service.NewStreamingService(guardrail, nil)

	hdl := handler.NewChatHandler(proxy, streamingSvc)

	payload := domain.ChatCompletionRequest{
		Model:    "gpt-4o",
		Stream:   true,
		Messages: []domain.ChatMessage{{Role: "user", Content: "Stream test"}},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, constants.RouteChatCompletions, bytes.NewReader(body))
	rec := httptest.NewRecorder()

	hdl.Complete(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for stream, got %d", rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("expected text/event-stream content type, got: %s", contentType)
	}

	resBody := rec.Body.String()
	if !strings.Contains(resBody, "data: [DONE]") {
		t.Errorf("expected stream to contain [DONE] sentinel")
	}
}
