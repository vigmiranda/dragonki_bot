package postgres

import "database/sql"

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY
	);`)
	if err != nil {
		return err
	}

	var current int
	err = db.QueryRow(`
		SELECT COALESCE(MAX(version), 0)
		FROM schema_migrations;
	`).Scan(&current)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version > current {
			if _, err := db.Exec(m.sql); err != nil {
				return err
			}
			_, err := db.Exec(
				`INSERT INTO schema_migrations(version) VALUES ($1)`,
				m.version,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
