package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// KeyRepository manages storage and rate limit counters for API keys.
type KeyRepository struct {
	rdb         *redis.Client
	mu          sync.RWMutex
	memoryKeys  map[string]domain.APIKey
	memoryLimit map[string][]time.Time
}

// NewKeyRepository constructs a KeyRepository instance.
func NewKeyRepository(rdb *redis.Client) *KeyRepository {
	return &KeyRepository{
		rdb:         rdb,
		memoryKeys:  make(map[string]domain.APIKey),
		memoryLimit: make(map[string][]time.Time),
	}
}

// SaveKey persists the API key metadata.
func (r *KeyRepository) SaveKey(ctx context.Context, key domain.APIKey) error {
	if r.rdb == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.memoryKeys[key.KeyHash] = key
		return nil
	}

	data, err := json.Marshal(key)
	if err != nil {
		return fmt.Errorf("keyRepo.Marshal: %w", err)
	}

	redisKey := fmt.Sprintf("aegis:keys:%s", key.KeyHash)
	return r.rdb.Set(ctx, redisKey, data, 0).Err()
}

// GetKeyByHash retrieves an API key by its cryptographic hash.
func (r *KeyRepository) GetKeyByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	if r.rdb == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		key, ok := r.memoryKeys[keyHash]
		if !ok {
			return nil, nil
		}
		return &key, nil
	}

	redisKey := fmt.Sprintf("aegis:keys:%s", keyHash)
	data, err := r.rdb.Get(ctx, redisKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("keyRepo.Get: %w", err)
	}

	var key domain.APIKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("keyRepo.Unmarshal: %w", err)
	}
	return &key, nil
}

// CheckAndIncrementRateLimit checks if a key exceeds RPM (Requests Per Minute).
func (r *KeyRepository) CheckAndIncrementRateLimit(ctx context.Context, keyID string, rpm int) (bool, error) {
	if rpm <= 0 {
		return true, nil // Unlimited
	}

	if r.rdb == nil {
		return r.checkMemoryRateLimit(keyID, rpm), nil
	}

	now := time.Now().Unix()
	windowKey := fmt.Sprintf("aegis:ratelimit:%s:%d", keyID, now/60)

	pipe := r.rdb.Pipeline()
	incr := pipe.Incr(ctx, windowKey)
	pipe.Expire(ctx, windowKey, 2*time.Minute)
	if _, err := pipe.Exec(ctx); err != nil {
		return true, nil // Fail-open gracefully
	}

	return incr.Val() <= int64(rpm), nil
}

func (r *KeyRepository) checkMemoryRateLimit(keyID string, rpm int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-time.Minute)

	timestamps := r.memoryLimit[keyID]
	var active []time.Time
	for _, t := range timestamps {
		if t.After(cutoff) {
			active = append(active, t)
		}
	}

	if len(active) >= rpm {
		r.memoryLimit[keyID] = active
		return false
	}

	active = append(active, now)
	r.memoryLimit[keyID] = active
	return true
}
