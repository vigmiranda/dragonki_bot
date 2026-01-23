package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) handleCallback(update tgbotapi.Update) {
	cb := update.CallbackQuery
	if cb == nil || cb.Message == nil {
		return
	}

	data := cb.Data
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID
	messageID := cb.Message.MessageID

	var extra string

	switch data {

	case "action_menu":
		msg, err := h.service.Status(userID)
		if err != nil {
			extra = "❌ Erro ao buscar menu"
			break
		}
		extra = msg

	case "action_cacar":
		msg, err := h.service.Battle(userID)
		if err != nil {
			extra = "❌ Erro ao caçar"
			break
		}
		extra = msg

	case "action_status":
		msg, err := h.service.Status(userID)
		if err != nil {
			extra = "❌ Erro ao buscar status"
			break
		}
		extra = msg

	case "action_transformar":
		msg, err := h.service.Transform(userID)
		if err != nil {
			extra = "❌ Não foi possível transformar"
			break
		}
		extra = msg

	case "action_energia":
		msg, err := h.service.Energy(userID)
		if err != nil {
			extra = "❌ Erro ao calcular energia"
			break
		}
		extra = msg
	}

	// 🔥 atualiza a tela principal
	h.editPanel(chatID, messageID, extra)

	// remove loading do botão
	h.bot.Request(
		tgbotapi.NewCallback(cb.ID, ""),
	)
}
