package application

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
	"github.com/vigmiranda/dragonki/internal/ports"
	"gorm.io/gorm"
)

type VIPService struct {
	playerRepo    ports.PlayerRepository
	configRepo    ports.VIPConfigRepository
	playerVIPRepo ports.PlayerVIPRepository
}

const vipNamePrefix = "👑 "

func NewVIPService(
	playerRepo ports.PlayerRepository,
	configRepo ports.VIPConfigRepository,
	playerVIPRepo ports.PlayerVIPRepository,
) *VIPService {
	return &VIPService{
		playerRepo:    playerRepo,
		configRepo:    configRepo,
		playerVIPRepo: playerVIPRepo,
	}
}

func (s *VIPService) Open(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	cfg, err := s.configRepo.GetEnabledConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return view.GameResponse{
				Text:    "VIP indisponivel no momento.",
				MenuKey: "EXCHANGE_MENU",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
			}, nil
		}
		return view.GameResponse{}, err
	}

	statusText := "Status VIP: inativo"
	now := time.Now().UTC()
	vip, vipErr := s.playerVIPRepo.GetByPlayerID(int64(player.ID))
	if vipErr == nil && vip != nil && vip.ExpiresAt.After(now) {
		statusText = fmt.Sprintf("Status VIP: ativo (%s restante)", formatVIPRemaining(vip.ExpiresAt.Sub(now)))
	}

	msg := fmt.Sprintf(
		"<b>DragonKi VIP</b>\n\n"+
			"Custo: %d Ki Crystals\n"+
			"Duracao: %d dias\n\n"+
			"Beneficios:\n"+
			"- Bonus de XP: +%s\n"+
			"- Bonus de Gold: +%s\n"+
			"- Bonus na venda de itens: +%s\n"+
			"- Desconto na loja: -%s\n"+
			"- Limite diario de ataques no Big Boss: %d\n\n"+
			"%s",
		cfg.VIPPriceKiCrystals,
		cfg.DurationDays,
		formatVIPPercent(cfg.XPBonusPercent),
		formatVIPPercent(cfg.GoldBonusPercent),
		formatVIPPercent(cfg.SellBonusPercent),
		formatVIPPercent(cfg.ShopDiscountPercent),
		cfg.BigBossDailyAttackLimit,
		statusText,
	)

	return view.GameResponse{
		Text:    msg,
		MenuKey: "VIP_MENU",
		Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga.png"},
	}, nil
}

func (s *VIPService) Purchase(userID int64) (view.GameResponse, error) {
	player, err := s.playerRepo.GetOrCreate(userID)
	if err != nil {
		return view.GameResponse{}, err
	}
	cfg, err := s.configRepo.GetEnabledConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return view.GameResponse{
				Text:    "VIP indisponivel no momento.",
				MenuKey: "EXCHANGE_MENU",
				Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
			}, nil
		}
		return view.GameResponse{}, err
	}

	if player.KiCrystals < cfg.VIPPriceKiCrystals {
		return view.GameResponse{
			Text:    "Voce nao possui Ki Crystals suficientes para comprar o VIP.",
			MenuKey: "VIP_MENU",
			Image:   &view.Image{Type: view.ImageAsset, Path: "assets/error.png"},
		}, nil
	}

	now := time.Now().UTC()
	duration := time.Duration(cfg.DurationDays) * 24 * time.Hour

	var vip *domain.PlayerVIP
	currentVIP, vipErr := s.playerVIPRepo.GetByPlayerID(int64(player.ID))
	if vipErr != nil && !errors.Is(vipErr, gorm.ErrRecordNotFound) {
		return view.GameResponse{}, vipErr
	}

	if currentVIP == nil || currentVIP.ExpiresAt.IsZero() {
		vip = &domain.PlayerVIP{
			PlayerID:  int64(player.ID),
			StartedAt: now,
			ExpiresAt: now.Add(duration),
			CreatedAt: now,
			UpdatedAt: now,
		}
	} else {
		vip = currentVIP
		if vip.ExpiresAt.After(now) {
			vip.ExpiresAt = vip.ExpiresAt.Add(duration)
		} else {
			vip.StartedAt = now
			vip.ExpiresAt = now.Add(duration)
		}
		vip.UpdatedAt = now
	}

	player.KiCrystals -= cfg.VIPPriceKiCrystals
	if err := s.playerRepo.Update(player); err != nil {
		return view.GameResponse{}, err
	}
	if err := s.playerVIPRepo.Upsert(vip); err != nil {
		return view.GameResponse{}, err
	}

	return view.GameResponse{
		Text: fmt.Sprintf(
			"VIP ativado com sucesso!\n\nValido ate: %s",
			vip.ExpiresAt.Format("02/01/2006 15:04"),
		),
		MenuKey: "VIP_MENU",
		Image:   &view.Image{Type: view.ImageAsset, Path: "assets/npcs/mestre_kalunga.png"},
	}, nil
}

func (s *VIPService) GetVIPStatusText(playerID int64) string {
	now := time.Now().UTC()
	vip := s.activeVIP(playerID, now)
	if vip == nil {
		return "VIP: inativo"
	}
	return fmt.Sprintf("VIP: ativo (%s restante)", formatVIPRemaining(vip.ExpiresAt.Sub(now)))
}

func (s *VIPService) ApplyBattleBonuses(playerID int64, xpBase int, goldBase int) (int, int) {
	if s == nil || s.configRepo == nil {
		return xpBase, goldBase
	}
	if s.activeVIP(playerID, time.Now().UTC()) == nil {
		return xpBase, goldBase
	}
	cfg, err := s.configRepo.GetEnabledConfig()
	if err != nil || cfg == nil {
		return xpBase, goldBase
	}

	xpFinal := xpBase
	goldFinal := goldBase
	if cfg.XPBonusPercent > 0 {
		xpFinal += int(math.Floor(float64(xpBase) * cfg.XPBonusPercent))
	}
	if cfg.GoldBonusPercent > 0 {
		goldFinal += int(math.Floor(float64(goldBase) * cfg.GoldBonusPercent))
	}
	return xpFinal, goldFinal
}

func (s *VIPService) GetShopPriceAndSellMultiplier(playerID int64) (priceMultiplier float64, sellMultiplier float64) {
	priceMultiplier = 1.0
	sellMultiplier = 1.0

	if s == nil || s.configRepo == nil {
		return priceMultiplier, sellMultiplier
	}
	if s.activeVIP(playerID, time.Now().UTC()) == nil {
		return priceMultiplier, sellMultiplier
	}
	cfg, err := s.configRepo.GetEnabledConfig()
	if err != nil || cfg == nil {
		return priceMultiplier, sellMultiplier
	}
	if cfg.ShopDiscountPercent > 0 {
		priceMultiplier = 1.0 - cfg.ShopDiscountPercent
		if priceMultiplier < 0.1 {
			priceMultiplier = 0.1
		}
	}
	if cfg.SellBonusPercent > 0 {
		sellMultiplier = 1.0 + cfg.SellBonusPercent
	}
	return priceMultiplier, sellMultiplier
}

func (s *VIPService) GetBigBossDailyAttackLimit(playerID int64) int {
	baseLimit := 1
	if s == nil || s.configRepo == nil {
		return baseLimit
	}
	if s.activeVIP(playerID, time.Now().UTC()) == nil {
		return baseLimit
	}
	cfg, err := s.configRepo.GetEnabledConfig()
	if err != nil || cfg == nil {
		return baseLimit
	}
	if cfg.BigBossDailyAttackLimit < 1 {
		return baseLimit
	}
	return cfg.BigBossDailyAttackLimit
}

func (s *VIPService) HasActiveVIP(playerID int64) bool {
	return s.activeVIP(playerID, time.Now().UTC()) != nil
}

func (s *VIPService) DecoratePlayerName(playerID int64, name string) string {
	if !s.HasActiveVIP(playerID) {
		return name
	}

	trimmed := strings.TrimSpace(name)
	if trimmed == "" || strings.HasPrefix(trimmed, vipNamePrefix) {
		return name
	}

	return vipNamePrefix + name
}

func (s *VIPService) activeVIP(playerID int64, now time.Time) *domain.PlayerVIP {
	if s == nil || s.playerVIPRepo == nil {
		return nil
	}

	vip, err := s.playerVIPRepo.GetByPlayerID(playerID)
	if err != nil || vip == nil {
		return nil
	}
	if !vip.ExpiresAt.After(now) {
		return nil
	}
	return vip
}

func formatVIPRemaining(d time.Duration) string {
	if d <= 0 {
		return "0h"
	}
	totalHours := int(d.Hours())
	days := totalHours / 24
	hours := totalHours % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	return fmt.Sprintf("%dh", hours)
}

func formatVIPPercent(value float64) string {
	return fmt.Sprintf("%.0f%%", value*100)
}
