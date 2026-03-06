package application

import (
	"context"
	"fmt"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PaymentService struct {
	gateway           ports.PaymentGateway
	paymentRepository ports.PaymentRepository
	cfg               appConfig
}

func NewPaymentService(
	gateway ports.PaymentGateway,
	paymentRepository ports.PaymentRepository,
	_ ports.PlayerRepository,
) *PaymentService {
	return &PaymentService{gateway: gateway, paymentRepository: paymentRepository, cfg: getAppConfig()}
}

func (s *PaymentService) CreatePix(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
	payment, qrCode, err := s.gateway.CreatePix(ctx, playerID, amount)
	if err != nil {
		return nil, "", err
	}

	if err := s.paymentRepository.Save(ctx, payment); err != nil {
		return nil, "", err
	}

	return payment, qrCode, nil
}

func (s *PaymentService) HandleBuyKi(userID int64) (view.GameResponse, error) {
	cfg := resolveAppConfig(s.cfg)
	amount := cfg.KiCrystalsValue

	payment, qrCode, err := s.CreatePix(context.Background(), userID, float64(amount))
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

func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID int64) error {
	cfg := resolveAppConfig(s.cfg)
	kiCrystalsMult := cfg.KiCrystalsMult

	settledNow, err := s.paymentRepository.SettlePendingByPartnerID(ctx, paymentID, kiCrystalsMult)
	if err != nil {
		return err
	}
	if !settledNow {
		return nil
	}
	return nil
}
