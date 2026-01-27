package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram"
	"github.com/vigmiranda/dragonki/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL não definida")
	}

	// 🔹 Abre conexão com o banco
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// 🔹 Roda migrations
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	// 🔹 Infra DB (wrapper)
	postgresDB := postgres.NewPostgres(db)

	// 🔹 Repositórios
	playerRepo := postgres.NewPlayerPostgresRepo(postgresDB)
	monsterRepo := postgres.NewMonsterPostgresRepo(postgresDB)
	regionRepo := postgres.NewRegionPostgresRepo(postgresDB)
	paymentRepo := postgres.NewPaymentRepo(postgresDB)

	// 🔹 Services (Application)
	playerService := application.NewPlayerService(playerRepo, regionRepo)
	battleService := application.NewBattleService(playerRepo, monsterRepo, regionRepo)

	// 🔹 Mercado Pago (Gateway)
	mpClient := mercadopago.NewClient()
	mpGateway := mercadopago.NewGateway(mpClient)

	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	// 🔹 Bot Telegram
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🐉 DragonKi online:", bot.Self.UserName)

	// 🔹 Handler Telegram
	handler := telegram.NewHandler(
		bot,
		playerService,
		battleService,
		paymentService, // 🔥 NOVO
	)

	// 🔹 Loop principal
	handler.Start()
}
