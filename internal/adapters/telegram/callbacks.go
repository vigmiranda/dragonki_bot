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

	msgResponse := MsgResponse{
		Message: "",
		Img:     "assets/maps/error.png",
	}

	switch data {

	case "select_man":
		msg, err := h.service.SetGender(userID, "MAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse(msg)
	case "select_woman":
		msg, err := h.service.SetGender(userID, "WOMAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse(msg)

	case "action_menu":
		msg, err := h.service.Menu(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse(msg)

	case "action_cacar":
		msg, err := h.service.Battle(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao caçar"
			break
		}
		msgResponse = MsgResponse(msg)

	case "action_status":
		msg, err := h.service.Status(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar status"
			break
		}
		msgResponse = MsgResponse(msg)

	case "action_transformar":
		msg, err := h.service.Transform(userID)
		if err != nil {
			msgResponse.Message = "❌ Não foi possível transformar"
			break
		}
		msgResponse = MsgResponse(msg)

	case "action_energia":
		msg, err := h.service.Energy(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		msgResponse = MsgResponse(msg)
	}

	// 🔥 atualiza a tela principal
	h.editPanel(chatID, messageID, msgResponse)

	// remove loading do botão
	h.bot.Request(
		tgbotapi.NewCallback(cb.ID, ""),
	)
}
