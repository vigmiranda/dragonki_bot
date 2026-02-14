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

	case "action_menu":
		inputTelegram = h.handleMenu(action, userID)

	case "action_status":
		inputTelegram = h.handlePlayer(action, userID)

	case "action_fight", "action_escape", "action_hunt":
		inputTelegram = h.handleBattle(action, userID, params)

	case "list_transform", "action_transform", "to_training":
		inputTelegram = h.handleTraining(action, userID)

	case "to_sleep", "action_sleep", "action_rest":
		inputTelegram = h.handleEnergy(action, userID)

	case "to_shop", "list_potions", "list_equipments", "buy_item":
		inputTelegram = h.handleShop(action, userID, params)

	case "to_exchange", "action_convert":
		inputTelegram = h.handleExchange(action, userID, params)

	case "open_inventory", "open_inventory_potions", "open_inventory_equipaments", "use_item":
		inputTelegram = h.handleInventory(action, userID, params)

	case "action_buy_ki":
		inputTelegram = h.handlePayment(action, userID)

	}

	// 🔥 atualiza a tela principal
	h.editPanel(chatID, messageID, inputTelegram)

	// remove loading do botão
	h.bot.Request(
		tgbotapi.NewCallback(cb.ID, ""),
	)
}
