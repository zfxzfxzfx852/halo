package detect

import (
	"regexp"

	"halo/internal/model"
)

func piiRules() []compiledRule {
	return []compiledRule{
		must("pii.email", model.RiskMedium, `[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`),
		must("pii.phone_cn", model.RiskHigh, `\b1[3-9]\d{9}\b`),
		{
			id:    "pii.id_card_cn",
			level: model.RiskHigh,
			re:    regexp.MustCompile(`\b\d{17}[\dXx]\b`),
			valid: validIDCard,
		},
		must("pii.credit_card", model.RiskHigh, `\b(?:\d{4}[ -]){3}\d{4}\b`),
		must("pii.ipv4", model.RiskLow, `\b(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\b`),
	}
}

// idWeights 为身份证前 17 位的加权因子（GB 11643-1999）。
var idWeights = [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// idCheckCodes 为加权和对 11 取模后对应的校验码。
var idCheckCodes = [11]byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// validIDCard 校验 18 位居民身份证号的校验位，过滤随机长数字的误报。
func validIDCard(s string) bool {
	if len(s) != 18 {
		return false
	}
	sum := 0
	for i := 0; i < 17; i++ {
		d := int(s[i] - '0')
		if d < 0 || d > 9 {
			return false
		}
		sum += d * idWeights[i]
	}
	last := s[17]
	if last == 'x' {
		last = 'X'
	}
	return idCheckCodes[sum%11] == last
}
