package telegram

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeBattlePlayerRepo struct {
	player *domain.Player
}

func (f *fakeBattlePlayerRepo) GetOrCreate(userID int64) (*domain.Player, error) {
	return f.player, nil
}

func (f *fakeBattlePlayerRepo) Update(player *domain.Player) error {
	if player == nil {
		return nil
	}
	copyPlayer := *player
	f.player = &copyPlayer
	return nil
}

func (f *fakeBattlePlayerRepo) CountOnlineLast15Min() (int, error) { return 0, nil }

func (f *fakeBattlePlayerRepo) GetClassesByLvl(lvl int) ([]domain.Class, error) { return nil, nil }

func (f *fakeBattlePlayerRepo) GetClassesByID(id int64) (*domain.Class, error) { return nil, nil }

func (f *fakeBattlePlayerRepo) FindClassEvolutionByBaseClassAndLevel(baseClass string, lvl int) (*domain.Class, error) {
	return nil, nil
}

type fakeBattleInventoryRepo struct {
	item      *domain.InventoryItem
	equipment *domain.PlayerEquipment
	listItems []domain.InventoryItem
	listTotal int64
}

func (f *fakeBattleInventoryRepo) UpdateInventoryItem(item *domain.InventoryItem) error {
	if item == nil {
		return nil
	}
	copyItem := *item
	f.item = &copyItem
	return nil
}

func (f *fakeBattleInventoryRepo) CreateOrUpdateInventory(userID int64, itemID uint, amount int) error {
	return nil
}

func (f *fakeBattleInventoryRepo) PurchaseItem(playerID int64, itemID uint, amount int, price int) error {
	return nil
}

func (f *fakeBattleInventoryRepo) SellInventoryItem(playerID int64, inventoryItemID uint, goldValue int) (*domain.InventoryItem, error) {
	if f.item == nil || f.item.ID != inventoryItemID {
		return nil, nil
	}
	copyItem := *f.item
	return &copyItem, nil
}

func (f *fakeBattleInventoryRepo) ListItemsByUserAndType(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	return f.listItems, f.listTotal, nil
}

func (f *fakeBattleInventoryRepo) ListEquipmentItemsByUserAndSlot(playerID int64, slot domain.EquipmentSlot, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	return nil, 0, nil
}

func (f *fakeBattleInventoryRepo) CountEquipmentItemsBySlot(playerID int64) (map[domain.EquipmentSlot]int64, error) {
	return map[domain.EquipmentSlot]int64{}, nil
}

func (f *fakeBattleInventoryRepo) GetInventoryItemByUser(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
	if f.item == nil || f.item.ID != inventoryItemID {
		return nil, nil
	}
	copyItem := *f.item
	return &copyItem, nil
}

func (f *fakeBattleInventoryRepo) GetOrCreatePlayerEquipment(playerID int64) (*domain.PlayerEquipment, error) {
	if f.equipment == nil {
		f.equipment = &domain.PlayerEquipment{PlayerID: playerID}
	}
	return f.equipment, nil
}

func (f *fakeBattleInventoryRepo) UpdatePlayerEquipment(equipment *domain.PlayerEquipment) error {
	if equipment == nil {
		return nil
	}
	copyEquipment := *equipment
	f.equipment = &copyEquipment
	return nil
}

func (f *fakeBattleInventoryRepo) SavePlayerAndEquipment(player *domain.Player, equipment *domain.PlayerEquipment) error {
	_ = player
	return f.UpdatePlayerEquipment(equipment)
}

func (f *fakeBattleInventoryRepo) SavePlayerAndInventoryItem(player *domain.Player, item *domain.InventoryItem) error {
	_ = player
	return f.UpdateInventoryItem(item)
}

type fakeBattleMonsterRepo struct {
	monster *domain.Monster
}

func (f *fakeBattleMonsterRepo) ListByRegion(userID int64) ([]domain.Monster, error) { return nil, nil }

func (f *fakeBattleMonsterRepo) GetByMonsterID(monsterID int64) (*domain.Monster, error) {
	if f.monster == nil || int64(f.monster.MonsterID) != monsterID {
		return nil, nil
	}
	copyMonster := *f.monster
	return &copyMonster, nil
}

func (f *fakeBattleMonsterRepo) GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error) {
	return nil, nil
}

type fakeBattleRepo struct{}

func (f *fakeBattleRepo) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	return nil
}

func (f *fakeBattleRepo) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	return nil
}

func (f *fakeBattleRepo) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	return nil, nil
}

type fakeBattleItemsRepo struct{}

func (f *fakeBattleItemsRepo) GetRandomList(itemType domain.ItemType, rarity domain.ItemRarity) ([]domain.Item, error) {
	return nil, nil
}

func (f *fakeBattleItemsRepo) GetRandomEquipmentsBySlot(slot domain.EquipmentSlot, limit int) ([]domain.Item, error) {
	return nil, nil
}

func (f *fakeBattleItemsRepo) GetItemById(id int) (*domain.Item, error) { return nil, nil }

type fakeBattleRegionRepo struct{}

func (f *fakeBattleRegionRepo) GetByID(regionID int64) (*domain.Region, error) { return nil, nil }

func (f *fakeBattleRegionRepo) ListByLevel(level int64) ([]domain.Region, error) { return nil, nil }

func TestHandleBattle_UseItemMergesNoticeIntoDuelPreview(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	inventoryRepo := &fakeBattleInventoryRepo{
		item: &domain.InventoryItem{
			ID:       7,
			ItemID:   2,
			Code:     "ITEM_HP_SMALL",
			Quantity: 2,
			Item: domain.Item{
				ID:   2,
				Name: "Elixir de Batalha",
				Type: domain.ItemConsumable,
			},
		},
	}
	monsterRepo := &fakeBattleMonsterRepo{
		monster: &domain.Monster{
			MonsterID: 9,
			Name:      "Saibaman",
			Tag:       "saibaman",
			Level:     5,
			HP:        40,
			ATK:       9,
			DEF:       3,
			RegionID:  1,
		},
	}

	h := &Handler{
		battleService:    application.NewBattleService(playerRepo, monsterRepo, &fakeBattleRegionRepo{}, &fakeBattleRepo{}, inventoryRepo),
		inventoryService: application.NewInventoryService(inventoryRepo, playerRepo, &fakeBattleItemsRepo{}),
	}

	resp := h.handleBattle("action_battle_use_item", 1, map[string]string{
		"monsterID":       "9",
		"inventoryItemID": "7",
	})

	if !strings.Contains(resp.Message, "Elixir de Batalha") {
		t.Fatalf("expected merged duel message to include item usage notice, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Saibaman") {
		t.Fatalf("expected merged duel message to include duel preview, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "70/100") {
		t.Fatalf("expected duel preview to reflect healed HP, got %q", resp.Message)
	}
	if resp.CameFrom != "action_battle_use_item" {
		t.Fatalf("expected cameFrom action_battle_use_item, got %q", resp.CameFrom)
	}
}

func TestHandleBattle_TransformMergesNoticeIntoDuelPreview(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	playerRepo.player.Level = 8
	playerRepo.player.Energy = 4
	monsterRepo := &fakeBattleMonsterRepo{
		monster: &domain.Monster{
			MonsterID: 9,
			Name:      "Saibaman",
			Tag:       "saibaman",
			Level:     5,
			HP:        40,
			ATK:       9,
			DEF:       3,
			RegionID:  1,
		},
	}
	inventoryRepo := &fakeBattleInventoryRepo{}

	h := &Handler{
		battleService:    application.NewBattleService(playerRepo, monsterRepo, &fakeBattleRegionRepo{}, &fakeBattleRepo{}, inventoryRepo),
		transformService: application.NewTransformService(playerRepo),
	}

	resp := h.handleBattle("action_battle_transform", 1, map[string]string{
		"monsterID": "9",
	})

	if !strings.Contains(resp.Message, "Custo: 1 energia") {
		t.Fatalf("expected merged duel message to include transform notice, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Estado: Transformado") {
		t.Fatalf("expected merged duel message to include transformed duel preview, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Saibaman") {
		t.Fatalf("expected merged duel message to include duel preview, got %q", resp.Message)
	}
	if resp.CameFrom != "action_battle_transform" {
		t.Fatalf("expected cameFrom action_battle_transform, got %q", resp.CameFrom)
	}
}

func TestHandleBattle_PostBattleUseItemMergesNoticeIntoMenu(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	inventoryRepo := &fakeBattleInventoryRepo{
		item: &domain.InventoryItem{
			ID:       7,
			ItemID:   2,
			Code:     "ITEM_HP_SMALL",
			Quantity: 2,
			Item: domain.Item{
				ID:   2,
				Name: "Elixir de Batalha",
				Type: domain.ItemConsumable,
			},
		},
	}

	h := &Handler{
		inventoryService: application.NewInventoryService(inventoryRepo, playerRepo, &fakeBattleItemsRepo{}),
	}

	resp := h.handleBattle("action_post_battle_use_item", 1, map[string]string{
		"inventoryItemID": "7",
	})

	if !strings.Contains(resp.Message, "Elixir de Batalha") {
		t.Fatalf("expected merged post-battle message to include item usage notice, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Escolha sua proxima acao.") {
		t.Fatalf("expected merged post-battle message to include post-battle menu prompt, got %q", resp.Message)
	}
	if resp.CameFrom != "action_post_battle_use_item" {
		t.Fatalf("expected cameFrom action_post_battle_use_item, got %q", resp.CameFrom)
	}
}

func TestHandleBattle_PostBattleTransformMergesNoticeIntoMenu(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	playerRepo.player.Level = 8
	playerRepo.player.Energy = 4

	h := &Handler{
		transformService: application.NewTransformService(playerRepo),
	}

	resp := h.handleBattle("action_post_battle_transform", 1, map[string]string{})

	if !strings.Contains(resp.Message, "Custo: 1 energia") {
		t.Fatalf("expected merged post-battle message to include transform notice, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Escolha sua proxima acao.") {
		t.Fatalf("expected merged post-battle message to include post-battle menu prompt, got %q", resp.Message)
	}
	if resp.CameFrom != "action_post_battle_transform" {
		t.Fatalf("expected cameFrom action_post_battle_transform, got %q", resp.CameFrom)
	}
}

func TestHandleBattle_OpenPotionsUsesCombatSpecificShortText(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	inventoryRepo := &fakeBattleInventoryRepo{
		listItems: []domain.InventoryItem{
			{ID: 1, Quantity: 1, Item: domain.Item{Name: "Pocao"}},
		},
		listTotal: 6,
	}

	h := &Handler{
		inventoryService: application.NewInventoryService(inventoryRepo, playerRepo, &fakeBattleItemsRepo{}),
	}

	resp := h.handleBattle("action_battle_open_potions", 1, map[string]string{
		"monsterID": "9",
		"page":      "2",
	})

	if !strings.Contains(resp.Message, "Duelo em andamento.") {
		t.Fatalf("expected battle potion text to mention combat context, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Pagina 2/2.") {
		t.Fatalf("expected battle potion text to include pagination, got %q", resp.Message)
	}
	if strings.Contains(resp.Message, "Arkanis") {
		t.Fatalf("expected battle potion text to avoid generic inventory narration, got %q", resp.Message)
	}
}

func TestHandleBattle_PostBattleOpenPotionsUsesShortRecoveryText(t *testing.T) {
	playerRepo := &fakeBattlePlayerRepo{
		player: testBattlePlayer(1),
	}
	inventoryRepo := &fakeBattleInventoryRepo{
		listItems: []domain.InventoryItem{
			{ID: 1, Quantity: 1, Item: domain.Item{Name: "Pocao"}},
		},
		listTotal: 1,
	}

	h := &Handler{
		inventoryService: application.NewInventoryService(inventoryRepo, playerRepo, &fakeBattleItemsRepo{}),
	}

	resp := h.handleBattle("action_post_battle_open_potions", 1, map[string]string{
		"page": "1",
	})

	if !strings.Contains(resp.Message, "Recupere-se antes da proxima cacada.") {
		t.Fatalf("expected post-battle potion text to mention recovery context, got %q", resp.Message)
	}
	if !strings.Contains(resp.Message, "Pagina 1/1.") {
		t.Fatalf("expected post-battle potion text to include pagination, got %q", resp.Message)
	}
	if strings.Contains(resp.Message, "Arkanis") {
		t.Fatalf("expected post-battle potion text to avoid generic inventory narration, got %q", resp.Message)
	}
}

func testBattlePlayer(userID int64) *domain.Player {
	return &domain.Player{
		ID:           1,
		UserID:       userID,
		Name:         "Goku",
		Gender:       "man",
		ClassID:      2,
		Class:        domain.Class{Name: "Ki Vanguard", Skin: "sayajin"},
		Level:        10,
		HP:           50,
		MaxHP:        100,
		ATK:          20,
		DEF:          15,
		Speed:        12,
		Energy:       3,
		MaxEnergy:    5,
		RegionID:     1,
		Form:         "Base",
		LastEnergyAt: time.Now().UTC(),
	}
}
