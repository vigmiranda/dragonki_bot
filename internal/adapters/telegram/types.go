package telegram

import (
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
	energyService     *application.EnergyService
	shopService       *application.ShopService
	exchangeService   *application.ExchangeService
}

func NewHandler(
	bot *tgbotapi.BotAPI,
	playerService *application.PlayerService,
	battleService *application.BattleService,
	paymentService *application.PaymentService,
	onboardingService *application.OnboardingService,
	transformService *application.TransformService,
	menuService *application.MenuService,
	energyService *application.EnergyService,
	shopService *application.ShopService,
	exchangeService *application.ExchangeService,

) *Handler {
	return &Handler{
		bot:               bot,
		playerService:     playerService,
		battleService:     battleService,
		paymentService:    paymentService,
		onboardingService: onboardingService,
		transformService:  transformService,
		menuService:       menuService,
		energyService:     energyService,
		shopService:       shopService,
		exchangeService:   exchangeService,
	}
}
