package constants

import "time"

// Environment variable keys.
const (
	EnvAppPort           = "PORT"
	EnvLegacyPort        = "APP_PORT"
	EnvAppEnv            = "APP_ENV"
	EnvLegacyEnv         = "ENV"
	EnvAppURL            = "APP_URL"
	EnvAllowedOrigins    = "ALLOWED_ORIGINS"
	EnvDbConnection      = "DB_CONNECTION"
	EnvDbDatabase        = "DB_DATABASE"
	EnvResendApiKey      = "RESEND_API_KEY"
	EnvResendFrom        = "RESEND_FROM"
	EnvResendFromName    = "RESEND_FROM_NAME"
	EnvRedisAddr         = "REDIS_ADDR"
	EnvRedisPass         = "REDIS_PASSWORD"
	EnvOpenAIKey         = "OPENAI_API_KEY"
	EnvGeminiKey         = "GEMINI_API_KEY"
	EnvPrimaryUpstream   = "PRIMARY_UPSTREAM_URL"
	EnvFallbackUpstream  = "FALLBACK_UPSTREAM_URL"
	EnvCircuitThreshold  = "CIRCUIT_FAILURE_THRESHOLD"
	EnvCircuitTimeoutSec = "CIRCUIT_RECOVERY_TIMEOUT_SEC"
)

// Default configuration fallback values.
const (
	DefaultPort              = "8080"
	DefaultEnv               = "development"
	DefaultAppURL            = "https://api-aegis.rhankbrguw.xyz"
	DefaultAllowedOrigins    = "http://localhost:8080,https://api-aegis.rhankbrguw.xyz,https://rhankbrguw.xyz"
	DefaultDbConnection      = "sqlite"
	DefaultDbDatabase        = "aegis.sqlite"
	DefaultResendFrom        = "no-reply@rhankbrguw.xyz"
	DefaultResendFromName    = "AegisAI Sentinel"
	DefaultRedisAddr         = "localhost:6379"
	DefaultPrimaryUpstream   = "https://api.openai.com/v1"
	DefaultFallbackUpstream  = "https://generativelanguage.googleapis.com/v1beta"
	DefaultCircuitThreshold  = 5
	DefaultCircuitTimeoutSec = 30
	DefaultReadTimeout       = 15 * time.Second
	DefaultWriteTimeout      = 60 * time.Second
	DefaultIdleTimeout       = 120 * time.Second
	DefaultCacheTTL          = 24 * time.Hour
)

