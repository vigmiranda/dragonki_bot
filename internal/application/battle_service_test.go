package application

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

// ---------- Fakes (test doubles) ----------

type fixedSource struct {
	values []int64
	idx    int
}

func (s *fixedSource) Int63() int64 {
	if len(s.values) == 0 {
		return 0
	}
	v := s.values[s.idx%len(s.values)]
	s.idx++
	if v < 0 {
		v = -v
	}
	return v & ((^int64(0)) >> 1)
}

func (s *fixedSource) Seed(seed int64) {
	s.idx = 0
}

type fakeMonsterRepo struct {
	listByRegionFn    func(regionID int64) ([]domain.Monster, error)
	getByMonsterIDFn  func(monsterID int64) (*domain.Monster, error)
	getMonsterDropsFn func(monsterID int) ([]domain.MonsterDrop, error)
}

func (f *fakeMonsterRepo) ListByRegion(regionID int64) ([]domain.Monster, error) {
	if f.listByRegionFn == nil {
		return nil, nil
	}
	return f.listByRegionFn(regionID)
}

func (f *fakeMonsterRepo) GetByMonsterID(monsterID int64) (*domain.Monster, error) {
	if f.getByMonsterIDFn == nil {
		return nil, nil
	}
	return f.getByMonsterIDFn(monsterID)
}

func (f *fakeMonsterRepo) GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error) {
	if f.getMonsterDropsFn == nil {
		return nil, nil
	}
	return f.getMonsterDropsFn(monsterID)
}

type fakeInventoryRepo struct {
	createOrUpdateCalls []struct {
		playerID int64
		itemID   uint
		amount   int
	}
	updateInventoryCalls   []*domain.InventoryItem
	updateFn               func(item *domain.InventoryItem) error
	listByTypeFn           func(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error)
	listEquipmentBySlotFn  func(playerID int64, slot domain.EquipmentSlot, page int, pageSize int) ([]domain.InventoryItem, int64, error)
	countEquipmentBySlotFn func(playerID int64) (map[domain.EquipmentSlot]int64, error)
	getByUserFn            func(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error)
	getOrCreateEquipFn     func(playerID int64) (*domain.PlayerEquipment, error)
	updateEquipFn          func(equipment *domain.PlayerEquipment) error
	updateEquipCalls       []*domain.PlayerEquipment
	err                    error
}

func (f *fakeInventoryRepo) UpdateInventoryItem(item *domain.InventoryItem) error {
	f.updateInventoryCalls = append(f.updateInventoryCalls, item)
	if f.updateFn != nil {
		return f.updateFn(item)
	}
	return nil
}

func (f *fakeInventoryRepo) CreateOrUpdateInventory(userID int64, itemID uint, amount int) error {
	f.createOrUpdateCalls = append(f.createOrUpdateCalls, struct {
		playerID int64
		itemID   uint
		amount   int
	}{playerID: userID, itemID: itemID, amount: amount})
	return f.err
}

func (f *fakeInventoryRepo) ListItemsByUserAndType(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	if f.listByTypeFn != nil {
		return f.listByTypeFn(playerID, itemType, page, pageSize)
	}
	return nil, 0, nil
}

func (f *fakeInventoryRepo) ListEquipmentItemsByUserAndSlot(playerID int64, slot domain.EquipmentSlot, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	if f.listEquipmentBySlotFn != nil {
		return f.listEquipmentBySlotFn(playerID, slot, page, pageSize)
	}
	return nil, 0, nil
}

func (f *fakeInventoryRepo) CountEquipmentItemsBySlot(playerID int64) (map[domain.EquipmentSlot]int64, error) {
	if f.countEquipmentBySlotFn != nil {
		return f.countEquipmentBySlotFn(playerID)
	}
	return map[domain.EquipmentSlot]int64{}, nil
}

func (f *fakeInventoryRepo) GetInventoryItemByUser(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
	if f.getByUserFn != nil {
		return f.getByUserFn(playerID, inventoryItemID)
	}
	return nil, nil
}

func (f *fakeInventoryRepo) GetOrCreatePlayerEquipment(playerID int64) (*domain.PlayerEquipment, error) {
	if f.getOrCreateEquipFn != nil {
		return f.getOrCreateEquipFn(playerID)
	}
	return &domain.PlayerEquipment{PlayerID: playerID}, nil
}

func (f *fakeInventoryRepo) UpdatePlayerEquipment(equipment *domain.PlayerEquipment) error {
	f.updateEquipCalls = append(f.updateEquipCalls, equipment)
	if f.updateEquipFn != nil {
		return f.updateEquipFn(equipment)
	}
	return nil
}

type fakePlayerRepo struct {
	getOrCreateFn        func(userID int64) (*domain.Player, error)
	updateFn             func(player *domain.Player) error
	countOnlineFn        func() (int, error)
	getClassesByLvlFn    func(lvl int) ([]domain.Class, error)
	getClassesByIDFn     func(id int64) (*domain.Class, error)
	updatedPlayers       []*domain.Player
	getOrCreateCallCount int
}

func (f *fakePlayerRepo) GetOrCreate(userID int64) (*domain.Player, error) {
	f.getOrCreateCallCount++
	if f.getOrCreateFn == nil {
		return nil, nil
	}
	return f.getOrCreateFn(userID)
}

func (f *fakePlayerRepo) Update(player *domain.Player) error {
	f.updatedPlayers = append(f.updatedPlayers, player)
	if f.updateFn == nil {
		return nil
	}
	return f.updateFn(player)
}

func (f *fakePlayerRepo) CountOnlineLast15Min() (int, error) {
	if f.countOnlineFn == nil {
		return 0, nil
	}
	return f.countOnlineFn()
}

func (f *fakePlayerRepo) GetClassesByLvl(lvl int) ([]domain.Class, error) {
	if f.getClassesByLvlFn == nil {
		return nil, nil
	}
	return f.getClassesByLvlFn(lvl)
}

func (f *fakePlayerRepo) GetClassesByID(id int64) (*domain.Class, error) {
	if f.getClassesByIDFn == nil {
		return nil, nil
	}
	return f.getClassesByIDFn(id)
}

type fakeRegionRepo struct {
	getByIDFn     func(regionID int64) (*domain.Region, error)
	listByLevelFn func(level int64) ([]domain.Region, error)
	getByIDCalls  int
	lastRegionID  int64
}

func (f *fakeRegionRepo) GetByID(regionID int64) (*domain.Region, error) {
	f.getByIDCalls++
	f.lastRegionID = regionID
	if f.getByIDFn == nil {
		return nil, nil
	}
	return f.getByIDFn(regionID)
}

func (f *fakeRegionRepo) ListByLevel(level int64) ([]domain.Region, error) {
	if f.listByLevelFn == nil {
		return nil, nil
	}
	return f.listByLevelFn(level)
}

type fakeBattleRepo struct {
	saveFn       func(battleLog *domain.BattleLog) error
	saveWithFn   func(battleLog *domain.BattleLog, player *domain.Player) error
	weeklyRankFn func(metric string, limit int) ([]domain.WeeklyBattleRank, error)
	logs         []*domain.BattleLog
	players      []*domain.Player
	saveCall     int
	saveWithCall int
}

func (f *fakeBattleRepo) SaveBattleLog(ctx context.Context, battleLog *domain.BattleLog) error {
	f.saveCall++
	f.logs = append(f.logs, battleLog)
	if f.saveFn == nil {
		return nil
	}
	return f.saveFn(battleLog)
}

func (f *fakeBattleRepo) SaveBattleLogAndPlayer(ctx context.Context, battleLog *domain.BattleLog, player *domain.Player) error {
	f.saveWithCall++
	f.logs = append(f.logs, battleLog)
	f.players = append(f.players, player)
	if f.saveWithFn == nil {
		return nil
	}
	return f.saveWithFn(battleLog, player)
}

func (f *fakeBattleRepo) GetWeeklyLeaderboard(metric string, limit int) ([]domain.WeeklyBattleRank, error) {
	if f.weeklyRankFn == nil {
		return nil, nil
	}
	return f.weeklyRankFn(metric, limit)
}

func basePlayer() *domain.Player {
	return &domain.Player{
		ID:        1,
		UserID:    99,
		Name:      "Tester",
		Gender:    "man",
		Class:     domain.Class{Skin: "basic", AtkBonus: 0.2, DefBonus: 0.2, HpBonus: 0.1, SpeedBonus: 0.1},
		Level:     1,
		XP:        0,
		HP:        30,
		MaxHP:     30,
		ATK:       10,
		DEF:       3,
		Speed:     0,
		Energy:    5,
		MaxEnergy: 10,
		Gold:      100,
		RegionID:  1,
		Form:      formBase,
	}
}

// ---------- Tests ----------

func TestNewBattleService_InitializesRNGAndConfig(t *testing.T) {
	resetAppConfigForTest()
	t.Setenv("ENERGY_REGEN_TIME", "9m")
	t.Setenv("LEVEL_MULT", "3")

	svc := NewBattleService(&fakePlayerRepo{}, nil, nil, &fakeBattleRepo{}, &fakeInventoryRepo{})
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
	if svc.rng == nil {
		t.Fatalf("expected rng initialized")
	}
	if svc.cfg.EnergyRegenTime != "9m" || svc.cfg.LevelMult != 3 {
		t.Fatalf("unexpected cfg: %+v", svc.cfg)
	}
}

func TestLoadBattleConfig_Defaults(t *testing.T) {
	cfg := loadBattleConfig(appConfig{EnergyRegenTime: "10m", LevelMult: 1})
	if cfg.LevelMult != 1 {
		t.Fatalf("expected LevelMult=1, got %d", cfg.LevelMult)
	}
	if cfg.EnergyRegenTime != "10m" {
		t.Fatalf("expected regen text 10m, got %q", cfg.EnergyRegenTime)
	}
}

func TestLoadBattleConfig_UsesConfiguredLevelMult(t *testing.T) {
	cfg := loadBattleConfig(appConfig{EnergyRegenTime: "10m", LevelMult: 3})
	if cfg.LevelMult != 3 {
		t.Fatalf("expected level mult 3, got %d", cfg.LevelMult)
	}
}

func TestClamp_Bounds(t *testing.T) {
	if got := clamp(1, 2, 5); got != 2 {
		t.Fatalf("expected clamp min 2, got %d", got)
	}
	if got := clamp(6, 2, 5); got != 5 {
		t.Fatalf("expected clamp max 5, got %d", got)
	}
	if got := clamp(4, 2, 5); got != 4 {
		t.Fatalf("expected value 4, got %d", got)
	}
}

func TestCalculateXPReward_MinimumIs5(t *testing.T) {
	got := CalculateXPReward(999, 1)
	if got < 5 {
		t.Fatalf("expected xp >= 5, got %d", got)
	}
}

func TestCalculateXPReward_LevelDiffBuckets(t *testing.T) {
	tests := []struct {
		name         string
		playerLevel  int
		monsterLevel int
		multiplier   float64
	}{
		{"diffLeMinus5", 15, 10, 0.25},
		{"diffLeMinus3", 13, 10, 0.5},
		{"diffLeMinus1", 11, 10, 0.8},
		{"diffEq0", 10, 10, 1.0},
		{"diffLe2", 8, 10, 1.15},
		{"diffLe4", 6, 10, 1.3},
		{"diffGt4", 4, 10, 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateXPReward(tt.playerLevel, tt.monsterLevel)
			xpBase := 25 * math.Pow(float64(tt.monsterLevel), 1.12)
			want := int(xpBase * tt.multiplier)
			if want < 5 {
				want = 5
			}
			if got != want {
				t.Fatalf("expected xp=%d, got %d", want, got)
			}
		})
	}
}

func TestGetBattle_PlayerRepoError(t *testing.T) {
	svc := &BattleService{
		playerRepo: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) {
				return nil, errors.New("db error")
			},
		},
	}

	_, err := svc.GetBattle(1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetBattle_NoEnergy_ShowsTired(t *testing.T) {
	player := basePlayer()
	player.Energy = 0
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return player, nil
		},
	}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: &fakeMonsterRepo{},
		cfg: battleConfig{
			EnergyRegenTime: "10m",
		},
		rng: rand.New(rand.NewSource(1)),
	}

	resp, err := svc.GetBattle(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MenuKey != menuMain {
		t.Fatalf("expected main menu, got %q", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Energia") {
		t.Fatalf("expected energy warning text, got %q", resp.Text)
	}
	if !strings.Contains(resp.Image.Path, "tired.png") {
		t.Fatalf("expected tired image, got %q", resp.Image.Path)
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one player update, got %d", len(playerRepo.updatedPlayers))
	}
}

func TestGetBattle_MonsterRepoError_ShowsErrorImage(t *testing.T) {
	player := basePlayer()
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	monsterRepo := &fakeMonsterRepo{
		listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
			return nil, errors.New("fail")
		},
	}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		cfg:         battleConfig{EnergyRegenTime: "10m"},
		rng:         rand.New(rand.NewSource(1)),
	}

	resp, err := svc.GetBattle(1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuMain {
		t.Fatalf("expected main menu, got %q", resp.MenuKey)
	}
	if resp.Image.Path != "assets/error.png" {
		t.Fatalf("expected error image, got %q", resp.Image.Path)
	}
}

func TestGetBattle_WeightedMonster_SetsBattleMenu(t *testing.T) {
	player := basePlayer()
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	monsterRepo := &fakeMonsterRepo{
		listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
			return []domain.Monster{
				{MonsterID: 1, Name: "M1", Tag: "M1", SpawnWeight: 0, RegionID: 1},
				{MonsterID: 2, Name: "M2", Tag: "M2", SpawnWeight: 10, RegionID: 1},
			}, nil
		},
	}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		cfg:         battleConfig{EnergyRegenTime: "10m"},
		rng:         rand.New(rand.NewSource(1)),
	}

	resp, err := svc.GetBattle(1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuBattle {
		t.Fatalf("expected battle menu, got %q", resp.MenuKey)
	}
	if resp.MenuParams["monsterID"] != 2 {
		t.Fatalf("expected monsterID=2, got %#v", resp.MenuParams["monsterID"])
	}
}

func TestGetBattle_WeightsZero_ShowsError(t *testing.T) {
	player := basePlayer()
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	monsterRepo := &fakeMonsterRepo{
		listByRegionFn: func(regionID int64) ([]domain.Monster, error) {
			return []domain.Monster{
				{MonsterID: 1, SpawnWeight: 0},
				{MonsterID: 2, SpawnWeight: 0},
			}, nil
		},
	}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		cfg:         battleConfig{EnergyRegenTime: "10m"},
		rng:         rand.New(rand.NewSource(1)),
	}

	resp, err := svc.GetBattle(1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuMain {
		t.Fatalf("expected main menu, got %q", resp.MenuKey)
	}
	if resp.Image.Path != "assets/error.png" {
		t.Fatalf("expected error image, got %q", resp.Image.Path)
	}
}

func TestBattle_MissingMonsterID(t *testing.T) {
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
	}
	svc := &BattleService{playerRepo: playerRepo}
	_, err := svc.Battle(1, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "monsterID") {
		t.Fatalf("expected missing monsterID error, got %v", err)
	}
}

func TestBattle_Win_PersistsPlayerAndBattleLog(t *testing.T) {
	player := basePlayer()
	player.XP = 10
	player.Gold = 20

	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{
				MonsterID: 1, Name: "Slime", Tag: "SLIME", RegionID: 1,
				HP: 1, ATK: 1, DEF: 0,
				MinXP: 7, MaxXP: 7, MinMoney: 3, MaxMoney: 3,
			}, nil
		},
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) { return nil, nil },
	}
	regionRepo := &fakeRegionRepo{
		getByIDFn: func(regionID int64) (*domain.Region, error) { return &domain.Region{RegionID: 1, Nerf: 10}, nil },
	}
	battleRepo := &fakeBattleRepo{}

	svc := &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          battleRepo,
		inventoryRepository: &fakeInventoryRepo{},
		rng:                 rand.New(rand.NewSource(1)),
	}

	resp, err := svc.Battle(99, map[string]string{"monsterID": "1"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(resp.Text, "Vit") {
		t.Fatalf("expected win text, got %q", resp.Text)
	}
	if player.Energy != 4 {
		t.Fatalf("expected energy decremented to 4, got %d", player.Energy)
	}
	if battleRepo.saveWithCall != 1 {
		t.Fatalf("expected one transactional battle save, got %d", battleRepo.saveWithCall)
	}
	if battleRepo.logs[0].Result != "WIN" {
		t.Fatalf("expected WIN log, got %q", battleRepo.logs[0].Result)
	}
	if len(battleRepo.players) != 1 {
		t.Fatalf("expected one player persisted with battle log, got %d", len(battleRepo.players))
	}
}

func TestBattle_Lose_AppliesNerfAndResetsHP(t *testing.T) {
	player := basePlayer()
	player.HP = 1
	player.DEF = -100
	player.XP = 10
	player.Gold = 10

	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{
				MonsterID: 7, Name: "Orc", Tag: "ORC", RegionID: 1,
				HP: 50, ATK: 10, DEF: 0, MinXP: 1, MaxXP: 1, MinMoney: 1, MaxMoney: 1,
			}, nil
		},
	}
	regionRepo := &fakeRegionRepo{
		getByIDFn: func(regionID int64) (*domain.Region, error) {
			return &domain.Region{RegionID: 1, Nerf: 50}, nil
		},
	}
	battleRepo := &fakeBattleRepo{}

	svc := &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          battleRepo,
		inventoryRepository: &fakeInventoryRepo{},
		rng:                 rand.New(rand.NewSource(1)),
	}

	resp, err := svc.Battle(99, map[string]string{"monsterID": "7"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.Contains(resp.Text, "Derrota") {
		t.Fatalf("expected defeat text, got %q", resp.Text)
	}
	if player.HP != 1 {
		t.Fatalf("expected HP fixed at 1, got %d", player.HP)
	}
	if player.XP != 5 || player.Gold != 5 {
		t.Fatalf("expected 50%% nerf on xp/gold, got xp=%d gold=%d", player.XP, player.Gold)
	}
	if battleRepo.logs[0].Result != "LOSE" {
		t.Fatalf("expected LOSE log, got %q", battleRepo.logs[0].Result)
	}
}

func TestBattle_Level15_TriggersChangeClassMenu(t *testing.T) {
	player := basePlayer()
	player.Level = 14
	player.XP = 0
	player.ATK = 100
	player.DEF = 100
	player.MaxHP = 100
	player.HP = 100

	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		getClassesByLvlFn: func(lvl int) ([]domain.Class, error) {
			return []domain.Class{{Name: "Wind", MinLevel: 15}}, nil
		},
	}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{
				MonsterID: 9, Name: "Boss", Tag: "BOSS", RegionID: 1,
				HP: 1, ATK: 1, DEF: 0,
				MinXP: 5000, MaxXP: 5000, MinMoney: 1, MaxMoney: 1,
			}, nil
		},
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) { return nil, nil },
	}
	regionRepo := &fakeRegionRepo{
		getByIDFn: func(regionID int64) (*domain.Region, error) { return &domain.Region{RegionID: 1, Nerf: 10}, nil },
	}

	svc := &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          &fakeBattleRepo{},
		inventoryRepository: &fakeInventoryRepo{},
		rng:                 rand.New(rand.NewSource(1)),
	}

	resp, err := svc.Battle(99, map[string]string{"monsterID": "9"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuChoiceCls {
		t.Fatalf("expected class choice menu, got %q", resp.MenuKey)
	}
	if player.Level != 15 {
		t.Fatalf("expected level up to 15, got %d", player.Level)
	}
	if player.XP != 0 {
		t.Fatalf("expected xp reset after level up, got %d", player.XP)
	}
}

func TestBattle_SaveBattleLogError_ReturnsError(t *testing.T) {
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
	}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{MonsterID: 1, Name: "M", Tag: "M", RegionID: 1, HP: 1, ATK: 1, MinXP: 1, MaxXP: 1, MinMoney: 1, MaxMoney: 1}, nil
		},
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) { return nil, nil },
	}
	regionRepo := &fakeRegionRepo{getByIDFn: func(regionID int64) (*domain.Region, error) {
		return &domain.Region{RegionID: 1, Nerf: 10}, nil
	}}
	battleRepo := &fakeBattleRepo{
		saveWithFn: func(battleLog *domain.BattleLog, player *domain.Player) error {
			return errors.New("log fail")
		},
	}

	svc := &BattleService{
		playerRepo:          playerRepo,
		monsterRepo:         monsterRepo,
		regionRepo:          regionRepo,
		battleRepo:          battleRepo,
		inventoryRepository: &fakeInventoryRepo{},
		rng:                 rand.New(rand.NewSource(1)),
	}

	_, err := svc.Battle(1, map[string]string{"monsterID": "1"})
	if err == nil {
		t.Fatalf("expected save battle log error")
	}
}

func TestEscape_MissingMonsterID(t *testing.T) {
	svc := &BattleService{
		playerRepo: &fakePlayerRepo{
			getOrCreateFn: func(userID int64) (*domain.Player, error) { return basePlayer(), nil },
		},
	}
	_, _, err := svc.Escape(1, map[string]string{})
	if err == nil || !strings.Contains(err.Error(), "monsterID") {
		t.Fatalf("expected missing monsterID error, got %v", err)
	}
}

func TestEscape_Success_MenuMain(t *testing.T) {
	player := basePlayer()
	playerRepo := &fakePlayerRepo{getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil }}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{MonsterID: 5, Name: "Wolf", Tag: "WOLF", RegionID: 1, ATK: 10}, nil
		},
	}
	regionRepo := &fakeRegionRepo{getByIDFn: func(regionID int64) (*domain.Region, error) {
		return &domain.Region{RegionID: 1, Nerf: 50}, nil
	}}
	battleRepo := &fakeBattleRepo{}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		regionRepo:  regionRepo,
		battleRepo:  battleRepo,
		rng:         rand.New(&fixedSource{values: []int64{9 << 32, 1 << 32}}),
	}

	resp, _, err := svc.Escape(1, map[string]string{"monsterID": "5"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuMain {
		t.Fatalf("expected main menu on successful escape, got %q", resp.MenuKey)
	}
	if battleRepo.logs[0].Result != "ESCAPE" {
		t.Fatalf("expected ESCAPE result, got %q", battleRepo.logs[0].Result)
	}
}

func TestEscape_FailButAlive_MenuBattle(t *testing.T) {
	player := basePlayer()
	player.HP = 20
	playerRepo := &fakePlayerRepo{getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil }}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{MonsterID: 5, Name: "Wolf", Tag: "WOLF", RegionID: 1, ATK: 10}, nil
		},
	}
	regionRepo := &fakeRegionRepo{getByIDFn: func(regionID int64) (*domain.Region, error) {
		return &domain.Region{RegionID: 1, Nerf: 50}, nil
	}}
	battleRepo := &fakeBattleRepo{}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		regionRepo:  regionRepo,
		battleRepo:  battleRepo,
		rng:         rand.New(&fixedSource{values: []int64{1 << 32, 9 << 32, 0 << 32}}),
	}

	resp, _, err := svc.Escape(1, map[string]string{"monsterID": "5"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuBattle {
		t.Fatalf("expected battle menu when fail and alive, got %q", resp.MenuKey)
	}
	if player.HP >= 20 {
		t.Fatalf("expected hp reduction, got %d", player.HP)
	}
}

func TestEscape_FailAndZeroHP_AppliesNerf(t *testing.T) {
	player := basePlayer()
	player.HP = 3
	player.XP = 10
	player.Gold = 10
	playerRepo := &fakePlayerRepo{getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil }}
	monsterRepo := &fakeMonsterRepo{
		getByMonsterIDFn: func(monsterID int64) (*domain.Monster, error) {
			return &domain.Monster{MonsterID: 5, Name: "Wolf", Tag: "WOLF", RegionID: 1, ATK: 10}, nil
		},
	}
	regionRepo := &fakeRegionRepo{getByIDFn: func(regionID int64) (*domain.Region, error) {
		return &domain.Region{RegionID: 1, Nerf: 50}, nil
	}}
	battleRepo := &fakeBattleRepo{}

	svc := &BattleService{
		playerRepo:  playerRepo,
		monsterRepo: monsterRepo,
		regionRepo:  regionRepo,
		battleRepo:  battleRepo,
		rng:         rand.New(&fixedSource{values: []int64{1 << 32, 9 << 32, 0 << 32}}),
	}

	resp, _, err := svc.Escape(1, map[string]string{"monsterID": "5"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != menuEnergy {
		t.Fatalf("expected energy menu after dropping to 0 HP, got %q", resp.MenuKey)
	}
	if player.HP != 1 {
		t.Fatalf("expected hp forced to 1, got %d", player.HP)
	}
	if player.XP != 5 || player.Gold != 5 {
		t.Fatalf("expected 50%% reduction on xp/gold, got xp=%d gold=%d", player.XP, player.Gold)
	}
	if battleRepo.logs[0].Result != "ESCAPE/LOSS" {
		t.Fatalf("expected ESCAPE/LOSS, got %q", battleRepo.logs[0].Result)
	}
}

func TestAurumKarunLuck_WinAndEscape(t *testing.T) {
	for _, tc := range []struct {
		name       string
		values     []int64
		wantResult string
		wantGold   int
	}{
		{name: "win", values: []int64{9 << 32, 1 << 32}, wantResult: "AURUMKARUN/WIN", wantGold: 80},
		{name: "escape", values: []int64{1 << 32, 14 << 32}, wantResult: "AURUMKARUN/ESCAPE", wantGold: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			player := basePlayer()
			player.Gold = 0
			battleRepo := &fakeBattleRepo{}
			playerRepo := &fakePlayerRepo{}
			svc := &BattleService{
				playerRepo: playerRepo,
				battleRepo: battleRepo,
				rng:        rand.New(&fixedSource{values: tc.values}),
			}

			monster := &domain.Monster{MonsterID: 99, Name: "Aurum", MaxMoney: 80}
			region := &domain.Region{Tag: "PLAINS"}

			_, err := svc.AurumKarunLuck(monster, region, player)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if battleRepo.logs[0].Result != tc.wantResult {
				t.Fatalf("expected result %q, got %q", tc.wantResult, battleRepo.logs[0].Result)
			}
			if player.Gold != tc.wantGold {
				t.Fatalf("expected gold %d, got %d", tc.wantGold, player.Gold)
			}
		})
	}
}

func TestChangeClass_ReturnsChoiceMenu(t *testing.T) {
	playerRepo := &fakePlayerRepo{
		getClassesByLvlFn: func(lvl int) ([]domain.Class, error) {
			return []domain.Class{{Name: "Wind", MinLevel: 15}}, nil
		},
	}
	svc := &BattleService{playerRepo: playerRepo}
	player := basePlayer()
	player.Level = 15

	_, _, menu, params, err := svc.ChangeClass(player)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if menu != menuChoiceCls {
		t.Fatalf("expected menu %q, got %q", menuChoiceCls, menu)
	}
	if _, ok := params["Classes"]; !ok {
		t.Fatalf("expected Classes in menu params")
	}
}

func TestPickMonsterByWeight_ReturnsNilWhenTotalWeightZero(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	monsters := []domain.Monster{
		{MonsterID: 1, SpawnWeight: 0},
		{MonsterID: 2, SpawnWeight: 0},
	}

	got := pickMonsterByWeight(rng, monsters)
	if got != nil {
		t.Fatalf("expected nil, got %+v", *got)
	}
}

func TestPickMonsterByWeight_ReturnsTheOnlyWeightedMonster(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	monsters := []domain.Monster{
		{MonsterID: 1, SpawnWeight: 0},
		{MonsterID: 2, SpawnWeight: 10},
		{MonsterID: 3, SpawnWeight: 0},
	}

	got := pickMonsterByWeight(rng, monsters)
	if got == nil {
		t.Fatalf("expected non-nil monster")
	}
	if got.MonsterID != 2 {
		t.Fatalf("expected monsterID=2, got %d", got.MonsterID)
	}
}

func TestRandBetweenInclusive_SwapsWhenMaxLessThanMin(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	for i := 0; i < 50; i++ {
		got := randBetweenInclusive(rng, 10, 3)
		if got < 3 || got > 10 {
			t.Fatalf("expected 3..10, got %d", got)
		}
	}
}

func TestRandBetweenInclusive_ReturnsMinWhenEqual(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	got := randBetweenInclusive(rng, 7, 7)
	if got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}
}

func TestHandleDrops_NoDrop_ReturnsZeroAndDoesNotTouchInventory(t *testing.T) {
	monsterRepo := &fakeMonsterRepo{
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) {
			return []domain.MonsterDrop{
				{
					MonsterID: monsterID,
					ItemID:    100,
					Item:      domain.Item{Name: "Item A"},
					DropRate:  0.0,
					MinQty:    1,
					MaxQty:    1,
				},
			}, nil
		},
	}
	invRepo := &fakeInventoryRepo{}

	svc := &BattleService{
		monsterRepo:         monsterRepo,
		inventoryRepository: invRepo,
		rng:                 rand.New(rand.NewSource(1)),
	}

	name, qty, err := svc.HandleDrops(123, 999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "" || qty != 0 {
		t.Fatalf("expected empty result, got name=%q qty=%d", name, qty)
	}
	if len(invRepo.createOrUpdateCalls) != 0 {
		t.Fatalf("expected 0 inventory calls, got %d", len(invRepo.createOrUpdateCalls))
	}
}

func TestHandleDrops_DropGuaranteed_CreatesInventoryAndReturnsItem(t *testing.T) {
	monsterRepo := &fakeMonsterRepo{
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) {
			return []domain.MonsterDrop{
				{
					MonsterID: monsterID,
					ItemID:    200,
					Item:      domain.Item{Name: "Potion"},
					DropRate:  1.0,
					MinQty:    2,
					MaxQty:    2,
				},
			}, nil
		},
	}
	invRepo := &fakeInventoryRepo{}

	svc := &BattleService{
		monsterRepo:         monsterRepo,
		inventoryRepository: invRepo,
		rng:                 rand.New(rand.NewSource(1)),
	}

	name, qty, err := svc.HandleDrops(55, 777)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Potion" || qty != 2 {
		t.Fatalf("expected (Potion,2), got (%q,%d)", name, qty)
	}
	if len(invRepo.createOrUpdateCalls) != 1 {
		t.Fatalf("expected 1 inventory call, got %d", len(invRepo.createOrUpdateCalls))
	}
	call := invRepo.createOrUpdateCalls[0]
	if call.playerID != 777 || call.itemID != 200 || call.amount != 2 {
		t.Fatalf("unexpected inventory call: %+v", call)
	}
}

func TestHandleDrops_GetMonsterDropsError(t *testing.T) {
	monsterRepo := &fakeMonsterRepo{
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) {
			return nil, errors.New("drop list failed")
		},
	}
	svc := &BattleService{
		monsterRepo:         monsterRepo,
		inventoryRepository: &fakeInventoryRepo{},
		rng:                 rand.New(rand.NewSource(1)),
	}

	_, _, err := svc.HandleDrops(1, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestHandleDrops_CreateOrUpdateInventoryError(t *testing.T) {
	monsterRepo := &fakeMonsterRepo{
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) {
			return []domain.MonsterDrop{
				{
					MonsterID: monsterID,
					ItemID:    100,
					Item:      domain.Item{Name: "Item A"},
					DropRate:  1.0,
					MinQty:    1,
					MaxQty:    1,
				},
			}, nil
		},
	}
	invRepo := &fakeInventoryRepo{err: errors.New("inventory fail")}
	svc := &BattleService{
		monsterRepo:         monsterRepo,
		inventoryRepository: invRepo,
		rng:                 rand.New(rand.NewSource(1)),
	}

	_, _, err := svc.HandleDrops(1, 1)
	if err == nil {
		t.Fatalf("expected inventory error")
	}
}

func TestHandleDrops_MultipleDrops_ReturnsLastDroppedItem(t *testing.T) {
	monsterRepo := &fakeMonsterRepo{
		getMonsterDropsFn: func(monsterID int) ([]domain.MonsterDrop, error) {
			return []domain.MonsterDrop{
				{
					MonsterID: monsterID,
					ItemID:    1,
					Item:      domain.Item{Name: "Item 1"},
					DropRate:  1.0,
					MinQty:    1,
					MaxQty:    1,
				},
				{
					MonsterID: monsterID,
					ItemID:    2,
					Item:      domain.Item{Name: "Item 2"},
					DropRate:  1.0,
					MinQty:    3,
					MaxQty:    3,
				},
			}, nil
		},
	}
	invRepo := &fakeInventoryRepo{}

	svc := &BattleService{
		monsterRepo:         monsterRepo,
		inventoryRepository: invRepo,
		rng:                 rand.New(rand.NewSource(1)),
	}

	name, qty, err := svc.HandleDrops(9, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if name != "Item 2" || qty != 3 {
		t.Fatalf("expected (Item 2,3), got (%q,%d)", name, qty)
	}

	if len(invRepo.createOrUpdateCalls) != 2 {
		t.Fatalf("expected 2 inventory calls, got %d", len(invRepo.createOrUpdateCalls))
	}
}
