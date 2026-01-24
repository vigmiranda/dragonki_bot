package domain

import "time"

const (
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

	LastUpdate time.Time
}

func (p *Player) RegenEnergy() {
	elapsed := time.Now().Sub(p.LastEnergyAt)
	energy := int(elapsed / EnergyRegenTime)
	p.Energy += energy
	if p.Energy > p.MaxEnergy {
		p.Energy = p.MaxEnergy
	}
	p.LastEnergyAt = time.Now()
}

func (p *Player) CanTransform() bool {
	return p.Level >= 3 && p.Energy >= 3 && p.Form == "Base"
}

func (p *Player) Transform() {
	p.Form = "Ki Awakened"
	p.FormBattlesLeft = 3
	p.Energy -= 3
}
