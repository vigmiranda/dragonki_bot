package application

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
)

type BattleService struct {
	playerRepo  ports.PlayerRepository
	monsterRepo ports.MonsterRepository
	regionRepo  ports.RegionRepository
	battleRepo  ports.BattleRepository
}

func NewBattleService(
	playerRepo ports.PlayerRepository,
	monsterRepo *postgres.MonsterPostgresRepo,
	regionRepo *postgres.RegionPostgres,
	battleRepo ports.BattleRepository,
) *BattleService {
	return &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		regionRepo:  regionRepo,
		battleRepo:  battleRepo,
	}
}

func (s *BattleService) GetBattle(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	recoveryStr := os.Getenv("ENERGY_REGEN_TIME")
	//todo tirar daqui
	player.SetTransformBase()
	err = s.playerRepo.Update(player)
	if err != nil {
		log.Fatal(err)
	}
	// ------------

	if player.Form == "Ki Awakened" {
		player = domain.KiAwakened(player)
	}

	msg := fmt.Sprintf("⚡ Energia esgotada. \n Aguarde %s para recarregar", recoveryStr)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_tired.png", domain.ConvertLowCase(player.Gender)),
	}
	menu := "MAIN_MENU"
	monsterID := 0
	if player.Energy > 0 {
		monsters, err := s.monsterRepo.ListByRegion(player.RegionID)
		if err != nil || len(monsters) == 0 {
			msg = fmt.Sprintf("Error ao encontra um monstro")
			img.Path = fmt.Sprintf("assets/error.png")
		}

		monster := pickMonsterByWeight(monsters)
		if monster == nil {
			msg = fmt.Sprintf("Error ao encontra um monstro")
			img.Path = fmt.Sprintf("assets/error.png")
		}

		msg = fmt.Sprintf(
			"🐉<b> %s apareceu! </b>\n\n"+
				"Você \n"+
				"💪 <b>Lv:</b> %d\n"+
				"❤️ <b>HP:</b> %d/%d\n"+
				"⚔️ <b>ATK:</b> %d  🛡️<b>DEF:</> %d \n\n\n"+
				"%s \n"+
				"💪 <b>Lv:</b> %d\n"+
				"❤️ <b>HP:</b> %d\n"+
				"⚔️ <b>ATK:</b> %d  🛡️<b>DEF:</> %d\n",
			monster.Name,
			player.Level,
			player.HP,
			player.MaxHP,
			player.ATK,
			player.DEF,
			monster.Name,
			monster.Level,
			monster.HP,
			monster.ATK,
			monster.DEF,
		)
		img.Path = fmt.Sprintf("assets/monsters/map_%d/%s.png", monster.RegionID, monster.Tag)
		monsterID = monster.MonsterID
		menu = "BATTLE_MENU"
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
		log.Fatal(err)
	}
	levelMulti, err := strconv.ParseInt(os.Getenv("LEVEL_MULT"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	monsterIDStr, _ := params["monsterID"]
	monsterId, err := strconv.ParseInt(monsterIDStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	monster, err := s.monsterRepo.GetByMonsterID(monsterId)
	if err != nil {
		log.Fatal(err)
	}

	region, err := s.regionRepo.GetByID(player.RegionID)
	if err != nil {
		log.Fatal(err)
	}

	msg := fmt.Sprintf("💀 Ocorreu um erro ao fim da batalha!!!")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/error.png"),
	}
	menu := "MAIN_MENU"

	resultBattle := "LOSE"

	monsterHP := monster.HP
	playerHp := player.HP
	player.Energy--
	atkTransform := 0
	defTransform := 0
	skin := ""
	if player.Form == "Ki Awakened" {
		atkTransform = int(float64(player.ATK) * 0.2)
		defTransform = int(float64(player.ATK) * 0.2)
		skin = "_ki_awakened"
	}
	atkTotalPlayer := 0
	//logica de dano
	for monsterHP > 0 && player.HP > 0 {
		monsterHP -= player.ATK + atkTransform
		atkTotalPlayer += player.ATK + atkTransform
		atkMonster := rand.Intn(monster.ATK+(monster.ATK/2)+1) - (player.DEF + defTransform)
		if atkMonster < 0 {
			atkMonster = 1
		}
		player.HP -= atkMonster
	}

	xpReceived := rand.Intn(monster.MaxXP) + 1
	moneyReceived := rand.Intn(monster.MaxMoney) + 1

	if player.HP > 0 {
		googNews := ""

		player.XP += xpReceived
		player.Gold += moneyReceived * 2

		if player.XP >= player.Level*int(levelMulti) {
			player.Level++
			player.ATK += int(float64(player.ATK) * 0.2)
			player.DEF += int(float64(player.DEF) * 0.2)
			player.MaxHP += int(float64(player.ATK) * 0.1)
			player.Energy = player.MaxEnergy
			player.XP = 0
			player.FormBattlesLeft = 0
			player.Form = "Base"
			player.HP = player.MaxHP
			googNews = "⬆️ <b>Subiu de nível!</b> \n"
			if player.FormBattlesLeft > 0 {
				player.FormBattlesLeft = 0
				player.Form = "Base"
			}

		}
		if player.FormBattlesLeft > 0 {
			player.FormBattlesLeft--
		}
		resultBattle = "WIN"
		msg = fmt.Sprintf(
			"🏆 Vitória! %s \n"+
				"%s"+
				"Você recebeu de dano: %d\n"+
				"⚔️ +%d XP\n"+
				"💰 +%d Gold - 🎉 Evento 2x 🎉\n"+
				"❤️ HP: %d/%d\n"+
				"⚡️ Energia: %d/%d\n"+
				"💪 Level: %d\n",
			player.Name,
			googNews,
			playerHp-player.HP,
			xpReceived,
			moneyReceived*2,
			player.HP,
			player.MaxHP,
			player.Energy,
			player.MaxEnergy,
			player.Level,
		)
		img.Path = fmt.Sprintf("assets/skins/%s%s_battle_win.png", domain.ConvertLowCase(player.Gender), skin)

	} else {
		if player.FormBattlesLeft > 0 {
			player.FormBattlesLeft = 0
			player.Form = "Base"
		}
		oldXP := player.XP
		oldGold := player.Gold
		nerf := float64(region.Nerf) / 100.0

		xpReduction := math.Round(float64(player.XP) * nerf)
		goldReduction := math.Round(float64(player.Gold) * nerf)

		player.XP -= int(xpReduction)
		player.Gold -= int(goldReduction)
		player.HP = 1
		player.Form = "Base"
		msg = fmt.Sprintf(
			"💀 Derrota! %s\n"+
				"⚔️ -%d%% de XP = %d\n"+
				"💰 -%d%% de Gold = %d\n"+
				"⚡️ Energia: %d/%d\n"+
				"💪 Level: %d\n",
			player.Name,
			region.Nerf,
			oldXP-player.XP,
			region.Nerf,
			oldGold-player.Gold,
			player.Energy,
			player.MaxEnergy,
			player.Level,
		)
		img.Path = fmt.Sprintf("assets/skins/%s_battle_loss.png", domain.ConvertLowCase(player.Gender))
		xpReceived = player.XP - oldXP
		moneyReceived = player.Gold - oldGold
	}

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          monster.MonsterID,
		Result:             resultBattle,
		XPGained:           xpReceived,
		GoldGained:         moneyReceived,
		PlayerDamageTotal:  atkTotalPlayer,
		MonsterDamageTotal: playerHp - player.HP,
		CreatedAt:          time.Now(),
	}

	err = s.battleRepo.SaveBattleLog(battleLog)
	if err != nil {
		return view.GameResponse{}, err
	}

	if err := s.playerRepo.Update(player); err != nil {
		msg = fmt.Sprintf("💀 Ocorreu um erro ao fim da batalha!!!")
		img.Path = fmt.Sprintf("assets/error.png")

	}
	viewResponse := view.GameResponse{
		Text:    msg,
		Image:   img,
		MenuKey: menu,
	}

	return viewResponse, nil
}
func (s *BattleService) Escape(userID int64, params map[string]string) (view.GameResponse, *domain.Monster, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		log.Fatal(err)
	}
	monsterIDStr, _ := params["monsterID"]
	monsterID, err := strconv.ParseInt(monsterIDStr, 10, 64)
	monster, err := s.monsterRepo.GetByMonsterID(monsterID)
	if err != nil {
		log.Fatal(err)
	}
	region, err := s.regionRepo.GetByID(player.RegionID)
	if err != nil {
		log.Fatal(err)
	}

	result := "ESCAPE"
	xpResult := 0
	goldResult := 0

	msg := fmt.Sprintf("💀 Você conseguiu fugir do %s!!!", monster.Name)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s_skin_basic.png", domain.ConvertLowCase(player.Gender)),
	}
	menu := "MAIN_MENU"

	escapePlayer := rand.Intn(10)
	escapeMonster := rand.Intn(10)
	damage := 0
	if escapePlayer <= escapeMonster {
		damage = rand.Intn(monster.ATK/2+1) + monster.ATK/2
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

			menu = "ENERGY_MENU"
			msg = fmt.Sprintf(
				"💀 Você zerou seu HP!!! \n\n"+
					"%s não é derrota, é aprendizado. \n"+
					"Venha, guerreiro. \nO descanso do Sir Soneca curará suas feridas. \n\n"+
					"O mundo ainda precisa da sua força. \n",
				player.Name,
			)
			img.Path = fmt.Sprintf("assets/npcs/mestre_kalunga.png")

			result = "ESCAPE/LOSS"
		} else {
			menu = "BATTLE_MENU"
			msg = fmt.Sprintf(
				"💀 Você não conseguiu fugir!!! \n"+
					"❤️ Seu HP agora é %d/%d \n"+
					"⚡️ E sua Energia é %d/%d \n",
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
		XPGained:           xpResult,
		GoldGained:         goldResult,
		PlayerDamageTotal:  0,
		MonsterDamageTotal: damage,
		CreatedAt:          time.Now(),
	}

	err = s.battleRepo.SaveBattleLog(battleLog)
	if err != nil {
		return view.GameResponse{}, monster, err
	}

	err = s.playerRepo.Update(player)
	if err != nil {
		log.Fatal(err)
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

func pickMonsterByWeight(monsters []domain.Monster) *domain.Monster {
	totalWeight := 0
	for _, m := range monsters {
		if m.SpawnWeight > 0 {
			totalWeight += m.SpawnWeight
		}
	}

	if totalWeight == 0 {
		return nil
	}

	roll := rand.Intn(totalWeight)
	acc := 0

	for _, m := range monsters {
		acc += m.SpawnWeight
		if roll < acc {
			return &m
		}
	}

	return nil
}
