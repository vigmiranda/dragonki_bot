package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	service *application.PlayerService
}

func NewHandler(bot *tgbotapi.BotAPI, service *application.PlayerService) *Handler {
	return &Handler{
		bot:     bot,
		service: service,
	}
}
