package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type ReferralRepository interface {
	Create(referral *domain.Referral) error
	GetByReferredPlayerID(referredPlayerID int64) (*domain.Referral, error)
	MarkRewarded(referralID uint) error
	GrantReward(referrerPlayerID int64, gold int, kiCrystals int) error
	CountByReferrer(referrerPlayerID int64) (pending int64, rewarded int64, err error)
	GetPlayerUserIDByPlayerID(playerID int64) (int64, error)
}
