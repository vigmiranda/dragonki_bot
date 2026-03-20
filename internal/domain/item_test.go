package domain

import "testing"

func TestAllowsClassByIDAndName(t *testing.T) {
	item := Item{RequiredLevel: 15, AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard"}

	if !item.AllowsClassByIDAndName(4, "Wind Striker") {
		t.Fatalf("expected class id 4 / Wind Striker to be allowed")
	}
	if !item.AllowsClassByIDAndName(2, "Qualquer Nome") {
		t.Fatalf("expected class id 2 to be allowed even with different name")
	}
	if item.AllowsClassByIDAndName(99, "Unknown") {
		t.Fatalf("expected unknown class to be blocked")
	}
}

func TestAllowsClassByIDAndName_IgnoresRestrictionBelowLevel15(t *testing.T) {
	item := Item{
		RequiredLevel:  14,
		AllowedClasses: "Wind Striker,Ki Vanguard",
	}

	if !item.AllowsClassByIDAndName(3, "Earth Sentinel") {
		t.Fatalf("expected class restriction to be ignored below level 15")
	}
}
