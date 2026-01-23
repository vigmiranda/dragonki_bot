package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram"
	"github.com/vigmiranda/dragonki/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// 🔹 Abre conexão com o banco
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	// 🔹 Roda migrations ANTES de tudo
	if err := postgres.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	// 🔹 Container de infraestrutura SQLite
	sqliteDB := postgres.NewPostgres(db)

	// 🔹 Repositórios
	playerRepo := postgres.NewPlayerPostgresRepo(sqliteDB)
	monsterRepo := postgres.NewMonsterPostgresRepo(sqliteDB)
	regionRepo := postgres.NewRegionPostgresRepo(sqliteDB)

	// 🔹 Services (casos de uso)
	playerService := application.NewPlayerService(playerRepo, regionRepo)
	battleService := application.NewBattleService(playerRepo, monsterRepo)

	// 🔹 Bot Telegram
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🐉 DragonKi online:", bot.Self.UserName)

	// 🔹 Handler (adapter Telegram)
	handler := telegram.NewHandler(bot, playerService, battleService)

	// 🔹 Loop principal
	handler.Start()
}
