//go:build cgo

package axdatahub

import (
	"errors"
	"os"
	"testing"
)

// On a machine without libdevice-data-hub-client-c (every host, and cameras
// before AXIS OS 12.10) the package must report ErrUnavailable from every
// entry point, never crash in a forwarder.
func TestWithoutTheLibraryEveryEntryPointReportsUnavailable(t *testing.T) {
	if _, err := os.Stat("/usr/lib/libdevice-data-hub-client-c.so.0"); err == nil {
		t.Skip("this machine has the Device Data Hub client")
	}
	if err := Load(); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Load: %v, want ErrUnavailable", err)
	}
	if _, err := NewClient("test", nil); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("NewClient: %v, want ErrUnavailable", err)
	}
	if _, err := Open(SubscriptionOptions{Topics: []string{TopicSceneFrame}}); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Open: %v, want ErrUnavailable", err)
	}
}

func TestOpenNeedsATopic(t *testing.T) {
	if _, err := Open(SubscriptionOptions{}); err == nil {
		t.Fatal("Open without topics succeeded")
	}
}

func TestDeniedCodes(t *testing.T) {
	for code, want := range map[ErrorCode]bool{
		ErrAuthenticationFailed: true,
		ErrAuthorizationFailed:  true,
		ErrCredential:           true,
		ErrConnection:           false,
		ErrInvalidTopic:         false,
		ErrMaxConnections:       false,
	} {
		if got := (&Error{Code: code}).Denied(); got != want {
			t.Errorf("code %d: Denied()=%t, want %t", code, got, want)
		}
	}
}

func TestMissingTopics(t *testing.T) {
	got := missingTopics([]string{TopicSceneFrame, TopicSceneObjectTrack}, []string{"x", TopicSceneObjectTrack})
	if len(got) != 1 || got[0] != TopicSceneFrame {
		t.Fatalf("got %v", got)
	}
}
