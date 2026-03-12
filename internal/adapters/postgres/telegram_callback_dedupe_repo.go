package postgres

import (
	"context"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelegramCallbackDedupeRepo struct {
	db *gorm.DB
}

func NewTelegramCallbackDedupeRepo(db *gorm.DB) *TelegramCallbackDedupeRepo {
	return &TelegramCallbackDedupeRepo{db: db}
}

func (r *TelegramCallbackDedupeRepo) Reserve(ctx context.Context, dedupe *domain.TelegramCallbackDedupe) (bool, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "dedupe_key"}},
			DoNothing: true,
		}).
		Create(dedupe)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *TelegramCallbackDedupeRepo) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("window_start < ?", before.UTC()).
		Delete(&domain.TelegramCallbackDedupe{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
