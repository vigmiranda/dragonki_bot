package webhook

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application"
)

type MercadoPagoWebhook struct {
	paymentService *application.PaymentService
}

func NewMercadoPagoWebhook(paymentService *application.PaymentService) *MercadoPagoWebhook {
	return &MercadoPagoWebhook{paymentService: paymentService}
}

func (w *MercadoPagoWebhook) Handle(wr http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload struct {
		Type string `json:"type"`
		Data struct {
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
