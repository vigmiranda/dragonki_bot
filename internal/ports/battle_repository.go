package ports

import "context"
import "github.com/vigmiranda/dragonki/internal/domain"

type BattleRepository interface {
	SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error
	SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error
}
