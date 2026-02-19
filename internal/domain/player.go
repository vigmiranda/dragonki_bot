package domain

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	UserID int64 `gorm:"uniqueIndex;not null"`

	Name      string
	Gender    string
	ClassID   uint
	ClassName string

	Level int
	XP    int

	HP    int
	MaxHP int
	ATK   int
	DEF   int

	Energy    int
	MaxEnergy int

	Gold       int
	KiCrystals int

	RegionID int64

	LastEnergyAt time.Time

	Form            string
	FormBattlesLeft int

	OnboardingStep string
	LastUpdate     time.Time

	UserIDTelegram int64
	NameTelegram   string
	UserName       string
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

func KiAwakened(player *Player) *Player {
	atk := int(float64(player.ATK) * 0.1)
	def := int(float64(player.DEF) * 0.1)
	if atk < 1 {
		atk = 1
	}
	if def < 1 {
		def = 1
	}
	player.ATK += atk
	player.DEF += def
	return player
}

func ConvertLowCase(text string) string {
	return strings.ToLower(text)
}

func (p *Player) RegenEnergy() {
	now := time.Now().UTC()

	// 🔹 fallback seguro
	recoveryStr := os.Getenv("ENERGY_REGEN_TIME")
	if recoveryStr == "" {
		recoveryStr = "10m"
	}

	// 🔹 se vier "10", assume minutos
	if _, err := strconv.Atoi(recoveryStr); err == nil {
		recoveryStr = recoveryStr + "m"
	}

	timeRecovery, err := time.ParseDuration(recoveryStr)
	if err != nil {
		log.Println("⚠️ ENERGY_REGEN_TIME inválido, usando default 10m:", err)
		timeRecovery = 10 * time.Minute
	}

	elapsed := now.Sub(p.LastEnergyAt)
	energyRecovered := int(elapsed / timeRecovery)

	if energyRecovered <= 0 {
		return
	}

	// 🔹 energia
	oldEnergy := p.Energy
	p.Energy += energyRecovered
	if p.Energy > p.MaxEnergy {
		p.Energy = p.MaxEnergy
	}

	// 🔹 HP só regenera se energia cheia
	if oldEnergy == p.MaxEnergy && p.HP < p.MaxHP {
		p.HP += energyRecovered
		if p.HP > p.MaxHP {
			p.HP = p.MaxHP
		}
	}

	// 🔹 avança apenas o tempo realmente consumido
	p.LastEnergyAt = p.LastEnergyAt.Add(
		time.Duration(energyRecovered) * timeRecovery,
	)
}
