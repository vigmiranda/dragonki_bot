package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func choosingGenderMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♂️ Homem", "select_man"),
			tgbotapi.NewInlineKeyboardButtonData("♀️ Mulher", "select_woman"),
		),
	)
	return &menu
}

func mainMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "action_menu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Caçar", "action_cacar"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar (Lv +3)", "action_transformar"),
			tgbotapi.NewInlineKeyboardButtonData("🛏️ Pousada", "to_sleep"),
		),
		//tgbotapi.NewInlineKeyboardRow(
		//	tgbotapi.NewInlineKeyboardButtonData("💎 Comprar 10 KiCrystals", "action_buy_ki"),
		//),
	)
	return &menu
}
func energyMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💤 Descansar (-500 💰 Gold)", "action_sleep"),
			tgbotapi.NewInlineKeyboardButtonData("📖 Menu", "action_menu"),
		),
	)
	return &menu
}
func battleMenu(monsterID int) *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Atacar ", fmt.Sprintf("action_fight?monsterID=%d", monsterID)),
			//tgbotapi.NewInlineKeyboardButtonData("🧪 Poções ", "use_potion"),
			tgbotapi.NewInlineKeyboardButtonData("🏃🏻 Escapar ", fmt.Sprintf("action_escape?monsterID=%d", monsterID)),
		),
	)
	return &menu
}
