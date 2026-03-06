package domain

type WeeklyBattleRank struct {
	Name         string
	TotalBattles int64
	Wins         int64
	Losses       int64
	Escape       int64
}
