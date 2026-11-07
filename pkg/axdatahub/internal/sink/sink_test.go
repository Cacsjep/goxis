package sink

import (
	"sync"
	"testing"
	"time"
)

func TestSendDelivers(t *testing.T) {
	ch := make(chan int, 1)
	s := New(ch)
	if !s.Send(7) {
		t.Fatal("Send on an open sink reported false")
	}
	if got := <-ch; got != 7 {
		t.Fatalf("got %d, want 7", got)
	}
}

// A callback blocked on a full channel must be released by Close, or the
// library's destroy (which waits for the callback) never returns.
func TestCloseReleasesABlockedSender(t *testing.T) {
	ch := make(chan int) // nobody reads
	s := New(ch)
	result := make(chan bool)
	go func() { result <- s.Send(1) }()
	time.Sleep(20 * time.Millisecond)
	s.Close()
	select {
	case ok := <-result:
		if ok {
			t.Fatal("a send released by Close reported delivery")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Close did not release the blocked sender")
	}
}

func TestSendAfterCloseDoesNotPanic(t *testing.T) {
	ch := make(chan int, 1)
	s := New(ch)
	s.Close()
	s.Close() // repeated close is a no-op
	if s.Send(1) {
		t.Fatal("Send after Close reported delivery")
	}
	if _, open := <-ch; open {
		t.Fatal("Close did not close the channel")
	}
}

func TestConcurrentSendAndClose(t *testing.T) {
	for i := 0; i < 200; i++ {
		ch := make(chan int, 2)
		s := New(ch)
		var wg sync.WaitGroup
		for j := 0; j < 8; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.Send(j)
			}()
		}
		go func() {
			for range ch {
			}
		}()
		s.Close()
		wg.Wait()
	}
}
