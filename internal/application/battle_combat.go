package application

import (
	"fmt"
	"math"
	"strconv"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type battleState struct {
	oldXP           int
	oldGold         int
	monsterHP       int
	playerATK       int
	playerDEF       int
	playerSpeed     int
	atkTransform    int
	defTransform    int
	skin            string
	atkTotalPlayer  int
	atkTotalMonster int
}

type battleOutcome struct {
	msg           string
	img           *view.Image
	menu          string
	resultBattle  string
	xpReceived    int
	moneyReceived int
}

func parseMonsterIDParam(params map[string]string) (int64, error) {
	monsterIDStr, ok := params["monsterID"]
	if !ok || monsterIDStr == "" {
		return 0, fmt.Errorf("missing param monsterID")
	}
	return strconv.ParseInt(monsterIDStr, 10, 64)
}

func initBattleState(player *domain.Player, monster *domain.Monster, stats effectivePlayerStats) *battleState {
	state := &battleState{
		oldXP:       player.XP,
		oldGold:     player.Gold,
		monsterHP:   monster.HP,
		playerATK:   stats.ATK,
		playerDEF:   stats.DEF,
		playerSpeed: stats.Speed,
	}

	if player.Energy > 0 {
		player.Energy--
	}

	if player.Form == formKiAwakened {
		state.atkTransform = int(float64(state.playerATK) * 0.2)
		state.defTransform = int(float64(state.playerDEF) * 0.2)
		state.skin = "ki_awakened_"
	}

	baseChange := float64(state.playerSpeed) / (float64(state.playerSpeed) + 100) * 35
	if baseChange > 35 {
		baseChange = 35
	}
	player.DodgeAccumulator += baseChange

	return state
}

func (s *BattleService) executeCombat(player *domain.Player, monster *domain.Monster, state *battleState) {
	for state.monsterHP > 0 && player.HP > 0 {
		roll := s.rng.Float64() * 100

		if roll < player.DodgeAccumulator {
			player.DodgeAccumulator = 0
			continue
		}

		damageToMonster := state.playerATK + state.atkTransform
		state.monsterHP -= damageToMonster
		state.atkTotalPlayer += damageToMonster

		atkMonster := s.rng.Intn(monster.ATK+(monster.ATK/2)+1) - (state.playerDEF + state.defTransform)
		if atkMonster < 0 {
			atkMonster = 1
		}

		player.HP -= atkMonster
		state.atkTotalMonster += atkMonster
	}
}

func (s *BattleService) resolveBattleOutcome(
	player *domain.Player,
	monster *domain.Monster,
	region *domain.Region,
	state *battleState,
	xpReceived int,
	moneyReceived int,
) (*battleOutcome, error) {
	if player.HP > 0 {
		return s.resolveWinOutcome(player, monster, state, xpReceived, moneyReceived)
	}
	return s.resolveLoseOutcome(player, region, state), nil
}

func (s *BattleService) resolveWinOutcome(
	player *domain.Player,
	monster *domain.Monster,
	state *battleState,
	xpReceived int,
	moneyReceived int,
) (*battleOutcome, error) {
	goodNews := ""

	itemName, qty, dropErr := s.HandleDrops(monster.MonsterID, int64(player.ID))
	if dropErr != nil {
		return nil, dropErr
	}

	player.XP += xpReceived
	player.Gold += moneyReceived

	if qty > 0 {
		goodNews += fmt.Sprintf("🥳 Dropou %dx %s \n", qty, itemName)
	}

	XPLvlUp := 45 * math.Pow(float64(player.Level), 1.65)
	if player.XP >= int(XPLvlUp) {
		player.Level++
		player.ATK += int(float64(player.ATK) * player.Class.AtkBonus)
		player.DEF += int(float64(player.DEF) * player.Class.DefBonus)
		player.MaxHP += int(float64(player.MaxHP) * player.Class.HpBonus)
		player.Speed += int(float64(player.Speed) * player.Class.SpeedBonus)
		player.Energy = player.MaxEnergy
		player.XP = 0
		player.FormBattlesLeft = 0
		player.Form = formBase
		player.HP = player.MaxHP
		goodNews += "⬆️ <b>Subiu de nível!</b> \n"
		xpReceived = int(XPLvlUp) - state.oldXP
	}

	if player.FormBattlesLeft > 0 {
		player.FormBattlesLeft--
	}

	msg := fmt.Sprintf(
		"🏆 Vitória! %s \n"+
			"%s"+
			"Você recebeu de dano: %d\n"+
			"⚔️ +%d XP\n"+
			"💰 +%d Gold\n"+
			"❤️ HP: %d/%d\n"+
			"⚡️ Energia: %d/%d\n"+
			"💪 Level: %d\n",
		player.Name,
		goodNews,
		state.atkTotalMonster,
		xpReceived,
		moneyReceived,
		player.HP,
		player.MaxHP,
		player.Energy,
		player.MaxEnergy,
		player.Level,
	)

	return &battleOutcome{
		msg: msg,
		img: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf(
				"assets/skins/%s/%s/%sbattle_win.png",
				domain.ConvertLowCase(player.Gender),
				player.Class.Skin,
				state.skin,
			),
		},
		menu:          menuMain,
		resultBattle:  "WIN",
		xpReceived:    xpReceived,
		moneyReceived: moneyReceived,
	}, nil
}

func (s *BattleService) resolveLoseOutcome(player *domain.Player, region *domain.Region, state *battleState) *battleOutcome {
	if player.FormBattlesLeft > 0 {
		player.FormBattlesLeft = 0
		player.Form = formBase
	}

	nerf := float64(region.Nerf) / 100.0
	xpReduction := math.Round(float64(player.XP) * nerf)
	goldReduction := math.Round(float64(player.Gold) * nerf)

	player.XP -= int(xpReduction)
	player.Gold -= int(goldReduction)
	player.HP = 1
	player.Form = formBase

	msg := fmt.Sprintf(
		"💀 Derrota! %s\n"+
			"⚔️ -%d%% de XP = %d\n"+
			"💰 -%d%% de Gold = %d\n"+
			"⚡️ Energia: %d/%d\n"+
			"💪 Level: %d\n",
		player.Name,
		region.Nerf,
		state.oldXP-player.XP,
		region.Nerf,
		state.oldGold-player.Gold,
		player.Energy,
		player.MaxEnergy,
		player.Level,
	)

	return &battleOutcome{
		msg: msg,
		img: &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf(
				"assets/skins/%s/%sbattle_loss.png",
				domain.ConvertLowCase(player.Gender),
				player.Class.Skin,
			),
		},
		menu:          menuMain,
		resultBattle:  "LOSE",
		xpReceived:    player.XP - state.oldXP,
		moneyReceived: player.Gold - state.oldGold,
	}
}

func CalculateXPReward(playerLevel, monsterLevel int) int {
	xpBase := 25 * math.Pow(float64(monsterLevel), 1.12)

	levelDiff := monsterLevel - playerLevel

	multiplier := 1.0

	switch {
	case levelDiff <= -5:
		multiplier = 0.25
	case levelDiff <= -3:
		multiplier = 0.5
	case levelDiff <= -1:
		multiplier = 0.8
	case levelDiff == 0:
		multiplier = 1.0
	case levelDiff <= 2:
		multiplier = 1.15
	case levelDiff <= 4:
		multiplier = 1.3
	default:
		multiplier = 1.5

	}
	xpReward := int(xpBase * multiplier)

	if xpReward < 5 {
		xpReward = 5
	}

	return xpReward
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
