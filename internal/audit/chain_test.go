package audit

import (
	"strings"
	"testing"
	"time"

	"halo/internal/model"
)

func mkEvent(id string) *model.Event {
	return &model.Event{
		ID: id, Type: model.EventRequest, Protocol: model.ProtocolHTTP,
		SessionID: "default", Timestamp: time.Now(), Method: "GET /x", Risk: model.RiskNone,
	}
}

func buildChain(t *testing.T, n int) ([]model.AuditRecord, map[string]*model.Event) {
	c := NewChain(0, "")
	events := map[string]*model.Event{}
	var recs []model.AuditRecord
	for i := 0; i < n; i++ {
		e := mkEvent("evt_" + string(rune('a'+i)))
		events[e.ID] = e
		rec, err := c.Append(e)
		if err != nil {
			t.Fatal(err)
		}
		recs = append(recs, rec)
	}
	return recs, events
}

func TestChainLinks(t *testing.T) {
	recs, _ := buildChain(t, 3)
	if recs[0].PrevHash != "" {
		t.Fatal("genesis prev must be empty")
	}
	if recs[1].PrevHash != recs[0].Hash || recs[2].PrevHash != recs[1].Hash {
		t.Fatal("hash links broken")
	}
	if recs[0].Seq != 1 || recs[2].Seq != 3 {
		t.Fatal("seq wrong")
	}
}

func TestVerifyOK(t *testing.T) {
	recs, events := buildChain(t, 3)
	rep := Verify(recs, events)
	if !rep.Valid || rep.Checked != 3 || rep.MissingEvents != 0 {
		t.Fatalf("unexpected report: %+v", rep)
	}
}

func TestVerifyTamperedEvent(t *testing.T) {
	recs, events := buildChain(t, 3)
	bad := *events["evt_b"]
	bad.Method = "GET /hacked"
	events["evt_b"] = &bad
	if rep := Verify(recs, events); rep.Valid {
		t.Fatal("expected tamper detection")
	}
}

func TestVerifyBrokenLink(t *testing.T) {
	recs, events := buildChain(t, 3)
	recs[1].PrevHash = "deadbeef"
	if rep := Verify(recs, events); rep.Valid {
		t.Fatal("expected broken link")
	}
}

func TestVerifyPurgedEvent(t *testing.T) {
	recs, events := buildChain(t, 3)
	delete(events, "evt_b")
	rep := Verify(recs, events)
	if !rep.Valid {
		t.Fatalf("purged event should keep chain valid, got %+v", rep)
	}
	if rep.MissingEvents != 1 || rep.Checked != 2 {
		t.Fatalf("expected 1 missing and 2 checked, got %+v", rep)
	}
}

func TestExportJSONL(t *testing.T) {
	recs, _ := buildChain(t, 3)
	data, err := ExportJSONL(recs)
	if err != nil {
		t.Fatal(err)
	}
	if lines := strings.Count(strings.TrimSpace(string(data)), "\n") + 1; lines != 3 {
		t.Fatalf("lines=%d", lines)
	}
}
