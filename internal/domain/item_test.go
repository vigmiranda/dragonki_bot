package domain

import "testing"

func TestAllowsClassByIDAndName(t *testing.T) {
	item := Item{AllowedClasses: "Wind Striker,Earth Sentinel,Ki Vanguard"}

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
