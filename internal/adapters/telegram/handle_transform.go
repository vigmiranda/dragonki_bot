package telegram

import "fmt"

func (h *Handler) handleTransform(action string, userID int64) InputTelegram {
	switch action {
	case "action_transformar":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			fmt.Printf("action_transformar -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transformar")

	}
	return InputTelegram{}
}
