package service

import (
	"fmt"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// GeminiRequest represents Google Gemini generateContent payload schema.
type GeminiRequest struct {
	Contents          []GeminiContent    `json:"contents"`
	SystemInstruction *GeminiInstruction `json:"systemInstruction,omitempty"`
}

// GeminiContent holds role and token parts.
type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiInstruction defines system instructions.
type GeminiInstruction struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart contains text payload.
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents the output returned by Gemini generateContent.
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []GeminiPart `json:"parts"`
			Role  string       `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
}

// MapOpenAIToGemini converts an OpenAI ChatCompletionRequest to Gemini schema.
func MapOpenAIToGemini(req domain.ChatCompletionRequest) GeminiRequest {
	var gContents []GeminiContent
	var gSys *GeminiInstruction

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			gSys = &GeminiInstruction{Parts: []GeminiPart{{Text: msg.Content}}}
			continue
		}
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		gContents = append(gContents, GeminiContent{
			Role:  role,
			Parts: []GeminiPart{{Text: msg.Content}},
		})
	}
	return GeminiRequest{Contents: gContents, SystemInstruction: gSys}
}

// MapGeminiToOpenAI converts a GeminiResponse into standard OpenAI format.
func MapGeminiToOpenAI(gResp GeminiResponse, model string) *domain.ChatCompletionResponse {
	text := "[Empty response]"
	if len(gResp.Candidates) > 0 && len(gResp.Candidates[0].Content.Parts) > 0 {
		text = gResp.Candidates[0].Content.Parts[0].Text
	}
	pTokens := gResp.UsageMetadata.PromptTokenCount
	cTokens := gResp.UsageMetadata.CandidatesTokenCount

	return &domain.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-gemini-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []domain.ChatChoice{{
			Index:        0,
			Message:      domain.ChatMessage{Role: "assistant", Content: text},
			FinishReason: "stop",
		}},
		Usage: domain.UsageMetrics{
			PromptTokens:     pTokens,
			CompletionTokens: cTokens,
			TotalTokens:      pTokens + cTokens,
		},
	}
}
