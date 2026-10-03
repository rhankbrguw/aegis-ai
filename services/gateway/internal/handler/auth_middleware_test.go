package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/handler"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func TestAuthMiddleware_MissingKey(t *testing.T) {
	repo := repository.NewKeyRepository(nil)
	authSvc := service.NewAuthService(repo)
	mw := handler.NewAuthMiddleware(authSvc)

	nextCalled := false
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	rec := httptest.NewRecorder()

	mw.RequireKey(dummyHandler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized || nextCalled {
		t.Fatalf("expected 401 Unauthorized for missing key, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidKeyAndRateLimit(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewKeyRepository(nil)
	authSvc := service.NewAuthService(repo)
	mw := handler.NewAuthMiddleware(authSvc)

	keyResp, err := authSvc.GenerateKey(ctx, domain.CreateKeyRequest{
		TenantID: "tenant-delta",
		Name:     "TDD Key",
		RPM:      1,
	})
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := r.Context().Value(handler.TenantContextKey)
		if tenantID != "tenant-delta" {
			t.Errorf("expected tenant-delta in context, got %v", tenantID)
		}
		w.WriteHeader(http.StatusOK)
	})

	// 1st request (Allowed)
	req1 := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	req1.Header.Set("X-Aegis-Key", keyResp.Secret)
	rec1 := httptest.NewRecorder()
	mw.RequireKey(dummyHandler).ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Fatalf("expected 200 OK on 1st request, got %d", rec1.Code)
	}

	// 2nd request (Rate limited - RPM: 1)
	req2 := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	req2.Header.Set("Authorization", "Bearer "+keyResp.Secret)
	rec2 := httptest.NewRecorder()
	mw.RequireKey(dummyHandler).ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests on 2nd request, got %d", rec2.Code)
	}
}
