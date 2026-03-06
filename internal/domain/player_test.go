package domain

import (
	"testing"
	"time"
)

func TestPlayer_SetTransformBase(t *testing.T) {
	p := &Player{Form: "Ki Awakened", FormBattlesLeft: 0}
	got := p.SetTransformBase()
	if got != p {
		t.Fatalf("expected same player pointer")
	}
	if p.Form != "Base" {
		t.Fatalf("expected Base form, got %q", p.Form)
	}
}

func TestPlayer_SetTransformBase_KeepFormWhenBattlesLeft(t *testing.T) {
	p := &Player{Form: "Ki Awakened", FormBattlesLeft: 2}
	p.SetTransformBase()
	if p.Form != "Ki Awakened" {
		t.Fatalf("expected form unchanged, got %q", p.Form)
	}
}

func TestPlayer_CanTransform(t *testing.T) {
	p := &Player{Level: 3, Energy: 3, Form: "Base"}
	if !p.CanTransform() {
		t.Fatalf("expected can transform")
	}

	p.Level = 2
	if p.CanTransform() {
		t.Fatalf("expected cannot transform with low level")
	}

	p.Level = 3
	p.Energy = 2
	if p.CanTransform() {
		t.Fatalf("expected cannot transform with low energy")
	}

	p.Energy = 3
	p.Form = "Ki Awakened"
	if p.CanTransform() {
		t.Fatalf("expected cannot transform when not in Base form")
	}
}

func TestPlayer_Transform(t *testing.T) {
	p := &Player{Form: "Base", Energy: 5}
	p.Transform()
	if p.Form != "Ki Awakened" {
		t.Fatalf("expected Ki Awakened, got %q", p.Form)
	}
	if p.FormBattlesLeft != 3 {
		t.Fatalf("expected 3 battles left, got %d", p.FormBattlesLeft)
	}
	if p.Energy != 2 {
		t.Fatalf("expected energy 2, got %d", p.Energy)
	}
}

func TestKiAwakened_PercentBonus(t *testing.T) {
	p := &Player{ATK: 10, DEF: 20}
	got := KiAwakened(p)
	if got != p {
		t.Fatalf("expected same player pointer")
	}
	if p.ATK != 12 || p.DEF != 24 {
		t.Fatalf("expected ATK=12 DEF=24, got ATK=%d DEF=%d", p.ATK, p.DEF)
	}
}

func TestKiAwakened_MinimumBonusIsOne(t *testing.T) {
	p := &Player{ATK: 1, DEF: 1}
	KiAwakened(p)
	if p.ATK != 2 || p.DEF != 2 {
		t.Fatalf("expected ATK=2 DEF=2, got ATK=%d DEF=%d", p.ATK, p.DEF)
	}
}

func TestConvertLowCase(t *testing.T) {
	if got := ConvertLowCase("WOMAN"); got != "woman" {
		t.Fatalf("expected woman, got %q", got)
	}
}

func TestPlayer_RegenEnergy_NoRecoveryWhenElapsedTooSmall(t *testing.T) {
	last := time.Now().UTC().Add(-1 * time.Minute)
	p := &Player{
		Energy:       1,
		MaxEnergy:    10,
		HP:           5,
		MaxHP:        10,
		LastEnergyAt: last,
	}

	p.RegenEnergy(10 * time.Minute)
	if p.Energy != 1 || p.HP != 5 {
		t.Fatalf("expected no changes, got energy=%d hp=%d", p.Energy, p.HP)
	}
	if !p.LastEnergyAt.Equal(last) {
		t.Fatalf("expected LastEnergyAt unchanged")
	}
}

func TestPlayer_RegenEnergy_DefaultEnvAndTimeAdvance(t *testing.T) {
	last := time.Now().UTC().Add(-25 * time.Minute)
	p := &Player{
		Energy:       1,
		MaxEnergy:    10,
		HP:           5,
		MaxHP:        10,
		LastEnergyAt: last,
	}

	p.RegenEnergy(10 * time.Minute)
	if p.Energy != 3 {
		t.Fatalf("expected energy 3, got %d", p.Energy)
	}
	expectedLast := last.Add(20 * time.Minute)
	if !p.LastEnergyAt.Equal(expectedLast) {
		t.Fatalf("expected LastEnergyAt %v, got %v", expectedLast, p.LastEnergyAt)
	}
}

func TestPlayer_RegenEnergy_UsesConfiguredDuration(t *testing.T) {
	last := time.Now().UTC().Add(-11 * time.Minute)
	p := &Player{
		Energy:       0,
		MaxEnergy:    10,
		LastEnergyAt: last,
	}

	p.RegenEnergy(5 * time.Minute)
	if p.Energy != 2 {
		t.Fatalf("expected energy 2, got %d", p.Energy)
	}
}

func TestPlayer_RegenEnergy_IgnoreInvalidDuration(t *testing.T) {
	last := time.Now().UTC().Add(-21 * time.Minute)
	p := &Player{
		Energy:       0,
		MaxEnergy:    10,
		LastEnergyAt: last,
	}

	p.RegenEnergy(0)
	if p.Energy != 0 {
		t.Fatalf("expected no energy change for invalid duration, got %d", p.Energy)
	}
}

func TestPlayer_RegenEnergy_CapsEnergyAtMax(t *testing.T) {
	p := &Player{
		Energy:       9,
		MaxEnergy:    10,
		LastEnergyAt: time.Now().UTC().Add(-30 * time.Minute),
	}

	p.RegenEnergy(10 * time.Minute)
	if p.Energy != 10 {
		t.Fatalf("expected capped energy 10, got %d", p.Energy)
	}
}

func TestPlayer_RegenEnergy_RegensHPOnlyWhenOldEnergyWasFull(t *testing.T) {
	pFull := &Player{
		Energy:       10,
		MaxEnergy:    10,
		HP:           5,
		MaxHP:        10,
		LastEnergyAt: time.Now().UTC().Add(-20 * time.Minute),
	}
	pFull.RegenEnergy(10 * time.Minute)
	if pFull.HP != 7 {
		t.Fatalf("expected hp regen to 7, got %d", pFull.HP)
	}

	pNotFull := &Player{
		Energy:       9,
		MaxEnergy:    10,
		HP:           5,
		MaxHP:        10,
		LastEnergyAt: time.Now().UTC().Add(-20 * time.Minute),
	}
	pNotFull.RegenEnergy(10 * time.Minute)
	if pNotFull.HP != 5 {
		t.Fatalf("expected no hp regen when energy was not full, got %d", pNotFull.HP)
	}
}
