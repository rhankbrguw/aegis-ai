package handler

import (
	"encoding/json"
	"net/http"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// HealthHandler handles system health status checks.
type HealthHandler struct{}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP writes the standard health envelope.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	payload := domain.NewSuccessResponse("OK", constants.MsgHealthOK, map[string]string{
		"service": "aegis-gateway",
		"status":  "healthy",
	}, r.Header.Get("X-Request-ID"))

	_ = json.NewEncoder(w).Encode(payload)
}
