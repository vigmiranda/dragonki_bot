package ports

import "context"
import "github.com/vigmiranda/dragonki/internal/domain"

type PaymentRepository interface {
	Save(ctx context.Context, payment *domain.Payment) error
	SettlePendingByPartnerID(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error)
	GetLatestPendingByPlayerID(ctx context.Context, playerID int64) (*domain.Payment, error)
}
