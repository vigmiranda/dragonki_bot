package application

import (
	"context"
	"errors"
	"fmt"
	"html"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
	"gorm.io/gorm"
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

	msg := fmt.Sprintf(
		"PIX gerado com sucesso!\n\n1) Toque em \"Copiar codigo PIX\".\n2) Cole no app do seu banco e finalize o pagamento.\n\nPagamento ID: <code>%s</code>",
		html.EscapeString(payment.ID),
	)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type:   view.ImageBase64,
			Base64: qrCode,
			Name:   "pix.png",
		},
		MenuKey: "PIX_PAYMENT_MENU",
	}
	return viewResponse, nil
}

func (s *PaymentService) HandleCopyPixCode(userID int64) (view.GameResponse, error) {
	payment, err := s.paymentRepository.GetLatestPendingByPlayerID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return view.GameResponse{
				Text:    "Nao encontrei um PIX pendente para copia.\nGere um novo pagamento no menu de cambio.",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
				MenuKey: "EXCHANGE_MENU",
			}, nil
		}
		return view.GameResponse{}, err
	}

	return view.GameResponse{
		Text: fmt.Sprintf(
			"Codigo PIX (copia e cola):\n\n<code>%s</code>\n\nDepois de pagar, a confirmacao e automatica.",
			html.EscapeString(payment.QRCode),
		),
		Image:   &view.Image{Type: view.ImageAsset, Path: "assets/ki_crystals.png"},
		MenuKey: "PIX_PAYMENT_MENU",
		MenuParams: map[string]any{
			"pixCode": payment.QRCode,
		},
	}, nil
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
