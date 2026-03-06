package application

import (
	"os"
	"strconv"
	"sync"
	"time"
)

type appConfig struct {
	EnergyRegenTime     string
	EnergyRegenDuration time.Duration
}

var (
	cfgOnce sync.Once
	cfg     appConfig
)

func getAppConfig() appConfig {
	cfgOnce.Do(func() {
		cfg = appConfig{
			EnergyRegenTime: os.Getenv("ENERGY_REGEN_TIME"),
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
	return getAppConfig()
}
