package webhook

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type paymentConfirmer interface {
	ConfirmPayment(paymentID int64) error
}

type MercadoPagoWebhook struct {
	paymentService paymentConfirmer
}

func NewMercadoPagoWebhook(paymentService paymentConfirmer) *MercadoPagoWebhook {
	return &MercadoPagoWebhook{paymentService: paymentService}
}

func (w *MercadoPagoWebhook) Handle(wr http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload struct {
		Type   string `json:"type"`
		Action string `json:"action"`
		Data   struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Println("❌ Webhook decode error:", err)
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	// Mercado Pago envia vários tipos de eventos
	if payload.Type != "payment" {
		wr.WriteHeader(http.StatusOK)
		return
	}

	log.Println("🔔 Webhook MP recebido:", payload.Data.ID)
	if payload.Action == "payment.created" {
		log.Println("Pagamento criado no mercadopago --->", payload.Data.ID)
		return
	}
	partnerId, err := strconv.ParseInt(payload.Data.ID, 10, 64)
	if err != nil {
		log.Println("❌ ID inválido no webhook:", payload.Data.ID)
		return
	}

	if err := w.paymentService.ConfirmPayment(partnerId); err != nil {
		log.Println("❌ Erro ao confirmar pagamento:", err)
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	wr.WriteHeader(http.StatusOK)
}
