package constants

// Application standard error codes.
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeUnauthenticated  = "UNAUTHENTICATED"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeCircuitTripped   = "CIRCUIT_TRIPPED"
	ErrCodeKillSwitchActive = "KILL_SWITCH_ACTIVE"
	ErrCodeBudgetExceeded   = "BUDGET_EXCEEDED"
	ErrCodeRateLimited      = "RATE_LIMITED"
	ErrCodePolicyViolation  = "SECURITY_POLICY_VIOLATION"
	ErrCodeUpstreamError    = "UPSTREAM_ERROR"
	ErrCodeInternalError    = "INTERNAL_ERROR"
)

// Standard error messages.
const (
	ErrMsgValidation       = "Input validation failed."
	ErrMsgUnauthenticated  = "Missing or invalid API authorization."
	ErrMsgCircuitTripped   = "Upstream circuit broken. Request routed to fallback."
	ErrMsgKillSwitchActive = "System operations paused via mobile kill-switch."
	ErrMsgBudgetExceeded   = "Daily FinOps token budget threshold exceeded."
	ErrMsgRateLimited      = "Rate limit exceeded. Please throttle requests."
	ErrMsgPolicyViolation  = "Request violates security policy or prompt safety rules."
	ErrMsgUpstreamError    = "Upstream provider returned an unhandled error."
	ErrMsgInternalError    = "An unexpected internal failure occurred."
)
