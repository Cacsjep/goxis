package axevent

/*
#cgo pkg-config: glib-2.0 axevent
#include <axsdk/axevent.h>
#include <glib.h>
#include <stdint.h>
extern void GoSubscriptionCallback(guint subscription, AXEvent *event, gpointer user_data);
extern void GoDeclarationCompleteCallback(guint declaration, gpointer user_data);
*/
import "C"
import (
	"fmt"
	"runtime/cgo"
	"time"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#ae40c17762e9ed663356e34c5a9ea05fe
type AXEventHandler struct {
	Ptr                        *C.AXEventHandler
	subscriptionHandles        map[int]cgo.Handle
	declarationCompleteHandles map[int]cgo.Handle
}

type Subcription uint

// Creates a new AXEventHandler.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#aeb60d443c4006c0deb4ea3763f896de2
func NewEventHandler() *AXEventHandler {
	return &AXEventHandler{
		Ptr:                        C.ax_event_handler_new(),
		subscriptionHandles:        make(map[int]cgo.Handle),
		declarationCompleteHandles: make(map[int]cgo.Handle),
	}
}

// Sends an event.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#ac45d53fe12862e3799e2abb133b50d7a
func (ev *AXEventHandler) SendEvent(declaration int, evt *AXEvent) error {
	var gerr *C.GError
	if int(C.ax_event_handler_send_event(ev.Ptr, C.guint(declaration), evt.Ptr, &gerr)) == 0 {
		return newEventError(gerr)
	}
	evt.Free()
	return nil
}

// Undeclares an event. Any pending callbacks associated with the declaration will be cancelled.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#a61d4484571aacc5547f9e5bf524fd01d
func (evt *AXEventHandler) Undeclare(declaration int) error {
	var gerr *C.GError
	if handle, exists := evt.declarationCompleteHandles[declaration]; exists {
		handle.Delete()
		delete(evt.declarationCompleteHandles, declaration)
	}
	if int(C.ax_event_handler_undeclare(evt.Ptr, C.guint(declaration), &gerr)) == 0 {
		return newEventError(gerr)
	}
	return nil
}

// This is the prototype of the callback function called whenever an event matching a subscription is received.
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#ad1bb63bc12366aefb50139ea6a8b3904
type SubscriptionCallback func(subscription int, event *AXEvent, userdata any)

type subscriptionCallbackData struct {
	Callback SubscriptionCallback
	Userdata any
}

//export GoSubscriptionCallback
func GoSubscriptionCallback(subscription C.guint, event *C.AXEvent, user_data unsafe.Pointer) {
	h := cgo.Handle(user_data)
	data := h.Value().(*subscriptionCallbackData)
	if data == nil {
		fmt.Println("Error: in value conv (GoSubscriptionCallback)")
		return
	}
	evt := &AXEvent{Ptr: event}
	data.Callback(int(subscription), evt, data.Userdata)
	evt.Free()
}

// Subscribes to an event or a set of events.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#a5688da60eeea59fd1d50f394a8177fc9
func (eh *AXEventHandler) Subscribe(kvs *AXEventKeyValueSet, callback SubscriptionCallback, userdata any) (subscription int, err error) {
	var csubscription C.guint
	var gerr *C.GError

	data := &subscriptionCallbackData{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)

	if int(C.ax_event_handler_subscribe(
		eh.Ptr,
		kvs.Ptr,
		&csubscription,
		(C.AXSubscriptionCallback)(C.GoSubscriptionCallback),
		(C.gpointer)(unsafe.Pointer(handle)),
		&gerr,
	)) == 0 {
		return 0, newEventError(gerr)
	}

	eh.subscriptionHandles[int(csubscription)] = handle
	return int(csubscription), nil
}

// This is the prototype of the callback function called when a declaration has registered with the event system.
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#a16d563cbc9c8974b72f296f8dfbdff3a
type DeclarationCompleteCallback func(subscription int, userdata any)

type declarationComplete struct {
	Callback DeclarationCompleteCallback
	Userdata any
}

//export GoDeclarationCompleteCallback
func GoDeclarationCompleteCallback(declaration C.guint, user_data unsafe.Pointer) {
	h := cgo.Handle(user_data)
	data := h.Value().(*declarationComplete)
	if data == nil {
		fmt.Println("Error: in value conv (GoDeclarationCompleteCallback)")
		return
	}
	data.Callback(int(declaration), data.Userdata)
}

// Declares a new event.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#a069c6140889e0b0aec24b8a1f9063ebc
func (eh *AXEventHandler) Declare(keyValueSet *AXEventKeyValueSet, stateless bool, callback DeclarationCompleteCallback, userdata any) (declaration int, err error) {
	var cdeclaration C.guint
	var gerr *C.GError
	data := &declarationComplete{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)

	if int(C.ax_event_handler_declare(
		eh.Ptr,
		keyValueSet.Ptr,
		C.gboolean(map[bool]int{true: 1, false: 0}[stateless]),
		&cdeclaration,
		(C.AXDeclarationCompleteCallback)(C.GoDeclarationCompleteCallback),
		(C.gpointer)(unsafe.Pointer(handle)),
		&gerr,
	)) == 0 {
		return 0, newEventError(gerr)
	}

	eh.declarationCompleteHandles[int(cdeclaration)] = handle
	return int(cdeclaration), nil
}

// Declares a new event based upon an event template
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#aa7d1fd47acc735ded14835b1b2c693f3
func (eh *AXEventHandler) DeclareFromTemplate(template string, keyValueSet *AXEventKeyValueSet, callback DeclarationCompleteCallback, userdata any) (int, error) {
	var declaration C.guint
	var gerr *C.GError
	data := &declarationComplete{Callback: callback, Userdata: userdata}
	handle := cgo.NewHandle(data)

	cTemplate := C.CString(template)
	defer C.free(unsafe.Pointer(cTemplate))

	if int(C.ax_event_handler_declare_from_template(
		eh.Ptr,
		cTemplate,
		keyValueSet.Ptr,
		&declaration,
		(C.AXDeclarationCompleteCallback)(C.GoDeclarationCompleteCallback),
		(C.gpointer)(unsafe.Pointer(handle)),
		&gerr,
	)) == 0 {
		return 0, newEventError(gerr)
	}

	eh.declarationCompleteHandles[int(declaration)] = handle
	return int(declaration), nil
}

// Unsubscribes from an event or a set of events.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#ae18cd4f8d25c6fc555d91fde187dac8d
func (eh *AXEventHandler) Unsubscribe(subscription int) error {
	var gerr *C.GError
	if handle, exists := eh.subscriptionHandles[subscription]; exists {
		handle.Delete()
		delete(eh.subscriptionHandles, subscription)
	}
	if int(C.ax_event_handler_unsubscribe(eh.Ptr, C.guint(subscription), &gerr)) == 0 {
		return newEventError(gerr)
	}
	return nil
}

// Destroys an AXEventHandler an deallocates all associated declarations and subscriptions.
// Any pending callbacks associated with the AXEventHandler will be cancelled.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__handler_8h.html#ac8fa0ee5cba77fffffad4153833b040d
func (eh *AXEventHandler) Free() {
	if (eh.Ptr) == nil {
		return
	}
	C.ax_event_handler_free(eh.Ptr)
	eh.Ptr = nil
}

type Event struct {
	Kvs       *AXEventKeyValueSet
	Timestamp time.Time
}

// OnEvent creates a subscription callback for the given event key value set.
func (eh *AXEventHandler) OnEvent(kvs *AXEventKeyValueSet, callback func(*Event)) (subscription int, err error) {
	subscription, err = eh.Subscribe(kvs, func(subscription int, event *AXEvent, userdata any) {
		callback(&Event{Kvs: event.GetKeyValueSet(), Timestamp: event.GetTimestamp()})
	}, nil)
	kvs.Free()
	return
}
