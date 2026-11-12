package server

import (
	"net/http"

	"halo/internal/audit"
	"halo/internal/model"
	"halo/internal/store"
)

func (s *Server) verifyAudit(w http.ResponseWriter, r *http.Request) {
	records, err := s.store.ListAuditRecords(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	events, _, err := s.store.ListEvents(r.Context(), store.Filter{Limit: 100000})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	em := make(map[string]*model.Event, len(events))
	for _, e := range events {
		em[e.ID] = e
	}

	if err := audit.Verify(records, em); err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"valid": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"valid": true, "records": len(records)})
}

func (s *Server) exportAudit(w http.ResponseWriter, r *http.Request) {
	records, err := s.store.ListAuditRecords(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	data, err := audit.ExportJSONL(records)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Content-Disposition", `attachment; filename="halo-audit.jsonl"`)
	_, _ = w.Write(data)
}
