package ports

import "github.com/vigmiranda/dragonki/internal/domain"

type BattleRepository interface {
	SaveBattleLog(battleLog *domain.BattleLog) error
}
