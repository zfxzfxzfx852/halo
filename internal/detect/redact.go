package detect

// Redact 把文本中所有命中规则的敏感片段替换为 [REDACTED]。
func (e *Engine) Redact(text string) string {
	out := text
	for _, r := range e.rules {
		out = r.re.ReplaceAllString(out, "[REDACTED]")
	}
	return out
}
