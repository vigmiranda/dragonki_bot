package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range h.bot.GetUpdatesChan(u) {

		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Println("❌ Panic no handler:", r)
				}
			}()

			// callbacks (inline buttons)
			if update.CallbackQuery != nil {
				h.handleCallback(update)
				return
			}

			// mensagens normais
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
					// 🔥 jogador já existe → volta pro menu principal
					h.newPanel(update.Message.Chat.ID)
				} else {
					// 🧑 primeiro acesso → escolher gênero
					h.sendInitialPanel(update.Message.Chat.ID)
				}

				return
			}
		}()
	}
}
