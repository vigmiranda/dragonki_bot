// internal/domain/equipment_slot.go
package domain

type EquipmentSlot string

const (
	SlotWeapon EquipmentSlot = "WEAPON"
	SlotArmor  EquipmentSlot = "ARMOR"
	SlotGloves EquipmentSlot = "GLOVES"
	SlotHelmet EquipmentSlot = "HELMET"
	SlotBoots  EquipmentSlot = "BOOTS"
	SlotRing   EquipmentSlot = "RING"
	SlotAmulet EquipmentSlot = "AMULET"

	// Slot legado para manter compatibilidade com dados antigos.
	SlotAccessory EquipmentSlot = "ACCESSORY"
)
