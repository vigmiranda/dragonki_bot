package utils

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func BuildMenuItems(itemsParams map[string]any) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	backAction := "to_shop"

	for name, id := range itemsParams {
		if strings.HasPrefix(name, "_") {
			if name == "_backAction" {
				backAction = fmt.Sprint(id)
			}
			continue
		}
		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s", name),
			fmt.Sprintf("buy_item?itemID=%s", id),
		)

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(btn))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", backAction),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func ShopEquipmentTypesMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Armas", "list_equipments_type?slot=WEAPON"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Armaduras", "list_equipments_type?slot=ARMOR"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Luvas", "list_equipments_type?slot=GLOVES"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Capacetes", "list_equipments_type?slot=HELMET"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Botas", "list_equipments_type?slot=BOOTS"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Aneis", "list_equipments_type?slot=RING"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Amuletos", "list_equipments_type?slot=AMULET"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "list_equipments"),
		),
	)
	return &menu
}

func ShopSellItemsMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	rawItems, ok := params["Items"].([]domain.InventoryItem)
	if !ok {
		fmt.Printf("ShopSellItemsMenu -> %v\\n", ok)
	}

	page := toInt(params["page"])
	if page < 1 {
		page = 1
	}
	totalPages := toInt(params["totalPages"])
	if totalPages < 1 {
		totalPages = 1
	}

	paginationAction := paramString(params, "paginationAction")
	if paginationAction == "" {
		paginationAction = "list_sell_equipments"
	}
	backAction := paramString(params, "backAction")
	if backAction == "" {
		backAction = "list_equipments"
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, item := range rawItems {
		price := int(float64(item.Item.Price) * 0.4)
		label := fmt.Sprintf("Vender %s (+%d Gold)", item.Item.Name, price)
		callback := fmt.Sprintf("sell_item?inventoryItemID=%d&page=%d", item.ID, page)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callback),
		))
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if page > 1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData("<", withPage(paginationAction, page-1)))
	}
	if page < totalPages {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(">", withPage(paginationAction, page+1)))
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", backAction),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func ChoiceClassMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	rawClasses, ok := params["Classes"].([]domain.Class)
	if !ok {
		fmt.Printf("ClassesMenu -> %v\\n", ok)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, class := range rawClasses {
		label := fmt.Sprintf("%s", class.Name)
		callback := fmt.Sprintf("to_class?classID=%d", class.ID)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callback),
			),
		)
	}
	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func RegionsMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	rawRegions, ok := params["Regions"].([]domain.Region)
	if !ok {
		fmt.Printf("RegionsMenu -> %v\\n", ok)
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, region := range rawRegions {
		label := fmt.Sprintf("%s - Lvl (%d)", region.Name, region.MinLevel)
		callback := fmt.Sprintf("to_region?regionID=%d", region.RegionID)

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callback),
			),
		)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "action_menu"),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func PotionMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	rawItems, ok := params["Items"].([]domain.InventoryItem)
	if !ok {
		fmt.Printf("PotionMenu -> %v\\n", ok)
	}

	page := toInt(params["page"])
	if page < 1 {
		page = 1
	}
	totalPages := toInt(params["totalPages"])
	if totalPages < 1 {
		totalPages = 1
	}

	itemType := paramString(params, "itemType")
	isEquipment := itemType == string(domain.ItemEquipment)

	paginationAction := paramString(params, "paginationAction")
	if paginationAction == "" {
		paginationAction = "open_inventory_potions"
	}
	backAction := paramString(params, "backAction")
	if backAction == "" {
		backAction = "open_inventory"
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, item := range rawItems {
		label := fmt.Sprintf("%dx %s", item.Quantity, item.Item.Name)
		callback := fmt.Sprintf("use_item?inventoryItemID=%d", item.ID)
		if isEquipment {
			callback = fmt.Sprintf("equip_item?inventoryItemID=%d", item.ID)
		}

		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(label, callback),
			),
		)
	}

	var navRow []tgbotapi.InlineKeyboardButton
	if page > 1 {
		navRow = append(navRow,
			tgbotapi.NewInlineKeyboardButtonData("<", withPage(paginationAction, page-1)),
		)
	}
	if page < totalPages {
		navRow = append(navRow,
			tgbotapi.NewInlineKeyboardButtonData(">", withPage(paginationAction, page+1)),
		)
	}
	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", backAction),
		),
	)

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func EquipmentTypesMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	weaponCount := toInt(params["weaponCount"])
	armorCount := toInt(params["armorCount"])
	glovesCount := toInt(params["glovesCount"])
	helmetCount := toInt(params["helmetCount"])
	bootsCount := toInt(params["bootsCount"])
	ringCount := toInt(params["ringCount"])
	amuletCount := toInt(params["amuletCount"])

	menu := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Armas (%d)", weaponCount), "open_inventory_equipaments_type?slot=WEAPON&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Armaduras (%d)", armorCount), "open_inventory_equipaments_type?slot=ARMOR&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Luvas (%d)", glovesCount), "open_inventory_equipaments_type?slot=GLOVES&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Capacetes (%d)", helmetCount), "open_inventory_equipaments_type?slot=HELMET&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Botas (%d)", bootsCount), "open_inventory_equipaments_type?slot=BOOTS&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Aneis (%d)", ringCount), "open_inventory_equipaments_type?slot=RING&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("Amuletos (%d)", amuletCount), "open_inventory_equipaments_type?slot=AMULET&page=1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Itens Equipados", "open_inventory_equipped"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "open_inventory"),
		),
	)
	return &menu
}

func EquippedItemsMenu(params map[string]any) *tgbotapi.InlineKeyboardMarkup {
	weaponAction := paramString(params, "weaponAction")
	helmetAction := paramString(params, "helmetAction")
	bootsAction := paramString(params, "bootsAction")
	armorAction := paramString(params, "armorAction")
	glovesAction := paramString(params, "glovesAction")
	ringAction := paramString(params, "ringAction")
	amuletAction := paramString(params, "amuletAction")
	unequipAllAction := paramString(params, "unequipAllAction")

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 10)
	addRow := func(label, callback string) {
		if strings.TrimSpace(callback) == "" {
			return
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(label, callback),
		))
	}

	addRow("Desequipar Arma", weaponAction)
	addRow("Desequipar Capacete", helmetAction)
	addRow("Desequipar Botas", bootsAction)
	addRow("Desequipar Armadura", armorAction)
	addRow("Desequipar Luva", glovesAction)
	addRow("Desequipar Anel", ringAction)
	addRow("Desequipar Amuleto", amuletAction)
	addRow("Desequipar Tudo", unequipAllAction)

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Voltar", "open_inventory"),
	))

	menu := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &menu
}

func withPage(action string, page int) string {
	separator := "?"
	if strings.Contains(action, "?") {
		separator = "&"
	}
	return fmt.Sprintf("%s%spage=%d", action, separator, page)
}

func paramString(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok || v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
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
