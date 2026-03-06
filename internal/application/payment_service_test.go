package application

import (
	"errors"
	"testing"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakePaymentGateway struct {
	createPixFn   func(playerID int64, amount float64) (*domain.Payment, string, error)
	lastPlayerID  int64
	lastAmount    float64
	createPixCall int
}

func (f *fakePaymentGateway) CreatePix(playerID int64, amount float64) (*domain.Payment, string, error) {
	f.createPixCall++
	f.lastPlayerID = playerID
	f.lastAmount = amount
	if f.createPixFn == nil {
		return nil, "", nil
	}
	return f.createPixFn(playerID, amount)
}

type fakePaymentRepo struct {
	saveFn            func(payment *domain.Payment) error
	findByPartnerIDFn func(id int64) (*domain.Payment, error)
	updateStatusFn    func(id string, status domain.PaymentStatus) error

	savedPayments      []*domain.Payment
	findByPartnerCalls int
	updateStatusCalls  int
	lastUpdateID       string
	lastUpdateStatus   domain.PaymentStatus
}

func (f *fakePaymentRepo) Save(payment *domain.Payment) error {
	f.savedPayments = append(f.savedPayments, payment)
	if f.saveFn == nil {
		return nil
	}
	return f.saveFn(payment)
}

func (f *fakePaymentRepo) FindByPartnerID(id int64) (*domain.Payment, error) {
	f.findByPartnerCalls++
	if f.findByPartnerIDFn == nil {
		return nil, nil
	}
	return f.findByPartnerIDFn(id)
}

func (f *fakePaymentRepo) UpdateStatus(id string, status domain.PaymentStatus) error {
	f.updateStatusCalls++
	f.lastUpdateID = id
	f.lastUpdateStatus = status
	if f.updateStatusFn == nil {
		return nil
	}
	return f.updateStatusFn(id, status)
}

func TestNewPaymentService(t *testing.T) {
	gateway := &fakePaymentGateway{}
	paymentRepo := &fakePaymentRepo{}
	playerRepo := &fakePlayerRepo{}

	svc := NewPaymentService(gateway, paymentRepo, playerRepo)
	if svc == nil {
		t.Fatalf("expected non-nil service")
	}
}

func TestPaymentService_CreatePix_GatewayError(t *testing.T) {
	svc := &PaymentService{
		gateway: &fakePaymentGateway{
			createPixFn: func(playerID int64, amount float64) (*domain.Payment, string, error) {
				return nil, "", errors.New("gateway fail")
			},
		},
		paymentRepository: &fakePaymentRepo{},
	}

	_, _, err := svc.CreatePix(1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_CreatePix_SaveError(t *testing.T) {
	payment := &domain.Payment{ID: "p1", PlayerID: 1, Amount: 10, Status: domain.PaymentPending}
	svc := &PaymentService{
		gateway: &fakePaymentGateway{
			createPixFn: func(playerID int64, amount float64) (*domain.Payment, string, error) {
				return payment, "YWJj", nil
			},
		},
		paymentRepository: &fakePaymentRepo{
			saveFn: func(payment *domain.Payment) error { return errors.New("save fail") },
		},
	}

	_, _, err := svc.CreatePix(1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_CreatePix_Success(t *testing.T) {
	payment := &domain.Payment{ID: "p1", PlayerID: 1, Amount: 10, Status: domain.PaymentPending}
	repo := &fakePaymentRepo{}
	svc := &PaymentService{
		gateway: &fakePaymentGateway{
			createPixFn: func(playerID int64, amount float64) (*domain.Payment, string, error) {
				return payment, "YWJj", nil
			},
		},
		paymentRepository: repo,
	}

	got, qr, err := svc.CreatePix(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != payment {
		t.Fatalf("expected same payment pointer")
	}
	if qr != "YWJj" {
		t.Fatalf("expected qr YWJj, got %q", qr)
	}
	if len(repo.savedPayments) != 1 {
		t.Fatalf("expected one saved payment, got %d", len(repo.savedPayments))
	}
}

func TestPaymentService_HandleBuyKi_DefaultAmount(t *testing.T) {
	t.Setenv("KI_CRYSTALS_VALUE", "")

	gateway := &fakePaymentGateway{
		createPixFn: func(playerID int64, amount float64) (*domain.Payment, string, error) {
			return &domain.Payment{
				ID:     "p1",
				Amount: amount,
				QRCode: "pix-code",
			}, "YWJj", nil
		},
	}
	repo := &fakePaymentRepo{}
	svc := &PaymentService{
		gateway:           gateway,
		paymentRepository: repo,
		playersRepository: &fakePlayerRepo{},
	}

	resp, err := svc.HandleBuyKi(77)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gateway.lastAmount != 1 {
		t.Fatalf("expected default amount 1, got %v", gateway.lastAmount)
	}
	if resp.MenuKey != "EXCHANGE_MENU" {
		t.Fatalf("expected EXCHANGE_MENU, got %q", resp.MenuKey)
	}
	if resp.Image == nil ||
		resp.Image.Name != "pix.png" ||
		resp.Image.Type != view.ImageBase64 ||
		resp.Image.Base64 == "" {
		t.Fatalf("unexpected image payload: %+v", resp.Image)
	}
}

func TestPaymentService_HandleBuyKi_InvalidEnv(t *testing.T) {
	t.Setenv("KI_CRYSTALS_VALUE", "abc")
	svc := &PaymentService{
		gateway:           &fakePaymentGateway{},
		paymentRepository: &fakePaymentRepo{},
		playersRepository: &fakePlayerRepo{},
	}

	_, err := svc.HandleBuyKi(1)
	if err == nil {
		t.Fatalf("expected env parse error")
	}
}

func TestPaymentService_HandleBuyKi_ValidEnvAmount(t *testing.T) {
	t.Setenv("KI_CRYSTALS_VALUE", "4")
	gateway := &fakePaymentGateway{
		createPixFn: func(playerID int64, amount float64) (*domain.Payment, string, error) {
			return &domain.Payment{
				ID:     "p2",
				Amount: amount,
				QRCode: "pix-code",
			}, "YWJj", nil
		},
	}
	svc := &PaymentService{
		gateway:           gateway,
		paymentRepository: &fakePaymentRepo{},
		playersRepository: &fakePlayerRepo{},
	}

	_, err := svc.HandleBuyKi(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gateway.lastAmount != 4 {
		t.Fatalf("expected amount 4, got %v", gateway.lastAmount)
	}
}

func TestPaymentService_ConfirmPayment_InvalidMultEnv(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "abc")
	svc := &PaymentService{
		paymentRepository: &fakePaymentRepo{},
		playersRepository: &fakePlayerRepo{},
	}
	err := svc.ConfirmPayment(1)
	if err == nil {
		t.Fatalf("expected env parse error")
	}
}

func TestPaymentService_ConfirmPayment_FindByPartnerError(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "2")
	svc := &PaymentService{
		paymentRepository: &fakePaymentRepo{
			findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
				return nil, errors.New("not found")
			},
		},
		playersRepository: &fakePlayerRepo{},
	}
	err := svc.ConfirmPayment(9)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_ConfirmPayment_AlreadyApproved(t *testing.T) {
	repo := &fakePaymentRepo{
		findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
			return &domain.Payment{
				ID:       "px",
				PlayerID: 1,
				Amount:   10,
				Status:   domain.PaymentApproved,
			}, nil
		},
	}
	playerRepo := &fakePlayerRepo{}
	svc := &PaymentService{
		paymentRepository: repo,
		playersRepository: playerRepo,
	}

	err := svc.ConfirmPayment(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updateStatusCalls != 0 {
		t.Fatalf("expected no update status call, got %d", repo.updateStatusCalls)
	}
	if playerRepo.getOrCreateCallCount != 0 {
		t.Fatalf("expected no player fetch for approved payment")
	}
}

func TestPaymentService_ConfirmPayment_UpdateStatusError(t *testing.T) {
	repo := &fakePaymentRepo{
		findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
			return &domain.Payment{
				ID:       "px",
				PlayerID: 1,
				Amount:   10,
				Status:   domain.PaymentPending,
			}, nil
		},
		updateStatusFn: func(id string, status domain.PaymentStatus) error {
			return errors.New("update fail")
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
		playersRepository: &fakePlayerRepo{},
	}

	err := svc.ConfirmPayment(5)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_ConfirmPayment_GetPlayerError(t *testing.T) {
	repo := &fakePaymentRepo{
		findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
			return &domain.Payment{
				ID:       "px",
				PlayerID: 1,
				Amount:   10,
				Status:   domain.PaymentPending,
			}, nil
		},
	}
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) {
			return nil, errors.New("player fail")
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
		playersRepository: playerRepo,
	}

	err := svc.ConfirmPayment(5)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_ConfirmPayment_UpdatePlayerError(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "3")
	repo := &fakePaymentRepo{
		findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
			return &domain.Payment{
				ID:       "px",
				PlayerID: 1,
				Amount:   2.9,
				Status:   domain.PaymentPending,
			}, nil
		},
	}
	player := basePlayer()
	player.KiCrystals = 0
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
		updateFn: func(player *domain.Player) error {
			return errors.New("save fail")
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
		playersRepository: playerRepo,
	}

	err := svc.ConfirmPayment(5)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_ConfirmPayment_Success(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "3")
	repo := &fakePaymentRepo{
		findByPartnerIDFn: func(id int64) (*domain.Payment, error) {
			return &domain.Payment{
				ID:       "pay-1",
				PlayerID: 42,
				Amount:   2.9,
				Status:   domain.PaymentPending,
			}, nil
		},
	}
	player := basePlayer()
	player.UserID = 42
	player.KiCrystals = 10
	playerRepo := &fakePlayerRepo{
		getOrCreateFn: func(userID int64) (*domain.Player, error) { return player, nil },
	}
	svc := &PaymentService{
		paymentRepository: repo,
		playersRepository: playerRepo,
	}

	err := svc.ConfirmPayment(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updateStatusCalls != 1 {
		t.Fatalf("expected one update status call, got %d", repo.updateStatusCalls)
	}
	if repo.lastUpdateID != "pay-1" || repo.lastUpdateStatus != domain.PaymentApproved {
		t.Fatalf("unexpected status update payload: id=%s status=%s", repo.lastUpdateID, repo.lastUpdateStatus)
	}
	// int(2.9) * 3 = 6
	if player.KiCrystals != 16 {
		t.Fatalf("expected ki crystals 16, got %d", player.KiCrystals)
	}
	if len(playerRepo.updatedPlayers) != 1 {
		t.Fatalf("expected one player update call, got %d", len(playerRepo.updatedPlayers))
	}
}
