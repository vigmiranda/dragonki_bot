package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func seedFilePath() string {
	if custom := strings.TrimSpace(os.Getenv("SQL_SEED_FILE")); custom != "" {
		return custom
	}
	return filepath.Join("internal", "adapters", "postgres", "SQL", "Inserts", "monster_drops.sql")
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	db := postgres.OpenDB()
	path := seedFilePath()
	if err := postgres.ApplySQLFile(db, path); err != nil {
		log.Fatal(err)
	}

	log.Printf("sql seed applied successfully: %s", path)
}
