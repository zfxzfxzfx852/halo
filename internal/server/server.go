package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"halo/internal/config"
	"halo/internal/store"
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
	proxyHandler http.Handler
}

// New 构建服务实例。
func New(cfg *config.Config, log *slog.Logger, st *store.Store, hub *Hub, proxyHandler http.Handler) *Server {
	mux := http.NewServeMux()
	s := &Server{
		cfg:          cfg,
		log:          log,
		store:        st,
		hub:          hub,
		proxyHandler: proxyHandler,
	}

	// 内部 API
	mux.HandleFunc("GET /healthz", s.healthz)
	mux.HandleFunc("GET /api/v1/version", s.version)
	mux.HandleFunc("GET /api/v1/events", s.listEvents)
	mux.HandleFunc("GET /api/v1/events/{id}", s.getEvent)
	mux.HandleFunc("GET /api/v1/stats", s.stats)
	mux.HandleFunc("GET /api/v1/stream", hub.ServeWS)

	if proxyHandler != nil {
		// 兜底转发；更具体的内部路由优先匹配。
		mux.Handle("/", proxyHandler)
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
