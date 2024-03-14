package clib

/*
#cgo pkg-config: glib-2.0
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// String wraps a C string (*C.char) to manage its lifecycle in Go.
// It provides methods for converting to Go strings and freeing allocated memory.
type String struct {
	Ptr *C.char
}

// NewString creates a new String instance.
// If `str` is not nil, it creates a C string from the Go string and stores its pointer.
// If `str` is nil, it initializes the String with a nil pointer, suitable for direct C library use.
func NewString(str *string) *String {
	if str == nil {
		return &String{Ptr: nil}
	}
	cstr := C.CString(*str)
	return &String{Ptr: cstr}
}

// ToGolang converts the C string to a Go string.
// It returns an empty string if the C string is nil.
func (cs *String) ToGolang() string {
	if cs.Ptr != nil {
		return C.GoString(cs.Ptr)
	}
	return ""
}

// Free releases the memory allocated for the C string.
// It is safe to call multiple times, but the String should not be used after being freed.
func (cs *String) Free() {
	if cs.Ptr != nil {
		C.free(unsafe.Pointer(cs.Ptr))
		cs.Ptr = nil
	}
}

// AllocatableCString wraps a pointer to a C string pointer to facilitate
// receiving strings allocated by C functions.
type AllocatableCString struct {
	Ptr **C.char
}

// NewAllocatableCString creates an AllocatableCString instance.
// This is useful when a C function expects to allocate a string and return it via a pointer argument.
func NewAllocatableCString() *AllocatableCString {
	var placeholder *C.char
	return &AllocatableCString{Ptr: &placeholder}
}

// ToGolang converts the allocated C string to a Go string and frees the C string.
// It returns an empty string if the C string is nil.
func (ms *AllocatableCString) ToGolang() string {
	if *ms.Ptr != nil {
		return C.GoString(*ms.Ptr)
	}
	return ""
}

// Free explicitly frees the memory allocated for the C string.
// Note: ToGolang already frees the memory, so calling Free after ToGolang is not necessary.
func (ms *AllocatableCString) Free() {
	if *ms.Ptr != nil {
		C.free(unsafe.Pointer(*ms.Ptr))
		*ms.Ptr = nil
	}
}
