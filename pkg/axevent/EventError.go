package axevent

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
*/
import "C"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__error_8h.html
type AXEventErrorCode int

// TODO: Create custom error types, currently only prepared but unused
const (
	AXEventErrorGeneric           AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_GENERIC)
	AXEventErrorInvalidArgument   AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_INVALID_ARGUMENT)
	AXEventErrorIncompatibleValue AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_INCOMPATIBLE_VALUE)
	AXEventErrorDeclaration       AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_DECLARATION)
	AXEventErrorUndeclare         AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_UNDECLARE)
	AXEventErrorSend              AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_SEND)
	AXEventErrorSubscription      AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_SUBSCRIPTION)
	AXEventErrorUnsubscribe       AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_UNSUBSCRIBE)
	AXEventErrorKeyNotFound       AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_KEY_NOT_FOUND)
	AXEventErrorEnd               AXEventErrorCode = AXEventErrorCode(C.AX_EVENT_ERROR_END)
)
