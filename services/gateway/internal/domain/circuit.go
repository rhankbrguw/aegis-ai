package domain

import "time"

// CircuitState represents the current operating mode of the breaker.
type CircuitState string

const (
	StateClosed   CircuitState = "CLOSED"    // Normal operation: routes to primary
	StateOpen     CircuitState = "OPEN"      // Tripped: routes to fallback
	StateHalfOpen CircuitState = "HALF_OPEN" // Trial recovery: probe primary with 1 request
)

// CircuitStatus snapshot for telemetry and mobile dashboard.
type CircuitStatus struct {
	State            CircuitState `json:"state"`
	FailureCount     int64        `json:"failure_count"`
	SuccessCount     int64        `json:"success_count"`
	LastStateChange  time.Time    `json:"last_state_change"`
	IsManualOverride bool         `json:"is_manual_override"`
	KillSwitchActive bool         `json:"kill_switch_active"`
}

// CircuitBreaker defines the resilience interface.
type CircuitBreaker interface {
	GetStatus() CircuitStatus
	CanExecute() bool
	RecordSuccess()
	RecordFailure()
	SetManualOverride(state CircuitState)
	SetKillSwitch(active bool)
}
