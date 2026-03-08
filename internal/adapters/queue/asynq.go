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
	if q == nil || q.client == nil {
		return errors.New("queue client not configured")
	}
	body, err := json.Marshal(payload)
	if err != nil {
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
		return nil
	}
	return err
}

type BigBossTaskHandler struct {
	service *application.WeeklyBossService
}

func NewBigBossTaskHandler(service *application.WeeklyBossService) *BigBossTaskHandler {
	return &BigBossTaskHandler{service: service}
}

func (h *BigBossTaskHandler) HandleDistributeRewardsTask(ctx context.Context, task *asynq.Task) error {
	if h == nil || h.service == nil {
		return errors.New("big boss service not configured")
	}
	var payload application.BigBossRewardsTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}
	return h.service.ProcessRewardDistribution(payload.BossID)
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
