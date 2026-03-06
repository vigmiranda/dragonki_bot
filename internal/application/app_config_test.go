package application

import (
	"sync"
	"testing"
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
}

func TestGetAppConfig_FromEnv(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "7m")

	got := getAppConfig()
	if got.EnergyRegenTime != "7m" {
		t.Fatalf("expected 7m, got %q", got.EnergyRegenTime)
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
}
