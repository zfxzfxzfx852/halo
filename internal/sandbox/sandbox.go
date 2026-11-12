package sandbox

import (
	"context"
	"errors"
	"sort"
	"sync"
)

// DriverProcess 是内置的进程级驱动名称。
const DriverProcess = "process"

// Command 描述待执行的隔离命令。
type Command struct {
	Name string
	Args []string
	Env  map[string]string
	Dir  string
}

// Result 是命令执行结果。
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Driver 是隔离执行驱动；Docker/gVisor/Firecracker/WASM 后续实现此接口。
type Driver interface {
	Name() string
	Run(ctx context.Context, cmd Command) (*Result, error)
}

// ErrNotImplemented 表示请求的驱动尚未实现。
var ErrNotImplemented = errors.New("sandbox driver not implemented")

// Registry 按名称注册驱动，预留多隔离后端扩展点。
type Registry struct {
	mu      sync.RWMutex
	drivers map[string]Driver
}

// NewRegistry 创建注册表并注册内置进程驱动。
func NewRegistry() *Registry {
	r := &Registry{drivers: map[string]Driver{}}
	r.Register(&ProcessDriver{UID: -1, GID: -1})
	return r
}

func (r *Registry) Register(d Driver) {
	r.mu.Lock()
	r.drivers[d.Name()] = d
	r.mu.Unlock()
}

func (r *Registry) Get(name string) (Driver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.drivers[name]
	return d, ok
}

// Names 返回已注册驱动名称，按字母序。
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.drivers))
	for n := range r.drivers {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Run 用指定驱动执行；name 为空时使用默认进程驱动。
func (r *Registry) Run(ctx context.Context, name string, cmd Command) (*Result, error) {
	if name == "" {
		name = DriverProcess
	}
	d, ok := r.Get(name)
	if !ok {
		return nil, ErrNotImplemented
	}
	return d.Run(ctx, cmd)
}
