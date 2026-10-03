package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

type contextKey string

const TenantContextKey contextKey = "aegis_tenant_id"

// AuthMiddleware enforces API Key validation and rate limiting.
type AuthMiddleware struct {
	authService *service.AuthService
}

// NewAuthMiddleware constructs the auth middleware.
func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireKey verifies authorization header before forwarding request.
func (m *AuthMiddleware) RequireKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("X-Aegis-Key")
		}

		key, err := m.authService.AuthenticateAndRateLimit(r.Context(), authHeader)
		if err != nil {
			m.respondError(w, r, err)
			return
		}

		ctx := context.WithValue(r.Context(), TenantContextKey, key.TenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) respondError(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusUnauthorized
	code := constants.ErrCodeUnauthenticated
	msg := constants.ErrMsgUnauthenticated

	if appErr, ok := err.(*domain.AppException); ok {
		status = appErr.StatusCode
		code = appErr.Code
		msg = appErr.Message
	}

	w.WriteHeader(status)
	resp := domain.NewErrorResponse(code, msg, nil, r.Header.Get("X-Request-ID"))
	_ = json.NewEncoder(w).Encode(resp)
}
