package axevent

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/Cacsjep/goxis/pkg/clib"
)

// https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax__event__key__value__set_8h.html
type AXEventKeyValueSet struct {
	Ptr      *C.AXEventKeyValueSet
	cStrings []*clib.String
}

// Creates a new AXEventKeyValueSet
func NewAXEventKeyValueSet() *AXEventKeyValueSet {
	return &AXEventKeyValueSet{
		Ptr: C.ax_event_key_value_set_new(),
	}
}

// Adds an key value to the event set
func (axEventKeyValueSet *AXEventKeyValueSet) AddKeyValue(key string, namespace *string, value interface{}, value_type AXEventValueType) error {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	cError := clib.NewError()
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	cValue, err := valueConverter(value, value_type)

	if err != nil {
		return err
	}

	success := C.ax_event_key_value_set_add_key_value(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		cValue,
		C.AXEventValueType(value_type),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to add key value to set"); err != nil {
		return err
	}
	return nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetValueType(key string, namespace *string) (AXEventValueType, error) {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	cError := clib.NewError()
	var cValueType C.AXEventValueType

	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	success := C.ax_event_key_value_set_get_value_type(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		&cValueType,
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get value type"); err != nil {
		return 0, err
	}
	return AXEventValueType(cValueType), nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetString(key string, namespace *string) (string, error) {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	cError := clib.NewError()
	cValue := clib.NewAllocatableCString()
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	success := C.ax_event_key_value_set_get_string(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(**C.char)(unsafe.Pointer(cValue.Ptr)),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get string value"); err != nil {
		return "", err
	}
	defer cValue.Free()
	return cValue.ToGolang(), nil
}

// Retrieve the integer value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetInteger(key string, namespace *string) (int, error) {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	cError := clib.NewError()
	cValue := clib.NewInt()
	success := C.ax_event_key_value_set_get_integer(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(*C.gint)(&cValue.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get integer value"); err != nil {
		return 0, err
	}

	return cValue.ToGolang(), nil
}

// Retrieve the boolean value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetBoolean(key string, namespace *string) (bool, error) {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	cError := clib.NewError()
	cValue := clib.NewBool()
	success := C.ax_event_key_value_set_get_boolean(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(*C.gboolean)(&cValue.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get boolean value"); err != nil {
		return false, err
	}

	return cValue.ToGolang(), nil
}

// Retrieve the double value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetDouble(key string, namespace *string) (float64, error) {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	cError := clib.NewError()
	cValue := clib.NewDouble()
	success := C.ax_event_key_value_set_get_double(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(*C.gdouble)(&cValue.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get double value"); err != nil {
		return 0, err
	}

	return cValue.ToGolang(), nil
}

// Mark a key in the AXEventKeyValueSet as a source. A source key is an identifier used to distinguish between
// multiple instances of the same event declaration. E.g. if a device has multiple I/O ports then event declarations
// that represent the state of each port will have the same keys but different values. The key that represents
// which port the event represents should be marked as source and the key which represents the state should be marked
// as data. Please note that although it is possible to mark more than one key as a source,
// only events with zero or one source keys can be used to trigger actions.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsSource(key string, namespace *string) error {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	cError := clib.NewError()

	success := C.ax_event_key_value_set_mark_as_source(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to mark key as source"); err != nil {
		return err
	}

	return nil
}

// Mark a key in the AXEventKeyValueSet as data. A data key is a key that represents the state of what the event represents.
// E.g. an event declaration that represents an I/O port should have a key marked as data which represents the state,
// high or low, of the port. Please note that although it is possible to mark more than one key as data,
// only events with one and only one data key can be used to trigger actions.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsData(key string, namespace *string) error {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	cError := clib.NewError()

	success := C.ax_event_key_value_set_mark_as_data(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to mark key as data"); err != nil {
		return err
	}

	return nil
}

// Mark a key in AXEventKeyValueSet with an user defined tag.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsUserDefined(key string, namespace *string, userTag string) error {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	cUserTag := clib.NewString(&userTag)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace, cUserTag)

	cError := clib.NewError()
	defer cError.Free()

	success := C.ax_event_key_value_set_mark_as_user_defined(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(*C.char)(cUserTag.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to mark key as user-defined"); err != nil {
		return err
	}

	return nil
}

// Removes a key and its associated value from an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) RemoveKey(key string, namespace *string) error {
	cKey := clib.NewString(&key)
	cNamespace := clib.NewString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	cError := clib.NewError()
	success := C.ax_event_key_value_set_remove_key(
		axEventKeyValueSet.Ptr,
		(*C.char)(cKey.Ptr),
		(*C.char)(cNamespace.Ptr),
		(**C.GError)(unsafe.Pointer(cError.Ptr)),
	)

	if err := cError.IsErrorOrNotSuccess(int(success), "Unable to get remove key"); err != nil {
		return err
	}

	return nil
}

// Convert a interface value to correct c type for the set
func valueConverter(value interface{}, value_type AXEventValueType) (C.gconstpointer, error) {
	if value != nil {
		switch value_type {
		case AXValueTypeInt:
			intval := value.(int)
			cvalue := C.int(intval)
			return C.gconstpointer(&cvalue), nil
		case AXValueTypeString:
			strval := value.(string)
			cvalue := C.CString(strval)
			return C.gconstpointer(cvalue), nil
		case AXValueTypeBool:
			var cvalue C.gboolean
			if value.(bool) {
				cvalue = C.gboolean(1)
			} else {
				cvalue = C.gboolean(0)
			}
			return C.gconstpointer(&cvalue), nil
		case AXValueTypeDouble:
			floatval := value.(float64)
			cvalue := C.double(floatval)
			return C.gconstpointer(&cvalue), nil
		default:
			return nil, fmt.Errorf("unexpected type")
		}
	} else {
		return nil, nil
	}
}

// Frees an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) Free() {
	for _, cs := range axEventKeyValueSet.cStrings {
		cs.Free()
	}
	C.ax_event_key_value_set_free(axEventKeyValueSet.Ptr)
}
