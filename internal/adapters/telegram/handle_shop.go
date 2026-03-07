package telegram

import (
	"fmt"
	"strconv"
)

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
	case "list_buy_equipments":
		gameResp, err := h.shopService.ListEquipmentsByType(params["slot"])
		if err != nil {
			fmt.Printf("list_buy_equipments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_buy_equipments")
	case "open_shop_equipment_types":
		gameResp, err := h.shopService.OpenBuyEquipmentTypes()
		if err != nil {
			fmt.Printf("open_shop_equipment_types -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "open_shop_equipment_types")
	case "list_equipments_type":
		gameResp, err := h.shopService.ListEquipmentsByType(params["slot"])
		if err != nil {
			fmt.Printf("list_equipments_type -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_equipments_type")
	case "list_sell_equipments":
		page := 1
		if pageStr, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(pageStr); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.shopService.ListSellEquipments(userID, page)
		if err != nil {
			fmt.Printf("list_sell_equipments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_sell_equipments")
	case "buy_item":
		gameResp, err := h.shopService.BuyItem(userID, params)
		if err != nil {
			fmt.Printf("list_equipments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "list_equipments")
	case "sell_item":
		gameResp, err := h.shopService.SellItem(userID, params)
		if err != nil {
			fmt.Printf("sell_item -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "sell_item")

	}

	return InputTelegram{}
}
