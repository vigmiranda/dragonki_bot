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

	switch action {
	case "to_start", "select_man", "select_woman":
		inputTelegram = h.handleOnboarding(action, userID)
	case "action_menu", "menu_adventure", "menu_character", "menu_inventory_scope", "menu_economy", "menu_events":
		inputTelegram = h.handleMenu(action, userID)
	case "list_map", "to_region":
		inputTelegram = h.handleMap(action, userID, params)
	case "action_status", "to_class":
		inputTelegram = h.handlePlayer(action, userID, params)
	case "action_fight", "action_escape", "action_hunt":
		inputTelegram = h.handleBattle(action, userID, params)
	case "open_rank", "rank_top_battles", "rank_top_wins", "rank_top_losses", "rank_top_escape":
		inputTelegram = h.handleRank(action, userID)
	case "list_transform", "action_transform", "to_training":
		inputTelegram = h.handleTraining(action, userID)
	case "to_sleep", "action_sleep", "action_rest":
		inputTelegram = h.handleEnergy(action, userID)
	case "to_shop", "list_potions", "list_equipments", "open_shop_equipment_types", "list_buy_equipments", "list_equipments_type", "list_sell_equipments", "buy_item", "sell_item":
		inputTelegram = h.handleShop(action, userID, params)
	case "to_exchange", "action_convert":
		inputTelegram = h.handleExchange(action, userID, params)
	case "to_big_boss", "big_boss_attack", "big_boss_ranking":
		inputTelegram = h.handleBigBoss(action, userID)
	case "open_inventory", "open_inventory_potions", "open_inventory_equipaments", "open_inventory_equipaments_type", "open_inventory_equipped", "use_item", "equip_item", "unequip_item", "unequip_all":
		inputTelegram = h.handleInventory(action, userID, params)
	case "action_buy_ki":
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
