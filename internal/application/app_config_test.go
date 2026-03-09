package application

import (
	"sync"
	"testing"
	"time"
)

func resetAppConfigForTest() {
	cfgOnce = sync.Once{}
	cfg = appConfig{}
}

func TestGetAppConfig_Default(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "")

	got := getAppConfig()
	if got.EnergyRegenTime != "10m" {
		t.Fatalf("expected default 10m, got %q", got.EnergyRegenTime)
	}
	if got.EnergyRegenDuration != 10*time.Minute {
		t.Fatalf("expected default duration 10m, got %s", got.EnergyRegenDuration)
	}
	if got.LevelMult != 1 || got.MaxPlayerLevel != 50 || got.KiCrystalsValue != 1 || got.KiCrystalsMult != 1 {
		t.Fatalf("expected default numeric config values")
	}
}

func TestGetAppConfig_FromEnv(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "7m")

	got := getAppConfig()
	if got.EnergyRegenTime != "7m" {
		t.Fatalf("expected 7m, got %q", got.EnergyRegenTime)
	}
	if got.EnergyRegenDuration != 7*time.Minute {
		t.Fatalf("expected duration 7m, got %s", got.EnergyRegenDuration)
	}
}

func TestGetAppConfig_OnceBehavior(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "9m")
	first := getAppConfig()

	t.Setenv("ENERGY_REGEN_TIME", "1m")
	second := getAppConfig()

	if first.EnergyRegenTime != "9m" || second.EnergyRegenTime != "9m" {
		t.Fatalf("expected once-cached 9m, got first=%q second=%q", first.EnergyRegenTime, second.EnergyRegenTime)
	}
	if first.EnergyRegenDuration != 9*time.Minute || second.EnergyRegenDuration != 9*time.Minute {
		t.Fatalf("expected once-cached duration 9m")
	}
}

func TestGetAppConfig_NumericMinutesAndInvalidFallback(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "5")
	got := getAppConfig()
	if got.EnergyRegenTime != "5m" || got.EnergyRegenDuration != 5*time.Minute {
		t.Fatalf("expected 5m/5m, got %q/%s", got.EnergyRegenTime, got.EnergyRegenDuration)
	}

	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "invalid")
	got = getAppConfig()
	if got.EnergyRegenTime != "10m" || got.EnergyRegenDuration != 10*time.Minute {
		t.Fatalf("expected fallback 10m/10m, got %q/%s", got.EnergyRegenTime, got.EnergyRegenDuration)
	}
}

func TestGetAppConfig_ReadsNumericConfigFromEnv(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("LEVEL_MULT", "4")
	t.Setenv("MAX_PLAYER_LEVEL", "60")
	t.Setenv("KI_CRYSTALS_VALUE", "3")
	t.Setenv("KI_CRYSTALS_MULT", "2")

	got := getAppConfig()
	if got.LevelMult != 4 || got.MaxPlayerLevel != 60 || got.KiCrystalsValue != 3 || got.KiCrystalsMult != 2 {
		t.Fatalf("unexpected numeric config: %+v", got)
	}
}
