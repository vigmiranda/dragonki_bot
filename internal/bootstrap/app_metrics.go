package bootstrap

import "sync"

type appMetricsRegistry struct {
	mu       sync.RWMutex
	counters map[string]int64
}

var defaultAppMetrics = &appMetricsRegistry{
	counters: make(map[string]int64),
}

func IncAppMetric(name string) {
	AddAppMetric(name, 1)
}

func AddAppMetric(name string, delta int64) {
	if name == "" || delta == 0 {
		return
	}

	defaultAppMetrics.mu.Lock()
	defer defaultAppMetrics.mu.Unlock()
	defaultAppMetrics.counters[name] += delta
}

func SnapshotAppMetrics() map[string]int64 {
	defaultAppMetrics.mu.RLock()
	defer defaultAppMetrics.mu.RUnlock()

	snapshot := make(map[string]int64, len(defaultAppMetrics.counters))
	for key, value := range defaultAppMetrics.counters {
		snapshot[key] = value
	}
	return snapshot
}

func ResetAppMetricsForTesting() {
	defaultAppMetrics.mu.Lock()
	defer defaultAppMetrics.mu.Unlock()
	defaultAppMetrics.counters = make(map[string]int64)
}
