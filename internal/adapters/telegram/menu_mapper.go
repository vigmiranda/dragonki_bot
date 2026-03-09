package telegram

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func mapMenu(
	key string,
	params map[string]any,
) *tgbotapi.InlineKeyboardMarkup {
	switch key {
	case "GENDER_MENU":
		return choosingGenderMenu()

	case "MAIN_MENU":
		return mainMenu()

	case "ENERGY_MENU":
		return energyMenu()

	case "SHOP_MENU":
		return shopMenu()

	case "EXCHANGE_MENU":
		return exchangeMenu()

	case "PIX_PAYMENT_MENU":
		return pixPaymentMenu()

	case "RANK_MENU":
		return rankMenu()

	case "BIG_BOSS_MENU":
		return bigBossMenu()

	case "BIG_BOSS_RANKING_MENU":
		return bigBossRankingMenu()

	case "ANCESTRAL_KI_CHAMBER_MENU":
		return ancestralKiChamberMenu()

	case "ANCESTRAL_KI_CHAMBER_LEVELS_6_10_MENU":
		return ancestralKiChamberLevels6To10Menu()

	case "ANCESTRAL_KI_CHAMBER_RUN_MENU":
		return ancestralKiChamberRunMenu()

	case "ANCESTRAL_KI_CHAMBER_POST_WIN_MENU":
		return ancestralKiChamberPostWinMenu()

	case "SHOP_ITEMS_MENU":
		return shopItemsMenu(params)

	case "SHOP_EQUIPMENT_TYPES_MENU":
		return shopEquipmentTypesMenu()

	case "SHOP_EQUIPMENT_ACTION_MENU":
		return shopEquipmentActionMenu()

	case "SHOP_SELL_ITEMS_MENU":
		return shopSellItemsMenu(params)

	case "INVENTORY_MENU":
		return inventoryMenu()

	case "INVENTORY_ITEMS_MENU":
		return potionMenu(params)

	case "INVENTORY_EQUIPMENT_TYPES_MENU":
		return equipmentTypesMenu(params)

	case "INVENTORY_EQUIPPED_MENU":
		return inventoryEquippedMenu(params)

	case "TRANSFORM_MENU":
		return transformMenu()

	case "REGIONS_MENU":
		return regionsMenu(params)

	case "CHOICE_CLASS":
		return choiceClassMenu(params)

	case "BATTLE_MENU":
		monsterID, ok := params["monsterID"]
		monster, err := strconv.ParseInt(fmt.Sprint(monsterID), 10, 64)
		if err != nil {
			fmt.Println("Error parsing string:", err)
		}
		if !ok {
			log.Println("❌ monsterID inválido no menu")
			return mainMenu()
		}
		return battleMenu(monster)

	default:
		return mainMenu()
	}
}
