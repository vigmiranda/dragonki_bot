package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
