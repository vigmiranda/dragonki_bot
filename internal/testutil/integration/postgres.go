package integrationtest

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
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresDB struct {
	DB          *gorm.DB
	DatabaseURL string
}

type PostgresOptions struct {
	DatabaseNamePrefix string
	MigrationsDir      string
	AutoMigrateModels  []any
	Prepare            func(db *gorm.DB) error
}

func NewPostgresDB(t *testing.T, opts PostgresOptions) *PostgresDB {
	t.Helper()

	testDBURL := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if testDBURL == "" {
		t.Skip("set TEST_DATABASE_URL to run postgres integration tests")
	}
	MustUseSafePostgresURL(t, testDBURL)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	prefix := strings.TrimSpace(opts.DatabaseNamePrefix)
	if prefix == "" {
		prefix = "dragonki_it"
	}
	dbName := fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())

	adminDB, adminURL := openAdminDatabase(t, testDBURL)
	createDatabase(t, ctx, adminDB, dbName)
	t.Cleanup(func() {
		dropDatabase(t, ctx, adminURL, dbName)
	})

	databaseURL := WithDatabaseName(t, testDBURL, dbName)
	gormDB, err := gorm.Open(pgdriver.Open(databaseURL), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open gorm db: %v", err)
	}

	models := opts.AutoMigrateModels
	if len(models) == 0 {
		models = defaultAutoMigrateModels()
	}
	if err := gormDB.AutoMigrate(models...); err != nil {
		t.Fatalf("auto migrate integration models: %v", err)
	}

	prepare := opts.Prepare
	if prepare == nil {
		prepare = PrepareLegacyColumnsForMigrations
	}
	if err := prepare(gormDB); err != nil {
		t.Fatalf("prepare integration schema: %v", err)
	}

	if err := bootstrap.RunMigrationsUp(databaseURL, opts.MigrationsDir); err != nil {
		t.Fatalf("run migrations up: %v", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("get sql db from gorm: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	return &PostgresDB{
		DB:          gormDB,
		DatabaseURL: databaseURL,
	}
}

func MustUseSafePostgresURL(t *testing.T, databaseURL string) {
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

func WithDatabaseName(t *testing.T, databaseURL, dbName string) string {
	t.Helper()

	parsed, err := url.Parse(databaseURL)
	if err != nil {
		t.Fatalf("parse database url: %v", err)
	}

	parsed.Path = "/" + dbName
	return parsed.String()
}

func PrepareLegacyColumnsForMigrations(db *gorm.DB) error {
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

func defaultAutoMigrateModels() []any {
	return []any{
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
	}
}

func openAdminDatabase(t *testing.T, databaseURL string) (*sql.DB, string) {
	t.Helper()

	adminURL := WithDatabaseName(t, databaseURL, "postgres")
	adminDB, err := sql.Open("pgx", adminURL)
	if err != nil {
		t.Fatalf("open admin database: %v", err)
	}

	if err := adminDB.Ping(); err != nil {
		_ = adminDB.Close()
		t.Fatalf("ping admin database: %v", err)
	}

	t.Cleanup(func() {
		_ = adminDB.Close()
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
