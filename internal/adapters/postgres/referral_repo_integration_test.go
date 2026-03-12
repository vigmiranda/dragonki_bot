package postgres

import (
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

func TestReferralPostgresRepoIntegration_PersistsCountsAndRewards(t *testing.T) {
	testDB := newIntegrationTestDB(t)

	playerRepo := NewPlayerPostgresRepo(testDB.DB)
	referralRepo := NewReferralPostgresRepo(testDB.DB)

	referrer, err := playerRepo.GetOrCreate(700001)
	if err != nil {
		t.Fatalf("create referrer player: %v", err)
	}

	referred, err := playerRepo.GetOrCreate(700002)
	if err != nil {
		t.Fatalf("create referred player: %v", err)
	}

	initialGold := referrer.Gold
	initialKi := referrer.KiCrystals

	referral := &domain.Referral{
		ReferrerPlayerID: int64(referrer.ID),
		ReferredPlayerID: int64(referred.ID),
		Status:           domain.ReferralStatusPending,
		RequiredLevel:    7,
	}

	if err := referralRepo.Create(referral); err != nil {
		t.Fatalf("create referral: %v", err)
	}

	if referral.ID == 0 {
		t.Fatal("expected referral id to be generated")
	}

	savedReferral, err := referralRepo.GetByReferredPlayerID(int64(referred.ID))
	if err != nil {
		t.Fatalf("get referral by referred player id: %v", err)
	}

	if savedReferral.ReferrerPlayerID != int64(referrer.ID) {
		t.Fatalf("unexpected referrer player id: got %d want %d", savedReferral.ReferrerPlayerID, referrer.ID)
	}
	if savedReferral.Status != domain.ReferralStatusPending {
		t.Fatalf("unexpected initial status: %s", savedReferral.Status)
	}

	pending, rewarded, err := referralRepo.CountByReferrer(int64(referrer.ID))
	if err != nil {
		t.Fatalf("count referrals by referrer: %v", err)
	}

	if pending != 1 || rewarded != 0 {
		t.Fatalf("unexpected counts before reward: pending=%d rewarded=%d", pending, rewarded)
	}

	granted, err := referralRepo.GrantRewardIfPending(referral.ID, int64(referrer.ID), 1000, 2)
	if err != nil {
		t.Fatalf("grant reward if pending: %v", err)
	}
	if !granted {
		t.Fatal("expected first reward grant to succeed")
	}

	granted, err = referralRepo.GrantRewardIfPending(referral.ID, int64(referrer.ID), 1000, 2)
	if err != nil {
		t.Fatalf("grant reward if pending second time: %v", err)
	}
	if granted {
		t.Fatal("expected second reward grant to be ignored")
	}

	userID, err := referralRepo.GetPlayerUserIDByPlayerID(int64(referrer.ID))
	if err != nil {
		t.Fatalf("get player user id by player id: %v", err)
	}

	if userID != referrer.UserID {
		t.Fatalf("unexpected user id for referrer: got %d want %d", userID, referrer.UserID)
	}

	reloadedReferrer, err := playerRepo.GetOrCreate(referrer.UserID)
	if err != nil {
		t.Fatalf("reload referrer player: %v", err)
	}

	if reloadedReferrer.Gold != initialGold+1000 {
		t.Fatalf("unexpected gold after reward: got %d want %d", reloadedReferrer.Gold, initialGold+1000)
	}
	if reloadedReferrer.KiCrystals != initialKi+2 {
		t.Fatalf("unexpected ki after reward: got %d want %d", reloadedReferrer.KiCrystals, initialKi+2)
	}

	rewardedReferral, err := referralRepo.GetByReferredPlayerID(int64(referred.ID))
	if err != nil {
		t.Fatalf("reload rewarded referral: %v", err)
	}

	if rewardedReferral.Status != domain.ReferralStatusRewarded {
		t.Fatalf("unexpected status after reward: %s", rewardedReferral.Status)
	}
	if rewardedReferral.CompletedAt == nil || rewardedReferral.RewardedAt == nil {
		t.Fatalf("expected completed_at and rewarded_at to be set: %+v", rewardedReferral)
	}

	pending, rewarded, err = referralRepo.CountByReferrer(int64(referrer.ID))
	if err != nil {
		t.Fatalf("count referrals by referrer after reward: %v", err)
	}

	if pending != 0 || rewarded != 1 {
		t.Fatalf("unexpected counts after reward: pending=%d rewarded=%d", pending, rewarded)
	}
}
