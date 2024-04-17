package axoverlay

/*
#cgo LDFLAGS: -laxoverlay
#cgo pkg-config: glib-2.0
#include <axoverlay.h>
#include <glib.h>
*/
import "C"
import "fmt"

type OverlayError struct {
	Ptr      *C.GError
	Code     AxOverlayErrorCode
	Message  string
	Expected bool
}

type AxOverlayErrorCode int

const (
	AxOverlayErrorInvalidValue       AxOverlayErrorCode = 1000
	AxOverlayErrorInternal           AxOverlayErrorCode = 2000
	AxOverlayErrorUnexpected         AxOverlayErrorCode = 3000
	AxOverlayErrorGeneric            AxOverlayErrorCode = 4000
	AxOverlayErrorInvalidArgument    AxOverlayErrorCode = 5000
	AxOverlayErrorServiceUnavailable AxOverlayErrorCode = 6000
	AxOverlayErrorBackend            AxOverlayErrorCode = 7000
)

func newOverlayError(gerr *C.GError) error {
	if gerr == nil {
		return nil
	}
	err := &OverlayError{Message: C.GoString(gerr.message), Code: AxOverlayErrorCode(gerr.code), Ptr: gerr}
	defer C.g_error_free(gerr)
	return err
}

func (e *OverlayError) Error() string {
	return fmt.Sprintf("%s, AxOverlayError: %s", e.Message, e.Code.ErrorName())
}

// ErrorName returns the string representation of the AxOverlayErrorCode.
func (code AxOverlayErrorCode) ErrorName() string {
	switch code {
	case AxOverlayErrorInvalidValue:
		return "AxOverlayErrorInvalidValue"
	case AxOverlayErrorInternal:
		return "AxOverlayErrorInternal"
	case AxOverlayErrorUnexpected:
		return "AxOverlayErrorUnexpected"
	case AxOverlayErrorGeneric:
		return "AxOverlayErrorGeneric"
	case AxOverlayErrorInvalidArgument:
		return "AxOverlayErrorInvalidArgument"
	case AxOverlayErrorServiceUnavailable:
		return "AxOverlayErrorServiceUnavailable"
	case AxOverlayErrorBackend:
		return "AxOverlayErrorBackend"
	default:
		return fmt.Sprintf("Unknown AxOverlayErrorCode: %d", code)
	}
}
