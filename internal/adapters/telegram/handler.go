package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("❌ Panic no handler:", r)
		}
	}()

	if update.CallbackQuery != nil {
		h.handleCallback(update)
		return
	}

	if update.Message == nil {
		return
	}

	if update.Message.Text == "/start" {
		completed, err := h.playerService.HasCompletedOnboarding(update.Message.Chat.ID)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Erro interno"))
			return
		}

		if completed {
			h.newPanel(update.Message.Chat.ID)
		} else {
			h.sendInitialPanel(update.Message.Chat.ID)
		}
		return
	}
	h.handleCallback(update)
}
