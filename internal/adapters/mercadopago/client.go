package mercadopago

import (
	"context"
	"log"
	"os"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
)

type Client struct {
	Payment payment.Client
	Ctx     context.Context
}

func NewClient() *Client {
	accessToken := os.Getenv("MP_ACCESS_TOKEN")
	if accessToken == "" {
		log.Fatal("❌ MP_ACCESS_TOKEN não definido")
	}

	cfg, err := config.New(accessToken)
	if err != nil {
		log.Fatal("❌ erro ao criar config Mercado Pago:", err)
	}

	return &Client{
		Payment: payment.NewClient(cfg), // ✅ valor, não ponteiro
		Ctx:     context.Background(),
	}
}
