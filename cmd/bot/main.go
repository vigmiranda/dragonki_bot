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
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("🐉 DragonKi online:", bot.Self.UserName)

	repo := sqlite.NewPlayerSQLiteRepo(db)
	service := application.NewPlayerService(repo)
	handler := telegram.NewHandler(bot, service)

	handler.Start()
}
