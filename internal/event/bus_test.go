package event

import (
	"context"
	"testing"
	"time"

	"halo/internal/model"
)

func TestPublishSubscribe(t *testing.T) {
	bus := New()
	sub := bus.Subscribe()
	defer sub.Close()

	bus.Publish(context.Background(), &model.Event{ID: "evt-1"})

	select {
	case got := <-sub.Events():
		if got.ID != "evt-1" {
			t.Fatalf("got %q, want evt-1", got.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event")
	}
}

func TestFanOutToMultiple(t *testing.T) {
	bus := New()
	s1 := bus.Subscribe()
	s2 := bus.Subscribe()
	defer s1.Close()
	defer s2.Close()

	bus.Publish(context.Background(), &model.Event{ID: "x"})

	for i, sub := range []Subscription{s1, s2} {
		select {
		case <-sub.Events():
		case <-time.After(time.Second):
			t.Fatalf("subscriber %d did not receive event", i)
		}
	}
}
