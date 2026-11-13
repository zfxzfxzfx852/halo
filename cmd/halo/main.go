package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"halo/internal/approval"
	"halo/internal/audit"
	"halo/internal/config"
	"halo/internal/detect"
	"halo/internal/event"
	"halo/internal/model"
	"halo/internal/policy"
	"halo/internal/proxy"
	"halo/internal/sandbox"
	"halo/internal/server"
	"halo/internal/store"
)

// version 在发布构建时通过 -ldflags -X main.version=... 注入。
var version = "dev"

func main() {
	// 子命令：halo init 进入交互式配置向导。
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := runInit(); err != nil {
			slog.Error("init failed", "err", err)
			os.Exit(1)
		}
		return
	}

	// 子命令：halo mcp 以 stdio 方式桥接 MCP server。
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		if err := runMCP(os.Args[2:]); err != nil {
			slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("mcp failed", "err", err)
			os.Exit(1)
		}
		return
	}

	configPath := flag.String("config", envOr("HALO_CONFIG", "halo.yaml"), "path to config file")
	flag.Parse()

	// 配置文件缺失时回退到默认值与环境变量，保证零配置可启动。
	path := *configPath
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = ""
	}
	cfg, err := config.Load(path)
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}
	log := newLogger(cfg.LogLevel)
	log.Info("halo starting", "version", version, "mode", cfg.Mode, "addr", cfg.Addr())

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		log.Error("create data dir failed", "err", err)
		os.Exit(1)
	}

	// 存储
	st, err := store.Open(filepath.Join(cfg.DataDir, "halo.db"))
	if err != nil {
		log.Error("open store failed", "err", err)
		os.Exit(1)
	}
	defer func() { _ = st.Close() }()

	// 事件总线 + WebSocket Hub
	bus := event.New()
	hub := server.NewHub(log)
	detector := detect.New()
	runCollector(bus, st, hub, detector, log)

	// 代理
	// 加载策略引擎（可选）。
	var engine *policy.Engine
	if cfg.PolicyFile != "" {
		pol, err := policy.LoadFile(cfg.PolicyFile)
		if err != nil {
			log.Error("load policy failed", "err", err)
			os.Exit(1)
		}
		engine = policy.New(pol)
		log.Info("policy loaded", "file", cfg.PolicyFile, "rules", len(pol.Rules))
	}

	// 审批管理器与沙箱注册表（沙箱能力就绪，供隔离执行与后续驱动扩展）。
	approvals := approval.NewManager(bus, 0)
	sandboxes := sandbox.NewRegistry()
	log.Info("sandbox ready", "drivers", sandboxDriverNames(sandboxes))

	var proxyHandler http.Handler
	if cfg.Upstream.BaseURL != "" {
		upstream, err := url.Parse(cfg.Upstream.BaseURL)
		if err != nil {
			log.Error("invalid upstream url", "err", err)
			os.Exit(1)
		}
		// 仅 guard 模式启用阻断与审批；observe 模式保持旁路观察。
		var evaluator proxy.PolicyEvaluator
		var approver proxy.Approver
		if cfg.Mode == config.ModeGuard {
			evaluator = engine
			approver = approvals
		}
		proxyHandler = proxy.New(upstream, bus, log, evaluator, approver)
	}

	srv := server.New(cfg, log, st, hub, approvals, proxyHandler)
	srv.SetPolicy(engine)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 数据保留：按配置周期清理过期事件（审计记录独立保留）。
	startRetentionJanitor(ctx, st, cfg.Storage.RetentionDays, log)

	select {
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		log.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("shutdown error", "err", err)
			os.Exit(1)
		}
	}
	log.Info("halo stopped")
}

// runCollector 订阅总线，检测富化后持久化、写入审计链并广播。
func runCollector(bus event.Bus, st *store.Store, hub *server.Hub, det *detect.Engine, log *slog.Logger) {
	// 从存储恢复哈希链头，重启后继续延伸而非重新开始。
	chain := audit.NewChain(0, "")
	if last, ok, err := st.LastAuditRecord(context.Background()); err != nil {
		log.Error("restore audit chain failed", "err", err)
	} else if ok {
		chain = audit.NewChain(last.Seq, last.Hash)
	}

	sub := bus.Subscribe()
	go func() {
		defer sub.Close()
		for e := range sub.Events() {
			enrich(e, det)
			if err := st.InsertEvent(context.Background(), e); err != nil {
				log.Error("persist event failed", "err", err)
			}
			if rec, err := chain.Append(e); err != nil {
				log.Error("audit append failed", "err", err)
			} else if err := st.InsertAuditRecord(context.Background(), rec); err != nil {
				log.Error("persist audit record failed", "err", err)
			}
			hub.Broadcast(e)
		}
	}()
}

// enrich 运行离线检测，附加命中规则并把事件风险提升到检测到的最高等级。
func enrich(e *model.Event, det *detect.Engine) {
	findings := det.Inspect(e)
	if len(findings) == 0 {
		return
	}
	e.Findings = append(e.Findings, findings...)
	if r := detect.RiskOf(findings); r.Severity() > e.Risk.Severity() {
		e.Risk = r
	}
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func sandboxDriverNames(r *sandbox.Registry) []string {
	return r.Names()
}

// startRetentionJanitor 启动数据保留清理循环；retentionDays<=0 表示永久保留。
func startRetentionJanitor(ctx context.Context, st *store.Store, retentionDays int, log *slog.Logger) {
	if retentionDays <= 0 {
		return
	}
	ticker := time.NewTicker(24 * time.Hour)
	purge := func() {
		cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
		n, err := st.PurgeBefore(ctx, cutoff)
		if err != nil {
			if ctx.Err() == nil {
				log.Error("retention purge failed", "err", err)
			}
			return
		}
		if n > 0 {
			log.Info("retention purge complete", "removed", n)
		}
	}
	go func() {
		defer ticker.Stop()
		purge() // 启动时先清理一次
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				purge()
			}
		}
	}()
}
