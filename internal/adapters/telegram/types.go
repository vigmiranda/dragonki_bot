package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type telegramCallbackReceiptStore interface {
	Reserve(ctx context.Context, receipt *domain.TelegramCallbackReceipt) (bool, error)
}

type telegramCallbackDedupeStore interface {
	Reserve(ctx context.Context, dedupe *domain.TelegramCallbackDedupe) (bool, error)
}

type Handler struct {
	bot               *tgbotapi.BotAPI
	delivery          telegramDelivery
	presenter         *telegramPresenter
	playerService     *application.PlayerService
	battleService     *application.BattleService
	paymentService    *application.PaymentService
	onboardingService *application.OnboardingService
	transformService  *application.TransformService
	menuService       *application.MenuService
	mapService        *application.MapService
	energyService     *application.EnergyService
	shopService       *application.ShopService
	exchangeService   *application.ExchangeService
	inventoryService  *application.InventoryService
	rankService       *application.RankService
	weeklyBossService *application.WeeklyBossService
	referralService   *application.ReferralService
	chamberService    *application.AncestralKiChamberService
	vipService        *application.VIPService
	callbackReceipts  telegramCallbackReceiptStore
	callbackDedupes   telegramCallbackDedupeStore
	callbackGuard     *telegramCallbackGuard
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	playerService *application.PlayerService,
	battleService *application.BattleService,
	paymentService *application.PaymentService,
	onboardingService *application.OnboardingService,
	transformService *application.TransformService,
	menuService *application.MenuService,
	mapService *application.MapService,
	energyService *application.EnergyService,
	shopService *application.ShopService,
	exchangeService *application.ExchangeService,
	inventoryService *application.InventoryService,
	rankService *application.RankService,
	weeklyBossService *application.WeeklyBossService,
	referralService *application.ReferralService,
	chamberService *application.AncestralKiChamberService,
	callbackReceipts telegramCallbackReceiptStore,
	callbackDedupes telegramCallbackDedupeStore,
	vipService ...*application.VIPService,
) *Handler {
	var v *application.VIPService
	if len(vipService) > 0 {
		v = vipService[0]
	}
	delivery := newTelegramDelivery(bot)
	return &Handler{
		bot:               bot,
		delivery:          delivery,
		presenter:         newTelegramPresenter(delivery),
		playerService:     playerService,
		battleService:     battleService,
		paymentService:    paymentService,
		onboardingService: onboardingService,
		transformService:  transformService,
		menuService:       menuService,
		mapService:        mapService,
		energyService:     energyService,
		shopService:       shopService,
		exchangeService:   exchangeService,
		inventoryService:  inventoryService,
		rankService:       rankService,
		weeklyBossService: weeklyBossService,
		referralService:   referralService,
		chamberService:    chamberService,
		vipService:        v,
		callbackReceipts:  callbackReceipts,
		callbackDedupes:   callbackDedupes,
		callbackGuard:     newTelegramCallbackGuard(telegramCallbackGuardTTL),
	}
}

func (h *Handler) Close() {
	if h == nil || h.delivery == nil {
		return
	}
	h.delivery.Close()
}
