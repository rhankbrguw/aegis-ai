package handler

import (
	"encoding/json"
	"net/http"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// CircuitHandler manages circuit status inspection and manual overrides.
type CircuitHandler struct {
	circuitBreaker domain.CircuitBreaker
}

// NewCircuitHandler constructs a CircuitHandler.
func NewCircuitHandler(cb domain.CircuitBreaker) *CircuitHandler {
	return &CircuitHandler{circuitBreaker: cb}
}

// GetState returns the current status snapshot.
func (h *CircuitHandler) GetState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := h.circuitBreaker.GetStatus()
	resp := domain.NewSuccessResponse("OK", "Circuit state retrieved", status, r.Header.Get("X-Request-ID"))
	_ = json.NewEncoder(w).Encode(resp)
}

// OverrideState handles manual circuit trips and kill-switch toggling.
func (h *CircuitHandler) OverrideState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		State      *domain.CircuitState `json:"state"`
		KillSwitch *bool                `json:"kill_switch"`
	}

	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := domain.NewErrorResponse(constants.ErrCodeValidation, constants.ErrMsgValidation, nil, r.Header.Get("X-Request-ID"))
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	if req.State != nil {
		h.circuitBreaker.SetManualOverride(*req.State)
	}
	if req.KillSwitch != nil {
		h.circuitBreaker.SetKillSwitch(*req.KillSwitch)
	}

	resp := domain.NewSuccessResponse("OK", "Circuit state updated", h.circuitBreaker.GetStatus(), r.Header.Get("X-Request-ID"))
	_ = json.NewEncoder(w).Encode(resp)
}
