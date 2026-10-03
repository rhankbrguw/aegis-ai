package constants

// HTTP and WebSocket route paths.
const (
	RouteHealth          = "/health"
	RouteChatCompletions = "/v1/chat/completions"
	RouteCircuitState    = "/v1/circuit/state"
	RouteCircuitOverride = "/v1/circuit/override"
	RouteMetricsStream   = "/v1/telemetry/stream"
	RouteAdminKillSwitch = "/v1/admin/kill-switch"
	RouteAdminKeys       = "/v1/admin/keys"
)
