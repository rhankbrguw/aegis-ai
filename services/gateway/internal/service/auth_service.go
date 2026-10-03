package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
)

// AuthService handles API key generation, validation, and rate limit checks.
type AuthService struct {
	repo *repository.KeyRepository
}

// NewAuthService constructs an AuthService instance.
func NewAuthService(repo *repository.KeyRepository) *AuthService {
	return &AuthService{repo: repo}
}

// GenerateKey creates a cryptographically secure API key and persists its hash.
func (s *AuthService) GenerateKey(ctx context.Context, req domain.CreateKeyRequest) (*domain.CreateKeyResponse, error) {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("authService.GenerateKey: %w", err)
	}

	secret := fmt.Sprintf("aegis_live_%s", hex.EncodeToString(bytes))
	hash := s.HashKey(secret)
	keyID := fmt.Sprintf("key_%d", time.Now().UnixNano())

	key := domain.APIKey{
		ID:        keyID,
		KeyHash:   hash,
		Prefix:    secret[:15] + "...",
		TenantID:  req.TenantID,
		Name:      req.Name,
		RPM:       req.RPM,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.SaveKey(ctx, key); err != nil {
		return nil, fmt.Errorf("authService.SaveKey: %w", err)
	}

	return &domain.CreateKeyResponse{APIKey: key, Secret: secret}, nil
}

// AuthenticateAndRateLimit validates the bearer token and checks quota.
func (s *AuthService) AuthenticateAndRateLimit(ctx context.Context, authHeader string) (*domain.APIKey, error) {
	rawKey := s.extractToken(authHeader)
	if rawKey == "" {
		return nil, domain.NewAppException(constants.ErrCodeUnauthenticated, constants.ErrMsgUnauthenticated, http.StatusUnauthorized, nil)
	}

	hash := s.HashKey(rawKey)
	key, err := s.repo.GetKeyByHash(ctx, hash)
	if err != nil || key == nil || !key.IsActive {
		return nil, domain.NewAppException(constants.ErrCodeUnauthenticated, constants.ErrMsgUnauthenticated, http.StatusUnauthorized, nil)
	}

	allowed, err := s.repo.CheckAndIncrementRateLimit(ctx, key.ID, key.RPM)
	if err != nil || !allowed {
		return nil, domain.NewAppException(constants.ErrCodeRateLimited, constants.ErrMsgRateLimited, http.StatusTooManyRequests, nil)
	}

	return key, nil
}

// HashKey computes the SHA-256 hash of an API key secret.
func (s *AuthService) HashKey(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) extractToken(header string) string {
	header = strings.TrimSpace(header)
	if strings.HasPrefix(strings.ToLower(header), "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}
