package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

func TestBattlePostgresRepoIntegration_SaveBattleLogAndPlayer_PersistsTogether(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	battleRepo := NewBattlePostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	player.HP = 22
	player.Energy = 4
	player.Gold = 180
	player.XP = 55

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          999,
		Result:             "WIN",
		OldXP:              10,
		XPGained:           45,
		XPCurrent:          55,
		OldGold:            100,
		GoldGained:         80,
		GoldCurrent:        180,
		PlayerDamageTotal:  30,
		MonsterDamageTotal: 8,
		CreatedAt:          time.Now().UTC(),
	}

	if err := battleRepo.SaveBattleLogAndPlayer(context.Background(), battleLog, player); err != nil {
		t.Fatalf("save battle log and player atomically: %v", err)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload player: %v", err)
	}
	if reloadedPlayer.HP != 22 || reloadedPlayer.Energy != 4 {
		t.Fatalf("unexpected persisted player vitality: hp=%d energy=%d", reloadedPlayer.HP, reloadedPlayer.Energy)
	}
	if reloadedPlayer.Gold != 180 || reloadedPlayer.XP != 55 {
		t.Fatalf("unexpected persisted player economy/progress: gold=%d xp=%d", reloadedPlayer.Gold, reloadedPlayer.XP)
	}

	var persistedLog domain.BattleLog
	if err := testDB.DB.Where("id = ?", battleLog.ID).Take(&persistedLog).Error; err != nil {
		t.Fatalf("reload persisted battle log: %v", err)
	}
	if persistedLog.Result != "WIN" || persistedLog.GoldGained != 80 || persistedLog.XPGained != 45 {
		t.Fatalf("unexpected persisted battle log: %+v", persistedLog)
	}
}

func TestBattlePostgresRepoIntegration_SaveBattleLogAndPlayer_RollsBackLogWhenPlayerMissing(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	battleRepo := NewBattlePostgresRepo(testDB.DB)
	playerRepo := NewPlayerPostgresRepo(testDB.DB)

	userID := time.Now().UnixNano()
	player, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("create player fixture: %v", err)
	}

	missingPlayer := *player
	missingPlayer.ID = player.ID + 999999
	missingPlayer.HP = 1

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          321,
		Result:             "LOSE",
		OldXP:              20,
		XPGained:           -5,
		XPCurrent:          15,
		OldGold:            80,
		GoldGained:         -10,
		GoldCurrent:        70,
		PlayerDamageTotal:  9,
		MonsterDamageTotal: 25,
		CreatedAt:          time.Now().UTC(),
	}

	err = battleRepo.SaveBattleLogAndPlayer(context.Background(), battleLog, &missingPlayer)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected record not found on missing player update, got %v", err)
	}

	var logCount int64
	if err := testDB.DB.Model(&domain.BattleLog{}).Where("monster_id = ?", 321).Count(&logCount).Error; err != nil {
		t.Fatalf("count rolled back battle logs: %v", err)
	}
	if logCount != 0 {
		t.Fatalf("expected battle log rollback on missing player, got %d rows", logCount)
	}

	reloadedPlayer, err := playerRepo.GetOrCreate(userID)
	if err != nil {
		t.Fatalf("reload original player: %v", err)
	}
	if reloadedPlayer.HP == 1 {
		t.Fatalf("expected original player state to remain unchanged after rollback")
	}
}
