package domain

// PIIType represents categorized personally identifiable information.
type PIIType string

const (
	PIIEmail      PIIType = "EMAIL"
	PIIPhone      PIIType = "PHONE"
	PIICreditCard PIIType = "CREDIT_CARD"
	PIINIK        PIIType = "ID_NIK"
	PIIApiKey     PIIType = "SECRET_KEY"
)

// PIIMapping stores a masked placeholder and its original sensitive value.
type PIIMapping struct {
	Placeholder string  `json:"placeholder"`
	Original    string  `json:"original"`
	Type        PIIType `json:"type"`
}

// GuardrailResult holds the inspection and sanitization outcome.
type GuardrailResult struct {
	SanitizedMessages []ChatMessage           `json:"sanitized_messages"`
	Redactions        map[string]string       `json:"redactions"`
	IsBlocked         bool                    `json:"is_blocked"`
	BlockReason       string                  `json:"block_reason,omitempty"`
}

// NewCleanGuardrailResult returns an unmodified clean result.
func NewCleanGuardrailResult(msgs []ChatMessage) *GuardrailResult {
	return &GuardrailResult{
		SanitizedMessages: msgs,
		Redactions:        make(map[string]string),
		IsBlocked:         false,
	}
}
