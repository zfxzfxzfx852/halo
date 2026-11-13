package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"halo/internal/approval"
	"halo/internal/config"
	"halo/internal/policy"
	"halo/internal/store"
	"halo/web"
)

// 构建期通过 -ldflags 注入。
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// Server 持有 HTTP 服务与全部依赖。
type Server struct {
	cfg          *config.Config
	log          *slog.Logger
	http         *http.Server
	store        *store.Store
	hub          *Hub
	approvals    *approval.Manager
	policy       *policy.Engine
	proxyHandler http.Handler
}

// New 构建服务实例。
func New(cfg *config.Config, log *slog.Logger, st *store.Store, hub *Hub, approvals *approval.Manager, proxyHandler http.Handler) *Server {
	mux := http.NewServeMux()
	s := &Server{
		cfg:          cfg,
		log:          log,
		store:        st,
		hub:          hub,
		approvals:    approvals,
		proxyHandler: proxyHandler,
	}

	// 内部 API
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /api/v1/version", s.version)
	mux.HandleFunc("GET /api/v1/events", s.listEvents)
	mux.HandleFunc("GET /api/v1/events/{id}", s.getEvent)
	mux.HandleFunc("GET /api/v1/stats", s.stats)
	mux.HandleFunc("GET /api/v1/approvals", s.listApprovals)
	mux.HandleFunc("POST /api/v1/approvals/{id}", s.resolveApproval)
	mux.HandleFunc("GET /api/v1/audit/verify", s.verifyAudit)
	mux.HandleFunc("GET /api/v1/audit/export", s.exportAudit)
	mux.HandleFunc("GET /api/v1/sessions", s.listSessions)
	mux.HandleFunc("GET /api/v1/findings", s.listFindings)
	mux.HandleFunc("GET /api/v1/policies", s.listPolicies)
	mux.HandleFunc("GET /api/v1/runtime", s.runtimeInfo)
	mux.HandleFunc("GET /api/v1/stream", hub.ServeWS)

	// 嵌入式控制台（iOS 风格 UI）。
	mux.HandleFunc("GET /console", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/console/", http.StatusTemporaryRedirect)
	})
	mux.Handle("GET /console/", http.StripPrefix("/console/", web.Handler()))

	if proxyHandler != nil {
		// 其余路径兜底转发；内部 API 与控制台优先匹配。
		mux.Handle("/", proxyHandler)
	} else {
		// 无上游时，根路径引导到控制台。
		mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/console/", http.StatusTemporaryRedirect)
				return
			}
			http.NotFound(w, r)
		})
	}

	s.http = &http.Server{
		Addr:              cfg.Addr(),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"version":   Version,
		"commit":    Commit,
		"buildTime": BuildTime,
		"mode":      s.cfg.Mode,
	})
}

// Start 启动监听，阻塞直到服务结束。
func (s *Server) Start() error {
	s.log.Info("halo starting",
		"addr", s.http.Addr, "mode", s.cfg.Mode,
		"proxy", s.proxyHandler != nil, "wsClients", s.hub.ClientCount())
	return s.http.ListenAndServe()
}

// Shutdown 优雅关闭。
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
