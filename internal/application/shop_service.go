package application

import (
	"fmt"
	"log"
	"strconv"

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
			"Seja bem-vindo %s à minha lojinha! \n\n"+
			"Eu sou Kiara, mercadora viajante e guardiã de artefatos raros. \n\n"+
			"Aqui você encontrará poções, equipamentos e itens especiais que vão fortalecer seu corpo, "+
			"renovar seu Ki e te ajudar a sobreviver aos desafios que te esperam lá fora. \n\n"+
			"💎 Cada item aqui foi escolhido para quem busca evoluir, não apenas lutar. \n"+
			"Prepare-se bem… a jornada de um verdadeiro guerreiro (a) começa com as escolhas certas. \n\n"+
			"O que posso fazer por você hoje? 😄 \n\n",
		player.Name)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/kiara_mercadora.png"),
		},
		MenuKey: "SHOP_MENU",
	}
	return viewResponse, nil
}

func (s *ShopService) ListPotions() (view.GameResponse, error) {
	potions, err := s.ItemRepository.GetRandomList(domain.ItemConsumable, domain.Common)
	if err != nil {
		log.Fatal(err)
	}

	menuParams := make(map[string]any)
	for _, potion := range potions {

		label := fmt.Sprintf(
			"%d x %s - %d Gold",
			potion.Amount,
			potion.Name,
			potion.Price,
		)

		menuParams[label] = strconv.Itoa(int(potion.ID))
	}

	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n" +
			"Certo estas são as poções que tenho disponíveis hoje!!  \n\n" +
			"Qual delas você vai querer?   \n\n",
	)

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/kiara_mercadora.png"),
		},
		MenuKey:    "SHOP_ITEMS_MENU",
		MenuParams: menuParams,
	}

	return viewResponse, nil
}

func (s *ShopService) ListEquipments() (view.GameResponse, error) {
	equipments, err := s.ItemRepository.GetRandomList(domain.ItemEquipment, domain.Common)
	if err != nil {
		log.Fatal(err)
	}

	menuParams := make(map[string]any)
	for _, equipament := range equipments {

		label := fmt.Sprintf(
			"%d x %s - %d Gold",
			equipament.Amount,
			equipament.Name,
			equipament.Price,
		)

		menuParams[label] = strconv.Itoa(int(equipament.ID))
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Sprintf("->", equipments)
	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n" +
			"Ainda não temos nenhum equipamento disponível aqui na lojinha!!   \n\n")

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/kiara_mercadora.png"),
		},
		MenuKey: "SHOP_ITEMS_MENU",
	}
	return viewResponse, nil
}
func (s *ShopService) BuyItem(userID int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.PlayerRepository.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	itemIDStr, _ := params["itemID"]
	itemID, err := strconv.ParseInt(itemIDStr, 10, 64)
	item, err := s.ItemRepository.GetItemById(int(itemID))
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf(
		"🪙 Kiara, a Mercadora 🪙 \n\n" +
			"Você não tem gold suficiente para este item, junte mais um pouco e retorne quando tiver!!\n\n")

	if player.Gold >= item.Price {
		err = s.InventoryRepository.CreateOrUpdateInventory(int64(player.ID), item.ID, item.Amount)
		if err != nil {
			return view.GameResponse{}, err
		}
		player.Gold -= item.Price
		err = s.PlayerRepository.Update(player)
		if err != nil {
			return view.GameResponse{}, err
		}
		msg = fmt.Sprintf(
			"🪙 Kiara, a Mercadora 🪙 \n\n" +
				"Obrigada por sua compara, caso precise pode retornar aqui!!   \n\n")
	}

	viewResponse := view.GameResponse{
		Text: msg,
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/npcs/kiara_mercadora.png"),
		},
		MenuKey: "SHOP_MENU",
	}

	return viewResponse, nil
}
