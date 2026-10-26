package handler

import (
	"encoding/json"
	"net/http"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

// ChatHandler handles OpenAI-compatible completions requests.
type ChatHandler struct {
	proxyEngine      *service.ProxyEngine
	streamingService *service.StreamingService
}

// NewChatHandler constructs a ChatHandler.
func NewChatHandler(pe *service.ProxyEngine, ss *service.StreamingService) *ChatHandler {
	return &ChatHandler{proxyEngine: pe, streamingService: ss}
}

// Complete handles POST /v1/chat/completions.
func (h *ChatHandler) Complete(w http.ResponseWriter, r *http.Request) {
	var req domain.ChatCompletionRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusUnprocessableEntity, constants.ErrCodeValidation, constants.ErrMsgValidation)
		return
	}
	if len(req.Messages) == 0 {
		h.respondError(w, r, http.StatusUnprocessableEntity, constants.ErrCodeValidation, "At least one message is required.")
		return
	}
	if req.Stream && h.streamingService != nil {
		h.handleStream(w, r, req)
		return
	}
	h.handleSync(w, r, req)
}

func (h *ChatHandler) handleSync(w http.ResponseWriter, r *http.Request, req domain.ChatCompletionRequest) {
	w.Header().Set("Content-Type", "application/json")
	result, err := h.proxyEngine.ProcessChat(r.Context(), req)
	if err != nil {
		h.respondError(w, r, http.StatusBadGateway, constants.ErrCodeUpstreamError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}

func (h *ChatHandler) handleStream(w http.ResponseWriter, r *http.Request, req domain.ChatCompletionRequest) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.respondError(w, r, http.StatusInternalServerError, constants.ErrCodeInternalError, "Streaming unsupported")
		return
	}
	_ = h.streamingService.StreamChat(r.Context(), req, func(chunk domain.ChatCompletionChunk) error {
		_, _ = w.Write(domain.FormatSSEChunk(chunk))
		flusher.Flush()
		return nil
	})
	_, _ = w.Write(domain.FormatSSEDone())
	flusher.Flush()
}

func (h *ChatHandler) respondError(w http.ResponseWriter, r *http.Request, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := domain.NewErrorResponse(code, msg, nil, r.Header.Get("X-Request-ID"))
	_ = json.NewEncoder(w).Encode(resp)
}

