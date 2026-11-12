package model

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"
)

// NewID 生成带前缀的随机标识，如 evt_3f9a...
func NewID(prefix string) string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand 失败极罕见；用纳秒时间戳保证退化路径仍具备唯一性
		return prefix + "_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return prefix + "_" + hex.EncodeToString(b)
}
