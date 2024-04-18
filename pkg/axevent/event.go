/*
Package axeventwrapper provides a Go wrapper around the axevent C library, facilitating the integration and manipulation of event systems within Axis devices.

Important:

	This package relies on CGo for integration with the axevent C library and the GLib library, necessitating a running GMainLoop for proper operation.
*/
package axevent

/*
#cgo pkg-config: glib-2.0 axevent
#include <axsdk/axevent.h>
#include <glib.h>

long ax_event_get_time_stamp_unix(GDateTime *gdateTime) {
    return g_date_time_to_unix(gdateTime);
}
*/
import "C"
import "time"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event_8h.html#a06778bcedc5cf3aaba11d40fba6bef33
type AXEvent struct {
	Ptr *C.AXEvent
}

// Creates a new AXEvent.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event_8h.html#a788775473fa91f3503b42616c7626e2b
func NewAxEvent(axEventKeyValueSet *AXEventKeyValueSet, datetime *time.Time) *AXEvent {
	var cDateTime *C.GDateTime
	if datetime != nil {
		unixTimestamp := datetime.Unix()
		cDateTime = C.g_date_time_new_from_unix_local(C.gint64(unixTimestamp))
	}

	defer axEventKeyValueSet.Free()
	return &AXEvent{
		Ptr: C.ax_event_new2(axEventKeyValueSet.Ptr, cDateTime),
	}
}

// Get the AXEventKeyValueSet associated with the AXEvent.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event_8h.html#abc27a691c703d11563ef0b0f338fc775
func (axEvent *AXEvent) GetKeyValueSet() *AXEventKeyValueSet {
	return &AXEventKeyValueSet{Ptr: C.ax_event_get_key_value_set(axEvent.Ptr)}
}

// Get the AXEventKeyValueSet associated with the AXEvent.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event_8h.html#a37fcd4106a9ed74e315bbbec24c941fa
func (axEvent *AXEvent) GetTimestamp() time.Time {
	gdateTime := C.ax_event_get_time_stamp2(axEvent.Ptr)
	unixTimestamp := C.ax_event_get_time_stamp_unix(gdateTime)
	return time.Unix(int64(unixTimestamp), 0)
}

// Free an AXEvent.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event_8h.html#a011c2d3b82c8e9cbcf0fab02610a5020
func (axEvent *AXEvent) Free() {
	C.ax_event_free(axEvent.Ptr)
}
