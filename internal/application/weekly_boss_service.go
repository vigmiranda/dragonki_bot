package application

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
	"gorm.io/gorm"
)

const (
	weeklyBossAttackEnergyCost = 5
	weeklyBossAttackGoldCost   = 200
)

type WeeklyBossService struct {
	bossRepo      ports.WeeklyBossRepository
	playerRepo    ports.PlayerRepository
	inventoryRepo ports.InventoryRepository
	rewardQueue   WeeklyBossRewardQueue
	vipService    *VIPService
	rng           *rand.Rand
}

func NewWeeklyBossService(
	bossRepo ports.WeeklyBossRepository,
	playerRepo ports.PlayerRepository,
	inventoryRepo ports.InventoryRepository,
	rewardQueue WeeklyBossRewardQueue,
) *WeeklyBossService {
	return &WeeklyBossService{
		bossRepo:      bossRepo,
		playerRepo:    playerRepo,
		inventoryRepo: inventoryRepo,
		rewardQueue:   rewardQueue,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *WeeklyBossService) SetVIPService(vipService *VIPService) {
	s.vipService = vipService
}

func (s *WeeklyBossService) EnsureCurrentWeeklyBoss(at time.Time) error {
	now := at.UTC()

	if err := s.bossRepo.ExpireActiveBefore(now); err != nil {
		return err
	}

	weekKey := isoWeekKey(now)
	_, err := s.bossRepo.GetByWeekKey(weekKey)
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	template, err := s.bossRepo.GetRandomTemplate()
	if err != nil {
		return err
	}

	startsAt := startOfISOWeekUTC(now)
	endsAt := startsAt.Add(7 * 24 * time.Hour).Add(-time.Nanosecond)

	boss := &domain.WeeklyBoss{
		Name:      template.Name,
		Tag:       template.Tag,
		WeekKey:   weekKey,
		MaxHP:     template.MaxHP,
		CurrentHP: template.MaxHP,
		ATK:       template.ATK,
		Status:    domain.WeeklyBossStatusActive,
		StartsAt:  startsAt,
		EndsAt:    endsAt,
	}
	return s.bossRepo.Create(boss)
}

func (s *WeeklyBossService) Open(userID int64) (view.GameResponse, error) {
	now := time.Now().UTC()
	if err := s.EnsureCurrentWeeklyBoss(now); err != nil {
		return view.GameResponse{}, err
	}

	boss, err := s.bossRepo.GetActive(now)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return view.GameResponse{}, err
	}
	if boss == nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return view.GameResponse{
			Text:    "Nao ha Big Boss ativo nesta semana.",
			Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga.png"},
			MenuKey: "MAIN_MENU",
		}, nil
	}

	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	attacksToday, err := s.bossRepo.CountPlayerAttacksOnDate(boss.ID, int64(player.ID), now)
	if err != nil {
		return view.GameResponse{}, err
	}
	attackLimit := s.dailyAttackLimit(int64(player.ID))
	canAttack := attacksToday < int64(attackLimit)

	totalDamage := int64(0)
	contribution, err := s.bossRepo.GetContribution(boss.ID, int64(player.ID))
	if err == nil && contribution != nil {
		totalDamage = contribution.TotalDamage
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return view.GameResponse{}, err
	}

	remaining := boss.EndsAt.Sub(now)
	if remaining < 0 {
		remaining = 0
	}

	text := fmt.Sprintf(
		"👑 <b>Big Boss da Semana</b>\n\n"+
			"Nome: %s\n"+
			"HP: %d / %d\n\n"+
			"Tempo restante: %s\n"+
			"Seu dano total: %d\n"+
			"Ataques hoje: %d/%d (%s)\n\n"+
			"Custo do ataque: -%d⚡ e -%d💰",
		boss.Name,
		boss.CurrentHP,
		boss.MaxHP,
		formatDuration(remaining),
		totalDamage,
		attacksToday,
		attackLimit,
		availabilityText(canAttack),
		weeklyBossAttackEnergyCost,
		weeklyBossAttackGoldCost,
	)

	return view.GameResponse{
		Text:    text,
		Image:   &view.Image{Type: view.ImageAsset, Path: weeklyBossImagePath(boss.Tag)},
		MenuKey: "BIG_BOSS_MENU",
	}, nil
}

func (s *WeeklyBossService) Attack(userID int64) (view.GameResponse, error) {
	now := time.Now().UTC()
	if err := s.EnsureCurrentWeeklyBoss(now); err != nil {
		return view.GameResponse{}, err
	}

	boss, err := s.bossRepo.GetActive(now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return view.GameResponse{
				Text:    "Nao ha Big Boss ativo nesta semana.",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
				MenuKey: "MAIN_MENU",
			}, nil
		}
		return view.GameResponse{}, err
	}

	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}

	if player.Energy < weeklyBossAttackEnergyCost || player.Gold < weeklyBossAttackGoldCost {
		return view.GameResponse{
			Text: fmt.Sprintf(
				"Recursos insuficientes para atacar o Big Boss.\nPrecisa de %d⚡ e %d💰.",
				weeklyBossAttackEnergyCost, weeklyBossAttackGoldCost,
			),
			Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga_insta.png"},
			MenuKey: "BIG_BOSS_MENU",
		}, nil
	}

	attacksToday, err := s.bossRepo.CountPlayerAttacksOnDate(boss.ID, int64(player.ID), now)
	if err != nil {
		return view.GameResponse{}, err
	}
	attackLimit := s.dailyAttackLimit(int64(player.ID))
	if attacksToday >= int64(attackLimit) {
		return view.GameResponse{
			Text:    "Voce ja atingiu seu limite diario de ataques no Big Boss.\nVolte amanha para atacar novamente.",
			Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga_insta.png"},
			MenuKey: "BIG_BOSS_MENU",
		}, nil
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

	damage := simulateWeeklyBossFight(player.HP, stats.ATK, stats.DEF, stats.Speed, boss.ATK, s.rng)

	attackDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if err := s.bossRepo.CreateAttack(&domain.WeeklyBossAttack{
		BossID:   boss.ID,
		PlayerID: int64(player.ID),
		AttackAt: attackDate,
		Damage:   damage,
	}); err != nil {
		if isDailyAttackUniqueViolation(err) {
			return view.GameResponse{
				Text:    "Voce ja atingiu seu limite diario de ataques no Big Boss.\nVolte amanha para atacar novamente.",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga_insta.png"},
				MenuKey: "BIG_BOSS_MENU",
			}, nil
		}
		return view.GameResponse{}, err
	}

	updatedBoss, err := s.bossRepo.ApplyDamage(boss.ID, damage, now)
	if err != nil {
		return view.GameResponse{}, err
	}

	if err := s.bossRepo.AddContributionDamage(boss.ID, int64(player.ID), damage, now); err != nil {
		return view.GameResponse{}, err
	}

	player.HP = 1
	player.Energy -= weeklyBossAttackEnergyCost
	player.Gold -= weeklyBossAttackGoldCost
	if err := s.playerRepo.Update(player); err != nil {
		return view.GameResponse{}, err
	}

	rewardMessage := ""
	if updatedBoss.Status == domain.WeeklyBossStatusDefeated {
		if s.rewardQueue != nil {
			enqueueErr := s.rewardQueue.EnqueueBigBossRewardDistribution(BigBossRewardsTaskPayload{
				BossID:      updatedBoss.ID,
				WeekKey:     updatedBoss.WeekKey,
				TriggeredAt: now,
			})
			if enqueueErr == nil {
				rewardMessage = "\n\nPremiacao do Top 3 em processamento."
			} else if rewardErr := s.distributeDefeatRewards(updatedBoss.ID, now); rewardErr != nil {
				rewardMessage = "\n\nPremiacao do Top 3 pendente. A equipe ira processar em instantes."
			} else {
				rewardMessage = "\n\nPremiacao do Top 3 distribuida com sucesso."
			}
		} else if rewardErr := s.distributeDefeatRewards(updatedBoss.ID, now); rewardErr != nil {
			rewardMessage = "\n\nPremiacao do Top 3 pendente. A equipe ira processar em instantes."
		} else {
			rewardMessage = "\n\nPremiacao do Top 3 distribuida com sucesso."
		}
	}

	totalDamage := damage
	contribution, err := s.bossRepo.GetContribution(boss.ID, int64(player.ID))
	if err == nil && contribution != nil {
		totalDamage = contribution.TotalDamage
	}

	return view.GameResponse{
		Text: fmt.Sprintf(
			"⚔️ Voce enfrentou %s!\n\n"+
				"💥 Dano causado: %d\n"+
				"🩸 HP restante do Boss: %d\n"+
				"💀 Voce foi derrotado.\n"+
				"📊 Seu dano total nesta semana: %d\n\n"+
				"Custo aplicado: -%d⚡ e -%d💰%s",
			boss.Name,
			damage,
			updatedBoss.CurrentHP,
			totalDamage,
			weeklyBossAttackEnergyCost,
			weeklyBossAttackGoldCost,
			rewardMessage,
		),
		Image:   &view.Image{Type: view.ImageAsset, Path: weeklyBossImagePath(boss.Tag)},
		MenuKey: "BIG_BOSS_MENU",
	}, nil
}

func (s *WeeklyBossService) Ranking() (view.GameResponse, error) {
	now := time.Now().UTC()
	if err := s.EnsureCurrentWeeklyBoss(now); err != nil {
		return view.GameResponse{}, err
	}

	boss, err := s.bossRepo.GetActive(now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return view.GameResponse{
				Text:    "Nao ha Big Boss ativo nesta semana.",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
				MenuKey: "MAIN_MENU",
			}, nil
		}
		return view.GameResponse{}, err
	}

	rows, err := s.bossRepo.ListTopContributions(boss.ID, 10)
	if err != nil {
		return view.GameResponse{}, err
	}

	var b strings.Builder
	b.WriteString("🏆 <b>Ranking Big Boss</b>\n\n")
	if len(rows) == 0 {
		b.WriteString("Nenhum dano registrado nesta semana.")
	} else {
		for i, row := range rows {
			fmt.Fprintf(&b, "%d. %s - %d", i+1, row.PlayerName, row.TotalDamage)
			if i < len(rows)-1 {
				b.WriteString("\n")
			}
		}
	}

	return view.GameResponse{
		Text:    b.String(),
		Image:   &view.Image{Type: view.ImageAsset, Path: weeklyBossImagePath(boss.Tag)},
		MenuKey: "BIG_BOSS_RANKING_MENU",
	}, nil
}

func simulateWeeklyBossFight(playerHP, playerATK, playerDEF, playerSpeed, bossATK int, rng *rand.Rand) int64 {
	if playerHP < 1 {
		playerHP = 1
	}
	if playerATK < 1 {
		playerATK = 1
	}
	if playerDEF < 0 {
		playerDEF = 0
	}
	if playerSpeed < 0 {
		playerSpeed = 0
	}

	bossDamage := bossATK - playerDEF
	if bossDamage < 1 {
		bossDamage = 1
	}

	dodgeChance := float64(playerSpeed) * 0.5
	if dodgeChance > 40 {
		dodgeChance = 40
	}

	totalDamage := int64(0)
	for playerHP > 0 {
		totalDamage += int64(playerATK)

		roll := rng.Float64() * 100
		if roll < dodgeChance {
			continue
		}

		playerHP -= bossDamage
	}
	return totalDamage
}

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "0h"
	}

	totalHours := int(math.Floor(d.Hours()))
	days := totalHours / 24
	hours := totalHours % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	return fmt.Sprintf("%dh", hours)
}

func availabilityText(available bool) string {
	if available {
		return "disponivel"
	}
	return "indisponivel"
}

func (s *WeeklyBossService) dailyAttackLimit(playerID int64) int {
	baseLimit := 1
	if s == nil || s.vipService == nil {
		return baseLimit
	}
	limit := s.vipService.GetBigBossDailyAttackLimit(playerID)
	if limit < 1 {
		return baseLimit
	}
	return limit
}

func isDailyAttackUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "uq_weekly_boss_attacks_daily") ||
		(strings.Contains(msg, "unique") && strings.Contains(msg, "attack_date"))
}

func isoWeekKey(at time.Time) string {
	year, week := at.UTC().ISOWeek()
	return fmt.Sprintf("%04d-W%02d", year, week)
}

func startOfISOWeekUTC(at time.Time) time.Time {
	t := at.UTC()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	daysBack := weekday - 1
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return start.AddDate(0, 0, -daysBack)
}

func weeklyBossImagePath(tag string) string {
	normalizedTag := strings.ToUpper(strings.TrimSpace(tag))
	if normalizedTag == "" {
		return "assets/capa.png"
	}
	return fmt.Sprintf("assets/monsters/big_boss/%s.png", normalizedTag)
}

func (s *WeeklyBossService) distributeDefeatRewards(bossID uint, at time.Time) error {
	rows, err := s.bossRepo.ListTopContributions(bossID, 3)
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		return nil
	}

	for idx, row := range rows {
		position := idx + 1
		reward, rewardErr := s.bossRepo.GetRewardByPosition(position)
		if rewardErr != nil {
			if errors.Is(rewardErr, gorm.ErrRecordNotFound) {
				continue
			}
			return rewardErr
		}

		payout := &domain.WeeklyBossPayout{
			BossID:   bossID,
			PlayerID: row.PlayerID,
			Position: position,
			GoldPaid: reward.Gold,
			Status:   domain.WeeklyBossPayoutPending,
		}

		created, createErr := s.bossRepo.CreatePayoutIfNotExists(payout)
		if createErr != nil {
			return createErr
		}
		if !created {
			continue
		}

		if reward.Gold > 0 {
			if goldErr := s.bossRepo.AddGoldToPlayer(row.PlayerID, reward.Gold); goldErr != nil {
				_ = s.bossRepo.MarkPayoutFailed(payout.ID, goldErr.Error())
				return goldErr
			}
		}

		for _, itemReward := range reward.Items {
			if itemReward.ItemID == 0 || itemReward.Quantity <= 0 {
				continue
			}
			if s.inventoryRepo == nil {
				err := errors.New("inventory repository not configured for item payouts")
				_ = s.bossRepo.MarkPayoutFailed(payout.ID, err.Error())
				return err
			}
			if invErr := s.inventoryRepo.CreateOrUpdateInventory(row.PlayerID, itemReward.ItemID, itemReward.Quantity); invErr != nil {
				_ = s.bossRepo.MarkPayoutFailed(payout.ID, invErr.Error())
				return invErr
			}
		}

		if paidErr := s.bossRepo.MarkPayoutPaid(payout.ID, at); paidErr != nil {
			return paidErr
		}
	}

	return nil
}

func (s *WeeklyBossService) ProcessRewardDistribution(bossID uint) error {
	if bossID == 0 {
		return errors.New("invalid boss id")
	}
	return s.distributeDefeatRewards(bossID, time.Now().UTC())
}
