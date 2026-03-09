package webhook

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type paymentConfirmer interface {
	ConfirmPayment(ctx context.Context, paymentID int64) error
}

type paymentWebhookAuditStore interface {
	Save(ctx context.Context, audit *domain.PaymentWebhookAudit) error
}

type MercadoPagoWebhook struct {
	paymentService paymentConfirmer
	secretToken    string
	auditStore     paymentWebhookAuditStore
}

func NewMercadoPagoWebhook(paymentService paymentConfirmer, secretToken string, auditStore ...paymentWebhookAuditStore) *MercadoPagoWebhook {
	var store paymentWebhookAuditStore
	if len(auditStore) > 0 {
		store = auditStore[0]
	}
	return &MercadoPagoWebhook{
		paymentService: paymentService,
		secretToken:    strings.TrimSpace(secretToken),
		auditStore:     store,
	}
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

	audit := &domain.PaymentWebhookAudit{
		ReceivedAt:       time.Now().UTC(),
		RequestMethod:    r.Method,
		RequestPath:      r.URL.Path,
		RemoteAddr:       r.RemoteAddr,
		HeadersJSON:      encodeJSONOrDefault(r.Header, "{}"),
		QueryParamsJSON:  encodeJSONOrDefault(r.URL.Query(), "{}"),
		ProcessingStatus: "RECEIVED",
	}

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		audit.RawBody = ""
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "READ_ERROR", http.StatusBadRequest, readErr.Error())
		wr.WriteHeader(http.StatusBadRequest)
		return
	}
	audit.RawBody = string(body)

	featureEnabled := bootstrap.IsFeatureEnabled("payments", true)
	audit.FeatureEnabled = featureEnabled
	if !featureEnabled {
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "FEATURE_DISABLED", http.StatusNotFound, "")
		wr.WriteHeader(http.StatusNotFound)
		return
	}

	audit.Authorized = w.authorize(r)
	if !audit.Authorized {
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "UNAUTHORIZED", http.StatusUnauthorized, "")
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

	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("webhook decode error:", err)
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "INVALID_JSON", http.StatusBadRequest, err.Error())
		wr.WriteHeader(http.StatusBadRequest)
		return
	}

	audit.EventType = strings.TrimSpace(payload.Type)
	audit.EventAction = strings.TrimSpace(payload.Action)
	audit.EventDataID = strings.TrimSpace(payload.Data.ID)

	if payload.Type != "payment" {
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "IGNORED_NON_PAYMENT", http.StatusOK, "")
		wr.WriteHeader(http.StatusOK)
		return
	}

	if payload.Action == "payment.created" {
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "IGNORED_PAYMENT_CREATED", http.StatusOK, "")
		wr.WriteHeader(http.StatusOK)
		return
	}

	partnerID, err := strconv.ParseInt(payload.Data.ID, 10, 64)
	if err != nil {
		log.Println("invalid payment id in webhook:", payload.Data.ID)
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "INVALID_PAYMENT_ID", http.StatusOK, err.Error())
		wr.WriteHeader(http.StatusOK)
		return
	}
	audit.PartnerPaymentID = partnerID

	if err := w.paymentService.ConfirmPayment(r.Context(), partnerID); err != nil {
		log.Println("error confirming payment:", err)
		finalizeWebhookAudit(r.Context(), w.auditStore, audit, "CONFIRM_ERROR", http.StatusInternalServerError, err.Error())
		wr.WriteHeader(http.StatusInternalServerError)
		return
	}

	finalizeWebhookAudit(r.Context(), w.auditStore, audit, "CONFIRMED", http.StatusOK, "")
	wr.WriteHeader(http.StatusOK)
}

func encodeJSONOrDefault(v any, fallback string) string {
	raw, err := json.Marshal(v)
	if err != nil {
		return fallback
	}
	return string(raw)
}

func finalizeWebhookAudit(
	ctx context.Context,
	store paymentWebhookAuditStore,
	audit *domain.PaymentWebhookAudit,
	status string,
	httpStatus int,
	errMsg string,
) {
	if audit == nil {
		return
	}
	audit.ProcessingStatus = status
	audit.HTTPStatus = httpStatus
	audit.ErrorMessage = errMsg

	if store == nil {
		return
	}
	if err := store.Save(ctx, audit); err != nil {
		log.Printf("webhook audit save error: %v", err)
	}
}
