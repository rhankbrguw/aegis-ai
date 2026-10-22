package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/samaele/aegis-ai/services/gateway/internal/config"
	"github.com/samaele/aegis-ai/services/gateway/internal/constants"
	"github.com/samaele/aegis-ai/services/gateway/internal/handler"
	"github.com/samaele/aegis-ai/services/gateway/internal/repository"
	"github.com/samaele/aegis-ai/services/gateway/internal/service"
)

func main() {
	cfg := config.Load()

	var redisClient *redis.Client
	if cfg.RedisAddr != "" {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
		})
	}

	cb := service.NewMemoryCircuitBreaker(int64(cfg.CircuitThreshold), cfg.CircuitTimeoutSec)
	fb := service.NewFallbackRouter(cfg.FallbackUpstream, cfg.GeminiKey, 10*time.Second)
	cacheRepo := repository.NewCacheRepository(redisClient)
	telemetryHub := repository.NewTelemetryHub(redisClient)
	finopsSvc := service.NewFinOpsService(redisClient, 50.0)
	guardrailSvc := service.NewGuardrailService()
	proxy := service.NewProxyEngine(cb, fb, cacheRepo, telemetryHub, finopsSvc, guardrailSvc, cfg.PrimaryUpstream, cfg.OpenAIKey, 15*time.Second)

	keyRepo := repository.NewKeyRepository(redisClient)
	authSvc := service.NewAuthService(keyRepo)
	keyHdl := handler.NewKeyHandler(authSvc)

	healthHdl := handler.NewHealthHandler()
	circuitHdl := handler.NewCircuitHandler(cb)
	streamingSvc := service.NewStreamingService(guardrailSvc, telemetryHub)
	chatHdl := handler.NewChatHandler(proxy, streamingSvc)
	telemetryHdl := handler.NewTelemetryHandler(telemetryHub)
	authMiddleware := handler.NewAuthMiddleware(authSvc)
	r := handler.NewRouter(healthHdl, circuitHdl, chatHdl, telemetryHdl, keyHdl, authMiddleware, cfg.AllowedOrigins)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("%s :%s (%s)", constants.MsgServerStarting, cfg.Port, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server startup failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println(constants.MsgServerShutdown)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
