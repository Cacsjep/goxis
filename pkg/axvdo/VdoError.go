package axvdo

/*
#cgo pkg-config: vdostream
#include "vdo-error.h"
*/
import "C"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-types_8h_source.html
// VdoError represents the error codes for VDO operations.
type VdoError int

const (
	VdoErrorNotFound         VdoError = 1
	VdoErrorExists           VdoError = 2
	VdoErrorInvalidArgument  VdoError = 3
	VdoErrorPermissionDenied VdoError = 4
	VdoErrorNotSupported     VdoError = 5
	VdoErrorClosed           VdoError = 6
	VdoErrorBusy             VdoError = 7
	VdoErrorIO               VdoError = 8
	VdoErrorHAL              VdoError = 9
	VdoErrorDBus             VdoError = 10
	VdoErrorOOM              VdoError = 11 // Out Of Memory
	VdoErrorIdle             VdoError = 12
	VdoErrorNoData           VdoError = 13
	VdoErrorNoBufferSpace    VdoError = 14
	VdoErrorBufferFailure    VdoError = 15
	VdoErrorInterfaceDown    VdoError = 16
	VdoErrorFailed           VdoError = 17
	VdoErrorFatal            VdoError = 18
	VdoErrorNotControlled    VdoError = 19
	VdoErrorNoEvent          VdoError = 20
)

// Check if error is expected.
// Expected errors typically occur as a result of force stopping 'com.axis.Vdo1.System'. This class of errors should not be logged as failures,
// when they occur the recipient is expected to either silently recover or exit informing the user that vdo is currently unavailable.
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-error_8h.html#ac748fae792da6c96a4cba4619a3a3d90
func VdoErrorIsExpected(gerr **C.GError) bool {
	return C.vdo_error_is_expected(gerr) != C.FALSE
}
