package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/webhook"
	"gorm.io/gorm"
)

type webhookAppConfig struct {
	DB                  *gorm.DB
	MPClient            *mercadopago.Client
	FeatureFlagResolver func(name string) (enabled bool, found bool)
}

type webhookApp struct {
	Handler http.Handler
	Cleanup func()
}

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

func parseBoolEnv(name string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		log.Printf("invalid bool for %s=%q, using fallback %t", name, raw, fallback)
		return fallback
	}
}

func buildWebhookApp(cfg webhookAppConfig) (*webhookApp, error) {
	cleanupFns := make([]func(), 0, 2)
	cleanup := func() {
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			cleanupFns[i]()
		}
	}

	db := cfg.DB
	if db == nil {
		db = postgres.OpenDB()
		sqlDB, err := db.DB()
		if err == nil {
			cleanupFns = append(cleanupFns, func() {
				_ = sqlDB.Close()
			})
		}
	}
	postgres.AttachObservability(db)

	paymentRepo := postgres.NewPaymentRepo(db)
	webhookAuditRepo := postgres.NewPaymentWebhookAuditRepo(db)
	playerRepo := postgres.NewPlayerPostgresRepo(db)
	featureFlagRepo := postgres.NewFeatureFlagPostgresRepo(db)

	resolver := cfg.FeatureFlagResolver
	if resolver == nil {
		resolver = featureFlagRepo.FindByName
	}
	restoreFeatureFlags := bootstrap.WithFeatureFlagResolver(resolver)
	cleanupFns = append(cleanupFns, restoreFeatureFlags)

	mpClient := cfg.MPClient
	if mpClient == nil {
		mpClient = mercadopago.NewClient()
	}
	mpGateway := mercadopago.NewGateway(mpClient)

	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	mpWebhook := webhook.NewMercadoPagoWebhook(paymentService, os.Getenv("MP_WEBHOOK_TOKEN"), webhookAuditRepo)

	mux := http.NewServeMux()
	registerMercadoPagoWebhook(mux, mpWebhook.Handle)

	observabilityEnabled := parseBoolEnv("OBSERVABILITY_ENABLED", true)
	if observabilityEnabled {
		observability := bootstrap.NewHTTPObservability("webhook")
		observability.RegisterDebugRoutes(mux, parseBoolEnv("PPROF_ENABLED", true))
		return &webhookApp{
			Handler: observability.Wrap(mux),
			Cleanup: cleanup,
		}, nil
	}

	return &webhookApp{
		Handler: mux,
		Cleanup: cleanup,
	}, nil
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	app, err := buildWebhookApp(webhookAppConfig{})
	if err != nil {
		log.Fatal(err)
	}
	defer app.Cleanup()

	log.Println("Mercado Pago Webhook rodando")
	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.Handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
