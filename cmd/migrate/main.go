package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func main() {
	action := flag.String("action", "up", "up | down | version | force")
	steps := flag.Int("steps", 0, "used by down: number of steps")
	forceVersion := flag.Int("version", 0, "used by force: version to force")
	flag.Parse()

	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	databaseURL := os.Getenv("DATABASE_URL")
	migrationsDir := os.Getenv("MIGRATIONS_DIR")

	switch strings.ToLower(strings.TrimSpace(*action)) {
	case "up":
		if err := bootstrap.RunMigrationsUp(databaseURL, migrationsDir); err != nil {
			log.Fatal(err)
		}
		log.Println("migrations up applied")

	case "down":
		if err := bootstrap.RunMigrationsDown(databaseURL, migrationsDir, *steps); err != nil {
			log.Fatal(err)
		}
		log.Printf("migrations down applied | steps=%d", *steps)

	case "version":
		version, dirty, err := bootstrap.MigrationVersion(databaseURL, migrationsDir)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("migration version=%d dirty=%v", version, dirty)

	case "force":
		if err := bootstrap.ForceMigrationVersion(databaseURL, migrationsDir, *forceVersion); err != nil {
			log.Fatal(err)
		}
		log.Printf("migration forced to version=%d", *forceVersion)

	default:
		log.Fatalf("invalid action: %s", *action)
	}
}
