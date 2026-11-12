package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"halo/internal/model"
)

// Parse 解析并校验 YAML 策略。
func Parse(data []byte) (*model.Policy, error) {
	var p model.Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parse policy: %w", err)
	}
	if err := Validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

// LoadFile 从文件加载策略。
func LoadFile(path string) (*model.Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read policy %s: %w", path, err)
	}
	return Parse(data)
}

// Validate 校验默认策略与每条规则的决策值合法。
func Validate(p *model.Policy) error {
	if p.Default == "" {
		return fmt.Errorf("policy.default is required")
	}
	if !validDecision(p.Default) {
		return fmt.Errorf("invalid default decision %q", p.Default)
	}
	for i, r := range p.Rules {
		if !validDecision(r.Decision) {
			return fmt.Errorf("rule #%d %q has invalid decision %q", i+1, r.Name, r.Decision)
		}
	}
	return nil
}

func validDecision(d model.Decision) bool {
	return d == model.DecisionAllow || d == model.DecisionDeny || d == model.DecisionAsk
}
