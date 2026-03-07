package bootstrap

import "testing"

func TestIsFeatureEnabled_UsesDefaultWhenMissing(t *testing.T) {
	SetFeatureFlagResolver(nil)
	t.Cleanup(func() { SetFeatureFlagResolver(nil) })

	if !IsFeatureEnabled("payments", true) {
		t.Fatalf("expected default true when flag is missing")
	}
	if IsFeatureEnabled("payments", false) {
		t.Fatalf("expected default false when flag is missing")
	}
}

func TestIsFeatureEnabled_ReadsResolverValue(t *testing.T) {
	SetFeatureFlagResolver(func(name string) (bool, bool) {
		switch name {
		case "PAYMENTS", "RANK":
			return true, true
		default:
			return false, false
		}
	})
	t.Cleanup(func() { SetFeatureFlagResolver(nil) })

	if !IsFeatureEnabled("payments", false) {
		t.Fatalf("expected payments enabled by resolver")
	}
	if !IsFeatureEnabled("RANK", false) {
		t.Fatalf("expected rank enabled by resolver")
	}
}

func TestIsFeatureEnabled_NormalizesNames(t *testing.T) {
	SetFeatureFlagResolver(func(name string) (bool, bool) {
		if name == "MY_NEW_FEATURE" {
			return true, true
		}
		return false, false
	})
	t.Cleanup(func() { SetFeatureFlagResolver(nil) })

	if !IsFeatureEnabled("my-new-feature", false) {
		t.Fatalf("expected normalized feature name to match")
	}
}

func TestIsFeatureEnabled_UsesDefaultWhenResolverMisses(t *testing.T) {
	SetFeatureFlagResolver(func(name string) (bool, bool) {
		return false, false
	})
	t.Cleanup(func() { SetFeatureFlagResolver(nil) })

	if !IsFeatureEnabled("not_found", true) {
		t.Fatalf("expected default true when resolver misses")
	}
	if IsFeatureEnabled("not_found", false) {
		t.Fatalf("expected default false when resolver misses")
	}
}
