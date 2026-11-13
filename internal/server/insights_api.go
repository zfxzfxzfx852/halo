package server

import (
	"net/http"

	"halo/internal/policy"
)

// SetPolicy 绑定策略引擎，使控制台可以展示当前生效规则。
func (s *Server) SetPolicy(engine *policy.Engine) {
	s.policy = engine
}

// listSessions 返回最近的会话聚合。
func (s *Server) listSessions(w http.ResponseWriter, r *http.Request) {
	limit := clampLimit(atoiDefault(r.URL.Query().Get("limit"), 50))
	sessions, err := s.store.Sessions(r.Context(), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if len(sessions) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "total": 0})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": sessions, "total": len(sessions)})
}

// listFindings 返回检测规则命中聚合。
func (s *Server) listFindings(w http.ResponseWriter, r *http.Request) {
	findings, err := s.store.Findings(r.Context(), 0)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": findings, "total": len(findings)})
}

// listPolicies 返回当前生效策略规则。
func (s *Server) listPolicies(w http.ResponseWriter, _ *http.Request) {
	if s.policy == nil {
		writeJSON(w, http.StatusOK, map[string]any{"default": "", "rules": []any{}})
		return
	}
	def, rules := s.policy.Snapshot()
	writeJSON(w, http.StatusOK, map[string]any{"default": def, "rules": rules})
}

// runtimeInfo 返回控制台需要的运行信息（不含敏感配置）。
func (s *Server) runtimeInfo(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"version":       Version,
		"commit":        Commit,
		"buildTime":     BuildTime,
		"mode":          s.cfg.Mode,
		"addr":          s.cfg.Addr(),
		"dataDir":       s.cfg.DataDir,
		"retentionDays": s.cfg.Storage.RetentionDays,
		"policyFile":    s.cfg.PolicyFile,
		"upstream":      s.cfg.Upstream.BaseURL,
		"wsClients":     s.hub.ClientCount(),
	})
}
