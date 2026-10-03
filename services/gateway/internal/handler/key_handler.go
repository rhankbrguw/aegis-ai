package handler

import (
	"encoding/json"
	"net/http"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

// KeyHandler handles API key generation and management requests.
type KeyHandler struct {
	authService *service.AuthService
}

// NewKeyHandler constructs a KeyHandler instance.
func NewKeyHandler(authService *service.AuthService) *KeyHandler {
	return &KeyHandler{authService: authService}
}

// CreateKey processes new API key generation requests.
func (h *KeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req domain.CreateKeyRequest

	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TenantID == "" || req.Name == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := domain.NewErrorResponse(constants.ErrCodeValidation, constants.ErrMsgValidation, nil, r.Header.Get("X-Request-ID"))
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	keyResp, err := h.authService.GenerateKey(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := domain.NewErrorResponse(constants.ErrCodeInternalError, constants.ErrMsgInternalError, nil, r.Header.Get("X-Request-ID"))
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := domain.NewSuccessResponse("CREATED", "API key generated successfully", keyResp, r.Header.Get("X-Request-ID"))
	_ = json.NewEncoder(w).Encode(resp)
}
