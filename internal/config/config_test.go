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
