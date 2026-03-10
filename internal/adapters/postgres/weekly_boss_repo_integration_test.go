package postgres

import (
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestWeeklyBossRepoIntegration_BossLifecycleContributionsAndPayouts(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	repo := NewWeeklyBossPostgresRepo(testDB.DB)

	firstPlayer, err := playerRepo.GetOrCreate(730001)
	if err != nil {
		t.Fatalf("create first player: %v", err)
	}
	secondPlayer, err := playerRepo.GetOrCreate(730002)
	if err != nil {
		t.Fatalf("create second player: %v", err)
	}

	firstPlayer.Name = "Alpha"
	if err := playerRepo.Update(firstPlayer); err != nil {
		t.Fatalf("update first player name: %v", err)
	}
	secondPlayer.Name = "Beta"
	if err := playerRepo.Update(secondPlayer); err != nil {
		t.Fatalf("update second player name: %v", err)
	}

	now := time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC)
	boss := &domain.WeeklyBoss{
		Name:      "Integration Boss",
		Tag:       "IT_BOSS",
		WeekKey:   "2026-W11",
		MaxHP:     1000,
		CurrentHP: 1000,
		ATK:       120,
		Status:    domain.WeeklyBossStatusActive,
		StartsAt:  now.Add(-1 * time.Hour),
		EndsAt:    now.Add(24 * time.Hour),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := repo.Create(boss); err != nil {
		t.Fatalf("create weekly boss: %v", err)
	}

	activeBoss, err := repo.GetActive(now)
	if err != nil {
		t.Fatalf("get active weekly boss: %v", err)
	}
	if activeBoss.ID != boss.ID {
		t.Fatalf("unexpected active boss id: got %d want %d", activeBoss.ID, boss.ID)
	}

	byWeekKey, err := repo.GetByWeekKey("2026-W11")
	if err != nil {
		t.Fatalf("get weekly boss by week key: %v", err)
	}
	if byWeekKey.ID != boss.ID {
		t.Fatalf("unexpected boss by week key id: got %d want %d", byWeekKey.ID, boss.ID)
	}

	hasActive, err := repo.HasActive()
	if err != nil {
		t.Fatalf("has active boss: %v", err)
	}
	if !hasActive {
		t.Fatal("expected active boss to exist")
	}

	attackDate := now
	firstAttack := &domain.WeeklyBossAttack{
		BossID:    boss.ID,
		PlayerID:  int64(firstPlayer.ID),
		AttackAt:  attackDate,
		Damage:    250,
		CreatedAt: now,
	}
	if err := repo.CreateAttack(firstAttack); err != nil {
		t.Fatalf("create first attack: %v", err)
	}

	count, err := repo.CountPlayerAttacksOnDate(boss.ID, int64(firstPlayer.ID), attackDate)
	if err != nil {
		t.Fatalf("count player attacks on date: %v", err)
	}
	if count != 1 {
		t.Fatalf("unexpected attack count: got %d want 1", count)
	}

	if err := repo.AddContributionDamage(boss.ID, int64(firstPlayer.ID), 250, now); err != nil {
		t.Fatalf("add first player contribution: %v", err)
	}
	if err := repo.AddContributionDamage(boss.ID, int64(firstPlayer.ID), 125, now.Add(1*time.Minute)); err != nil {
		t.Fatalf("add first player second contribution: %v", err)
	}
	if err := repo.AddContributionDamage(boss.ID, int64(secondPlayer.ID), 300, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("add second player contribution: %v", err)
	}

	firstContribution, err := repo.GetContribution(boss.ID, int64(firstPlayer.ID))
	if err != nil {
		t.Fatalf("get first player contribution: %v", err)
	}
	if firstContribution.TotalDamage != 375 {
		t.Fatalf("unexpected first player total damage: got %d want 375", firstContribution.TotalDamage)
	}

	ranking, err := repo.ListTopContributions(boss.ID, 2)
	if err != nil {
		t.Fatalf("list top contributions: %v", err)
	}
	if len(ranking) != 2 {
		t.Fatalf("unexpected ranking size: got %d want 2", len(ranking))
	}
	if ranking[0].PlayerID != int64(firstPlayer.ID) || ranking[0].TotalDamage != 375 {
		t.Fatalf("unexpected first ranking row: %+v", ranking[0])
	}
	if ranking[0].PlayerName != "Alpha" {
		t.Fatalf("unexpected first ranking player name: %q", ranking[0].PlayerName)
	}
	if ranking[1].PlayerID != int64(secondPlayer.ID) || ranking[1].TotalDamage != 300 {
		t.Fatalf("unexpected second ranking row: %+v", ranking[1])
	}

	damagedBoss, err := repo.ApplyDamage(boss.ID, 999, now.Add(3*time.Minute))
	if err != nil {
		t.Fatalf("apply damage without defeat: %v", err)
	}
	if damagedBoss.CurrentHP != 1 {
		t.Fatalf("unexpected boss hp after non-lethal hit: got %d want 1", damagedBoss.CurrentHP)
	}
	if damagedBoss.Status != domain.WeeklyBossStatusActive {
		t.Fatalf("unexpected boss status after non-lethal hit: %s", damagedBoss.Status)
	}

	defeatedBoss, err := repo.ApplyDamage(boss.ID, 10, now.Add(4*time.Minute))
	if err != nil {
		t.Fatalf("apply lethal damage: %v", err)
	}
	if defeatedBoss.CurrentHP != 0 {
		t.Fatalf("unexpected boss hp after lethal hit: got %d want 0", defeatedBoss.CurrentHP)
	}
	if defeatedBoss.Status != domain.WeeklyBossStatusDefeated {
		t.Fatalf("unexpected boss status after lethal hit: %s", defeatedBoss.Status)
	}
	if defeatedBoss.DefeatedAt == nil {
		t.Fatal("expected defeated_at to be set after lethal hit")
	}

	latestDefeated, err := repo.FindLatestByStatus(domain.WeeklyBossStatusDefeated)
	if err != nil {
		t.Fatalf("find latest defeated boss: %v", err)
	}
	if latestDefeated.ID != boss.ID {
		t.Fatalf("unexpected latest defeated boss id: got %d want %d", latestDefeated.ID, boss.ID)
	}

	reward, err := repo.GetRewardByPosition(1)
	if err != nil {
		t.Fatalf("get reward by position: %v", err)
	}
	if reward.Position != 1 {
		t.Fatalf("unexpected reward position: got %d want 1", reward.Position)
	}

	payout := &domain.WeeklyBossPayout{
		BossID:    boss.ID,
		PlayerID:  int64(firstPlayer.ID),
		Position:  1,
		GoldPaid:  reward.Gold,
		Status:    domain.WeeklyBossPayoutPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	created, err := repo.CreatePayoutIfNotExists(payout)
	if err != nil {
		t.Fatalf("create payout if not exists: %v", err)
	}
	if !created {
		t.Fatal("expected first payout creation to insert row")
	}

	created, err = repo.CreatePayoutIfNotExists(&domain.WeeklyBossPayout{
		BossID:    boss.ID,
		PlayerID:  int64(firstPlayer.ID),
		Position:  1,
		GoldPaid:  reward.Gold,
		Status:    domain.WeeklyBossPayoutPending,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("create duplicate payout if not exists: %v", err)
	}
	if created {
		t.Fatal("expected duplicate payout creation to be ignored")
	}

	if err := repo.MarkPayoutFailed(payout.ID, "queue timeout"); err != nil {
		t.Fatalf("mark payout failed: %v", err)
	}

	var failedPayout domain.WeeklyBossPayout
	if err := testDB.DB.Where("id = ?", payout.ID).First(&failedPayout).Error; err != nil {
		t.Fatalf("load failed payout: %v", err)
	}
	if failedPayout.Status != domain.WeeklyBossPayoutFailed || failedPayout.ErrorMessage != "queue timeout" {
		t.Fatalf("unexpected failed payout state: %+v", failedPayout)
	}

	paidAt := now.Add(5 * time.Minute)
	if err := repo.MarkPayoutPaid(payout.ID, paidAt); err != nil {
		t.Fatalf("mark payout paid: %v", err)
	}

	var paidPayout domain.WeeklyBossPayout
	if err := testDB.DB.Where("id = ?", payout.ID).First(&paidPayout).Error; err != nil {
		t.Fatalf("load paid payout: %v", err)
	}
	if paidPayout.Status != domain.WeeklyBossPayoutPaid {
		t.Fatalf("unexpected paid payout status: %s", paidPayout.Status)
	}
	if paidPayout.PaidAt == nil {
		t.Fatal("expected paid_at to be set")
	}

	initialGold := firstPlayer.Gold
	if err := repo.AddGoldToPlayer(int64(firstPlayer.ID), reward.Gold); err != nil {
		t.Fatalf("add gold to player: %v", err)
	}

	reloadedFirstPlayer, err := playerRepo.GetOrCreate(firstPlayer.UserID)
	if err != nil {
		t.Fatalf("reload first player after gold payout: %v", err)
	}
	if reloadedFirstPlayer.Gold != initialGold+int(reward.Gold) {
		t.Fatalf("unexpected player gold after payout: got %d want %d", reloadedFirstPlayer.Gold, initialGold+int(reward.Gold))
	}
}
