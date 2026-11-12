package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"halo/internal/approval"
	"halo/internal/config"
	"halo/internal/detect"
	"halo/internal/event"
	mcpproxy "halo/internal/proxy/mcp"
	"halo/internal/server"
	"halo/internal/store"
)

// runMCP 以 stdio 桥接方式包装一个 MCP server，同时提供 HTTP 控制台。
func runMCP(args []string) error {
	fs := flag.NewFlagSet("mcp", flag.ExitOnError)
	configPath := fs.String("config", envOr("HALO_CONFIG", "halo.yaml"), "path to config file")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return errors.New("missing MCP server command: halo mcp -- <command> [args...]")
	}

	path := *configPath
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = ""
	}
	cfg, err := config.Load(path)
	if err != nil {
		return err
	}
	// stdio 承载 JSON-RPC，日志必须写到 stderr，避免污染协议。
	log := newStderrLogger(cfg.LogLevel)

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return err
	}
	st, err := store.Open(filepath.Join(cfg.DataDir, "halo.db"))
	if err != nil {
		return err
	}
	defer st.Close()

	bus := event.New()
	hub := server.NewHub(log)
	runCollector(bus, st, hub, detect.New(), log)
	approvals := approval.NewManager(bus, 0)

	srv := server.New(cfg, log, st, hub, approvals, nil)
	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("console server error", "err", err)
		}
	}()

	conn, err := mcpproxy.DialProcess(cmdArgs[0], cmdArgs[1:]...)
	if err != nil {
		return err
	}
	bridge := mcpproxy.NewStdioBridge(os.Stdin, os.Stdout, conn, bus, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Info("mcp bridge ready", "server", cmdArgs[0], "console", cfg.Addr())
	return bridge.Run(ctx)
}

func newStderrLogger(level string) *slog.Logger {
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
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}
