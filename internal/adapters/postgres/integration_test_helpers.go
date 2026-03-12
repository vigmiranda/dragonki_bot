package postgres

import (
	"testing"

	integrationtest "github.com/vigmiranda/dragonki/internal/testutil/integration"
	"gorm.io/gorm"
)

type integrationTestDB struct {
	DB          *gorm.DB
	DatabaseURL string
}

func newIntegrationTestDB(t *testing.T) *integrationTestDB {
	t.Helper()

	testDB := integrationtest.NewPostgresDB(t, integrationtest.PostgresOptions{
		DatabaseNamePrefix: "dragonki_it",
	})

	return &integrationTestDB{
		DB:          testDB.DB,
		DatabaseURL: testDB.DatabaseURL,
	}
}
