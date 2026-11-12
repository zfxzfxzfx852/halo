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

	"halo/internal/event"
	"halo/internal/model"
	"halo/internal/proxy/mcp"
)

const maxBodySize = 16 << 20 // 16 MiB

// PolicyEvaluator 由策略引擎实现，用于转发前决策。
type PolicyEvaluator interface {
	Decide(ev *model.Event) model.Decision
}

// Approver 处理需要人工确认的调用，挂起直到得到决策。
type Approver interface {
	Approve(ctx context.Context, ev *model.Event) model.Decision
}

// Handler 是透明代理处理器，负责转发并记录事件。
type Handler struct {
	upstream *url.URL
	client   *http.Client
	bus      event.Bus
	log      *slog.Logger
	policy   PolicyEvaluator
	approver Approver
	now      func() time.Time
}

// New 构建代理处理器。
func New(upstream *url.URL, bus event.Bus, log *slog.Logger, pe PolicyEvaluator, ap Approver) *Handler {
	return &Handler{
		upstream: upstream,
		client:   &http.Client{},
		bus:      bus,
		log:      log,
		policy:   pe,
		approver: ap,
		now:      time.Now,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := h.now()
	sessionID := sessionIDFrom(r)
	eventID := model.NewID("evt")

	reqBody, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		http.Error(w, "halo: cannot read request", http.StatusBadGateway)
		return
	}

	target := *h.upstream
	target.Path = r.URL.Path
	target.RawQuery = r.URL.RawQuery

	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, target.String(), bytes.NewReader(reqBody))
	if err != nil {
		http.Error(w, "halo: invalid request", http.StatusBadGateway)
		return
	}
	copyHeaders(outReq.Header, r.Header)
	outReq.Host = r.Host

	reqEvent := buildRequestEvent(eventID, r, reqBody, sessionID, ts)
	h.bus.Publish(context.Background(), reqEvent)
	for _, e := range mcp.EventsFromRequest(reqBody, sessionID, ts) {
		h.bus.Publish(context.Background(), e)
	}

	// 转发前决策：deny 直接阻断；ask 挂起等待人工审批，未批准则阻断。
	if h.policy != nil {
		switch h.policy.Decide(reqEvent) {
		case model.DecisionDeny:
			h.block(w, eventID)
			return
		case model.DecisionAsk:
			decision := model.DecisionDeny
			if h.approver != nil {
				decision = h.approver.Approve(r.Context(), reqEvent)
			}
			if decision != model.DecisionAllow {
				h.block(w, eventID)
				return
			}
		}
	}

	resp, err := h.client.Do(outReq)
	if err != nil {
		h.log.Error("upstream error", "err", err, "path", r.URL.Path)
		http.Error(w, "halo: upstream unavailable", http.StatusBadGateway)
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

// block 返回 403，请求不会被转发到上游。
func (h *Handler) block(w http.ResponseWriter, eventID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte(`{"error":"blocked_by_halo_policy","eventId":"` + eventID + `"}`))
}

func sessionIDFrom(r *http.Request) string {
	if v := r.Header.Get("X-Halo-Session"); v != "" {
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
