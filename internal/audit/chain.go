package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"halo/internal/model"
)

// Chain 维护哈希链状态，生成可校验的审计记录。
type Chain struct {
	lastHash string
	seq      int64
}

// NewChain 从已有链头恢复；空链从 genesis（空哈希、seq=0）开始。
func NewChain(lastSeq int64, lastHash string) *Chain {
	return &Chain{lastHash: lastHash, seq: lastSeq}
}

// Append 对事件计算哈希并生成下一条审计记录。
func (c *Chain) Append(e *model.Event) (model.AuditRecord, error) {
	payload, err := canonicalEvent(e)
	if err != nil {
		return model.AuditRecord{}, err
	}
	sum := sha256.Sum256(append([]byte(c.lastHash), payload...))
	hash := hex.EncodeToString(sum[:])

	c.seq++
	rec := model.AuditRecord{
		Seq:       c.seq,
		EventID:   e.ID,
		Hash:      hash,
		PrevHash:  c.lastHash,
		Timestamp: e.Timestamp,
	}
	c.lastHash = hash
	return rec, nil
}

// Verify 校验审计链与对应事件的完整性，任何篡改或断链即返回错误。
func Verify(records []model.AuditRecord, events map[string]*model.Event) error {
	prev := ""
	for _, rec := range records {
		if rec.PrevHash != prev {
			return fmt.Errorf("seq %d: previous hash mismatch", rec.Seq)
		}
		e, ok := events[rec.EventID]
		if !ok {
			return fmt.Errorf("seq %d: missing event %s", rec.Seq, rec.EventID)
		}
		payload, err := canonicalEvent(e)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(append([]byte(prev), payload...))
		if got := hex.EncodeToString(sum[:]); got != rec.Hash {
			return fmt.Errorf("seq %d: hash mismatch (record may be tampered)", rec.Seq)
		}
		prev = rec.Hash
	}
	return nil
}

func canonicalEvent(e *model.Event) ([]byte, error) {
	return json.Marshal(e)
}
