package axvdo

/*
#cgo pkg-config: glib-2.0 gio-2.0 gio-unix-2.0 vdostream
#include "vdo-stream.h"
*/
import "C"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html
type VdoStream struct {
	Ptr *C.VdoStream
}

// Create a new VdoStream
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a36fa0021eb58d482c494163db9b22b61
func NewVdoStream(settings *VdoMap) (*VdoStream, error) {
	return nil, nil
}

// Get an existing video stream
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a07e12d4c5d79563413711dcdc6085171
func (v *VdoStream) Get(id int) (*VdoStream, error) {
	return nil, nil
}

// Gets all existing video streams
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#ac21dc2bb1e463b20cca05213739e505c
func (v *VdoStream) GetAll() (*[]VdoStream, error) {
	return nil, nil
}

// Returns the id of this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#afbd3cb015a9d186123d534068db46749
func (v *VdoStream) GetId() int {
	return 0
}

// Returns a file descriptor representing the underlying socket connection.
// The file descriptor returned by this function represents the underlying socket based connection to the vdo service.
// The returned file descriptor can by used as an event source in an event loop to handle asynchronous I/O events.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a81faad176c49398f0a0b9826fb7c31f8
func (v *VdoStream) GetFd() (int, error) {
	return 0, nil
}

// Returns a file descriptor for prioritized events.
// This requires VDO_INTENT_EVENTFD and is intended to be used together with the vdo_stream_get_event function.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#acbdc222c597329e543ba85f330526b12
func (v *VdoStream) GetEventFd() (int, error) {
	return 0, nil
}

// Get the info for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a6368c8c989cbea997947d433dce6a3cb
// VdoMap must unref by user
func (v *VdoStream) GetInfo() (*VdoMap, error) {
	return nil, nil
}

// Get the settings for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#ad1279257f13b9a406a0df2558f10123e
// VdoMap must unref by user
func (v *VdoStream) GetSettings() (*VdoMap, error) {
	return nil, nil
}

// Update the settings for this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#affa9aec868d9c2fd1f39f44aa63adaf2
func (v *VdoStream) SetSettings(settings *VdoMap) error {
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
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#aba5b4264502272caae02a621f3bad63c
func (v *VdoStream) Attach(intent *VdoMap) {}

// Start this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5a366f51af1a7171a6739d191ca1e113
func (v *VdoStream) Start() error {
	return nil
}

// Stop this video stream.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5e28776ff99b3ecf0b996630eacb4f89
func (v *VdoStream) Stop() {}

// Forces this video stream to insert a key frame.
// This function is invoked in order to force a key frame into a video stream.
// Invoking this functionon a video stream with a format that is not a video format (e.g. H.264 or H.265) will have no effect.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a32ba781b100c13c7b70e0f700bbce268
func (v *VdoStream) ForceKeyFrame() error {
	return nil
}

// Allocates a new buffer for this stream.
// Invoke this function in order to allocate a new buffer for this video stream.
// The vdo service performs the actual memory allocation and owns the buffer.
// A handle to the allocated buffer is returned in the form of a VdoBuffer
//
// Note:
//	This function is synchronous and will block until a response from the vdo service is received.
//	This function can only be invoked for non-encoded video streams.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a18ff54d35650a8fa4a00da2198b4b2d3
func (v *VdoStream) BufferAlloc(opaque C.gpointer) (*VdoBuffer, error) {
	return nil, nil
}

// Decreases the reference count for the specified buffer.
// The buffer is freed by the vdo service when the reference count reaches 0.
// Note:
// 	This function is synchronous and will block until a response from the vdo service is received.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a5c13ae89ffee889aebf1a9b6c3bc3594
func (v *VdoStream) BufferUnref(buffer *VdoBuffer) error {
	return nil
}

// Enqueue a buffer for this video stream.
// Invoking this function equeues the specified VdoBuffer in the internal queue of this video stream in order to be subsequently filled with frame data.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a897b1e6d50e00aa8974511b4625cbad3
func (v *VdoStream) BufferEnqueue(buffer *VdoBuffer) error {
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
	return nil, nil
}

// Create and start a new stream to already existing file descriptors.
// This is a convenience function without VdoBuffer support, instead the stream will output data and metadata to two separate file descriptors.
// NOTE:
// 	metadata is not implemented yet, set meta_fd to -1.
// When done with the stream, free it with unref before closing the file descriptors.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#aa327d3bc31376dd3edebc15a4491a13e
func (v *VdoStream) ToFd(settings *VdoMap, data_fd int, meta_fd int) (*VdoStream, error) {
	return nil, nil
}

// Fetches a single VdoBuffer containing a frame.
// Convenience function for fetching a single frame.
// Free the buffer with unref.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#a4b4c0f2124280bde265491c08bd2f47c
func (v *VdoStream) Snapshot(settings *VdoMap) (*VdoBuffer, error) {
	return nil, nil
}

// Fetches the next Event.
// See VdoStreamEvent for the list of possible events, the current event is returned in the VdoMap "event" field.
// The following errors are transient: VDO_ERROR_NO_EVENT The equivalent of EAGAIN, i.e. there's currently no event.
// The following errors are expected to occur during maintenance: VDO_ERROR_INTERFACE_DOWN Recover by waiting for the service to restart.
// Complete VdoStream reinitialization is necessary. All remaining errors are fatal: Complete VdoStream reinitialization is necessary.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-stream_8h.html#af416bd3ec1edf4055c554e0f78b7b9f9
func (v *VdoStream) GetEvent() error {
	return nil
}
