package application

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type weeklyBossRepoMock struct {
	activeBoss      *domain.WeeklyBoss
	activeErr       error
	weekBoss        *domain.WeeklyBoss
	weekErr         error
	attackedToday   bool
	attackedErr     error
	contribution    *domain.WeeklyBossContribution
	contributionErr error
	createAttackErr error
	applyBoss       *domain.WeeklyBoss
	applyErr        error
	addContribErr   error
	rankRows        []domain.WeeklyBossRankingEntry
	rankErr         error
	createBossErr   error
	expireErr       error
	randomTemplate  *domain.WeeklyBossTemplate
	randomErr       error
	rewardByPos     map[int]*domain.WeeklyBossReward
	payoutCreated   bool
}

func (m *weeklyBossRepoMock) GetActive(at time.Time) (*domain.WeeklyBoss, error) {
	if m.activeErr != nil {
		return nil, m.activeErr
	}
	return m.activeBoss, nil
}
func (m *weeklyBossRepoMock) GetByID(id uint) (*domain.WeeklyBoss, error) { return m.activeBoss, nil }
func (m *weeklyBossRepoMock) GetByWeekKey(weekKey string) (*domain.WeeklyBoss, error) {
	if m.weekErr != nil {
		return nil, m.weekErr
	}
	return m.weekBoss, nil
}
func (m *weeklyBossRepoMock) Create(boss *domain.WeeklyBoss) error { return m.createBossErr }
func (m *weeklyBossRepoMock) ExpireActiveBefore(at time.Time) error {
	return m.expireErr
}
func (m *weeklyBossRepoMock) GetRandomTemplate() (*domain.WeeklyBossTemplate, error) {
	if m.randomErr != nil {
		return nil, m.randomErr
	}
	if m.randomTemplate != nil {
		return m.randomTemplate, nil
	}
	return &domain.WeeklyBossTemplate{
		Name:  "Astarion, o Senhor da Ruina",
		Tag:   "ASTARION",
		MaxHP: 2150000,
		ATK:   190,
	}, nil
}
func (m *weeklyBossRepoMock) ApplyDamage(id uint, damage int64, at time.Time) (*domain.WeeklyBoss, error) {
	if m.applyErr != nil {
		return nil, m.applyErr
	}
	if m.applyBoss != nil {
		return m.applyBoss, nil
	}
	return m.activeBoss, nil
}
func (m *weeklyBossRepoMock) HasPlayerAttackedOnDate(bossID uint, playerID int64, attackDate time.Time) (bool, error) {
	if m.attackedErr != nil {
		return false, m.attackedErr
	}
	return m.attackedToday, nil
}
func (m *weeklyBossRepoMock) CreateAttack(attack *domain.WeeklyBossAttack) error {
	return m.createAttackErr
}
func (m *weeklyBossRepoMock) GetContribution(bossID uint, playerID int64) (*domain.WeeklyBossContribution, error) {
	if m.contributionErr != nil {
		return nil, m.contributionErr
	}
	return m.contribution, nil
}
func (m *weeklyBossRepoMock) AddContributionDamage(bossID uint, playerID int64, damage int64, attackedAt time.Time) error {
	return m.addContribErr
}
func (m *weeklyBossRepoMock) ListTopContributions(bossID uint, limit int) ([]domain.WeeklyBossRankingEntry, error) {
	if m.rankErr != nil {
		return nil, m.rankErr
	}
	return m.rankRows, nil
}
func (m *weeklyBossRepoMock) GetRewardByPosition(position int) (*domain.WeeklyBossReward, error) {
	if m.rewardByPos != nil {
		if r, ok := m.rewardByPos[position]; ok && r != nil {
			return r, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (m *weeklyBossRepoMock) CreatePayoutIfNotExists(payout *domain.WeeklyBossPayout) (bool, error) {
	if payout != nil {
		payout.ID = 1
	}
	if m.payoutCreated {
		return false, nil
	}
	m.payoutCreated = true
	return true, nil
}
func (m *weeklyBossRepoMock) MarkPayoutPaid(payoutID uint, paidAt time.Time) error { return nil }
func (m *weeklyBossRepoMock) MarkPayoutFailed(payoutID uint, reason string) error  { return nil }
func (m *weeklyBossRepoMock) AddGoldToPlayer(playerID int64, gold int64) error     { return nil }

type weeklyBossPlayerRepoMock struct {
	player    *domain.Player
	getErr    error
	updateErr error
	updated   *domain.Player
}

func (m *weeklyBossPlayerRepoMock) GetOrCreate(userID int64) (*domain.Player, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.player, nil
}
func (m *weeklyBossPlayerRepoMock) Update(player *domain.Player) error {
	m.updated = player
	return m.updateErr
}
func (m *weeklyBossPlayerRepoMock) CountOnlineLast15Min() (int, error) { return 0, nil }
func (m *weeklyBossPlayerRepoMock) GetClassesByLvl(lvl int) ([]domain.Class, error) {
	return nil, nil
}
func (m *weeklyBossPlayerRepoMock) GetClassesByID(id int64) (*domain.Class, error) {
	return nil, nil
}

type weeklyBossInventoryRepoMock struct {
	calls []struct {
		userID int64
		itemID uint
		amount int
	}
	createErr error
}

func (m *weeklyBossInventoryRepoMock) UpdateInventoryItem(item *domain.InventoryItem) error {
	return nil
}
func (m *weeklyBossInventoryRepoMock) CreateOrUpdateInventory(userID int64, itemID uint, amount int) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.calls = append(m.calls, struct {
		userID int64
		itemID uint
		amount int
	}{userID: userID, itemID: itemID, amount: amount})
	return nil
}
func (m *weeklyBossInventoryRepoMock) ListItemsByUserAndType(playerID int64, itemType domain.ItemType, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	return nil, 0, nil
}
func (m *weeklyBossInventoryRepoMock) ListEquipmentItemsByUserAndSlot(playerID int64, slot domain.EquipmentSlot, page int, pageSize int) ([]domain.InventoryItem, int64, error) {
	return nil, 0, nil
}
func (m *weeklyBossInventoryRepoMock) CountEquipmentItemsBySlot(playerID int64) (map[domain.EquipmentSlot]int64, error) {
	return map[domain.EquipmentSlot]int64{}, nil
}
func (m *weeklyBossInventoryRepoMock) GetInventoryItemByUser(playerID int64, inventoryItemID uint) (*domain.InventoryItem, error) {
	return nil, nil
}
func (m *weeklyBossInventoryRepoMock) GetOrCreatePlayerEquipment(playerID int64) (*domain.PlayerEquipment, error) {
	return &domain.PlayerEquipment{PlayerID: playerID}, nil
}
func (m *weeklyBossInventoryRepoMock) UpdatePlayerEquipment(equipment *domain.PlayerEquipment) error {
	return nil
}

type weeklyBossQueueMock struct {
	payloads []BigBossRewardsTaskPayload
	err      error
}

func (m *weeklyBossQueueMock) EnqueueBigBossRewardDistribution(payload BigBossRewardsTaskPayload) error {
	if m.err != nil {
		return m.err
	}
	m.payloads = append(m.payloads, payload)
	return nil
}

func TestWeeklyBossService_Open_NoActiveBoss(t *testing.T) {
	repo := &weeklyBossRepoMock{activeErr: gorm.ErrRecordNotFound}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	resp, err := svc.Open(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != "MAIN_MENU" {
		t.Fatalf("expected MAIN_MENU, got %s", resp.MenuKey)
	}
}

func TestWeeklyBossService_Attack_InsufficientResources(t *testing.T) {
	repo := &weeklyBossRepoMock{activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 1000, MaxHP: 1000}}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, Energy: 1, Gold: 10, HP: 100, ATK: 20, DEF: 10, Speed: 10}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	resp, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != "BIG_BOSS_MENU" {
		t.Fatalf("expected BIG_BOSS_MENU, got %s", resp.MenuKey)
	}
}

func TestWeeklyBossService_Attack_AlreadyAttackedToday(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss:    &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		attackedToday: true,
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, Energy: 50, Gold: 1000, HP: 100, ATK: 20, DEF: 10, Speed: 10}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	resp, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != "BIG_BOSS_MENU" {
		t.Fatalf("expected BIG_BOSS_MENU, got %s", resp.MenuKey)
	}
}

func TestWeeklyBossService_Ranking_Success(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		rankRows: []domain.WeeklyBossRankingEntry{
			{PlayerID: 11, PlayerName: "Vitor", TotalDamage: 1000},
			{PlayerID: 12, PlayerName: "Kael", TotalDamage: 900},
		},
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	resp, err := svc.Ranking()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != "BIG_BOSS_RANKING_MENU" {
		t.Fatalf("expected BIG_BOSS_RANKING_MENU, got %s", resp.MenuKey)
	}
}

func TestWeeklyBossService_EnsureCurrentWeeklyBoss_CreatesWhenMissing(t *testing.T) {
	repo := &weeklyBossRepoMock{
		weekErr: gorm.ErrRecordNotFound,
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	err := svc.EnsureCurrentWeeklyBoss(time.Date(2026, 3, 7, 10, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestWeeklyBossService_Attack_PropagatesCreateAttackError(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss:      &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		createAttackErr: errors.New("insert failed"),
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, Energy: 50, Gold: 1000, HP: 100, ATK: 20, DEF: 10, Speed: 10}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	_, err := svc.Attack(10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestWeeklyBossService_Attack_SetsPlayerHPToOne(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		applyBoss:  &domain.WeeklyBoss{ID: 1, Name: "Surtr", ATK: 100, CurrentHP: 900, MaxHP: 1000},
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, HP: 120, Energy: 50, Gold: 1000, ATK: 20, DEF: 10, Speed: 10}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	_, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if playerRepo.updated == nil {
		t.Fatalf("expected player update")
	}
	if playerRepo.updated.HP != 1 {
		t.Fatalf("expected HP to be 1 after big boss duel, got %d", playerRepo.updated.HP)
	}
}

func TestWeeklyBossService_Attack_DistributesRewardsWhenBossDefeated(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		applyBoss:  &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 0, MaxHP: 1000, Status: domain.WeeklyBossStatusDefeated},
		rankRows: []domain.WeeklyBossRankingEntry{
			{PlayerID: 1, PlayerName: "Vitor", TotalDamage: 1200},
		},
		rewardByPos: map[int]*domain.WeeklyBossReward{
			1: {ID: 1, Position: 1, Gold: 1000, Enabled: true},
		},
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, HP: 120, Energy: 50, Gold: 1000, ATK: 20, DEF: 10, Speed: 10}}
	svc := NewWeeklyBossService(repo, playerRepo, nil, nil)

	resp, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if resp.MenuKey != "BIG_BOSS_MENU" {
		t.Fatalf("expected BIG_BOSS_MENU, got %s", resp.MenuKey)
	}
	if !strings.Contains(resp.Text, "Premiacao do Top 3 distribuida com sucesso") {
		t.Fatalf("expected reward success message, got %q", resp.Text)
	}
}

func TestWeeklyBossService_Attack_DistributesRewardItemsToInventory(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		applyBoss:  &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 0, MaxHP: 1000, Status: domain.WeeklyBossStatusDefeated},
		rankRows: []domain.WeeklyBossRankingEntry{
			{PlayerID: 99, PlayerName: "Top1", TotalDamage: 5000},
		},
		rewardByPos: map[int]*domain.WeeklyBossReward{
			1: {
				ID:       1,
				Position: 1,
				Gold:     1000,
				Enabled:  true,
				Items: []domain.WeeklyBossRewardItem{
					{RewardID: 1, ItemID: 42, Quantity: 2},
					{RewardID: 1, ItemID: 51, Quantity: 1},
				},
			},
		},
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, HP: 120, Energy: 50, Gold: 1000, ATK: 20, DEF: 10, Speed: 10}}
	invRepo := &weeklyBossInventoryRepoMock{}
	svc := NewWeeklyBossService(repo, playerRepo, invRepo, nil)

	_, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if len(invRepo.calls) != 2 {
		t.Fatalf("expected 2 inventory reward calls, got %d", len(invRepo.calls))
	}
	if invRepo.calls[0].userID != 99 || invRepo.calls[0].itemID != 42 || invRepo.calls[0].amount != 2 {
		t.Fatalf("unexpected first reward call: %+v", invRepo.calls[0])
	}
	if invRepo.calls[1].userID != 99 || invRepo.calls[1].itemID != 51 || invRepo.calls[1].amount != 1 {
		t.Fatalf("unexpected second reward call: %+v", invRepo.calls[1])
	}
}

func TestWeeklyBossService_Attack_EnqueuesRewardWhenQueueIsConfigured(t *testing.T) {
	repo := &weeklyBossRepoMock{
		activeBoss: &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", ATK: 100, CurrentHP: 1000, MaxHP: 1000},
		applyBoss:  &domain.WeeklyBoss{ID: 1, Name: "Surtr", Tag: "SURTR", WeekKey: "2026-W10", ATK: 100, CurrentHP: 0, MaxHP: 1000, Status: domain.WeeklyBossStatusDefeated},
		rankRows: []domain.WeeklyBossRankingEntry{
			{PlayerID: 99, PlayerName: "Top1", TotalDamage: 5000},
		},
	}
	playerRepo := &weeklyBossPlayerRepoMock{player: &domain.Player{ID: 1, HP: 120, Energy: 50, Gold: 1000, ATK: 20, DEF: 10, Speed: 10}}
	queue := &weeklyBossQueueMock{}
	svc := NewWeeklyBossService(repo, playerRepo, nil, queue)

	resp, err := svc.Attack(10)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(queue.payloads) != 1 {
		t.Fatalf("expected 1 enqueued payload, got %d", len(queue.payloads))
	}
	if queue.payloads[0].BossID != 1 || queue.payloads[0].WeekKey != "2026-W10" {
		t.Fatalf("unexpected queue payload: %+v", queue.payloads[0])
	}
	if !strings.Contains(resp.Text, "em processamento") {
		t.Fatalf("expected queue processing message, got %q", resp.Text)
	}
}
