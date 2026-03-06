package application

import (
	"fmt"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

func (s *BattleService) HandleDrops(monsterID int, playerID int64) (string, int, error) {
	drops, err := s.monsterRepo.GetMonsterDrops(monsterID)
	if err != nil {
		return "", 0, err
	}

	itemName := ""
	itemQty := 0

	for _, drop := range drops {
		roll := s.rng.Float64()

		if roll <= drop.DropRate {
			qty := drop.MinQty
			if drop.MaxQty > drop.MinQty {
				qty = s.rng.Intn(drop.MaxQty-drop.MinQty+1) + drop.MinQty
			}

			if err := s.inventoryRepository.CreateOrUpdateInventory(playerID, drop.ItemID, qty); err != nil {
				return "", 0, err
			}

			// Mantive o comportamento: se dropar mais de um item na mesma luta, devolve o ultimo que caiu.
			itemQty = qty
			itemName = drop.Item.Name
		}
	}

	return itemName, itemQty, nil
}

func (s *BattleService) ChangeClass(player *domain.Player) (string, *view.Image, string, map[string]any, error) {
	classes, err := s.playerRepo.GetClassesByLvl(player.Level)
	if err != nil {
		return "", nil, "", nil, err
	}

	menu := menuChoiceCls
	msg := fmt.Sprintf(
		"“Você chegou até aqui… e agora seu verdadeiro caminho começa.” \n\n" +
			"No nível 15, todo guerreiro deve decidir como sua energia irá se manifestar. \n" +
			"O Ki não é apenas poder… \n" +
			"É direção. É propósito. É identidade.\n" +
			"Você poderá trilhar um dos três caminhos: \n\n" +
			"🌪️<b> Wind Striker </b>— Velocidade, esquiva e ataques cortantes como o vento.\n" +
			"🔥<b> Ki Vanguarda </b>— Ofensiva intensa, explosões de energia e domínio do combate.\n" +
			"🛡️<b> Earth Sentinel </b>— Resistência absoluta, defesa impenetrável e força inabalável.\n\n" +
			"Escolha com sabedoria… \n" +
			"Pois após esta decisão, seu Ki jamais será o mesmo.\n\n" +
			"Qual será o seu legado em Aethoria?",
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: "assets/npcs/mestre_kalunga.png",
	}
	menuParams := map[string]any{
		"Classes": classes,
	}

	return msg, img, menu, menuParams, nil
}
