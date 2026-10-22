package constants

// Application and response string messages.
const (
	MsgServerStarting   = "AegisAI Gateway starting on port"
	MsgServerShutdown   = "AegisAI Gateway shutting down gracefully"
	MsgHealthOK         = "AegisAI Gateway operational"
	MsgChatSuccess      = "Chat completion processed successfully"
	MsgCircuitOpen      = "Circuit breaker OPEN: upstream failover engaged"
	MsgKillSwitchActive = "Emergency kill-switch is active"
	MsgCacheHit         = "Response served from semantic cache"
)
