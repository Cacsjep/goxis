/*
Package axdatahub is a Go client for the AXIS Device Data Hub (AXIS OS 12.10
and later), the publish/subscribe service that replaces the Message Broker
(axmdb) on AXIS OS 13. AXIS Scene Metadata is published there on the topics
com.axis.scene.frame.v1 and com.axis.scene.object_track.v1.

The package loads libdevice-data-hub-client-c at runtime instead of linking
it, so an application built with it still starts on cameras that do not have
the library (AXIS OS 12.9 and earlier). Call Load to find out whether the
camera has it. The package builds with any ACAP Native SDK, including those
that ship no Device Data Hub headers.

Device Data Hub access is declared in the manifest:

	"resources": { "deviceDataHub": { "enabled": true } }

(manifest schema 1.11.0 or 2.2.0 and later). Caveat, measured on AXIS OS
13.0.114: the declaration fails the install (acapctl code 10) when the appName
has capitals, because the plugin compares the lower-cased dynamic user with the
appName case-sensitively. Without the declaration the public topics (such as
AXIS Scene Metadata) are still readable on 12.11 and 13.0.114, since the
Message Broker plugin registers every application with the hub.

Library callbacks run on the library's own threads. Client.Destroy (and so
Subscription.Close) waits for in-flight callbacks, so never call it from inside
a callback. A callback that arrives after its owner was destroyed is dropped.
*/
package axdatahub

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include "axdatahub.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime/cgo"
	"sync"
	"time"
	"unsafe"
)

// Load opens the Device Data Hub client library. Every constructor calls it;
// call it directly to probe. It is safe to call repeatedly and from any
// goroutine; the result of the first call is final.
func Load() error {
	if C.axdatahub_load() == 1 {
		return nil
	}
	return fmt.Errorf("%w: %s", ErrUnavailable, C.GoString(C.axdatahub_load_error()))
}

// ErrorCode is a Device Data Hub error code (datahub/common.h).
type ErrorCode uint32

// Error codes from datahub/common.h (SDK 12.11).
const (
	ErrInvalidParams         ErrorCode = 0
	ErrDataTooBig            ErrorCode = 1
	ErrNotConnected          ErrorCode = 2
	ErrAlreadyConnected      ErrorCode = 3
	ErrUnknown               ErrorCode = 4
	ErrInternal              ErrorCode = 5
	ErrInvalidID             ErrorCode = 6
	ErrInvalidTopic          ErrorCode = 7
	ErrNotInitialized        ErrorCode = 8
	ErrMaxClients            ErrorCode = 9
	ErrMaxConnections        ErrorCode = 10
	ErrAuthenticationFailed  ErrorCode = 11
	ErrTopicExists           ErrorCode = 12
	ErrMaxTopic              ErrorCode = 13
	ErrConnection            ErrorCode = 14
	ErrInvalidData           ErrorCode = 15
	ErrInvalidInstance       ErrorCode = 16
	ErrInvalidKeys           ErrorCode = 17
	ErrInvalidSubscription   ErrorCode = 18
	ErrInstanceExists        ErrorCode = 19
	ErrInstancesNotSupported ErrorCode = 20
	ErrMaxInstances          ErrorCode = 21
	ErrMaxSubscriptions      ErrorCode = 22
	ErrAuthorizationFailed   ErrorCode = 23
	ErrMaxProductions        ErrorCode = 24
	ErrInvalidProduction     ErrorCode = 25
	ErrInvalidRequest        ErrorCode = 26
	ErrCredential            ErrorCode = 27
)

// Error is an error reported by the Device Data Hub client.
type Error struct {
	Op      string    // the call that failed, e.g. "dh_client_connect"
	Code    ErrorCode // the library's error code
	Name    string    // the code's name, e.g. "DH_ERR_AUTHENTICATION_FAILED"
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s (%d): %s", e.Op, e.Name, e.Code, e.Message)
}

// Denied reports whether the hub refused this application access: the manifest
// does not grant Device Data Hub access, or the application's user has no
// credentials. Retrying cannot cure it.
func (e *Error) Denied() bool {
	switch e.Code {
	case ErrAuthenticationFailed, ErrAuthorizationFailed, ErrCredential:
		return true
	}
	return false
}

// takeError converts and frees a library error. It returns nil for nil.
func takeError(op string, cErr *C.DHError) error {
	if cErr == nil {
		return nil
	}
	code := C.dh_error_get_code(cErr)
	e := &Error{
		Op:      op,
		Code:    ErrorCode(code),
		Name:    C.GoString(C.dh_error_code_to_string(code)),
		Message: C.GoString(C.dh_error_get_message(cErr)),
	}
	C.dh_error_destroy(cErr)
	return e
}

// failure returns the library error, or a generic one when the call failed
// without setting it.
func failure(op string, cErr *C.DHError) error {
	if err := takeError(op, cErr); err != nil {
		return err
	}
	return &Error{Op: op, Code: ErrUnknown, Name: "DH_ERR_UNKNOWN_ERROR", Message: "the call failed without an error"}
}

// liveHandleValue returns the value behind a cgo.Handle that a C callback
// passed back as user data, or false when the handle was already deleted. A
// deleted handle is never reused, so a late callback is dropped instead of
// panicking ("misuse of an invalid Handle").
func liveHandleValue(h cgo.Handle) (v any, ok bool) {
	defer func() {
		if recover() != nil {
			v, ok = nil, false
		}
	}()
	return h.Value(), true
}

// ConnectionCallback receives connection state changes.
type ConnectionCallback func(connected bool)

// Client is a connection to the Device Data Hub.
type Client struct {
	mu        sync.Mutex
	ptr       *C.DHClient
	onState   cgo.Handle
	hasState  bool
	destroyed bool
}

// NewClient creates a client. onState, if not nil, receives every connection
// state change on a library thread; it must not block and must not call
// Destroy. The client is not connected until Connect.
func NewClient(name string, onState ConnectionCallback) (*Client, error) {
	if err := Load(); err != nil {
		return nil, err
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cErr *C.DHError
	ptr := C.dh_client_create(cName, &cErr)
	if ptr == nil {
		return nil, failure("dh_client_create", cErr)
	}
	_ = takeError("dh_client_create", cErr)
	c := &Client{ptr: ptr}
	// The library's default log target is the console; keep only warnings, in
	// the system log with the application's own messages. Best effort.
	cErr = nil
	C.dh_client_set_logging(ptr, C.AXDATAHUB_LOG_WARNING, C.AXDATAHUB_LOG_TARGET_SYSLOG, &cErr)
	_ = takeError("dh_client_set_logging", cErr)
	if onState != nil {
		c.onState = cgo.NewHandle(onState)
		c.hasState = true
		cErr = nil
		if !C.axdatahub_set_connection_handle(ptr, C.uintptr_t(c.onState), &cErr) {
			err := failure("dh_client_set_connection_update_callback", cErr)
			c.Destroy()
			return nil, err
		}
	}
	return c, nil
}

// Connect connects the client synchronously. A *Error with Denied() true means
// the application has no Device Data Hub access.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.destroyed {
		return errors.New("axdatahub: client destroyed")
	}
	var cErr *C.DHError
	if !C.dh_client_connect(c.ptr, &cErr) {
		return failure("dh_client_connect", cErr)
	}
	return takeError("dh_client_connect", cErr)
}

// TopicNames lists the topics on the hub that this client may see.
func (c *Client) TopicNames() ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.destroyed {
		return nil, errors.New("axdatahub: client destroyed")
	}
	var cErr *C.DHError
	list := C.dh_client_get_topic_list(c.ptr, &cErr)
	if list == nil {
		return nil, failure("dh_client_get_topic_list", cErr)
	}
	_ = takeError("dh_client_get_topic_list", cErr)
	defer C.dh_topic_list_destroy(list)
	n := uint32(C.dh_topic_list_get_count(list))
	names := make([]string, 0, n)
	for i := uint32(0); i < n; i++ {
		if name := C.dh_topic_list_get_name(list, C.uint32_t(i)); name != nil {
			names = append(names, C.GoString(name))
		}
	}
	return names, nil
}

// Destroy disconnects and frees the client. It waits for in-flight callbacks,
// so it must not run inside one. Destroy every Subscriber of the client first.
func (c *Client) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.destroyed {
		return
	}
	c.destroyed = true
	if c.ptr != nil {
		C.dh_client_destroy(c.ptr)
		c.ptr = nil
	}
	// No callback runs after dh_client_destroy returns.
	if c.hasState {
		c.onState.Delete()
	}
}

// Sample is one data sample received on a subscribed topic.
type Sample struct {
	Topic      string
	JSON       string
	SampleID   uint64
	Timestamp  time.Time // the producer's timestamp; zero when absent
	Historical bool      // delivered from history, not live
}

// DataCallback receives samples on a library thread. It may block briefly
// (back-pressure), but a Destroy waits for it.
type DataCallback func(Sample)

// Subscriber receives data samples for one subscription.
type Subscriber struct {
	mu        sync.Mutex
	ptr       *C.DHSubscriber
	onData    cgo.Handle
	destroyed bool
}

// Filter selects the samples a Subscriber receives. Topics is required. Each
// entry of Instances is one set of integer instance keys that must all match
// (for example {"channel_id": 1}); entries are OR-ed. No entries means every
// instance.
type Filter struct {
	Topics    []string
	Instances []map[string]int64
}

// NewSubscriber creates a subscriber on c that delivers samples to onData.
func (c *Client) NewSubscriber(name string, onData DataCallback) (*Subscriber, error) {
	if onData == nil {
		return nil, errors.New("axdatahub: nil data callback")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.destroyed {
		return nil, errors.New("axdatahub: client destroyed")
	}
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	var cErr *C.DHError
	ptr := C.dh_client_create_subscriber(c.ptr, cName, &cErr)
	if ptr == nil {
		return nil, failure("dh_client_create_subscriber", cErr)
	}
	_ = takeError("dh_client_create_subscriber", cErr)
	s := &Subscriber{ptr: ptr, onData: cgo.NewHandle(onData)}
	cErr = nil
	if !C.axdatahub_set_data_handle(ptr, C.uintptr_t(s.onData), &cErr) {
		err := failure("dh_subscriber_set_data_callback", cErr)
		s.Destroy()
		return nil, err
	}
	return s, nil
}

// Subscribe starts delivery for f. A *Error with Code ErrInvalidTopic means a
// topic does not exist on this camera.
func (s *Subscriber) Subscribe(f Filter) error {
	if len(f.Topics) == 0 {
		return errors.New("axdatahub: no topics to subscribe to")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.destroyed {
		return errors.New("axdatahub: subscriber destroyed")
	}
	filter := C.dh_filter_create()
	if filter == nil {
		return errors.New("axdatahub: dh_filter_create failed")
	}
	defer C.dh_filter_destroy(filter)
	for _, topic := range f.Topics {
		cTopic := C.CString(topic)
		var cErr *C.DHError
		ok := C.dh_filter_add_topic_name(filter, cTopic, &cErr)
		C.free(unsafe.Pointer(cTopic))
		if !ok {
			return failure("dh_filter_add_topic_name", cErr)
		}
		_ = takeError("dh_filter_add_topic_name", cErr)
	}
	for _, keys := range f.Instances {
		if err := addInstance(filter, keys); err != nil {
			return err
		}
	}
	options := C.dh_subscribe_options_create()
	if options == nil {
		return errors.New("axdatahub: dh_subscribe_options_create failed")
	}
	defer C.dh_subscribe_options_destroy(options)
	var cErr *C.DHError
	if !C.dh_subscribe_options_add_filter(options, filter, &cErr) {
		return failure("dh_subscribe_options_add_filter", cErr)
	}
	_ = takeError("dh_subscribe_options_add_filter", cErr)
	C.dh_subscribe_options_set_enable_data_updates(options, C.bool(true))
	cErr = nil
	if !C.dh_subscriber_subscribe(s.ptr, options, &cErr) {
		return failure("dh_subscriber_subscribe", cErr)
	}
	return takeError("dh_subscriber_subscribe", cErr)
}

func addInstance(filter *C.DHFilter, keys map[string]int64) error {
	if len(keys) == 0 {
		return nil
	}
	cKeys := C.dh_instance_keys_create()
	if cKeys == nil {
		return errors.New("axdatahub: dh_instance_keys_create failed")
	}
	defer C.dh_instance_keys_destroy(cKeys)
	for name, value := range keys {
		cName := C.CString(name)
		var cErr *C.DHError
		ok := C.dh_instance_keys_add_integer(cKeys, cName, C.int64_t(value), &cErr)
		C.free(unsafe.Pointer(cName))
		if !ok {
			return failure("dh_instance_keys_add_integer", cErr)
		}
		_ = takeError("dh_instance_keys_add_integer", cErr)
	}
	var cErr *C.DHError
	if !C.dh_filter_add_instance(filter, cKeys, &cErr) {
		return failure("dh_filter_add_instance", cErr)
	}
	return takeError("dh_filter_add_instance", cErr)
}

// Destroy unsubscribes and frees the subscriber. It must not run inside a
// callback.
func (s *Subscriber) Destroy() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.destroyed {
		return
	}
	s.destroyed = true
	if s.ptr != nil {
		C.dh_subscriber_destroy(s.ptr)
		s.ptr = nil
	}
	s.onData.Delete()
}

//export axdatahubOnConnectionUpdate
func axdatahubOnConnectionUpdate(update C.DHConnectionState, handle C.uintptr_t) {
	v, live := liveHandleValue(cgo.Handle(handle))
	cb, ok := v.(ConnectionCallback)
	if !live || !ok {
		return
	}
	cb(update == C.AXDATAHUB_CONN_CONNECTED)
}

//export axdatahubOnData
func axdatahubOnData(sample *C.DHTopicSample, handle C.uintptr_t) {
	v, live := liveHandleValue(cgo.Handle(handle))
	cb, ok := v.(DataCallback)
	if !live || !ok || sample == nil {
		return
	}
	out := Sample{
		SampleID:   uint64(C.dh_topic_sample_get_sample_id(sample)),
		Historical: bool(C.dh_topic_sample_is_historical(sample)),
	}
	if name := C.dh_topic_sample_get_topic_name(sample); name != nil {
		out.Topic = C.GoString(name)
	}
	if data := C.dh_topic_sample_get_data(sample); data != nil {
		if js := C.dh_topic_data_get_json_data(data); js != nil {
			out.JSON = C.GoString(js)
		}
	}
	if ts := C.dh_topic_sample_get_timestamp(sample); ts != nil {
		if ms := int64(C.dh_timestamp_get_as_ms(ts)); ms > 0 {
			out.Timestamp = time.UnixMilli(ms)
		}
	}
	cb(out)
}
