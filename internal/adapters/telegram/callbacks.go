package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
)

func (h *Handler) handleCallback(update tgbotapi.Update) {
	cb := update.CallbackQuery
	if cb == nil || cb.Message == nil {
		return
	}

	data := cb.Data
	chatID := cb.Message.Chat.ID
	userID := cb.From.ID
	messageID := cb.Message.MessageID
	action, params := utils.ParseCallbackData(data)
	inputTelegram := InputTelegram{}
	log.Printf("Callback recebido | action=%s | params=%v", action, params)

	if action == "to_exchange" {
		h.clearPixCodeMessage(userID, chatID)
	}

	switch action {
	case "to_start", "select_man", "select_woman":
		inputTelegram = h.handleOnboarding(action, userID)
	case "action_menu", "menu_adventure", "menu_character", "menu_inventory_scope", "menu_economy", "menu_events":
		inputTelegram = h.handleMenu(action, userID)
	case "list_map", "to_region":
		inputTelegram = h.handleMap(action, userID, params)
	case "action_status", "to_class":
		inputTelegram = h.handlePlayer(action, userID, params)
	case "action_fight", "action_escape", "action_hunt",
		"action_post_battle_open_potions", "action_post_battle_use_item", "action_post_battle_transform", "action_post_battle_menu",
		"action_battle_open_potions", "action_battle_use_item", "action_battle_transform", "action_battle_menu":
		inputTelegram = h.handleBattle(action, userID, params)
	case "open_rank", "rank_top_battles", "rank_top_wins", "rank_top_losses", "rank_top_escape":
		inputTelegram = h.handleRank(action, userID)
	case "list_transform", "action_transform", "to_training":
		inputTelegram = h.handleTraining(action, userID)
	case "to_sleep", "action_sleep", "action_rest":
		inputTelegram = h.handleEnergy(action, userID)
	case "to_shop", "list_potions", "list_equipments", "open_shop_equipment_types", "list_buy_equipments", "list_equipments_type", "list_sell_equipments", "buy_item", "sell_item":
		inputTelegram = h.handleShop(action, userID, params)
	case "to_exchange", "action_convert", "to_vip", "confirm_buy_vip":
		inputTelegram = h.handleExchange(action, userID, params)
	case "to_big_boss", "big_boss_attack", "big_boss_ranking":
		inputTelegram = h.handleBigBoss(action, userID)
	case "to_ancestral_ki_chamber", "ancestral_ki_levels_6_10", "start_ancestral_ki_run", "abandon_ancestral_ki_run", "ancestral_ki_chamber_continue", "chamber_open_potions", "chamber_use_item", "chamber_transform":
		inputTelegram = h.handleAncestralKiChamber(action, userID, params)
	case "open_referral":
		inputTelegram = h.handleReferral(action, userID)
	case "open_inventory", "open_inventory_potions", "open_inventory_equipaments", "open_inventory_equipaments_type", "open_inventory_equipped", "use_item", "equip_item", "unequip_item", "unequip_all":
		inputTelegram = h.handleInventory(action, userID, params)
	case "action_buy_ki", "action_copy_pix_code":
		inputTelegram = h.handlePayment(action, userID)
	}

	if inputTelegram.Message == "" || inputTelegram.Img == nil {
		log.Printf("callback fallback response used | action=%s | params=%v", action, params)
		inputTelegram = InputTelegram{
			Message:  "Nao foi possivel processar esta acao agora.",
			Img:      tgbotapi.FilePath("assets/error.png"),
			Menu:     mainMenu(),
			ShowMenu: true,
		}
	}

	h.editPanel(chatID, messageID, inputTelegram)

	h.bot.Request(tgbotapi.NewCallback(cb.ID, ""))
}
