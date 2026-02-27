package domain

type ClassTier string

const (
	TierBase     ClassTier = "BASE"
	TierAdvanced ClassTier = "ADVANCED"
	TierElite    ClassTier = "ELITE"
	TierMaster   ClassTier = "MASTER"
)

type Class struct {
	ID         int64
	Name       string
	BaseClass  string
	Tier       ClassTier
	MinLevel   int
	AtkBonus   float64 `gorm:"default:0.2"`
	DefBonus   float64 `gorm:"default:0.2"`
	SpeedBonus float64 `gorm:"default:0"`
	HpBonus    float64 `gorm:"default:0.1"`
	Skin       string  `gorm:"default:basic"`
}
