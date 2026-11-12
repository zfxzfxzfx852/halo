package proxy

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vigil/internal/event"
	"vigil/internal/model"
)

const maxBodySize = 16 << 20 // 16 MiB

// Handler 是透明代理处理器，负责转发并记录事件。
type Handler struct {
	upstream *url.URL
	client   *http.Client
	bus      event.Bus
	log      *slog.Logger
	now      func() time.Time
}

// New 构建代理处理器。
func New(upstream *url.URL, bus event.Bus, log *slog.Logger) *Handler {
	return &Handler{
		upstream: upstream,
		client:   &http.Client{},
		bus:      bus,
		log:      log,
		now:      time.Now,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := h.now()
	sessionID := sessionIDFrom(r)
	eventID := model.NewID("evt")

	reqBody, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, "vigil: cannot read request", http.StatusBadGateway)
		return
	}

	target := *h.upstream
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery

	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "vigil: invalid request", http.StatusBadGateway)
		return
	}
	copyHeaders(outReq.Header, r.Header)
	outReq.Host = r.Host

	h.bus.Publish(context.Background(), buildRequestEvent(eventID, r, reqBody, sessionID, ts))

	resp, err := h.client.Do(outReq)
	if err != nil {
		h.log.Error("upstream error", "err", err, "path", r.URL.Path)
		http.Error(w, "vigil: upstream unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	if isStream(resp.Header) {
		h.proxyStream(w, resp, eventID, sessionID)
		return
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
	if err != nil {
		h.log.Error("read response error", "err", err)
		return
	}
	if _, err := w.Write(body); err != nil {
		h.log.Warn("write response error", "err", err)
	}
	h.bus.Publish(context.Background(), buildResponseEvent(eventID, resp.Header, resp.StatusCode, body, sessionID, h.now()))
}

// proxyStream 转发 SSE 流式响应，边转发边累积，不破坏实时性。
func (h *Handler) proxyStream(w http.ResponseWriter, resp *http.Response, eventID, sessionID string) {
	flusher, _ := w.(http.Flusher)
	var buf bytes.Buffer
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), maxBodySize)

	for scanner.Scan() {
		line := scanner.Bytes()
		if _, err := w.Write(append(line, '\n')); err != nil {
			break
		}
		buf.Write(line)
		buf.WriteByte('\n')
		if flusher != nil {
			flusher.Flush()
		}
	}
	h.bus.Publish(context.Background(), buildResponseEvent(eventID, resp.Header, resp.StatusCode, buf.Bytes(), sessionID, h.now()))
}

func isStream(h http.Header) bool {
	return strings.Contains(strings.ToLower(h.Get("Content-Type")), "text/event-stream")
}

func sessionIDFrom(r *http.Request) string {
	if v := r.Header.Get("X-Vigil-Session"); v != "" {
		return v
	}
	return "default"
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
