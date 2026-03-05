package telegram

import "fmt"

func (h *Handler) handleEnergy(action string, userID int64) InputTelegram {
	switch action {
	case "to_sleep":
		gameResp, err := h.energyService.Energy(userID)
		if err != nil {
			fmt.Printf("to_sleep -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_sleep")
	case "action_sleep":
		gameResp, err := h.energyService.Sleep(userID)
		if err != nil {
			fmt.Printf("action_sleep -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_sleep")
	case "action_rest":
		gameResp, err := h.energyService.Rest(userID)
		if err != nil {
			fmt.Printf("action_rest -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_rest")
	}
	return InputTelegram{}
}
