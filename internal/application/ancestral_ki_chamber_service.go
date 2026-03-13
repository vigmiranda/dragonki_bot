package application

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
	"gorm.io/gorm"
)

const (
	ancestralKiChamberEntryEnergyCost = 1
	ancestralKiChamberEntryGoldCost   = 150
	ancestralKiChamberMaxLevel        = 10
	ancestralKiChamberTotalBattles    = 5
	ancestralKiChamberSuccessPct      = 0.40
)

type AncestralKiChamberService struct {
	repo          ports.AncestralKiChamberRepository
	playerRepo    ports.PlayerRepository
	monsterRepo   ports.MonsterRepository
	inventoryRepo ports.InventoryRepository
	rng           *rand.Rand
}

func NewAncestralKiChamberService(
	repo ports.AncestralKiChamberRepository,
	playerRepo ports.PlayerRepository,
	monsterRepo ports.MonsterRepository,
	inventoryRepo ports.InventoryRepository,
) *AncestralKiChamberService {
	return &AncestralKiChamberService{
		repo:          repo,
		playerRepo:    playerRepo,
		monsterRepo:   monsterRepo,
		inventoryRepo: inventoryRepo,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *AncestralKiChamberService) Open(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	progress, err := s.repo.GetOrCreateProgress(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	unlocked := progress.HighestCompletedLevel + 1
	if unlocked < 1 {
		unlocked = 1
	}
	if unlocked > ancestralKiChamberMaxLevel {
		unlocked = ancestralKiChamberMaxLevel
	}

	baseMsg := fmt.Sprintf(
		"Camara do Ki Ancestral\n\n"+
			"Custo de entrada: %d Energia e %d Gold\n"+
			"Nivel maximo liberado: %d\n"+
			"Cada nivel exige 5 vitorias consecutivas.",
		ancestralKiChamberEntryEnergyCost,
		ancestralKiChamberEntryGoldCost,
		unlocked,
	)

	if run, runErr := s.repo.GetActiveRunByPlayer(int64(player.ID)); runErr == nil && run != nil {
		return view.GameResponse{
			Text:    "Run ativa detectada.\n\n" + chamberRunProgressLine(run) + "\n\nUse 'Atacar Batalha Atual' para continuar.",
			MenuKey: "ANCESTRAL_KI_CHAMBER_RUN_MENU",
			Image: &view.Image{
				Type: view.ImageAsset,
				Path: "assets/maps/ki_ancestral_hall.png",
			},
		}, nil
	}

	return view.GameResponse{
		Text:    baseMsg,
		MenuKey: "ANCESTRAL_KI_CHAMBER_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/mestre_kalunga.png",
		},
	}, nil
}

func (s *AncestralKiChamberService) StartRun(userID int64, chamberLevel int) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if chamberLevel < 1 || chamberLevel > ancestralKiChamberMaxLevel {
		return chamberErrorResponse("Nivel da camara invalido."), nil
	}

	progress, err := s.repo.GetOrCreateProgress(int64(player.ID))
	if err != nil {
		return view.GameResponse{}, err
	}

	maxUnlocked := progress.HighestCompletedLevel + 1
	if maxUnlocked < 1 {
		maxUnlocked = 1
	}
	if chamberLevel > maxUnlocked {
		return chamberErrorResponse("Este nivel ainda esta bloqueado. Conclua os niveis anteriores."), nil
	}

	run, err := s.repo.StartRunIfEligible(
		int64(player.ID),
		chamberLevel,
		ancestralKiChamberEntryEnergyCost,
		ancestralKiChamberEntryGoldCost,
		time.Now().UTC(),
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrAncestralKiChamberRunActive):
			return chamberErrorResponse("Voce ja possui uma run ativa na camara."), nil
		case errors.Is(err, domain.ErrAncestralKiChamberInsufficientResource):
			return chamberErrorResponse("Recursos insuficientes para entrar na camara."), nil
		}
		return view.GameResponse{}, err
	}

	return view.GameResponse{
		Text: fmt.Sprintf(
			"Voce entrou na Camara do Ki Ancestral!\n\n"+
				"%s\n\n"+
				"Use Continuar para iniciar a batalha.",
			chamberRunProgressLine(run),
		),
		MenuKey: "ANCESTRAL_KI_CHAMBER_RUN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/maps/ki_ancestral_hall.png",
		},
	}, nil
}

func (s *AncestralKiChamberService) ContinueRun(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	run, err := s.repo.GetActiveRunByPlayer(int64(player.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chamberErrorResponse("Voce nao possui run ativa na camara."), nil
		}
		return view.GameResponse{}, err
	}

	monster, monsterLevel, err := s.pickMonsterForRun(player, run)
	if err != nil {
		return chamberErrorResponse("Nao foi possivel encontrar um adversario para esta batalha."), nil
	}

	stats, statsErr := calculateEffectivePlayerStats(player, s.inventoryRepo)
	if statsErr != nil {
		stats = effectivePlayerStats{
			ATK:       player.ATK,
			DEF:       player.DEF,
			Speed:     player.Speed,
			MaxHP:     player.MaxHP,
			MaxEnergy: player.MaxEnergy,
		}
	}

	playerHP, monsterHP, totalMonsterDamage := simulateChamberFight(
		player.HP,
		stats.ATK,
		stats.DEF,
		stats.Speed,
		monster.HP,
		monster.ATK,
		s.rng,
	)
	player.HP = playerHP
	if player.HP <= 0 {
		player.HP = 1
	}
	if err := s.playerRepo.Update(player); err != nil {
		return view.GameResponse{}, err
	}

	won := playerHP > 0 && monsterHP <= 0
	xpReward := clamp(CalculateXPReward(player.Level, monsterLevel), monster.MinXP, monster.MaxXP)
	goldReward := randBetweenInclusive(s.rng, monster.MinMoney, monster.MaxMoney)
	resultResp, err := s.ResolveBattleResult(userID, int64(monster.MonsterID), monsterLevel, won, xpReward, goldReward)
	if err != nil {
		return view.GameResponse{}, err
	}

	prefix := fmt.Sprintf(
		"Camara - Batalha %d/%d (Nivel %d)\n\nInimigo: %s (Lv %d)\nDano recebido: %d\nHP atual: %d/%d\n\n",
		run.CurrentBattle,
		ancestralKiChamberTotalBattles,
		run.ChamberLevel,
		monster.Name,
		monsterLevel,
		totalMonsterDamage,
		player.HP,
		player.MaxHP,
	)
	resultResp.Text = prefix + resultResp.Text
	if monster.Tag != "" {
		resultResp.Image = &view.Image{
			Type: view.ImageAsset,
			Path: fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.RegionID, monster.Tag),
		}
	}
	return resultResp, nil
}

func (s *AncestralKiChamberService) ResolveBattleResult(
	userID int64,
	monsterID int64,
	monsterLevel int,
	won bool,
	xpReward int,
	goldReward int,
) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	resolution, err := s.repo.ResolveBattleResult(
		int64(player.ID),
		monsterID,
		monsterLevel,
		won,
		xpReward,
		goldReward,
		ancestralKiChamberTotalBattles,
		ancestralKiChamberSuccessPct,
		time.Now().UTC(),
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chamberErrorResponse("Voce nao possui run ativa na camara."), nil
		}
		return view.GameResponse{}, err
	}
	run := resolution.Run

	if !won {
		return chamberErrorResponse("Voce foi derrotado. A run da camara foi encerrada.\n\nTente novamente a partir da batalha 1."), nil
	}
	if run.Victories >= ancestralKiChamberTotalBattles {
		return view.GameResponse{
			Text: fmt.Sprintf(
				"Camara concluida!\n\n"+
					"Voce venceu as %d batalhas do nivel %d.\n"+
					"Recompensa de conclusao: %d XP e %d Gold (40%% do acumulado da run).\n"+
					"Acumulado da run: %d XP e %d Gold.\n"+
					"Proximo nivel desbloqueado.\n\nVolte ao menu da Camara para iniciar o proximo desafio.",
				ancestralKiChamberTotalBattles,
				run.ChamberLevel,
				resolution.XPPayout,
				resolution.GoldPayout,
				run.AccumulatedXP,
				run.AccumulatedGold,
			),
			MenuKey: "ANCESTRAL_KI_CHAMBER_MENU",
			Image: &view.Image{
				Type: view.ImageAsset,
				Path: "assets/npcs/mestre_kalunga.png",
			},
		}, nil
	}

	return view.GameResponse{
		Text: fmt.Sprintf(
			"Vitoria!\n\n%s\n\nAgora voce pode:\n- Continuar\n- Usar pocoes\n- Transformar\n- Sair da camara",
			chamberRunProgressLine(run),
		),
		MenuKey: "ANCESTRAL_KI_CHAMBER_POST_WIN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/maps/plain_lv1.png",
		},
	}, nil
}

func (s *AncestralKiChamberService) AbandonRun(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	_, err = s.repo.AbandonActiveRun(int64(player.ID), time.Now().UTC())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return chamberErrorResponse("Voce nao possui run ativa para abandonar."), nil
		}
		return view.GameResponse{}, err
	}

	return view.GameResponse{
		Text:    "Voce saiu da Camara do Ki Ancestral. O progresso da run atual foi perdido.",
		MenuKey: "ANCESTRAL_KI_CHAMBER_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/npcs/mestre_kalunga.png",
		},
	}, nil
}

func GetChamberMonsterLevelRange(playerLevel int, chamberLevel int) (int, int) {
	if chamberLevel <= 1 {
		min := playerLevel - 2
		if min < 1 {
			min = 1
		}
		return min, playerLevel + 2
	}
	return playerLevel, playerLevel + chamberLevel
}

func (s *AncestralKiChamberService) pickMonsterForRun(
	player *domain.Player,
	run *domain.AncestralKiChamberRun,
) (*domain.Monster, int, error) {
	if player == nil || run == nil {
		return nil, 0, gorm.ErrInvalidData
	}
	if s.monsterRepo == nil {
		return nil, 0, gorm.ErrInvalidData
	}
	monsters, err := s.monsterRepo.ListAll()
	if err != nil {
		return nil, 0, err
	}
	if len(monsters) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}
	allowedMonsters := make([]domain.Monster, 0, len(monsters))
	for _, m := range monsters {
		if m.Tag == monsterAurumKarunTag {
			continue
		}
		allowedMonsters = append(allowedMonsters, m)
	}
	if len(allowedMonsters) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	minLv, maxLv := GetChamberMonsterLevelRange(player.Level, run.ChamberLevel)
	inRangeMonsters := make([]domain.Monster, 0, len(allowedMonsters))
	for _, m := range allowedMonsters {
		baseLv := int(m.Level)
		if baseLv >= minLv && baseLv <= maxLv {
			inRangeMonsters = append(inRangeMonsters, m)
		}
	}
	if len(inRangeMonsters) == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	targetLevel := minLv
	if maxLv > minLv {
		targetLevel = minLv + s.rng.Intn((maxLv-minLv)+1)
	}

	closestIdx := 0
	closestDiff := int(^uint(0) >> 1)
	for i, m := range inRangeMonsters {
		diff := int(math.Abs(float64(int(m.Level) - targetLevel)))
		if diff < closestDiff {
			closestDiff = diff
			closestIdx = i
		}
	}

	scaled := scaleChamberMonster(inRangeMonsters[closestIdx], targetLevel)
	if scaled.RegionID == 0 {
		scaled.RegionID = player.RegionID
	}
	return &scaled, targetLevel, nil
}

func scaleChamberMonster(base domain.Monster, targetLevel int) domain.Monster {
	baseLevel := int(base.Level)
	if baseLevel < 1 {
		baseLevel = 1
	}
	if targetLevel < 1 {
		targetLevel = 1
	}

	ratio := float64(targetLevel) / float64(baseLevel)
	scaled := base
	scaled.Level = int64(targetLevel)
	scaled.HP = maxInt(10, int(math.Round(float64(base.HP)*ratio)))
	scaled.ATK = maxInt(1, int(math.Round(float64(base.ATK)*ratio)))
	scaled.DEF = maxInt(1, int(math.Round(float64(base.DEF)*ratio)))
	return scaled
}

func simulateChamberFight(
	playerHP, playerATK, playerDEF, playerSpeed int,
	monsterHP, monsterATK int,
	rng *rand.Rand,
) (int, int, int) {
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	totalMonsterDamage := 0
	dodgeChance := float64(playerSpeed) * 0.5
	if dodgeChance > 35 {
		dodgeChance = 35
	}

	for playerHP > 0 && monsterHP > 0 {
		if rng.Float64()*100 >= dodgeChance {
			monsterDamage := rng.Intn(monsterATK+(monsterATK/3)+1) - playerDEF
			if monsterDamage < 1 {
				monsterDamage = 1
			}
			playerHP -= monsterDamage
			totalMonsterDamage += monsterDamage
		}

		if playerHP <= 0 {
			break
		}

		playerDamage := playerATK + rng.Intn(maxInt(2, playerATK/3))
		if playerDamage < 1 {
			playerDamage = 1
		}
		monsterHP -= playerDamage
	}

	return playerHP, monsterHP, totalMonsterDamage
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func chamberRunProgressLine(run *domain.AncestralKiChamberRun) string {
	if run == nil {
		return "Nivel 1 | Batalha 1/5 | Vitorias 0/5"
	}
	return fmt.Sprintf(
		"Nivel %d | Batalha %d/%d | Vitorias %d/%d",
		run.ChamberLevel,
		run.CurrentBattle,
		ancestralKiChamberTotalBattles,
		run.Victories,
		ancestralKiChamberTotalBattles,
	)
}

func chamberErrorResponse(msg string) view.GameResponse {
	return view.GameResponse{
		Text:    msg,
		MenuKey: "ANCESTRAL_KI_CHAMBER_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/error.png",
		},
	}
}
