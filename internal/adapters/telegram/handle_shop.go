package telegram

import "fmt"

func (h *Handler) handleShop(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "to_shop":
		gameResp, err := h.shopService.OpenShop(userID)
		if err != nil {
			fmt.Printf("to_shop -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_shop")
	case "list_potions":
		gameResp, err := h.shopService.ListPotions()
		if err != nil {
			fmt.Printf("list_potions -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_potions")
	case "list_equipments":
		gameResp, err := h.shopService.ListEquipments()
		if err != nil {
			fmt.Printf("list_equipments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_equipments")
	case "buy_item":
		gameResp, err := h.shopService.BuyItem(userID, params)
		if err != nil {
			fmt.Printf("list_equipments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_equipments")

	}

	return InputTelegram{}
}
