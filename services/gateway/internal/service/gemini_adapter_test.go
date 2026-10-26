package service_test

import (
	"testing"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func TestGeminiAdapter_MapOpenAIToGemini(t *testing.T) {
	req := domain.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []domain.ChatMessage{
			{Role: "system", Content: "You are a tactical sentinel."},
			{Role: "user", Content: "Report perimeter status."},
			{Role: "assistant", Content: "Perimeter secure."},
		},
	}

	gReq := service.MapOpenAIToGemini(req)

	if gReq.SystemInstruction == nil || len(gReq.SystemInstruction.Parts) == 0 {
		t.Fatalf("expected system instruction to be populated")
	}
	if gReq.SystemInstruction.Parts[0].Text != "You are a tactical sentinel." {
		t.Errorf("system instruction mismatch")
	}

	if len(gReq.Contents) != 2 {
		t.Fatalf("expected 2 content messages, got %d", len(gReq.Contents))
	}
	if gReq.Contents[0].Role != "user" || gReq.Contents[1].Role != "model" {
		t.Errorf("expected role user and model, got %s and %s", gReq.Contents[0].Role, gReq.Contents[1].Role)
	}
}

func TestGeminiAdapter_MapGeminiToOpenAI(t *testing.T) {
	var gResp service.GeminiResponse
	gResp.Candidates = append(gResp.Candidates, struct {
		Content struct {
			Parts []service.GeminiPart `json:"parts"`
			Role  string               `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	}{
		Content: struct {
			Parts []service.GeminiPart `json:"parts"`
			Role  string               `json:"role"`
		}{
			Parts: []service.GeminiPart{{Text: "All sensors operational."}},
			Role:  "model",
		},
		FinishReason: "STOP",
	})
	gResp.UsageMetadata.PromptTokenCount = 12
	gResp.UsageMetadata.CandidatesTokenCount = 28
	gResp.UsageMetadata.TotalTokenCount = 40

	oResp := service.MapGeminiToOpenAI(gResp, "gemini-2.5-flash")

	if len(oResp.Choices) == 0 {
		t.Fatalf("expected non-empty choices")
	}
	if oResp.Choices[0].Message.Content != "All sensors operational." {
		t.Errorf("unexpected content: %s", oResp.Choices[0].Message.Content)
	}
	if oResp.Usage.TotalTokens != 40 {
		t.Errorf("expected total tokens 40, got %d", oResp.Usage.TotalTokens)
	}
}
