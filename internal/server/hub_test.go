package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"

	"halo/internal/model"
)

func TestHubBroadcast(t *testing.T) {
	hub := NewHub(slog.Default())
	srv := httptest.NewServer(http.HandlerFunc(hub.ServeWS))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, srv.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.CloseNow()

	// 等待连接注册到 Hub。
	deadline := time.Now().Add(2 * time.Second)
	for hub.ClientCount() == 0 && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if hub.ClientCount() != 1 {
		t.Fatalf("client count = %d, want 1", hub.ClientCount())
	}

	hub.Broadcast(&model.Event{ID: "evt-x", Type: model.EventAlert})

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var got model.Event
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ID != "evt-x" {
		t.Fatalf("got id %q, want evt-x", got.ID)
	}
}
