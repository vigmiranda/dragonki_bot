package application

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type appConfig struct {
	EnergyRegenTime                    string
	EnergyRegenDuration                time.Duration
	LevelMult                          int
	MaxPlayerLevel                     int
	KiCrystalsValue                    int
	KiCrystalsMult                     int
	WeeklyBossDamageMultAbove20Pct     float64
	WeeklyBossDamageMultAtOrBelow20Pct float64
}

var (
	cfgOnce sync.Once
	cfg     appConfig
)

func getAppConfig() appConfig {
	cfgOnce.Do(func() {
		cfg = appConfig{
			EnergyRegenTime:                    os.Getenv("ENERGY_REGEN_TIME"),
			LevelMult:                          1,
			MaxPlayerLevel:                     50,
			KiCrystalsValue:                    1,
			KiCrystalsMult:                     1,
			WeeklyBossDamageMultAbove20Pct:     5,
			WeeklyBossDamageMultAtOrBelow20Pct: 2,
		}
		if cfg.EnergyRegenTime == "" {
			cfg.EnergyRegenTime = "10m"
		}
		if _, err := strconv.Atoi(cfg.EnergyRegenTime); err == nil {
			cfg.EnergyRegenTime += "m"
		}
		cfg.EnergyRegenDuration, _ = time.ParseDuration(cfg.EnergyRegenTime)
		if cfg.EnergyRegenDuration <= 0 {
			cfg.EnergyRegenTime = "10m"
			cfg.EnergyRegenDuration = 10 * time.Minute
		}
		if v := os.Getenv("LEVEL_MULT"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				cfg.LevelMult = parsed
			}
		}
		if v := os.Getenv("KI_CRYSTALS_VALUE"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				cfg.KiCrystalsValue = parsed
			}
		}
		if v := os.Getenv("KI_CRYSTALS_MULT"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				cfg.KiCrystalsMult = parsed
			}
		}
		if v := os.Getenv("MAX_PLAYER_LEVEL"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				cfg.MaxPlayerLevel = parsed
			}
		}
		if v := os.Getenv("BIG_BOSS_DAMAGE_MULT_ABOVE_20_PCT"); v != "" {
			if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
				cfg.WeeklyBossDamageMultAbove20Pct = parsed
			}
		}
		if v := os.Getenv("BIG_BOSS_DAMAGE_MULT_AT_OR_BELOW_20_PCT"); v != "" {
			if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
				cfg.WeeklyBossDamageMultAtOrBelow20Pct = parsed
			}
		}
	})
	return cfg
}

func resolveAppConfig(cfg appConfig) appConfig {
	if cfg.EnergyRegenDuration > 0 && cfg.EnergyRegenTime != "" {
		return cfg
	}
	if cfg.EnergyRegenTime != "" {
		if _, err := strconv.Atoi(cfg.EnergyRegenTime); err == nil {
			cfg.EnergyRegenTime += "m"
		}
		parsed, err := time.ParseDuration(cfg.EnergyRegenTime)
		if err == nil && parsed > 0 {
			cfg.EnergyRegenDuration = parsed
			return cfg
		}
	}
	base := getAppConfig()
	if cfg.LevelMult > 0 {
		base.LevelMult = cfg.LevelMult
	}
	if cfg.MaxPlayerLevel > 0 {
		base.MaxPlayerLevel = cfg.MaxPlayerLevel
	}
	if cfg.KiCrystalsValue > 0 {
		base.KiCrystalsValue = cfg.KiCrystalsValue
	}
	if cfg.KiCrystalsMult > 0 {
		base.KiCrystalsMult = cfg.KiCrystalsMult
	}
	if cfg.WeeklyBossDamageMultAbove20Pct > 0 {
		base.WeeklyBossDamageMultAbove20Pct = cfg.WeeklyBossDamageMultAbove20Pct
	}
	if cfg.WeeklyBossDamageMultAtOrBelow20Pct > 0 {
		base.WeeklyBossDamageMultAtOrBelow20Pct = cfg.WeeklyBossDamageMultAtOrBelow20Pct
	}
	return base
}
