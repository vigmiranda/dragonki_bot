package telegram

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/vigmiranda/dragonki/internal/application"
)

type Handler struct {
	bot               *tgbotapi.BotAPI
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

	pixCodeMsgMu       sync.Mutex
	pixCodeMessageByID map[int64]int
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

) *Handler {
	return &Handler{
		bot:                bot,
		playerService:      playerService,
		battleService:      battleService,
		paymentService:     paymentService,
		onboardingService:  onboardingService,
		transformService:   transformService,
		menuService:        menuService,
		mapService:         mapService,
		energyService:      energyService,
		shopService:        shopService,
		exchangeService:    exchangeService,
		inventoryService:   inventoryService,
		rankService:        rankService,
		weeklyBossService:  weeklyBossService,
		referralService:    referralService,
		chamberService:     chamberService,
		pixCodeMessageByID: make(map[int64]int),
	}
}
