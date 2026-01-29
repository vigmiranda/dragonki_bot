package domain

import (
	"time"
)

type Player struct {
	UserID int64

	// Identidade
	Name      string
	Gender    string
	ClassID   int64
	ClassName string

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

func (p *Player) SetTransformBase() *Player {
	if p.FormBattlesLeft == 0 {
		p.Form = "Base"
	}
	return p
}

func (p *Player) CanTransform() bool {
	return p.Level >= 3 && p.Energy >= 3 && p.Form == "Base"
}

func (p *Player) Transform() {
	p.Form = "Ki Awakened"
	p.FormBattlesLeft = 3
	p.Energy -= 3
}
