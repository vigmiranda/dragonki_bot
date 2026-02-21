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

func (r *InventoryPostgres) UpdateInventoryItem(item *domain.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *InventoryPostgres) ListItemsByUserAndType(
	playerID int64,
	itemType domain.ItemType,
	page int,
	pageSize int,
) ([]domain.InventoryItem, int64, error) {

	var inventory domain.Inventory
	err := r.db.
		Where("player_id = ?", playerID).
		First(&inventory).
		Error

	if err != nil {
		return nil, 0, err
	}

	var items []domain.InventoryItem
	var total int64

	offset := (page - 1) * pageSize

	// 🔢 Total de itens
	err = r.db.
		Model(&domain.InventoryItem{}).
		Where("inventory_id = ?", inventory.ID).
		Count(&total).
		Error

	if err != nil {
		return nil, 0, err
	}

	// 📦 Busca paginada + preload do item
	err = r.db.
		Model(&domain.InventoryItem{}).
		Joins("JOIN items ON items.id = inventory_items.item_id").
		Where("inventory_items.inventory_id = ?", inventory.ID).
		Where("items.type = ?", itemType).
		Where("inventory_items.quantity > 0 ").
		Preload("Item").
		Order("inventory_items.id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&items).
		Error

	if err != nil {
		return nil, 0, err
	}

	return items, total, nil

	// items, total, err := inventoryRepo.ListItemsByUser(userID, page, 5)
	// totalPages :=int(math.Ceil(float64(total)/float(pageSize)))
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
func (r *InventoryPostgres) GetInventoryItemByUser(
	playerID int64,
	inventoryItemID uint,
) (*domain.InventoryItem, error) {

	var invItem domain.InventoryItem

	err := r.db.
		Joins("JOIN inventories ON inventories.id = inventory_items.inventory_id").
		Preload("Item").
		Where("inventory_items.id = ?", inventoryItemID).
		Where("inventories.player_id = ?", playerID).
		First(&invItem).
		Error

	if err != nil {
		return nil, err
	}
	return &invItem, nil
}
