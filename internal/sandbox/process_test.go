package sandbox

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunBasic(t *testing.T) {
	r := NewRegistry()
	res, err := r.Run(context.Background(), "", Command{Name: "echo", Args: []string{"hello"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 0 || !strings.Contains(res.Stdout, "hello") {
		t.Fatalf("got %+v", res)
	}
}

func TestExitCode(t *testing.T) {
	r := NewRegistry()
	res, err := r.Run(context.Background(), "", Command{Name: "sh", Args: []string{"-c", "exit 3"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.ExitCode != 3 {
		t.Fatalf("exit %d", res.ExitCode)
	}
}

func TestTimeout(t *testing.T) {
	r := NewRegistry()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err := r.Run(ctx, "", Command{Name: "sleep", Args: []string{"5"}})
	if err != ErrTimeout {
		t.Fatalf("expected timeout, got %v", err)
	}
}

func TestUnknownDriver(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Run(context.Background(), "gvisor", Command{Name: "x"}); err != ErrNotImplemented {
		t.Fatalf("got %v", err)
	}
}

func TestEmptyCommand(t *testing.T) {
	r := NewRegistry()
	if _, err := r.Run(context.Background(), "", Command{}); err == nil {
		t.Fatal("expected error")
	}
}
