package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	h := Handler()
	cases := []struct {
		in   string
		want int
	}{
		{"", http.StatusOK},                              // 根目录输出 index.html
		{"assets/index-CADYqMlY.js", http.StatusOK},      // 静态资源
		{"missing/spa/route", http.StatusOK},             // SPA 回退
	}
	for _, c := range cases {
		req := httptest.NewRequest(http.MethodGet, "/"+c.in, nil)
		req.URL.Path = c.in
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != c.want {
			t.Fatalf("in=%q got %d want %d (loc=%q)", c.in, rec.Code, c.want,
				rec.Header().Get("Location"))
		}
	}
}
