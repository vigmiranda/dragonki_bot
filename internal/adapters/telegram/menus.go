package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func mainMenu() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "action_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Caçar", "action_cacar"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "action_transformar"),
			tgbotapi.NewInlineKeyboardButtonData("⚡ Energia", "action_energia"),
		),
	)
}
