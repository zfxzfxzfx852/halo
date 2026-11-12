package detect

import "halo/internal/model"

func piiRules() []compiledRule {
	return []compiledRule{
		must("pii.email", model.RiskMedium, `[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`),
		must("pii.phone_cn", model.RiskHigh, `\b1[3-9]\d{9}\b`),
		must("pii.id_card_cn", model.RiskHigh, `\b\d{17}[\dXx]\b`),
		must("pii.credit_card", model.RiskHigh, `\b(?:\d{4}[ -]){3}\d{4}\b`),
		must("pii.ipv4", model.RiskLow, `\b\d{1,3}(?:\.\d{1,3}){3}\b`),
	}
}
