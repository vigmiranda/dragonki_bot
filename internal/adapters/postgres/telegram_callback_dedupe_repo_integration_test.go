package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestTelegramCallbackDedupeRepo_ReserveIsIdempotent(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewTelegramCallbackDedupeRepo(testDB.DB)

	entry := &domain.TelegramCallbackDedupe{
		DedupeKey:   "win:77:123:99:action_fight:100",
		Action:      "action_fight",
		WindowStart: time.Now().UTC().Truncate(3 * time.Second),
	}

	reserved, err := repo.Reserve(context.Background(), entry)
	if err != nil {
		t.Fatalf("reserve first dedupe: %v", err)
	}
	if !reserved {
		t.Fatal("expected first reserve to insert dedupe")
	}

	reserved, err = repo.Reserve(context.Background(), entry)
	if err != nil {
		t.Fatalf("reserve duplicate dedupe: %v", err)
	}
	if reserved {
		t.Fatal("expected duplicate reserve to be ignored")
	}

	var count int64
	if err := testDB.DB.Model(&domain.TelegramCallbackDedupe{}).Count(&count).Error; err != nil {
		t.Fatalf("count dedupes: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 dedupe row, got %d", count)
	}
}

func TestTelegramCallbackDedupeRepo_DeleteOlderThan(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewTelegramCallbackDedupeRepo(testDB.DB)

	oldEntry := &domain.TelegramCallbackDedupe{
		DedupeKey:   "old-key",
		Action:      "action_fight",
		WindowStart: time.Now().UTC().Add(-10 * time.Minute),
	}
	newEntry := &domain.TelegramCallbackDedupe{
		DedupeKey:   "new-key",
		Action:      "big_boss_attack",
		WindowStart: time.Now().UTC(),
	}

	if _, err := repo.Reserve(context.Background(), oldEntry); err != nil {
		t.Fatalf("reserve old dedupe: %v", err)
	}
	if _, err := repo.Reserve(context.Background(), newEntry); err != nil {
		t.Fatalf("reserve new dedupe: %v", err)
	}

	deleted, err := repo.DeleteOlderThan(context.Background(), time.Now().UTC().Add(-5*time.Minute))
	if err != nil {
		t.Fatalf("delete old dedupes: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted dedupe, got %d", deleted)
	}

	var count int64
	if err := testDB.DB.Model(&domain.TelegramCallbackDedupe{}).Count(&count).Error; err != nil {
		t.Fatalf("count dedupes: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 remaining dedupe row, got %d", count)
	}
}
