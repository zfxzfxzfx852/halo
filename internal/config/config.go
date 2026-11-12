package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// 运行模式。
const (
	ModeObserve = "observe"
	ModeGuard   = "guard"
)

// UpstreamConfig 指向被代理的模型/API 地址。
type UpstreamConfig struct {
	BaseURL string `yaml:"baseUrl"`
}

// StorageConfig 存储相关配置。
type StorageConfig struct {
	RetentionDays int `yaml:"retentionDays"`
}

// Config 全局配置。
type Config struct {
	Host       string         `yaml:"host"`
	Port       string         `yaml:"port"`
	DataDir    string         `yaml:"dataDir"`
	Mode       string         `yaml:"mode"`
	LogLevel   string         `yaml:"logLevel"`
	PolicyFile string         `yaml:"policyFile"`
	Upstream   UpstreamConfig `yaml:"upstream"`
	Storage    StorageConfig  `yaml:"storage"`
}

func defaults() *Config {
	return &Config{
		Host:     "127.0.0.1",
		Port:     "8787",
		DataDir:  "./data",
		Mode:     ModeObserve,
		LogLevel: "info",
		Upstream: UpstreamConfig{BaseURL: ""},
		Storage:  StorageConfig{RetentionDays: 30},
	}
}

// Load 按「默认值 → YAML 文件 → 环境变量」顺序加载并校验。
// path 为空时跳过文件加载。
func Load(path string) (*Config, error) {
	cfg := defaults()
	if path != "" {
		if err := loadFile(cfg, path); err != nil {
			return nil, err
		}
	}
	if err := applyEnv(cfg); err != nil {
		return nil, err
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}
func loadFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}

func applyEnv(cfg *Config) error {
	if v := os.Getenv("HALO_HOST"); v != "" {
		cfg.Host = v
	}
	if v := os.Getenv("HALO_PORT"); v != "" {
		cfg.Port = v
	}
	if v := os.Getenv("HALO_DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("HALO_MODE"); v != "" {
		cfg.Mode = v
	}
	if v := os.Getenv("HALO_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("HALO_POLICY_FILE"); v != "" {
		cfg.PolicyFile = v
	}
	if v := os.Getenv("HALO_UPSTREAM_URL"); v != "" {
		cfg.Upstream.BaseURL = v
	}
	if v := os.Getenv("HALO_RETENTION_DAYS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid HALO_RETENTION_DAYS %q: %w", v, err)
		}
		cfg.Storage.RetentionDays = n
	}
	return nil
}

func (c *Config) validate() error {
	if c.Mode != ModeObserve && c.Mode != ModeGuard {
		return fmt.Errorf("invalid mode %q: must be observe or guard", c.Mode)
	}
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %q: must be an integer in 1-65535", c.Port)
	}
	if c.Storage.RetentionDays < 0 {
		return fmt.Errorf("storage.retentionDays must be >= 0")
	}
	return nil
}

// Addr 返回监听地址。
func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}
