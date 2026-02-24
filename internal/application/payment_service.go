package application

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
	"github.com/vigmiranda/dragonki/internal/application/view"
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
	paymentRepository ports.PaymentRepository,
	playerRepository ports.PlayerRepository,
) *PaymentService {
	return &PaymentService{gateway, paymentRepository, playerRepository}
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

	amountEnv, err := strconv.ParseInt(os.Getenv("KI_CRYSTALS_VALUE"), 10, 64)
	if err != nil {
		log.Fatal(err)
	} else {
		amount = int(amountEnv)
	}
	log.Println(fmt.Sprintf("Valor do KiDragons: %d Atual: %d", amountEnv, amount))
	payment, qrCode, err := s.CreatePix(userID, float64(amount))
	msg := fmt.Sprintf(
		"%s",
		payment.QRCode,
	)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageBytes,
			Data: utils.DecodeBase64(qrCode),
			Name: "pix.png",
		},
		MenuKey: "EXCHANGE_MENU",
	}
	return viewResponse, nil
}

func (s *PaymentService) ConfirmPayment(paymentID int64) error {
	kiCrystalsMult, err := strconv.Atoi(os.Getenv("KI_CRYSTALS_MULT"))
	if err != nil {
		log.Fatal(err)
	}
	payment, err := s.paymentRepository.FindByPartnerID(paymentID)
	if err != nil {
		return err
	}
	log.Println("kiCrystalsMult -> ", kiCrystalsMult)
	if payment.Status == domain.PaymentApproved {
		return nil
	}

	payment.Status = domain.PaymentApproved
	log.Println("UpdateStatus -- Atualizando ki crystals -> ", paymentID)
	if err := s.paymentRepository.UpdateStatus(payment.ID, payment.Status); err != nil {
		return err
	}

	player, err := s.playersRepository.GetOrCreate(payment.PlayerID)
	if err != nil {
		return err
	}

	player.KiCrystals += int(payment.Amount) * kiCrystalsMult
	log.Println("------ Player name: ", player.Name, "Update(player) -- Atualizando ki crystals => ", int(payment.Amount)*kiCrystalsMult, player.KiCrystals)
	return s.playersRepository.Update(player)
}
