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
		msg, err := h.playerService.SetGender(userID, "MAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}

	case "select_woman":
		msg, err := h.playerService.SetGender(userID, "WOMAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}

	case "action_menu":
		msg, err := h.playerService.Menu(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}

	case "action_cacar":
		menu := battleMenu()
		msg, err := h.battleService.GetBattle(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao caçar"
			break
		}
		if msg.Conditions {
			menu = mainMenu()
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    menu,
		}

	case "action_status":
		msg, err := h.playerService.Status(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar status"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}

	case "action_transformar":
		msg, err := h.playerService.Transform(userID)
		if err != nil {
			msgResponse.Message = "❌ Não foi possível transformar"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}

	case "action_energia":
		msg, err := h.playerService.Energy(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    energyMenu(),
		}
	case "action_sleep":
		menu := energyMenu()
		msg, err := h.playerService.Recovery(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		if msg.Conditions {
			menu = mainMenu()
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    menu,
		}
	}

	// 🔥 atualiza a tela principal
	h.editPanel(chatID, messageID, msgResponse)

	// remove loading do botão
	h.bot.Request(
		tgbotapi.NewCallback(cb.ID, ""),
	)
}
