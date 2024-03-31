package axvdo

/*
#cgo pkg-config: vdostream
#include "vdo-stream.h"
*/
import "C"
import (
	"unsafe"

	"github.com/Cacsjep/goxis/pkg/glib"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html
type VdoStream struct {
	Ptr *C.VdoStream
}

// Create a new VdoStream
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a36fa0021eb58d482c494163db9b22b61
func NewStream(settings *VdoMap) (*VdoStream, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_new(settings.Ptr, nil, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoStream{Ptr: ptr}, nil
}

// Get an existing video stream
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a07e12d4c5d79563413711dcdc6085171
func StreamGet(id int) (*VdoStream, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_get(C.guint(id), &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoStream{Ptr: ptr}, nil
}

// Gets all existing video streams
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#ac21dc2bb1e463b20cca05213739e505c
func StreamGetAll() ([]*VdoStream, error) {
	var streams []*VdoStream
	var gerr *C.GError
	list_ptr := C.vdo_stream_get_all(&gerr)
	if list_ptr == nil {
		return nil, newVdoError(gerr)
	}
	vdoStreamsPtr := uintptr(unsafe.Pointer(list_ptr))
	vdoStreamsList := glib.WrapList(vdoStreamsPtr)
	vdoStreamsList.DataWrapper(wrapVdoStream)
	vdoStreamsList.Foreach(func(item interface{}) {
		vdoStream, ok := item.(*VdoStream)
		if !ok {
			panic("VdoChannelGetAll: item is not of type *VdoStream")
		}
		streams = append(streams, vdoStream)
	})
	vdoStreamsList.Free()
	return streams, nil
}

// Returns the id of this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#afbd3cb015a9d186123d534068db46749
func (v *VdoStream) GetId() int {
	return int(C.vdo_stream_get_id(v.Ptr))
}

// Returns a file descriptor representing the underlying socket connection.
// The file descriptor returned by this function represents the underlying socket based connection to the vdo service.
// The returned file descriptor can by used as an event source in an event loop to handle asynchronous I/O events.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a81faad176c49398f0a0b9826fb7c31f8
func (v *VdoStream) GetFd() (int, error) {
	var gerr *C.GError
	id := int(C.vdo_stream_get_fd(v.Ptr, &gerr))
	if id == -1 {
		return 0, newVdoError(gerr)
	}
	return id, nil
}

// Returns a file descriptor for prioritized events.
// This requires VDO_INTENT_EVENTFD and is intended to be used together with the vdo_stream_get_event function.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#acbdc222c597329e543ba85f330526b12
func (v *VdoStream) GetEventFd() (int, error) {
	var gerr *C.GError
	id := int(C.vdo_stream_get_event_fd(v.Ptr, &gerr))
	if id == -1 {
		return 0, newVdoError(gerr)
	}
	return id, nil
}

// Get the info for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a6368c8c989cbea997947d433dce6a3cb
// VdoMap must unref by user
func (v *VdoStream) GetInfo() (*VdoMap, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_get_info(v.Ptr, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoMap{Ptr: ptr}, nil
}

// Get the settings for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#ad1279257f13b9a406a0df2558f10123e
// VdoMap must unref by user
func (v *VdoStream) GetSettings() (*VdoMap, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_get_settings(v.Ptr, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoMap{Ptr: ptr}, nil
}

// Update the settings for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#affa9aec868d9c2fd1f39f44aa63adaf2
func (v *VdoStream) SetSettings(settings *VdoMap) error {
	var gerr *C.GError
	if int(C.vdo_stream_set_settings(v.Ptr, settings.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Update the framerate for this video stream.
// This function is invoked in order to update the framerate of a video stream.
// For this API to be used, the stream needs one or more of the following stream settings to be applied before the stream is started:
// 1) Dynamic framerate is enabled (dynamic.framerate = TRUE)
// 2) Zipstream fps mode is set to dynamic (zip.fps_mode = 1)
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#af43f6acfe327b99037ea1e046788e7b5
func (v *VdoStream) SetFramerate(framerate float64) error {
	var gerr *C.GError
	if int(C.vdo_stream_set_framerate(v.Ptr, C.gdouble(framerate), &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Attach to a Stream.
// This function is intended to be used with vdo_stream_get,
// it is redundant to attach to a stream which was created by NewVdoStream.
// VDO_INTENT_CONTROL	Grants start, stop, info, settings, keyframe
// VDO_INTENT_MONITOR	Monitor events using g_signal
// VDO_INTENT_CONSUME	The client intends to stream
// VDO_INTENT_PRODUCE	The client intends to inject frames
// VDO_INTENT_DEFAULT	Grants CONSUME and CONTROL
// VDO_INTENT_EVENTFD	Monitor events using file descriptors
// VDO_INTENT_UNIVERSE	Everything except VDO_INTENT_EVENTFD
//
// vdoMap.SetUint32("intent", VdoIntentEventFD)
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#aba5b4264502272caae02a621f3bad63c
func (v *VdoStream) Attach(intent *VdoMap) error {
	var gerr *C.GError
	if int(C.vdo_stream_attach(v.Ptr, intent.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Start this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5a366f51af1a7171a6739d191ca1e113
func (v *VdoStream) Start() error {
	var gerr *C.GError
	if int(C.vdo_stream_start(v.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Stop this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5e28776ff99b3ecf0b996630eacb4f89
func (v *VdoStream) Stop() {
	C.vdo_stream_stop(v.Ptr)
}

// Forces this video stream to insert a key frame.
// This function is invoked in order to force a key frame into a video stream.
// Invoking this functionon a video stream with a format that is not a video format (e.g. H.264 or H.265) will have no effect.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a32ba781b100c13c7b70e0f700bbce268
func (v *VdoStream) ForceKeyFrame() error {
	var gerr *C.GError
	if int(C.vdo_stream_force_key_frame(v.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Allocates a new buffer for this stream.
// Invoke this function in order to allocate a new buffer for this video stream.
// The vdo service performs the actual memory allocation and owns the buffer.
// A handle to the allocated buffer is returned in the form of a VdoBuffer
//
// Note:
//
//	This function is synchronous and will block until a response from the vdo service is received.
//	This function can only be invoked for non-encoded video streams.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a18ff54d35650a8fa4a00da2198b4b2d3
// TODO: add opaque
func (v *VdoStream) BufferAlloc() (*VdoBuffer, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_buffer_alloc(v.Ptr, nil, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoBuffer{Ptr: ptr}, nil
}

// Decreases the reference count for the specified buffer.
// The buffer is freed by the vdo service when the reference count reaches 0.
// Note:
//
//	This function is synchronous and will block until a response from the vdo service is received.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5c13ae89ffee889aebf1a9b6c3bc3594
func (v *VdoStream) BufferUnref(buffer *VdoBuffer) error {
	var gerr *C.GError
	if int(C.vdo_stream_buffer_unref(v.Ptr, &buffer.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Enqueue a buffer for this video stream.
// Invoking this function equeues the specified VdoBuffer in the internal queue of this video stream in order to be subsequently filled with frame data.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a897b1e6d50e00aa8974511b4625cbad3
func (v *VdoStream) BufferEnqueue(buffer *VdoBuffer) error {
	var gerr *C.GError
	if int(C.vdo_stream_buffer_enqueue(v.Ptr, buffer.Ptr, &gerr)) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Fetches a VdoBuffer containing a frame.
// Stream global settings which control the behavior of this function: "socket.blocking" (default: TRUE) "socket.timeout_ms"
// The following errors are transient: VDO_ERROR_NO_DATA Recover by fetching the next buffer.
// The following errors are expected to occur during maintenance: VDO_ERROR_INTERFACE_DOWN Recover by waiting for the service to restart.
// Complete VdoStream reinitialization is necessary. All remaining errors are fatal: Complete VdoStream reinitialization is necessary.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a021f68451699a9ca04aa0e67d9b2917e
func (v *VdoStream) GetBuffer() (*VdoBuffer, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_get_buffer(v.Ptr, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoBuffer{Ptr: ptr}, nil
}

// Fetches a single VdoBuffer containing a frame.
// Convenience function for fetching a single frame.
// Free the buffer with unref.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a4b4c0f2124280bde265491c08bd2f47c
func Snapshot(settings *VdoMap) (*VdoBuffer, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_snapshot(settings.Ptr, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoBuffer{Ptr: ptr}, nil
}

// Fetches the next Event.
// See VdoStreamEvent for the list of possible events, the current event is returned in the VdoMap "event" field.
// The following errors are transient: VDO_ERROR_NO_EVENT The equivalent of EAGAIN, i.e. there's currently no event.
// The following errors are expected to occur during maintenance: VDO_ERROR_INTERFACE_DOWN Recover by waiting for the service to restart.
// Complete VdoStream reinitialization is necessary. All remaining errors are fatal: Complete VdoStream reinitialization is necessary.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#af416bd3ec1edf4055c554e0f78b7b9f9
func (v *VdoStream) GetEvent() (*VdoMap, error) {
	var gerr *C.GError
	ptr := C.vdo_stream_get_event(v.Ptr, &gerr)
	if ptr == nil {
		return nil, newVdoError(gerr)
	}
	return &VdoMap{Ptr: ptr}, nil
}

// Unref/Free the VdoStream
func (v *VdoStream) Unref() {
	if v.Ptr != nil {
		C.g_object_unref(C.gpointer(v.Ptr))
		v.Ptr = nil
	}
}
