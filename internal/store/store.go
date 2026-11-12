package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite" // 纯 Go SQLite 驱动

	"halo/internal/model"
)

// ErrNotFound 表示事件不存在。
var ErrNotFound = errors.New("event not found")

// Filter 定义事件查询条件。
type Filter struct {
	Type      string
	Protocol  string
	SessionID string
	RiskMin   model.RiskLevel
	Limit     int
	Offset    int
}

// Stats 聚合统计。
type Stats struct {
	Total  int64
	ByType map[string]int64
	ByRisk map[string]int64
}

// Store 是 SQLite 持久化层。
type Store struct {
	db *sql.DB
}

// Open 打开数据库并执行迁移。
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// SQLite 单写者，限制连接数避免 database locked。
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS events (
	id TEXT PRIMARY KEY,
	type TEXT NOT NULL,
	protocol TEXT NOT NULL,
	session_id TEXT NOT NULL,
	method TEXT,
	parent_id TEXT,
	risk TEXT NOT NULL,
	agent TEXT,
	payload TEXT NOT NULL,
	created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_events_created ON events(created_at);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(type);
CREATE INDEX IF NOT EXISTS idx_events_session ON events(session_id);
CREATE INDEX IF NOT EXISTS idx_events_risk ON events(risk);
	`)
	return err
}

// InsertEvent 写入或替换一条事件。
func (s *Store) InsertEvent(ctx context.Context, e *model.Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
INSERT OR REPLACE INTO events
(id, type, protocol, session_id, method, parent_id, risk, agent, payload, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		e.ID, string(e.Type), string(e.Protocol), e.SessionID, e.Method,
		e.ParentID, string(e.Risk), e.Agent, string(payload), e.Timestamp.UnixNano())
	return err
}

// ListEvents 按条件分页返回事件及总数。
func (s *Store) ListEvents(ctx context.Context, f Filter) ([]*model.Event, int, error) {
	where, args := buildWhere(f)

	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events"+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limit := f.Limit
	if limit <= 0 {
		limit = 50
	}
	query := "SELECT payload FROM events" + where +
		" ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, f.Offset)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var events []*model.Event
	for rows.Next() {
		var payload string
		if err := rows.Scan(&payload); err != nil {
			return nil, 0, err
		}
		var e model.Event
		if err := json.Unmarshal([]byte(payload), &e); err != nil {
			return nil, 0, err
		}
		events = append(events, &e)
	}
	return events, total, rows.Err()
}

// GetEvent 按 ID 返回单条事件。
func (s *Store) GetEvent(ctx context.Context, id string) (*model.Event, error) {
	var payload string
	err := s.db.QueryRowContext(ctx, "SELECT payload FROM events WHERE id = ?", id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	var e model.Event
	if err := json.Unmarshal([]byte(payload), &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Stats 返回总量与分类聚合。
func (s *Store) Stats(ctx context.Context) (Stats, error) {
	st := Stats{ByType: map[string]int64{}, ByRisk: map[string]int64{}}
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM events").Scan(&st.Total); err != nil {
		return st, err
	}
	if err := s.groupCount(ctx, "type", st.ByType); err != nil {
		return st, err
	}
	if err := s.groupCount(ctx, "risk", st.ByRisk); err != nil {
		return st, err
	}
	return st, nil
}

func (s *Store) groupCount(ctx context.Context, col string, dst map[string]int64) error {
	rows, err := s.db.QueryContext(ctx, "SELECT "+col+", COUNT(*) FROM events GROUP BY "+col)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var k string
		var n int64
		if err := rows.Scan(&k, &n); err != nil {
			return err
		}
		dst[k] = n
	}
	return rows.Err()
}

// PurgeBefore 删除指定时间之前的事件，返回删除条数。
func (s *Store) PurgeBefore(ctx context.Context, before time.Time) (int64, error) {
	res, err := s.db.ExecContext(ctx, "DELETE FROM events WHERE created_at < ?", before.UnixNano())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Close 关闭数据库。
func (s *Store) Close() error { return s.db.Close() }

func buildWhere(f Filter) (string, []any) {
	var conds []string
	var args []any
	if f.Type != "" {
		conds = append(conds, "type = ?")
		args = append(args, f.Type)
	}
	if f.Protocol != "" {
		conds = append(conds, "protocol = ?")
		args = append(args, f.Protocol)
	}
	if f.SessionID != "" {
		conds = append(conds, "session_id = ?")
		args = append(args, f.SessionID)
	}
	if f.RiskMin != "" && f.RiskMin != model.RiskNone {
		conds = append(conds, "risk IN (?, ?, ?, ?)")
		args = append(args,
			string(model.RiskLow), string(model.RiskMedium),
			string(model.RiskHigh), string(model.RiskCritical))
	}
	if len(conds) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}
