package application

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
	"gorm.io/gorm"
)

const (
	defaultReferralRequiredLevel = 7
	defaultReferralRewardGold    = 1000
	defaultReferralRewardKi      = 1
)

type ReferralService struct {
	referralRepo    ports.ReferralRepository
	playerRepo      ports.PlayerRepository
	botUsername     string
	requiredLevel   int
	rewardGold      int
	rewardKiCrystal int
}

type ReferralRegisterResult struct {
	Linked           bool
	ReferrerPlayerID int64
	ReferrerUserID   int64
	RequiredLevel    int
	RewardGold       int
	RewardKiCrystal  int
}

type ReferralRewardResult struct {
	Granted          bool
	ReferrerPlayerID int64
	ReferrerUserID   int64
	RequiredLevel    int
	RewardGold       int
	RewardKiCrystal  int
}

func NewReferralService(referralRepo ports.ReferralRepository, playerRepo ports.PlayerRepository, botUsername string) *ReferralService {
	return &ReferralService{
		referralRepo:    referralRepo,
		playerRepo:      playerRepo,
		botUsername:     strings.TrimPrefix(strings.TrimSpace(botUsername), "@"),
		requiredLevel:   envIntOrDefault("REFERRAL_REQUIRED_LEVEL", defaultReferralRequiredLevel),
		rewardGold:      envIntOrDefault("REFERRAL_REWARD_GOLD", defaultReferralRewardGold),
		rewardKiCrystal: envIntOrDefault("REFERRAL_REWARD_KI", defaultReferralRewardKi),
	}
}

func (s *ReferralService) OpenReferral(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	link := s.GenerateReferralLink(int64(player.ID))
	pending, rewarded, err := s.referralRepo.CountByReferrer(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	msg := fmt.Sprintf(
		"🎁 Convide amigos para o DragonKi!\n\n"+
			"Seu link:\n%s\n\n"+
			"Quando o indicado atingir nivel %d, voce recebe:\n"+
			"- %d Gold\n"+
			"- %d Ki Crystal(s)\n\n"+
			"Indicacoes pendentes: %d\n"+
			"Indicacoes recompensadas: %d",
		link,
		s.requiredLevel,
		s.rewardGold,
		s.rewardKiCrystal,
		pending,
		rewarded,
	)

	return view.GameResponse{
		Text:    msg,
		MenuKey: "MAIN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/mestre_kalunga.png",
		},
	}, nil
}

func (s *ReferralService) GenerateReferralLink(referrerPlayerID int64) string {
	if s.botUsername == "" {
		return fmt.Sprintf("Use /start ref_%d para convidar (bot username nao configurado).", referrerPlayerID)
	}
	return fmt.Sprintf("https://t.me/%s?start=ref_%d", s.botUsername, referrerPlayerID)
}

func (s *ReferralService) RegisterByStartParam(startParam string, referredPlayerID int64, referredOnboardingStep string) (*ReferralRegisterResult, error) {
	if s.referralRepo == nil {
		return nil, nil
	}
	if strings.TrimSpace(referredOnboardingStep) != "ONBOARDING" {
		return nil, nil
	}

	referrerPlayerID, ok := parseReferralStartParam(startParam)
	if !ok || referrerPlayerID <= 0 || referrerPlayerID == referredPlayerID {
		return nil, nil
	}

	existing, err := s.referralRepo.GetByReferredPlayerID(referredPlayerID)
	if err == nil && existing != nil {
		return nil, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.referralRepo.Create(&domain.Referral{
		ReferrerPlayerID: referrerPlayerID,
		ReferredPlayerID: referredPlayerID,
		Status:           domain.ReferralStatusPending,
		RequiredLevel:    s.requiredLevel,
	}); err != nil {
		return nil, err
	}

	referrerUserID, err := s.referralRepo.GetPlayerUserIDByPlayerID(referrerPlayerID)
	if err != nil {
		referrerUserID = 0
	}

	return &ReferralRegisterResult{
		Linked:           true,
		ReferrerPlayerID: referrerPlayerID,
		ReferrerUserID:   referrerUserID,
		RequiredLevel:    s.requiredLevel,
		RewardGold:       s.rewardGold,
		RewardKiCrystal:  s.rewardKiCrystal,
	}, nil
}

func (s *ReferralService) ProcessLevelUp(referredPlayerID int64, newLevel int) (*ReferralRewardResult, error) {
	if s.referralRepo == nil {
		return nil, nil
	}
	referral, err := s.referralRepo.GetByReferredPlayerID(referredPlayerID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if referral == nil || referral.Status == domain.ReferralStatusRewarded {
		return nil, nil
	}
	required := referral.RequiredLevel
	if required <= 0 {
		required = s.requiredLevel
	}
	if newLevel < required {
		return nil, nil
	}

	if err := s.referralRepo.GrantReward(referral.ReferrerPlayerID, s.rewardGold, s.rewardKiCrystal); err != nil {
		return nil, err
	}
	if err := s.referralRepo.MarkRewarded(referral.ID); err != nil {
		return nil, err
	}

	referrerUserID, err := s.referralRepo.GetPlayerUserIDByPlayerID(referral.ReferrerPlayerID)
	if err != nil {
		referrerUserID = 0
	}

	return &ReferralRewardResult{
		Granted:          true,
		ReferrerPlayerID: referral.ReferrerPlayerID,
		ReferrerUserID:   referrerUserID,
		RequiredLevel:    required,
		RewardGold:       s.rewardGold,
		RewardKiCrystal:  s.rewardKiCrystal,
	}, nil
}

func parseReferralStartParam(raw string) (int64, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	if strings.HasPrefix(trimmed, "/start") {
		parts := strings.Fields(trimmed)
		if len(parts) < 2 {
			return 0, false
		}
		trimmed = parts[1]
	}
	if !strings.HasPrefix(strings.ToLower(trimmed), "ref_") {
		return 0, false
	}
	idRaw := strings.TrimSpace(trimmed[4:])
	id, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func envIntOrDefault(key string, def int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return def
	}
	return v
}
