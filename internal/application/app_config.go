package application

import (
	"os"
	"sync"
)

type appConfig struct {
	EnergyRegenTime string
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
	})
	return cfg
}
