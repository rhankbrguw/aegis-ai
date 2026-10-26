package config_test

import (
	"testing"
	"time"

	"github.com/samaele/aegis-ai/services/gateway/internal/config"
	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
)

func TestLoadDefaults(t *testing.T) {
	cfg := config.Load()

	if cfg.Port != constants.DefaultPort {
		t.Errorf("expected port %s, got %s", constants.DefaultPort, cfg.Port)
	}
	if cfg.CircuitThreshold != constants.DefaultCircuitThreshold {
		t.Errorf("expected threshold %d, got %d", constants.DefaultCircuitThreshold, cfg.CircuitThreshold)
	}
	if cfg.AppURL != constants.DefaultAppURL {
		t.Errorf("expected app url %s, got %s", constants.DefaultAppURL, cfg.AppURL)
	}
	if len(cfg.AllowedOrigins) == 0 {
		t.Errorf("expected allowed origins not empty")
	}
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("expected read timeout 15s, got %v", cfg.ReadTimeout)
	}
}

func TestLoadCustomEnv(t *testing.T) {
	t.Setenv(constants.EnvAppPort, "9090")
	t.Setenv(constants.EnvCircuitThreshold, "10")
	t.Setenv(constants.EnvAllowedOrigins, "https://el-ngadu.vercel.app,https://rhankbrguw.xyz")
	t.Setenv(constants.EnvAppURL, "https://api.rhankbrguw.xyz")

	cfg := config.Load()
	if cfg.Port != "9090" {
		t.Errorf("expected custom port 9090, got %s", cfg.Port)
	}
	if cfg.CircuitThreshold != 10 {
		t.Errorf("expected custom threshold 10, got %d", cfg.CircuitThreshold)
	}
	if cfg.AppURL != "https://api.rhankbrguw.xyz" {
		t.Errorf("expected custom app url, got %s", cfg.AppURL)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Errorf("expected 2 allowed origins, got %d", len(cfg.AllowedOrigins))
	}
}

