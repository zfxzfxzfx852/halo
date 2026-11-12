package policy

import (
	"fmt"
	"sync"

	"halo/internal/model"
)

// Engine 持有当前生效策略，支持并发评估与原子热更新。
type Engine struct {
	mu     sync.RWMutex
	policy *model.Policy
}

func New(p *model.Policy) *Engine {
	return &Engine{policy: p}
}

// Replace 原子替换策略，实现热更新。
func (e *Engine) Replace(p *model.Policy) {
	e.mu.Lock()
	e.policy = p
	e.mu.Unlock()
}

// Result 是一次策略评估的结论。
type Result struct {
	Decision model.Decision
	RuleName string
	Reason   string
}

// Decide 仅返回决策，供代理层接口适配。
func (e *Engine) Decide(ev *model.Event) model.Decision {
	return e.Evaluate(ev).Decision
}

// Evaluate 按规则顺序匹配，第一个命中的规则生效；否则回退默认策略。
func (e *Engine) Evaluate(ev *model.Event) Result {
	e.mu.RLock()
	p := e.policy
	e.mu.RUnlock()

	if p == nil {
		// 策略缺失时 fail-closed，避免在无策略情况下放行。
		return Result{Decision: model.DecisionDeny, Reason: "no policy loaded; fail closed"}
	}

	s := extract(ev)
	for _, r := range p.Rules {
		// 风险超过规则上限时该规则不适用。
		if r.MaxRisk != "" && s.risk.Severity() > r.MaxRisk.Severity() {
			continue
		}
		if ruleMatches(r, s) {
			return Result{
				Decision: r.Decision,
				RuleName: r.Name,
				Reason:   fmt.Sprintf("rule %q matched", r.Name),
			}
		}
	}
	return Result{
		Decision: p.Default,
		Reason:   "no rule matched; using default",
	}
}

// ruleMatches 对规则的所有非空条件做逻辑与（AND）。
func ruleMatches(r model.PolicyRule, s subject) bool {
	if r.Tool != "" && !match(r.Tool, s.tool) {
		return false
	}
	if r.Method != "" && !match(r.Method, s.method) {
		return false
	}
	if r.Command != "" && !match(r.Command, s.command) {
		return false
	}
	if len(r.Paths) > 0 {
		hit := false
		for _, p := range s.paths {
			if matchAny(r.Paths, p) {
				hit = true
				break
			}
		}
		if !hit {
			return false
		}
	}
	for key, pattern := range r.Args {
		raw, ok := s.args[key]
		if !ok {
			return false
		}
		if !match(pattern, fmt.Sprintf("%v", raw)) {
			return false
		}
	}
	return true
}
