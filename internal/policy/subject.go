package policy

import (
	"encoding/json"

	"halo/internal/model"
)

// subject 是从事件归一化出的策略评估对象。
type subject struct {
	tool    string
	method  string
	command string
	paths   []string
	args    map[string]any
	risk    model.RiskLevel
}

func extract(e *model.Event) subject {
	s := subject{
		method: e.Method,
		risk:   e.Risk,
		args:   map[string]any{},
	}
	if e.ToolCall != nil {
		s.tool = e.ToolCall.Name
		if len(e.ToolCall.Arguments) > 0 {
			_ = json.Unmarshal(e.ToolCall.Arguments, &s.args)
		}
	}
	s.command, s.paths = pullCommandAndPaths(s.args)
	return s
}

// pullCommandAndPaths 从参数中识别命令与文件路径类字段。
func pullCommandAndPaths(args map[string]any) (string, []string) {
	var command string
	var paths []string
	for k, v := range args {
		val, ok := v.(string)
		if !ok {
			continue
		}
		switch k {
		case "command", "cmd":
			command = val
		case "path", "file", "filename", "dir", "directory":
			paths = append(paths, val)
		}
	}
	return command, paths
}
