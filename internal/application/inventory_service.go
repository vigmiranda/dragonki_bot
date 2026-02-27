package application

import (
	"fmt"
	"log"
	"math"
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
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		log.Fatal(err)
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
	menu := "POTION_MENU"
	page := 1

	items, total, err := s.InventoryRepository.ListItemsByUserAndType(int64(player.ID), itemType, page, 5)
	totalPages := int(math.Ceil(float64(total) / float64(5)))

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
		MenuParams: map[string]any{
			"Items":      items,
			"page":       page,
			"totalPages": totalPages,
		},
	}

	return viewResponse, nil
}

func (s *InventoryService) UseItem(playerID int64, inventoryItemID uint) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(playerID)
	if err != nil {
		return view.GameResponse{}, err
	}
	invItem, err := s.InventoryRepository.GetInventoryItemByUser(int64(player.ID), inventoryItemID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if invItem.Quantity <= 0 {
		return view.GameResponse{
			Text:    "❌ Você não possui mais essa poção.",
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
	err = s.InventoryRepository.UpdateInventoryItem(invItem)
	if err != nil {
		return view.GameResponse{}, err
	}
	err = s.PlayerRepository.Update(player)
	if err != nil {
		return view.GameResponse{}, err
	}

	return view.GameResponse{
		Text: fmt.Sprintf("🧪 Você usou %s!", invItem.Item.Name),
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/skins/%s_skin_%s.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
		},
		MenuKey: "MAIN_MENU",
	}, nil
}
