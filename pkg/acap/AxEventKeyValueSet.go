package acap

/*
#cgo pkg-config: axevent
#include <axsdk/axevent.h>
*/
import "C"
import (
	"fmt"
)

// https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax__event__key__value__set_8h.html
type AXEventKeyValueSet struct {
	Ptr      *C.AXEventKeyValueSet
	cStrings []*cString
}

// Creates a new AXEventKeyValueSet
func NewAXEventKeyValueSet() *AXEventKeyValueSet {
	return &AXEventKeyValueSet{
		Ptr: C.ax_event_key_value_set_new(),
	}
}

// Adds an key value to the event set
func (axEventKeyValueSet *AXEventKeyValueSet) AddKeyValue(key string, namespace *string, value interface{}, value_type AXEventValueType) error {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	var gerr *C.GError
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	cValue, err := valueConverter(value, value_type)

	if err != nil {
		return err
	}

	success := C.ax_event_key_value_set_add_key_value(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		cValue,
		C.AXEventValueType(value_type),
		&gerr,
	)

	if int(success) == 0 {
		return newVdoError(gerr)
	}
	return nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetValueType(key string, namespace *string) (AXEventValueType, error) {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	var gerr *C.GError
	var cValueType C.AXEventValueType

	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	success := C.ax_event_key_value_set_get_value_type(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&cValueType,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newVdoError(gerr)
	}
	return AXEventValueType(cValueType), nil
}

// Retrieve the value type of the value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetString(key string, namespace *string) (string, error) {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	var gerr *C.GError
	cValue := newAllocatableCString()
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	success := C.ax_event_key_value_set_get_string(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		cValue.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return "", newVdoError(gerr)
	}
	defer cValue.Free()
	return cValue.ToGolang(), nil
}

// Retrieve the integer value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetInteger(key string, namespace *string) (int, error) {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	var gerr *C.GError
	cValue := newInt()
	success := C.ax_event_key_value_set_get_integer(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&cValue.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newVdoError(gerr)
	}

	return cValue.ToGolang(), nil
}

// Retrieve the boolean value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetBoolean(key string, namespace *string) (bool, error) {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	var gerr *C.GError
	cValue := newBool()
	success := C.ax_event_key_value_set_get_boolean(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&cValue.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return false, newVdoError(gerr)
	}

	return cValue.ToGolang(), nil
}

// Retrieve the double value associated with a key.
func (axEventKeyValueSet *AXEventKeyValueSet) GetDouble(key string, namespace *string) (float64, error) {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	var gerr *C.GError
	cValue := newDouble()
	success := C.ax_event_key_value_set_get_double(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&cValue.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return 0, newVdoError(gerr)
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
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	var gerr *C.GError

	success := C.ax_event_key_value_set_mark_as_source(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return newVdoError(gerr)
	}

	return nil
}

// Mark a key in the AXEventKeyValueSet as data. A data key is a key that represents the state of what the event represents.
// E.g. an event declaration that represents an I/O port should have a key marked as data which represents the state,
// high or low, of the port. Please note that although it is possible to mark more than one key as data,
// only events with one and only one data key can be used to trigger actions.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsData(key string, namespace *string) error {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)
	var gerr *C.GError

	success := C.ax_event_key_value_set_mark_as_data(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return newVdoError(gerr)
	}

	return nil
}

// Mark a key in AXEventKeyValueSet with an user defined tag.
func (axEventKeyValueSet *AXEventKeyValueSet) MarkAsUserDefined(key string, namespace *string, userTag string) error {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	cUserTag := newString(&userTag)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace, cUserTag)

	var gerr *C.GError
	success := C.ax_event_key_value_set_mark_as_user_defined(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		cUserTag.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return newVdoError(gerr)
	}

	return nil
}

// Removes a key and its associated value from an AXEventKeyValueSet.
func (axEventKeyValueSet *AXEventKeyValueSet) RemoveKey(key string, namespace *string) error {
	cKey := newString(&key)
	cNamespace := newString(namespace)
	axEventKeyValueSet.cStrings = append(axEventKeyValueSet.cStrings, cKey, cNamespace)

	var gerr *C.GError
	success := C.ax_event_key_value_set_remove_key(
		axEventKeyValueSet.Ptr,
		cKey.Ptr,
		cNamespace.Ptr,
		&gerr,
	)

	if int(success) == 0 {
		return newVdoError(gerr)
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
