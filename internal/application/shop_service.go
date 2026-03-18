package application

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type ShopService struct {
	PlayerRepository    ports.PlayerRepository
	ItemRepository      ports.ItemsRepository
	InventoryRepository ports.InventoryRepository
	vipService          *VIPService
}

func NewShopService(repo ports.PlayerRepository, itemRepo ports.ItemsRepository, inventoryRepo ports.InventoryRepository) *ShopService {
	return &ShopService{
		PlayerRepository:    repo,
		ItemRepository:      itemRepo,
		InventoryRepository: inventoryRepo,
	}
}

func (s *ShopService) SetVIPService(vipService *VIPService) {
	s.vipService = vipService
}

func (s *ShopService) OpenShop(userID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	playerName := displayPlayerName(player, s.vipService)
	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n"+
			"Seja bem-vindo(a) %s à minha lojinha! \n\n"+
			"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
			"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
			"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
			"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
			"Prepare-se bem… a jornada de um verdadeiro guerreiro(a) começa com as escolhas certas. \n\n"+
			"O que posso fazer por você hoje? 😄 \n\n",
		playerName)

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey: "SHOP_MENU",
	}, nil
}

func (s *ShopService) ListPotions() (view.GameResponse, error) {
	potions, err := s.ItemRepository.GetRandomList(domain.ItemConsumable, domain.Common)
	if err != nil {
		log.Fatal(err)
	}

	menuParams := make(map[string]any)
	for _, potion := range potions {
		label := fmt.Sprintf("%d x %s - %d Gold", potion.Amount, potion.Name, potion.Price)
		menuParams[label] = strconv.Itoa(int(potion.ID))
	}

	msg := "Kiara, a Mercadora\n\nEssas sao as pocoes disponiveis hoje.\n\n"

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey:    "SHOP_ITEMS_MENU",
		MenuParams: menuParams,
	}, nil
}

func (s *ShopService) ListEquipments() (view.GameResponse, error) {
	msg := "Kiara, a Mercadora\n\nEquipamentos:\nEscolha se deseja comprar ou vender.\n\n"

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey: "SHOP_EQUIPMENT_ACTION_MENU",
	}, nil
}

func (s *ShopService) OpenBuyEquipmentTypes() (view.GameResponse, error) {
	msg := "Kiara, a Mercadora\n\nEscolha o tipo de equipamento para compra.\n\n"
	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey: "SHOP_EQUIPMENT_TYPES_MENU",
	}, nil
}

func (s *ShopService) ListEquipmentsByType(rawSlot string) (view.GameResponse, error) {
	slot, ok := normalizeShopEquipmentSlot(rawSlot)
	if !ok {
		return view.GameResponse{Text: "Tipo de equipamento invalido.", MenuKey: "SHOP_EQUIPMENT_TYPES_MENU"}, nil
	}

	equipments, err := s.ItemRepository.GetRandomEquipmentsBySlot(slot, 3)
	if err != nil {
		log.Fatal(err)
	}

	menuParams := make(map[string]any)
	for _, equipment := range equipments {
		label := fmt.Sprintf("%d x %s - %d Gold", equipment.Amount, equipment.Name, equipmentPurchasePrice(&equipment, 1))
		menuParams[label] = strconv.Itoa(int(equipment.ID))
	}
	menuParams["_backAction"] = "list_equipments"

	msg := fmt.Sprintf("Kiara, a Mercadora\n\nOpcoes de %s de hoje.\n\n", shopEquipmentTypeLabel(slot))
	for _, equipment := range equipments {
		msg += fmt.Sprintf("• %s\n", equipment.Name)
		msg += fmt.Sprintf("  Requisitos: %s\n", itemRequirementText(&equipment))
		if bonusText := itemPositiveBonusText(&equipment); bonusText != "" {
			msg += fmt.Sprintf("  Bonus: %s\n", bonusText)
		}
		msg += "\n"
	}

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey:    "SHOP_ITEMS_MENU",
		MenuParams: menuParams,
	}, nil
}

func (s *ShopService) BuyItem(userID int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	itemIDStr := params["itemID"]
	itemID, _ := strconv.ParseInt(itemIDStr, 10, 64)
	item, err := s.ItemRepository.GetItemById(int(itemID))
	if err != nil {
		log.Fatal(err)
	}

	priceMultiplier, _ := s.vipMultipliers(int64(player.ID))
	price := equipmentPurchasePrice(item, priceMultiplier)
	msg := "Kiara, a Mercadora\n\nVoce nao tem gold suficiente para este item.\n\n"

	if err := s.InventoryRepository.PurchaseItem(int64(player.ID), item.ID, item.Amount, price); err != nil {
		if !errors.Is(err, domain.ErrInventoryInsufficientGold) {
			return view.GameResponse{}, err
		}
	} else {
		msg = fmt.Sprintf("Kiara, a Mercadora\n\nObrigada pela compra. Volte sempre!\n\nItem: %s\n", item.Name)
		msg += fmt.Sprintf("Requisitos: %s\n", itemRequirementText(item))
		if bonusText := itemPositiveBonusText(item); bonusText != "" {
			msg += fmt.Sprintf("Bonus: %s\n", bonusText)
		}
		msg += "\n"
	}

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey: "SHOP_MENU",
	}, nil
}

func (s *ShopService) ListSellEquipments(userID int64, page int) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	if page < 1 {
		page = 1
	}

	items, total, err := s.InventoryRepository.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, page, 5)
	if err != nil {
		return view.GameResponse{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(5)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
		items, total, err = s.InventoryRepository.ListItemsByUserAndType(int64(player.ID), domain.ItemEquipment, page, 5)
		if err != nil {
			return view.GameResponse{}, err
		}
		totalPages = int(math.Ceil(float64(total) / float64(5)))
		if totalPages == 0 {
			totalPages = 1
		}
	}

	msg := "Kiara, a Mercadora\n\nEscolha um equipamento para vender.\nItens equipados nao aparecem nesta lista.\n\n"
	if len(items) == 0 {
		msg = "Kiara, a Mercadora\n\nVoce nao possui equipamentos disponiveis para venda.\n\n"
	}

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/kiara_mercadora.png",
		},
		MenuKey: "SHOP_SELL_ITEMS_MENU",
		MenuParams: map[string]any{
			"Items":            items,
			"page":             page,
			"totalPages":       totalPages,
			"paginationAction": "list_sell_equipments",
			"backAction":       "list_equipments",
		},
	}, nil
}

func (s *ShopService) SellItem(userID int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	inventoryItemID, convErr := strconv.Atoi(strings.TrimSpace(params["inventoryItemID"]))
	if convErr != nil || inventoryItemID <= 0 {
		return view.GameResponse{
			Text: "Kiara, a Mercadora\n\nItem invalido para venda.\n\n",
			Image: &view.Image{
				Type: view.ImageAsset,
				Path: "assets/npcs/kiara_mercadora.png",
			},
			MenuKey: "SHOP_MENU",
		}, nil
	}

	page, _ := strconv.Atoi(strings.TrimSpace(params["page"]))
	if page < 1 {
		page = 1
	}

	invItem, err := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), uint(inventoryItemID))
	if err != nil {
		return view.GameResponse{}, err
	}
	if invItem == nil || invItem.Quantity <= 0 || invItem.Item.Type != domain.ItemEquipment {
		listResp, listErr := s.ListSellEquipments(userID, page)
		if listErr != nil {
			return view.GameResponse{}, listErr
		}
		listResp.Text = "Kiara, a Mercadora\n\nNao foi possivel vender este item.\n\n" + listResp.Text
		return listResp, nil
	}

	_, sellMultiplier := s.vipMultipliers(int64(player.ID))
	sellValue := equipmentSellPrice(&invItem.Item, sellMultiplier)
	soldItem, err := s.InventoryRepository.SellInventoryItem(int64(player.ID), invItem.ID, sellValue)
	if err != nil {
		if errors.Is(err, domain.ErrInventoryItemNotSellable) || errors.Is(err, domain.ErrInventoryItemEquipped) {
			listResp, listErr := s.ListSellEquipments(userID, page)
			if listErr != nil {
				return view.GameResponse{}, listErr
			}
			listResp.Text = "Kiara, a Mercadora\n\nNao foi possivel vender este item.\n\n" + listResp.Text
			return listResp, nil
		}
		return view.GameResponse{}, err
	}
	if soldItem != nil {
		invItem = soldItem
	}

	listResp, listErr := s.ListSellEquipments(userID, page)
	if listErr != nil {
		return view.GameResponse{}, listErr
	}
	listResp.Text = fmt.Sprintf(
		"Kiara, a Mercadora\n\n%s vendido por %d Gold.\n\n%s",
		invItem.Item.Name,
		sellValue,
		listResp.Text,
	)
	return listResp, nil
}

func equipmentPurchasePrice(item *domain.Item, multiplier float64) int {
	if item == nil {
		return 0
	}
	if multiplier <= 0 {
		multiplier = 1
	}
	value := int(math.Floor(float64(item.Price) * multiplier))
	if value < 1 {
		value = 1
	}
	return value
}

func equipmentSellPrice(item *domain.Item, multiplier float64) int {
	if item == nil || item.Price <= 0 {
		return 0
	}
	if multiplier <= 0 {
		multiplier = 1
	}
	value := int(math.Floor(float64(item.Price) * 0.4 * multiplier))
	if value < 1 {
		value = 1
	}
	return value
}

func (s *ShopService) vipMultipliers(playerID int64) (priceMultiplier float64, sellMultiplier float64) {
	priceMultiplier = 1.0
	sellMultiplier = 1.0
	if s == nil || s.vipService == nil {
		return priceMultiplier, sellMultiplier
	}
	return s.vipService.GetShopPriceAndSellMultiplier(playerID)
}

func normalizeShopEquipmentSlot(raw string) (domain.EquipmentSlot, bool) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "WEAPON":
		return domain.SlotWeapon, true
	case "ARMOR":
		return domain.SlotArmor, true
	case "GLOVES":
		return domain.SlotGloves, true
	case "HELMET":
		return domain.SlotHelmet, true
	case "BOOTS":
		return domain.SlotBoots, true
	case "RING", "ACCESSORY":
		return domain.SlotRing, true
	case "AMULET":
		return domain.SlotAmulet, true
	default:
		return "", false
	}
}

func shopEquipmentTypeLabel(slot domain.EquipmentSlot) string {
	switch slot {
	case domain.SlotWeapon:
		return "armas"
	case domain.SlotArmor:
		return "armaduras"
	case domain.SlotGloves:
		return "luvas"
	case domain.SlotHelmet:
		return "capacetes"
	case domain.SlotBoots:
		return "botas"
	case domain.SlotRing:
		return "aneis"
	case domain.SlotAmulet:
		return "amuletos"
	case domain.SlotAccessory:
		return "acessorios"
	default:
		return string(slot)
	}
}

func itemPositiveBonusText(item *domain.Item) string {
	if item == nil {
		return ""
	}
	bonuses := make([]string, 0, 5)
	if item.AtkBonus > 0 {
		bonuses = append(bonuses, fmt.Sprintf("ATK +%d", item.AtkBonus))
	}
	if item.DefBonus > 0 {
		bonuses = append(bonuses, fmt.Sprintf("DEF +%d", item.DefBonus))
	}
	if item.SpeedBonus > 0 {
		bonuses = append(bonuses, fmt.Sprintf("VEL +%d", item.SpeedBonus))
	}
	if item.HPBonus > 0 {
		bonuses = append(bonuses, fmt.Sprintf("HP +%d", item.HPBonus))
	}
	if item.EnergyBonus > 0 {
		bonuses = append(bonuses, fmt.Sprintf("ENERGIA +%d", item.EnergyBonus))
	}
	return strings.Join(bonuses, " | ")
}

func itemRequirementText(item *domain.Item) string {
	if item == nil {
		return "N/A"
	}

	levelText := "Nivel 1+"
	if item.RequiredLevel > 0 {
		levelText = fmt.Sprintf("Nivel %d+", item.RequiredLevel)
	}

	classText := "Todas as classes"
	if strings.TrimSpace(item.AllowedClasses) != "" {
		classText = strings.TrimSpace(item.AllowedClasses)
	}

	return fmt.Sprintf("%s | Classes: %s", levelText, classText)
}
