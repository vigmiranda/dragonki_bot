package telegram

import (
	"fmt"
	"log"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func (h *Handler) handleInventory(action string, userID int64, params map[string]string) InputTelegram {
	switch action {
	case "open_inventory":
		gameResp, err := h.inventoryService.Open(userID)
		if err != nil {
			fmt.Printf("to_exchange -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")

	case "open_inventory_potions":
		page := 1
		if pageStr, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(pageStr); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.inventoryService.ListItemsByUserAndTypePage(userID, domain.ItemConsumable, page)
		if err != nil {
			fmt.Printf("open_inventory_potions -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")

	case "open_inventory_equipaments":
		gameResp, err := h.inventoryService.OpenEquipmentTypes(userID)
		if err != nil {
			fmt.Printf("open_inventory_equipaments -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")
	case "open_inventory_equipped":
		gameResp, err := h.inventoryService.ViewEquippedItems(userID)
		if err != nil {
			fmt.Printf("open_inventory_equipped -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")
	case "open_inventory_equipaments_type":
		slot := params["slot"]
		page := 1
		if pageStr, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(pageStr); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.inventoryService.ListEquipmentItemsBySlot(userID, slot, page)
		if err != nil {
			fmt.Printf("open_inventory_equipaments_type -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, "to_exchange")
	case "use_item", "equip_item":
		idStr := params["inventoryItemID"]
		id, convErr := strconv.Atoi(idStr)
		if convErr != nil {
			log.Printf("inventory callback parse error | action=%s | userID=%d | rawID=%q | err=%v", action, userID, idStr, convErr)
		}
		log.Printf("inventory callback dispatch | action=%s | userID=%d | inventoryItemID=%d", action, userID, id)
		gameResp, err := h.inventoryService.UseItem(userID, uint(id))
		if err != nil {
			log.Printf("inventory service error | action=%s | userID=%d | inventoryItemID=%d | err=%v", action, userID, id, err)
			gameResp = h.inventoryServiceErrorResponse()
		}
		return ToTelegramResponse(gameResp, true, action)
	case "unequip_item":
		slot := params["slot"]
		log.Printf("inventory callback dispatch | action=%s | userID=%d | slot=%s", action, userID, slot)
		gameResp, err := h.inventoryService.UnequipItem(userID, slot)
		if err != nil {
			log.Printf("inventory service error | action=%s | userID=%d | slot=%s | err=%v", action, userID, slot, err)
			gameResp = h.inventoryServiceErrorResponse()
		}
		return ToTelegramResponse(gameResp, true, action)
	case "unequip_all":
		log.Printf("inventory callback dispatch | action=%s | userID=%d", action, userID)
		gameResp, err := h.inventoryService.UnequipAll(userID)
		if err != nil {
			log.Printf("inventory service error | action=%s | userID=%d | err=%v", action, userID, err)
			gameResp = h.inventoryServiceErrorResponse()
		}
		return ToTelegramResponse(gameResp, true, action)

	}
	return InputTelegram{}
}

func (h *Handler) inventoryServiceErrorResponse() view.GameResponse {
	return view.GameResponse{
		Text:    "Nao foi possivel processar este item agora. Tente novamente.",
		MenuKey: "MAIN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/error.png",
		},
	}
}
