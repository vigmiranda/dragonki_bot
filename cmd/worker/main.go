package main

import (
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/queue"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

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

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL not configured")
	}

	db := postgres.OpenDB()
	weeklyBossRepo := postgres.NewWeeklyBossPostgresRepo(db)
	playerRepo := postgres.NewPlayerPostgresRepo(db)
	inventoryRepo := postgres.NewInventoryPostgresRepo(db)

	service := application.NewWeeklyBossService(weeklyBossRepo, playerRepo, inventoryRepo, nil)
	taskHandler := queue.NewBigBossTaskHandler(service)
	mux := queue.NewServeMux(taskHandler)

	concurrency := envInt("QUEUE_CONCURRENCY", 10)
	srv, err := queue.NewAsynqServer(redisURL, concurrency)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("worker started | concurrency=%d", concurrency)
	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
