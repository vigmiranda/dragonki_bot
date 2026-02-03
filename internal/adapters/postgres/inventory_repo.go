package postgres

import (
	"errors"
	"fmt"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type InventoryPostgres struct {
	db *gorm.DB
}

func NewInventoryPostgresRepo(db *gorm.DB) *InventoryPostgres {
	return &InventoryPostgres{
		db: db,
	}
}

func (r *InventoryPostgres) CreateOrUpdateInventory(
	userID int64,
	itemID uint,
	amount int,
) error {

	// 🔹 1. Carrega o item (para saber o tipo)
	var item domain.Item
	if err := r.db.First(&item, itemID).Error; err != nil {
		return err
	}

	// 🔹 2. Busca ou cria o inventário do player
	var inventory domain.Inventory
	err := r.db.
		Where("player_id = ?", userID).
		First(&inventory).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		inventory = domain.Inventory{
			PlayerID: userID,
		}
		if err := r.db.Create(&inventory).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// 🔹 3. Lógica por tipo de item
	switch item.Type {

	// ==========================
	// 🧪 CONSUMABLES
	// ==========================
	case domain.ItemConsumable:

		var invItem domain.InventoryItem
		err := r.db.
			Where("inventory_id = ? AND code = ?", inventory.ID, item.Code).
			First(&invItem).
			Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ➕ Não existe → cria
			invItem = domain.InventoryItem{
				InventoryID: inventory.ID,
				ItemID:      item.ID,
				Code:        item.Code,
				Quantity:    amount,
			}
			return r.db.Create(&invItem).Error
		}

		if err != nil {
			return err
		}

		// ➕ Existe → soma quantidade
		invItem.Quantity += amount
		return r.db.Save(&invItem).Error

	// ==========================
	// 🗡️ EQUIPAMENTOS
	// ==========================
	case domain.ItemEquipment:

		invItem := domain.InventoryItem{
			InventoryID: inventory.ID,
			ItemID:      item.ID,
			Code:        item.Code,
			Quantity:    1, // sempre 1
		}

		return r.db.Create(&invItem).Error

	// ==========================
	// 🗝️ ITENS DE CHAVE / ESPECIAIS
	// ==========================
	case domain.ItemKey:

		invItem := domain.InventoryItem{
			InventoryID: inventory.ID,
			ItemID:      item.ID,
			Code:        item.Code,
			Quantity:    1,
		}

		return r.db.Create(&invItem).Error

	default:
		return fmt.Errorf("tipo de item inválido: %s", item.Type)
	}
}
