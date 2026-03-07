package postgres

import (
	"errors"
	"strings"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type FeatureFlagPostgresRepo struct {
	db *gorm.DB
}

func NewFeatureFlagPostgresRepo(db *gorm.DB) *FeatureFlagPostgresRepo {
	return &FeatureFlagPostgresRepo{db: db}
}

func (r *FeatureFlagPostgresRepo) FindByName(name string) (bool, bool) {
	var flag domain.FeatureFlag
	err := r.db.
		Model(&domain.FeatureFlag{}).
		Where("name = ?", strings.ToUpper(strings.TrimSpace(name))).
		First(&flag).
		Error

	if err == nil {
		return flag.Enabled, true
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, false
	}

	// If table is unavailable (e.g. migration not applied yet), fallback to default behavior.
	if strings.Contains(strings.ToLower(err.Error()), "does not exist") {
		return false, false
	}

	return false, false
}
