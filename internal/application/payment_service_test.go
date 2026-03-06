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
	settlePendingFn   func(partnerID int64, kiCrystalsMult int) (bool, error)
	savedPayments     []*domain.Payment
	settleCalls       int
	lastPartnerID     int64
	lastKiCrystalsMul int
}

func (f *fakePaymentRepo) Save(payment *domain.Payment) error {
	f.savedPayments = append(f.savedPayments, payment)
	if f.saveFn == nil {
		return nil
	}
	return f.saveFn(payment)
}

func (f *fakePaymentRepo) SettlePendingByPartnerID(partnerID int64, kiCrystalsMult int) (bool, error) {
	f.settleCalls++
	f.lastPartnerID = partnerID
	f.lastKiCrystalsMul = kiCrystalsMult
	if f.settlePendingFn == nil {
		return false, nil
	}
	return f.settlePendingFn(partnerID, kiCrystalsMult)
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
	}
	err := svc.ConfirmPayment(1)
	if err == nil {
		t.Fatalf("expected env parse error")
	}
}

func TestPaymentService_ConfirmPayment_SettleError(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "2")
	expectedErr := errors.New("settle fail")
	svc := &PaymentService{
		paymentRepository: &fakePaymentRepo{
			settlePendingFn: func(partnerID int64, kiCrystalsMult int) (bool, error) {
				return false, expectedErr
			},
		},
	}
	err := svc.ConfirmPayment(9)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected settle error, got %v", err)
	}
}

func TestPaymentService_ConfirmPayment_AlreadySettled_IsIdempotent(t *testing.T) {
	repo := &fakePaymentRepo{
		settlePendingFn: func(partnerID int64, kiCrystalsMult int) (bool, error) {
			return false, nil
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
	}

	err := svc.ConfirmPayment(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.settleCalls != 1 {
		t.Fatalf("expected one settle call, got %d", repo.settleCalls)
	}
}

func TestPaymentService_ConfirmPayment_Success(t *testing.T) {
	t.Setenv("KI_CRYSTALS_MULT", "3")
	repo := &fakePaymentRepo{
		settlePendingFn: func(partnerID int64, kiCrystalsMult int) (bool, error) {
			if partnerID != 99 {
				t.Fatalf("expected partnerID=99, got %d", partnerID)
			}
			if kiCrystalsMult != 3 {
				t.Fatalf("expected mult=3, got %d", kiCrystalsMult)
			}
			return true, nil
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
	}

	err := svc.ConfirmPayment(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.settleCalls != 1 {
		t.Fatalf("expected one settle call, got %d", repo.settleCalls)
	}
	if repo.lastPartnerID != 99 || repo.lastKiCrystalsMul != 3 {
		t.Fatalf("unexpected settle payload: partnerID=%d mult=%d", repo.lastPartnerID, repo.lastKiCrystalsMul)
	}
}
