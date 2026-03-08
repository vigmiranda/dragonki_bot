package telegram

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func (h *Handler) handleAncestralKiChamber(action string, userID int64, params map[string]string) InputTelegram {
	if !bootstrap.IsFeatureEnabled("ancestral_ki_chamber", false) {
		return InputTelegram{
			Message:  "Camara do Ki Ancestral indisponivel no momento.",
			Img:      tgbotapi.FilePath("assets/error.png"),
			Menu:     mainMenu(),
			ShowMenu: true,
		}
	}
	if h.chamberService == nil {
		return InputTelegram{
			Message:  "Modulo Camara do Ki Ancestral indisponivel no momento.",
			Img:      tgbotapi.FilePath("assets/error.png"),
			Menu:     mainMenu(),
			ShowMenu: true,
		}
	}

	switch action {
	case "to_ancestral_ki_chamber":
		gameResp, err := h.chamberService.Open(userID)
		if err != nil {
			fmt.Printf("to_ancestral_ki_chamber -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "ancestral_ki_levels_6_10":
		return InputTelegram{
			Message:  "Camara do Ki Ancestral\nSelecione o nivel desejado.",
			Img:      tgbotapi.FilePath("assets/npcs/mestre_kalunga.png"),
			Menu:     ancestralKiChamberLevels6To10Menu(),
			ShowMenu: true,
		}
	case "start_ancestral_ki_run":
		level, _ := strconv.Atoi(params["level"])
		gameResp, err := h.chamberService.StartRun(userID, level)
		if err != nil {
			fmt.Printf("start_ancestral_ki_run -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "abandon_ancestral_ki_run":
		gameResp, err := h.chamberService.AbandonRun(userID)
		if err != nil {
			fmt.Printf("abandon_ancestral_ki_run -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "ancestral_ki_chamber_continue":
		gameResp, err := h.chamberService.ContinueRun(userID)
		if err != nil {
			fmt.Printf("ancestral_ki_chamber_continue -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "chamber_open_potions":
		page := 1
		if rawPage, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(rawPage); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.inventoryService.ListItemsByUserAndTypePage(userID, domain.ItemConsumable, page)
		if err != nil {
			fmt.Printf("chamber_open_potions -> %v\\n", err)
			break
		}
		if gameResp.MenuParams == nil {
			gameResp.MenuParams = map[string]any{}
		}
		gameResp.MenuParams["paginationAction"] = "chamber_open_potions"
		gameResp.MenuParams["backAction"] = "ancestral_ki_chamber_continue"
		gameResp.MenuParams["itemAction"] = "chamber_use_item"
		input := ToTelegramResponse(gameResp, true, action)
		if input.Menu == nil {
			input.Menu = ancestralKiChamberPostWinMenu()
		}
		if input.Img == nil {
			input.Img = tgbotapi.FilePath("assets/npcs/arkanis_inventario.png")
		}
		return input
	case "chamber_use_item":
		id, parseErr := strconv.Atoi(params["inventoryItemID"])
		if parseErr != nil {
			break
		}
		gameResp, err := h.inventoryService.UseItem(userID, uint(id))
		if err != nil {
			fmt.Printf("chamber_use_item -> %v\\n", err)
			break
		}
		input := ToTelegramResponse(gameResp, true, action)
		input.Menu = ancestralKiChamberPostWinMenu()
		if input.Img == nil {
			input.Img = tgbotapi.FilePath("assets/npcs/arkanis_inventario.png")
		}
		input.Message = gameResp.Text + "\n\nDeseja continuar ou sair da camara?"
		return input
	case "chamber_transform":
		gameResp, err := h.transformService.Transform(userID)
		if err != nil {
			fmt.Printf("chamber_transform -> %v\\n", err)
			break
		}
		input := ToTelegramResponse(gameResp, true, action)
		input.Menu = ancestralKiChamberPostWinMenu()
		if input.Img == nil {
			input.Img = tgbotapi.FilePath("assets/npcs/mestre_kalunga.png")
		}
		input.Message = gameResp.Text + "\n\nDeseja continuar ou sair da camara?"
		return input
	}
	return InputTelegram{}
}
