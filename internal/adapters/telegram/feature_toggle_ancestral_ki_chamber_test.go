package telegram

import (
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func TestHandleAncestralKiChamber_BlockedWhenFeatureDisabled(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "ANCESTRAL_KI_CHAMBER" {
			return false, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	h := &Handler{}
	got := h.handleAncestralKiChamber("to_ancestral_ki_chamber", 1, nil)
	if !strings.Contains(strings.ToLower(got.Message), "indisponivel") {
		t.Fatalf("expected disabled message, got %q", got.Message)
	}
}

func TestHandleTraining_UsesChamberFlowWhenFeatureEnabled(t *testing.T) {
	bootstrap.SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "ANCESTRAL_KI_CHAMBER" {
			return true, true
		}
		return false, false
	})
	t.Cleanup(func() { bootstrap.SetFeatureFlagResolver(nil) })

	h := &Handler{}
	got := h.handleTraining("to_training", 1)
	if !strings.Contains(strings.ToLower(got.Message), "indisponivel") {
		t.Fatalf("expected chamber unavailable message, got %q", got.Message)
	}
}
