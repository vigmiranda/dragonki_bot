package telegram

import "fmt"
import "github.com/vigmiranda/dragonki/internal/bootstrap"
import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (h *Handler) handleExchange(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "to_exchange":
		gameResp, err := h.exchangeService.Open()
		if err != nil {
			fmt.Printf("to_exchange -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")
	case "action_convert":
		gameResp, err := h.exchangeService.Convert(userID, params)
		if err != nil {
			fmt.Printf("action_convert -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_convert")
	case "to_vip":
		if !bootstrap.IsFeatureEnabled("vip", false) || h.vipService == nil {
			return InputTelegram{
				Message:  "VIP indisponivel no momento.",
				Img:      tgbotapi.FilePath("assets/error.png"),
				Menu:     economyMenu(),
				ShowMenu: true,
			}
		}
		gameResp, err := h.vipService.Open(userID)
		if err != nil {
			fmt.Printf("to_vip -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_vip")
	case "confirm_buy_vip":
		if !bootstrap.IsFeatureEnabled("vip", false) || h.vipService == nil {
			return InputTelegram{
				Message:  "VIP indisponivel no momento.",
				Img:      tgbotapi.FilePath("assets/error.png"),
				Menu:     economyMenu(),
				ShowMenu: true,
			}
		}
		gameResp, err := h.vipService.Purchase(userID)
		if err != nil {
			fmt.Printf("confirm_buy_vip -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "confirm_buy_vip")
	}
	return InputTelegram{}
}
