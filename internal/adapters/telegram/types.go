package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vigmiranda/dragonki/internal/application"
)

type Handler struct {
	bot           *tgbotapi.BotAPI
	playerService *application.PlayerService
	battleService *application.BattleService
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	playerService *application.PlayerService,
	battleService *application.BattleService,
) *Handler {
	return &Handler{
		bot:           bot,
		playerService: playerService,
		battleService: battleService,
	}
}
