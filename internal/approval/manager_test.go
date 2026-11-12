package approval

import (
	"context"
	"testing"
	"time"

	"halo/internal/model"
)

func sampleEvent() *model.Event {
	return &model.Event{ID: "evt_1", Method: "POST /x", Risk: model.RiskHigh}
}

func waitPendingID(m *Manager) string {
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if p := m.Pending(); len(p) > 0 {
			return p[0].ID
		}
		time.Sleep(5 * time.Millisecond)
	}
	return ""
}

func TestApproveAllow(t *testing.T) {
	m := NewManager(nil, time.Second)
	got := make(chan model.Decision, 1)
	go func() { got <- m.Approve(context.Background(), sampleEvent()) }()

	id := waitPendingID(m)
	if id == "" {
		t.Fatal("no pending request")
	}
	if err := m.Resolve(id, model.DecisionAllow); err != nil {
		t.Fatal(err)
	}
	if d := <-got; d != model.DecisionAllow {
		t.Fatalf("got %v", d)
	}
}

func TestResolveDeny(t *testing.T) {
	m := NewManager(nil, time.Second)
	got := make(chan model.Decision, 1)
	go func() { got <- m.Approve(context.Background(), sampleEvent()) }()

	if err := m.Resolve(waitPendingID(m), model.DecisionDeny); err != nil {
		t.Fatal(err)
	}
	if d := <-got; d != model.DecisionDeny {
		t.Fatalf("got %v", d)
	}
}

func TestTimeoutDenies(t *testing.T) {
	m := NewManager(nil, 30*time.Millisecond)
	start := time.Now()
	if d := m.Approve(context.Background(), sampleEvent()); d != model.DecisionDeny {
		t.Fatalf("got %v", d)
	}
	if time.Since(start) < 30*time.Millisecond {
		t.Fatal("returned too early")
	}
}

func TestResolveNotFound(t *testing.T) {
	m := NewManager(nil, time.Second)
	if err := m.Resolve("nope", model.DecisionAllow); err != ErrNotFound {
		t.Fatalf("got %v", err)
	}
}
