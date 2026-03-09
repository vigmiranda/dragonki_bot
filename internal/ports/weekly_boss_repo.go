package ports

import (
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type WeeklyBossRepository interface {
	GetActive(at time.Time) (*domain.WeeklyBoss, error)
	GetByID(id uint) (*domain.WeeklyBoss, error)
	GetByWeekKey(weekKey string) (*domain.WeeklyBoss, error)
	Create(boss *domain.WeeklyBoss) error
	ExpireActiveBefore(at time.Time) error
	GetRandomTemplate() (*domain.WeeklyBossTemplate, error)
	ApplyDamage(id uint, damage int64, at time.Time) (*domain.WeeklyBoss, error)

	CountPlayerAttacksOnDate(bossID uint, playerID int64, attackDate time.Time) (int64, error)
	CreateAttack(attack *domain.WeeklyBossAttack) error

	GetContribution(bossID uint, playerID int64) (*domain.WeeklyBossContribution, error)
	AddContributionDamage(bossID uint, playerID int64, damage int64, attackedAt time.Time) error
	ListTopContributions(bossID uint, limit int) ([]domain.WeeklyBossRankingEntry, error)

	GetRewardByPosition(position int) (*domain.WeeklyBossReward, error)
	CreatePayoutIfNotExists(payout *domain.WeeklyBossPayout) (bool, error)
	MarkPayoutPaid(payoutID uint, paidAt time.Time) error
	MarkPayoutFailed(payoutID uint, reason string) error
	AddGoldToPlayer(playerID int64, gold int64) error
}
