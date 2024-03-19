package acap

/*
#cgo pkg-config: vdostream
#include <vdo-frame.h>
*/
import "C"
import (
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-frame_8h.html
type VdoFrame struct {
	Ptr *C.VdoFrame
}

// Unref/Free the VdoFrame
func (b *VdoFrame) Unref() {
	if b.Ptr != nil {
		C.g_object_unref(C.gpointer(b.Ptr))
	}
}

// GetFrameType returns the type of this frame.
func (f *VdoFrame) GetFrameType() VdoFrameType {
	return VdoFrameType(C.vdo_frame_get_frame_type(f.Ptr))
}

// GetSequenceNbr returns the sequence number of this frame.
func (f *VdoFrame) GetSequenceNbr() uint {
	return uint(C.vdo_frame_get_sequence_nbr(f.Ptr))
}

// GetTimestamp returns the timestamp of this frame.
func (f *VdoFrame) GetTimestamp() uint64 {
	return uint64(C.vdo_frame_get_timestamp(f.Ptr))
}

// GetCustomTimestamp returns a custom timestamp for this frame.
func (f *VdoFrame) GetCustomTimestamp() int64 {
	return int64(C.vdo_frame_get_custom_timestamp(f.Ptr))
}

// GetSize returns the size of this frame.
func (f *VdoFrame) GetSize() uint {
	return uint(C.vdo_frame_get_size(f.Ptr))
}

// GetHeaderSize returns the size of any H264/H265 headers contained in this frame.
func (f *VdoFrame) GetHeaderSize() int {
	return int(C.vdo_frame_get_header_size(f.Ptr))
}

// GetFd returns a file descriptor for this frame.
func (f *VdoFrame) GetFd() int {
	return int(C.vdo_frame_get_fd(f.Ptr))
}

// GetExtraInfo returns the extra info of this frame.
func (f *VdoFrame) GetExtraInfo() *VdoMap {
	return &VdoMap{Ptr: C.vdo_frame_get_extra_info(f.Ptr)}
}

// GetOpaque returns a pointer to the opaque data of this frame.
func (f *VdoFrame) GetOpaque() unsafe.Pointer {
	return unsafe.Pointer(C.vdo_frame_get_opaque(f.Ptr))
}

// GetIsLastBuffer tests whether this frame is the last buffer.
func (f *VdoFrame) GetIsLastBuffer() bool {
	return C.vdo_frame_get_is_last_buffer(f.Ptr) != C.FALSE
}

// SetSize sets the size of this frame.
func (f *VdoFrame) SetSize(size uint) {
	C.vdo_frame_set_size(f.Ptr, C.gsize(size))
}

// SetFrameType sets the type of this frame.
func (f *VdoFrame) SetFrameType(frameType VdoFrameType) {
	C.vdo_frame_set_frame_type(f.Ptr, C.VdoFrameType(frameType))
}

// SetSequenceNbr sets the sequence number of this frame.
func (f *VdoFrame) SetSequenceNbr(seqNum uint) {
	C.vdo_frame_set_sequence_nbr(f.Ptr, C.guint(seqNum))
}

// SetTimestamp sets the timestamp of this frame.
func (f *VdoFrame) SetTimestamp(timestamp uint64) {
	C.vdo_frame_set_timestamp(f.Ptr, C.guint64(timestamp))
}

// SetCustomTimestamp sets a custom timestamp for this frame.
func (f *VdoFrame) SetCustomTimestamp(timestamp int64) {
	C.vdo_frame_set_custom_timestamp(f.Ptr, C.gint64(timestamp))
}

// SetIsLastBuffer marks this frame as the last buffer.
func (f *VdoFrame) SetIsLastBuffer(isLastBuffer bool) {
	C.vdo_frame_set_is_last_buffer(f.Ptr, goBooleanToC(isLastBuffer))
}

// SetExtraInfo sets the extra info of this frame.
func (f *VdoFrame) SetExtraInfo(extraInfo *VdoMap) {
	C.vdo_frame_set_extra_info(f.Ptr, extraInfo.Ptr)
}

// SetHeaderSize sets the header size of this frame, used for H26x frames.
func (f *VdoFrame) SetHeaderSize(size int) {
	C.vdo_frame_set_header_size(f.Ptr, C.gssize(size))
}

// Memmap maps the frame into current process memory. Returns a pointer to the data.
func (f *VdoFrame) Memmap() unsafe.Pointer {
	return unsafe.Pointer(C.vdo_frame_memmap(f.Ptr))
}

// Unmap unmaps the frame from current process memory.
func (f *VdoFrame) Unmap() {
	C.vdo_frame_unmap(f.Ptr)
}

// TODO: vdo_frame_take_chunk
