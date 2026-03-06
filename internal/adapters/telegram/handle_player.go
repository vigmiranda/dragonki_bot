package telegram

import "fmt"

func (h *Handler) handlePlayer(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "action_status":
		gameResp, err := h.playerService.Status(userID)
		if err != nil {
			fmt.Printf("action_status -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_status")
	case "to_class":
		gameResp, err := h.playerService.UpdateClass(userID, params)
		if err != nil {
			fmt.Printf("action_status -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_status")
	}
	return InputTelegram{}
}
