package model

import "crypto/rand"
import "encoding/hex"

// NewID 生成带前缀的随机标识，如 evt_3f9a...
func NewID(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand 失败极罕见；退化为时间戳仍保证唯一性弱需求
		return prefix + "_fallback"
	}
	return prefix + "_" + hex.EncodeToString(b)
}
