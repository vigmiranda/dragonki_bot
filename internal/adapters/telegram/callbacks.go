package telegram

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
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
	now := time.Now().UTC()
	bootstrap.IncAppMetric("telegram.callback.received_total")

	if h.callbackReceipts != nil && shouldGuardTelegramCallback(action) && cb.ID != "" {
		reserved, err := h.callbackReceipts.Reserve(context.Background(), &domain.TelegramCallbackReceipt{
			CallbackID: cb.ID,
			UserID:     userID,
			ChatID:     chatID,
			MessageID:  messageID,
			Action:     action,
			RawData:    data,
			ReceivedAt: now,
		})
		if err != nil {
			log.Printf("callback receipt reserve error | action=%s | callbackID=%s | err=%v", action, cb.ID, err)
			bootstrap.IncAppMetric("telegram.callback.receipt_reserve_error_total")
		} else if !reserved {
			log.Printf("duplicate callback receipt ignored | action=%s | callbackID=%s | chatID=%d | messageID=%d | userID=%d", action, cb.ID, chatID, messageID, userID)
			bootstrap.IncAppMetric("telegram.callback.duplicate_receipt_total")
			_, _ = h.presenter.Render(newAckCallbackRender(cb.ID))
			return
		}
	}

	if h.callbackDedupes != nil && shouldGuardTelegramCallback(action) && messageID > 0 {
		dedupeKey, windowStart := telegramCallbackSemanticDedupeKey(chatID, messageID, userID, action, params, now, telegramCallbackGuardTTL)
		reserved, err := h.callbackDedupes.Reserve(context.Background(), &domain.TelegramCallbackDedupe{
			DedupeKey:   dedupeKey,
			Action:      action,
			WindowStart: windowStart,
		})
		if err != nil {
			log.Printf("callback dedupe reserve error | action=%s | key=%s | err=%v", action, dedupeKey, err)
			bootstrap.IncAppMetric("telegram.callback.semantic_reserve_error_total")
		} else if !reserved {
			log.Printf("duplicate semantic callback ignored | action=%s | chatID=%d | messageID=%d | userID=%d", action, chatID, messageID, userID)
			bootstrap.IncAppMetric("telegram.callback.duplicate_semantic_total")
			_, _ = h.presenter.Render(newAckCallbackRender(cb.ID))
			return
		}
	}

	if h.callbackGuard != nil && !h.callbackGuard.TryReserve(chatID, messageID, userID, action, params) {
		log.Printf("duplicate callback ignored | action=%s | params=%v | chatID=%d | messageID=%d | userID=%d", action, params, chatID, messageID, userID)
		bootstrap.IncAppMetric("telegram.callback.duplicate_local_total")
		_, _ = h.presenter.Render(newAckCallbackRender(cb.ID))
		return
	}

	switch action {
	case "to_start", "select_man", "select_woman":
		inputTelegram = h.handleOnboarding(action, userID)
	case "action_menu", "menu_adventure", "menu_character", "menu_economy", "menu_events":
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
	case "menu_inventory_scope", "open_inventory", "open_inventory_potions", "open_inventory_equipaments", "open_inventory_equipaments_type", "open_inventory_equipped", "use_item", "equip_item", "unequip_item", "unequip_all":
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

	_, _ = h.presenter.Render(newEditPanelRender(chatID, messageID, inputTelegram))

	_, _ = h.presenter.Render(newAckCallbackRender(cb.ID))
}
