package domain

import "time"

const (
	MaxEnergy       = 10
	EnergyRegenTime = 10 * time.Minute
)

type Player struct {
	UserID int64

	// Identidade
	Name   string
	Gender string
	Class  string

	// Progressão
	Level int
	XP    int

	// Combate
	HP    int
	MaxHP int

	ATK int
	DEF int

	// Energia
	Energy    int
	MaxEnergy int

	// Economia
	Gold       int
	KiCrystals int

	// Mundo
	RegionID int64

	// Estados
	LastEnergyAt    time.Time
	Form            string
	FormBattlesLeft int

	OnboardingStep string
}

func (p *Player) RegenEnergy() {
	now := time.Now()
	regenPoints := int(now.Sub(p.LastEnergyAt) / EnergyRegenTime)
	if regenPoints > 0 {
		p.Energy += regenPoints
		if p.Energy > MaxEnergy {
			p.Energy = MaxEnergy
		}
		p.LastEnergyAt = p.LastEnergyAt.Add(time.Duration(regenPoints) * EnergyRegenTime)
	}
}

func (p *Player) CanTransform() bool {
	return p.Level >= 3 && p.Energy >= 3 && p.Form == "Base"
}

func (p *Player) Transform() {
	p.Form = "Ki Awakened"
	p.FormBattlesLeft = 3
	p.Energy -= 3
}
