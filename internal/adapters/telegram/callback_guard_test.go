package telegram

import (
	"context"
	"sync"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakeTelegramDelivery struct {
	sendCalls    int
	requestCalls int
	lastSend     tgbotapi.Chattable
	lastRequest  tgbotapi.Chattable
}

func (f *fakeTelegramDelivery) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.sendCalls++
	f.lastSend = c
	return tgbotapi.Message{}, nil
}

func (f *fakeTelegramDelivery) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	f.requestCalls++
	f.lastRequest = c
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (f *fakeTelegramDelivery) Close() {}

type fakeTelegramCallbackReceiptStore struct {
	mu       sync.Mutex
	reserved map[string]struct{}
}

func (f *fakeTelegramCallbackReceiptStore) Reserve(_ context.Context, receipt *domain.TelegramCallbackReceipt) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.reserved == nil {
		f.reserved = make(map[string]struct{})
	}
	if _, exists := f.reserved[receipt.CallbackID]; exists {
		return false, nil
	}
	f.reserved[receipt.CallbackID] = struct{}{}
	return true, nil
}

type fakeTelegramCallbackDedupeStore struct {
	mu       sync.Mutex
	reserved map[string]struct{}
}

func (f *fakeTelegramCallbackDedupeStore) Reserve(_ context.Context, dedupe *domain.TelegramCallbackDedupe) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.reserved == nil {
		f.reserved = make(map[string]struct{})
	}
	if _, exists := f.reserved[dedupe.DedupeKey]; exists {
		return false, nil
	}
	f.reserved[dedupe.DedupeKey] = struct{}{}
	return true, nil
}

type fakeGuardBattleRegionRepo struct {
	region *domain.Region
}

func (f *fakeGuardBattleRegionRepo) GetByID(regionID int64) (*domain.Region, error) {
	if f.region == nil {
		return &domain.Region{RegionID: regionID, Nerf: 10}, nil
	}
	return f.region, nil
}

func (f *fakeGuardBattleRegionRepo) ListByLevel(level int64) ([]domain.Region, error) {
	return nil, nil
}

type fakeGuardBattleRepo struct{}

func (f *fakeGuardBattleRepo) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	return nil
}

func (f *fakeGuardBattleRepo) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	return nil
}

func (f *fakeGuardBattleRepo) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	return nil, nil
}

func TestHandleCallback_IgnoresDuplicateFightCallback(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	playerRepo := &fakeBattlePlayerRepo{
		player: &domain.Player{
			ID:           1,
			UserID:       99,
			Name:         "Goku",
			Gender:       "man",
			ClassID:      2,
			Class:        domain.Class{Name: "Ki Vanguard", Skin: "sayajin", AtkBonus: 0.2, DefBonus: 0.2, HpBonus: 0.1, SpeedBonus: 0.1},
			Level:        10,
			XP:           15,
			HP:           100,
			MaxHP:        100,
			ATK:          30,
			DEF:          12,
			Speed:        8,
			Energy:       3,
			MaxEnergy:    5,
			Gold:         50,
			RegionID:     1,
			Form:         "Base",
			LastEnergyAt: time.Now().UTC(),
		},
	}
	monsterRepo := &fakeBattleMonsterRepo{
		monster: &domain.Monster{
			MonsterID: 9,
			Name:      "Saibaman",
			Tag:       "saibaman",
			Level:     5,
			HP:        20,
			ATK:       5,
			DEF:       1,
			RegionID:  1,
			MinXP:     10,
			MaxXP:     20,
			MinMoney:  5,
			MaxMoney:  10,
		},
	}
	inventoryRepo := &fakeBattleInventoryRepo{}
	delivery := &fakeTelegramDelivery{}

	h := &Handler{
		battleService: application.NewBattleService(
			playerRepo,
			monsterRepo,
			&fakeGuardBattleRegionRepo{region: &domain.Region{RegionID: 1, Nerf: 10, Tag: "plain_lv1"}},
			&fakeGuardBattleRepo{},
			inventoryRepo,
		),
		presenter:     newTelegramPresenter(delivery),
		callbackGuard: newTelegramCallbackGuard(3 * time.Second),
	}

	update := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb-1",
			From: &tgbotapi.User{ID: 99},
			Data: "action_fight?monsterID=9",
			Message: &tgbotapi.Message{
				MessageID: 123,
				Chat:      &tgbotapi.Chat{ID: 77},
			},
		},
	}

	h.handleCallback(update)
	if delivery.sendCalls != 1 {
		t.Fatalf("expected first callback to edit panel once, got %d sends", delivery.sendCalls)
	}
	if delivery.requestCalls != 1 {
		t.Fatalf("expected first callback to ack once, got %d requests", delivery.requestCalls)
	}

	update.CallbackQuery.ID = "cb-2"
	h.handleCallback(update)

	if delivery.sendCalls != 1 {
		t.Fatalf("expected duplicate callback to skip second edit, got %d sends", delivery.sendCalls)
	}
	if delivery.requestCalls != 2 {
		t.Fatalf("expected duplicate callback to still ack, got %d requests", delivery.requestCalls)
	}
	if got := bootstrap.SnapshotAppMetrics()["telegram.callback.duplicate_local_total"]; got != 1 {
		t.Fatalf("expected local duplicate counter 1, got %d", got)
	}
}

func TestTelegramCallbackGuard_ReservesNewMutatingActions(t *testing.T) {
	guard := newTelegramCallbackGuard(3 * time.Second)

	testCases := []struct {
		name   string
		action string
		params map[string]string
	}{
		{name: "big boss attack", action: "big_boss_attack"},
		{name: "start chamber run", action: "start_ancestral_ki_run", params: map[string]string{"level": "3"}},
		{name: "abandon chamber run", action: "abandon_ancestral_ki_run"},
		{name: "continue chamber run", action: "ancestral_ki_chamber_continue"},
		{name: "use chamber item", action: "chamber_use_item", params: map[string]string{"inventoryItemID": "7"}},
		{name: "transform chamber", action: "chamber_transform"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !guard.TryReserve(77, 123, 99, tc.action, tc.params) {
				t.Fatalf("expected first reservation for %s to succeed", tc.action)
			}
			if guard.TryReserve(77, 123, 99, tc.action, tc.params) {
				t.Fatalf("expected duplicate reservation for %s to be blocked", tc.action)
			}
		})
	}
}

func TestHandleCallback_IgnoresDuplicateFightCallbackAcrossHandlersByReceipt(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	playerRepo := &fakeBattlePlayerRepo{
		player: &domain.Player{
			ID:           1,
			UserID:       99,
			Name:         "Goku",
			Gender:       "man",
			ClassID:      2,
			Class:        domain.Class{Name: "Ki Vanguard", Skin: "sayajin", AtkBonus: 0.2, DefBonus: 0.2, HpBonus: 0.1, SpeedBonus: 0.1},
			Level:        10,
			XP:           15,
			HP:           100,
			MaxHP:        100,
			ATK:          30,
			DEF:          12,
			Speed:        8,
			Energy:       3,
			MaxEnergy:    5,
			Gold:         50,
			RegionID:     1,
			Form:         "Base",
			LastEnergyAt: time.Now().UTC(),
		},
	}
	monsterRepo := &fakeBattleMonsterRepo{
		monster: &domain.Monster{
			MonsterID: 9,
			Name:      "Saibaman",
			Tag:       "saibaman",
			Level:     5,
			HP:        20,
			ATK:       5,
			DEF:       1,
			RegionID:  1,
			MinXP:     10,
			MaxXP:     20,
			MinMoney:  5,
			MaxMoney:  10,
		},
	}
	inventoryRepo := &fakeBattleInventoryRepo{}
	receiptStore := &fakeTelegramCallbackReceiptStore{}
	delivery1 := &fakeTelegramDelivery{}
	delivery2 := &fakeTelegramDelivery{}

	buildHandler := func(delivery *fakeTelegramDelivery) *Handler {
		return &Handler{
			battleService: application.NewBattleService(
				playerRepo,
				monsterRepo,
				&fakeGuardBattleRegionRepo{region: &domain.Region{RegionID: 1, Nerf: 10, Tag: "plain_lv1"}},
				&fakeGuardBattleRepo{},
				inventoryRepo,
			),
			presenter:        newTelegramPresenter(delivery),
			callbackReceipts: receiptStore,
		}
	}

	update := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb-shared",
			From: &tgbotapi.User{ID: 99},
			Data: "action_fight?monsterID=9",
			Message: &tgbotapi.Message{
				MessageID: 123,
				Chat:      &tgbotapi.Chat{ID: 77},
			},
		},
	}

	buildHandler(delivery1).handleCallback(update)
	buildHandler(delivery2).handleCallback(update)

	if delivery1.sendCalls != 1 {
		t.Fatalf("expected first handler to edit panel once, got %d sends", delivery1.sendCalls)
	}
	if delivery1.requestCalls != 1 {
		t.Fatalf("expected first handler to ack once, got %d requests", delivery1.requestCalls)
	}
	if delivery2.sendCalls != 0 {
		t.Fatalf("expected second handler to skip edit on duplicate receipt, got %d sends", delivery2.sendCalls)
	}
	if delivery2.requestCalls != 1 {
		t.Fatalf("expected second handler to ack duplicate receipt once, got %d requests", delivery2.requestCalls)
	}
	if got := bootstrap.SnapshotAppMetrics()["telegram.callback.duplicate_receipt_total"]; got != 1 {
		t.Fatalf("expected receipt duplicate counter 1, got %d", got)
	}
}

func TestHandleCallback_IgnoresDuplicateFightCallbackAcrossHandlersBySemanticWindow(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	playerRepo := &fakeBattlePlayerRepo{
		player: &domain.Player{
			ID:           1,
			UserID:       99,
			Name:         "Goku",
			Gender:       "man",
			ClassID:      2,
			Class:        domain.Class{Name: "Ki Vanguard", Skin: "sayajin", AtkBonus: 0.2, DefBonus: 0.2, HpBonus: 0.1, SpeedBonus: 0.1},
			Level:        10,
			XP:           15,
			HP:           100,
			MaxHP:        100,
			ATK:          30,
			DEF:          12,
			Speed:        8,
			Energy:       3,
			MaxEnergy:    5,
			Gold:         50,
			RegionID:     1,
			Form:         "Base",
			LastEnergyAt: time.Now().UTC(),
		},
	}
	monsterRepo := &fakeBattleMonsterRepo{
		monster: &domain.Monster{
			MonsterID: 9,
			Name:      "Saibaman",
			Tag:       "saibaman",
			Level:     5,
			HP:        20,
			ATK:       5,
			DEF:       1,
			RegionID:  1,
			MinXP:     10,
			MaxXP:     20,
			MinMoney:  5,
			MaxMoney:  10,
		},
	}
	inventoryRepo := &fakeBattleInventoryRepo{}
	dedupeStore := &fakeTelegramCallbackDedupeStore{}
	delivery1 := &fakeTelegramDelivery{}
	delivery2 := &fakeTelegramDelivery{}

	buildHandler := func(delivery *fakeTelegramDelivery) *Handler {
		return &Handler{
			battleService: application.NewBattleService(
				playerRepo,
				monsterRepo,
				&fakeGuardBattleRegionRepo{region: &domain.Region{RegionID: 1, Nerf: 10, Tag: "plain_lv1"}},
				&fakeGuardBattleRepo{},
				inventoryRepo,
			),
			presenter:       newTelegramPresenter(delivery),
			callbackDedupes: dedupeStore,
		}
	}

	update1 := tgbotapi.Update{
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "cb-a",
			From: &tgbotapi.User{ID: 99},
			Data: "action_fight?monsterID=9",
			Message: &tgbotapi.Message{
				MessageID: 123,
				Chat:      &tgbotapi.Chat{ID: 77},
			},
		},
	}
	update2 := update1
	update2.CallbackQuery = &tgbotapi.CallbackQuery{
		ID:   "cb-b",
		From: &tgbotapi.User{ID: 99},
		Data: "action_fight?monsterID=9",
		Message: &tgbotapi.Message{
			MessageID: 123,
			Chat:      &tgbotapi.Chat{ID: 77},
		},
	}

	buildHandler(delivery1).handleCallback(update1)
	buildHandler(delivery2).handleCallback(update2)

	if delivery1.sendCalls != 1 {
		t.Fatalf("expected first handler to edit panel once, got %d sends", delivery1.sendCalls)
	}
	if delivery1.requestCalls != 1 {
		t.Fatalf("expected first handler to ack once, got %d requests", delivery1.requestCalls)
	}
	if delivery2.sendCalls != 0 {
		t.Fatalf("expected second handler to skip edit on semantic duplicate, got %d sends", delivery2.sendCalls)
	}
	if delivery2.requestCalls != 1 {
		t.Fatalf("expected second handler to ack semantic duplicate once, got %d requests", delivery2.requestCalls)
	}
	if got := bootstrap.SnapshotAppMetrics()["telegram.callback.duplicate_semantic_total"]; got != 1 {
		t.Fatalf("expected semantic duplicate counter 1, got %d", got)
	}
}
