package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

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
	case "action_battle_open_potions":
		monsterID, ok := parseMonsterID(params)
		if !ok {
			break
		}
		page := 1
		if pageStr, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(pageStr); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.inventoryService.ListItemsByUserAndTypePage(userID, domain.ItemConsumable, page)
		if err != nil {
			fmt.Printf("action_battle_open_potions -> %v\\n", err)
			break
		}
		if gameResp.MenuParams == nil {
			gameResp.MenuParams = map[string]any{}
		}
		gameResp.Text = battlePotionListText(gameResp, false)
		gameResp.MenuParams["paginationAction"] = fmt.Sprintf("action_battle_open_potions?monsterID=%d", monsterID)
		gameResp.MenuParams["backAction"] = fmt.Sprintf("action_battle_menu?monsterID=%d", monsterID)
		gameResp.MenuParams["itemAction"] = fmt.Sprintf("action_battle_use_item?monsterID=%d", monsterID)
		return ToTelegramResponse(gameResp, true, action)
	case "action_battle_use_item":
		monsterID, ok := parseMonsterID(params)
		if !ok {
			break
		}
		id, convErr := strconv.Atoi(params["inventoryItemID"])
		if convErr != nil || id <= 0 {
			break
		}
		gameResp, err := h.inventoryService.UseItem(userID, uint(id))
		if err != nil {
			fmt.Printf("action_battle_use_item -> %v\\n", err)
			break
		}
		duelResp, duelErr := h.battleService.OpenCurrentDuel(userID, monsterID)
		if duelErr != nil {
			fmt.Printf("action_battle_use_item(open_duel) -> %v\\n", duelErr)
			break
		}
		duelResp = mergeBattleNotice(gameResp, duelResp)
		return ToTelegramResponse(duelResp, true, action)
	case "action_battle_transform":
		monsterID, ok := parseMonsterID(params)
		if !ok {
			break
		}
		gameResp, err := h.transformService.TransformSingleBattle(userID)
		if err != nil {
			fmt.Printf("action_battle_transform -> %v\\n", err)
			break
		}
		duelResp, duelErr := h.battleService.OpenCurrentDuel(userID, monsterID)
		if duelErr != nil {
			fmt.Printf("action_battle_transform(open_duel) -> %v\\n", duelErr)
			break
		}
		duelResp = mergeBattleNotice(gameResp, duelResp)
		return ToTelegramResponse(duelResp, true, action)
	case "action_battle_menu":
		monsterID, ok := parseMonsterID(params)
		if !ok {
			break
		}
		gameResp, err := h.battleService.OpenCurrentDuel(userID, monsterID)
		if err != nil {
			fmt.Printf("action_battle_menu -> %v\\n", err)
			break
		}
		return ToTelegramResponse(gameResp, true, action)
	case "action_post_battle_open_potions":
		page := 1
		if pageStr, ok := params["page"]; ok {
			if parsed, parseErr := strconv.Atoi(pageStr); parseErr == nil && parsed > 0 {
				page = parsed
			}
		}
		gameResp, err := h.inventoryService.ListItemsByUserAndTypePage(userID, domain.ItemConsumable, page)
		if err != nil {
			fmt.Printf("action_post_battle_open_potions -> %v\\n", err)
			break
		}
		if gameResp.MenuParams == nil {
			gameResp.MenuParams = map[string]any{}
		}
		gameResp.Text = battlePotionListText(gameResp, true)
		gameResp.MenuParams["paginationAction"] = "action_post_battle_open_potions"
		gameResp.MenuParams["backAction"] = "action_post_battle_menu"
		gameResp.MenuParams["itemAction"] = "action_post_battle_use_item"
		return ToTelegramResponse(gameResp, true, action)
	case "action_post_battle_use_item":
		id, convErr := strconv.Atoi(params["inventoryItemID"])
		if convErr != nil || id <= 0 {
			break
		}
		gameResp, err := h.inventoryService.UseItem(userID, uint(id))
		if err != nil {
			fmt.Printf("action_post_battle_use_item -> %v\\n", err)
			break
		}
		menuResp := mergeBattleNotice(gameResp, postBattleMenuResponse())
		return ToTelegramResponse(menuResp, true, action)
	case "action_post_battle_transform":
		gameResp, err := h.transformService.TransformSingleBattle(userID)
		if err != nil {
			fmt.Printf("action_post_battle_transform -> %v\\n", err)
			break
		}
		menuResp := mergeBattleNotice(gameResp, postBattleMenuResponse())
		return ToTelegramResponse(menuResp, true, action)
	case "action_post_battle_menu":
		return ToTelegramResponse(postBattleMenuResponse(), true, action)
	}
	return InputTelegram{}
}

func battleMenuResponse(monsterID int64) view.GameResponse {
	return view.GameResponse{
		Text:    "Duelo em andamento. Escolha sua acao.",
		MenuKey: "BATTLE_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/mestre_kalunga.png",
		},
		MenuParams: map[string]any{
			"monsterID": monsterID,
		},
	}
}

func postBattleMenuResponse() view.GameResponse {
	return view.GameResponse{
		Text:    "Escolha sua proxima acao.",
		MenuKey: "POST_BATTLE_WIN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/mestre_kalunga.png",
		},
	}
}

func parseMonsterID(params map[string]string) (int64, bool) {
	raw := params["monsterID"]
	if raw == "" {
		return 0, false
	}
	monsterID, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || monsterID <= 0 {
		return 0, false
	}
	return monsterID, true
}

func mergeBattleNotice(notice view.GameResponse, duel view.GameResponse) view.GameResponse {
	noticeText := strings.TrimSpace(notice.Text)
	if noticeText == "" {
		return duel
	}
	if strings.TrimSpace(duel.Text) == "" {
		duel.Text = noticeText
		return duel
	}
	duel.Text = fmt.Sprintf("%s\n\n%s", noticeText, duel.Text)
	return duel
}

func battlePotionListText(gameResp view.GameResponse, postBattle bool) string {
	page := menuParamInt(gameResp.MenuParams, "page", 1)
	totalPages := menuParamInt(gameResp.MenuParams, "totalPages", 1)

	if postBattle {
		return fmt.Sprintf(
			"Recupere-se antes da proxima cacada.\nEscolha uma pocao para usar.\n\nPagina %d/%d.",
			page,
			totalPages,
		)
	}

	return fmt.Sprintf(
		"Duelo em andamento.\nEscolha uma pocao para usar e voltar ao combate.\n\nPagina %d/%d.",
		page,
		totalPages,
	)
}

func menuParamInt(params map[string]any, key string, fallback int) int {
	if params == nil {
		return fallback
	}
	raw, ok := params[key]
	if !ok {
		return fallback
	}
	switch v := raw.(type) {
	case int:
		if v > 0 {
			return v
		}
	case int64:
		if v > 0 {
			return int(v)
		}
	case float64:
		if v > 0 {
			return int(v)
		}
	case string:
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			return parsed
		}
	}
	return fallback
}
