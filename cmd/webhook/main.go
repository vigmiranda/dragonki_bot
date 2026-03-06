package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/webhook"
)

func registerMercadoPagoWebhook(mux *http.ServeMux, handleFunc http.HandlerFunc) {
	mux.HandleFunc("/webhook/mercadopago", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		handleFunc(w, r)
	})
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	// 🔹 Abre conexão com o banco
	db := postgres.OpenDB()

	paymentRepo := postgres.NewPaymentRepo(db)
	playerRepo := postgres.NewPlayerPostgresRepo(db)

	mpClient := mercadopago.NewClient()
	mpGateway := mercadopago.NewGateway(mpClient)

	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	// Webhook
	mpWebhook := webhook.NewMercadoPagoWebhook(paymentService, os.Getenv("MP_WEBHOOK_TOKEN"))

	mux := http.NewServeMux()
	registerMercadoPagoWebhook(mux, mpWebhook.Handle)

	log.Println("🚀 Mercado Pago Webhook rodando")
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
