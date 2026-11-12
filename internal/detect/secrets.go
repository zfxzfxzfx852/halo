package detect

import (
	"regexp"

	"halo/internal/model"
)

func must(id string, level model.RiskLevel, pattern string) compiledRule {
	return compiledRule{id: id, level: level, re: regexp.MustCompile(pattern)}
}

func secretRules() []compiledRule {
	return []compiledRule{
		must("secret.aws_access_key", model.RiskCritical, `AKIA[0-9A-Z]{16}`),
		must("secret.openai_key", model.RiskCritical, `sk-[A-Za-z0-9]{20,}`),
		must("secret.github_token", model.RiskCritical, `gh[pousr]_[A-Za-z0-9]{36,}`),
		must("secret.private_key", model.RiskCritical, `-----BEGIN (?:RSA |EC |OPENSSH |DSA )?PRIVATE KEY-----`),
		must("secret.bearer_token", model.RiskHigh, `Bearer [A-Za-z0-9._\-]{20,}`),
		must("secret.generic_api_key", model.RiskHigh, `(?i)(?:api[_-]?key|secret|token)["'=:\s]{1,5}[A-Za-z0-9_\-]{16,}`),
	}
}
