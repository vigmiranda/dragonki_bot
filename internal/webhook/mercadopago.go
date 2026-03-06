package webhook

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type paymentConfirmer interface {
	ConfirmPayment(ctx context.Context, paymentID int64) error
}

type MercadoPagoWebhook struct {
	paymentService paymentConfirmer
	secretToken    string
}

func NewMercadoPagoWebhook(paymentService paymentConfirmer, secretToken string) *MercadoPagoWebhook {
	return &MercadoPagoWebhook{paymentService: paymentService, secretToken: strings.TrimSpace(secretToken)}
}

func (w *MercadoPagoWebhook) authorize(r *http.Request) bool {
	if w.secretToken == "" {
		return true
	}

	if r.URL.Query().Get("token") == w.secretToken {
		return true
	}

	if r.Header.Get("X-Webhook-Token") == w.secretToken {
		return true
	}

	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		return token == w.secretToken
	}

	return false
}

func (w *MercadoPagoWebhook) Handle(wr http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if !w.authorize(r) {
		wr.WriteHeader(http.StatusUnauthorized)
		return
	}

	var payload struct {
		Type   string `json:"type"`
		Action string `json:"action"`
		Data   struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("webhook decode error:", err)
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	if payload.Type != "payment" {
		wr.WriteHeader(http.StatusOK)
		return
	}

	if payload.Action == "payment.created" {
		wr.WriteHeader(http.StatusOK)
		return
	}

	partnerID, err := strconv.ParseInt(payload.Data.ID, 10, 64)
	if err != nil {
		log.Println("invalid payment id in webhook:", payload.Data.ID)
		wr.WriteHeader(http.StatusOK)
		return
	}

	if err := w.paymentService.ConfirmPayment(r.Context(), partnerID); err != nil {
		log.Println("error confirming payment:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusOK)
}
