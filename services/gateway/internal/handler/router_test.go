package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/handler"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func setupTestRouter(withAuth bool) (http.Handler, string) {
	cb := service.NewMemoryCircuitBreaker(3, time.Minute)
	fb := service.NewFallbackRouter("https://generativelanguage.googleapis.com", "", time.Second)
	cacheRepo := repository.NewCacheRepository(nil)
	telemetryHub := repository.NewTelemetryHub(nil)
	guardrail := service.NewGuardrailService()
	proxy := service.NewProxyEngine(cb, fb, cacheRepo, telemetryHub, nil, guardrail, "https://api.openai.com/v1", "", time.Second)

	keyRepo := repository.NewKeyRepository(nil)
	authSvc := service.NewAuthService(keyRepo)
	keyHdl := handler.NewKeyHandler(authSvc)

	health := handler.NewHealthHandler()
	circuit := handler.NewCircuitHandler(cb)
	streamingSvc := service.NewStreamingService(guardrail, telemetryHub)
	chat := handler.NewChatHandler(proxy, streamingSvc)
	telemetry := handler.NewTelemetryHandler(telemetryHub)

	var authMiddleware *handler.AuthMiddleware
	var token string
	if withAuth {
		authMiddleware = handler.NewAuthMiddleware(authSvc)
		genKey, _ := authSvc.GenerateKey(nil, domain.CreateKeyRequest{TenantID: "test-tenant", Name: "Key", RPM: 60})
		if genKey != nil {
			token = genKey.Secret
		}
	}

	return handler.NewRouter(health, circuit, chat, telemetry, keyHdl, authMiddleware), token
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupTestRouter(false)
	req := httptest.NewRequest(http.MethodGet, constants.RouteHealth, nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Errorf("expected security header X-Content-Type-Options: nosniff")
	}

	var resp domain.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !resp.Success || resp.Code != "OK" {
		t.Errorf("expected OK success envelope")
	}
}

func TestChatCompletionsAuth(t *testing.T) {
	r, token := setupTestRouter(true)
	payload := domain.ChatCompletionRequest{
		Model:    "gpt-4o",
		Messages: []domain.ChatMessage{{Role: "user", Content: "Hello from test"}},
	}
	body, _ := json.Marshal(payload)

	reqUnauth := httptest.NewRequest(http.MethodPost, constants.RouteChatCompletions, bytes.NewReader(body))
	recUnauth := httptest.NewRecorder()
	r.ServeHTTP(recUnauth, reqUnauth)
	if recUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for unauthenticated request, got %d", recUnauth.Code)
	}

	reqAuth := httptest.NewRequest(http.MethodPost, constants.RouteChatCompletions, bytes.NewReader(body))
	reqAuth.Header.Set("Authorization", "Bearer "+token)
	recAuth := httptest.NewRecorder()
	r.ServeHTTP(recAuth, reqAuth)
	if recAuth.Code != http.StatusOK {
		t.Fatalf("expected status 200 with valid key, got %d", recAuth.Code)
	}
}

func TestChatCompletionsValidation(t *testing.T) {
	r, _ := setupTestRouter(false)
	payload := domain.ChatCompletionRequest{
		Model:    "gpt-4o",
		Messages: []domain.ChatMessage{},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, constants.RouteChatCompletions, bytes.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422 for empty messages, got %d", rec.Code)
	}
}

func TestAdminCreateKeyEndpoint(t *testing.T) {
	r, _ := setupTestRouter(false)
	payload := domain.CreateKeyRequest{
		TenantID: "tenant-beta",
		Name:     "Test Key",
		RPM:      60,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, constants.RouteAdminKeys, bytes.NewReader(body))
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 for key creation, got %d", rec.Code)
	}

	var resp domain.APIResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if !resp.Success || resp.Code != "CREATED" {
		t.Errorf("expected CREATED response envelope")
	}
}
