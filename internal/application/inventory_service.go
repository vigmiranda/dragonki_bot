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

type InventoryService struct {
	InventoryRepository ports.InventoryRepository
	ItemsRepository     ports.ItemsRepository
	PlayerRepository    ports.PlayerRepository
}

func NewInventoryService(
	inventoryRepository ports.InventoryRepository,
	playerRepo ports.PlayerRepository,
	itemsRepository ports.ItemsRepository,
) *InventoryService {
	return &InventoryService{
		InventoryRepository: inventoryRepository,
		PlayerRepository:    playerRepo,
		ItemsRepository:     itemsRepository,
	}
}

func (s InventoryService) Open(userID int64) (view.GameResponse, error) {
	msg := fmt.Sprintf(
		"🌌 Arkanis, o Guardião dos Ecos \n\n" +
			"“Eu sou Arkanis… aquele que vigia o que você conquista.\n" +
			"Cada item que carrega conta a história das suas batalhas.”\n\n" +
			"✨ “Organize-se. Prepare-se.\n" +
			"Um(a) guerreiro(a) forte não é feito só de Ki… mas de estratégia.”\n\n" +
			"📦 Abra seu inventário e escolha com sabedoria.\n\n",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/arkanis_inventario.png"),
	}
	menu := "INVENTORY_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}
	return viewResponse, nil
}

func (s InventoryService) ListItemsByUserAndType(playerID int64, itemType domain.ItemType) (view.GameResponse, error) {
	return s.ListItemsByUserAndTypePage(playerID, itemType, 1)
}

func (s InventoryService) ListItemsByUserAndTypePage(playerID int64, itemType domain.ItemType, page int) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		log.Fatal(err)
	}
	if page < 1 {
		page = 1
	}

	typeMessageItem := ""
	if itemType == "CONSUMABLE" {
		typeMessageItem = "suas 🧪 poções..\n\n"
	} else {
		typeMessageItem = "seus 📿 Equipamentos..\n\n"
	}

	msg := fmt.Sprintf(
		"🌌 Arkanis, o Guardião dos Ecos \n\n"+
			"“Aqui esta %s",
		typeMessageItem,
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/arkanis_inventario.png"),
	}
	menu := "INVENTORY_ITEMS_MENU"

	items, total, err := s.InventoryRepository.ListItemsByUserAndType(int64(player.ID), itemType, page, 5)
	if err != nil {
		return view.GameResponse{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(5)))
	if totalPages == 0 {
		totalPages = 1
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
		MenuParams: map[string]any{
			"Items":      items,
			"page":       page,
			"totalPages": totalPages,
			"itemType":   string(itemType),
			"paginationAction": func() string {
				if itemType == domain.ItemConsumable {
					return "open_inventory_potions"
				}
				return "open_inventory_equipaments"
			}(),
			"backAction": "open_inventory",
		},
	}

	return viewResponse, nil
}

func (s InventoryService) OpenEquipmentTypes(playerID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	counts, err := s.InventoryRepository.CountEquipmentItemsBySlot(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}
	if counts == nil {
		counts = map[domain.EquipmentSlot]int64{}
	}

	msg := "Escolha o tipo de equipamento para ver os itens do inventario:"
	img := &view.Image{
		Type: view.ImageAsset,
		Path: "assets/npcs/arkanis_inventario.png",
	}

	return view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: "INVENTORY_EQUIPMENT_TYPES_MENU",
		MenuParams: map[string]any{
			"weaponCount": counts[domain.SlotWeapon],
			"armorCount":  counts[domain.SlotArmor],
			"glovesCount": counts[domain.SlotGloves],
			"helmetCount": counts[domain.SlotHelmet],
			"bootsCount":  counts[domain.SlotBoots],
			"ringCount":   counts[domain.SlotRing] + counts[domain.SlotAccessory],
			"amuletCount": counts[domain.SlotAmulet],
		},
	}, nil
}

func (s InventoryService) ViewEquippedItems(playerID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	pe, err := s.InventoryRepository.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	describeSlot := func(inventoryItemID *uint) string {
		if inventoryItemID == nil || *inventoryItemID == 0 {
			return "vazio"
		}

		invItem, invErr := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), *inventoryItemID)
		if invErr != nil {
			log.Printf("view_equipped load inventory_item error | playerID=%d | inventoryItemID=%d | err=%v", player.ID, *inventoryItemID, invErr)
			return fmt.Sprintf("erro ao carregar (inventory_item_id=%d)", *inventoryItemID)
		}
		if invItem == nil {
			return fmt.Sprintf("nao encontrado (inventory_item_id=%d)", *inventoryItemID)
		}
		item := invItem.Item
		if item.ID == 0 && invItem.ItemID > 0 {
			dbItem, itemErr := s.ItemsRepository.GetItemById(int(invItem.ItemID))
			if itemErr != nil {
				log.Printf("view_equipped load item metadata error | playerID=%d | inventoryItemID=%d | itemID=%d | err=%v", player.ID, invItem.ID, invItem.ItemID, itemErr)
				return fmt.Sprintf("erro ao carregar item (inventory_item_id=%d)", invItem.ID)
			}
			if dbItem == nil {
				return fmt.Sprintf("item nao encontrado (inventory_item_id=%d)", invItem.ID)
			}
			item = *dbItem
		}
		return fmt.Sprintf("%s | %s", item.Name, equipmentBonusOnly(&item))
	}

	ringID := pe.RingInventoryItemID
	if (ringID == nil || *ringID == 0) && pe.AccessoryInventoryItemID != nil && *pe.AccessoryInventoryItemID > 0 {
		ringID = pe.AccessoryInventoryItemID
	}
	slotAction := func(slot string) string {
		return fmt.Sprintf("unequip_item?slot=%s", slot)
	}

	msg := strings.Join([]string{
		"Equipamentos atuais:",
		fmt.Sprintf("Arma: %s", describeSlot(pe.WeaponInventoryItemID)),
		fmt.Sprintf("Capacete: %s", describeSlot(pe.HelmetInventoryItemID)),
		fmt.Sprintf("Botas: %s", describeSlot(pe.BootsInventoryItemID)),
		"Armadura:",
		fmt.Sprintf("  Armadura: %s", describeSlot(pe.ArmorInventoryItemID)),
		fmt.Sprintf("  Luva: %s", describeSlot(pe.GlovesInventoryItemID)),
		"Acessorio:",
		fmt.Sprintf("  Anel: %s", describeSlot(ringID)),
		fmt.Sprintf("  Amuleto: %s", describeSlot(pe.AmuletInventoryItemID)),
	}, "\n")

	return view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/arkanis_inventario.png",
		},
		MenuKey: "INVENTORY_EQUIPPED_MENU",
		MenuParams: map[string]any{
			"weaponAction":     slotAction("WEAPON"),
			"helmetAction":     slotAction("HELMET"),
			"bootsAction":      slotAction("BOOTS"),
			"armorAction":      slotAction("ARMOR"),
			"glovesAction":     slotAction("GLOVES"),
			"ringAction":       slotAction("RING"),
			"amuletAction":     slotAction("AMULET"),
			"unequipAllAction": "unequip_all",
		},
	}, nil
}

func (s *InventoryService) UnequipItem(playerID int64, rawSlot string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	slot, ok := normalizeEquipmentSlot(rawSlot)
	if !ok {
		return view.GameResponse{
			Text:    "Slot de equipamento invalido.",
			MenuKey: "INVENTORY_EQUIPPED_MENU",
		}, nil
	}

	playerEquipment, err := s.InventoryRepository.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	equippedInventoryItemID := playerEquipment.GetSlotInventoryItemID(slot)
	if equippedInventoryItemID == nil || *equippedInventoryItemID == 0 {
		resp, viewErr := s.ViewEquippedItems(playerID)
		if viewErr != nil {
			return view.GameResponse{}, viewErr
		}
		resp.Text = "Nenhum item equipado neste slot.\n\n" + resp.Text
		return resp, nil
	}

	invItem, err := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), *equippedInventoryItemID)
	if err != nil {
		return view.GameResponse{}, err
	}
	if invItem == nil {
		resp, viewErr := s.ViewEquippedItems(playerID)
		if viewErr != nil {
			return view.GameResponse{}, viewErr
		}
		resp.Text = "Item equipado nao encontrado no inventario.\n\n" + resp.Text
		return resp, nil
	}

	item := invItem.Item
	if item.ID == 0 && invItem.ItemID > 0 {
		dbItem, itemErr := s.ItemsRepository.GetItemById(int(invItem.ItemID))
		if itemErr != nil {
			return view.GameResponse{}, itemErr
		}
		if dbItem != nil {
			item = *dbItem
		}
	}

	player.RemoveEquipmentBonus(&item)

	switch slot {
	case domain.SlotRing, domain.SlotAccessory:
		playerEquipment.RingInventoryItemID = nil
		playerEquipment.AccessoryInventoryItemID = nil
	default:
		if !playerEquipment.SetSlotInventoryItemID(slot, nil) {
			return view.GameResponse{}, fmt.Errorf("nao foi possivel desequipar no slot %s", slot)
		}
	}

	if err := s.InventoryRepository.SavePlayerAndEquipment(player, playerEquipment); err != nil {
		return view.GameResponse{}, err
	}

	resp, viewErr := s.ViewEquippedItems(playerID)
	if viewErr != nil {
		return view.GameResponse{}, viewErr
	}
	resp.Text = fmt.Sprintf("%s desequipado com sucesso.\n\n%s", item.Name, resp.Text)
	return resp, nil
}

func (s *InventoryService) UnequipAll(playerID int64) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}

	playerEquipment, err := s.InventoryRepository.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	seen := map[uint]struct{}{}
	equippedIDs := make([]uint, 0, 8)
	for _, ptr := range []*uint{
		playerEquipment.WeaponInventoryItemID,
		playerEquipment.ArmorInventoryItemID,
		playerEquipment.GlovesInventoryItemID,
		playerEquipment.HelmetInventoryItemID,
		playerEquipment.BootsInventoryItemID,
		playerEquipment.RingInventoryItemID,
		playerEquipment.AmuletInventoryItemID,
		playerEquipment.AccessoryInventoryItemID,
	} {
		if ptr == nil || *ptr == 0 {
			continue
		}
		if _, exists := seen[*ptr]; exists {
			continue
		}
		seen[*ptr] = struct{}{}
		equippedIDs = append(equippedIDs, *ptr)
	}

	if len(equippedIDs) == 0 {
		resp, viewErr := s.ViewEquippedItems(playerID)
		if viewErr != nil {
			return view.GameResponse{}, viewErr
		}
		resp.Text = "Voce nao tem equipamentos para desequipar.\n\n" + resp.Text
		return resp, nil
	}

	for _, inventoryItemID := range equippedIDs {
		invItem, getErr := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), inventoryItemID)
		if getErr != nil {
			return view.GameResponse{}, getErr
		}
		if invItem == nil {
			continue
		}

		item := invItem.Item
		if item.ID == 0 && invItem.ItemID > 0 {
			dbItem, itemErr := s.ItemsRepository.GetItemById(int(invItem.ItemID))
			if itemErr != nil {
				return view.GameResponse{}, itemErr
			}
			if dbItem != nil {
				item = *dbItem
			}
		}
		player.RemoveEquipmentBonus(&item)
	}

	playerEquipment.WeaponInventoryItemID = nil
	playerEquipment.ArmorInventoryItemID = nil
	playerEquipment.GlovesInventoryItemID = nil
	playerEquipment.HelmetInventoryItemID = nil
	playerEquipment.BootsInventoryItemID = nil
	playerEquipment.RingInventoryItemID = nil
	playerEquipment.AmuletInventoryItemID = nil
	playerEquipment.AccessoryInventoryItemID = nil

	if err := s.InventoryRepository.SavePlayerAndEquipment(player, playerEquipment); err != nil {
		return view.GameResponse{}, err
	}

	resp, viewErr := s.ViewEquippedItems(playerID)
	if viewErr != nil {
		return view.GameResponse{}, viewErr
	}
	resp.Text = "Todos os equipamentos foram desequipados.\n\n" + resp.Text
	return resp, nil
}

func (s InventoryService) ListEquipmentItemsBySlot(playerID int64, rawSlot string, page int) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}
	if page < 1 {
		page = 1
	}

	slot, ok := parseEquipmentSlot(rawSlot)
	if !ok {
		return view.GameResponse{
			Text:    "Slot de equipamento invalido.",
			MenuKey: "INVENTORY_EQUIPMENT_TYPES_MENU",
		}, nil
	}

	items, total, err := s.InventoryRepository.ListEquipmentItemsByUserAndSlot(int64(player.ID), slot, page, 5)
	if err != nil {
		return view.GameResponse{}, err
	}
	totalPages := int(math.Ceil(float64(total) / float64(5)))
	if totalPages == 0 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
		items, total, err = s.InventoryRepository.ListEquipmentItemsByUserAndSlot(int64(player.ID), slot, page, 5)
		if err != nil {
			return view.GameResponse{}, err
		}
		totalPages = int(math.Ceil(float64(total) / float64(5)))
		if totalPages == 0 {
			totalPages = 1
		}
	}

	msg := fmt.Sprintf("Equipamentos do tipo %s:", equipmentSlotLabel(slot))
	img := &view.Image{
		Type: view.ImageAsset,
		Path: "assets/npcs/arkanis_inventario.png",
	}

	return view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: "INVENTORY_ITEMS_MENU",
		MenuParams: map[string]any{
			"Items":            items,
			"page":             page,
			"totalPages":       totalPages,
			"itemType":         string(domain.ItemEquipment),
			"paginationAction": fmt.Sprintf("open_inventory_equipaments_type?slot=%s", slot),
			"backAction":       "open_inventory_equipaments",
		},
	}, nil
}

func (s *InventoryService) UseItem(playerID int64, inventoryItemID uint) (view.GameResponse, error) {
	log.Printf("use_item start | playerID=%d | inventoryItemID=%d", playerID, inventoryItemID)
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		log.Printf("use_item get player error | playerID=%d | err=%v", playerID, err)
		return view.GameResponse{}, err
	}
	if player == nil {
		log.Printf("use_item player nil | playerID=%d", playerID)
		return view.GameResponse{}, errors.New("player nao encontrado")
	}
	log.Printf("use_item player loaded | playerID=%d | internalPlayerID=%d | level=%d | class=%s", playerID, player.ID, player.Level, player.Class.Name)
	invItem, err := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), inventoryItemID)
	if err != nil {
		log.Printf("use_item get inventory item error | internalPlayerID=%d | inventoryItemID=%d | err=%v", player.ID, inventoryItemID, err)
		return view.GameResponse{}, err
	}
	if invItem == nil {
		log.Printf("use_item inventory item nil | internalPlayerID=%d | inventoryItemID=%d", player.ID, inventoryItemID)
		return view.GameResponse{}, errors.New("item nao encontrado no inventario")
	}
	log.Printf("use_item inventory item loaded | inventoryItemID=%d | itemID=%d | itemType=%s | qty=%d | slot=%s", invItem.ID, invItem.ItemID, invItem.Item.Type, invItem.Quantity, invItem.Item.Slot)

	if invItem.Item.Type == domain.ItemEquipment {
		log.Printf("use_item dispatch equip | inventoryItemID=%d | itemID=%d", invItem.ID, invItem.ItemID)
		return s.EquipItem(player, invItem)
	}
	if invItem.Item.Type == "" && invItem.ItemID > 0 {
		log.Printf("use_item item preload missing type | inventoryItemID=%d | itemID=%d", invItem.ID, invItem.ItemID)
		item, itemErr := s.ItemsRepository.GetItemById(int(invItem.ItemID))
		if itemErr == nil && item != nil && item.Type == domain.ItemEquipment {
			invItem.Item = *item
			log.Printf("use_item dispatch equip after reload | inventoryItemID=%d | itemID=%d | slot=%s", invItem.ID, invItem.ItemID, invItem.Item.Slot)
			return s.EquipItem(player, invItem)
		}
		log.Printf("use_item item reload result | itemID=%d | found=%v | err=%v", invItem.ItemID, item != nil, itemErr)
	}

	if invItem.Quantity <= 0 {
		return view.GameResponse{
			Text:    "❌ Você não possui mais essa poção.",
			Image:   inventoryGuideImage(),
			MenuKey: "MAIN_MENU",
		}, nil
	}
	if strings.Contains(invItem.Code, "HP") {
		size := strings.Split(invItem.Code, "_")

		switch size[2] {
		case "SMALL":
			player.HP += int(float64(player.MaxHP) * 0.2)
		case "MEDIUM":
			player.HP += int(float64(player.MaxHP) * 0.5)
		case "LARGE":
			player.HP = player.MaxHP
		}

		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}
	}

	if strings.Contains(invItem.Code, "ENERGY") {
		size := strings.Split(invItem.Code, "_")

		switch size[2] {
		case "SMALL":
			player.Energy += int(float64(player.MaxEnergy) * 0.2)
		case "MEDIUM":
			player.Energy += int(float64(player.MaxEnergy) * 0.5)
		case "LARGE":
			player.Energy = player.MaxEnergy
		}

		if player.Energy > player.MaxEnergy {
			player.Energy = player.MaxEnergy
		}
	}

	invItem.Quantity--
	err = s.InventoryRepository.SavePlayerAndInventoryItem(player, invItem)
	if err != nil {
		return view.GameResponse{}, err
	}

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%splayer.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
	}

	return view.GameResponse{
		Text:    fmt.Sprintf("🧪 Você usou %s!", invItem.Item.Name),
		Image:   img,
		MenuKey: "MAIN_MENU",
	}, nil
}

func normalizeEquipmentSlot(raw string) (domain.EquipmentSlot, bool) {
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

func parseEquipmentSlot(raw string) (domain.EquipmentSlot, bool) {
	return normalizeEquipmentSlot(raw)
}

func parsePage(raw string) int {
	page, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func equipmentSlotLabel(slot domain.EquipmentSlot) string {
	switch slot {
	case domain.SlotWeapon:
		return "Armas"
	case domain.SlotArmor:
		return "Armaduras"
	case domain.SlotGloves:
		return "Luvas"
	case domain.SlotHelmet:
		return "Capacetes"
	case domain.SlotBoots:
		return "Botas"
	case domain.SlotRing:
		return "Aneis"
	case domain.SlotAmulet:
		return "Amuletos"
	case domain.SlotAccessory:
		return "Acessorios"
	default:
		return string(slot)
	}
}

func equipmentBonusOnly(item *domain.Item) string {
	if item == nil {
		return "sem bonus"
	}
	parts := make([]string, 0, 5)
	if item.AtkBonus != 0 {
		parts = append(parts, fmt.Sprintf("ATK %+d", item.AtkBonus))
	}
	if item.DefBonus != 0 {
		parts = append(parts, fmt.Sprintf("DEF %+d", item.DefBonus))
	}
	if item.SpeedBonus != 0 {
		parts = append(parts, fmt.Sprintf("Speed %+d", item.SpeedBonus))
	}
	if item.HPBonus != 0 {
		parts = append(parts, fmt.Sprintf("MaxHP %+d", item.HPBonus))
	}
	if item.EnergyBonus != 0 {
		parts = append(parts, fmt.Sprintf("Energy %+d", item.EnergyBonus))
	}
	if len(parts) == 0 {
		return "sem bonus"
	}
	return strings.Join(parts, " | ")
}

func (s *InventoryService) EquipItem(player *domain.Player, invItem *domain.InventoryItem) (view.GameResponse, error) {
	if player == nil {
		log.Printf("equip_item invalid player nil")
		return view.GameResponse{}, errors.New("player invalido")
	}
	if invItem == nil {
		log.Printf("equip_item invalid inventory item nil | playerID=%d", player.ID)
		return view.GameResponse{}, errors.New("item invalido")
	}
	log.Printf("equip_item start | playerID=%d | inventoryItemID=%d | itemID=%d", player.ID, invItem.ID, invItem.ItemID)

	if invItem.ItemID == 0 && invItem.Item.ID > 0 {
		invItem.ItemID = invItem.Item.ID
	}
	if invItem.ItemID == 0 {
		log.Printf("equip_item invalid itemID | inventoryItemID=%d", invItem.ID)
		return view.GameResponse{}, errors.New("item invalido")
	}
	if invItem.Quantity <= 0 {
		log.Printf("equip_item out of stock | inventoryItemID=%d | itemID=%d", invItem.ID, invItem.ItemID)
		return view.GameResponse{
			Text:    "Item nao encontrado no inventario.",
			Image:   inventoryGuideImage(),
			MenuKey: "MAIN_MENU",
		}, nil
	}

	item := invItem.Item
	if (item.ID == 0 || item.Slot == "" || item.Type == "") && invItem.ItemID > 0 {
		log.Printf("equip_item loading metadata | itemID=%d", invItem.ItemID)
		dbItem, itemErr := s.ItemsRepository.GetItemById(int(invItem.ItemID))
		if itemErr != nil {
			log.Printf("equip_item metadata load error | itemID=%d | err=%v", invItem.ItemID, itemErr)
			return view.GameResponse{}, itemErr
		}
		if dbItem == nil {
			log.Printf("equip_item metadata not found | itemID=%d", invItem.ItemID)
			return view.GameResponse{}, errors.New("item nao encontrado")
		}
		item = *dbItem
	}

	log.Printf("equip_item metadata ready | itemID=%d | type=%s | slot=%s | requiredLevel=%d | class=%s", item.ID, item.Type, item.Slot, item.RequiredLevel, player.Class.Name)
	if item.Type != domain.ItemEquipment {
		log.Printf("equip_item invalid type | itemID=%d | type=%s", item.ID, item.Type)
		return view.GameResponse{}, errors.New("item nao e equipamento")
	}
	if player.Level < item.RequiredLevel {
		log.Printf("equip_item blocked by level | playerLevel=%d | requiredLevel=%d", player.Level, item.RequiredLevel)
		return view.GameResponse{
			Text:    fmt.Sprintf("Voce precisa ser nivel %d para equipar este item.", item.RequiredLevel),
			Image:   inventoryGuideImage(),
			MenuKey: "MAIN_MENU",
		}, nil
	}
	if !item.AllowsClassByIDAndName(player.ClassID, player.Class.Name) {
		log.Printf("equip_item blocked by class | class=%s | allowed=%q", player.Class.Name, item.AllowedClasses)
		return view.GameResponse{
			Text:    fmt.Sprintf("%s nao pode equipar este item.", player.Class.Name),
			Image:   inventoryGuideImage(),
			MenuKey: "MAIN_MENU",
		}, nil
	}

	slot, ok := normalizeEquipmentSlot(item.Slot)
	if !ok {
		log.Printf("equip_item invalid slot | itemID=%d | slot=%q", item.ID, item.Slot)
		return view.GameResponse{}, fmt.Errorf("slot invalido para item: %s", item.Slot)
	}
	log.Printf("equip_item slot resolved | itemID=%d | slot=%s", item.ID, slot)

	playerEquipment, err := s.InventoryRepository.GetOrCreatePlayerEquipment(int64(player.ID))
	if err != nil {
		log.Printf("equip_item get/create player_equipment error | playerID=%d | err=%v", player.ID, err)
		return view.GameResponse{}, err
	}
	log.Printf(
		"equip_item player_equipment loaded | weapon=%v | armor=%v | gloves=%v | helmet=%v | boots=%v | ring=%v | amulet=%v | accessory_legacy=%v",
		playerEquipment.WeaponInventoryItemID,
		playerEquipment.ArmorInventoryItemID,
		playerEquipment.GlovesInventoryItemID,
		playerEquipment.HelmetInventoryItemID,
		playerEquipment.BootsInventoryItemID,
		playerEquipment.RingInventoryItemID,
		playerEquipment.AmuletInventoryItemID,
		playerEquipment.AccessoryInventoryItemID,
	)

	if equippedInventoryItemID := playerEquipment.GetSlotInventoryItemID(slot); equippedInventoryItemID != nil && *equippedInventoryItemID > 0 {
		log.Printf("equip_item replacing existing slot item | slot=%s | oldInventoryItemID=%d | newInventoryItemID=%d", slot, *equippedInventoryItemID, invItem.ID)
		oldInvItem, oldInvErr := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), *equippedInventoryItemID)
		if oldInvErr != nil {
			log.Printf("equip_item old inventory item load error | oldInventoryItemID=%d | err=%v", *equippedInventoryItemID, oldInvErr)
			return view.GameResponse{}, oldInvErr
		}
		if oldInvItem != nil {
			oldItem := oldInvItem.Item
			if oldItem.ID == 0 && oldInvItem.ItemID > 0 {
				dbItem, oldErr := s.ItemsRepository.GetItemById(int(oldInvItem.ItemID))
				if oldErr != nil {
					log.Printf("equip_item old item metadata load error | oldInventoryItemID=%d | itemID=%d | err=%v", *equippedInventoryItemID, oldInvItem.ItemID, oldErr)
					return view.GameResponse{}, oldErr
				}
				if dbItem != nil {
					oldItem = *dbItem
				}
			}
			player.RemoveEquipmentBonus(&oldItem)
		}
	}

	player.ApplyEquipmentBonus(&item)
	log.Printf("equip_item bonus applied | playerID=%d | atk=%d | def=%d | speed=%d | maxHP=%d | maxEnergy=%d", player.ID, player.ATK, player.DEF, player.Speed, player.MaxHP, player.MaxEnergy)

	newInventoryItemID := invItem.ID
	if !playerEquipment.SetSlotInventoryItemID(slot, &newInventoryItemID) {
		log.Printf("equip_item set slot failed | slot=%s | newInventoryItemID=%d", slot, newInventoryItemID)
		return view.GameResponse{}, fmt.Errorf("nao foi possivel equipar no slot %s", slot)
	}

	if err := s.InventoryRepository.SavePlayerAndEquipment(player, playerEquipment); err != nil {
		log.Printf("equip_item save player+equipment error | playerID=%d | err=%v", player.ID, err)
		return view.GameResponse{}, err
	}
	log.Printf("equip_item success | playerID=%d | inventoryItemID=%d | itemID=%d | slot=%s", player.ID, invItem.ID, item.ID, slot)

	bonusText := fmt.Sprintf(
		"ATK %+d | DEF %+d | Speed %+d | MaxHP %+d | Energy %+d",
		item.AtkBonus,
		item.DefBonus,
		item.SpeedBonus,
		item.HPBonus,
		item.EnergyBonus,
	)
	msg := fmt.Sprintf(
		"%s equipada!\n\n%s\n\nAtributos atuais:\nATK %d | DEF %d | Speed %d | MaxHP %d",
		item.Name,
		bonusText,
		player.ATK,
		player.DEF,
		player.Speed,
		player.MaxHP,
	)

	return view.GameResponse{
		Text:       msg,
		Image:      &view.Image{Type: view.ImageAsset, Path: fmt.Sprintf("assets/skins/%s/%splayer.png", domain.ConvertLowCase(player.Gender), player.Class.Skin)},
		MenuKey:    "MAIN_MENU",
		MenuParams: map[string]any{"stats_updated": map[string]int{"atk": player.ATK, "def": player.DEF, "speed": player.Speed, "max_hp": player.MaxHP}},
	}, nil
}

func inventoryGuideImage() *view.Image {
	return &view.Image{
		Type: view.ImageAsset,
		Path: "assets/npcs/arkanis_inventario.png",
	}
}
