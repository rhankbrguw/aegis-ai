package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// CacheRepository manages prompt response caching in Redis.
type CacheRepository struct {
	client *redis.Client
}

// NewCacheRepository constructs a CacheRepository.
func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{client: client}
}

// ComputeCacheKey hashes the model and messages deterministically.
func (r *CacheRepository) ComputeCacheKey(req domain.ChatCompletionRequest) string {
	hasher := sha256.New()
	hasher.Write([]byte(req.Model))
	for _, msg := range req.Messages {
		hasher.Write([]byte(fmt.Sprintf("%s:%s;", msg.Role, msg.Content)))
	}
	return fmt.Sprintf("aegis:cache:%s", hex.EncodeToString(hasher.Sum(nil)))
}

// Get retrieves a cached response if available.
func (r *CacheRepository) Get(ctx context.Context, req domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, bool, error) {
	if r.client == nil {
		return nil, false, nil
	}
	key := r.ComputeCacheKey(req)
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("cacheRepository.Get: %w", err)
	}

	var resp domain.ChatCompletionResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, false, fmt.Errorf("cacheRepository.Unmarshal: %w", err)
	}
	return &resp, true, nil
}

// Set saves a chat response to Redis with expiration.
func (r *CacheRepository) Set(ctx context.Context, req domain.ChatCompletionRequest, resp *domain.ChatCompletionResponse, ttl time.Duration) error {
	if r.client == nil || resp == nil {
		return nil
	}
	key := r.ComputeCacheKey(req)
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("cacheRepository.Marshal: %w", err)
	}
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cacheRepository.Set: %w", err)
	}
	return nil
}
