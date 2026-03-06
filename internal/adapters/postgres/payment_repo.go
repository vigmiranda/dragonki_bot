package postgres

import (
	"context"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Save(ctx context.Context, p *domain.Payment) error {
	return r.db.WithContext(ctx).Create(&p).Error
}

func (r *PaymentRepo) SettlePendingByPartnerID(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error) {
	var p domain.Payment
	settledNow := false

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("partner_id = ?", partnerID).
			First(&p).
			Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return domain.ErrPaymentNotFound
			}
			return err
		}

		if p.Status != domain.PaymentPending {
			return nil
		}

		updates := map[string]any{
			"status":      domain.PaymentApproved,
			"last_update": time.Now().UTC(),
		}

		if err := tx.
			Model(&domain.Payment{}).
			Where("id = ? AND status = ?", p.ID, domain.PaymentPending).
			Updates(updates).
			Error; err != nil {
			return err
		}

		crystalsToAdd := int(p.Amount) * kiCrystalsMult
		result := tx.
			Model(&domain.Player{}).
			Where("user_id = ?", p.PlayerID).
			Update("ki_crystals", gorm.Expr("ki_crystals + ?", crystalsToAdd))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return domain.ErrPaymentNotFound
		}

		p.Status = domain.PaymentApproved
		settledNow = true
		return nil
	})
	if err != nil {
		return false, err
	}

	return settledNow, nil
}
