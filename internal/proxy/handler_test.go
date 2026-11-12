package proxy

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"halo/internal/event"
	"halo/internal/model"
)

func newTestHandler(t *testing.T, upstream http.HandlerFunc) (*Handler, event.Bus) {
	t.Helper()
	up := httptest.NewServer(upstream)
	t.Cleanup(up.Close)
	u, err := url.Parse(up.URL)
	if err != nil {
		t.Fatal(err)
	}
	bus := event.New()
	return New(u, bus, slog.Default()), bus
}

func TestProxyJSONAndEvents(t *testing.T) {
	h, bus := newTestHandler(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "ping" {
			t.Errorf("upstream body = %q, want ping", body)
		}
		if r.URL.Path != "/v1/chat/completions" || r.URL.RawQuery != "q=1" {
			t.Errorf("upstream path/query = %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	sub := bus.Subscribe()
	defer sub.Close()

	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions?q=1", strings.NewReader("ping"))
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK || rec.Body.String() != `{"ok":true}` {
		t.Fatalf("forward failed: code=%d body=%q", rec.Code, rec.Body.String())
	}

	count := 0
	deadline := time.After(time.Second)
	for count < 2 {
		select {
		case e := <-sub.Events():
			count++
			if e.Type == model.EventRequest && e.Request.Headers["Authorization"] != "[redacted]" {
				t.Error("authorization header was not redacted")
			}
		case <-deadline:
			t.Fatalf("expected 2 events, got %d", count)
		}
	}
}

func TestProxySSEStream(t *testing.T) {
	h, _ := newTestHandler(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		f, _ := w.(http.Flusher)
		for i := 0; i < 3; i++ {
			_, _ = io.WriteString(w, "data: chunk\n\n")
			f.Flush()
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !strings.Contains(rec.Body.String(), "data: chunk") {
		t.Fatalf("stream not forwarded: %q", rec.Body.String())
	}
}

func TestUpstreamUnavailable(t *testing.T) {
	u, _ := url.Parse("http://127.0.0.1:1")
	h := New(u, event.New(), slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}
