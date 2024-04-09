package axevent

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
*/
import "C"
import "fmt"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axevent/html/ax__event__error_8h.html
type AXEventErrorCode int

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

type EventError struct {
	Ptr      *C.GError
	Code     AXEventErrorCode
	Message  string
	Expected bool
}

func newEventError(gerr *C.GError) error {
	if gerr == nil {
		return nil
	}
	err := &EventError{Message: C.GoString(gerr.message), Code: AXEventErrorCode(gerr.code), Ptr: gerr}
	defer C.g_error_free(gerr)
	return err
}

func (e *EventError) Error() string {
	return fmt.Sprintf("%s, AxEventError: %s", e.Message, e.Code.ErrorName())
}

// ErrorName returns the string representation of the AxOverlayErrorCode.
func (code AXEventErrorCode) ErrorName() string {
	switch code {
	case AXEventErrorGeneric:
		return "AXEventErrorGeneric"
	case AXEventErrorInvalidArgument:
		return "AXEventErrorInvalidArgument"
	case AXEventErrorIncompatibleValue:
		return "AXEventErrorIncompatibleValue"
	case AXEventErrorDeclaration:
		return "AXEventErrorDeclaration"
	case AXEventErrorUndeclare:
		return "AXEventErrorUndeclare"
	case AXEventErrorSend:
		return "AXEventErrorSend"
	case AXEventErrorSubscription:
		return "AXEventErrorSubscription"
	case AXEventErrorUnsubscribe:
		return "AXEventErrorUnsubscribe"
	case AXEventErrorKeyNotFound:
		return "AXEventErrorKeyNotFound"
	case AXEventErrorEnd:
		return "AXEventErrorEnd"
	default:
		return fmt.Sprintf("Unknown AXEventError: %d", code)
	}
}
