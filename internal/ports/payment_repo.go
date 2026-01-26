package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	FindByID(id string) (*domain.Payment, error)
	UpdateStatus(id string, status domain.PaymentStatus) error
}
