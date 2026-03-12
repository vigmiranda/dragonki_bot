package queue

import (
	"context"
	"errors"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/vigmiranda/dragonki/internal/application"
	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

type fakeRewardDistributionProcessor struct {
	processFn func(bossID uint) error
}

func (f *fakeRewardDistributionProcessor) ProcessRewardDistribution(bossID uint) error {
	if f.processFn == nil {
		return nil
	}
	return f.processFn(bossID)
}

func TestBigBossTaskHandler_HandleDistributeRewardsTask_SuccessRecordsMetric(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	handler := &BigBossTaskHandler{
		service: &fakeRewardDistributionProcessor{
			processFn: func(bossID uint) error {
				if bossID != 7 {
					t.Fatalf("expected bossID=7, got %d", bossID)
				}
				return nil
			},
		},
	}

	task := asynq.NewTask(application.BigBossRewardsTaskType, []byte(`{"boss_id":7}`))
	if err := handler.HandleDistributeRewardsTask(context.Background(), task); err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}

	snapshot := bootstrap.SnapshotAppMetrics()
	if snapshot["worker.big_boss.job_received_total"] != 1 {
		t.Fatalf("expected job received counter 1, got %d", snapshot["worker.big_boss.job_received_total"])
	}
	if snapshot["worker.big_boss.job_success_total"] != 1 {
		t.Fatalf("expected job success counter 1, got %d", snapshot["worker.big_boss.job_success_total"])
	}
	if snapshot["worker.big_boss.job_duration_ms_total"] < 0 {
		t.Fatalf("expected non-negative duration metric, got %d", snapshot["worker.big_boss.job_duration_ms_total"])
	}
}

func TestBigBossTaskHandler_HandleDistributeRewardsTask_ErrorRecordsMetric(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	handler := &BigBossTaskHandler{
		service: &fakeRewardDistributionProcessor{
			processFn: func(bossID uint) error {
				return errors.New("process fail")
			},
		},
	}

	task := asynq.NewTask(application.BigBossRewardsTaskType, []byte(`{"boss_id":7}`))
	if err := handler.HandleDistributeRewardsTask(context.Background(), task); err == nil {
		t.Fatal("expected handler error")
	}

	snapshot := bootstrap.SnapshotAppMetrics()
	if snapshot["worker.big_boss.job_error_total"] != 1 {
		t.Fatalf("expected job error counter 1, got %d", snapshot["worker.big_boss.job_error_total"])
	}
	if snapshot["worker.big_boss.job_duration_ms_total"] < 0 {
		t.Fatalf("expected non-negative duration metric, got %d", snapshot["worker.big_boss.job_duration_ms_total"])
	}
}
