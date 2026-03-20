package domain

type PlayerEquipment struct {
	ID       uint  `gorm:"primaryKey"`
	PlayerID int64 `gorm:"uniqueIndex"`

	WeaponInventoryItemID *uint `gorm:"column:weapon_inventory_item_id"`

	ArmorInventoryItemID *uint `gorm:"column:armor_inventory_item_id"`

	GlovesInventoryItemID *uint `gorm:"column:gloves_inventory_item_id"`

	HelmetInventoryItemID *uint `gorm:"column:helmet_inventory_item_id"`

	BootsInventoryItemID *uint `gorm:"column:boots_inventory_item_id"`

	RingInventoryItemID *uint `gorm:"column:ring_inventory_item_id"`

	AmuletInventoryItemID *uint `gorm:"column:amulet_inventory_item_id"`

	AccessoryInventoryItemID *uint `gorm:"column:accessory_inventory_item_id"`
}

func (PlayerEquipment) TableName() string {
	return "player_equipments"
}

func (pe *PlayerEquipment) GetSlotInventoryItemID(slot EquipmentSlot) *uint {
	switch slot {
	case SlotWeapon:
		return pe.WeaponInventoryItemID
	case SlotArmor:
		return pe.ArmorInventoryItemID
	case SlotGloves:
		return pe.GlovesInventoryItemID
	case SlotHelmet:
		return pe.HelmetInventoryItemID
	case SlotBoots:
		return pe.BootsInventoryItemID
	case SlotRing:
		return pe.RingInventoryItemID
	case SlotAmulet:
		return pe.AmuletInventoryItemID
	case SlotAccessory:
		return pe.AccessoryInventoryItemID
	default:
		return nil
	}
}

func (pe *PlayerEquipment) SetSlotInventoryItemID(slot EquipmentSlot, id *uint) bool {
	switch slot {
	case SlotWeapon:
		pe.WeaponInventoryItemID = id
	case SlotArmor:
		pe.ArmorInventoryItemID = id
	case SlotGloves:
		pe.GlovesInventoryItemID = id
	case SlotHelmet:
		pe.HelmetInventoryItemID = id
	case SlotBoots:
		pe.BootsInventoryItemID = id
	case SlotRing:
		pe.RingInventoryItemID = id
	case SlotAccessory:
		pe.AccessoryInventoryItemID = id
	case SlotAmulet:
		pe.AmuletInventoryItemID = id
	default:
		return false
	}
	return true
}
