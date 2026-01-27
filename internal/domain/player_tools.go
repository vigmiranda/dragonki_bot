package domain

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

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

func BuildNumericPaymentID(userID int64) string {
	timestamp := time.Now().UnixMilli()
	return fmt.Sprintf("%d%d", userID, timestamp)
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
