package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "8787" {
		t.Errorf("default port = %q, want 8787", cfg.Port)
	}
	if cfg.Mode != ModeObserve {
		t.Errorf("default mode = %q, want observe", cfg.Mode)
	}
	if cfg.Storage.RetentionDays != 30 {
		t.Errorf("retention = %d, want 30", cfg.Storage.RetentionDays)
	}
}

func TestLoadYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "halo.yaml")
	content := []byte("port: 9999\nmode: guard\nstorage:\n  retentionDays: 7\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != "9999" || cfg.Mode != ModeGuard || cfg.Storage.RetentionDays != 7 {
		t.Errorf("yaml not applied: %+v", cfg)
	}
}

func TestEnvOverridesYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "halo.yaml")
	if err := os.WriteFile(path, []byte("port: 9999\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HALO_PORT", "7777")
	t.Setenv("HALO_MODE", "guard")
	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "7777" {
		t.Errorf("env should override, got port %q", cfg.Port)
	}
}

func TestInvalidMode(t *testing.T) {
	t.Setenv("HALO_MODE", "bogus")
	if _, err := Load(""); err == nil {
		t.Fatal("expected validation error for invalid mode")
	}
}

func TestAddr(t *testing.T) {
	cfg, _ := Load("")
	if cfg.Addr() != "127.0.0.1:8787" {
		t.Fatalf("addr = %q", cfg.Addr())
	}
}

func TestLoadFileErrors(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.yaml")); err == nil {
		t.Fatal("missing file should error")
	}
	bad := filepath.Join(t.TempDir(), "bad.yaml")
	if err := os.WriteFile(bad, []byte(":::not yaml:::\n  - ["), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(bad); err == nil {
		t.Fatal("bad yaml should error")
	}
}

func TestAllEnvBranches(t *testing.T) {
	t.Setenv("HALO_HOST", "0.0.0.0")
	t.Setenv("HALO_DATA_DIR", "/tmp/x")
	t.Setenv("HALO_LOG_LEVEL", "debug")
	t.Setenv("HALO_POLICY_FILE", "/tmp/p.yaml")
	t.Setenv("HALO_UPSTREAM_URL", "https://x.example")
	t.Setenv("HALO_RETENTION_DAYS", "12")
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "0.0.0.0" || cfg.DataDir != "/tmp/x" || cfg.LogLevel != "debug" ||
		cfg.PolicyFile != "/tmp/p.yaml" || cfg.Upstream.BaseURL != "https://x.example" ||
		cfg.Storage.RetentionDays != 12 {
		t.Fatalf("env not applied: %+v", cfg)
	}
}

func TestBadRetentionEnv(t *testing.T) {
	t.Setenv("HALO_RETENTION_DAYS", "abc")
	cfg, err := Load("")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Storage.RetentionDays != 30 {
		t.Fatalf("bad retention should keep default, got %d", cfg.Storage.RetentionDays)
	}
}

func TestNegativeRetention(t *testing.T) {
	path := filepath.Join(t.TempDir(), "h.yaml")
	if err := os.WriteFile(path, []byte("storage:\n  retentionDays: -1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path); err == nil {
		t.Fatal("negative retention should fail validation")
	}
}
