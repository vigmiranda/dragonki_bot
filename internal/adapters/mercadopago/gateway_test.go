package mercadopago

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/vigmiranda/dragonki/internal/domain"
)

type fakePaymentClient struct {
	createFn func(ctx context.Context, request payment.Request) (*payment.Response, error)

	lastCtx     context.Context
	lastRequest payment.Request
}

func (f *fakePaymentClient) Create(ctx context.Context, request payment.Request) (*payment.Response, error) {
	f.lastCtx = ctx
	f.lastRequest = request
	if f.createFn == nil {
		return nil, nil
	}
	return f.createFn(ctx, request)
}

func (f *fakePaymentClient) Search(ctx context.Context, request payment.SearchRequest) (*payment.SearchResponse, error) {
	return nil, nil
}

func (f *fakePaymentClient) Get(ctx context.Context, id int) (*payment.Response, error) {
	return nil, nil
}

func (f *fakePaymentClient) Cancel(ctx context.Context, id int) (*payment.Response, error) {
	return nil, nil
}

func (f *fakePaymentClient) Capture(ctx context.Context, id int) (*payment.Response, error) {
	return nil, nil
}

func (f *fakePaymentClient) CaptureAmount(ctx context.Context, id int, amount float64) (*payment.Response, error) {
	return nil, nil
}

func TestNewGateway(t *testing.T) {
	g := NewGateway(&Client{})
	if g == nil {
		t.Fatalf("expected non-nil gateway")
	}
}

func TestNewClient_WithToken(t *testing.T) {
	t.Setenv("MP_ACCESS_TOKEN", "test-token")
	c := NewClient()
	if c == nil {
		t.Fatalf("expected non-nil client")
	}
	if c.Ctx == nil {
		t.Fatalf("expected context initialized")
	}
}

func TestGateway_CreatePix_CreateError(t *testing.T) {
	fakeClient := &fakePaymentClient{
		createFn: func(ctx context.Context, request payment.Request) (*payment.Response, error) {
			return nil, errors.New("create fail")
		},
	}
	g := &Gateway{
		client: &Client{Payment: fakeClient, Ctx: context.Background()},
	}

	_, _, err := g.CreatePix(context.Background(), 10, 5)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGateway_CreatePix_MissingQRCodeBase64(t *testing.T) {
	fakeClient := &fakePaymentClient{
		createFn: func(ctx context.Context, request payment.Request) (*payment.Response, error) {
			return &payment.Response{
				ID: 1,
				PointOfInteraction: payment.PointOfInteractionResponse{
					TransactionData: payment.TransactionDataResponse{
						QRCodeBase64: "",
						QRCode:       "code",
					},
				},
			}, nil
		},
	}
	g := &Gateway{
		client: &Client{Payment: fakeClient, Ctx: context.Background()},
	}

	_, _, err := g.CreatePix(context.Background(), 10, 5)
	if err == nil || !strings.Contains(err.Error(), "qrBase64") {
		t.Fatalf("expected qrBase64 error, got %v", err)
	}
}

func TestGateway_CreatePix_MissingQRCode(t *testing.T) {
	fakeClient := &fakePaymentClient{
		createFn: func(ctx context.Context, request payment.Request) (*payment.Response, error) {
			return &payment.Response{
				ID: 1,
				PointOfInteraction: payment.PointOfInteractionResponse{
					TransactionData: payment.TransactionDataResponse{
						QRCodeBase64: "YWJj",
						QRCode:       "",
					},
				},
			}, nil
		},
	}
	g := &Gateway{
		client: &Client{Payment: fakeClient, Ctx: context.Background()},
	}

	_, _, err := g.CreatePix(context.Background(), 10, 5)
	if err == nil || !strings.Contains(err.Error(), "qrCode") {
		t.Fatalf("expected qrCode error, got %v", err)
	}
}

func TestGateway_CreatePix_ZeroID(t *testing.T) {
	fakeClient := &fakePaymentClient{
		createFn: func(ctx context.Context, request payment.Request) (*payment.Response, error) {
			return &payment.Response{
				ID: 0,
				PointOfInteraction: payment.PointOfInteractionResponse{
					TransactionData: payment.TransactionDataResponse{
						QRCodeBase64: "YWJj",
						QRCode:       "pix",
					},
				},
			}, nil
		},
	}
	g := &Gateway{
		client: &Client{Payment: fakeClient, Ctx: context.Background()},
	}

	_, _, err := g.CreatePix(context.Background(), 10, 5)
	if err == nil || !strings.Contains(err.Error(), "ID - PIX") {
		t.Fatalf("expected ID error, got %v", err)
	}
}

func TestGateway_CreatePix_Success(t *testing.T) {
	t.Setenv("MP_NOTIFICATION_URL", "https://example.com/hook")

	fakeClient := &fakePaymentClient{
		createFn: func(ctx context.Context, request payment.Request) (*payment.Response, error) {
			return &payment.Response{
				ID: 123,
				PointOfInteraction: payment.PointOfInteractionResponse{
					TransactionData: payment.TransactionDataResponse{
						QRCodeBase64: "YWJj",
						QRCode:       "pix-copy-paste",
					},
				},
			}, nil
		},
	}
	g := &Gateway{
		client: &Client{Payment: fakeClient, Ctx: context.Background()},
	}

	pay, qr, err := g.CreatePix(context.Background(), 99, 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if qr != "YWJj" {
		t.Fatalf("expected qr base64 YWJj, got %q", qr)
	}
	if pay == nil {
		t.Fatalf("expected payment")
	}
	if pay.Status != domain.PaymentPending {
		t.Fatalf("expected pending status, got %s", pay.Status)
	}
	if pay.PlayerID != 99 || pay.Amount != 7 || pay.PartnerId != 123 || pay.QRCode != "pix-copy-paste" {
		t.Fatalf("unexpected payment payload: %+v", pay)
	}

	if fakeClient.lastRequest.PaymentMethodID != "pix" {
		t.Fatalf("expected pix payment method")
	}
	if fakeClient.lastRequest.TransactionAmount != 7 {
		t.Fatalf("expected amount 7, got %v", fakeClient.lastRequest.TransactionAmount)
	}
	if fakeClient.lastRequest.NotificationURL != "https://example.com/hook" {
		t.Fatalf("unexpected notification url %q", fakeClient.lastRequest.NotificationURL)
	}
	if fakeClient.lastRequest.Payer == nil || fakeClient.lastRequest.Payer.Email != "99_@dragonki.com" {
		t.Fatalf("unexpected payer: %+v", fakeClient.lastRequest.Payer)
	}

	deviceID, _ := fakeClient.lastCtx.Value("X-Device-Id").(string)
	if !strings.HasPrefix(deviceID, "telegram_99_") {
		t.Fatalf("expected device id prefix telegram_99_, got %q", deviceID)
	}
}
