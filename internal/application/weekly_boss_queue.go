package application

import "time"

const BigBossRewardsTaskType = "big_boss_rewards.distribute"

type BigBossRewardsTaskPayload struct {
	BossID      uint      `json:"boss_id"`
	WeekKey     string    `json:"week_key,omitempty"`
	TriggeredAt time.Time `json:"triggered_at,omitempty"`
}

type WeeklyBossRewardQueue interface {
	EnqueueBigBossRewardDistribution(payload BigBossRewardsTaskPayload) error
}
