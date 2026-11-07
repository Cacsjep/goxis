// Package sink is the pure-Go delivery path of axdatahub.Subscription. It is
// pure Go so the host can test it.
//
// Device Data Hub callbacks run on library threads the owner does not control,
// and can fire while, or after, the owner closes the channel they deliver to.
// Destroying the subscriber waits for an in-flight callback, so a callback
// blocked on a full channel nobody reads would hang the teardown forever.
// Sink makes a send safe against a concurrent Close and releases a blocked
// sender when Close begins.
package sink

import "sync"

// Sink sends on a channel that only Close closes.
type Sink[T any] struct {
	mu     sync.RWMutex
	closed bool
	done   chan struct{}
	once   sync.Once
	ch     chan T
}

// New wraps ch.
func New[T any](ch chan T) *Sink[T] {
	return &Sink[T]{ch: ch, done: make(chan struct{})}
}

// Send delivers v, waiting while ch is full. It reports false, and never
// panics, once Close has begun.
func (s *Sink[T]) Send(v T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return false
	}
	select {
	case s.ch <- v:
		return true
	case <-s.done:
		return false
	}
}

// Close releases blocked senders, waits for in-flight ones to leave, then
// closes the channel. Repeated calls do nothing.
func (s *Sink[T]) Close() {
	s.once.Do(func() {
		close(s.done)
		s.mu.Lock()
		s.closed = true
		close(s.ch)
		s.mu.Unlock()
	})
}
