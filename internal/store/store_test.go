package store

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"halo/internal/model"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func sampleEvent(id string, typ model.EventType, risk model.RiskLevel, ts time.Time) *model.Event {
	return &model.Event{
		ID:        id,
		Type:      typ,
		Protocol:  model.ProtocolHTTP,
		SessionID: "s1",
		Timestamp: ts,
		Risk:      risk,
	}
}

func TestInsertAndGet(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	want := sampleEvent("e1", model.EventRequest, model.RiskNone, time.Now())
	if err := s.InsertEvent(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := s.GetEvent(ctx, "e1")
	if err != nil || got.ID != "e1" {
		t.Fatalf("get failed: %v %+v", err, got)
	}
}

func TestGetNotFound(t *testing.T) {
	s := newTestStore(t)
	if _, err := s.GetEvent(context.Background(), "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListFilterAndPaging(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	base := time.Now()
	for i := 0; i < 5; i++ {
		if err := s.InsertEvent(ctx, sampleEvent(
			string(rune('a'+i)), model.EventRequest, model.RiskNone, base.Add(time.Duration(i)*time.Second))); err != nil {
			t.Fatal(err)
		}
	}
	events, total, err := s.ListEvents(ctx, Filter{Limit: 2, Offset: 0})
	if err != nil || total != 5 || len(events) != 2 {
		t.Fatalf("list failed: total=%d len=%d err=%v", total, len(events), err)
	}
	// 降序：最新的在前
	if events[0].ID != "e" {
		t.Fatalf("order wrong, first id = %q", events[0].ID)
	}
}

func TestStats(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	_ = s.InsertEvent(ctx, sampleEvent("1", model.EventRequest, model.RiskNone, time.Now()))
	_ = s.InsertEvent(ctx, sampleEvent("2", model.EventResponse, model.RiskHigh, time.Now()))

	st, err := s.Stats(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if st.Total != 2 || st.ByType["request"] != 1 || st.ByRisk["high"] != 1 {
		t.Fatalf("stats wrong: %+v", st)
	}
}

func TestPurge(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	old := time.Now().Add(-time.Hour)
	recent := time.Now()
	_ = s.InsertEvent(ctx, sampleEvent("old", model.EventRequest, model.RiskNone, old))
	_ = s.InsertEvent(ctx, sampleEvent("new", model.EventRequest, model.RiskNone, recent))

	n, err := s.PurgeBefore(ctx, time.Now().Add(-30*time.Minute))
	if err != nil || n != 1 {
		t.Fatalf("purge n=%d err=%v", n, err)
	}
	if _, err := s.GetEvent(ctx, "old"); !errors.Is(err, ErrNotFound) {
		t.Fatal("old event should be purged")
	}
}

func TestAuditRecords(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()

	if _, ok, err := s.LastAuditRecord(ctx); err != nil || ok {
		t.Fatalf("empty last: ok=%v err=%v", ok, err)
	}

	recs := []model.AuditRecord{
		{Seq: 1, EventID: "e1", Hash: "h1", PrevHash: "", Timestamp: time.Now()},
		{Seq: 2, EventID: "e2", Hash: "h2", PrevHash: "h1", Timestamp: time.Now()},
	}
	for _, r := range recs {
		if err := s.InsertAuditRecord(ctx, r); err != nil {
			t.Fatal(err)
		}
	}

	got, err := s.ListAuditRecords(ctx)
	if err != nil || len(got) != 2 {
		t.Fatalf("list audit: len=%d err=%v", len(got), err)
	}
	if got[0].Hash != "h1" || got[1].PrevHash != "h1" {
		t.Fatalf("audit link wrong: %+v", got)
	}

	last, ok, err := s.LastAuditRecord(ctx)
	if err != nil || !ok || last.Seq != 2 {
		t.Fatalf("last: %+v ok=%v err=%v", last, ok, err)
	}
}

func TestBuildWhere(t *testing.T) {
	if w, args := buildWhere(Filter{}); w != "" || len(args) != 0 {
		t.Fatalf("empty filter: w=%q args=%v", w, args)
	}
	w, args := buildWhere(Filter{
		Type: "request", Protocol: "http", SessionID: "s1", RiskMin: model.RiskLow,
	})
	if !strings.Contains(w, "WHERE") || len(args) != 7 {
		t.Fatalf("full filter: w=%q args=%d", w, len(args))
	}
}

func TestBuildWhereRiskMinThreshold(t *testing.T) {
	levels := riskLevelsAtOrAbove(model.RiskHigh)
	if len(levels) != 2 || levels[0] != model.RiskHigh || levels[1] != model.RiskCritical {
		t.Fatalf("high threshold should be [high critical], got %v", levels)
	}
	w, args := buildWhere(Filter{RiskMin: model.RiskHigh})
	if strings.Count(w, "?") != 2 || len(args) != 2 {
		t.Fatalf("high filter should use 2 placeholders, got w=%q args=%v", w, args)
	}
	crit := riskLevelsAtOrAbove(model.RiskCritical)
	if len(crit) != 1 || crit[0] != model.RiskCritical {
		t.Fatalf("critical threshold should only include critical, got %v", crit)
	}
}
