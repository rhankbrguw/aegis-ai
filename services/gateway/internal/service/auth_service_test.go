package service_test

import (
	"context"
	"testing"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func TestAuthService_GenerateAndAuthenticate(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewKeyRepository(nil)
	authSvc := service.NewAuthService(repo)

	// Generate Key
	req := domain.CreateKeyRequest{TenantID: "tenant-acme", Name: "Production Key", RPM: 2}
	resp, err := authSvc.GenerateKey(ctx, req)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	if resp.Secret == "" || resp.APIKey.TenantID != "tenant-acme" {
		t.Fatalf("invalid generated key response")
	}

	// 1st request (Allowed)
	key, err := authSvc.AuthenticateAndRateLimit(ctx, "Bearer "+resp.Secret)
	if err != nil || key == nil {
		t.Fatalf("expected successful auth, got err: %v", err)
	}

	// 2nd request (Allowed)
	_, err = authSvc.AuthenticateAndRateLimit(ctx, resp.Secret)
	if err != nil {
		t.Fatalf("expected 2nd request allowed, got: %v", err)
	}

	// 3rd request (Rate Limited - RPM: 2)
	_, err = authSvc.AuthenticateAndRateLimit(ctx, resp.Secret)
	if err == nil {
		t.Fatalf("expected rate limit error on 3rd request")
	}
}

func TestAuthService_InvalidToken(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewKeyRepository(nil)
	authSvc := service.NewAuthService(repo)

	_, err := authSvc.AuthenticateAndRateLimit(ctx, "Bearer invalid-token-xyz")
	if err == nil {
		t.Fatalf("expected auth error for invalid token")
	}
}
