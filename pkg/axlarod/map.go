package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import "unsafe"

// LarodMap wraps a pointer to a larodMap, which is a key-value store for parameters used in Larod operations.
type LarodMap struct {
	ptr *C.larodMap
}

// LarodMapValueType enumerates the types of values that can be stored in a LarodMap.
type LarodMapValueType int

const (
	LarodMapValueTypeInt LarodMapValueType = iota
	LarodMapValueTypeStr
	LarodMapValueTypeIntArr2
	LarodMapValueTypeIntArr4
)

// LarodMapEntries represents a key-value pair to be used in a LarodMap.
type LarodMapEntries struct {
	Key       string            // Key of the entry.
	Value     interface{}       // Value of the entry, can be of different types based on ValueType.
	ValueType LarodMapValueType // Type of the value stored.
}

// NewLarodMap creates a new LarodMap instance, initializing the underlying larodMap.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#ac7cb26c7b2ef5c99507edd44e7727b0e
func NewLarodMap() (*LarodMap, error) {
	var cError *C.larodError
	cMap := C.larodCreateMap(&cError)
	if cMap == nil {
		return nil, newLarodError(cError)
	}
	return &LarodMap{ptr: cMap}, nil
}

// NewLarodMapWithEntries creates a new LarodMap and populates it with the provided entries.
func NewLarodMapWithEntries(entries []*LarodMapEntries) (*LarodMap, error) {
	lmap, err := NewLarodMap()
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		switch entry.ValueType {
		case LarodMapValueTypeInt:
			if err := lmap.SetInt(entry.Key, entry.Value.(int64)); err != nil {
				return nil, err
			}
		case LarodMapValueTypeStr:
			if err := lmap.SetStr(entry.Key, entry.Value.(string)); err != nil {
				return nil, err
			}
		case LarodMapValueTypeIntArr2:
			if err := lmap.SetIntArr2(entry.Key, entry.Value.([2]int64)); err != nil {
				return nil, err
			}
		case LarodMapValueTypeIntArr4:
			if err := lmap.SetIntArr4(entry.Key, entry.Value.([4]int64)); err != nil {
				return nil, err
			}
		}
	}
	return lmap, nil
}

// SetInt stores an int64 value in the map with the provided key.
func (m *LarodMap) SetInt(key string, value int64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError

	if C.larodMapSetInt(m.ptr, cKey, C.int64_t(value), &cError) == C.bool(false) {
		return newLarodError(cError)
	}

	return nil
}

// GetStr retrieves a string value from the map by its key.
func (m *LarodMap) GetStr(key string) (string, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError

	cValue := C.larodMapGetStr(m.ptr, cKey, &cError)
	if cValue == nil {
		return "", newLarodError(cError)
	}

	return C.GoString(cValue), nil
}

// SetStr stores a string value in the map with the provided key.
func (m *LarodMap) SetStr(key, value string) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))

	var cError *C.larodError
	if C.larodMapSetStr(m.ptr, cKey, cValue, &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}

// GetIntArr2 retrieves an array of two int64 values from the map by its key.
func (m *LarodMap) GetIntArr2(key string) ([2]int64, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError
	cValues := C.larodMapGetIntArr2(m.ptr, cKey, &cError)
	if cValues == nil {
		return [2]int64{}, newLarodError(cError)
	}

	goValues := (*[2]C.int64_t)(unsafe.Pointer(cValues))[:2:2]

	return [2]int64{int64(goValues[0]), int64(goValues[1])}, nil
}

// SetIntArr2 stores an array of two int64 values in the map with the provided key.
func (m *LarodMap) SetIntArr2(key string, value [2]int64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError
	if C.larodMapSetIntArr2(m.ptr, cKey, C.int64_t(value[0]), C.int64_t(value[1]), &cError) == C.bool(false) {
		return newLarodError(cError)
	}
	return nil
}

// GetIntArr4 retrieves an array of four int64 values from the map by its key.
func (m *LarodMap) GetIntArr4(key string) ([4]int64, error) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError
	cValues := C.larodMapGetIntArr4(m.ptr, cKey, &cError)
	if cValues == nil {
		return [4]int64{}, newLarodError(cError)
	}

	goValues := (*[4]C.int64_t)(unsafe.Pointer(cValues))[:4:4]

	return [4]int64{int64(goValues[0]), int64(goValues[1]), int64(goValues[2]), int64(goValues[3])}, nil
}

// SetIntArr4 stores an array of four int64 values in the map with the provided key.
func (m *LarodMap) SetIntArr4(key string, value [4]int64) error {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var cError *C.larodError
	if C.larodMapSetIntArr4(m.ptr, cKey, C.int64_t(value[0]), C.int64_t(value[1]), C.int64_t(value[2]), C.int64_t(value[3]), &cError) == C.bool(false) {
		return newLarodError(cError)
	}

	return nil
}

// Destroy cleans up resources associated with the LarodMap.
func (m *LarodMap) Destroy() {
	if m == nil {
		return
	}
	if m.ptr == nil {
		return
	}
	C.larodDestroyMap(&m.ptr)
}
