package domain

type ClassTier string

const (
	TierBase     ClassTier = "BASE"
	TierAdvanced ClassTier = "ADVANCED"
	TierElite    ClassTier = "ELITE"
	TierMaster   ClassTier = "MASTER"
)

type Class struct {
	ID          int64
	Name        string
	BaseClass   string // Blade Monk, Sword Disciple etc
	Tier        ClassTier
	MinLevel    int
	AtkBonus    int
	HpBonus     int
	EnergyBonus int
}
