package audit

import (
	"bytes"
	"encoding/json"

	"halo/internal/model"
)

// ExportJSONL 把审计记录序列化为 JSONL（每行一条）。
func ExportJSONL(records []model.AuditRecord) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, r := range records {
		if err := enc.Encode(r); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
