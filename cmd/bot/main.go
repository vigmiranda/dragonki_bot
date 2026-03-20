package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vigmiranda/dragonki/internal/adapters/mercadopago"
	"github.com/vigmiranda/dragonki/internal/adapters/postgres"
	"github.com/vigmiranda/dragonki/internal/adapters/queue"
	"github.com/vigmiranda/dragonki/internal/adapters/telegram"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
	"github.com/vigmiranda/dragonki/internal/domain"
	"gorm.io/gorm"
)

type telegramUpdateHandler interface {
	HandleUpdate(update tgbotapi.Update)
}

type botAppConfig struct {
	DB                  *gorm.DB
	Bot                 *tgbotapi.BotAPI
	MPClient            *mercadopago.Client
	FeatureFlagResolver func(name string) (enabled bool, found bool)
	RewardQueue         application.WeeklyBossRewardQueue
	AutoMigrate         bool
	RunMigrations       bool
	StartScheduler      bool
	BigBossQueueEnabled bool
	DatabaseURL         string
	MigrationsDir       string
	BotToken            string
	RedisURL            string
	QueueName           string
}

type botApp struct {
	Bot     *tgbotapi.BotAPI
	Handler http.Handler
	Cleanup func()
}

type telegramCallbackCleanupRepository interface {
	DeleteOlderThan(ctx context.Context, before time.Time) (int64, error)
}

func decodeTelegramUpdate(r io.Reader, update *tgbotapi.Update) error {
	return json.NewDecoder(r).Decode(update)
}

func registerTelegramWebhook(
	mux *http.ServeMux,
	handler telegramUpdateHandler,
	decodeFn func(r io.Reader, update *tgbotapi.Update) error,
) {
	mux.HandleFunc("/webhook/telegram", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		defer r.Body.Close()

		var update tgbotapi.Update
		if err := decodeFn(r.Body, &update); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		handler.HandleUpdate(update)
		w.WriteHeader(http.StatusOK)
	})
}

func shouldRunMigrations() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("RUN_MIGRATIONS")), "true")
}

func shouldRunAutoMigrate() bool {
	raw := strings.TrimSpace(os.Getenv("RUN_AUTO_MIGRATE"))
	if raw == "" {
		return true
	}
	return strings.EqualFold(raw, "true")
}

func isBigBossQueueEnabled() bool {
	raw := strings.TrimSpace(os.Getenv("BIG_BOSS_QUEUE_ENABLED"))
	if raw == "" {
		return false
	}
	return strings.EqualFold(raw, "true")
}

func parseDurationEnv(name string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil || parsed <= 0 {
		log.Printf("invalid duration for %s=%q, using fallback %s", name, raw, fallback)
		return fallback
	}
	return parsed
}

func parseBoolEnv(name string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	switch strings.ToLower(raw) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		log.Printf("invalid bool for %s=%q, using fallback %t", name, raw, fallback)
		return fallback
	}
}

func startWeeklyBossScheduler(service *application.WeeklyBossService) func() {
	if service == nil {
		return func() {}
	}

	const interval = time.Hour
	stop := make(chan struct{})
	go func() {
		run := func() {
			if err := service.EnsureCurrentWeeklyBoss(time.Now().UTC()); err != nil {
				log.Printf("weekly boss scheduler error: %v", err)
				return
			}
			log.Printf("weekly boss scheduler sync ok")
		}

		run()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				run()
			case <-stop:
				return
			}
		}
	}()

	return func() {
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
}

func startTelegramCallbackCleanupScheduler(
	receiptRepo telegramCallbackCleanupRepository,
	dedupeRepo telegramCallbackCleanupRepository,
	interval time.Duration,
	receiptRetention time.Duration,
	dedupeRetention time.Duration,
) func() {
	if receiptRepo == nil || dedupeRepo == nil {
		return func() {}
	}
	if interval <= 0 {
		interval = time.Hour
	}
	if receiptRetention <= 0 {
		receiptRetention = 7 * 24 * time.Hour
	}
	if dedupeRetention <= 0 {
		dedupeRetention = time.Hour
	}

	stop := make(chan struct{})
	go func() {
		run := func() {
			now := time.Now().UTC()

			receiptDeleted, receiptErr := receiptRepo.DeleteOlderThan(context.Background(), now.Add(-receiptRetention))
			if receiptErr != nil {
				log.Printf("telegram callback receipt cleanup error: %v", receiptErr)
			} else if receiptDeleted > 0 {
				log.Printf("telegram callback receipt cleanup removed %d rows", receiptDeleted)
			}

			dedupeDeleted, dedupeErr := dedupeRepo.DeleteOlderThan(context.Background(), now.Add(-dedupeRetention))
			if dedupeErr != nil {
				log.Printf("telegram callback dedupe cleanup error: %v", dedupeErr)
			} else if dedupeDeleted > 0 {
				log.Printf("telegram callback dedupe cleanup removed %d rows", dedupeDeleted)
			}
		}

		run()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				run()
			case <-stop:
				return
			}
		}
	}()

	return func() {
		select {
		case <-stop:
		default:
			close(stop)
		}
	}
}

func defaultBotAppConfig() botAppConfig {
	return botAppConfig{
		AutoMigrate:         shouldRunAutoMigrate(),
		RunMigrations:       shouldRunMigrations(),
		StartScheduler:      true,
		BigBossQueueEnabled: isBigBossQueueEnabled(),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		MigrationsDir:       os.Getenv("MIGRATIONS_DIR"),
		BotToken:            os.Getenv("BOT_TOKEN"),
		RedisURL:            os.Getenv("REDIS_URL"),
		QueueName:           strings.TrimSpace(os.Getenv("BIG_BOSS_QUEUE_NAME")),
	}
}

func autoMigrateBotModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Player{},
		&domain.Class{},
		&domain.Monster{},
		&domain.Region{},
		&domain.Payment{},
		&domain.PaymentWebhookAudit{},
		&domain.Item{},
		&domain.EquipmentSet{},
		&domain.EquipmentSetBonus{},
		&domain.Equipment{},
		&domain.Inventory{},
		&domain.InventoryItem{},
		&domain.PlayerEquipment{},
		&domain.BattleLog{},
		&domain.MonsterDrop{},
		&domain.FeatureFlag{},
		&domain.TelegramCallbackReceipt{},
		&domain.TelegramCallbackDedupe{},
		&domain.AncestralKiChamberProgress{},
		&domain.AncestralKiChamberRun{},
		&domain.AncestralKiChamberBattleLog{},
	)
}

func buildBotApp(cfg botAppConfig) (*botApp, error) {
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	}

	cleanupFns := make([]func(), 0, 4)
	cleanup := func() {
		for i := len(cleanupFns) - 1; i >= 0; i-- {
			cleanupFns[i]()
		}
	}

	db := cfg.DB
	if db == nil {
		db = postgres.OpenDB()
		sqlDB, err := db.DB()
		if err == nil {
			cleanupFns = append(cleanupFns, func() {
				_ = sqlDB.Close()
			})
		}
	}
	postgres.AttachObservability(db)

	if cfg.AutoMigrate {
		if err := autoMigrateBotModels(db); err != nil {
			cleanup()
			return nil, err
		}
	}

	if cfg.RunMigrations {
		if err := bootstrap.RunMigrationsUp(cfg.DatabaseURL, cfg.MigrationsDir); err != nil {
			cleanup()
			return nil, err
		}
	}

	playerRepo := postgres.NewPlayerPostgresRepo(db)
	monsterRepo := postgres.NewMonsterPostgresRepo(db)
	regionRepo := postgres.NewRegionPostgresRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)
	vipConfigRepo := postgres.NewVIPConfigRepo(db)
	playerVIPRepo := postgres.NewPlayerVIPRepo(db)
	featureFlagRepo := postgres.NewFeatureFlagPostgresRepo(db)
	referralRepo := postgres.NewReferralPostgresRepo(db)
	itemRepo := postgres.NewItemsPostgresRepo(db)
	inventoryRepo := postgres.NewInventoryPostgresRepo(db)
	battleRepo := postgres.NewBattlePostgresRepo(db)
	weeklyBossRepo := postgres.NewWeeklyBossPostgresRepo(db)
	chamberRepo := postgres.NewAncestralKiChamberPostgresRepo(db)
	callbackReceiptRepo := postgres.NewTelegramCallbackReceiptRepo(db)
	callbackDedupeRepo := postgres.NewTelegramCallbackDedupeRepo(db)

	resolver := cfg.FeatureFlagResolver
	if resolver == nil {
		resolver = featureFlagRepo.FindByName
	}
	restoreFeatureFlags := bootstrap.WithFeatureFlagResolver(resolver)
	cleanupFns = append(cleanupFns, restoreFeatureFlags)

	vipService := application.NewVIPService(playerRepo, vipConfigRepo, playerVIPRepo)
	playerService := application.NewPlayerService(playerRepo, regionRepo, inventoryRepo, vipService)
	onboardingService := application.NewOnboardingService(playerRepo)
	transformService := application.NewTransformService(playerRepo)
	menuService := application.NewMenuService(playerRepo, regionRepo, vipService)
	mapService := application.NewMapService(playerRepo, regionRepo)
	energyService := application.NewEnergyService(playerRepo, vipService)
	shopService := application.NewShopService(playerRepo, itemRepo, inventoryRepo)
	shopService.SetVIPService(vipService)
	exchangeService := application.NewExchangeService(playerRepo, vipService)
	inventoryService := application.NewInventoryService(inventoryRepo, playerRepo, itemRepo)
	rankService := application.NewRankService(battleRepo, vipService)
	chamberService := application.NewAncestralKiChamberService(chamberRepo, playerRepo, monsterRepo, inventoryRepo)

	rewardQueue := cfg.RewardQueue
	if rewardQueue == nil && cfg.BigBossQueueEnabled {
		q, err := queue.NewAsynqBigBossQueue(cfg.RedisURL, cfg.QueueName)
		if err != nil {
			log.Printf("big boss queue init failed, fallback inline payout | err=%v", err)
		} else {
			rewardQueue = q
			cleanupFns = append(cleanupFns, func() {
				_ = q.Close()
			})
		}
	}

	weeklyBossService := application.NewWeeklyBossService(weeklyBossRepo, playerRepo, inventoryRepo, rewardQueue)
	weeklyBossService.SetVIPService(vipService)
	if cfg.StartScheduler && bootstrap.IsFeatureEnabled("big_boss_weekly", false) {
		cleanupFns = append(cleanupFns, startWeeklyBossScheduler(weeklyBossService))
	}
	if cfg.StartScheduler {
		cleanupFns = append(cleanupFns, startTelegramCallbackCleanupScheduler(
			callbackReceiptRepo,
			callbackDedupeRepo,
			parseDurationEnv("TELEGRAM_CALLBACK_CLEANUP_INTERVAL", time.Hour),
			parseDurationEnv("TELEGRAM_CALLBACK_RECEIPT_RETENTION", 7*24*time.Hour),
			parseDurationEnv("TELEGRAM_CALLBACK_DEDUPE_RETENTION", time.Hour),
		))
	}

	mpClient := cfg.MPClient
	if mpClient == nil {
		mpClient = mercadopago.NewClient()
	}
	mpGateway := mercadopago.NewGateway(mpClient)
	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	bot := cfg.Bot
	if bot == nil {
		var err error
		bot, err = tgbotapi.NewBotAPI(cfg.BotToken)
		if err != nil {
			cleanup()
			return nil, err
		}
	}
	if bot == nil {
		cleanup()
		return nil, errors.New("telegram bot not configured")
	}

	referralService := application.NewReferralService(referralRepo, playerRepo, bot.Self.UserName)
	battleService := application.NewBattleService(playerRepo, monsterRepo, regionRepo, battleRepo, inventoryRepo, referralService)
	battleService.SetVIPService(vipService)

	handler := telegram.NewHandler(
		bot,
		playerService,
		battleService,
		paymentService,
		onboardingService,
		transformService,
		menuService,
		mapService,
		energyService,
		shopService,
		exchangeService,
		inventoryService,
		rankService,
		weeklyBossService,
		referralService,
		chamberService,
		callbackReceiptRepo,
		callbackDedupeRepo,
		vipService,
	)
	cleanupFns = append(cleanupFns, handler.Close)

	mux := http.NewServeMux()
	registerTelegramWebhook(mux, handler, decodeTelegramUpdate)

	observabilityEnabled := parseBoolEnv("OBSERVABILITY_ENABLED", true)
	if observabilityEnabled {
		observability := bootstrap.NewHTTPObservability("bot")
		observability.RegisterDebugRoutes(mux, parseBoolEnv("PPROF_ENABLED", true))
		return &botApp{
			Bot:     bot,
			Handler: observability.Wrap(mux),
			Cleanup: cleanup,
		}, nil
	}

	return &botApp{
		Bot:     bot,
		Handler: mux,
		Cleanup: cleanup,
	}, nil
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	app, err := buildBotApp(defaultBotAppConfig())
	if err != nil {
		log.Fatal(err)
	}
	defer app.Cleanup()

	log.Println("DragonKi bot iniciado:", app.Bot.Self.UserName)

	log.Println("Telegram Webhook escutando na porta 8080")
	server := &http.Server{
		Addr:              ":8080",
		Handler:           app.Handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
