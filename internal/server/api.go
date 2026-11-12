package server

import (
	"errors"
	"net/http"
	"strconv"

	"halo/internal/model"
	"halo/internal/store"
)

func (s *Server) listEvents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := store.Filter{
		Type:      q.Get("type"),
		Protocol:  q.Get("protocol"),
		SessionID: q.Get("session"),
		RiskMin:   model.RiskLevel(q.Get("riskMin")),
		Limit:  clampLimit(atoiDefault(q.Get("limit"), 50)),
		Offset: nonNegative(atoiDefault(q.Get("offset"), 0), 0),
	}
	events, total, err := s.store.ListEvents(r.Context(), f)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": events, "total": total})
}

func (s *Server) getEvent(w http.ResponseWriter, r *http.Request) {
	e, err := s.store.GetEvent(r.Context(), r.PathValue("id"))
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, e)
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	st, err := s.store.Stats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, st)
}

func atoiDefault(v string, def int) int {
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// maxPageSize 限制单次查询返回条数，避免一次拉取过多。
const maxPageSize = 1000

// clampLimit 把分页大小钳制到 [1, maxPageSize]。
func clampLimit(v int) int {
	if v < 1 {
		return 50
	}
	if v > maxPageSize {
		return maxPageSize
	}
	return v
}

// nonNegative 把负值回退到 def，避免非法分页参数导致 SQL 错误。
func nonNegative(v, def int) int {
	if v < 0 {
		return def
	}
	return v
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
