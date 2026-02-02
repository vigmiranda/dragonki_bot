package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/domain"

	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram"
	"github.com/vigmiranda/dragonki/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	// 🔹 Abre conexão com o banco
	db := postgres.OpenDB()

	// 🔹 Roda migrations
	err := db.AutoMigrate(
		&domain.Player{},
		&domain.Class{},
		&domain.Monster{},
		&domain.Region{},
		&domain.Payment{},
		&domain.Item{},
		&domain.Equipment{},
		&domain.Inventory{},
		&domain.InventoryItem{},
		&domain.PlayerEquipment{},
	)
	if err != nil {
		log.Fatal("❌ error ao rodar automigrate", err)
	}

	// 🔹 Repositórios
	playerRepo := postgres.NewPlayerPostgresRepo(db)
	monsterRepo := postgres.NewMonsterPostgresRepo(db)
	regionRepo := postgres.NewRegionPostgresRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)

	// 🔹 Services (Application)
	playerService := application.NewPlayerService(playerRepo, regionRepo)
	battleService := application.NewBattleService(playerRepo, monsterRepo, regionRepo)
	onboardingService := application.NewOnboardingService(playerRepo)
	transformService := application.NewTransformService(playerRepo)
	menuService := application.NewMenuService(playerRepo, regionRepo)
	energyService := application.NewEnergyService(playerRepo)
	shopService := application.NewShopService(playerRepo)

	// 🔹 Mercado Pago (Gateway)
	mpClient := mercadopago.NewClient()
	mpGateway := mercadopago.NewGateway(mpClient)
	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	// 🔹 Telegram Bot
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🐉 DragonKi bot iniciado:", bot.Self.UserName)

	// 🔹 Handler Telegram
	handler := telegram.NewHandler(
		bot,
		playerService,
		battleService,
		paymentService,
		onboardingService,
		transformService,
		menuService,
		energyService,
		shopService,
	)

	// 🔹 Telegram Webhook HTTP
	http.HandleFunc("/webhook/telegram", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var update tgbotapi.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler.HandleUpdate(update)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("🌐 Telegram Webhook escutando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
