package telegram

import (
	"fmt"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram/utils"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
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
			tgbotapi.NewInlineKeyboardButtonData("🏆 Rank", "open_rank"),
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
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar -3⚡️", "action_transform"),
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
func shopEquipmentTypesMenu() *tgbotapi.InlineKeyboardMarkup {
	return utils.ShopEquipmentTypesMenu()
}
func potionMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.PotionMenu(itemsParams)
}
func equipmentTypesMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.EquipmentTypesMenu(itemsParams)
}
func inventoryEquippedMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.EquippedItemsMenu(itemsParams)
}
func regionsMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.RegionsMenu(itemsParams)
}
func choiceClassMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.ChoiceClassMenu(itemsParams)
}

func exchangeMenu() *tgbotapi.InlineKeyboardMarkup {
	kiCrystalsValue := getenvIntDefault("KI_CRYSTALS_VALUE", 1)
	kiCrystalsConvert := getenvIntDefault("KI_CRYSTALS_CONVERT", 1)
	kiCrystalsMult := getenvIntDefault("KI_CRYSTALS_MULT", 1)

	rows := [][]tgbotapi.InlineKeyboardButton{
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
	}

	if bootstrap.IsFeatureEnabled("payments", true) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Comprar %d 💎 por R$%d,00", kiCrystalsMult*kiCrystalsValue, kiCrystalsValue),
				"action_buy_ki"),
		))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func rankMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Top 5 Batalhas", "rank_top_battles"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Top 5 Vitorias", "rank_top_wins"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Top 5 Derrotas", "rank_top_losses"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Top 5 Escapes", "rank_top_escape"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "action_menu"),
		),
	)
	return &menu
}
func getenvIntDefault(key string, def int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return def
	}
	return v
}
