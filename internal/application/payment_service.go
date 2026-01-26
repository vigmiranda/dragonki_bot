package application

import (
	"errors"

	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PaymentService struct {
	gateway ports.PaymentGateway
	repo    ports.PaymentRepository
	players ports.PlayerRepository
}

func NewPaymentService(
	gateway ports.PaymentGateway,
	repo ports.PaymentRepository,
	players ports.PlayerRepository,
) *PaymentService {
	return &PaymentService{gateway, repo, players}
}

func (s *PaymentService) CreatePix(playerID int64, amount float64) (*domain.Payment, string, error) {
	payment, qrCode, err := s.gateway.CreatePix(playerID, amount)
	if err != nil {
		return nil, "", err
	}

	if err := s.repo.Save(payment); err != nil {
		return nil, "", err
	}

	return payment, qrCode, nil
}

func (s *PaymentService) ConfirmPayment(paymentID string) error {
	payment, err := s.repo.FindByID(paymentID)
	if err != nil {
		return err
	}

	if payment.Status == domain.PaymentApproved {
		return nil
	}

	status, err := s.gateway.GetStatus(paymentID)
	if err != nil {
		return err
	}

	if status != domain.PaymentApproved {
		return errors.New("payment not approved")
	}

	if err := s.repo.UpdateStatus(paymentID, status); err != nil {
		return err
	}

	player, err := s.players.GetOrCreate(payment.PlayerID)
	if err != nil {
		return err
	}

	player.KiCrystals += int(payment.Amount)
	return s.players.Save(player)
}
