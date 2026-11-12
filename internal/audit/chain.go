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

// VerificationReport 描述一次审计校验的结果。
type VerificationReport struct {
	Valid         bool   `json:"valid"`
	Checked       int    `json:"checked"`       // 内容成功重算校验的记录数
	MissingEvents int    `json:"missingEvents"` // 对应事件已被保留策略清理的记录数
	BrokenSeq     int64  `json:"brokenSeq,omitempty"`
	Reason        string `json:"reason,omitempty"`
}

// Verify 校验审计链与对应事件的完整性。断链或现存事件被篡改即判 invalid；
// 对应事件被保留策略清理时，链的连续性仍可校验，仅计入 MissingEvents。
func Verify(records []model.AuditRecord, events map[string]*model.Event) VerificationReport {
	report := VerificationReport{Valid: true}
	prev := ""
	for _, rec := range records {
		if rec.PrevHash != prev {
			report.Valid = false
			report.BrokenSeq = rec.Seq
			report.Reason = fmt.Sprintf("seq %d: previous hash mismatch", rec.Seq)
			return report
		}
		e, ok := events[rec.EventID]
		if !ok {
			// 事件已被保留策略清理：哈希链仍连续，仅无法重算该条内容。
			report.MissingEvents++
			prev = rec.Hash
			continue
		}
		payload, err := canonicalEvent(e)
		if err != nil {
			report.Valid = false
			report.BrokenSeq = rec.Seq
			report.Reason = err.Error()
			return report
		}
		sum := sha256.Sum256(append([]byte(prev), payload...))
		if got := hex.EncodeToString(sum[:]); got != rec.Hash {
			report.Valid = false
			report.BrokenSeq = rec.Seq
			report.Reason = fmt.Sprintf("seq %d: hash mismatch (record may be tampered)", rec.Seq)
			return report
		}
		report.Checked++
		prev = rec.Hash
	}
	return report
}

func canonicalEvent(e *model.Event) ([]byte, error) {
	return json.Marshal(e)
}
