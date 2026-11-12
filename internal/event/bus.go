package event

import (
	"context"
	"sync"
	"sync/atomic"

	"vigil/internal/model"
)

const subscriberBuffer = 256

// Bus 是进程内的事件发布订阅总线。
type Bus interface {
	Publish(ctx context.Context, e *model.Event)
	Subscribe() Subscription
}

// Subscription 代表一个订阅句柄。
type Subscription interface {
	Events() <-chan *model.Event
	Dropped() int64
	Close()
}

type memoryBus struct {
	mu   sync.RWMutex
	subs map[int]*subscriber
	next int
}

type subscriber struct {
	ch      chan *model.Event
	dropped int64
	closed  bool
}

func (s *subscriber) Events() <-chan *model.Event { return s.ch }

func (s *subscriber) Dropped() int64 { return atomic.LoadInt64(&s.dropped) }

func (s *subscriber) Close() { s.closed = true }

// New 创建内存事件总线。
func New() Bus {
	return &memoryBus{subs: make(map[int]*subscriber)}
}

func (b *memoryBus) Subscribe() Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := b.next
	b.next++
	sub := &subscriber{ch: make(chan *model.Event, subscriberBuffer)}
	b.subs[id] = sub
	return &subscriptionHandle{bus: b, id: id, sub: sub}
}

func (b *memoryBus) remove(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if sub, ok := b.subs[id]; ok {
		close(sub.ch)
		delete(b.subs, id)
	}
}

// Publish 向所有订阅者非阻塞广播，慢消费者丢弃并计数，绝不拖慢代理。
func (b *memoryBus) Publish(ctx context.Context, e *model.Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, sub := range b.subs {
		if sub.closed {
			continue
		}
		select {
		case sub.ch <- e:
		case <-ctx.Done():
			return
		default:
			atomic.AddInt64(&sub.dropped, 1)
		}
	}
}

type subscriptionHandle struct {
	bus *memoryBus
	id  int
	sub *subscriber
}

func (h *subscriptionHandle) Events() <-chan *model.Event { return h.sub.ch }

func (h *subscriptionHandle) Dropped() int64 { return h.sub.Dropped() }

func (h *subscriptionHandle) Close() { h.bus.remove(h.id) }
