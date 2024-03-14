package axvdo

/*
#cgo pkg-config: glib-2.0 gio-2.0 gio-unix-2.0 vdostream
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
