package telegram

import "fmt"

func (h *Handler) handleShop(action string, userID int64) InputTelegram {
	switch action {
	case "to_shop":
		gameResp, err := h.shopService.OpenShop(userID)
		if err != nil {
			fmt.Printf("to_shop -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_shop")
	}

	return InputTelegram{}
}
