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
	inputTelegram := InputTelegram{}
	log.Printf("Callback recebido | action=%s | params=%v", action, params)

	switch action {

	case "select_man":
		gameResp, err := h.onboardingService.SetGender(userID, "MAN")
		if err != nil {
			fmt.Printf("select_man -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, false, "select_woman")

	case "select_woman":
		gameResp, err := h.onboardingService.SetGender(userID, "WOMAN")
		if err != nil {
			fmt.Printf("select_woman -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, false, "select_woman")

	case "action_menu":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			fmt.Printf("action_menu -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_menu")

	case "action_status":
		gameResp, err := h.playerService.Status(userID)
		if err != nil {
			fmt.Printf("action_status -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_status")

	case "action_fight":
		gameResp, err := h.battleService.Battle(userID, params)
		if err != nil {
			fmt.Printf("action_fight -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_fight")

	case "action_escape":
		gameResp, err, _ := h.battleService.Escape(userID, params)
		if err != nil {
			fmt.Printf("action_escape -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_escape")

	case "action_cacar":
		gameResp, err := h.battleService.GetBattle(userID)
		if err != nil {
			fmt.Printf("action_cacar -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_cacar")

	case "action_transformar":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			fmt.Printf("action_transformar -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_transformar")

	case "to_sleep":
		gameResp, err := h.energyService.Energy(userID)
		if err != nil {
			fmt.Printf("to_sleep -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "to_sleep")

	case "action_sleep":
		gameResp, err := h.energyService.Sleep(userID)
		if err != nil {
			fmt.Printf("action_sleep -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_sleep")

	case "action_rest":
		gameResp, err := h.energyService.Rest(userID)
		if err != nil {
			fmt.Printf("action_rest -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_rest")

	case "to_shop":
		gameResp, err := h.shopService.OpenShop(userID)
		if err != nil {
			fmt.Printf("to_shop -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "to_shop")

	case "action_buy_ki":
		gameResp, err := h.paymentService.HandleBuyKi(userID)
		if err != nil {
			fmt.Printf("action_buy_ki -> ", err)
			break
		}
		inputTelegram = ToTelegramResponse(gameResp, true, "action_buy_ki")
	}

	// 🔥 atualiza a tela principal
	h.editPanel(chatID, messageID, inputTelegram)

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
