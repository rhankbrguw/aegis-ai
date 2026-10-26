package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// FinOpsService enforces daily spending limits and tracks token expenditures.
type FinOpsService struct {
	rdb            *redis.Client
	maxDailyUSD    float64
	mu             sync.RWMutex
	inMemorySpend  float64
}

// NewFinOpsService creates a thread-safe FinOps budget tracking service.
func NewFinOpsService(rdb *redis.Client, maxDailyUSD float64) *FinOpsService {
	if maxDailyUSD <= 0 {
		maxDailyUSD = 50.0 // Default $50/day spend limit
	}
	return &FinOpsService{
		rdb:         rdb,
		maxDailyUSD: maxDailyUSD,
	}
}

// CheckBudget validates if the current spending is within configured bounds.
func (s *FinOpsService) CheckBudget(ctx context.Context) error {
	currentSpend, err := s.GetCurrentSpend(ctx)
	if err != nil {
		return nil // Fail-open gracefully on Redis read errors
	}
	if currentSpend >= s.maxDailyUSD {
		return domain.NewAppException(
			constants.ErrCodeBudgetExceeded,
			constants.ErrMsgBudgetExceeded,
			http.StatusTooManyRequests,
			nil,
		)
	}
	return nil
}

// RecordSpend accumulates token expenditure for the current day.
func (s *FinOpsService) RecordSpend(ctx context.Context, costUSD float64) error {
	if costUSD <= 0 {
		return nil
	}
	if s.rdb == nil {
		s.mu.Lock()
		s.inMemorySpend += costUSD
		s.mu.Unlock()
		return nil
	}
	key := s.getTodayKey()
	pipe := s.rdb.Pipeline()
	pipe.IncrByFloat(ctx, key, costUSD)
	pipe.Expire(ctx, key, 48*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

// GetCurrentSpend returns the total accumulated USD spend for today.
func (s *FinOpsService) GetCurrentSpend(ctx context.Context) (float64, error) {
	if s.rdb == nil {
		s.mu.RLock()
		defer s.mu.RUnlock()
		return s.inMemorySpend, nil
	}
	key := s.getTodayKey()
	val, err := s.rdb.Get(ctx, key).Float64()
	if err == redis.Nil {
		return 0.0, nil
	}
	return val, err
}

func (s *FinOpsService) getTodayKey() string {
	return fmt.Sprintf("aegis:finops:spend:%s", time.Now().UTC().Format("2006-01-02"))
}
