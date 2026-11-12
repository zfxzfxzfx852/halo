package mcp

import (
	"testing"
	"time"

	"halo/internal/model"
)

func TestEventsFromSingleRequest(t *testing.T) {
	body := []byte(`{"jsonrpc":"2.0","id":7,"method":"tools/call","params":{"name":"query","arguments":{"k":"v"}}}`)
	events := EventsFromRequest(body, "sess", time.Now())
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].ToolCall.Name != "query" || events[0].Protocol != model.ProtocolMCPHTTP {
		t.Fatalf("unexpected event: %+v", events[0])
	}
}

func TestEventsFromBatch(t *testing.T) {
	body := []byte(`[
		{"jsonrpc":"2.0","id":1,"method":"tools/list"},
		{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"a"}},
		{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"b"}}
	]`)
	events := EventsFromRequest(body, "sess", time.Now())
	if len(events) != 2 {
		t.Fatalf("expected 2 tool call events, got %d", len(events))
	}
}

func TestNonMCPReturnsNil(t *testing.T) {
	body := []byte(`{"model":"gpt","messages":[]}`)
	if events := EventsFromRequest(body, "sess", time.Now()); len(events) != 0 {
		t.Fatalf("expected no events for non-MCP payload, got %d", len(events))
	}
}
