package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vigmiranda/dragonki/internal/bootstrap"
)

func TestObservabilityLogger_TraceRecordsQueryMetrics(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	logger := &observabilityLogger{}
	logger.Trace(context.Background(), time.Now().Add(-250*time.Millisecond), func() (string, int64) {
		return "SELECT * FROM players WHERE id = 1", 1
	}, nil)

	snapshot := bootstrap.SnapshotAppMetrics()
	if snapshot["db.query.total"] != 1 {
		t.Fatalf("expected db.query.total 1, got %d", snapshot["db.query.total"])
	}
	if snapshot["db.query.select_total"] != 1 {
		t.Fatalf("expected db.query.select_total 1, got %d", snapshot["db.query.select_total"])
	}
	if snapshot["db.query.slow_total"] != 1 {
		t.Fatalf("expected db.query.slow_total 1, got %d", snapshot["db.query.slow_total"])
	}
}

func TestObservabilityLogger_TraceRecordsErrors(t *testing.T) {
	bootstrap.ResetAppMetricsForTesting()

	logger := &observabilityLogger{}
	logger.Trace(context.Background(), time.Now(), func() (string, int64) {
		return "UPDATE players SET gold = gold + 1", 1
	}, errors.New("db fail"))

	snapshot := bootstrap.SnapshotAppMetrics()
	if snapshot["db.query.error_total"] != 1 {
		t.Fatalf("expected db.query.error_total 1, got %d", snapshot["db.query.error_total"])
	}
	if snapshot["db.query.update_total"] != 1 {
		t.Fatalf("expected db.query.update_total 1, got %d", snapshot["db.query.update_total"])
	}
}
