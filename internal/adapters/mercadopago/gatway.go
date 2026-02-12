package mercadopago

import (
	"fmt"
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
	payer := &payment.PayerRequest{
		Email: email,
	}

	req := payment.Request{
		TransactionAmount:   amount,
		Description:         "DragonKi - Ki Divino",
		PaymentMethodID:     "pix",
		ExternalReference:   externalID,
		StatementDescriptor: "DragonKi",
		AdditionalInfo: &payment.AdditionalInfoRequest{
			Items: []payment.ItemRequest{
				{
					EventDate:   fmt.Sprintf("%d", time.Now().UTC()),
					ID:          fmt.Sprintf("ki_%d", amount),
					Title:       fmt.Sprintf("%d Ki Crystals", amount),
					Description: "Moeda premium DragonKi",
					CategoryID:  "digital_goods",
					UnitPrice:   float64(amount),
					Quantity:    1,
				},
			},
		},
		Payer: payer,
	}

	resp, err := g.client.Payment.Create(
		g.client.Ctx, // ✅ context obrigatório
		req,
	)
	if err != nil {
		return nil, "", err
	}

	qrBase64 := resp.PointOfInteraction.TransactionData.QRCodeBase64
	if qrBase64 == "" {
		return nil, "", fmt.Errorf(
			"qrBase64- PIX não gerado. qrBase64=%s",
			resp.PointOfInteraction.TransactionData.QRCodeBase64,
		)
	}

	qrCode := resp.PointOfInteraction.TransactionData.QRCode
	if qrCode == "" {
		return nil, "", fmt.Errorf(
			"qrCode - PIX não gerado. QRCode=%s",
			resp.PointOfInteraction.TransactionData.QRCode,
		)
	}
	ID := resp.ID
	if ID == 0 {
		return nil, "", fmt.Errorf(
			"ID - PIX não gerado. ID=%d",
			resp.ID,
		)
	}

	paymentDomain := &domain.Payment{
		ID:           externalID,
		PlayerID:     playerID,
		Amount:       amount,
		Status:       domain.PaymentPending,
		QRCode:       qrCode,
		QRCodeBase64: qrBase64,
		PartnerId:    int64(ID),
		CreatedAt:    time.Now().UTC(),
	}

	return paymentDomain, qrBase64, nil
}
