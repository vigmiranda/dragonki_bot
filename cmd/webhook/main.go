package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/webhook"
)

func main() {
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
	mpWebhook := webhook.NewMercadoPagoWebhook(paymentService)

	http.HandleFunc("/webhook/mercadopago", mpWebhook.Handle)

	log.Println("🚀 Mercado Pago Webhook rodando")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
