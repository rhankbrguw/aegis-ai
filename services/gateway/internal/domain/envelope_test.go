package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/domain"
)

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"model": "gpt-4o"}
	resp := domain.NewSuccessResponse("OK", constants.MsgChatSuccess, data, "req-123")

	if !resp.Success {
		t.Errorf("expected success true, got false")
	}
	if resp.Code != "OK" {
		t.Errorf("expected code OK, got %s", resp.Code)
	}
	if resp.Meta.RequestID != "req-123" {
		t.Errorf("expected request_id req-123, got %s", resp.Meta.RequestID)
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}
	if len(bytes) == 0 {
		t.Errorf("expected non-empty json")
	}
}

func TestNewErrorResponse(t *testing.T) {
	errs := map[string]string{"prompt": "required"}
	resp := domain.NewErrorResponse(constants.ErrCodeValidation, constants.ErrMsgValidation, errs, "req-456")

	if resp.Success {
		t.Errorf("expected success false, got true")
	}
	if resp.Code != constants.ErrCodeValidation {
		t.Errorf("expected validation error code, got %s", resp.Code)
	}
}
