package server

import (
	"context"
	"log/slog"
	"net/http/httptest"
	"testing"
	"time"

	"halo/internal/config"
	"halo/internal/model"
	"halo/internal/policy"
)

func TestSessionsAggregation(t *testing.T) {
	env := newTestEnv(t)
	now := time.Now()
	same := []*model.Event{
		{ID: "a1", Type: model.EventRequest, Protocol: model.ProtocolHTTP, SessionID: "alpha", Agent: "agent-a", Timestamp: now.Add(-2 * time.Minute), Risk: model.RiskNone},
		{ID: "a2", Type: model.EventToolCall, Protocol: model.ProtocolMCPStdio, SessionID: "alpha", Agent: "agent-a", Timestamp: now.Add(-time.Minute), Risk: model.RiskHigh},
		{ID: "b1", Type: model.EventRequest, Protocol: model.ProtocolHTTP, SessionID: "beta", Agent: "agent-b", Timestamp: now, Risk: model.RiskLow},
	}
	for _, e := range same {
		if err := env.st.InsertEvent(context.Background(), e); err != nil {
			t.Fatal(err)
		}
	}

	m := getJSON(t, env.base+"/api/v1/sessions?limit=10")
	items := m["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("sessions = %d, want 2", len(items))
	}
	first := items[0].(map[string]any)
	if first["sessionId"] != "beta" {
		t.Fatalf("most recent session = %v, want beta", first["sessionId"])
	}
	second := items[1].(map[string]any)
	if second["maxRisk"] != "high" || int(second["events"].(float64)) != 2 {
		t.Fatalf("alpha aggregate = %v", second)
	}
	if int(second["toolCalls"].(float64)) != 1 {
		t.Fatalf("alpha toolCalls = %v", second["toolCalls"])
	}
}

func TestFindingsAggregation(t *testing.T) {
	env := newTestEnv(t)
	ev := &model.Event{ID: "f1", Type: model.EventRequest, Protocol: model.ProtocolHTTP,
		SessionID: "s", Timestamp: time.Now(), Risk: model.RiskHigh,
		Findings: []model.Finding{
			{Rule: "injection.ignore_instructions", Level: model.RiskHigh, Message: "x"},
			{Rule: "secret.aws_key", Level: model.RiskCritical, Message: "y"},
		}}
	if err := env.st.InsertEvent(context.Background(), ev); err != nil {
		t.Fatal(err)
	}

	m := getJSON(t, env.base+"/api/v1/findings")
	items := m["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("findings = %d, want 2", len(items))
	}
	top := items[0].(map[string]any)
	if top["maxRisk"] != "critical" && top["maxRisk"] != "high" {
		t.Fatalf("finding risk = %v", top["maxRisk"])
	}
}

func TestPoliciesEndpoint(t *testing.T) {
	env := newTestEnv(t)
	// 未绑定策略引擎时返回空规则集，而不是报错。
	m := getJSON(t, env.base+"/api/v1/policies")
	if m["rules"] == nil {
		t.Fatalf("rules should be empty array, got %v", m)
	}
}

func TestPoliciesWithEngine(t *testing.T) {
	env := newTestEnv(t)
	cfg := &config.Config{Host: "127.0.0.1", Port: "0", Mode: config.ModeGuard}
	s := New(cfg, slog.Default(), env.st, NewHub(slog.Default()), env.mgr, nil)
	s.SetPolicy(policy.New(&model.Policy{Default: model.DecisionAllow, Rules: []model.PolicyRule{
		{Name: "block-rm", Command: "rm *", Decision: model.DecisionDeny},
	}}))
	ts := httptest.NewServer(s.http.Handler)
	t.Cleanup(ts.Close)

	m := getJSON(t, ts.URL+"/api/v1/policies")
	if m["default"] != "allow" {
		t.Fatalf("default = %v", m["default"])
	}
	rules := m["rules"].([]any)
	if len(rules) != 1 {
		t.Fatalf("rules = %d, want 1", len(rules))
	}
	if rules[0].(map[string]any)["name"] != "block-rm" {
		t.Fatalf("rule = %v", rules[0])
	}
}

func TestRuntimeInfo(t *testing.T) {
	env := newTestEnv(t)
	m := getJSON(t, env.base+"/api/v1/runtime")
	if m["mode"] != "observe" {
		t.Fatalf("mode = %v", m["mode"])
	}
	if _, ok := m["retentionDays"]; !ok {
		t.Fatalf("retentionDays missing: %v", m)
	}
}
