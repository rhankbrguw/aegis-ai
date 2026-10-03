package service

import (
	"sync"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

// MemoryCircuitBreaker implements domain.CircuitBreaker in-memory with thread safety.
type MemoryCircuitBreaker struct {
	mu               sync.RWMutex
	state            domain.CircuitState
	failureCount     int64
	successCount     int64
	failureThreshold int64
	recoveryTimeout  time.Duration
	lastStateChange  time.Time
	isManualOverride bool
	killSwitchActive bool
}

// NewMemoryCircuitBreaker constructs a new circuit breaker instance.
func NewMemoryCircuitBreaker(threshold int64, recoveryTimeout time.Duration) *MemoryCircuitBreaker {
	return &MemoryCircuitBreaker{
		state:            domain.StateClosed,
		failureThreshold: threshold,
		recoveryTimeout:  recoveryTimeout,
		lastStateChange:  time.Now().UTC(),
	}
}

// GetStatus returns the current snapshot of the breaker.
func (cb *MemoryCircuitBreaker) GetStatus() domain.CircuitStatus {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	return domain.CircuitStatus{
		State:            cb.state,
		FailureCount:     cb.failureCount,
		SuccessCount:     cb.successCount,
		LastStateChange:  cb.lastStateChange,
		IsManualOverride: cb.isManualOverride,
		KillSwitchActive: cb.killSwitchActive,
	}
}

// CanExecute checks if primary route can be attempted.
func (cb *MemoryCircuitBreaker) CanExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.killSwitchActive {
		return false
	}
	if cb.isManualOverride {
		return cb.state == domain.StateClosed
	}
	if cb.state == domain.StateOpen {
		if time.Since(cb.lastStateChange) > cb.recoveryTimeout {
			cb.state = domain.StateHalfOpen
			cb.lastStateChange = time.Now().UTC()
			return true
		}
		return false
	}
	return true
}

// RecordSuccess records a healthy upstream response.
func (cb *MemoryCircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.isManualOverride {
		return
	}
	cb.successCount++
	cb.failureCount = 0
	if cb.state == domain.StateHalfOpen {
		cb.state = domain.StateClosed
		cb.lastStateChange = time.Now().UTC()
	}
}

// RecordFailure records an upstream failure or 429 rate limit.
func (cb *MemoryCircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.isManualOverride {
		return
	}
	cb.failureCount++
	if cb.state == domain.StateHalfOpen || cb.failureCount >= cb.failureThreshold {
		cb.state = domain.StateOpen
		cb.lastStateChange = time.Now().UTC()
	}
}

// SetManualOverride enforces a state manually from mobile/dashboard.
func (cb *MemoryCircuitBreaker) SetManualOverride(state domain.CircuitState) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.isManualOverride = true
	cb.state = state
	cb.lastStateChange = time.Now().UTC()
}

// SetKillSwitch activates or deactivates emergency stop.
func (cb *MemoryCircuitBreaker) SetKillSwitch(active bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.killSwitchActive = active
	cb.lastStateChange = time.Now().UTC()
}
