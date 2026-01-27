package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByPartnerID(id int64) (*domain.Payment, error)
	UpdateStatus(id string, status domain.PaymentStatus) error
}
