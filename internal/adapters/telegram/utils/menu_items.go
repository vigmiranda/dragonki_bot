package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func BuildMenuItems(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for name, id := range itemsParams {
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", name),
			fmt.Sprintf("buy_item?itemID=%s", id),
		)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	// botão de cancelar / voltar
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "to_shop"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}
func RegionsMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	rawRegions, ok := params["Regions"].([]domain.Region)
	if !ok {
		fmt.Printf("RegionsMenu -> ", ok)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, region := range rawRegions {

		label := fmt.Sprintf(
			"%s - Lvl (%d)",
			region.Name,
			region.MinLevel,
		)

		callback := fmt.Sprintf(
			"to_region?regionID=%d",
			region.RegionID,
		)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callback),
			),
		)
	}

	// 🔹 Botão voltar
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Voltar", "action_menu"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}
func PotionMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {

	rawItems, ok := params["Items"].([]domain.InventoryItem)
	if !ok {
		fmt.Printf("PotionMenu -> ", ok)
	}

	page := toInt(params["page"])
	totalPages := toInt(params["totalPages"])

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, item := range rawItems {

		label := fmt.Sprintf(
			"%dx %s",
			item.Quantity,
			item.Item.Name,
		)

		callback := fmt.Sprintf(
			"use_item?inventoryItemID=%d",
			item.ID,
		)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callback),
			),
		)
	}

	// 🔹 Paginação
	var navRow []tgbotapi.InlineKeyboardButton

	if page > 1 {
		navRow = append(navRow,
			tgbotapi.NewInlineKeyboardButtonData(
				"⬅️",
				fmt.Sprintf("inventory?page=%d", page-1),
			),
		)
	}

	if page < totalPages {
		navRow = append(navRow,
			tgbotapi.NewInlineKeyboardButtonData(
				"➡️",
				fmt.Sprintf("inventory?page=%d", page+1),
			),
		)
	}

	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	// 🔹 Botão voltar
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Voltar", "open_inventory"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func toInt(v any) int {
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	default:
		return 0
	}
}
