package axvdo

/*
#cgo pkg-config: glib-2.0 gio-2.0 gio-unix-2.0 vdostream
#include <vdo-buffer.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html
type VdoBuffer struct {
	Ptr *C.VdoBuffer
}

// Unref/Free the VdoStream
func (b *VdoBuffer) Unref() {
	if b.Ptr != nil {
		C.g_object_unref(C.gpointer(b.Ptr))
	}
}

// NewVdoBuffer creates a buffer owned by an external framework.
// TODO: opaque
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#ad35e3382fa9f9bb4024af62cec90b034
func NewVdoBuffer(fd int, capacity uint, offset uint64) *VdoBuffer {
	return &VdoBuffer{Ptr: C.vdo_buffer_new(C.gint(fd), C.gsize(capacity), C.guint64(offset), nil)}
}

// NewVdoBufferFull creates a buffer with custom properties.
// Default: VDO_BUFFER_ACCESS_ANY_RW (suitable for producers) The fd will not be closed by unref.
// TODO: opaque
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a097843481f6052a45fa46a6ca359750a
func NewVdoBufferFull(fd int, capacity uint, offset uint64, settings *VdoMap) *VdoBuffer {
	return &VdoBuffer{Ptr: C.vdo_buffer_new_full(C.gint(fd), C.gsize(capacity), C.guint64(offset), nil, settings.Ptr)}
}

// GetID returns an ID representing the VdoBuffer.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a5e5ecb68a24fedece4e2b3afb0cac3a7
func (b *VdoBuffer) GetId() int {
	return int(C.vdo_buffer_get_id(b.Ptr))
}

// GetFd returns a file descriptor representing the VdoBuffer.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#ae10562fcb528fc9337767556d6904da8
func (b *VdoBuffer) GetFd() (int, error) {
	id := int(C.vdo_buffer_get_fd(b.Ptr))
	if id == -1 {
		return -1, errors.New("Error on getting fd (vdo_stream_get_fd returns -1)")
	}
	return id, nil
}

// GetOffset returns the file offset to the VdoBuffer.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a00af1cf2ead704a04ca611a7a5f4213c
func (b *VdoBuffer) GetOffset() int64 {
	return int64(C.vdo_buffer_get_offset(b.Ptr))
}

// GetCapacity returns the entire buffer capacity of the VdoBuffer.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a7c5d9bc38c54bb580dab91542ddf3aa2
func (b *VdoBuffer) GetCapacity() uint {
	return uint(C.vdo_buffer_get_capacity(b.Ptr))
}

// IsComplete indicates whether the buffer is complete.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#ab0826cd3642b0cf662c69d84cbaaa2a0
func (b *VdoBuffer) IsComplete() bool {
	return C.vdo_buffer_is_complete(b.Ptr) != C.FALSE
}

// GetOpaque returns user-provided custom information.
// The opaque pointer has no predefined meaning inside the vdo framework itself,
// it is meant to facilitate interoperability with other existing frameworks as well as caching the buffer data pointer(see vdo_buffer_get_data).
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a231bbe5da68795986d481253c879bc11
func (b *VdoBuffer) GetOpaque() unsafe.Pointer {
	return unsafe.Pointer(C.vdo_buffer_get_opaque(b.Ptr))
}

// GetData returns a pointer to the underlying buffer data.
// Note: Beware this pointer is only valid for as long as the VdoBuffer itself is valid.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a3d39e2466d23b62e52a0fbfd10a968a3
func (b *VdoBuffer) GetData() unsafe.Pointer {
	return unsafe.Pointer(C.vdo_buffer_get_data(b.Ptr))
}

// GetFrame returns a pointer to the underlying frame.
// Note: The returned frame type and handling would depend on further definitions and usage.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-buffer_8h.html#a5367ff8edb99976f0b22a2809d44c097
func (b *VdoBuffer) GetFrame() *VdoFrame {
	return &VdoFrame{Ptr: C.vdo_buffer_get_frame(b.Ptr)}
}
