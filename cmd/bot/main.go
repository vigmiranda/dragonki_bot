package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/vigmiranda/dragonki/internal/adapters/sqlite"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram"
	"github.com/vigmiranda/dragonki/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "modernc.org/sqlite"
)

func main() {
	// 🔹 Abre conexão com o banco
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	// 🔹 Roda migrations ANTES de tudo
	if err := sqlite.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	// 🔹 Container de infraestrutura SQLite
	sqliteDB := sqlite.NewSQLite(db)

	// 🔹 Repositórios
	playerRepo := sqlite.NewPlayerSQLiteRepo(sqliteDB)
	monsterRepo := sqlite.NewMonsterSQLiteRepo(sqliteDB)

	// 🔹 Services (casos de uso)
	playerService := application.NewPlayerService(playerRepo)
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
