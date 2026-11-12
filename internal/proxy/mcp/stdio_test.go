package mcp

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"halo/internal/event"
	"halo/internal/model"
)

type fakeConn struct {
	stdin  io.Writer
	stdout io.Reader
	once   sync.Once
	closed chan struct{}
}

func (f *fakeConn) Stdin() io.Writer  { return f.stdin }
func (f *fakeConn) Stdout() io.Reader { return f.stdout }
func (f *fakeConn) Wait() error       { <-f.closed; return nil }
func (f *fakeConn) Close() error {
	f.once.Do(func() { close(f.closed) })
	return nil
}

func TestStdioBridgeToolCall(t *testing.T) {
	clientInR, clientInW := io.Pipe()
	clientOutR, clientOutW := io.Pipe()
	serverStdinR, serverStdinW := io.Pipe()
	serverStdoutR, serverStdoutW := io.Pipe()

	conn := &fakeConn{
		stdin:  serverStdinW,
		stdout: serverStdoutR,
		closed: make(chan struct{}),
	}
	bus := event.New()
	sub := bus.Subscribe()
	defer sub.Close()

	bridge := NewStdioBridge(clientInR, clientOutW, conn, bus, slog.Default())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	runDone := make(chan struct{})
	go func() {
		_ = bridge.Run(ctx)
		close(runDone)
	}()

	reqLine := `{"jsonrpc":"2.0","id":42,"method":"tools/call","params":{"name":"run_cmd","arguments":{"cmd":"ls"}}}`
	// 1. agent 发请求
	go func() {
		_, _ = io.WriteString(clientInW, reqLine+"\n")
	}()

	// 2. MCP server 侧应收到同样的请求
	sc := bufio.NewScanner(serverStdinR)
	if !sc.Scan() || sc.Text() != reqLine {
		t.Fatalf("server did not receive request: %q", sc.Text())
	}

	// 3. server 回响应
	respLine := `{"jsonrpc":"2.0","id":42,"result":{"content":[]}}`
	go func() {
		_, _ = io.WriteString(serverStdoutW, respLine+"\n")
	}()

	// 4. agent 侧应收到响应
	sc2 := bufio.NewScanner(clientOutR)
	if !sc2.Scan() || sc2.Text() != respLine {
		t.Fatalf("client did not receive response: %q", sc2.Text())
	}

	// 5. 总线上应有关联的两个事件
	want := map[model.EventType]bool{model.EventToolCall: false, model.EventToolResult: false}
	var parentID string
	deadline := time.After(2 * time.Second)
	for got := 0; got < 2; {
		select {
		case e := <-sub.Events():
			if _, ok := want[e.Type]; ok {
				want[e.Type] = true
				got++
				if e.Type == model.EventToolCall {
					if e.ToolCall.Name != "run_cmd" {
						t.Fatalf("tool name = %q", e.ToolCall.Name)
					}
					parentID = e.ID
				}
				if e.Type == model.EventToolResult && e.ParentID != parentID {
					t.Fatalf("result parent = %q, want %q", e.ParentID, parentID)
				}
			}
		case <-deadline:
			t.Fatalf("missing events: %+v", want)
		}
	}

	// 关闭管道，让 pump 结束
	_ = clientInW.Close()
	_ = serverStdoutW.Close()
	cancel()
	select {
	case <-runDone:
	case <-time.After(2 * time.Second):
		t.Fatal("bridge did not stop")
	}

	if !strings.Contains(reqLine, "tools/call") {
		t.Fatal("sanity")
	}
}
