package application

import (
	"context"
	"errors"
	"testing"

	"github.com/vigmiranda/dragonki/internal/application/view"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakePaymentGateway struct {
	createPixFn   func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error)
	lastPlayerID  int64
	lastAmount    float64
	createPixCall int
}

func (f *fakePaymentGateway) CreatePix(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
	f.createPixCall++
	f.lastPlayerID = playerID
	f.lastAmount = amount
	if f.createPixFn == nil {
		return nil, "", nil
	}
	return f.createPixFn(ctx, playerID, amount)
}

type fakePaymentRepo struct {
	saveFn            func(ctx context.Context, payment *domain.Payment) error
	settlePendingFn   func(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error)
	savedPayments     []*domain.Payment
	settleCalls       int
	lastPartnerID     int64
	lastKiCrystalsMul int
}

func (f *fakePaymentRepo) Save(ctx context.Context, payment *domain.Payment) error {
	f.savedPayments = append(f.savedPayments, payment)
	if f.saveFn == nil {
		return nil
	}
	return f.saveFn(ctx, payment)
}

func (f *fakePaymentRepo) SettlePendingByPartnerID(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error) {
	f.settleCalls++
	f.lastPartnerID = partnerID
	f.lastKiCrystalsMul = kiCrystalsMult
	if f.settlePendingFn == nil {
		return false, nil
	}
	return f.settlePendingFn(ctx, partnerID, kiCrystalsMult)
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
			createPixFn: func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
				return nil, "", errors.New("gateway fail")
			},
		},
		paymentRepository: &fakePaymentRepo{},
	}

	_, _, err := svc.CreatePix(context.Background(), 1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_CreatePix_SaveError(t *testing.T) {
	payment := &domain.Payment{ID: "p1", PlayerID: 1, Amount: 10, Status: domain.PaymentPending}
	svc := &PaymentService{
		gateway: &fakePaymentGateway{
			createPixFn: func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
				return payment, "YWJj", nil
			},
		},
		paymentRepository: &fakePaymentRepo{
			saveFn: func(ctx context.Context, payment *domain.Payment) error { return errors.New("save fail") },
		},
	}

	_, _, err := svc.CreatePix(context.Background(), 1, 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestPaymentService_CreatePix_Success(t *testing.T) {
	payment := &domain.Payment{ID: "p1", PlayerID: 1, Amount: 10, Status: domain.PaymentPending}
	repo := &fakePaymentRepo{}
	svc := &PaymentService{
		gateway: &fakePaymentGateway{
			createPixFn: func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
				return payment, "YWJj", nil
			},
		},
		paymentRepository: repo,
	}

	got, qr, err := svc.CreatePix(context.Background(), 1, 10)
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
	gateway := &fakePaymentGateway{
		createPixFn: func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
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
		cfg:               appConfig{KiCrystalsValue: 1, KiCrystalsMult: 1},
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

func TestPaymentService_HandleBuyKi_ConfigAmount(t *testing.T) {
	gateway := &fakePaymentGateway{
		createPixFn: func(ctx context.Context, playerID int64, amount float64) (*domain.Payment, string, error) {
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
		cfg:               appConfig{KiCrystalsValue: 4, KiCrystalsMult: 1},
	}

	_, err := svc.HandleBuyKi(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gateway.lastAmount != 4 {
		t.Fatalf("expected amount 4, got %v", gateway.lastAmount)
	}
}

func TestPaymentService_ConfirmPayment_SettleError(t *testing.T) {
	expectedErr := errors.New("settle fail")
	svc := &PaymentService{
		paymentRepository: &fakePaymentRepo{
			settlePendingFn: func(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error) {
				return false, expectedErr
			},
		},
		cfg: appConfig{KiCrystalsMult: 2, KiCrystalsValue: 1},
	}
	err := svc.ConfirmPayment(context.Background(), 9)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected settle error, got %v", err)
	}
}

func TestPaymentService_ConfirmPayment_AlreadySettled_IsIdempotent(t *testing.T) {
	repo := &fakePaymentRepo{
		settlePendingFn: func(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error) {
			return false, nil
		},
	}
	svc := &PaymentService{
		paymentRepository: repo,
		cfg:               appConfig{KiCrystalsMult: 1, KiCrystalsValue: 1},
	}

	err := svc.ConfirmPayment(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.settleCalls != 1 {
		t.Fatalf("expected one settle call, got %d", repo.settleCalls)
	}
}

func TestPaymentService_ConfirmPayment_Success(t *testing.T) {
	repo := &fakePaymentRepo{
		settlePendingFn: func(ctx context.Context, partnerID int64, kiCrystalsMult int) (bool, error) {
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
		cfg:               appConfig{KiCrystalsMult: 3, KiCrystalsValue: 1},
	}

	err := svc.ConfirmPayment(context.Background(), 99)
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
