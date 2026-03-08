package application

import (
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type fakeReferralRepo struct {
	referral            *domain.Referral
	getErr              error
	createCalls         int
	grantCalls          int
	markCalls           int
	lastGrantReferrerID int64
	lastGrantGold       int
	lastGrantKi         int
	referrerUserIDByID  map[int64]int64
}

func (f *fakeReferralRepo) Create(referral *domain.Referral) error {
	f.createCalls++
	copyRef := *referral
	if copyRef.ID == 0 {
		copyRef.ID = 1
	}
	f.referral = &copyRef
	return nil
}

func (f *fakeReferralRepo) GetByReferredPlayerID(referredPlayerID int64) (*domain.Referral, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.referral == nil || f.referral.ReferredPlayerID != referredPlayerID {
		return nil, gorm.ErrRecordNotFound
	}
	copyRef := *f.referral
	return &copyRef, nil
}

func (f *fakeReferralRepo) MarkRewarded(referralID uint) error {
	f.markCalls++
	if f.referral != nil && f.referral.ID == referralID {
		f.referral.Status = domain.ReferralStatusRewarded
	}
	return nil
}

func (f *fakeReferralRepo) GrantReward(referrerPlayerID int64, gold int, kiCrystals int) error {
	f.grantCalls++
	f.lastGrantReferrerID = referrerPlayerID
	f.lastGrantGold = gold
	f.lastGrantKi = kiCrystals
	return nil
}

func (f *fakeReferralRepo) CountByReferrer(referrerPlayerID int64) (int64, int64, error) {
	return 0, 0, nil
}

func (f *fakeReferralRepo) GetPlayerUserIDByPlayerID(playerID int64) (int64, error) {
	if f.referrerUserIDByID == nil {
		return 0, nil
	}
	return f.referrerUserIDByID[playerID], nil
}

func TestReferralService_RegisterByStartParam_SelfReferralIgnored(t *testing.T) {
	repo := &fakeReferralRepo{}
	svc := NewReferralService(repo, &fakePlayerRepo{}, "dragonki_bot")

	result, err := svc.RegisterByStartParam("/start ref_10", 10, "ONBOARDING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for self referral")
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected no create call, got %d", repo.createCalls)
	}
}

func TestReferralService_RegisterByStartParam_DuplicateIgnored(t *testing.T) {
	repo := &fakeReferralRepo{
		referral: &domain.Referral{
			ID:               10,
			ReferrerPlayerID: 1,
			ReferredPlayerID: 2,
			Status:           domain.ReferralStatusPending,
			RequiredLevel:    7,
		},
	}
	svc := NewReferralService(repo, &fakePlayerRepo{}, "dragonki_bot")

	result, err := svc.RegisterByStartParam("/start ref_1", 2, "ONBOARDING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result for duplicate referral")
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected no create call, got %d", repo.createCalls)
	}
}

func TestReferralService_ProcessLevelUp_RewardOnlyOnce(t *testing.T) {
	repo := &fakeReferralRepo{
		referral: &domain.Referral{
			ID:               20,
			ReferrerPlayerID: 99,
			ReferredPlayerID: 33,
			Status:           domain.ReferralStatusPending,
			RequiredLevel:    7,
		},
		referrerUserIDByID: map[int64]int64{
			99: 12345,
		},
	}
	svc := NewReferralService(repo, &fakePlayerRepo{}, "dragonki_bot")

	first, err := svc.ProcessLevelUp(33, 7)
	if err != nil {
		t.Fatalf("unexpected error on first reward: %v", err)
	}
	if first == nil || !first.Granted {
		t.Fatalf("expected granted reward on first level up")
	}
	if repo.grantCalls != 1 || repo.markCalls != 1 {
		t.Fatalf("expected one reward+mark call, got grant=%d mark=%d", repo.grantCalls, repo.markCalls)
	}

	second, err := svc.ProcessLevelUp(33, 8)
	if err != nil {
		t.Fatalf("unexpected error on second reward: %v", err)
	}
	if second != nil {
		t.Fatalf("expected nil result on second level up (already rewarded)")
	}
	if repo.grantCalls != 1 || repo.markCalls != 1 {
		t.Fatalf("expected reward to run once only, got grant=%d mark=%d", repo.grantCalls, repo.markCalls)
	}
}
