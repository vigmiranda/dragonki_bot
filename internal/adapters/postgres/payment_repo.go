package postgres

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(db *gorm.DB) *PaymentRepo {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Save(p *domain.Payment) error {
	err := r.db.Create(&p).Error
	if err != nil {
		return err
	}
	return err
}

func (r *PaymentRepo) FindByPartnerID(id int64) (*domain.Payment, error) {
	var p domain.Payment

	err := r.db.Where("partner_id = ? and status = 'pending'", id).Scan(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, err
}

func (r *PaymentRepo) UpdateStatus(id string, status domain.PaymentStatus) error {
	updates := map[string]interface{}{
		"status":      status,
		"last_update": time.Now(),
	}

	return r.db.
		Model(&domain.Payment{}).
		Where("id = ?", id).
		Updates(updates).
		Error
}
