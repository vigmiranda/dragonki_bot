package telegram

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
)

func onboardingMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Iniciar", "to_start"),
		),
	)
	return &menu
}

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
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Caçar", "action_hunt"),
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "list_transform"),
			tgbotapi.NewInlineKeyboardButtonData("⛩ Câmara do Ki Ancestral", "to_training"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏪 Lojinha", "to_shop"),
			tgbotapi.NewInlineKeyboardButtonData("🛏️ Pousada", "to_sleep"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧳 Inventory", "open_inventory"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Casa de Câmbio", "to_exchange"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				"💬 Entrar na Comunidade",
				"https://t.me/dragonki_community",
			),
		),
	)
	return &menu
}
func transformMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "action_transform"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
func inventoryMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧪 Poções", "open_inventory_potions"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📿 Equipamentos", "open_inventory_equipaments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
func energyMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💤 Dormir -500 💰Gold / +100% ❤️HP", "action_sleep"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥱 Descansar -1 ⚡️Energia / +30% ❤️HP", "action_rest"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
func battleMenu(monsterID int64) *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Atacar ", fmt.Sprintf("action_fight?monsterID=%d", monsterID)),
			//tgbotapi.NewInlineKeyboardButtonData("🧪 Poções ", "use_potion"),
			tgbotapi.NewInlineKeyboardButtonData("🏃🏻 Escapar ", fmt.Sprintf("action_escape?monsterID=%d", monsterID)),
		),
	)
	return &menu
}
func shopMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧪 Poções", "list_potions"),
			tgbotapi.NewInlineKeyboardButtonData("📿 Equipamentos", "list_equipments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
func shopItemsMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.BuildMenuItems(itemsParams)
}
func potionMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.PotionMenu(itemsParams)
}
func exchangeMenu() *tgbotapi.InlineKeyboardMarkup {
	kiCrystalsValue, err := strconv.Atoi(os.Getenv("KI_CRYSTALS_VALUE"))
	if err != nil {
		log.Fatal(err)
	}
	kiCrystalsConvert, err := strconv.Atoi(os.Getenv("KI_CRYSTALS_CONVERT"))
	if err != nil {
		log.Fatal(err)
	}
	kiCrystalsMult, err := strconv.Atoi(os.Getenv("KI_CRYSTALS_MULT"))
	if err != nil {
		log.Fatal(err)
	}

	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Converter 1 Ki Crystals 💎 em %d de Gold 💰", kiCrystalsConvert),
				"action_convert?amountConvert=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Converter 5 Ki Crystals 💎 em %d de Gold 💰", 5*kiCrystalsConvert),
				"action_convert?amountConvert=5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Converter 10 Ki Crystals 💎 em %d de Gold 💰", 10*kiCrystalsConvert),
				"action_convert?amountConvert=10"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("💎 Comprar %d Ki Crystals por R$%d,00", kiCrystalsMult*kiCrystalsValue, kiCrystalsValue),
				"action_buy_ki"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
