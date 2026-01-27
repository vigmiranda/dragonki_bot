package application

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type PaymentService struct {
	gateway           ports.PaymentGateway
	paymentRepository ports.PaymentRepository
	playersRepository ports.PlayerRepository
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

	if err := s.paymentRepository.Save(payment); err != nil {
		return nil, "", err
	}

	return payment, qrCode, nil
}

func (s *PaymentService) ConfirmPayment(paymentID string) error {
	payment, err := s.paymentRepository.FindByPartnerID(paymentID)
	if err != nil {
		return err
	}

	if payment.Status == domain.PaymentApproved {
		return nil
	}

	if err := s.paymentRepository.UpdateStatus(payment.ID, payment.Status); err != nil {
		return err
	}

	player, err := s.playersRepository.GetOrCreate(payment.PlayerID)
	if err != nil {
		return err
	}

	player.KiCrystals += int(payment.Amount)
	return s.playersRepository.Save(player)
}

func (s *PaymentService) HandleBuyKi(chatID int64, userID int64, messageID int) (MsgResponse, error) {

	amount := 1

	amountEnv, err := strconv.ParseInt(os.Getenv("KI_DRAGONS_VALUE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	} else {
		amount = int(amountEnv)
	}
	log.Println(fmt.Sprintf("Valor do KiDragons: %d Atual: %d", amountEnv, amount))

	var msgResponse MsgResponse

	media := tgbotapi.FilePath("assets/error.png")
	msgResponse = MsgResponse{
		Message: fmt.Sprintf(
			"",
		),
		Img:        media,
		Conditions: false,
	}

	payment, qrCode, err := s.CreatePix(userID, float64(amount))
	if err != nil {
		log.Fatal(chatID, "❌ Erro ao gerar pagamento", err)
	}
	log.Println(fmt.Sprintf("Codigo do pagamento: %d", payment))
	msgResponse.Message = fmt.Sprintf(
		"💎 *Comprar Ki Divino*\n\n"+
			"Escaneie o QR Code ou copie o código PIX abaixo.\n\n"+
			"%s\n\n"+
			"<i>Após o pagamento, o Ki será creditado automaticamente.</i>",
		payment.QRCode)

	// Envia QR Code como imagem
	msgResponse.Img = tgbotapi.FileBytes{
		Name:  "pix.png",
		Bytes: utils.DecodeBase64(qrCode),
	}

	return msgResponse, nil
}
