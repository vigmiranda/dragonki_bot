package telegram

import (
	"fmt"
	"os"
	"sync"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type telegramAPIMock struct {
	mu       sync.Mutex
	calls    []string
	nextID   int
	response *tgbotapi.APIResponse
}

func (m *telegramAPIMock) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextID++
	m.calls = append(m.calls, fmt.Sprintf("%T", c))
	return tgbotapi.Message{MessageID: m.nextID}, nil
}

func (m *telegramAPIMock) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, fmt.Sprintf("%T", c))
	if m.response != nil {
		return m.response, nil
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func TestQueuedTelegramDelivery_SerializesSendAndRequest(t *testing.T) {
	api := &telegramAPIMock{}
	delivery := newTelegramDelivery(api)
	defer delivery.Close()

	msg, err := delivery.Send(tgbotapi.NewMessage(1, "hello"))
	if err != nil {
		t.Fatalf("send message: %v", err)
	}
	if msg.MessageID != 1 {
		t.Fatalf("unexpected message id: got %d want 1", msg.MessageID)
	}

	resp, err := delivery.Request(tgbotapi.NewCallback("cb-id", ""))
	if err != nil {
		t.Fatalf("request callback: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("unexpected callback response: %+v", resp)
	}

	api.mu.Lock()
	defer api.mu.Unlock()
	if len(api.calls) != 2 {
		t.Fatalf("expected 2 delivery calls, got %d", len(api.calls))
	}
	if api.calls[0] != "tgbotapi.MessageConfig" || api.calls[1] != "tgbotapi.CallbackConfig" {
		t.Fatalf("unexpected delivery order: %v", api.calls)
	}
}

func TestNewTelegramDelivery_NoopMode(t *testing.T) {
	t.Setenv("TELEGRAM_DELIVERY_MODE", "noop")
	defer os.Unsetenv("TELEGRAM_DELIVERY_MODE")

	delivery := newTelegramDelivery(&telegramAPIMock{})
	defer delivery.Close()

	msg, err := delivery.Send(tgbotapi.NewMessage(1, "hello"))
	if err != nil {
		t.Fatalf("noop send: %v", err)
	}
	if msg.MessageID != 1 {
		t.Fatalf("expected noop message id 1, got %d", msg.MessageID)
	}

	resp, err := delivery.Request(tgbotapi.NewCallback("cb-id", ""))
	if err != nil {
		t.Fatalf("noop request: %v", err)
	}
	if resp == nil || !resp.Ok {
		t.Fatalf("unexpected noop response: %+v", resp)
	}
}
