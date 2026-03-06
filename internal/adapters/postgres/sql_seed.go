package postgres

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

func ApplySQLFile(db *gorm.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sqlText := strings.TrimSpace(string(content))
	if sqlText == "" {
		return nil
	}

	if err := db.Exec(sqlText).Error; err != nil {
		return fmt.Errorf("apply sql file %s: %w", path, err)
	}
	return nil
}
