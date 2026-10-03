package domain

import (
	"encoding/json"
	"fmt"
)

// ChatCompletionChunk represents an SSE chunk payload following OpenAI streaming spec.
type ChatCompletionChunk struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []ChatChunkChoice `json:"choices"`
}

// ChatChunkChoice represents an individual choice candidate in an SSE stream chunk.
type ChatChunkChoice struct {
	Index        int              `json:"index"`
	Delta        ChatMessageDelta `json:"delta"`
	FinishReason *string          `json:"finish_reason"`
}

// ChatMessageDelta contains incremental content tokens emitted during streaming.
type ChatMessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// FormatSSEChunk formats a chunk into standard Server-Sent Event data line.
func FormatSSEChunk(chunk ChatCompletionChunk) []byte {
	bytes, _ := json.Marshal(chunk)
	return []byte(fmt.Sprintf("data: %s\n\n", string(bytes)))
}

// FormatSSEDone formats the stream termination signal.
func FormatSSEDone() []byte {
	return []byte("data: [DONE]\n\n")
}
