package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/webhook"
)

func main() {
	// Infra
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL não definida")
	}

	// 🔹 Abre conexão com o banco
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// 🔹 Infra DB (wrapper)
	postgresDB := postgres.NewPostgres(db)

	paymentRepo := postgres.NewPaymentRepo(postgresDB)
	playerRepo := postgres.NewPlayerPostgresRepo(postgresDB)

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
