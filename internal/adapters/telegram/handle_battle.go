package telegram

import "fmt"

func (h *Handler) handleBattle(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "action_fight":
		gameResp, err := h.battleService.Battle(userID, params)
		if err != nil {
			fmt.Printf("action_fight -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_fight")

	case "action_escape":
		gameResp, _, err := h.battleService.Escape(userID, params)
		if err != nil {
			fmt.Printf("action_escape -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_escape")

	case "action_hunt":
		gameResp, err := h.battleService.GetBattle(userID)
		if err != nil {
			fmt.Printf("action_hunt -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "action_hunt")
	}
	return InputTelegram{}
}
