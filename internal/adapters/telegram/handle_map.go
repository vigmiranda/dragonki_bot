package telegram

import "fmt"

func (h *Handler) handleMap(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "list_map":
		gameResp, err := h.mapService.ListMap(userID)
		if err != nil {
			fmt.Printf("list_map -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_map")
	case "to_region":
		gameResp, err := h.mapService.ToRegion(userID, params)
		if err != nil {
			fmt.Printf("list_map -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_map")
	}
	return InputTelegram{}
}
