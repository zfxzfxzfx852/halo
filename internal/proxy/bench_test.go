package proxy

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"vigil/internal/event"
)

// BenchmarkProxyThroughput 测量代理在普通 JSON 转发下的单请求开销与内存分配。
func BenchmarkProxyThroughput(b *testing.B) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer up.Close()

	u, err := url.Parse(up.URL)
	if err != nil {
		b.Fatal(err)
	}
	h := New(u, event.New(), slog.New(slog.NewTextHandler(io.Discard, nil)))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader("{}"))
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			b.Fatalf("status = %d", rec.Code)
		}
	}
}
