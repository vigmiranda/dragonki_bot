package domain

import "strings"

type ItemType string

const (
	ItemConsumable ItemType = "CONSUMABLE"
	ItemEquipment  ItemType = "EQUIPMENT"
	ItemKey        ItemType = "KEY"
)

type Item struct {
	ID uint `gorm:"primaryKey"`

	Code        string
	Name        string
	Description string
	Type        ItemType
	Rarity      string
	Amount      int
	Price       int

	Slot          string `gorm:"-"`
	RequiredLevel int    `gorm:"-"`
	AtkBonus      int    `gorm:"-"`
	DefBonus      int    `gorm:"-"`
	SpeedBonus    int    `gorm:"-"`
	HPBonus       int    `gorm:"-"`
	EnergyBonus   int    `gorm:"-"`

	// CSV de classes permitidas. Ex.: "Earth Sentinel,Wind Striker,Ki Vanguarda"
	AllowedClasses string `gorm:"-"`

	EquipmentID *uint
	Equipment   *Equipment `gorm:"foreignKey:EquipmentID"`
}

func (i *Item) HydrateEquipmentMetadata() {
	if i == nil || i.Type != ItemEquipment || i.Equipment == nil {
		return
	}

	i.Slot = i.Equipment.Slot
	i.RequiredLevel = i.Equipment.RequiredLevel
	i.AtkBonus = i.Equipment.AtkBonus
	i.DefBonus = i.Equipment.DefBonus
	i.SpeedBonus = i.Equipment.SpeedBonus
	i.HPBonus = i.Equipment.HPBonus
	i.EnergyBonus = i.Equipment.EnergyBonus
	i.AllowedClasses = i.Equipment.AllowedClasses
}

func (i Item) AllowsClass(className string) bool {
	return i.AllowsClassByIDAndName(0, className)
}

func (i Item) AllowsClassByIDAndName(classID uint, className string) bool {
	raw := strings.TrimSpace(i.AllowedClasses)
	if raw == "" {
		return true
	}

	targets := classTargets(classID, className)
	for _, entry := range strings.Split(raw, ",") {
		normalizedEntry := normalizeClassName(entry)
		if normalizedEntry == "" {
			continue
		}
		if containsString(targets, normalizedEntry) {
			return true
		}
		if numericClassID(normalizedEntry) == classID && classID > 0 {
			return true
		}
	}
	return false
}

func classTargets(classID uint, className string) []string {
	targetSet := map[string]struct{}{}
	add := func(v string) {
		if v == "" {
			return
		}
		targetSet[v] = struct{}{}
	}

	add(normalizeClassName(className))
	for _, alias := range classAliasesByID[classID] {
		add(normalizeClassName(alias))
	}

	out := make([]string, 0, len(targetSet))
	for k := range targetSet {
		out = append(out, k)
	}
	return out
}

var classAliasesByID = map[uint][]string{
	2: {"Ki Vanguard", "Ki Vanguarda"},
	3: {"Earth Sentinel"},
	4: {"Wind Striker", "Win Striker"},
}

func normalizeClassName(v string) string {
	n := strings.ToLower(strings.TrimSpace(v))
	switch n {
	case "win striker":
		return "wind striker"
	case "ki vanguarda":
		return "ki vanguard"
	default:
		return n
	}
}

func containsString(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}

func numericClassID(v string) uint {
	if v == "" {
		return 0
	}
	var n uint
	for _, ch := range v {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + uint(ch-'0')
	}
	return n
}
