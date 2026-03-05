package postgres

import (
	"testing"

	"gorm.io/gorm"
)

func TestPostgresConstructors(t *testing.T) {
	db := &gorm.DB{}

	if got := NewBattlePostgresRepo(db); got == nil {
		t.Fatalf("expected NewBattlePostgresRepo non-nil")
	}
	if got := NewInventoryPostgresRepo(db); got == nil {
		t.Fatalf("expected NewInventoryPostgresRepo non-nil")
	}
	if got := NewItemsPostgresRepo(db); got == nil {
		t.Fatalf("expected NewItemsPostgresRepo non-nil")
	}
	if got := NewMonsterPostgresRepo(db); got == nil {
		t.Fatalf("expected NewMonsterPostgresRepo non-nil")
	}
	if got := NewPaymentRepo(db); got == nil {
		t.Fatalf("expected NewPaymentRepo non-nil")
	}
	if got := NewPlayerPostgresRepo(db); got == nil {
		t.Fatalf("expected NewPlayerPostgresRepo non-nil")
	}
	if got := NewRegionPostgresRepo(db); got == nil {
		t.Fatalf("expected NewRegionPostgresRepo non-nil")
	}
}
