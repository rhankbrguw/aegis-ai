package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
)

// NewRouter sets up the Chi HTTP router and mounts all endpoints.
func NewRouter(health *HealthHandler, circuit *CircuitHandler, chat *ChatHandler, telemetry *TelemetryHandler, key *KeyHandler, auth *AuthMiddleware, allowedOrigins ...[]string) http.Handler {
	origins := []string{"*"}
	if len(allowedOrigins) > 0 && len(allowedOrigins[0]) > 0 {
		origins = allowedOrigins[0]
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(SecurityHeadersMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Aegis-Key"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get(constants.RouteHealth, health.ServeHTTP)
	r.Get(constants.RouteCircuitState, circuit.GetState)
	r.Post(constants.RouteCircuitOverride, circuit.OverrideState)
	r.Get(constants.RouteMetricsStream, telemetry.StreamUpstream)

	if auth != nil {
		r.With(auth.RequireKey).Post(constants.RouteChatCompletions, chat.Complete)
	} else {
		r.Post(constants.RouteChatCompletions, chat.Complete)
	}

	if key != nil {
		r.Post(constants.RouteAdminKeys, key.CreateKey)
	}

	return r
}
