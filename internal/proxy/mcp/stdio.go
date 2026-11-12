package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"time"

	"halo/internal/event"
	"halo/internal/model"
)

const maxMessageSize = 16 << 20

// Direction 表示消息流向。
type Direction string

const (
	DirInbound  Direction = "inbound"  // agent -> MCP server
	DirOutbound Direction = "outbound" // MCP server -> agent
)

// ServerConn 抽象 MCP server 的双向通道，便于测试与替换驱动。
type ServerConn interface {
	Stdin() io.Writer
	Stdout() io.Reader
	Wait() error
	Close() error
}

type processConn struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// DialProcess 启动一个 MCP server 子进程并返回其连接。
func DialProcess(command string, args ...string) (ServerConn, error) {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &processConn{cmd: cmd, stdin: stdin, stdout: stdout}, nil
}

func (p *processConn) Stdin() io.Writer  { return p.stdin }
func (p *processConn) Stdout() io.Reader { return p.stdout }
func (p *processConn) Wait() error       { return p.cmd.Wait() }
func (p *processConn) Close() error      { return p.cmd.Process.Kill() }

// StdioBridge 在 agent 与 MCP server 之间做透明桥接并发布事件。
type StdioBridge struct {
	clientIn  io.Reader
	clientOut io.Writer
	server    ServerConn
	bus       event.Bus
	log       *slog.Logger
	now       func() time.Time

	mu        sync.Mutex
	idToEvent map[string]string
}

// NewStdioBridge 构建桥接器。
func NewStdioBridge(clientIn io.Reader, clientOut io.Writer, server ServerConn, bus event.Bus, log *slog.Logger) *StdioBridge {
	return &StdioBridge{
		clientIn:  clientIn,
		clientOut: clientOut,
		server:    server,
		bus:       bus,
		log:       log,
		now:       time.Now,
		idToEvent: make(map[string]string),
	}
}

// Run 启动双向 pump，阻塞直到任一方向结束。
func (b *StdioBridge) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		b.pump(b.clientIn, b.server.Stdin(), DirInbound)
	}()
	go func() {
		defer wg.Done()
		b.pump(b.server.Stdout(), b.clientOut, DirOutbound)
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		_ = b.server.Close()
		<-done
	}
	return b.server.Wait()
}

func (b *StdioBridge) pump(r io.Reader, w io.Writer, dir Direction) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxMessageSize)

	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		if _, err := w.Write(append(line, '\n')); err != nil {
			b.log.Warn("forward write error", "dir", dir, "err", err)
			return
		}
		b.inspect(line, dir)
	}
	if err := scanner.Err(); err != nil {
		b.log.Warn("scanner error", "dir", dir, "err", err)
	}
}

func (b *StdioBridge) inspect(line []byte, dir Direction) {
	if dir == DirInbound {
		var req Request
		if json.Unmarshal(line, &req) != nil || req.Method != MethodToolsCall {
			return
		}
		b.publishToolCall(req)
		return
	}

	var resp Response
	if json.Unmarshal(line, &resp) != nil || len(resp.ID) == 0 {
		return
	}
	b.publishToolResult(resp)
}

func (b *StdioBridge) publishToolCall(req Request) {
	var params ToolCallParams
	_ = json.Unmarshal(req.Params, &params)

	eventID := model.NewID("evt")
	if len(req.ID) > 0 {
		b.mu.Lock()
		b.idToEvent[string(req.ID)] = eventID
		b.mu.Unlock()
	}

	b.bus.Publish(context.Background(), &model.Event{
		ID:        eventID,
		Type:      model.EventToolCall,
		Protocol:  model.ProtocolMCPStdio,
		SessionID: "stdio",
		Timestamp: b.now(),
		Method:    req.Method,
		ToolCall: &model.ToolCall{
			ID:        string(req.ID),
			Name:      params.Name,
			Arguments: params.Arguments,
		},
	})
}

func (b *StdioBridge) publishToolResult(resp Response) {
	b.mu.Lock()
	parentID := b.idToEvent[string(resp.ID)]
	delete(b.idToEvent, string(resp.ID))
	b.mu.Unlock()

	risk := model.RiskNone
	if resp.Error != nil {
		risk = model.RiskLow
	}
	b.bus.Publish(context.Background(), &model.Event{
		ID:        model.NewID("evt"),
		Type:      model.EventToolResult,
		Protocol:  model.ProtocolMCPStdio,
		SessionID: "stdio",
		Timestamp: b.now(),
		ParentID:  parentID,
		Risk:      risk,
		Response: &model.Payload{
			Body:   mustMarshal(resp),
			Status: statusFor(resp),
		},
	})
}

func statusFor(resp Response) int {
	if resp.Error != nil {
		return resp.Error.Code
	}
	return 200
}

func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return data
}
