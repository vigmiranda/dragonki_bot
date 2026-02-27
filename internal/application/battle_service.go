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
	playerRepo          ports.PlayerRepository
	monsterRepo         ports.MonsterRepository
	regionRepo          ports.RegionRepository
	battleRepo          ports.BattleRepository
	inventoryRepository ports.InventoryRepository
}

func NewBattleService(
	playerRepo ports.PlayerRepository,
	monsterRepo *postgres.MonsterPostgresRepo,
	regionRepo *postgres.RegionPostgres,
	battleRepo ports.BattleRepository,
	inventoryRepo ports.InventoryRepository,
) *BattleService {
	return &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          battleRepo,
		inventoryRepository: inventoryRepo,
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
		Path: fmt.Sprintf("assets/skins/%s/%s/tired.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
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

	if monster.Tag == "AURUM_KARUN" {
		return s.AurumKarunLuck(monster, region, player, err)
	}

	msg := fmt.Sprintf("💀 Ocorreu um erro ao fim da batalha!!!")
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/error.png"),
	}
	menu := "MAIN_MENU"

	resultBattle := "LOSE"

	monsterHP := monster.HP
	player.Energy--
	atkTransform := 0
	defTransform := 0
	skin := ""
	oldXP := player.XP
	oldGold := player.Gold

	if player.Form == "Ki Awakened" {
		atkTransform = int(float64(player.ATK) * 0.2)
		defTransform = int(float64(player.ATK) * 0.2)
		skin = "ki_awakened_"
	}
	atkTotalPlayer := 0
	atkTotalMonster := 0
	//logica de dano
	for monsterHP > 0 && player.HP > 0 {
		monsterHP -= player.ATK + atkTransform
		atkTotalPlayer += player.ATK + atkTransform
		atkMonster := rand.Intn(monster.ATK+(monster.ATK/2)+1) - (player.DEF + defTransform)
		if atkMonster < 0 {
			atkMonster = 1
		}
		player.HP -= atkMonster
		atkTotalMonster += atkMonster
	}

	xpReceived := rand.Intn(monster.MaxXP-monster.MinXP) + monster.MinXP
	moneyReceived := rand.Intn(monster.MaxMoney-monster.MinMoney) + monster.MinMoney

	if player.HP > 0 {
		googNews := ""
		itemName, qty, _ := s.HandleDrops(monster.MonsterID, int64(player.ID))
		player.XP += xpReceived
		player.Gold += moneyReceived

		if qty > 0 {
			googNews += fmt.Sprintf("🥳 Dropou %dx %s \n", qty, itemName)
		}

		if player.XP >= player.Level*int(levelMulti) {
			player.Level++
			player.ATK += int(float64(player.ATK) * player.Class.AtkBonus)
			player.DEF += int(float64(player.DEF) * player.Class.DefBonus)
			player.MaxHP += int(float64(player.MaxHP) * player.Class.HpBonus)
			player.Speed += int(float64(player.Speed) * player.Class.SpeedBonus)
			player.Energy = player.MaxEnergy
			player.XP = 0
			player.FormBattlesLeft = 0
			player.Form = "Base"
			player.HP = player.MaxHP
			googNews += "⬆️ <b>Subiu de nível!</b> \n"
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
				"💰 +%d Gold\n"+
				"❤️ HP: %d/%d\n"+
				"⚡️ Energia: %d/%d\n"+
				"💪 Level: %d\n",
			player.Name,
			googNews,
			atkTotalMonster,
			xpReceived,
			moneyReceived,
			player.HP,
			player.MaxHP,
			player.Energy,
			player.MaxEnergy,
			player.Level,
		)
		img.Path = fmt.Sprintf("assets/skins/%s/%s/%sbattle_win.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
			skin)

	} else {
		if player.FormBattlesLeft > 0 {
			player.FormBattlesLeft = 0
			player.Form = "Base"
		}
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
		img.Path = fmt.Sprintf("assets/skins/%s/%sbattle_loss.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin)

		xpReceived = player.XP - oldXP
		moneyReceived = player.Gold - oldGold
	}

	battleLog := &domain.BattleLog{
		PlayerID:           int64(player.ID),
		MonsterID:          monster.MonsterID,
		Result:             resultBattle,
		OldXP:              oldXP,
		XPGained:           xpReceived,
		XPCurrent:          player.XP,
		OldGold:            oldGold,
		GoldGained:         moneyReceived,
		GoldCurrent:        player.Gold,
		PlayerDamageTotal:  atkTotalPlayer,
		MonsterDamageTotal: atkTotalMonster,
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
	menuParams := make(map[string]any)
	if (player.Level == 15 || player.Level == 35 || player.Level == 50) && player.XP == 0 {
		msg, img, menu, menuParams, err = s.ChangeClass(player)
		if err != nil {
			return view.GameResponse{}, err
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

func (s *BattleService) AurumKarunLuck(monster *domain.Monster, region *domain.Region, player *domain.Player, err error) (view.GameResponse, error) {
	txt := fmt.Sprintf("\n\n"+
		"\n\n <b> %s foi mais rápido e escapou!!! </b>\n\n "+
		"\n\n", monster.Name)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/maps/%s.png", domain.ConvertLowCase(region.Tag)),
	}

	skin := ""

	if player.Form == "Ki Awakened" {
		skin = "ki_awakened_"
	}
	result := "AURUMKARUN/ESCAPE"

	oldGold := player.Gold

	escapePlayer := rand.Intn(10)
	escapeAurumKarun := rand.Intn(15)

	if escapePlayer >= escapeAurumKarun {
		txt = fmt.Sprintf(
			"🏆 Vitória! %s \n"+
				"💰 +%d Gold\n",
			player.Name,
			monster.MaxMoney,
		)
		img.Path = fmt.Sprintf("assets/skins/%s/%s/%sbattle_win.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
			skin)
		result = "AURUMKARUN/WIN"
		player.Gold += monster.MaxMoney
	}

	viewResponse := view.GameResponse{
		Text:    txt,
		Image:   img,
		MenuKey: "MAIN_MENU",
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

	err = s.battleRepo.SaveBattleLog(battleLog)
	if err != nil {
		return view.GameResponse{}, err
	}
	err = s.playerRepo.Update(player)
	if err != nil {
		return view.GameResponse{}, err
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
	oldXP := player.XP
	oldGold := player.Gold

	msg := fmt.Sprintf("💀 Você conseguiu fugir do %s!!!", monster.Name)
	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/skins/%s/%splayer.png", domain.ConvertLowCase(player.Gender), player.Class.Skin),
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
					"⚔️ -%d%% de XP = %d\n"+
					"💰 -%d%% de Gold = %d\n\n"+
					"%s não é derrota, é aprendizado. \n"+
					"Venha, guerreiro. \nO descanso do Sir Soneca curará suas feridas. \n\n"+
					"O mundo ainda precisa da sua força. \n",
				region.Nerf,
				int(xpReduction),
				region.Nerf,
				int(goldReduction),
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

func (s *BattleService) HandleDrops(monsterID int, playerID int64) (string, int, error) {
	var drops []domain.MonsterDrop
	itemName := ""
	itemQty := 0

	drops, err := s.monsterRepo.GetMonsterDrops(monsterID)
	if err != nil {
		return "", 0, err
	}

	for _, drop := range drops {
		roll := rand.Float64()

		if roll <= drop.DropRate {
			qty := drop.MinQty
			if drop.MaxQty > drop.MaxQty {
				qty = rand.Intn(drop.MaxQty-drop.MinQty+1) + drop.MinQty

			}
			err = s.inventoryRepository.CreateOrUpdateInventory(playerID, drop.ItemID, qty)
			itemQty = qty
			itemName = drop.Item.Name
			if err != nil {
				return "", 0, err
			}
		}

	}
	return itemName, itemQty, nil
}

func (s *BattleService) ChangeClass(player *domain.Player) (string, *view.Image, string, map[string]any, error) {
	classes, err := s.playerRepo.GetClassesByLvl(player.Level)
	if err != nil {
		log.Fatal(err)
	}

	menu := "CHOICE_CLASS"
	msg := fmt.Sprintf(
		"“Você chegou até aqui… e agora seu verdadeiro caminho começa.” \n\n" +
			"No nível 15, todo guerreiro deve decidir como sua energia irá se manifestar. \n" +
			"O Ki não é apenas poder… \n" +
			"É direção. É propósito. É identidade.\n" +
			"Você poderá trilhar um dos três caminhos: \n\n" +
			"🌪️<b> Wind Striker </b>— Velocidade, esquiva e ataques cortantes como o vento.\n" +
			"🔥<b> Ki Vanguarda </b>— Ofensiva intensa, explosões de energia e domínio do combate.\n" +
			"🛡️<b> Earth Sentinel </b>— Resistência absoluta, defesa impenetrável e força inabalável.\n\n" +
			"Escolha com sabedoria… \n" +
			"Pois após esta decisão, seu Ki jamais será o mesmo.\n\n" +
			"Qual será o seu legado em Aethoria?")

	img := &view.Image{
		Type: view.ImageAsset,
		Path: fmt.Sprintf("assets/npcs/mestre_kalunga.png"),
	}
	menuParams := map[string]any{
		"Classes": classes,
	}
	return msg, img, menu, menuParams, nil
}
