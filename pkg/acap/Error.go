package acap

/*
#cgo pkg-config: vdostream
#include "vdo-error.h"
*/
import "C"
import "fmt"

type GError struct {
	Ptr      *C.GError
	Code     int
	Message  string
	Expected bool
}

func newGError(gerr *C.GError) *GError {
	if gerr == nil {
		return nil
	}
	return &GError{Message: C.GoString(gerr.message), Code: int(gerr.code), Ptr: gerr}
}

func (e *GError) Error() string {
	return fmt.Sprintf("%s, ErrorCode: %d", e.Message, e.Code)
}

func (e *GError) Free() {
	C.g_error_free(e.Ptr)
}
