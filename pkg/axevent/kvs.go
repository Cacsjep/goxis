package axevent

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
#include <stdbool.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax__event__key__value__set_8h.html
type AXEventKeyValueSet struct {
	Ptr *C.AXEventKeyValueSet
}

// Creates a new AXEventKeyValueSet
func NewAXEventKeyValueSet() *AXEventKeyValueSet {
	return &AXEventKeyValueSet{
		Ptr: C.ax_event_key_value_set_new(),
	}
}

func namespacePtr(namespace string) *string {
	if namespace == "" {
		return nil
	}
	return &namespace
}

// Adds an key value to the event set
func (axEventKeyValueSet *AXEventKeyValueSet) AddKeyValue(key string, namespace *string, value interface{}, value_type AXEventValueType) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	cValue, err := valueConverter(value, value_type)

	if err != nil {
		return err
	}

	success := C.ax_event_key_value_set_add_key_value(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		cValue,
		C.AXEventValueType(value_type),
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}
	return nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetValueType(key string, namespace *string) (AXEventValueType, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	var cValueType C.AXEventValueType

	success := C.ax_event_key_value_set_get_value_type(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&cValueType,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newEventError(gerr)
	}
	return AXEventValueType(cValueType), nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetString(key string, namespace *string) (string, error) {
	var gerr *C.GError

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var cValue *C.char
	defer func() {
		if cValue != nil {
			C.free(unsafe.Pointer(cValue))
		}
	}()

	success := C.ax_event_key_value_set_get_string(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&cValue,
		&gerr,
	)

	if int(success) == 0 {
		return "", newEventError(gerr)
	}
	return C.GoString(cValue), nil
}

// Retrieve the integer value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetInteger(key string, namespace *string) (int, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	var cValue C.gint
	success := C.ax_event_key_value_set_get_integer(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&cValue,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newEventError(gerr)
	}

	return int(cValue), nil
}

// Retrieve the boolean value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetBoolean(key string, namespace *string) (bool, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	var cValue C.gboolean
	success := C.ax_event_key_value_set_get_boolean(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&cValue,
		&gerr,
	)

	if int(success) == 0 {
		return false, newEventError(gerr)
	}

	return cValue != C.FALSE, nil
}

// Retrieve the double value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetDouble(key string, namespace *string) (float64, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	var cValue C.gdouble
	success := C.ax_event_key_value_set_get_double(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&cValue,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newEventError(gerr)
	}

	return float64(cValue), nil
}

// Mark a key in the AXEventKeyValueSet as a source. A source key is an identifier used to distinguish between
// multiple instances of the same event declaration. E.g. if a device has multiple I/O ports then event declarations
// that represent the state of each port will have the same keys but different values. The key that represents
// which port the event represents should be marked as source and the key which represents the state should be marked
// as data. Please note that although it is possible to mark more than one key as a source,
// only events with zero or one source keys can be used to trigger actions.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsSource(key string, namespace *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError

	success := C.ax_event_key_value_set_mark_as_source(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}

	return nil
}

// Mark a key in the AXEventKeyValueSet as data. A data key is a key that represents the state of what the event represents.
// E.g. an event declaration that represents an I/O port should have a key marked as data which represents the state,
// high or low, of the port. Please note that although it is possible to mark more than one key as data,
// only events with one and only one data key can be used to trigger actions.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsData(key string, namespace *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError

	success := C.ax_event_key_value_set_mark_as_data(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}

	return nil
}

// Mark a key in AXEventKeyValueSet with an user defined tag.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsUserDefined(key string, namespace *string, userTag *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	cUserTag := nilOrCString(userTag)
	if cUserTag != nil {
		defer C.free(unsafe.Pointer(cUserTag))
	}

	var gerr *C.GError
	success := C.ax_event_key_value_set_mark_as_user_defined(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		cUserTag,
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}

	return nil
}

// AddNiceNames sets human-readable names for a key/value pair in an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) AddNiceNames(key string, namespace *string, keyNiceName *string, valueNiceName *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	cKeyNiceName := nilOrCString(keyNiceName)
	if cKeyNiceName != nil {
		defer C.free(unsafe.Pointer(cKeyNiceName))
	}

	cValueNiceName := nilOrCString(valueNiceName)
	if cValueNiceName != nil {
		defer C.free(unsafe.Pointer(cValueNiceName))
	}

	var gerr *C.GError
	success := C.ax_event_key_value_set_add_nice_names(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		cKeyNiceName,
		cValueNiceName,
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}

	return nil
}

// Removes a key and its associated value from an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) RemoveKey(key string, namespace *string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	cNamespace := nilOrCString(namespace)
	if cNamespace != nil {
		defer C.free(unsafe.Pointer(cNamespace))
	}

	var gerr *C.GError
	success := C.ax_event_key_value_set_remove_key(
		axEventKeyValueSet.Ptr,
		cKey,
		cNamespace,
		&gerr,
	)

	if int(success) == 0 {
		return newEventError(gerr)
	}

	return nil
}

// Convert a interface value to correct c type for the set
func valueConverter(value interface{}, value_type AXEventValueType) (C.gconstpointer, error) {
	if value == nil {
		return nil, nil
	}

	switch value_type {
	case AXValueTypeInt:
		var cvalue C.int
		switch v := value.(type) {
		case *int:
			if v == nil {
				return C.gconstpointer(C.NULL), nil
			}
			cvalue = C.int(*v)
		case int:
			cvalue = C.int(v)
		default:
			return nil, fmt.Errorf("unexpected type for AXValueTypeInt, got %T", value)
		}
		return C.gconstpointer(&cvalue), nil

	case AXValueTypeString:
		var cvalue *C.char
		switch v := value.(type) {
		case *string:
			if v == nil {
				return C.gconstpointer(C.NULL), nil
			}
			cvalue = C.CString(*v)
		case string:
			cvalue = C.CString(v)
		default:
			return nil, fmt.Errorf("unexpected type for AXValueTypeString, got %T", value)
		}
		return C.gconstpointer(cvalue), nil

	case AXValueTypeBool:
		var cvalue C.gboolean
		switch v := value.(type) {
		case *bool:
			if v == nil {
				return C.gconstpointer(C.NULL), nil
			}
			if *v {
				cvalue = C.gboolean(1)
			} else {
				cvalue = C.gboolean(0)
			}
		case bool:
			if v {
				cvalue = C.gboolean(1)
			} else {
				cvalue = C.gboolean(0)
			}
		default:
			return nil, fmt.Errorf("unexpected type for AXValueTypeBool, got %T", value)
		}
		return C.gconstpointer(&cvalue), nil

	case AXValueTypeDouble:
		var cvalue C.double
		switch v := value.(type) {
		case *float64:
			if v == nil {
				return C.gconstpointer(C.NULL), nil
			}
			cvalue = C.double(*v)
		case float64:
			cvalue = C.double(v)
		default:
			return nil, fmt.Errorf("unexpected type for AXValueTypeDouble, got %T", value)
		}
		return C.gconstpointer(&cvalue), nil

	default:
		return nil, fmt.Errorf("unexpected value type: %v", value_type)
	}
}

// Frees an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) Free() {
	C.ax_event_key_value_set_free(axEventKeyValueSet.Ptr)
}

// nilOrCString safely converts a Go string pointer to a C string.
// It returns nil if the input pointer is nil, mimicking optional string behavior in C.
func nilOrCString(goStrPtr *string) *C.char {
	if goStrPtr == nil {
		return nil // Return nil to comply with C functions that accept NULL
	}
	return C.CString(*goStrPtr)
}
