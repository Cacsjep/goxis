package axparameter

/*
#cgo pkg-config: glib-2.0
#include <glib.h>
*/
import "C"
import "fmt"

type GError struct {
	Ptr      *C.GError
	Code     AXParameterErrorCode
	Message  string
	Expected bool
}

func newGError(gerr *C.GError) error {
	if gerr == nil {
		return nil
	}
	err := &GError{Message: C.GoString(gerr.message), Code: AXParameterErrorCode(gerr.code), Ptr: gerr}
	defer C.g_error_free(gerr)
	return err
}

func (e *GError) Error() string {
	return fmt.Sprintf("%s, ErrorCode: %d", e.Message, e.Code)
}
