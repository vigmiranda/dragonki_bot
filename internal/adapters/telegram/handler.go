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

	status, err := h.onboardingService.HasCompletedOnboarding(update.Message.Chat.ID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Erro interno"))
		return
	}
	if update.Message.Text == "/start" {
		switch status {
		case "DONE":
			h.newPanel(update.Message.Chat.ID, "OLD_PLAYER")
		default:
			h.sendInitialPanel(update.Message.Chat.ID)
		}

		return
	}
	if status == "CHOOSING_NAME" {
		_, err := h.handleNameInput(update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			h.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Erro interno"))
			return
		}
		h.newPanel(update.Message.Chat.ID, "NEW_PLAYER")
	}

	h.handleCallback(update)
}
