package approval

import (
	"context"
	"errors"
	"sync"
	"time"

	"halo/internal/event"
	"halo/internal/model"
)

// ErrNotFound 表示审批请求不存在或已被处理。
var ErrNotFound = errors.New("approval request not found")

// Request 是一次等待人工决策的请求。
type Request struct {
	ID        string
	EventID   string
	Summary   string
	Risk      model.RiskLevel
	CreatedAt time.Time
	decision  chan model.Decision
}

// Manager 管理挂起的人工审批请求。
type Manager struct {
	mu      sync.Mutex
	pending map[string]*Request
	bus     event.Bus
	timeout time.Duration
}

// NewManager 创建审批管理器；timeout 为等待上限，超时默认拒绝。
func NewManager(bus event.Bus, timeout time.Duration) *Manager {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Manager{pending: map[string]*Request{}, bus: bus, timeout: timeout}
}

// Approve 实现 proxy.Approver：挂起直到人工决策、超时或请求取消。
func (m *Manager) Approve(ctx context.Context, ev *model.Event) model.Decision {
	req := &Request{
		ID:        model.NewID("apr"),
		EventID:   ev.ID,
		Summary:   ev.Method,
		Risk:      ev.Risk,
		CreatedAt: time.Now(),
		decision:  make(chan model.Decision, 1),
	}

	m.mu.Lock()
	m.pending[req.ID] = req
	m.mu.Unlock()

	if m.bus != nil {
		m.bus.Publish(context.Background(), eventFor(req))
	}

	timer := time.NewTimer(m.timeout)
	defer timer.Stop()

	select {
	case d := <-req.decision:
		return d
	case <-timer.C:
		m.remove(req.ID)
		return model.DecisionDeny
	case <-ctx.Done():
		m.remove(req.ID)
		return model.DecisionDeny
	}
}

// Resolve 由 REST API 调用，对挂起请求给出决策。
func (m *Manager) Resolve(id string, d model.Decision) error {
	m.mu.Lock()
	req, ok := m.pending[id]
	if ok {
		delete(m.pending, id)
	}
	m.mu.Unlock()
	if !ok {
		return ErrNotFound
	}
	select {
	case req.decision <- d:
	default:
	}
	return nil
}

// Pending 返回当前挂起的请求列表。
func (m *Manager) Pending() []*Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Request, 0, len(m.pending))
	for _, r := range m.pending {
		out = append(out, r)
	}
	return out
}

func (m *Manager) remove(id string) {
	m.mu.Lock()
	delete(m.pending, id)
	m.mu.Unlock()
}

func eventFor(req *Request) *model.Event {
	return &model.Event{
		ID:        model.NewID("evt"),
		Type:      model.EventDecision,
		Protocol:  model.ProtocolHTTP,
		Timestamp: time.Now(),
		Method:    req.Summary,
		Risk:      req.Risk,
		Decision: &model.DecisionRecord{
			Decision: model.DecisionAsk,
			Reason:   "awaiting human approval",
			Policy:   req.ID,
		},
	}
}
