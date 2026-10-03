package service

import (
	"context"
	"testing"
)

func TestFinOpsService_InMemoryLimit(t *testing.T) {
	ctx := context.Background()
	svc := NewFinOpsService(nil, 1.0) // $1.00 daily limit

	// Check within limit
	if err := svc.CheckBudget(ctx); err != nil {
		t.Fatalf("expected budget check to pass, got: %v", err)
	}

	// Spend $0.50
	if err := svc.RecordSpend(ctx, 0.50); err != nil {
		t.Fatalf("failed to record spend: %v", err)
	}

	spend, err := svc.GetCurrentSpend(ctx)
	if err != nil || spend != 0.50 {
		t.Fatalf("expected spend 0.50, got %f", spend)
	}

	// Spend $0.60 (total $1.10 -> exceeds $1.00 limit)
	_ = svc.RecordSpend(ctx, 0.60)

	// Check budget should now fail with BUDGET_EXCEEDED
	if err := svc.CheckBudget(ctx); err == nil {
		t.Fatal("expected budget exceeded error, got nil")
	}
}
