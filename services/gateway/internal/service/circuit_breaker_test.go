package service_test

import (
	"testing"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func TestCircuitBreakerStateTransitions(t *testing.T) {
	timeout := 50 * time.Millisecond
	cb := service.NewMemoryCircuitBreaker(3, timeout)

	if !cb.CanExecute() || cb.GetStatus().State != domain.StateClosed {
		t.Fatalf("expected initial state CLOSED")
	}

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.GetStatus().State != domain.StateClosed {
		t.Fatalf("expected state CLOSED below threshold")
	}

	cb.RecordFailure()
	if cb.GetStatus().State != domain.StateOpen {
		t.Fatalf("expected state OPEN after 3 failures")
	}
	if cb.CanExecute() {
		t.Fatalf("expected CanExecute false when OPEN")
	}

	time.Sleep(60 * time.Millisecond)
	if !cb.CanExecute() {
		t.Fatalf("expected CanExecute true after timeout")
	}
	if cb.GetStatus().State != domain.StateHalfOpen {
		t.Fatalf("expected state HALF_OPEN")
	}

	cb.RecordSuccess()
	if cb.GetStatus().State != domain.StateClosed {
		t.Fatalf("expected recovery to CLOSED on success")
	}
}

func TestKillSwitch(t *testing.T) {
	cb := service.NewMemoryCircuitBreaker(3, time.Second)
	cb.SetKillSwitch(true)

	if cb.CanExecute() {
		t.Errorf("expected CanExecute false when kill switch is active")
	}
	if !cb.GetStatus().KillSwitchActive {
		t.Errorf("expected KillSwitchActive true")
	}

	cb.SetKillSwitch(false)
	if !cb.CanExecute() {
		t.Errorf("expected CanExecute true after releasing kill switch")
	}
}
