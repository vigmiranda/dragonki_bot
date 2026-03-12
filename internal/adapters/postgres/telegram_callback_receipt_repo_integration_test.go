package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestTelegramCallbackReceiptRepo_ReserveIsIdempotent(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewTelegramCallbackReceiptRepo(testDB.DB)

	receipt := &domain.TelegramCallbackReceipt{
		CallbackID: "cb-123",
		UserID:     99,
		ChatID:     77,
		MessageID:  55,
		Action:     "big_boss_attack",
		RawData:    "big_boss_attack",
		ReceivedAt: time.Now().UTC(),
	}

	reserved, err := repo.Reserve(context.Background(), receipt)
	if err != nil {
		t.Fatalf("reserve first receipt: %v", err)
	}
	if !reserved {
		t.Fatal("expected first reserve to insert receipt")
	}

	reserved, err = repo.Reserve(context.Background(), receipt)
	if err != nil {
		t.Fatalf("reserve duplicate receipt: %v", err)
	}
	if reserved {
		t.Fatal("expected duplicate reserve to be ignored")
	}

	var count int64
	if err := testDB.DB.Model(&domain.TelegramCallbackReceipt{}).Count(&count).Error; err != nil {
		t.Fatalf("count receipts: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 receipt row, got %d", count)
	}
}

func TestTelegramCallbackReceiptRepo_DeleteOlderThan(t *testing.T) {
	testDB := newIntegrationTestDB(t)
	repo := NewTelegramCallbackReceiptRepo(testDB.DB)

	oldReceipt := &domain.TelegramCallbackReceipt{
		CallbackID: "cb-old",
		UserID:     99,
		ChatID:     77,
		MessageID:  55,
		Action:     "action_fight",
		RawData:    "action_fight?monsterID=9",
		ReceivedAt: time.Now().UTC().Add(-48 * time.Hour),
	}
	newReceipt := &domain.TelegramCallbackReceipt{
		CallbackID: "cb-new",
		UserID:     99,
		ChatID:     77,
		MessageID:  56,
		Action:     "big_boss_attack",
		RawData:    "big_boss_attack",
		ReceivedAt: time.Now().UTC(),
	}

	if _, err := repo.Reserve(context.Background(), oldReceipt); err != nil {
		t.Fatalf("reserve old receipt: %v", err)
	}
	if _, err := repo.Reserve(context.Background(), newReceipt); err != nil {
		t.Fatalf("reserve new receipt: %v", err)
	}

	deleted, err := repo.DeleteOlderThan(context.Background(), time.Now().UTC().Add(-24*time.Hour))
	if err != nil {
		t.Fatalf("delete old receipts: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("expected 1 deleted receipt, got %d", deleted)
	}

	var count int64
	if err := testDB.DB.Model(&domain.TelegramCallbackReceipt{}).Count(&count).Error; err != nil {
		t.Fatalf("count receipts: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 remaining receipt row, got %d", count)
	}
}
