package domain

import (
	"log"
	"os"
	"time"
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
	now := time.Unix(time.Now().Unix(), 0)
	elapsed := now.Sub(p.LastEnergyAt.UTC())
	timeRecovery, err := time.ParseDuration(os.Getenv("ENERGY_REGEN_TIME"))
	if err != nil {
		log.Fatal(err)
	}
	energy := int(elapsed / timeRecovery)
	currentEnergy := p.Energy
	p.Energy += energy
	if p.Energy > p.MaxEnergy {
		p.Energy = p.MaxEnergy
	}
	if currentEnergy == p.MaxEnergy && p.HP < p.MaxHP {
		p.HP += energy
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
	}
	p.LastEnergyAt = time.Now().UTC()
}

func (p *Player) CanTransform() bool {
	return p.Level >= 3 && p.Energy >= 3 && p.Form == "Base"
}

func (p *Player) Transform() {
	p.Form = "Ki Awakened"
	p.FormBattlesLeft = 3
	p.Energy -= 3
}
