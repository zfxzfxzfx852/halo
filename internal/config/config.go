package config

import "os"

// Mode 运行模式：observe 仅观察，guard 启用管控。
const (
	ModeObserve = "observe"
	ModeGuard   = "guard"
)

// Config 全局配置。W1 先用环境变量 + 默认值，后续扩展 YAML。
type Config struct {
	Host    string
	Port    string
	DataDir string
	Mode    string
}

// Default 从环境变量加载配置，缺省值保证零配置可启动。
func Default() *Config {
	return &Config{
		Host:    getEnv("VIGIL_HOST", "127.0.0.1"),
		Port:    getEnv("VIGIL_PORT", "8787"),
		DataDir: getEnv("VIGIL_DATA_DIR", "./data"),
		Mode:    getEnv("VIGIL_MODE", ModeObserve),
	}
}

// Addr 返回监听地址。
func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
