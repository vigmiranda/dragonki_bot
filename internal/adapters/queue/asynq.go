package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hibiken/asynq"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

type AsynqBigBossQueue struct {
	client    *asynq.Client
	queueName string
}

func NewAsynqBigBossQueue(redisURL string, queueName string) (*AsynqBigBossQueue, error) {
	redisOpt, err := redisClientOptFromURL(redisURL)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(queueName) == "" {
		queueName = "big_boss_rewards"
	}
	return &AsynqBigBossQueue{
		client:    asynq.NewClient(redisOpt),
		queueName: queueName,
	}, nil
}

func (q *AsynqBigBossQueue) Close() error {
	if q == nil || q.client == nil {
		return nil
	}
	return q.client.Close()
}

func (q *AsynqBigBossQueue) EnqueueBigBossRewardDistribution(payload application.BigBossRewardsTaskPayload) error {
	start := time.Now().UTC()
	if q == nil || q.client == nil {
		bootstrap.IncAppMetric("worker.big_boss.enqueue_error_total")
		return errors.New("queue client not configured")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		bootstrap.IncAppMetric("worker.big_boss.enqueue_error_total")
		return err
	}
	task := asynq.NewTask(application.BigBossRewardsTaskType, body)
	_, err = q.client.Enqueue(task,
		asynq.Queue(q.queueName),
		asynq.MaxRetry(10),
		asynq.Timeout(2*time.Minute),
		asynq.Retention(24*time.Hour),
		asynq.TaskID(fmt.Sprintf("boss_reward_%d", payload.BossID)),
	)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "task id conflicts") {
		bootstrap.IncAppMetric("worker.big_boss.enqueue_idempotent_total")
		bootstrap.AddAppMetric("worker.big_boss.enqueue_duration_ms_total", time.Since(start).Milliseconds())
		return nil
	}
	if err != nil {
		bootstrap.IncAppMetric("worker.big_boss.enqueue_error_total")
		bootstrap.AddAppMetric("worker.big_boss.enqueue_duration_ms_total", time.Since(start).Milliseconds())
		return err
	}
	bootstrap.IncAppMetric("worker.big_boss.enqueue_total")
	bootstrap.AddAppMetric("worker.big_boss.enqueue_duration_ms_total", time.Since(start).Milliseconds())
	return err
}

type BigBossTaskHandler struct {
	service interface {
		ProcessRewardDistribution(bossID uint) error
	}
}

func NewBigBossTaskHandler(service *application.WeeklyBossService) *BigBossTaskHandler {
	return &BigBossTaskHandler{service: service}
}

func (h *BigBossTaskHandler) HandleDistributeRewardsTask(ctx context.Context, task *asynq.Task) error {
	start := time.Now().UTC()
	if h == nil || h.service == nil {
		bootstrap.IncAppMetric("worker.big_boss.job_error_total")
		return errors.New("big boss service not configured")
	}
	bootstrap.IncAppMetric("worker.big_boss.job_received_total")
	if retryCount, ok := asynq.GetRetryCount(ctx); ok && retryCount > 0 {
		bootstrap.IncAppMetric("worker.big_boss.job_retried_total")
		bootstrap.AddAppMetric("worker.big_boss.job_retry_attempt_total", int64(retryCount))
	}
	var payload application.BigBossRewardsTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		bootstrap.IncAppMetric("worker.big_boss.job_invalid_payload_total")
		bootstrap.AddAppMetric("worker.big_boss.job_duration_ms_total", time.Since(start).Milliseconds())
		return err
	}
	if err := h.service.ProcessRewardDistribution(payload.BossID); err != nil {
		bootstrap.IncAppMetric("worker.big_boss.job_error_total")
		bootstrap.AddAppMetric("worker.big_boss.job_duration_ms_total", time.Since(start).Milliseconds())
		return err
	}
	bootstrap.IncAppMetric("worker.big_boss.job_success_total")
	bootstrap.AddAppMetric("worker.big_boss.job_duration_ms_total", time.Since(start).Milliseconds())
	return nil
}

func NewAsynqServer(redisURL string, concurrency int) (*asynq.Server, error) {
	redisOpt, err := redisClientOptFromURL(redisURL)
	if err != nil {
		return nil, err
	}
	if concurrency <= 0 {
		concurrency = 10
	}
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: concurrency,
		Queues: map[string]int{
			"big_boss_rewards": 10,
		},
	})
	return srv, nil
}

func NewServeMux(handler *BigBossTaskHandler) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(application.BigBossRewardsTaskType, handler.HandleDistributeRewardsTask)
	return mux
}

func redisClientOptFromURL(raw string) (asynq.RedisClientOpt, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return asynq.RedisClientOpt{}, err
	}
	if u.Scheme != "redis" && u.Scheme != "rediss" {
		return asynq.RedisClientOpt{}, fmt.Errorf("invalid redis url scheme: %s", u.Scheme)
	}

	addr := u.Host
	if !strings.Contains(addr, ":") {
		addr = addr + ":6379"
	}

	password := ""
	if u.User != nil {
		if p, ok := u.User.Password(); ok {
			password = p
		}
	}

	db := 0
	dbPath := strings.TrimPrefix(u.Path, "/")
	if dbPath != "" {
		parsed, parseErr := strconv.Atoi(dbPath)
		if parseErr == nil && parsed >= 0 {
			db = parsed
		}
	}

	return asynq.RedisClientOpt{
		Addr:     addr,
		Password: password,
		DB:       db,
	}, nil
}
