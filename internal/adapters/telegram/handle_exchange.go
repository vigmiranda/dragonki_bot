package telegram

import "fmt"

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
	}
	return InputTelegram{}
}
