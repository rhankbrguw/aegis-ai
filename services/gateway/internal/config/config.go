package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
)

// AppConfig holds validated runtime settings.
type AppConfig struct {
	Port               string
	Env                string
	AppURL             string
	AllowedOrigins     []string
	DBConnection       string
	DBDatabase         string
	ResendAPIKey       string
	ResendFrom         string
	ResendFromName     string
	RedisAddr          string
	RedisPassword      string
	OpenAIKey          string
	GeminiKey          string
	PrimaryUpstream    string
	FallbackUpstream   string
	CircuitThreshold   int
	CircuitTimeoutSec  time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleTimeout        time.Duration
}

// Load reads and validates configuration from environment.
func Load() *AppConfig {
	envVal := getFirstEnv([]string{constants.EnvAppEnv, constants.EnvLegacyEnv}, constants.DefaultEnv)
	portVal := getFirstEnv([]string{constants.EnvAppPort, constants.EnvLegacyPort}, constants.DefaultPort)
	originsStr := getEnv(constants.EnvAllowedOrigins, constants.DefaultAllowedOrigins)

	return &AppConfig{
		Port:              portVal,
		Env:               envVal,
		AppURL:            getEnv(constants.EnvAppURL, constants.DefaultAppURL),
		AllowedOrigins:    parseOrigins(originsStr),
		DBConnection:      getEnv(constants.EnvDbConnection, constants.DefaultDbConnection),
		DBDatabase:        getEnv(constants.EnvDbDatabase, constants.DefaultDbDatabase),
		ResendAPIKey:      getEnv(constants.EnvResendApiKey, ""),
		ResendFrom:        getEnv(constants.EnvResendFrom, constants.DefaultResendFrom),
		ResendFromName:    getEnv(constants.EnvResendFromName, constants.DefaultResendFromName),
		RedisAddr:         getEnv(constants.EnvRedisAddr, constants.DefaultRedisAddr),
		RedisPassword:     getEnv(constants.EnvRedisPass, ""),
		OpenAIKey:         getEnv(constants.EnvOpenAIKey, ""),
		GeminiKey:         getEnv(constants.EnvGeminiKey, ""),
		PrimaryUpstream:   getEnv(constants.EnvPrimaryUpstream, constants.DefaultPrimaryUpstream),
		FallbackUpstream:  getEnv(constants.EnvFallbackUpstream, constants.DefaultFallbackUpstream),
		CircuitThreshold:  getEnvAsInt(constants.EnvCircuitThreshold, constants.DefaultCircuitThreshold),
		CircuitTimeoutSec: time.Duration(getEnvAsInt(constants.EnvCircuitTimeoutSec, constants.DefaultCircuitTimeoutSec)) * time.Second,
		ReadTimeout:       constants.DefaultReadTimeout,
		WriteTimeout:      constants.DefaultWriteTimeout,
		IdleTimeout:       constants.DefaultIdleTimeout,
	}
}

// parseOrigins splits a comma-delimited origins string into a slice.
func parseOrigins(raw string) []string {
	if raw == "" || raw == "*" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}

// getFirstEnv returns the first non-empty environment variable from keys.
func getFirstEnv(keys []string, fallback string) string {
	for _, key := range keys {
		if val, ok := os.LookupEnv(key); ok && val != "" {
			return val
		}
	}
	return fallback
}

// getEnv returns the environment variable or fallback value.
func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// getEnvAsInt parses an integer env var with fallback.
func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

