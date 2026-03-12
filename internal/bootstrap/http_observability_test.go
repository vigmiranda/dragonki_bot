package bootstrap

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPObservability_WrapRecordsRouteMetrics(t *testing.T) {
	obs := NewHTTPObservability("bot")

	handler := obs.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	}))

	req := httptest.NewRequest(http.MethodPost, "/ok", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}

	snapshot := obs.snapshot()
	route := snapshot.Routes["/ok"]
	if route.RequestCount != 1 {
		t.Fatalf("expected request count 1, got %d", route.RequestCount)
	}
	if route.MethodCount[http.MethodPost] != 1 {
		t.Fatalf("expected method count 1, got %d", route.MethodCount[http.MethodPost])
	}
	if route.StatusCount["202"] != 1 {
		t.Fatalf("expected status count 1, got %d", route.StatusCount["202"])
	}
}

func TestHTTPObservability_RegisterDebugRoutesExposesVarsAndPprof(t *testing.T) {
	ResetAppMetricsForTesting()
	IncAppMetric("telegram.callback.received_total")

	obs := NewHTTPObservability("webhook")
	mux := http.NewServeMux()
	obs.RegisterDebugRoutes(mux, true)

	varsReq := httptest.NewRequest(http.MethodGet, "/debug/vars", nil)
	varsRec := httptest.NewRecorder()
	mux.ServeHTTP(varsRec, varsReq)

	if varsRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /debug/vars, got %d", varsRec.Code)
	}

	var payload httpsObservabilitySnapshot
	if err := json.Unmarshal(varsRec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode vars payload: %v", err)
	}
	if payload.Service != "webhook" {
		t.Fatalf("expected service webhook, got %q", payload.Service)
	}
	if payload.Counters["telegram.callback.received_total"] != 1 {
		t.Fatalf("expected app counter in vars payload, got %d", payload.Counters["telegram.callback.received_total"])
	}

	pprofReq := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	pprofRec := httptest.NewRecorder()
	mux.ServeHTTP(pprofRec, pprofReq)

	if pprofRec.Code != http.StatusOK {
		t.Fatalf("expected 200 from /debug/pprof/, got %d", pprofRec.Code)
	}
}
