package policy

import (
	"encoding/json"
	"testing"

	"halo/internal/model"
)

func toolEvent(name string, args map[string]any, risk model.RiskLevel) *model.Event {
	raw, _ := json.Marshal(args)
	return &model.Event{
		Type:     model.EventToolCall,
		ToolCall: &model.ToolCall{Name: name, Arguments: raw},
		Risk:     risk,
	}
}

func TestDefaultFallback(t *testing.T) {
	p := &model.Policy{Default: model.DecisionDeny}
	got := New(p).Evaluate(toolEvent("anything", nil, model.RiskNone))
	if got.Decision != model.DecisionDeny || got.RuleName != "" {
		t.Fatalf("got %+v", got)
	}
}

func TestToolMatch(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionAllow,
		Rules:   []model.PolicyRule{{Name: "block rm", Tool: "rm", Decision: model.DecisionDeny}},
	}
	e := New(p)
	if got := e.Evaluate(toolEvent("rm", nil, model.RiskNone)); got.Decision != model.DecisionDeny {
		t.Fatal("rm should deny")
	}
	if got := e.Evaluate(toolEvent("ls", nil, model.RiskNone)); got.Decision != model.DecisionAllow {
		t.Fatal("ls should allow")
	}
}

func TestWildcard(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionAllow,
		Rules:   []model.PolicyRule{{Name: "fs", Tool: "fs_*", Decision: model.DecisionAsk}},
	}
	if got := New(p).Evaluate(toolEvent("fs_write", nil, model.RiskNone)); got.Decision != model.DecisionAsk {
		t.Fatal("wildcard should ask")
	}
}

func TestPathMatch(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionAllow,
		Rules: []model.PolicyRule{{
			Name: "protect etc", Tool: "read_file",
			Paths: []string{"/etc/*"}, Decision: model.DecisionDeny,
		}},
	}
	e := New(p)
	if got := e.Evaluate(toolEvent("read_file", map[string]any{"path": "/etc/passwd"}, model.RiskNone)); got.Decision != model.DecisionDeny {
		t.Fatal("etc path deny")
	}
	if got := e.Evaluate(toolEvent("read_file", map[string]any{"path": "/tmp/a"}, model.RiskNone)); got.Decision != model.DecisionAllow {
		t.Fatal("tmp path allow")
	}
}

func TestCommandMatch(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionAllow,
		Rules:   []model.PolicyRule{{Name: "no curl", Command: "curl *", Decision: model.DecisionDeny}},
	}
	if got := New(p).Evaluate(toolEvent("shell", map[string]any{"command": "curl http://x"}, model.RiskNone)); got.Decision != model.DecisionDeny {
		t.Fatal("command deny")
	}
}

func TestArgsMatch(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionAllow,
		Rules: []model.PolicyRule{{
			Name: "url", Tool: "http_request",
			Args: map[string]string{"url": "http://*"}, Decision: model.DecisionAsk,
		}},
	}
	if got := New(p).Evaluate(toolEvent("http_request", map[string]any{"url": "http://x"}, model.RiskNone)); got.Decision != model.DecisionAsk {
		t.Fatal("args ask")
	}
}

func TestMaxRisk(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionDeny,
		Rules: []model.PolicyRule{{
			Name: "low only", Tool: "x", MaxRisk: model.RiskLow, Decision: model.DecisionAllow,
		}},
	}
	e := New(p)
	if got := e.Evaluate(toolEvent("x", nil, model.RiskLow)); got.Decision != model.DecisionAllow {
		t.Fatal("low allow")
	}
	if got := e.Evaluate(toolEvent("x", nil, model.RiskHigh)); got.Decision != model.DecisionDeny {
		t.Fatal("high falls back to deny")
	}
}

func TestFirstMatchWins(t *testing.T) {
	p := &model.Policy{
		Default: model.DecisionDeny,
		Rules: []model.PolicyRule{
			{Name: "first", Tool: "x", Decision: model.DecisionAsk},
			{Name: "second", Tool: "x", Decision: model.DecisionAllow},
		},
	}
	if got := New(p).Evaluate(toolEvent("x", nil, model.RiskNone)); got.Decision != model.DecisionAsk {
		t.Fatal("first rule wins")
	}
}

func TestHotReplace(t *testing.T) {
	e := New(&model.Policy{Default: model.DecisionAllow})
	if got := e.Evaluate(toolEvent("x", nil, model.RiskNone)); got.Decision != model.DecisionAllow {
		t.Fatal("initial allow")
	}
	e.Replace(&model.Policy{Default: model.DecisionDeny})
	if got := e.Evaluate(toolEvent("x", nil, model.RiskNone)); got.Decision != model.DecisionDeny {
		t.Fatal("after replace deny")
	}
}

func TestParseAndValidate(t *testing.T) {
	yamlDoc := []byte(`
version: "1"
default: deny
rules:
  - name: safe
    tool: ls
    decision: allow
`)
	p, err := Parse(yamlDoc)
	if err != nil {
		t.Fatal(err)
	}
	if p.Default != model.DecisionDeny || len(p.Rules) != 1 {
		t.Fatalf("parse %+v", p)
	}
	if _, err := Parse([]byte("default: bogus")); err == nil {
		t.Fatal("expected validation error")
	}
}
