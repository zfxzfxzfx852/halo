package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"halo/internal/approval"
	"halo/internal/config"
	"halo/internal/model"
	"halo/internal/store"
)

type testEnv struct {
	st   *store.Store
	mgr  *approval.Manager
	base string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	cfg := &config.Config{Host: "127.0.0.1", Port: "0", Mode: config.ModeObserve}
	hub := NewHub(slog.Default())
	mgr := approval.NewManager(nil, time.Minute)
	s := New(cfg, slog.Default(), st, hub, mgr, nil)

	ts := httptest.NewServer(s.http.Handler)
	t.Cleanup(ts.Close)
	return &testEnv{st: st, mgr: mgr, base: ts.URL}
}

func TestHealthAndVersion(t *testing.T) {
	env := newTestEnv(t)
	m := getJSON(t, env.base+"/healthz")
	if m["status"] != "ok" {
		t.Fatalf("healthz = %v", m)
	}
	if code := getStatus(t, env.base+"/api/v1/version"); code != http.StatusOK {
		t.Fatalf("version status = %d", code)
	}
}

func TestListGetEvents(t *testing.T) {
	env := newTestEnv(t)
	ev1 := &model.Event{ID: "evt-1", Type: model.EventRequest, Protocol: model.ProtocolHTTP,
		SessionID: "s1", Timestamp: time.Now(), Method: "GET /a"}
	ev2 := &model.Event{ID: "evt-2", Type: model.EventResponse, Protocol: model.ProtocolHTTP,
		SessionID: "s1", Timestamp: time.Now(), Risk: model.RiskHigh}
	if err := env.st.InsertEvent(context.Background(), ev1); err != nil {
		t.Fatal(err)
	}
	if err := env.st.InsertEvent(context.Background(), ev2); err != nil {
		t.Fatal(err)
	}

	m := getJSON(t, env.base+"/api/v1/events?limit=10")
	if int(m["total"].(float64)) != 2 {
		t.Fatalf("total = %v, want 2", m["total"])
	}
	if len(m["items"].([]any)) != 2 {
		t.Fatalf("items len = %d", len(m["items"].([]any)))
	}

	if code := getStatus(t, env.base+"/api/v1/events/evt-2"); code != http.StatusOK {
		t.Fatalf("existing event status = %d", code)
	}
	if code := getStatus(t, env.base+"/api/v1/events/nope"); code != http.StatusNotFound {
		t.Fatalf("missing event status = %d, want 404", code)
	}
}

func TestStats(t *testing.T) {
	env := newTestEnv(t)
	if code := getStatus(t, env.base+"/api/v1/stats"); code != http.StatusOK {
		t.Fatalf("stats status = %d", code)
	}
}

func TestApprovals(t *testing.T) {
	env := newTestEnv(t)

	done := make(chan struct{})
	go func() {
		env.mgr.Approve(context.Background(),
			&model.Event{ID: "evt-ap", Method: "POST /danger", Risk: model.RiskHigh})
		close(done)
	}()
	deadline := time.Now().Add(2 * time.Second)
	for len(env.mgr.Pending()) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	pending := env.mgr.Pending()
	if len(pending) != 1 {
		t.Fatalf("pending = %d, want 1", len(pending))
	}
	id := pending[0].ID

	m := getJSON(t, env.base+"/api/v1/approvals")
	if int(m["total"].(float64)) != 1 {
		t.Fatalf("approvals total = %v", m["total"])
	}

	if code := postJSON(t, env.base+"/api/v1/approvals/"+id, map[string]string{"decision": "maybe"}); code != http.StatusBadRequest {
		t.Fatalf("invalid decision status = %d, want 400", code)
	}
	if code := postRaw(t, env.base+"/api/v1/approvals/"+id, "{bad"); code != http.StatusBadRequest {
		t.Fatalf("bad body status = %d, want 400", code)
	}
	if code := postJSON(t, env.base+"/api/v1/approvals/"+id, map[string]string{"decision": "allow"}); code != http.StatusOK {
		t.Fatalf("resolve status = %d, want 200", code)
	}
	<-done
	if code := postJSON(t, env.base+"/api/v1/approvals/"+id, map[string]string{"decision": "deny"}); code != http.StatusNotFound {
		t.Fatalf("re-resolve status = %d, want 404", code)
	}
}

func TestAuditEndpoints(t *testing.T) {
	env := newTestEnv(t)
	m := getJSON(t, env.base+"/api/v1/audit/verify")
	if valid, _ := m["valid"].(bool); !valid {
		t.Fatalf("empty chain should be valid: %v", m)
	}
	if code := getStatus(t, env.base+"/api/v1/audit/export"); code != http.StatusOK {
		t.Fatalf("export status = %d", code)
	}
}

func TestRootRedirect(t *testing.T) {
	env := newTestEnv(t)
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	resp, err := client.Get(env.base + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusTemporaryRedirect {
		t.Fatalf("root status = %d, want 307", resp.StatusCode)
	}
}

// ---- helpers ----

func getStatus(t *testing.T, url string) int {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

func getJSON(t *testing.T, url string) map[string]any {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var m map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatal(err)
	}
	return m
}

func postJSON(t *testing.T, url string, v any) int {
	t.Helper()
	data, _ := json.Marshal(v)
	return postRaw(t, url, string(data))
}

func postRaw(t *testing.T, url, raw string) int {
	t.Helper()
	resp, err := http.Post(url, "application/json", strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}
