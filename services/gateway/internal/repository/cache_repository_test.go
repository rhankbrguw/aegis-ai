package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
)

func TestComputeCacheKey(t *testing.T) {
	repo := repository.NewCacheRepository(nil)
	req1 := domain.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []domain.ChatMessage{
			{Role: "user", Content: "Hello world"},
		},
	}
	req2 := domain.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []domain.ChatMessage{
			{Role: "user", Content: "Hello world"},
		},
	}
	req3 := domain.ChatCompletionRequest{
		Model: "gpt-4o",
		Messages: []domain.ChatMessage{
			{Role: "user", Content: "Different prompt"},
		},
	}

	key1 := repo.ComputeCacheKey(req1)
	key2 := repo.ComputeCacheKey(req2)
	key3 := repo.ComputeCacheKey(req3)

	if key1 != key2 {
		t.Errorf("expected deterministic cache keys, got %s vs %s", key1, key2)
	}
	if key1 == key3 {
		t.Errorf("expected different keys for different prompts")
	}
}

func TestNilRedisGraceful(t *testing.T) {
	repo := repository.NewCacheRepository(nil)
	req := domain.ChatCompletionRequest{Model: "gpt-4o"}

	resp, hit, err := repo.Get(context.Background(), req)
	if err != nil || hit || resp != nil {
		t.Errorf("expected nil result on nil client without error")
	}

	err = repo.Set(context.Background(), req, &domain.ChatCompletionResponse{}, time.Hour)
	if err != nil {
		t.Errorf("expected nil error on set with nil client")
	}
}
