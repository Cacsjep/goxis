//go:build goxis_vdo_frame_legacy

package axvdo

// ACAP Native SDK 13 (AXIS OS 13) removed vdo_frame_get_fd,
// vdo_frame_get_opaque and vdo_frame_get_is_last_buffer, which earlier SDKs
// already marked deprecated. Use the VdoBuffer methods instead: GetFd,
// GetOpaque and IsComplete (VdoBuffer.GetFrame returns the frame of a buffer).
//
// These wrappers build only with the goxis_vdo_frame_legacy build tag, and only
// against an SDK older than 13.

/*
#cgo pkg-config: vdostream
#include <vdo-frame.h>
*/
import "C"
import "unsafe"

// GetFd returns a file descriptor for this frame.
//
// Deprecated: removed in ACAP Native SDK 13; use VdoBuffer.GetFd.
func (f *VdoFrame) GetFd() int {
	return int(C.vdo_frame_get_fd(f.Ptr))
}

// GetOpaque returns a pointer to the opaque data of this frame.
//
// Deprecated: removed in ACAP Native SDK 13; use VdoBuffer.GetOpaque.
func (f *VdoFrame) GetOpaque() unsafe.Pointer {
	return unsafe.Pointer(C.vdo_frame_get_opaque(f.Ptr))
}

// GetIsLastBuffer tests whether this frame is the last buffer.
//
// Deprecated: removed in ACAP Native SDK 13; use VdoBuffer.IsComplete.
func (f *VdoFrame) GetIsLastBuffer() bool {
	return C.vdo_frame_get_is_last_buffer(f.Ptr) != C.FALSE
}
