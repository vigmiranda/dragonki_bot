package application

import (
	"fmt"
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
}

type BattleService struct {
	playerRepo          ports.PlayerRepository
	monsterRepo         ports.MonsterRepository
	regionRepo          ports.RegionRepository
	battleRepo          ports.BattleRepository
	inventoryRepository ports.InventoryRepository

	cfg battleConfig
	rng *rand.Rand
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
		cfg:                 loadBattleConfigFromEnv(),
		rng:                 rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func loadBattleConfigFromEnv() battleConfig {
	cfg := battleConfig{
		LevelMult:       1,
		EnergyRegenTime: os.Getenv("ENERGY_REGEN_TIME"),
	}

	if cfg.EnergyRegenTime == "" {
		cfg.EnergyRegenTime = "alguns instantes"
	}

	if v := os.Getenv("LEVEL_MULT"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			cfg.LevelMult = parsed
		}
	}

	return cfg
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

	msg := fmt.Sprintf("⚡ Energia esgotada. \n Aguarde %s para recarregar", s.cfg.EnergyRegenTime)
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

	monsterIDStr, ok := params["monsterID"]
	if !ok || monsterIDStr == "" {
		return view.GameResponse{}, fmt.Errorf("missing param monsterID")
	}
	monsterId, err := strconv.ParseInt(monsterIDStr, 10, 64)
	if err != nil {
		return view.GameResponse{}, err
	}

	monster, err := s.monsterRepo.GetByMonsterID(monsterId)
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

	msg := "💀 Ocorreu um erro ao fim da batalha!!!"
	img := &view.Image{
		Type: view.ImageAsset,
		Path: "assets/error.png",
	}
	menu := menuMain

	resultBattle := "LOSE"

	monsterHP := monster.HP
	if player.Energy > 0 {
		player.Energy--
	}

	atkTransform := 0
	defTransform := 0
	skin := ""
	oldXP := player.XP
	oldGold := player.Gold

	if player.Form == formKiAwakened {
		atkTransform = int(float64(player.ATK) * 0.2)
		defTransform = int(float64(player.ATK) * 0.2)
		skin = "ki_awakened_"
	}

	atkTotalPlayer := 0
	atkTotalMonster := 0

	baseChange := float64(player.Speed) / (float64(player.Speed) + 100) * 35
	if baseChange > 35 {
		baseChange = 35
	}
	player.DodgeAccumulator += baseChange

	// lógica de dano (mantive igual, só usando rng injetado)
	for monsterHP > 0 && player.HP > 0 {
		roll := s.rng.Float64() * 100

		if roll < player.DodgeAccumulator {
			player.DodgeAccumulator = 0
			continue
		}

		monsterHP -= player.ATK + atkTransform
		atkTotalPlayer += player.ATK + atkTransform

		atkMonster := s.rng.Intn(monster.ATK+(monster.ATK/2)+1) - (player.DEF + defTransform)
		if atkMonster < 0 {
			atkMonster = 1
		}

		player.HP -= atkMonster
		atkTotalMonster += atkMonster
	}

	xpReceived := clamp(CalculateXPReward(player.Level, int(monster.Level)), monster.MinXP, monster.MaxXP) //randBetweenInclusive(s.rng, monster.MinXP, monster.MaxXP)
	moneyReceived := randBetweenInclusive(s.rng, monster.MinMoney, monster.MaxMoney)

	if player.HP > 0 {
		goodNews := ""

		itemName, qty, dropErr := s.HandleDrops(monster.MonsterID, int64(player.ID))
		if dropErr != nil {
			return view.GameResponse{}, dropErr
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
			xpReceived = int(XPLvlUp) - oldXP
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
			goodNews,
			atkTotalMonster,
			xpReceived,
			moneyReceived,
			player.HP,
			player.MaxHP,
			player.Energy,
			player.MaxEnergy,
			player.Level,
		)

		img.Path = fmt.Sprintf(
			"assets/skins/%s/%s/%sbattle_win.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
			skin,
		)
	} else {
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

		img.Path = fmt.Sprintf(
			"assets/skins/%s/%sbattle_loss.png",
			domain.ConvertLowCase(player.Gender),
			player.Class.Skin,
		)

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

	if err := s.battleRepo.SaveBattleLog(battleLog); err != nil {
		return view.GameResponse{}, err
	}

	if err := s.playerRepo.Update(player); err != nil {
		return view.GameResponse{}, err
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

func (s *BattleService) AurumKarunLuck(monster *domain.Monster, region *domain.Region, player *domain.Player) (view.GameResponse, error) {
	txt := fmt.Sprintf(
		"\n\n\n\n <b> %s foi mais rápido e escapou!!! </b>\n\n \n\n",
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
			"🏆 Vitória! %s \n"+
				"💰 +%d Gold\n",
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

	if err := s.battleRepo.SaveBattleLog(battleLog); err != nil {
		return view.GameResponse{}, err
	}
	if err := s.playerRepo.Update(player); err != nil {
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

	msg := fmt.Sprintf("💀 Você conseguiu fugir do %s!!!", monster.Name)
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
			img.Path = "assets/npcs/mestre_kalunga.png"

			result = "ESCAPE/LOSS"
		} else {
			menu = menuBattle
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

	if err := s.battleRepo.SaveBattleLog(battleLog); err != nil {
		return view.GameResponse{}, monster, err
	}

	if err := s.playerRepo.Update(player); err != nil {
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

func pickMonsterByWeight(rng *rand.Rand, monsters []domain.Monster) *domain.Monster {
	totalWeight := 0
	for i := range monsters {
		if monsters[i].SpawnWeight > 0 {
			totalWeight += monsters[i].SpawnWeight
		}
	}
	if totalWeight == 0 {
		return nil
	}

	roll := rng.Intn(totalWeight)
	acc := 0

	for i := range monsters {
		acc += monsters[i].SpawnWeight
		if roll < acc {
			return &monsters[i]
		}
	}

	return nil
}

func randBetweenInclusive(rng *rand.Rand, min, max int) int {
	if max < min {
		min, max = max, min
	}
	if max == min {
		return min
	}
	return rng.Intn(max-min+1) + min
}

func (s *BattleService) HandleDrops(monsterID int, playerID int64) (string, int, error) {
	drops, err := s.monsterRepo.GetMonsterDrops(monsterID)
	if err != nil {
		return "", 0, err
	}

	itemName := ""
	itemQty := 0

	for _, drop := range drops {
		roll := s.rng.Float64()

		if roll <= drop.DropRate {
			qty := drop.MinQty
			if drop.MaxQty > drop.MinQty {
				qty = s.rng.Intn(drop.MaxQty-drop.MinQty+1) + drop.MinQty
			}

			if err := s.inventoryRepository.CreateOrUpdateInventory(playerID, drop.ItemID, qty); err != nil {
				return "", 0, err
			}

			// Mantive o comportamento: se dropar mais de um item na mesma luta, devolve o "último" que caiu.
			itemQty = qty
			itemName = drop.Item.Name
		}
	}

	return itemName, itemQty, nil
}

func (s *BattleService) ChangeClass(player *domain.Player) (string, *view.Image, string, map[string]any, error) {
	classes, err := s.playerRepo.GetClassesByLvl(player.Level)
	if err != nil {
		return "", nil, "", nil, err
	}

	menu := menuChoiceCls
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
			"Qual será o seu legado em Aethoria?",
	)

	img := &view.Image{
		Type: view.ImageAsset,
		Path: "assets/npcs/mestre_kalunga.png",
	}
	menuParams := map[string]any{
		"Classes": classes,
	}

	return msg, img, menu, menuParams, nil
}
