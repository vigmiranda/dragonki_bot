package application

import (
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestCalculateEffectivePlayerStats_AppliesCumulativeSetBonuses(t *testing.T) {
	player := basePlayer()

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			weaponID := uint(1)
			helmetID := uint(2)
			armorID := uint(3)
			bootsID := uint(4)
			return &domain.PlayerEquipment{
				PlayerID:              playerID,
				WeaponInventoryItemID: &weaponID,
				HelmetInventoryItemID: &helmetID,
				ArmorInventoryItemID:  &armorID,
				BootsInventoryItemID:  &bootsID,
			}, nil
		},
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			items := map[uint]domain.Item{
				1: {ID: 101, SetID: 1, AtkBonus: 2},
				2: {ID: 102, SetID: 1, DefBonus: 3},
				3: {ID: 103, SetID: 1, HPBonus: 5},
				4: {ID: 104, SetID: 1, SpeedBonus: 1},
			}
			return &domain.InventoryItem{
				ID:     inventoryItemID,
				ItemID: items[inventoryItemID].ID,
				Item:   items[inventoryItemID],
			}, nil
		},
		getSetBonusesFn: func(setIDs []uint) ([]domain.EquipmentSetBonus, error) {
			return []domain.EquipmentSetBonus{
				{SetID: 1, PiecesRequired: 2, SpeedBonus: 4},
				{SetID: 1, PiecesRequired: 3, SpeedBonus: 6},
				{SetID: 1, PiecesRequired: 4, AtkBonus: 4, SpeedBonus: 8},
			}, nil
		},
	}

	stats, err := calculateEffectivePlayerStats(player, invRepo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.Bonus.ATK != 6 {
		t.Fatalf("expected atk bonus 6, got %d", stats.Bonus.ATK)
	}
	if stats.Bonus.DEF != 3 {
		t.Fatalf("expected def bonus 3, got %d", stats.Bonus.DEF)
	}
	if stats.Bonus.Speed != 19 {
		t.Fatalf("expected speed bonus 19, got %d", stats.Bonus.Speed)
	}
	if stats.Bonus.MaxHP != 5 {
		t.Fatalf("expected hp bonus 5, got %d", stats.Bonus.MaxHP)
	}
}

func TestCalculateEffectivePlayerStats_IgnoresSetBonusesWithoutEnoughPieces(t *testing.T) {
	player := basePlayer()

	invRepo := &fakeInventoryRepo{
		getOrCreateEquipFn: func(playerID int64) (*domain.PlayerEquipment, error) {
			weaponID := uint(1)
			helmetID := uint(2)
			return &domain.PlayerEquipment{
				PlayerID:              playerID,
				WeaponInventoryItemID: &weaponID,
				HelmetInventoryItemID: &helmetID,
			}, nil
		},
		getByUserFn: func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
			items := map[uint]domain.Item{
				1: {ID: 101, SetID: 1, AtkBonus: 2},
				2: {ID: 102, SetID: 2, DefBonus: 3},
			}
			return &domain.InventoryItem{
				ID:     inventoryItemID,
				ItemID: items[inventoryItemID].ID,
				Item:   items[inventoryItemID],
			}, nil
		},
		getSetBonusesFn: func(setIDs []uint) ([]domain.EquipmentSetBonus, error) {
			return []domain.EquipmentSetBonus{
				{SetID: 1, PiecesRequired: 2, SpeedBonus: 4},
				{SetID: 2, PiecesRequired: 2, DefBonus: 6},
			}, nil
		},
	}

	stats, err := calculateEffectivePlayerStats(player, invRepo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.Bonus.ATK != 2 || stats.Bonus.DEF != 3 || stats.Bonus.Speed != 0 {
		t.Fatalf("expected only item bonuses, got %+v", stats.Bonus)
	}
}
