package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	SettlePendingByPartnerID(partnerID int64, kiCrystalsMult int) (bool, error)
}
