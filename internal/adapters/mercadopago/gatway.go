package mercadopago

import (
	"fmt"
	"log"
	"time"

	"github.com/mercadopago/sdk-go/pkg/payment"

	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type Gateway struct {
	client *Client
}

var _ ports.PaymentGateway = (*Gateway)(nil)

func NewGateway(client *Client) *Gateway {
	return &Gateway{client: client}
}

func (g *Gateway) CreatePix(
	playerID int64,
	amount float64,
) (*domain.Payment, string, error) {

	externalID := domain.BuildNumericPaymentID(playerID)
	email := fmt.Sprintf("%d_@dragonki.com", playerID)

	req := payment.Request{
		TransactionAmount: amount,
		Description:       "DragonKi - Ki Divino",
		PaymentMethodID:   "pix",
		ExternalReference: externalID,
		Payer: &payment.PayerRequest{
			Email: email,
		},
	}

	resp, err := g.client.Payment.Create(
		g.client.Ctx, // ✅ context obrigatório
		req,
	)
	log.Println(resp)
	if err != nil {
		return nil, "", err
	}

	qrBase64 := resp.PointOfInteraction.TransactionData.QRCodeBase64
	if qrBase64 == "" {
		return nil, "", fmt.Errorf(
			"qrBase64- PIX não gerado. status=%s method=%s",
			resp.Status,
			resp.PaymentMethodID,
		)
	}

	qrCode := resp.PointOfInteraction.TransactionData.QRCode
	if qrCode == "" {
		return nil, "", fmt.Errorf(
			"qrCode - PIX não gerado. status=%s method=%s",
			resp.Status,
			resp.PaymentMethodID,
		)
	}
	paymentDomain := &domain.Payment{
		ID:           externalID,
		PlayerID:     playerID,
		Amount:       amount,
		Status:       domain.PaymentPending,
		QRCode:       qrCode,
		QRCodeBase64: qrBase64,
		PartnerId:    externalID,
		CreatedAt:    time.Now().UTC(),
	}

	return paymentDomain, qrBase64, nil
}
