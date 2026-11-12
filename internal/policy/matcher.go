package policy

import (
	"regexp"
	"sync"
)

var (
	cacheMu sync.Mutex
	cache   = map[string]*regexp.Regexp{}
)

// compile 把 glob 转成正则并缓存：* 匹配任意字符序列，? 匹配单字符。
func compile(pattern string) *regexp.Regexp {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	if re, ok := cache[pattern]; ok {
		return re
	}

	var b []byte
	b = append(b, '^')
	for _, r := range pattern {
		switch r {
		case '*':
			b = append(b, ".*"...)
		case '?':
			b = append(b, '.')
		default:
			b = append(b, regexp.QuoteMeta(string(r))...)
		}
	}
	b = append(b, '$')

	re := regexp.MustCompile(string(b))
	cache[pattern] = re
	return re
}

// match 用 glob 匹配；空 pattern 视为通配，始终成立。
func match(pattern, value string) bool {
	if pattern == "" {
		return true
	}
	return compile(pattern).MatchString(value)
}

// matchAny 判断 value 是否命中 patterns 中任意一条。
func matchAny(patterns []string, value string) bool {
	for _, p := range patterns {
		if match(p, value) {
			return true
		}
	}
	return false
}
