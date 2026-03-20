package application

import (
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type equipmentSetBonusLoader interface {
	GetEquipmentSetBonuses(setIDs []uint) ([]domain.EquipmentSetBonus, error)
}

type equipmentBonus struct {
	ATK       int
	DEF       int
	Speed     int
	MaxHP     int
	MaxEnergy int
}

type effectivePlayerStats struct {
	ATK       int
	DEF       int
	Speed     int
	MaxHP     int
	MaxEnergy int
	Bonus     equipmentBonus
}

func calculateEffectivePlayerStats(player *domain.Player, inventoryRepo ports.InventoryRepository) (effectivePlayerStats, error) {
	stats := effectivePlayerStats{
		ATK:       player.ATK,
		DEF:       player.DEF,
		Speed:     player.Speed,
		MaxHP:     player.MaxHP,
		MaxEnergy: player.MaxEnergy,
	}

	if inventoryRepo == nil {
		return stats, nil
	}

	bonus, err := loadEquippedBonus(inventoryRepo, int64(player.ID))
	if err != nil {
		return stats, err
	}

	stats.Bonus = bonus

	return stats, nil
}

func loadEquippedBonus(inventoryRepo ports.InventoryRepository, playerID int64) (equipmentBonus, error) {
	if inventoryRepo == nil {
		return equipmentBonus{}, nil
	}

	pe, err := inventoryRepo.GetOrCreatePlayerEquipment(playerID)
	if err != nil {
		return equipmentBonus{}, err
	}
	if pe == nil {
		return equipmentBonus{}, nil
	}

	seen := map[uint]struct{}{}
	itemIDs := make([]uint, 0, 8)
	for _, ptr := range []*uint{
		pe.WeaponInventoryItemID,
		pe.ArmorInventoryItemID,
		pe.GlovesInventoryItemID,
		pe.HelmetInventoryItemID,
		pe.BootsInventoryItemID,
		pe.RingInventoryItemID,
		pe.AmuletInventoryItemID,
		pe.AccessoryInventoryItemID,
	} {
		if ptr == nil || *ptr == 0 {
			continue
		}
		if _, exists := seen[*ptr]; exists {
			continue
		}
		seen[*ptr] = struct{}{}
		itemIDs = append(itemIDs, *ptr)
	}

	total := equipmentBonus{}
	setCounts := map[uint]int{}
	for _, inventoryItemID := range itemIDs {
		invItem, getErr := inventoryRepo.GetInventoryItemByUser(playerID, inventoryItemID)
		if getErr != nil {
			return equipmentBonus{}, getErr
		}
		if invItem == nil {
			continue
		}

		item := invItem.Item
		total.ATK += item.AtkBonus
		total.DEF += item.DefBonus
		total.Speed += item.SpeedBonus
		total.MaxHP += item.HPBonus
		total.MaxEnergy += item.EnergyBonus
		if item.SetID > 0 {
			setCounts[item.SetID]++
		}
	}

	if len(setCounts) == 0 {
		return total, nil
	}

	loader, ok := inventoryRepo.(equipmentSetBonusLoader)
	if !ok {
		return total, nil
	}

	setIDs := make([]uint, 0, len(setCounts))
	for setID := range setCounts {
		setIDs = append(setIDs, setID)
	}

	setBonuses, err := loader.GetEquipmentSetBonuses(setIDs)
	if err != nil {
		return equipmentBonus{}, err
	}
	for _, bonus := range setBonuses {
		if setCounts[bonus.SetID] < bonus.PiecesRequired {
			continue
		}
		total.ATK += bonus.AtkBonus
		total.DEF += bonus.DefBonus
		total.Speed += bonus.SpeedBonus
		total.MaxHP += bonus.HPBonus
	}

	return total, nil
}
