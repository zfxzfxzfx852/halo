package detect

import (
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"halo/internal/model"
)

// jsonEscapes 还原 Go 默认 JSON 编码转义的 HTML 字符，便于规则匹配。
var jsonEscapes = strings.NewReplacer(
	`\u003c`, "<",
	`\u003e`, ">",
	`\u0026`, "&",
)

// compiledRule 是一条可执行的检测规则。
type compiledRule struct {
	id    string
	level model.RiskLevel
	re    *regexp.Regexp
	// valid 为可选的命中校验；返回 false 表示该候选命中应被丢弃。
	valid func(match string) bool
}

// Engine 用离线规则扫描事件，输出命中的 Finding。
type Engine struct {
	rules []compiledRule
}

// New 构建内置离线规则引擎（密钥 + PII + 注入/越狱）。
func New() *Engine {
	return &Engine{rules: buildRules()}
}

// Inspect 扫描事件并返回去重后的命中列表。
func (e *Engine) Inspect(ev *model.Event) []model.Finding {
	text := eventText(ev)
	var findings []model.Finding
	seen := map[string]bool{}

	for _, r := range e.rules {
		for _, hit := range r.re.FindAllString(text, -1) {
			if r.valid != nil && !r.valid(hit) {
				continue
			}
			key := r.id + "|" + hit
			if seen[key] {
				continue
			}
			seen[key] = true
			findings = append(findings, model.Finding{
				Rule:    r.id,
				Level:   r.level,
				Message: r.id + " detected",
				Field:   "payload",
			})
		}
	}
	return findings
}

// RiskOf 返回命中结果中的最高风险等级；无命中为 none。
func RiskOf(findings []model.Finding) model.RiskLevel {
	high := model.RiskNone
	for _, f := range findings {
		if f.Level.Severity() > high.Severity() {
			high = f.Level
		}
	}
	return high
}

// eventText 汇总事件中所有可扫描的文本载荷。
func eventText(ev *model.Event) string {
	var parts []string
	add := func(b json.RawMessage) {
		if len(b) > 0 {
			parts = append(parts, string(b))
		}
	}
	if ev.Request != nil {
		add(ev.Request.Body)
	}
	if ev.Response != nil {
		add(ev.Response.Body)
	}
	if ev.ToolCall != nil {
		add(ev.ToolCall.Arguments)
	}
	if len(parts) == 0 {
		return ""
	}
	out := ""
	for _, p := range parts {
		out += p + "\n"
	}
	return jsonEscapes.Replace(out)
}

// buildRules 聚合各分类规则并按 id 排序，保证输出稳定。
func buildRules() []compiledRule {
	all := append(append(secretRules(), piiRules()...), injectionRules()...)
	sort.Slice(all, func(i, j int) bool { return all[i].id < all[j].id })
	return all
}
