package postgres

import (
	"context"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelegramCallbackReceiptRepo struct {
	db *gorm.DB
}

func NewTelegramCallbackReceiptRepo(db *gorm.DB) *TelegramCallbackReceiptRepo {
	return &TelegramCallbackReceiptRepo{db: db}
}

func (r *TelegramCallbackReceiptRepo) Reserve(ctx context.Context, receipt *domain.TelegramCallbackReceipt) (bool, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "callback_id"}},
			DoNothing: true,
		}).
		Create(receipt)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *TelegramCallbackReceiptRepo) DeleteOlderThan(ctx context.Context, before time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("received_at < ?", before.UTC()).
		Delete(&domain.TelegramCallbackReceipt{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
