package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"syscall"
)

// ErrTimeout 表示命令因超时被终止。
var ErrTimeout = errors.New("command timed out")

// ProcessDriver 是进程级隔离驱动：独立进程组、可选权限降级、超时整组终止。
type ProcessDriver struct {
	// UID/GID >= 0 时启用权限降级；-1 表示不设置。
	UID int
	GID int
}

func (p *ProcessDriver) Name() string { return DriverProcess }

func (p *ProcessDriver) Run(ctx context.Context, c Command) (*Result, error) {
	if c.Name == "" {
		return nil, errors.New("empty command")
	}

	cmd := exec.CommandContext(ctx, c.Name, c.Args...)
	cmd.Dir = c.Dir
	cmd.Env = envList(c.Env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if p.UID >= 0 || p.GID >= 0 {
		uid, gid := 0, 0
		if p.UID >= 0 {
			uid = p.UID
		}
		if p.GID >= 0 {
			gid = p.GID
		}
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
	}

	// 超时/取消时杀掉整个进程组，避免子进程残留。
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	res := &Result{}
	err := cmd.Run()
	res.Stdout = stdout.String()
	res.Stderr = stderr.String()

	// 超时优先判定：进程被 SIGKILL 终止时也会带 ExitError。
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return res, ErrTimeout
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		res.ExitCode = exitErr.ExitCode()
		return res, nil
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

func envList(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}
