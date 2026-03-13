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
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Aventura", "menu_adventure"),
			tgbotapi.NewInlineKeyboardButtonData("🗺 Personagem", "menu_character"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Inventario", "open_inventory"),
			tgbotapi.NewInlineKeyboardButtonData("💰 Economia", "menu_economy"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧳 Eventos", "menu_events"),
			tgbotapi.NewInlineKeyboardButtonData("🏆 Rankings", "open_rank"),
		),
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL(
			"💬 Entrar na Comunidade",
			"https://t.me/dragonki_community",
		),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
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
			tgbotapi.NewInlineKeyboardButtonData("📿 Equipados", "open_inventory_equipped"),
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
			tgbotapi.NewInlineKeyboardButtonData("🧪 Usar Pocao", fmt.Sprintf("action_battle_open_potions?monsterID=%d&page=1", monsterID)),
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", fmt.Sprintf("action_battle_transform?monsterID=%d", monsterID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Atacar", fmt.Sprintf("action_fight?monsterID=%d", monsterID)),
			tgbotapi.NewInlineKeyboardButtonData("🏃🏻 Escapar", fmt.Sprintf("action_escape?monsterID=%d", monsterID)),
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
func shopEquipmentActionMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛒 Compra", "open_shop_equipment_types"),
			tgbotapi.NewInlineKeyboardButtonData("💰 Venda", "list_sell_equipments?page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "to_shop"),
		),
	)
	return &menu
}
func shopSellItemsMenu(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	return utils.ShopSellItemsMenu(itemsParams)
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

func adventureMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Cacar", "action_hunt"),
			tgbotapi.NewInlineKeyboardButtonData("🗺 Explorar", "list_map"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⛩ Ki Hall", "to_training"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}

func characterMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 Status", "action_status"),
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "list_transform"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🛏️ Pousada", "to_sleep"),
			tgbotapi.NewInlineKeyboardButtonData("🎁 Convites", "open_referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}

func inventoryScopeMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎒 Abrir Inventario", "open_inventory"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📿 Equipados", "open_inventory_equipped"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}

func economyMenu() *tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏪 Lojinha", "to_shop"),
			tgbotapi.NewInlineKeyboardButtonData("🔄 Cambio", "to_exchange"),
		),
	}
	if bootstrap.IsFeatureEnabled("vip", false) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👑 VIP", "to_vip"),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
	))
	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func vipMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👑 Ativar VIP", "confirm_buy_vip"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "menu_economy"),
		),
	)
	return &menu
}

func eventsMenu() *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)

	if bootstrap.IsFeatureEnabled("big_boss_weekly", false) {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Big Boss Semanal", "to_big_boss"),
		))
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Voltar", "action_menu"),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
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

func pixPaymentMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Copiar codigo PIX", "action_copy_pix_code"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "to_exchange"),
		),
	)
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

func bigBossMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Atacar", "big_boss_attack"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏆 Ranking", "big_boss_ranking"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "action_menu"),
		),
	)
	return &menu
}

func bigBossRankingMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Voltar", "to_big_boss"),
		),
	)
	return &menu
}

func ancestralKiChamberMenu() *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 7)

	for i := 1; i <= 5; i++ {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Iniciar Nivel %d", i),
				fmt.Sprintf("start_ancestral_ki_run?level=%d", i),
			),
		))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Niveis 6 a 10", "ancestral_ki_levels_6_10"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "action_menu"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func ancestralKiChamberLevels6To10Menu() *tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 7)

	for i := 6; i <= 10; i++ {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("Iniciar Nivel %d", i),
				fmt.Sprintf("start_ancestral_ki_run?level=%d", i),
			),
		))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Niveis 1 a 5", "to_ancestral_ki_chamber"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "menu_events"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func ancestralKiChamberRunMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Atacar Batalha Atual", "ancestral_ki_chamber_continue"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Sair da Camara", "abandon_ancestral_ki_run"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "to_ancestral_ki_chamber"),
		),
	)
	return &menu
}

func ancestralKiChamberPostWinMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Continuar", "ancestral_ki_chamber_continue"),
			tgbotapi.NewInlineKeyboardButtonData("Sair", "abandon_ancestral_ki_run"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Usar Pocoes", "chamber_open_potions?page=1"),
			tgbotapi.NewInlineKeyboardButtonData("Transformar", "chamber_transform"),
		),
	)
	return &menu
}
func postBattleWinMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚔️ Cacar Novamente", "action_hunt"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧪 Usar Pocao", "action_post_battle_open_potions?page=1"),
			tgbotapi.NewInlineKeyboardButtonData("🔥 Transformar", "action_post_battle_transform"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Menu Principal", "action_menu"),
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
