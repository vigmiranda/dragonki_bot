package application

import (
	"fmt"
	"log"
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
}

func NewShopService(repo ports.PlayerRepository, itemRepo ports.ItemsRepository, inventoryRepo ports.InventoryRepository) *ShopService {
	return &ShopService{
		PlayerRepository:    repo,
		ItemRepository:      itemRepo,
		InventoryRepository: inventoryRepo,
	}
}

func (s *ShopService) OpenShop(userID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n"+
			"Seja bem-vindo(a) %s à minha lojinha! \n\n"+
			"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
			"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
			"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
			"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
			"Prepare-se bem… a jornada de um verdadeiro guerreiro(a) começa com as escolhas certas. \n\n"+
			"O que posso fazer por você hoje? 😄 \n\n",
		player.Name)

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
	msg := "Kiara, a Mercadora\n\nEscolha o tipo de equipamento.\nMostro ate 3 opcoes sortidas por tipo.\n\n"

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
		label := fmt.Sprintf("%d x %s - %d Gold", equipment.Amount, equipment.Name, equipmentPurchasePrice(&equipment))
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

	price := equipmentPurchasePrice(item)
	msg := "Kiara, a Mercadora\n\nVoce nao tem gold suficiente para este item.\n\n"

	if player.Gold >= price {
		if err := s.InventoryRepository.CreateOrUpdateInventory(int64(player.ID), item.ID, item.Amount); err != nil {
			return view.GameResponse{}, err
		}
		player.Gold -= price
		if err := s.PlayerRepository.Update(player); err != nil {
			return view.GameResponse{}, err
		}
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

func equipmentPurchasePrice(item *domain.Item) int {
	if item == nil {
		return 0
	}
	if item.Type == domain.ItemEquipment && item.Price < 1500 {
		return 1500
	}
	return item.Price
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
