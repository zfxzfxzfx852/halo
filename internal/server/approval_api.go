package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"halo/internal/approval"
	"halo/internal/model"
)

func (s *Server) listApprovals(w http.ResponseWriter, _ *http.Request) {
	pending := s.approvals.Pending()
	writeJSON(w, http.StatusOK, map[string]any{"items": pending, "total": len(pending)})
}

func (s *Server) resolveApproval(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Decision string `json:"decision"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	d := model.Decision(body.Decision)
	if d != model.DecisionAllow && d != model.DecisionDeny {
		writeError(w, http.StatusBadRequest, errors.New("decision must be allow or deny"))
		return
	}

	err := s.approvals.Resolve(r.PathValue("id"), d)
	if errors.Is(err, approval.ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
