package telegram

import (
	"fmt"
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
		Img:     tgbotapi.FilePath(fmt.Sprintf("assets/error.png")),
	}
	switch action {

	case "select_man":
		_, err := h.onboardingService.SetGender(userID, "MAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao definir gênero"
			break
		}
		msgResponse = MsgResponse{
			Message:  "✍️ Ótima escolha!\n\nAgora digite o nome do seu personagem:",
			Img:      tgbotapi.FilePath(fmt.Sprintf("assets/skins/man_skin_basic.png")),
			Menu:     nil,
			CameFrom: "select_man",
		}
	case "select_woman":
		_, err := h.onboardingService.SetGender(userID, "WOMAN")
		if err != nil {
			msgResponse.Message = "❌ Erro ao definir gênero"
			break
		}

		msgResponse = MsgResponse{
			Message:  "✍️ Ótima escolha!\n\nAgora digite o nome do seu personagem:",
			Img:      tgbotapi.FilePath(fmt.Sprintf("assets/skins/woman_skin_basic.png")),
			Menu:     nil, // ❗ sem menu enquanto digita nome
			CameFrom: "select_woman",
		}
	case "action_menu":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_menu",
		}
	case "action_status":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_status",
		}
	case "action_fight":
		gameResp, err := h.battleService.Battle(userID, params)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "to_sleep",
		}

	case "action_escape":
		gameResp, err, _ := h.battleService.Escape(userID, params)
		if err != nil {
			msgResponse.Message = "❌ Erro ao buscar menu"
			gameResp.Image.Path = fmt.Sprintf("assets/error.png")
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "to_sleep",
		}
	case "action_cacar":
		gameResp, err := h.battleService.GetBattle(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao caçar"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "to_sleep",
		}
	case "action_transformar":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			msgResponse.Message = "❌ Não foi possível transformar"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "to_sleep",
		}
	case "to_sleep":
		gameResp, err := h.energyService.Energy(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "to_sleep",
		}
	case "action_sleep":
		gameResp, err := h.energyService.Sleep(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_rest",
		}
	case "action_rest":
		gameResp, err := h.energyService.Rest(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_rest",
		}
	case "to_shop":
		gameResp, err := h.shopService.OpenShop(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao abrir lojinha"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_buy_ki",
		}
	case "action_buy_ki":
		gameResp, err := h.paymentService.HandleBuyKi(userID)
		if err != nil {
			msgResponse.Message = "❌ Erro ao calcular energia"
			break
		}
		text, photo, menu := ToTelegramResponse(gameResp)
		msgResponse = MsgResponse{
			Message:  text,
			Img:      photo,
			Menu:     menu,
			CameFrom: "action_buy_ki",
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
