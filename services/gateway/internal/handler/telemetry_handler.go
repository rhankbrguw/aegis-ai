package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow cross-origin WebSocket for mobile app
	},
}

// TelemetryHandler handles WebSocket streaming connections.
type TelemetryHandler struct {
	hub *repository.TelemetryHub
}

// NewTelemetryHandler constructs a TelemetryHandler.
func NewTelemetryHandler(hub *repository.TelemetryHub) *TelemetryHandler {
	return &TelemetryHandler{hub: hub}
}

// StreamUpstream upgrades connection and pipes telemetry to client.
func (h *TelemetryHandler) StreamUpstream(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	h.hub.RegisterClient(conn)
	defer h.hub.UnregisterClient(conn)

	// Keep alive reader
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}
