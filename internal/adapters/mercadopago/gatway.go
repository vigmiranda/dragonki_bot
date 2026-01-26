package mercadopago

import (
	"time"

	"github.com/google/uuid"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type Gateway struct {
	client *Client
}

func NewGateway(client *Client) *Gateway {
	return &Gateway{client: client}
}

func (g *Gateway) CreatePix(playerID int64, amount float64) (*domain.Payment, string, error) {
	id := uuid.NewString()

	payload := map[string]any{
		"transaction_amount": amount,
		"description":        "DragonKi - Ki Divino",
		"payment_method_id":  "pix",
		"external_reference": id,
		"payer": map[string]string{
			"email": "player@dragonki.com",
		},
	}

	var resp struct {
		PointOfInteraction struct {
			TransactionData struct {
				QRCode string `json:"qr_code"`
			} `json:"transaction_data"`
		} `json:"point_of_interaction"`
	}

	err := g.client.post("https://api.mercadopago.com/v1/payments", payload, &resp)
	if err != nil {
		return nil, "", err
	}

	payment := &domain.Payment{
		ID:        id,
		PlayerID:  playerID,
		Amount:    amount,
		Status:    domain.PaymentPending,
		CreatedAt: time.Now().UTC(),
	}

	return payment, resp.PointOfInteraction.TransactionData.QRCode, nil
}

func (g *Gateway) GetStatus(paymentID string) (domain.PaymentStatus, error) {
	// Simplificado (ideal: GET /v1/payments/{id})
	return domain.PaymentApproved, nil
}
