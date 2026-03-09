package main

import (
	"encoding/json"
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
)

type telegramUpdateHandler interface {
	HandleUpdate(update tgbotapi.Update)
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

func startWeeklyBossScheduler(service *application.WeeklyBossService) {
	if service == nil {
		return
	}

	const interval = time.Hour
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

		for range ticker.C {
			run()
		}
	}()
}

func main() {
	appEnv, loadedEnvFiles, err := bootstrap.LoadEnvByEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("env loaded | APP_ENV=%s | files=%v", appEnv, loadedEnvFiles)

	db := postgres.OpenDB()

	if shouldRunAutoMigrate() {
		err = db.AutoMigrate(
			&domain.Player{},
			&domain.Class{},
			&domain.Monster{},
			&domain.Region{},
			&domain.Payment{},
			&domain.PaymentWebhookAudit{},
			&domain.Item{},
			&domain.Equipment{},
			&domain.Inventory{},
			&domain.InventoryItem{},
			&domain.PlayerEquipment{},
			&domain.BattleLog{},
			&domain.MonsterDrop{},
			&domain.FeatureFlag{},
			&domain.AncestralKiChamberProgress{},
			&domain.AncestralKiChamberRun{},
			&domain.AncestralKiChamberBattleLog{},
		)
		if err != nil {
			log.Fatal("error ao rodar automigrate", err)
		}
	}

	if shouldRunMigrations() {
		if err := bootstrap.RunMigrationsUp(os.Getenv("DATABASE_URL"), os.Getenv("MIGRATIONS_DIR")); err != nil {
			log.Fatal("error ao rodar migrations", err)
		}
	}

	playerRepo := postgres.NewPlayerPostgresRepo(db)
	monsterRepo := postgres.NewMonsterPostgresRepo(db)
	regionRepo := postgres.NewRegionPostgresRepo(db)
	paymentRepo := postgres.NewPaymentRepo(db)
	featureFlagRepo := postgres.NewFeatureFlagPostgresRepo(db)
	referralRepo := postgres.NewReferralPostgresRepo(db)
	itemRepo := postgres.NewItemsPostgresRepo(db)
	inventoryRepo := postgres.NewInventoryPostgresRepo(db)
	battleRepo := postgres.NewBattlePostgresRepo(db)
	weeklyBossRepo := postgres.NewWeeklyBossPostgresRepo(db)
	chamberRepo := postgres.NewAncestralKiChamberPostgresRepo(db)
	bootstrap.SetFeatureFlagResolver(featureFlagRepo.FindByName)

	playerService := application.NewPlayerService(playerRepo, regionRepo, inventoryRepo)
	onboardingService := application.NewOnboardingService(playerRepo)
	transformService := application.NewTransformService(playerRepo)
	menuService := application.NewMenuService(playerRepo, regionRepo)
	mapService := application.NewMapService(playerRepo, regionRepo)
	energyService := application.NewEnergyService(playerRepo)
	shopService := application.NewShopService(playerRepo, itemRepo, inventoryRepo)
	exchangeService := application.NewExchangeService(playerRepo)
	inventoryService := application.NewInventoryService(inventoryRepo, playerRepo, itemRepo)
	rankService := application.NewRankService(battleRepo)
	chamberService := application.NewAncestralKiChamberService(chamberRepo, playerRepo, monsterRepo, inventoryRepo)
	var rewardQueue application.WeeklyBossRewardQueue
	if isBigBossQueueEnabled() {
		queueName := strings.TrimSpace(os.Getenv("BIG_BOSS_QUEUE_NAME"))
		q, qErr := queue.NewAsynqBigBossQueue(os.Getenv("REDIS_URL"), queueName)
		if qErr != nil {
			log.Printf("big boss queue init failed, fallback inline payout | err=%v", qErr)
		} else {
			rewardQueue = q
			defer q.Close()
		}
	}

	weeklyBossService := application.NewWeeklyBossService(weeklyBossRepo, playerRepo, inventoryRepo, rewardQueue)
	if bootstrap.IsFeatureEnabled("big_boss_weekly", false) {
		startWeeklyBossScheduler(weeklyBossService)
	}

	mpClient := mercadopago.NewClient()
	mpGateway := mercadopago.NewGateway(mpClient)
	paymentService := application.NewPaymentService(
		mpGateway,
		paymentRepo,
		playerRepo,
	)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DragonKi bot iniciado:", bot.Self.UserName)

	referralService := application.NewReferralService(referralRepo, playerRepo, bot.Self.UserName)
	battleService := application.NewBattleService(playerRepo, monsterRepo, regionRepo, battleRepo, inventoryRepo, referralService)

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
	)

	mux := http.NewServeMux()
	registerTelegramWebhook(mux, handler, decodeTelegramUpdate)

	log.Println("Telegram Webhook escutando na porta 8080")
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
