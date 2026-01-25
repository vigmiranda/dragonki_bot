package telegram

import (
	"log"
	"net/url"
	"strings"

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
	action, params := parseCallbackData(data)

	log.Printf("Callback recebido | action=%s | params=%v", action, params)

	msgResponse := MsgResponse{
		Message: "",
		Img:     "assets/error.png",
	}

	switch action {

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
	case "action_fight":
		msg, err := h.battleService.Battle(userID, params)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    mainMenu(),
		}
	case "action_escape":
		menu := mainMenu()
		msg, monster, err := h.battleService.Escape(userID, params)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			msg.Img = "assts/error.png"
			break
		}
		if msg.Conditions {
			menu = battleMenu(monster.MonsterID)
		}

		msgResponse = MsgResponse{
			Message: msg.Message,
			Img:     msg.Img,
			Menu:    menu,
		}

	case "action_cacar":
		menu := mainMenu()
		msg, monster, err := h.battleService.GetBattle(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao caçar"
			break
		}
		if msg.Conditions {
			menu = battleMenu(monster)
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

// ===============================
// CALLBACK PARSER
// ===============================

// parseCallbackData transforma:
// "action_fight?monsterID=12"
// em:
// action = "action_fight"
// params = map[string]string{"monsterID":"12"}
func parseCallbackData(data string) (string, map[string]string) {
	result := make(map[string]string)

	parts := strings.SplitN(data, "?", 2)
	action := parts[0]

	if len(parts) < 2 {
		return action, result
	}

	values, err := url.ParseQuery(parts[1])
	if err != nil {
		log.Println("❌ Erro ao parsear callback:", err)
		return action, result
	}

	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	return action, result
}
