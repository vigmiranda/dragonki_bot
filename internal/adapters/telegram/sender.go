package telegram

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const defaultTelegramSendInterval = 35 * time.Millisecond
const telegramDeliveryModeNoop = "noop"

var errTelegramDeliveryClosed = errors.New("telegram delivery closed")

type telegramAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}

type telegramDelivery interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	Close()
}

type telegramDeliveryMode int

const (
	telegramDeliverySend telegramDeliveryMode = iota + 1
	telegramDeliveryRequest
)

type telegramSendResult struct {
	message tgbotapi.Message
	err     error
}

type telegramRequestResult struct {
	response *tgbotapi.APIResponse
	err      error
}

type telegramDeliveryJob struct {
	mode      telegramDeliveryMode
	chattable tgbotapi.Chattable
	sendCh    chan telegramSendResult
	requestCh chan telegramRequestResult
}

type queuedTelegramDelivery struct {
	api       telegramAPI
	interval  time.Duration
	queue     chan telegramDeliveryJob
	done      chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

func newTelegramDelivery(api telegramAPI) telegramDelivery {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("TELEGRAM_DELIVERY_MODE")), telegramDeliveryModeNoop) {
		return &noopTelegramDelivery{}
	}

	d := &queuedTelegramDelivery{
		api:      api,
		interval: resolveTelegramSendInterval(),
		queue:    make(chan telegramDeliveryJob, 128),
		done:     make(chan struct{}),
	}
	d.wg.Add(1)
	go d.loop()
	return d
}

type noopTelegramDelivery struct{}

func (d *noopTelegramDelivery) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	return tgbotapi.Message{MessageID: 1}, nil
}

func (d *noopTelegramDelivery) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func (d *noopTelegramDelivery) Close() {}

func resolveTelegramSendInterval() time.Duration {
	raw := strings.TrimSpace(os.Getenv("TELEGRAM_SEND_MIN_INTERVAL_MS"))
	if raw == "" {
		return defaultTelegramSendInterval
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return defaultTelegramSendInterval
	}
	return time.Duration(ms) * time.Millisecond
}

func (d *queuedTelegramDelivery) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if d == nil || d.api == nil {
		return tgbotapi.Message{}, errTelegramDeliveryClosed
	}

	job := telegramDeliveryJob{
		mode:      telegramDeliverySend,
		chattable: c,
		sendCh:    make(chan telegramSendResult, 1),
	}
	if err := d.enqueue(job); err != nil {
		return tgbotapi.Message{}, err
	}

	result, ok := <-job.sendCh
	if !ok {
		return tgbotapi.Message{}, errTelegramDeliveryClosed
	}
	return result.message, result.err
}

func (d *queuedTelegramDelivery) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	if d == nil || d.api == nil {
		return nil, errTelegramDeliveryClosed
	}

	job := telegramDeliveryJob{
		mode:      telegramDeliveryRequest,
		chattable: c,
		requestCh: make(chan telegramRequestResult, 1),
	}
	if err := d.enqueue(job); err != nil {
		return nil, err
	}

	result, ok := <-job.requestCh
	if !ok {
		return nil, errTelegramDeliveryClosed
	}
	return result.response, result.err
}

func (d *queuedTelegramDelivery) Close() {
	if d == nil {
		return
	}
	d.closeOnce.Do(func() {
		close(d.done)
		d.wg.Wait()
	})
}

func (d *queuedTelegramDelivery) enqueue(job telegramDeliveryJob) error {
	select {
	case <-d.done:
		return errTelegramDeliveryClosed
	case d.queue <- job:
		return nil
	}
}

func (d *queuedTelegramDelivery) loop() {
	defer d.wg.Done()

	lastSent := time.Time{}
	for {
		select {
		case <-d.done:
			return
		case job := <-d.queue:
			if d.interval > 0 && !lastSent.IsZero() {
				wait := d.interval - time.Since(lastSent)
				if wait > 0 {
					timer := time.NewTimer(wait)
					select {
					case <-d.done:
						timer.Stop()
						return
					case <-timer.C:
					}
				}
			}

			switch job.mode {
			case telegramDeliverySend:
				msg, err := d.api.Send(job.chattable)
				job.sendCh <- telegramSendResult{message: msg, err: err}
				close(job.sendCh)
			case telegramDeliveryRequest:
				resp, err := d.api.Request(job.chattable)
				job.requestCh <- telegramRequestResult{response: resp, err: err}
				close(job.requestCh)
			}

			lastSent = time.Now()
		}
	}
}
