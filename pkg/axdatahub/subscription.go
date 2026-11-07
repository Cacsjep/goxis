//go:build cgo

package axdatahub

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/Cacsjep/goxis/pkg/axdatahub/internal/sink"
)

// SubscriptionOptions configures Open.
type SubscriptionOptions struct {
	// Name identifies the client and subscriber on the hub (and in its logs).
	Name string
	// Topics to subscribe to; at least one.
	Topics []string
	// Instances restricts delivery to these instance key sets (OR-ed), for
	// example []map[string]int64{{"channel_id": 1}}. Empty means all.
	Instances []map[string]int64
	// Buffer is the capacity of the Samples channel (default 16). When it is
	// full the library thread waits (back-pressure) until Close.
	Buffer int
	// RequireTopics makes Open fail with ErrInvalidTopic when one of Topics
	// is not on the hub yet, instead of subscribing to a topic that may never
	// appear.
	RequireTopics bool
	// OptionalTopics are subscribed to as well, but only those the hub lists
	// at Open; a missing one is skipped, not an error.
	OptionalTopics []string
}

// ErrDisconnected is the cause of Failed when the hub dropped the connection.
var ErrDisconnected = errors.New("axdatahub: disconnected from the Device Data Hub")

// Subscription is one client, one subscriber and one Go channel of samples:
// the common way to consume topics. A connection loss is reported on Failed;
// the owner then calls Close and opens a new one.
type Subscription struct {
	samples chan Sample
	out     *sink.Sink[Sample]
	client  *Client
	sub     *Subscriber

	failOnce  sync.Once
	failed    chan struct{}
	fatal     error
	closeOnce sync.Once

	delivered atomic.Uint64
	dropped   atomic.Uint64

	// HubTopics lists the topics the hub reported at Open (only with
	// RequireTopics or OptionalTopics). Diagnostic.
	HubTopics []string
	// Topics are the topics actually subscribed to.
	Topics []string
}

// Open connects, subscribes and starts delivery. On error nothing is left
// open. A *Error with Denied() true means the application has no Device Data
// Hub access (the manifest does not enable it); retrying will not help.
func Open(opts SubscriptionOptions) (*Subscription, error) {
	if len(opts.Topics) == 0 {
		return nil, errors.New("axdatahub: no topics to subscribe to")
	}
	if opts.Buffer <= 0 {
		opts.Buffer = 16
	}
	if opts.Name == "" {
		opts.Name = "goxis axdatahub"
	}
	s := &Subscription{
		samples: make(chan Sample, opts.Buffer),
		failed:  make(chan struct{}),
	}
	s.out = sink.New(s.samples)

	client, err := NewClient(opts.Name, func(connected bool) {
		if !connected {
			s.fail(ErrDisconnected)
		}
	})
	if err != nil {
		s.out.Close()
		return nil, err
	}
	s.client = client
	if err := client.Connect(); err != nil {
		s.Close()
		return nil, err
	}
	topics := append([]string(nil), opts.Topics...)
	if opts.RequireTopics || len(opts.OptionalTopics) > 0 {
		names, err := client.TopicNames()
		if err != nil {
			s.Close()
			return nil, err
		}
		s.HubTopics = names
		missing := missingTopics(opts.Topics, names)
		if opts.RequireTopics && len(missing) > 0 {
			s.Close()
			return nil, &Error{
				Op:      "topic check",
				Code:    ErrInvalidTopic,
				Name:    "DH_ERR_INVALID_TOPIC",
				Message: fmt.Sprintf("not on the hub: %s (hub has %d topics)", strings.Join(missing, ", "), len(names)),
			}
		}
		absent := missingTopics(opts.OptionalTopics, names)
		for _, topic := range opts.OptionalTopics {
			if !contains(absent, topic) && !contains(topics, topic) {
				topics = append(topics, topic)
			}
		}
	}
	s.Topics = topics
	sub, err := client.NewSubscriber(opts.Name, func(sample Sample) {
		if s.out.Send(sample) {
			s.delivered.Add(1)
		} else {
			s.dropped.Add(1)
		}
	})
	if err != nil {
		s.Close()
		return nil, err
	}
	s.sub = sub
	if err := sub.Subscribe(Filter{Topics: topics, Instances: opts.Instances}); err != nil {
		s.Close()
		return nil, err
	}
	return s, nil
}

func missingTopics(want, have []string) []string {
	present := make(map[string]struct{}, len(have))
	for _, name := range have {
		present[name] = struct{}{}
	}
	var missing []string
	for _, name := range want {
		if _, ok := present[name]; !ok {
			missing = append(missing, name)
		}
	}
	return missing
}

func contains(list []string, v string) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

// Samples delivers the received samples. It is closed by Close.
func (s *Subscription) Samples() <-chan Sample { return s.samples }

// Failed is closed when the subscription stopped delivering for good (the
// connection was lost). Close it and open a new one.
func (s *Subscription) Failed() <-chan struct{} { return s.failed }

// Err is the reason Failed closed, or nil.
func (s *Subscription) Err() error {
	select {
	case <-s.failed:
		return s.fatal
	default:
		return nil
	}
}

// Delivered and Dropped count samples handed to Samples, and samples that
// arrived while the subscription was closing.
func (s *Subscription) Delivered() uint64 { return s.delivered.Load() }

// Dropped counts samples that arrived while the subscription was closing.
func (s *Subscription) Dropped() uint64 { return s.dropped.Load() }

// fail runs on a library thread: it only signals.
func (s *Subscription) fail(err error) {
	s.failOnce.Do(func() {
		s.fatal = err
		close(s.failed)
	})
}

// Close releases everything. Call it from a Go goroutine, never from a
// callback. It is safe to call more than once.
func (s *Subscription) Close() {
	s.closeOnce.Do(func() {
		// Release a callback blocked on a full Samples channel before the
		// destroys wait for it.
		s.out.Close()
		if s.sub != nil {
			s.sub.Destroy()
		}
		if s.client != nil {
			s.client.Destroy()
		}
	})
}
