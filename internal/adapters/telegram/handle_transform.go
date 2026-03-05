package telegram

import "fmt"

func (h *Handler) handleTraining(action string, userID int64) InputTelegram {
	switch action {
	case "to_training":
		gameResp, err := h.transformService.Training()
		if err != nil {
			fmt.Printf("to_training -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")
	case "list_transform":
		gameResp, err := h.transformService.ListTransform(userID)
		if err != nil {
			fmt.Printf("list_transform -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")

	case "action_transform":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			fmt.Printf("action_transform -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_transform")

	}
	return InputTelegram{}
}
