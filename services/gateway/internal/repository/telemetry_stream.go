package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// TelemetryHub handles Redis Streams and WebSocket client broadcasting.
type TelemetryHub struct {
	redisClient *redis.Client
	mu          sync.RWMutex
	wsClients   map[*websocket.Conn]bool
}

// NewTelemetryHub constructs a TelemetryHub.
func NewTelemetryHub(redisClient *redis.Client) *TelemetryHub {
	return &TelemetryHub{
		redisClient: redisClient,
		wsClients:   make(map[*websocket.Conn]bool),
	}
}

// RegisterClient adds a new WebSocket subscriber.
func (h *TelemetryHub) RegisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.wsClients[conn] = true
}

// UnregisterClient removes a disconnected WebSocket client.
func (h *TelemetryHub) UnregisterClient(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.wsClients, conn)
	_ = conn.Close()
}

// Publish broadcasts event to WebSocket clients and publishes to Redis Streams.
func (h *TelemetryHub) Publish(ctx context.Context, event domain.TelemetryEvent) error {
	h.broadcastWebSocket(event)
	if h.redisClient == nil {
		return nil
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("telemetryHub.Marshal: %w", err)
	}

	err = h.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "aegis:telemetry:stream",
		Values: map[string]interface{}{"payload": string(data)},
	}).Err()
	if err != nil {
		return fmt.Errorf("telemetryHub.XAdd: %w", err)
	}
	return nil
}

func (h *TelemetryHub) broadcastWebSocket(event domain.TelemetryEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.wsClients {
		if err := client.WriteJSON(event); err != nil {
			_ = client.Close()
			delete(h.wsClients, client)
		}
	}
}
