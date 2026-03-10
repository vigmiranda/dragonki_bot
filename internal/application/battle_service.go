package application

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

const (
	menuMain      = "MAIN_MENU"
	menuBattle    = "BATTLE_MENU"
	menuEnergy    = "ENERGY_MENU"
	menuChoiceCls = "CHOICE_CLASS"

	formKiAwakened = "Ki Awakened"
	formBase       = "Base"

	monsterAurumKarunTag = "AURUM_KARUN"
)

type battleConfig struct {
	LevelMult       int
	EnergyRegenTime string
	MaxPlayerLevel  int
}

type BattleService struct {
	playerRepo          ports.PlayerRepository
	monsterRepo         ports.MonsterRepository
	regionRepo          ports.RegionRepository
	battleRepo          ports.BattleRepository
	inventoryRepository ports.InventoryRepository
	referralService     *ReferralService
	vipService          *VIPService

	cfg battleConfig
	rng *rand.Rand
}

func NewBattleService(
	playerRepo ports.PlayerRepository,
	monsterRepo ports.MonsterRepository,
	regionRepo ports.RegionRepository,
	battleRepo ports.BattleRepository,
	inventoryRepo ports.InventoryRepository,
	referralService ...*ReferralService,
) *BattleService {
	var refService *ReferralService
	if len(referralService) > 0 {
		refService = referralService[0]
	}
	return &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          battleRepo,
		inventoryRepository: inventoryRepo,
		referralService:     refService,
		cfg:                 loadBattleConfig(resolveAppConfig(getAppConfig())),
		rng:                 rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func loadBattleConfig(appCfg appConfig) battleConfig {
	cfg := battleConfig{
		LevelMult:       appCfg.LevelMult,
		EnergyRegenTime: appCfg.EnergyRegenTime,
		MaxPlayerLevel:  appCfg.MaxPlayerLevel,
	}
	return cfg
}

func (s *BattleService) maxPlayerLevel() int {
	if s.cfg.MaxPlayerLevel > 0 {
		return s.cfg.MaxPlayerLevel
	}
	return 50
}

func (s *BattleService) SetVIPService(vipService *VIPService) {
	s.vipService = vipService
}

func (s *BattleService) GetBattle(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	// todo: tirar daqui (mantive comportamento, mas sem derrubar o processo)
	player.SetTransformBase()
	if err := s.playerRepo.Update(player); err != nil {
		return view.GameResponse{}, err
	}
	// ------------

	if player.Form == formKiAwakened {
		player = domain.KiAwakened(player)
	}

	stats, statsErr := calculateEffectivePlayerStats(player, s.inventoryRepository)
	if statsErr != nil {
		log.Printf("battle preview equipment bonus unavailable | playerID=%d | err=%v", player.ID, statsErr)
		stats = effectivePlayerStats{
			ATK:       player.ATK,
			DEF:       player.DEF,
			Speed:     player.Speed,
			MaxHP:     player.MaxHP,
			MaxEnergy: player.MaxEnergy,
		}
	}

	msg := fmt.Sprintf("Energia esgotada.\nAguarde %s para recarregar.", s.cfg.EnergyRegenTime)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%s/tired.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
	}
	menu := menuMain
	monsterID := 0

	if player.Energy > 0 {
		monsters, err := s.monsterRepo.ListByRegion(player.RegionID)
		if err != nil || len(monsters) == 0 {
			msg = "Error ao encontra um monstro"
			img.Path = "assets/error.png"
		} else {
			monster := pickMonsterByWeight(s.rng, monsters)
			if monster == nil {
				msg = "Error ao encontra um monstro"
				img.Path = "assets/error.png"
			} else {
				msg = fmt.Sprintf(
					"🐉<b>%s apareceu!</b>\n\n"+
						"Voce\n"+
						"💪 <b>Lv:</b> %d\n"+
						"❤️ <b>HP:</b> %d/%d\n"+
						"⚔️ <b>ATK:</b> %d  🛡️<b>DEF:</b> %d  🏃🏻<b>VEL:</b> %d\n"+
						"Bonus Equip: ATK %+d | DEF %+d | VEL %+d\n\n"+
						"%s \n"+
						"💪 <b>Lv:</b> %d\n"+
						"❤️ <b>HP:</b> %d\n"+
						"⚔️ <b>ATK:</b> %d  🛡️<b>DEF:</> %d\n",
					monster.Name,
					player.Level,
					player.HP,
					player.MaxHP,
					stats.ATK,
					stats.DEF,
					stats.Speed,
					stats.Bonus.ATK,
					stats.Bonus.DEF,
					stats.Bonus.Speed,
					monster.Name,
					monster.Level,
					monster.HP,
					monster.ATK,
					monster.DEF,
				)
				img.Path = fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.RegionID, monster.Tag)
				monsterID = monster.MonsterID
				menu = menuBattle
			}
		}
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
		MenuParams: map[string]any{
			"monsterID": monsterID,
		},
	}

	return viewResponse, nil
}

func (s *BattleService) Battle(userID int64, params map[string]string) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	monsterID, err := parseMonsterIDParam(params)
	if err != nil {
		return view.GameResponse{}, err
	}

	monster, err := s.monsterRepo.GetByMonsterID(monsterID)
	if err != nil {
		return view.GameResponse{}, err
	}

	region, err := s.regionRepo.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if monster.Tag == monsterAurumKarunTag {
		return s.AurumKarunLuck(monster, region, player)
	}

	stats, statsErr := calculateEffectivePlayerStats(player, s.inventoryRepository)
	if statsErr != nil {
		log.Printf("battle equipment bonus unavailable | playerID=%d | err=%v", player.ID, statsErr)
		stats = effectivePlayerStats{
			ATK:       player.ATK,
			DEF:       player.DEF,
			Speed:     player.Speed,
			MaxHP:     player.MaxHP,
			MaxEnergy: player.MaxEnergy,
		}
	}

	state := initBattleState(player, monster, stats)
	s.executeCombat(player, monster, state)

	xpBase := clamp(CalculateXPReward(player.Level, int(monster.Level)), monster.MinXP, monster.MaxXP)
	goldBase := randBetweenInclusive(s.rng, monster.MinMoney, monster.MaxMoney)
	xpReceived := xpBase
	moneyReceived := goldBase
	vipBonusXP := 0
	vipBonusGold := 0
	if s.vipService != nil {
		xpReceived, moneyReceived = s.vipService.ApplyBattleBonuses(int64(player.ID), xpReceived, moneyReceived)
		vipBonusXP = xpReceived - xpBase
		vipBonusGold = moneyReceived - goldBase
	}

	outcome, err := s.resolveBattleOutcome(player, monster, region, state, xpReceived, moneyReceived, vipBonusXP, vipBonusGold)
	if err != nil {
		return view.GameResponse{}, err
	}
	if s.referralService != nil && player.Level > state.oldLevel {
		rewardResult, referralErr := s.referralService.ProcessLevelUp(int64(player.ID), player.Level)
		if referralErr != nil {
			log.Printf("referral level-up processing failed | playerID=%d | level=%d | err=%v", player.ID, player.Level, referralErr)
		} else if rewardResult != nil && rewardResult.Granted {
			outcome.msg += fmt.Sprintf(
				"\n\n🎁 Bônus de indicação liberado!\nSeu convidador recebeu %d Gold e %d Ki Crystal(s).",
				rewardResult.RewardGold,
				rewardResult.RewardKiCrystal,
			)
		}
	}

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          monster.MonsterID,
		Result:             outcome.resultBattle,
		OldXP:              state.oldXP,
		XPGained:           outcome.xpReceived,
		XPCurrent:          player.XP,
		OldGold:            state.oldGold,
		GoldGained:         outcome.moneyReceived,
		GoldCurrent:        player.Gold,
		PlayerDamageTotal:  state.atkTotalPlayer,
		MonsterDamageTotal: state.atkTotalMonster,
		CreatedAt:          time.Now(),
	}

	if err := s.battleRepo.SaveBattleLogAndPlayer(context.Background(), battleLog, player); err != nil {
		return view.GameResponse{}, err
	}

	msg := outcome.msg
	img := outcome.img
	menu := outcome.menu
	menuParams := make(map[string]any)
	if player.Level > state.oldLevel && player.XP == 0 {
		switch player.Level {
		case 15:
			msg, img, menu, menuParams, err = s.ChangeClass(player)
			if err != nil {
				return view.GameResponse{}, err
			}
		case 35, 50:
			evolveMsg, evolveImg, evolveErr := s.AutoEvolveClass(player)
			if evolveErr != nil {
				return view.GameResponse{}, evolveErr
			}
			if evolveMsg != "" {
				msg = evolveMsg
				if evolveImg != nil {
					img = evolveImg
				}
				menu = menuMain
				menuParams = map[string]any{}
			}
		}
	}

	viewResponse := view.GameResponse{
		Text:       msg,
		Image:      img,
		MenuKey:    menu,
		MenuParams: menuParams,
	}

	return viewResponse, nil
}

func (s *BattleService) AurumKarunLuck(monster *domain.Monster, region *domain.Region, player *domain.Player) (view.GameResponse, error) {
	txt := fmt.Sprintf(
		"\n\n\n\n <b>%s foi mais rapido e escapou!</b>\n\n",
		monster.Name,
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag)),
	}

	skin := ""
	if player.Form == formKiAwakened {
		skin = "ki_awakened_"
	}

	result := "AURUMKARUN/ESCAPE"
	oldGold := player.Gold

	escapePlayer := s.rng.Intn(10)
	escapeAurumKarun := s.rng.Intn(15)

	if escapePlayer >= escapeAurumKarun {
		txt = fmt.Sprintf(
			"Vitoria! %s\n"+
				"+%d Gold\n",
			player.Name,
			monster.MaxMoney,
		)

		img.Path = fmt.Sprintf(
			"assets/skins/%s/%s/%sbattle_win.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
			skin,
		)

		result = "AURUMKARUN/WIN"
		player.Gold += monster.MaxMoney
	}

	viewResponse := view.GameResponse{
		Text:    txt,
		Image:   img,
		MenuKey: menuMain,
	}

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          monster.MonsterID,
		Result:             result,
		OldXP:              player.XP,
		XPGained:           0,
		XPCurrent:          player.XP,
		OldGold:            oldGold,
		GoldGained:         player.Gold - oldGold,
		GoldCurrent:        player.Gold,
		PlayerDamageTotal:  0,
		MonsterDamageTotal: 0,
		CreatedAt:          time.Now(),
	}

	if err := s.battleRepo.SaveBattleLogAndPlayer(context.Background(), battleLog, player); err != nil {
		return view.GameResponse{}, err
	}

	return viewResponse, nil
}

func (s *BattleService) Escape(userID int64, params map[string]string) (view.GameResponse, *domain.Monster, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, nil, err
	}

	monsterIDStr, ok := params["monsterID"]
	if !ok || monsterIDStr == "" {
		return view.GameResponse{}, nil, fmt.Errorf("missing param monsterID")
	}

	monsterID, err := strconv.ParseInt(monsterIDStr, 10, 64)
	if err != nil {
		return view.GameResponse{}, nil, err
	}

	monster, err := s.monsterRepo.GetByMonsterID(monsterID)
	if err != nil {
		return view.GameResponse{}, nil, err
	}

	region, err := s.regionRepo.GetByID(player.RegionID)
	if err != nil {
		return view.GameResponse{}, monster, err
	}

	result := "ESCAPE"
	xpResult := 0
	goldResult := 0
	oldXP := player.XP
	oldGold := player.Gold

	msg := fmt.Sprintf("Voce conseguiu fugir do %s!", monster.Name)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%splayer.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
	}
	menu := menuMain

	escapePlayer := s.rng.Intn(10)
	escapeMonster := s.rng.Intn(10)
	damage := 0

	if escapePlayer <= escapeMonster {
		damage = s.rng.Intn(monster.ATK/2+1) + monster.ATK/2
		player.HP -= damage

		if player.HP <= 0 {
			player.HP = 1

			nerf := float64(region.Nerf) / 100.0
			xpReduction := math.Round(float64(player.XP) * nerf)
			goldReduction := math.Round(float64(player.Gold) * nerf)

			player.XP -= int(xpReduction)
			player.Gold -= int(goldReduction)

			xpResult = -int(xpReduction)
			goldResult = -int(goldReduction)

			menu = menuEnergy
			msg = fmt.Sprintf(
				"Voce zerou seu HP!\n\n"+
					"-%d%% de XP = %d\n"+
					"-%d%% de Gold = %d\n\n"+
					"%s nao e derrota, e aprendizado.\n"+
					"Venha, guerreiro.\nO descanso do Sir Soneca curara suas feridas.\n\n"+
					"O mundo ainda precisa da sua forca.\n",
				region.Nerf,
				int(xpReduction),
				region.Nerf,
				int(goldReduction),
				player.Name,
			)
			img.Path = "assets/npcs/mestre_kalunga.png"

			result = "ESCAPE/LOSS"
		} else {
			menu = menuBattle
			msg = fmt.Sprintf(
				"Voce nao conseguiu fugir!\n"+
					"Seu HP agora e %d/%d\n"+
					"E sua Energia e %d/%d\n",
				player.HP,
				player.MaxHP,
				player.Energy,
				player.MaxEnergy,
			)
			img.Path = fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.RegionID, monster.Tag)
		}
	}

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          monster.MonsterID,
		Result:             result,
		OldXP:              oldXP,
		XPGained:           xpResult,
		XPCurrent:          player.XP,
		OldGold:            oldGold,
		GoldGained:         goldResult,
		GoldCurrent:        player.Gold,
		PlayerDamageTotal:  0,
		MonsterDamageTotal: damage,
		CreatedAt:          time.Now(),
	}

	if err := s.battleRepo.SaveBattleLogAndPlayer(context.Background(), battleLog, player); err != nil {
		return view.GameResponse{}, monster, err
	}

	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
		MenuParams: map[string]any{
			"monsterID": monster.MonsterID,
		},
	}

	return viewResponse, monster, nil
}
