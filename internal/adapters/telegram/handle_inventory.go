package telegram

import (
	"fmt"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func (h *Handler) handleInventory(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "open_inventory":
		gameResp, err := h.inventoryService.Open(userID)
		if err != nil {
			fmt.Printf("to_exchange -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")

	case "open_inventory_potions":
		gameResp, err := h.inventoryService.ListItemsByUserAndType(userID, domain.ItemConsumable)
		if err != nil {
			fmt.Printf("open_inventory_potions -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")

	case "open_inventory_equipaments":
		gameResp, err := h.inventoryService.ListItemsByUserAndType(userID, domain.ItemEquipment)
		if err != nil {
			fmt.Printf("open_inventory_equipaments -> ", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")
	case "use_item":
		idStr := params["inventoryItemID"]
		id, _ := strconv.Atoi(idStr)
		gameResp, err := h.inventoryService.UseItem(userID, uint(id))
		if err != nil {
			break
		}
		return ToTelegramResponse(gameResp, true, "use_item")

	}
	return InputTelegram{}
}
