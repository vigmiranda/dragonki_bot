package application

import (
	"math/rand"
	"testing"

	"github.com/vigmiranda/dragonki/internal/domain"
)

// ---------- Fakes (test doubles) ----------

type fakeMonsterRepo struct {
	listByRegionFn    func(regionID int64) ([]domain.Monster, error)
	getByMonsterIDFn  func(monsterID int64) (*domain.Monster, error)
	getMonsterDropsFn func(monsterID int) ([]domain.MonsterDrop, error)
}

func (f *fakeMonsterRepo) ListByRegion(regionID int64) ([]domain.Monster, error) {
	return f.listByRegionFn(regionID)
}

func (f *fakeMonsterRepo) GetByMonsterID(monsterID int64) (*domain.Monster, error) {
	return f.getByMonsterIDFn(monsterID)
}

func (f *fakeMonsterRepo) GetMonsterDrops(monsterID int) ([]domain.MonsterDrop, error) {
	return f.getMonsterDropsFn(monsterID)
}

type fakeInventoryRepo struct {
	createOrUpdateCalls []struct {
		playerID int64
		itemID   uint
		amount   int
	}
	err error
}

func (f *fakeInventoryRepo) UpdateInventoryItem(item *domain.InventoryItem) error {
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
	return nil, 0, nil
}

func (f *fakeInventoryRepo) GetInventoryItemByUser(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
	return nil, nil
}

// ---------- Tests ----------

func TestCalculateXPReward_MinimumIs5(t *testing.T) {
	got := CalculateXPReward(999, 1) // diferença bem negativa => multiplicador baixo
	if got < 5 {
		t.Fatalf("expected xp >= 5, got %d", got)
	}
}

func TestCalculateXPReward_LevelDiffBuckets(t *testing.T) {
	tests := []struct {
		name         string
		playerLevel  int
		monsterLevel int
		wantMin      int
	}{
		{"sameLevel", 10, 10, 5},
		{"monsterMuchHigher", 1, 50, 5},
		{"monsterLower", 50, 1, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateXPReward(tt.playerLevel, tt.monsterLevel)
			if got < tt.wantMin {
				t.Fatalf("expected xp >= %d, got %d", tt.wantMin, got)
			}
		})
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
					DropRate:  0.0, // nunca dropa
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
					Item:      domain.Item{Name: "Poção"},
					DropRate:  1.0, // sempre dropa
					MinQty:    2,
					MaxQty:    2, // qty determinístico
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
	if name != "Poção" || qty != 2 {
		t.Fatalf("expected (Poção,2), got (%q,%d)", name, qty)
	}
	if len(invRepo.createOrUpdateCalls) != 1 {
		t.Fatalf("expected 1 inventory call, got %d", len(invRepo.createOrUpdateCalls))
	}
	call := invRepo.createOrUpdateCalls[0]
	if call.playerID != 777 || call.itemID != 200 || call.amount != 2 {
		t.Fatalf("unexpected inventory call: %+v", call)
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

	// comportamento atual: devolve o último item que caiu
	if name != "Item 2" || qty != 3 {
		t.Fatalf("expected (Item 2,3), got (%q,%d)", name, qty)
	}

	if len(invRepo.createOrUpdateCalls) != 2 {
		t.Fatalf("expected 2 inventory calls, got %d", len(invRepo.createOrUpdateCalls))
	}
}

// Observação:
// Dá para testar GetBattle/Battle/Escape também, mas aí precisamos de mais fakes
// (PlayerRepository/RegionRepository/BattleRepository) e montar Player/Region/Monster com cuidado
// pra não cair em caminhos aleatórios. Se você quiser, eu monto essa bateria também.
