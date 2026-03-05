package telegram

import "fmt"

func (h *Handler) handleMenu(action string, userID int64) InputTelegram {
	switch action {
	case "action_menu":
		gameResp, err := h.menuService.Menu(userID)
		if err != nil {
			fmt.Printf("action_menu -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_menu")
	}
	return InputTelegram{}
}
