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
			tgbotapi.NewInlineKeyboardButtonData("🗺 Explorar", "list_map"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "list_transform"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧳 Inventory", "open_inventory"),
			tgbotapi.NewInlineKeyboardButtonData("⛩ Câm. Ki Ancestral", "to_training"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏪 Lojinha", "to_shop"),
			tgbotapi.NewInlineKeyboardButtonData("🛏️ Pousada", "to_sleep"),
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
			tgbotapi.NewInlineKeyboardButtonData("💤 -500💰 / +100%❤️", "action_sleep"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🥱 -1⚡️/ +30%❤️", "action_rest"),
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
func regionsMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.RegionsMenu(itemsParams)
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
				fmt.Sprintf("🔄 1💎 em %d de 💰", kiCrystalsConvert),
				"action_convert?amountConvert=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🔄 5💎 em %d de 💰", 5*kiCrystalsConvert),
				"action_convert?amountConvert=5"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("🔄 10💎 em %d de 💰", 10*kiCrystalsConvert),
				"action_convert?amountConvert=10"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Comprar %d 💎 por R$%d,00", kiCrystalsMult*kiCrystalsValue, kiCrystalsValue),
				"action_buy_ki"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}
