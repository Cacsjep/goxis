package axstorage

/*
#cgo LDFLAGS: -laxstorage
#cgo pkg-config: glib-2.0
#include <glib.h>
#include <axsdk/axstorage.h>
*/
import "C"
import "fmt"

type StorageError struct {
	Ptr      *C.GError
	Code     AXStorageStatusEventId
	Message  string
	Expected bool
}

// AXStorageStatusEventId represents the list of events for AXStorage.
type AXStorageStatusEventId int

const (
	AXStorageAvailableEvent   AXStorageStatusEventId = C.AX_STORAGE_AVAILABLE_EVENT
	AXStorageExitingEvent     AXStorageStatusEventId = C.AX_STORAGE_EXITING_EVENT
	AXStorageWritableEvent    AXStorageStatusEventId = C.AX_STORAGE_WRITABLE_EVENT
	AXStorageFullEvent        AXStorageStatusEventId = C.AX_STORAGE_FULL_EVENT
	AXStorageStatusEventIDEnd AXStorageStatusEventId = C.AX_STORAGE_STATUS_EVENT_ID_END
)

func newStorageError(gerr *C.GError) error {
	if gerr == nil {
		return nil
	}
	err := &StorageError{Message: C.GoString(gerr.message), Code: AXStorageStatusEventId(gerr.code), Ptr: gerr}
	defer C.g_error_free(gerr)
	return err
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("%s, AxStorageError: %d", e.Message, e.Code.ErrorName())
}

// ErrorName returns the string representation of the AXStorageStatusEventId.
func (code AXStorageStatusEventId) ErrorName() string {
	switch code {
	case AXStorageAvailableEvent:
		return "AXStorageAvailableEvent"
	case AXStorageExitingEvent:
		return "AXStorageExitingEvent"
	case AXStorageWritableEvent:
		return "AXStorageWritableEvent"
	case AXStorageFullEvent:
		return "AXStorageFullEvent"
	case AXStorageStatusEventIDEnd:
		return "AXStorageStatusEventIDEnd"
	default:
		return fmt.Sprintf("Unknown AXStorageStatusEvent nor Error code: %d", code)
	}
}
