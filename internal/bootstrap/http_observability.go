package bootstrap

import (
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"
)

type HTTPObservability struct {
	service   string
	startedAt time.Time

	mu         sync.RWMutex
	totalCount int64
	routes     map[string]*httpRouteMetrics
}

type httpRouteMetrics struct {
	RequestCount    int64
	Inflight        int64
	DurationMsTotal float64
	MethodCount     map[string]int64
	StatusCount     map[string]int64
}

type httpRouteMetricsSnapshot struct {
	RequestCount    int64            `json:"request_count"`
	Inflight        int64            `json:"inflight"`
	DurationMsTotal float64          `json:"duration_ms_total"`
	DurationMsAvg   float64          `json:"duration_ms_avg"`
	MethodCount     map[string]int64 `json:"method_count"`
	StatusCount     map[string]int64 `json:"status_count"`
}

type httpsObservabilitySnapshot struct {
	Service       string                              `json:"service"`
	StartedAt     time.Time                           `json:"started_at"`
	UptimeSeconds float64                             `json:"uptime_seconds"`
	Now           time.Time                           `json:"now"`
	TotalRequests int64                               `json:"total_requests"`
	Goroutines    int                                 `json:"goroutines"`
	Counters      map[string]int64                    `json:"counters"`
	Routes        map[string]httpRouteMetricsSnapshot `json:"routes"`
}

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCapturingResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func NewHTTPObservability(service string) *HTTPObservability {
	return &HTTPObservability{
		service:   service,
		startedAt: time.Now().UTC(),
		routes:    make(map[string]*httpRouteMetrics),
	}
}

func (o *HTTPObservability) Wrap(next http.Handler) http.Handler {
	if o == nil || next == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		path := r.URL.Path

		o.markInflight(path, 1)
		defer o.markInflight(path, -1)

		rec := &statusCapturingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(rec, r)

		o.record(path, r.Method, rec.statusCode, time.Since(start))
	})
}

func (o *HTTPObservability) RegisterDebugRoutes(mux *http.ServeMux, enablePprof bool) {
	if o == nil || mux == nil {
		return
	}

	mux.Handle("/debug/vars", o)
	if !enablePprof {
		return
	}

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func (o *HTTPObservability) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if o == nil {
		http.Error(w, "observability disabled", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(o.snapshot())
}

func (o *HTTPObservability) snapshot() httpsObservabilitySnapshot {
	now := time.Now().UTC()

	o.mu.RLock()
	defer o.mu.RUnlock()

	routes := make(map[string]httpRouteMetricsSnapshot, len(o.routes))
	routeNames := make([]string, 0, len(o.routes))
	for path := range o.routes {
		routeNames = append(routeNames, path)
	}
	sort.Strings(routeNames)

	for _, path := range routeNames {
		metrics := o.routes[path]
		if metrics == nil {
			continue
		}

		methodCount := cloneMetricMap(metrics.MethodCount)
		statusCount := cloneMetricMap(metrics.StatusCount)
		avgDuration := 0.0
		if metrics.RequestCount > 0 {
			avgDuration = metrics.DurationMsTotal / float64(metrics.RequestCount)
		}

		routes[path] = httpRouteMetricsSnapshot{
			RequestCount:    metrics.RequestCount,
			Inflight:        metrics.Inflight,
			DurationMsTotal: metrics.DurationMsTotal,
			DurationMsAvg:   avgDuration,
			MethodCount:     methodCount,
			StatusCount:     statusCount,
		}
	}

	return httpsObservabilitySnapshot{
		Service:       o.service,
		StartedAt:     o.startedAt,
		UptimeSeconds: now.Sub(o.startedAt).Seconds(),
		Now:           now,
		TotalRequests: o.totalCount,
		Goroutines:    runtime.NumGoroutine(),
		Counters:      SnapshotAppMetrics(),
		Routes:        routes,
	}
}

func (o *HTTPObservability) markInflight(path string, delta int64) {
	o.mu.Lock()
	defer o.mu.Unlock()

	metrics := o.routeMetricsLocked(path)
	metrics.Inflight += delta
	if metrics.Inflight < 0 {
		metrics.Inflight = 0
	}
}

func (o *HTTPObservability) record(path string, method string, statusCode int, duration time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.totalCount++
	metrics := o.routeMetricsLocked(path)
	metrics.RequestCount++
	metrics.DurationMsTotal += float64(duration) / float64(time.Millisecond)
	metrics.MethodCount[method]++
	metrics.StatusCount[strconv.Itoa(statusCode)]++
}

func (o *HTTPObservability) routeMetricsLocked(path string) *httpRouteMetrics {
	metrics := o.routes[path]
	if metrics != nil {
		return metrics
	}

	metrics = &httpRouteMetrics{
		MethodCount: make(map[string]int64),
		StatusCount: make(map[string]int64),
	}
	o.routes[path] = metrics
	return metrics
}

func cloneMetricMap(src map[string]int64) map[string]int64 {
	if len(src) == 0 {
		return map[string]int64{}
	}

	dst := make(map[string]int64, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}
