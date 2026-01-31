package telegram

import "fmt"

func (h *Handler) handlePlayer(action string, userID int64) InputTelegram {
	switch action {
	case "action_status":
		gameResp, err := h.playerService.Status(userID)
		if err != nil {
			fmt.Printf("action_status -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_status")
	}
	return InputTelegram{}
}
