package acap

/*
#cgo pkg-config: vdostream
#include "vdo-error.h"
*/
import "C"
import "fmt"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-types_8h_source.html
// VdoError represents the error codes for VDO operations.
// Used internaly
type VdoErrorCode int

const (
	VdoErrorCodeNotFound VdoErrorCode = iota + 1
	VdoErrorCodeExists
	VdoErrorCodeInvalidArgument
	VdoErrorCodePermissionDenied
	VdoErrorCodeNotSupported
	VdoErrorCodeClosed
	VdoErrorCodeBusy
	VdoErrorCodeIO
	VdoErrorCodeHAL
	VdoErrorCodeDBus
	VdoErrorCodeOOM
	VdoErrorCodeIdle
	VdoErrorCodeNoData
	VdoErrorCodeNoBufferSpace
	VdoErrorCodeBufferFailure
	VdoErrorCodeInterfaceDown
	VdoErrorCodeFailed
	VdoErrorCodeFatal
	VdoErrorCodeNotControlled
	VdoErrorCodeNoEvent
)

// Check if error is expected.
// Expected errors typically occur as a result of force stopping 'com.axis.Vdo1.System'. This class of errors should not be logged as failures,
// when they occur the recipient is expected to either silently recover or exit informing the user that vdo is currently unavailable.
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-error_8h.html#ac748fae792da6c96a4cba4619a3a3d90
func VdoErrorIsExpected(gerr **C.GError) bool {
	return C.vdo_error_is_expected(gerr) != C.FALSE
}

// VdoStream Error, a glib error but with vdo error codes
type VdoError struct {
	Code     VdoErrorCode
	Message  string
	Expected bool
}

func newVdoError(gerr *C.GError) *VdoError {
	if gerr == nil {
		return nil
	}
	return &VdoError{Message: C.GoString(gerr.message), Code: VdoErrorCode(gerr.code), Expected: VdoErrorIsExpected(&gerr)}
}

func (e *VdoError) Error() string {
	return fmt.Sprintf("VdoError: %s, VdoErrorCode: %d, Expected: %t", e.Message, e.Code, e.Expected)
}

func (e *VdoError) ErrorCodeString() string {
	return [...]string{
		"NotFound",
		"Exists",
		"InvalidArgument",
		"PermissionDenied",
		"NotSupported",
		"Closed",
		"Busy",
		"IO",
		"HAL",
		"DBus",
		"OOM",
		"Idle",
		"NoData",
		"NoBufferSpace",
		"BufferFailure",
		"InterfaceDown",
		"Failed",
		"Fatal",
		"NotControlled",
		"NoEvent",
	}[e.Code-1]
}
