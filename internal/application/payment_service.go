package application

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PaymentService struct {
	gateway           ports.PaymentGateway
	paymentRepository ports.PaymentRepository
}

func NewPaymentService(
	gateway ports.PaymentGateway,
	paymentRepository ports.PaymentRepository,
	_ ports.PlayerRepository,
) *PaymentService {
	return &PaymentService{gateway, paymentRepository}
}

func (s *PaymentService) CreatePix(playerID int64, amount float64) (*domain.Payment, string, error) {
	payment, qrCode, err := s.gateway.CreatePix(playerID, amount)
	if err != nil {
		return nil, "", err
	}

	if err := s.paymentRepository.Save(payment); err != nil {
		return nil, "", err
	}

	return payment, qrCode, nil
}

func (s *PaymentService) HandleBuyKi(userID int64) (view.GameResponse, error) {
	amount := 1

	if raw := os.Getenv("KI_CRYSTALS_VALUE"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			amount = parsed
		} else if err != nil {
			return view.GameResponse{}, err
		}
	}

	payment, qrCode, err := s.CreatePix(userID, float64(amount))
	if err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf("%s", payment.QRCode)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type:   view.ImageBase64,
			Base64: qrCode,
			Name:   "pix.png",
		},
		MenuKey: "EXCHANGE_MENU",
	}
	return viewResponse, nil
}

func (s *PaymentService) ConfirmPayment(paymentID int64) error {
	kiCrystalsMult := 1

	if raw := os.Getenv("KI_CRYSTALS_MULT"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		if parsed > 0 {
			kiCrystalsMult = parsed
		}
	}

	settledNow, err := s.paymentRepository.SettlePendingByPartnerID(paymentID, kiCrystalsMult)
	if err != nil {
		return err
	}
	if !settledNow {
		return nil
	}
	return nil
}
