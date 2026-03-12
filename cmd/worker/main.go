package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/queue"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"gorm.io/gorm"
)

type workerAppConfig struct {
	DB                   *gorm.DB
	RedisURL             string
	Concurrency          int
	ObservabilityEnabled bool
	PprofEnabled         bool
	ObservabilityAddr    string
}

type workerApp struct {
	Server               *asynq.Server
	Handler              *asynq.ServeMux
	ObservabilityHandler http.Handler
	ObservabilityServer  *http.Server
	Cleanup              func()
}

func envInt(key string, def int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

func envBool(key string, def bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return def
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return def
	}
}

func defaultWorkerAppConfig() workerAppConfig {
	return workerAppConfig{
		RedisURL:             os.Getenv("REDIS_URL"),
		Concurrency:          envInt("QUEUE_CONCURRENCY", 10),
		ObservabilityEnabled: envBool("OBSERVABILITY_ENABLED", true),
		PprofEnabled:         envBool("PPROF_ENABLED", true),
		ObservabilityAddr:    strings.TrimSpace(os.Getenv("WORKER_OBSERVABILITY_ADDR")),
	}
}

func buildWorkerApp(cfg workerAppConfig) (*workerApp, error) {
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL not configured")
	}
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 10
	}
	if cfg.ObservabilityEnabled && cfg.ObservabilityAddr == "" {
		cfg.ObservabilityAddr = ":9091"
	}

	cleanupFns := make([]func(), 0, 1)
	cleanup := func() {
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			cleanupFns[i]()
		}
	}

	db := cfg.DB
	if db == nil {
		db = postgres.OpenDB()
		sqlDB, err := db.DB()
		if err == nil {
			cleanupFns = append(cleanupFns, func() {
				_ = sqlDB.Close()
			})
		}
	}
	postgres.AttachObservability(db)

	weeklyBossRepo := postgres.NewWeeklyBossPostgresRepo(db)
	playerRepo := postgres.NewPlayerPostgresRepo(db)
	inventoryRepo := postgres.NewInventoryPostgresRepo(db)

	service := application.NewWeeklyBossService(weeklyBossRepo, playerRepo, inventoryRepo, nil)
	taskHandler := queue.NewBigBossTaskHandler(service)
	mux := queue.NewServeMux(taskHandler)

	srv, err := queue.NewAsynqServer(cfg.RedisURL, cfg.Concurrency)
	if err != nil {
		cleanup()
		return nil, err
	}

	var observabilityHandler http.Handler
	var observabilityServer *http.Server
	if cfg.ObservabilityEnabled {
		mux := http.NewServeMux()
		observability := bootstrap.NewHTTPObservability("worker")
		observability.RegisterDebugRoutes(mux, cfg.PprofEnabled)
		observabilityHandler = observability.Wrap(mux)
		observabilityServer = &http.Server{
			Addr:              cfg.ObservabilityAddr,
			Handler:           observabilityHandler,
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       30 * time.Second,
		}
		cleanupFns = append(cleanupFns, func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = observabilityServer.Shutdown(ctx)
		})
	}

	return &workerApp{
		Server:               srv,
		Handler:              mux,
		ObservabilityHandler: observabilityHandler,
		ObservabilityServer:  observabilityServer,
		Cleanup:              cleanup,
	}, nil
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	app, err := buildWorkerApp(defaultWorkerAppConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer app.Cleanup()

	log.Printf("worker started | concurrency=%d", defaultWorkerAppConfig().Concurrency)
	if app.ObservabilityServer != nil {
		go func() {
			log.Printf("worker observability listening on %s", app.ObservabilityServer.Addr)
			if err := app.ObservabilityServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("worker observability server error: %v", err)
			}
		}()
	}
	if err := app.Server.Run(app.Handler); err != nil {
		log.Fatal(err)
	}
}
