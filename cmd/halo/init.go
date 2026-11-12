package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"halo/internal/config"
)

// runInit 交互式生成 halo.yaml，帮助新用户快速完成配置。
func runInit() error {
	reader := bufio.NewReader(os.Stdin)
	ask := func(prompt, def string) string {
		if def != "" {
			fmt.Printf("%s [%s]: ", prompt, def)
		} else {
			fmt.Printf("%s: ", prompt)
		}
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			return def
		}
		return line
	}

	fmt.Println("Halo 配置向导，按回车使用默认值。")
	port := ask("监听端口", "8787")
	upstream := ask("上游模型/API 地址（可留空，稍后手动填）", "")
	mode := ask("运行模式（observe 仅观察 / guard 启用管控）", config.ModeObserve)
	dataDir := ask("数据目录", "./data")
	policyFile := ask("策略文件路径（guard 模式可留空）", "")

	if mode != config.ModeObserve && mode != config.ModeGuard {
		return fmt.Errorf("invalid mode %q: must be observe or guard", mode)
	}

	var b strings.Builder
	b.WriteString("# 由 `halo init` 生成\n")
	fmt.Fprintf(&b, "host: 127.0.0.1\nport: %s\ndataDir: %s\nmode: %s\nlogLevel: info\n", port, dataDir, mode)
	fmt.Fprintf(&b, "policyFile: %s\n", policyFile)
	b.WriteString("upstream:\n")
	fmt.Fprintf(&b, "  baseUrl: %q\n", upstream)
	b.WriteString("storage:\n  retentionDays: 30\n")

	const target = "halo.yaml"
	if _, err := os.Stat(target); err == nil {
		answer := ask(target+" 已存在，是否覆盖？(y/N)", "n")
		if strings.ToLower(answer) != "y" {
			return errors.New("已取消，未覆盖现有配置")
		}
	}
	if err := os.WriteFile(target, []byte(b.String()), 0o644); err != nil {
		return err
	}
	fmt.Printf("已生成 %s。运行 `halo` 启动，控制台访问 http://127.0.0.1:%s/console/\n", target, port)
	return nil
}
