package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type integrationTestDB struct {
	DB          *gorm.DB
	DatabaseURL string
}

func newIntegrationTestDB(t *testing.T) *integrationTestDB {
	t.Helper()

	testDBURL := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if testDBURL == "" {
		t.Skip("set TEST_DATABASE_URL to run postgres integration tests")
	}
	mustUseSafeTestDatabaseURL(t, testDBURL)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	dbName := fmt.Sprintf("dragonki_it_%d", time.Now().UnixNano())
	adminDB, adminURL := openAdminDatabase(t, testDBURL)
	createDatabase(t, ctx, adminDB, dbName)
	t.Cleanup(func() {
		dropDatabase(t, ctx, adminURL, dbName)
	})

	databaseURL := withDatabaseName(t, testDBURL, dbName)
	gormDB, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open gorm db: %v", err)
	}

	if err := autoMigrateIntegrationModels(gormDB); err != nil {
		t.Fatalf("auto migrate integration models: %v", err)
	}

	if err := prepareLegacyColumnsForMigrations(gormDB); err != nil {
		t.Fatalf("prepare legacy schema for migrations: %v", err)
	}

	if err := bootstrap.RunMigrationsUp(databaseURL, ""); err != nil {
		t.Fatalf("run migrations up: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("get sql db from gorm: %v", err)
	}

	t.Cleanup(func() {
		sqlDB.Close()
	})

	return &integrationTestDB{
		DB:          gormDB,
		DatabaseURL: databaseURL,
	}
}

func mustUseSafeTestDatabaseURL(t *testing.T, databaseURL string) {
	t.Helper()

	parsed, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatalf("parse TEST_DATABASE_URL: %v", err)
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		t.Fatal("TEST_DATABASE_URL must include a hostname")
	}

	allowedHosts := map[string]struct{}{
		"localhost":         {},
		"127.0.0.1":         {},
		"::1":               {},
		"postgres":          {},
		"dragonki_postgres": {},
	}

	if _, ok := allowedHosts[host]; ok {
		return
	}

	if strings.HasSuffix(host, ".local") || strings.HasSuffix(host, ".internal") {
		return
	}

	t.Fatalf("refusing to run integration tests against non-local TEST_DATABASE_URL host %q", host)
}

func autoMigrateIntegrationModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Player{},
		&domain.Class{},
		&domain.Monster{},
		&domain.Region{},
		&domain.Payment{},
		&domain.Item{},
		&domain.Equipment{},
		&domain.Inventory{},
		&domain.InventoryItem{},
		&domain.PlayerEquipment{},
		&domain.BattleLog{},
		&domain.MonsterDrop{},
	)
}

func prepareLegacyColumnsForMigrations(db *gorm.DB) error {
	return db.Exec(`
		ALTER TABLE IF EXISTS public.items
			ADD COLUMN IF NOT EXISTS slot VARCHAR(32),
			ADD COLUMN IF NOT EXISTS required_level INTEGER,
			ADD COLUMN IF NOT EXISTS atk_bonus INTEGER,
			ADD COLUMN IF NOT EXISTS def_bonus INTEGER,
			ADD COLUMN IF NOT EXISTS speed_bonus INTEGER,
			ADD COLUMN IF NOT EXISTS hp_bonus INTEGER,
			ADD COLUMN IF NOT EXISTS energy_bonus INTEGER,
			ADD COLUMN IF NOT EXISTS allowed_classes TEXT;

		ALTER TABLE IF EXISTS public.items
			ALTER COLUMN slot DROP NOT NULL,
			ALTER COLUMN required_level DROP NOT NULL,
			ALTER COLUMN atk_bonus DROP NOT NULL,
			ALTER COLUMN def_bonus DROP NOT NULL,
			ALTER COLUMN speed_bonus DROP NOT NULL,
			ALTER COLUMN hp_bonus DROP NOT NULL,
			ALTER COLUMN energy_bonus DROP NOT NULL,
			ALTER COLUMN allowed_classes DROP NOT NULL;
	`).Error
}

func openAdminDatabase(t *testing.T, databaseURL string) (*sql.DB, string) {
	t.Helper()

	adminURL := withDatabaseName(t, databaseURL, "postgres")
	adminDB, err := sql.Open("pgx", adminURL)
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}

	if err := adminDB.Ping(); err != nil {
		adminDB.Close()
		t.Fatalf("ping admin database: %v", err)
	}

	t.Cleanup(func() {
		adminDB.Close()
	})

	return adminDB, adminURL
}

func createDatabase(t *testing.T, ctx context.Context, adminDB *sql.DB, dbName string) {
	t.Helper()

	if _, err := adminDB.ExecContext(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)); err != nil {
		t.Fatalf("create test database %s: %v", dbName, err)
	}
}

func dropDatabase(t *testing.T, ctx context.Context, adminURL, dbName string) {
	t.Helper()

	adminDB, err := sql.Open("pgx", adminURL)
	if err != nil {
		t.Fatalf("reopen admin database for cleanup: %v", err)
	}
	defer adminDB.Close()

	if _, err := adminDB.ExecContext(ctx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, dbName)); err != nil {
		t.Fatalf("drop test database %s: %v", dbName, err)
	}
}

func withDatabaseName(t *testing.T, databaseURL, dbName string) string {
	t.Helper()

	parsed, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatalf("parse database url: %v", err)
	}

	parsed.Path = "/" + dbName
	return parsed.String()
}
