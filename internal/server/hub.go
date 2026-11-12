package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"

	"halo/internal/model"
)

// Hub 管理 WebSocket 客户端并广播事件。
type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]context.CancelFunc
	log     *slog.Logger
}

// NewHub 创建 Hub。
func NewHub(log *slog.Logger) *Hub {
	return &Hub{clients: make(map[*websocket.Conn]context.CancelFunc), log: log}
}

// Broadcast 向所有连接客户端推送事件，失败则剔除。
func (h *Hub) Broadcast(e *model.Event) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn, cancel := range h.clients {
		writeCtx, cancelWrite := context.WithTimeout(context.Background(), 5*time.Second)
		err := conn.Write(writeCtx, websocket.MessageText, data)
		cancelWrite()
		if err != nil {
			h.log.Warn("ws write failed, removing client", "err", err)
			cancel()
			_ = conn.Close(websocket.StatusInternalError, "write failed")
			delete(h.clients, conn)
		}
	}
}

// ServeWS 处理 WebSocket 升级与连接生命周期。
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// 本地工具默认允许同源与本地来源；后续可按配置收紧。
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	// 不能用 r.Context()：handler 返回后它会被取消，进而误关长连接。
	ctx, cancel := context.WithCancel(context.Background())
	h.mu.Lock()
	h.clients[conn] = cancel
	h.mu.Unlock()

	go func() {
		defer h.remove(conn)
		for {
			if _, _, err := conn.Reader(ctx); err != nil {
				return
			}
		}
	}()
}

func (h *Hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if cancel, ok := h.clients[conn]; ok {
		cancel()
		delete(h.clients, conn)
	}
	_ = conn.Close(websocket.StatusNormalClosure, "closed")
}

// ClientCount 返回当前连接数。
func (h *Hub) ClientCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}
