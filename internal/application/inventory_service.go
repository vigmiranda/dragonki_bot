package application

import (
	"fmt"
	"log"
	"math"

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
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf(
		"Inventory.........\n\n",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
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

	msg := fmt.Sprintf(
		"Inventory.........\n\n",
	)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
	}
	menu := "INVENTORY_MENU"

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	items, total, err := s.InventoryRepository.ListItemsByUserAndType(int64(player.ID), itemType, 1, 5)
	totalPages := int(math.Ceil(float64(total) / float64(5)))
	fmt.Printf("items ->", items, totalPages)

	return viewResponse, nil
}
