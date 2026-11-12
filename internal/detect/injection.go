package detect

import "halo/internal/model"

func injectionRules() []compiledRule {
	return []compiledRule{
		// 提示注入 / 越狱
		must("injection.ignore_instructions", model.RiskHigh, `(?i)ignore (?:all |previous |prior )*instructions`),
		must("injection.disregard", model.RiskHigh, `(?i)disregard (?:all|previous|prior|the) (?:instructions|rules|prompts)`),
		must("injection.you_are_now", model.RiskHigh, `(?i)you are now [a-z]`),
		must("injection.do_anything_now", model.RiskHigh, `(?i)do anything now`),
		must("injection.role_tag", model.RiskHigh, `</?(?:system|user|assistant)>|\[(?:SYSTEM|INST|DAN)\]`),

		// 命令注入
		must("injection.shell_command", model.RiskCritical, `[;&|]\s*(?:rm|cat|curl|wget|bash|sh|nc)\b`),
		must("injection.command_substitution", model.RiskCritical, `\$\([^)]*\)|`+"`[^`]*`"),

		// SQL 注入
		must("injection.sql_tautology", model.RiskHigh, `(?i)['"]\s*OR\s+['"]?\d*['"]?\s*=\s*['"]?\d*`),
		must("injection.sql_union", model.RiskHigh, `(?i)UNION\s+(?:ALL\s+)?SELECT`),

		// 路径穿越与敏感文件
		must("injection.path_traversal", model.RiskHigh, `\.\./`),
		must("injection.sensitive_file", model.RiskHigh, `/etc/(?:passwd|shadow)`),
	}
}
