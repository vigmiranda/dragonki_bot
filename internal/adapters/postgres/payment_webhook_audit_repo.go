package postgres

import (
	"context"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type PaymentWebhookAuditRepo struct {
	db *gorm.DB
}

func NewPaymentWebhookAuditRepo(db *gorm.DB) *PaymentWebhookAuditRepo {
	return &PaymentWebhookAuditRepo{db: db}
}

func (r *PaymentWebhookAuditRepo) Save(ctx context.Context, audit *domain.PaymentWebhookAudit) error {
	if audit == nil {
		return gorm.ErrInvalidData
	}
	return r.db.WithContext(ctx).Create(audit).Error
}
