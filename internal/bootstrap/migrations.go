package bootstrap

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const defaultMigrationsDir = "db/migrations"

func ResolveMigrationsDir(raw string) (string, error) {
	dir := strings.TrimSpace(raw)
	if dir == "" {
		dir = defaultMigrationsDir
	}

	if filepath.IsAbs(dir) {
		return dir, nil
	}

	rootDir, err := findProjectRoot()
	if err != nil {
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			return "", err
		}
		return filepath.Join(wd, dir), nil
	}
	return filepath.Join(rootDir, dir), nil
}

func RunMigrationsUp(databaseURL, rawDir string) error {
	m, err := newMigrator(databaseURL, rawDir)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func RunMigrationsDown(databaseURL, rawDir string, steps int) error {
	m, err := newMigrator(databaseURL, rawDir)
	if err != nil {
		return err
	}
	defer m.Close()

	if steps > 0 {
		if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
		return nil
	}

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func MigrationVersion(databaseURL, rawDir string) (uint, bool, error) {
	m, err := newMigrator(databaseURL, rawDir)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return version, dirty, err
}

func ForceMigrationVersion(databaseURL, rawDir string, version int) error {
	m, err := newMigrator(databaseURL, rawDir)
	if err != nil {
		return err
	}
	defer m.Close()

	return m.Force(version)
}

func newMigrator(databaseURL, rawDir string) (*migrate.Migrate, error) {
	dsn := strings.TrimSpace(databaseURL)
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL nao definida")
	}

	migrationsDir, err := ResolveMigrationsDir(rawDir)
	if err != nil {
		return nil, err
	}

	sourceURL := "file://" + filepath.ToSlash(migrationsDir)
	return migrate.New(sourceURL, dsn)
}
