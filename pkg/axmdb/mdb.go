/*
Package axmdb provides a Go interface to the Axis Message Broker API,
showcasing how to use cgo.Handle for callback user data.
*/
package axmdb

/*
#cgo pkg-config: mdb
#include "mdb.h"
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"time"
	"unsafe"
)

// ------------------------------------------------------------
// Types and callbacks
// ------------------------------------------------------------

// MDBConnection wraps the mdb_connection_t type along with
// a cgo.Handle that references the Go callback (ErrorCallback).
type MDBConnection struct {
	ptr       *C.mdb_connection_t // pointer to C connection
	errHandle cgo.Handle          // handle to the Go error callback
	destroyed bool
}

// MDBSubscriberConfig wraps the mdb_subscriber_config_t type along with
// a cgo.Handle for the Go message callback.
type MDBSubscriberConfig struct {
	ptr       *C.mdb_subscriber_config_t
	msgHandle cgo.Handle
	destroyed bool
}

// MDBSubscriber wraps the mdb_subscriber_t type.
type MDBSubscriber struct {
	ptr        *C.mdb_subscriber_t
	doneHandle cgo.Handle // optional handle for the done callback
	destroyed  bool
}

// MDBError wraps the mdb_error_t type. Shown here for reference.
type MDBError struct {
	ptr *C.mdb_error_t
}

// ErrorCallback defines the signature of the error callback function.
type ErrorCallback func(err error)

// Message represents a structured message received from the Message Broker.
type Message struct {
	Timestamp time.Time // Timestamp of the message
	Payload   string    // Payload data
}

// MessageCallback defines the signature for the callback handling messages.
type MessageCallback func(msg *Message)

// DoneCallback defines the signature for a "done" callback,
// e.g. after asynchronously creating a subscriber, we get an optional error.
type DoneCallback func(err error)

// ------------------------------------------------------------
// Utility for converting a C error to Go error
// ------------------------------------------------------------

// NewMDBError creates a Go error from an MDBError pointer.
func NewMDBError(cErr *C.mdb_error_t) error {
	if cErr == nil {
		return nil
	}
	goErr := errors.New(C.GoString(cErr.message))
	// Destroy the C error object so it doesn't leak.
	C.mdb_error_destroy(&cErr)
	return goErr
}

// ------------------------------------------------------------
// Connection error callback
// ------------------------------------------------------------

//export onConnectionErrorCallback
func onConnectionErrorCallback(cError *C.mdb_error_t, userData unsafe.Pointer) {
	// Convert userData to cgo.Handle
	h := cgo.Handle(userData)
	// Attempt to extract our Go ErrorCallback from it
	callback, ok := h.Value().(ErrorCallback)
	if !ok {
		// If not found or the type is wrong, we can do nothing
		return
	}

	var goError error
	if cError != nil {
		goError = errors.New(C.GoString(cError.message))
		C.mdb_error_destroy(&cError)
	}

	// Invoke the Go callback
	callback(goError)
}

// ------------------------------------------------------------
// "Done" Callback for Async Subscriber Creation
// ------------------------------------------------------------

//export onSubscriberCreateDoneCallback
func onSubscriberCreateDoneCallback(cError *C.mdb_error_t, userData unsafe.Pointer) {
	h := cgo.Handle(userData)
	callback, ok := h.Value().(DoneCallback)
	if !ok {
		return
	}

	var goErr error
	if cError != nil {
		goErr = errors.New(C.GoString(cError.message))
		C.mdb_error_destroy(&cError)
	}
	// Invoke the user's Go callback
	callback(goErr)
}

// ------------------------------------------------------------
// Message callback
// ------------------------------------------------------------

//export onMessageCallback
func onMessageCallback(cMessage *C.mdb_message_t, userData unsafe.Pointer) {
	// Convert userData to cgo.Handle
	h := cgo.Handle(userData)
	// Attempt to extract our Go MessageCallback from it
	callback, ok := h.Value().(MessageCallback)
	if !ok {
		return
	}

	// Extract the timestamp from the message
	cTimestamp := C.mdb_message_get_timestamp(cMessage)
	timestamp := time.Unix(
		int64(cTimestamp.tv_sec),
		int64(cTimestamp.tv_nsec),
	)

	// Extract the payload from the message
	cPayload := C.mdb_message_get_payload(cMessage)
	payload := C.GoStringN(
		(*C.char)(unsafe.Pointer(cPayload.data)),
		C.int(cPayload.size),
	)

	// Construct the Go Message
	msg := &Message{
		Timestamp: timestamp,
		Payload:   payload,
	}

	// Invoke the Go callback
	callback(msg)
}

// ------------------------------------------------------------
// Public API for connection creation
// ------------------------------------------------------------

// MDBConnectionCreate creates a new MDB connection, storing
// the given onErr callback in a cgo.Handle.
func MDBConnectionCreate(onErr ErrorCallback) (*MDBConnection, error) {
	// Create a handle from the given Go callback
	errHandle := cgo.NewHandle(onErr)

	var cErr *C.mdb_error_t
	// Call mdb_connection_create with our wrapper
	// so that the signature is consistent.
	connPtr := C.mdb_connection_create(
		(*[0]byte)(C.onConnectionErrorCallback), // function pointer
		unsafe.Pointer(errHandle),               // user data
		&cErr,                                   // out error
	)
	// If an error occurred in C land, handle it
	if cErr != nil {
		// Cleanup the handle, because we won't use it if there's an error
		errHandle.Delete()
		return nil, NewMDBError(cErr)
	}

	// Return the new MDBConnection with the handle
	return &MDBConnection{
		ptr:       connPtr,
		errHandle: errHandle,
	}, nil
}

// Destroy cleans up the MDBConnection, calling the underlying C function.
// Also deletes the cgo.Handle if it hasn't already been destroyed.
func (conn *MDBConnection) Destroy() {
	if conn.destroyed {
		return
	}
	conn.destroyed = true

	if conn.ptr != nil {
		C.mdb_connection_destroy(&conn.ptr)
		conn.ptr = nil
	}
	// Free the cgo.Handle
	conn.errHandle.Delete()
}

// ------------------------------------------------------------
// Subscriber config creation
// ------------------------------------------------------------

// MDBSubscriberConfigCreate creates a subscriber configuration that uses
// a Go callback for messages, storing it in a cgo.Handle.
func MDBSubscriberConfigCreate(topic string, source string, onMessage MessageCallback) (*MDBSubscriberConfig, error) {
	// Convert the message callback into a cgo.Handle
	msgHandle := cgo.NewHandle(onMessage)

	cTopic := C.CString(topic)
	defer C.free(unsafe.Pointer(cTopic))

	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	var cErr *C.mdb_error_t
	configPtr := C.mdb_subscriber_config_create(
		cTopic,
		cSource,
		(*[0]byte)(C.onMessageCallback), // function pointer
		unsafe.Pointer(msgHandle),       // user data
		&cErr,                           // out error
	)
	if cErr != nil {
		msgHandle.Delete()
		return nil, NewMDBError(cErr)
	}
	return &MDBSubscriberConfig{
		ptr:       configPtr,
		msgHandle: msgHandle,
	}, nil
}

// Destroy cleans up the MDBSubscriberConfig, deleting its handle as well.
func (config *MDBSubscriberConfig) Destroy() {
	if config.destroyed {
		return
	}
	config.destroyed = true

	if config.ptr != nil {
		C.mdb_subscriber_config_destroy(&config.ptr)
		config.ptr = nil
	}
	config.msgHandle.Delete()
}

// ------------------------------------------------------------
// MDBSubscriber creation
// ------------------------------------------------------------

// MDBSubscriberCreateAsync creates an async subscriber using the provided
// connection and subscriber config. The user can also pass a pointer to a
// completion callback (`onDoneFunc`) if needed.
func MDBSubscriberCreateAsync(
	conn *MDBConnection,
	config *MDBSubscriberConfig,
	onDone DoneCallback, // Our Go callback to be invoked once creation is done
) (*MDBSubscriber, error) {

	doneHandle := cgo.NewHandle(onDone)

	var cErr *C.mdb_error_t
	// We pass onSubscriberCreateDoneCallback as a function pointer,
	// and the doneHandle as the user data.
	subscriberPtr := C.mdb_subscriber_create_async(
		conn.ptr,
		config.ptr,
		(*[0]byte)(C.onSubscriberCreateDoneCallback), // done callback
		unsafe.Pointer(doneHandle),                   // user data
		&cErr,
	)
	if cErr != nil {
		// On error, free the handle we made
		doneHandle.Delete()
		return nil, NewMDBError(cErr)
	}

	return &MDBSubscriber{
		ptr:        subscriberPtr,
		doneHandle: doneHandle,
	}, nil
}

// Destroy cleans up the MDBSubscriber.
func (subscriber *MDBSubscriber) Destroy() {
	if subscriber.destroyed {
		return
	}
	subscriber.destroyed = true

	if subscriber.ptr != nil {
		C.mdb_subscriber_destroy(&subscriber.ptr)
		subscriber.ptr = nil
	}
	// Free the handle for the done callback
	subscriber.doneHandle.Delete()
}
