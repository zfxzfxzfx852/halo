package model

// Policy 是声明式策略集合，W2 定义结构，W7 实现决策引擎。
type Policy struct {
	Version string       `json:"version" yaml:"version"`
	Default Decision     `json:"default" yaml:"default"`
	Rules   []PolicyRule `json:"rules" yaml:"rules"`
}

// PolicyRule 描述单条匹配规则；空字段表示通配。
type PolicyRule struct {
	Name     string            `json:"name" yaml:"name"`
	Tool     string            `json:"tool,omitempty" yaml:"tool,omitempty"`
	Method   string            `json:"method,omitempty" yaml:"method,omitempty"`
	Command  string            `json:"command,omitempty" yaml:"command,omitempty"`
	Paths    []string          `json:"paths,omitempty" yaml:"paths,omitempty"`
	Args     map[string]string `json:"args,omitempty" yaml:"args,omitempty"`
	Decision Decision          `json:"decision" yaml:"decision"`
	MaxRisk  RiskLevel         `json:"maxRisk,omitempty" yaml:"maxRisk,omitempty"`
}
