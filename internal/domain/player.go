package domain

import "time"

const (
	MaxEnergy       = 10
	EnergyRegenTime = 10 * time.Second
)

type Player struct {
	UserID          int64
	Class           string
	HP              int
	ATK             int
	XP              int
	Level           int
	Energy          int
	LastEnergyAt    time.Time
	Form            string
	FormBattlesLeft int
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
