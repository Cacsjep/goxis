package axvdo

/*
#cgo pkg-config: vdostream
#include "vdo-map.h"
*/
import "C"
import (
	"github.com/Cacsjep/goxis/pkg/clib"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/classVdoMap.html
type VdoMap struct {
	Ptr *C.VdoMap
}

// Creates a new VdoMap
func NewVdoMap() *VdoMap {
	return &VdoMap{Ptr: C.vdo_map_new()}
}

// Creates a new VdoMap from C
func NewVdoMapFromC(ptr *C.VdoMap) *VdoMap {
	return &VdoMap{Ptr: ptr}
}

// Checks if this map is empty.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a41324ec82e93a622f9073cf708a8545e
func (v *VdoMap) Empty() bool {
	b := C.vdo_map_empty(v.Ptr)
	if int(b) == 1 {
		return true
	}
	return false
}

// Returns the number of entries in this map.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#aa4f3655d8be32e7082f0a8b820ec0651
func (v *VdoMap) Size() int {
	return int(C.vdo_map_size(v.Ptr))
}

// Swaps the contents of two maps.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a3d7c695daa085566e776e5c96dc3d7ee
func (v *VdoMap) Swap(rhs *VdoMap) {
	C.vdo_map_swap(v.Ptr, rhs.Ptr)
}

// Checks if this map contains the specified key.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#afba46021d63a6747ecd968146a8bfc93
func (v *VdoMap) Contains(name string) bool {
	cName := clib.NewString(&name)
	defer cName.Free()
	b := C.vdo_map_contains(v.Ptr, (*C.char)(cName.Ptr))
	if int(b) == 1 {
		return true
	}
	return false
}

// Checks if all entries in this map and the specified map are equal.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a2c2af191d0695a8e3378fd8e591e7cbd
func (v *VdoMap) Equal(amp *VdoMap) bool {
	b := C.vdo_map_equals(v.Ptr, amp.Ptr)
	if int(b) == 1 {
		return true
	}
	return false
}

// Merges the specified map into this map.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#ae3d064d00161bf66423fa7f86b8b1d87
func (v *VdoMap) Merge(amp *VdoMap) {
	C.vdo_map_merge(v.Ptr, amp.Ptr)
}

// Removes the entry with the specified key from this map.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a6f649828e6de0e8147db2a0e741de978
func (v *VdoMap) Remove(name string) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_remove(v.Ptr, (*C.char)(cName.Ptr))
}

// Print a string representation of this map to stdout.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a98b42a54524038a2067ef2c7015c070b
func (v *VdoMap) Dump() {
	C.vdo_map_dump(v.Ptr)
}

// Unref/Free the map.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#a98b42a54524038a2067ef2c7015c070b
func (v *VdoMap) Unref() {
	if v.Ptr != nil {
		C.g_object_unref(C.gpointer(v.Ptr))
		v.Ptr = nil
	}
}

// Removes all of the entries from this map.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-map_8h.html#ab942ab99d8733221a504807227c55c78
func (v *VdoMap) Clear() {
	C.vdo_map_clear(v.Ptr)
}

//TODOD: Implement: func (m *VdoMap) GetVariant(name string) {}

// GetByte gets a byte value by name from VdoMap.
func (m *VdoMap) GetByte(name string, def byte) byte {
	cName := clib.NewString(&name)
	defer cName.Free()
	return byte(C.vdo_map_get_byte(m.Ptr, (*C.char)(cName.Ptr), C.guchar(def)))
}

// GetBoolean gets a bool value by name from VdoMap.
func (m *VdoMap) GetBoolean(name string, def bool) bool {
	cName := clib.NewString(&name)
	defer cName.Free()
	bo := clib.GoBooleanToC(def)
	gb := C.vdo_map_get_boolean(m.Ptr, (*C.char)(cName.Ptr), C.gboolean(bo))
	return gb != C.FALSE
}

// GetInt16 gets a int16 value by name from VdoMap.
func (m *VdoMap) GetInt16(name string, def int16) int16 {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_get_int16(m.Ptr, (*C.char)(cName.Ptr), C.gint16(def))
	return int16(C.vdo_map_get_int16(m.Ptr, (*C.char)(cName.Ptr), C.gint16(def)))
}

// GetUint16 gets a uint16 value by name from VdoMap.
func (m *VdoMap) GetUint16(name string, def uint16) uint16 {
	cName := clib.NewString(&name)
	defer cName.Free()
	return uint16(C.vdo_map_get_uint16(m.Ptr, (*C.char)(cName.Ptr), C.guint16(def)))
}

// GetInt32 gets a int32 value by name from VdoMap.
func (m *VdoMap) GetInt32(name string, def int32) int32 {
	cName := clib.NewString(&name)
	defer cName.Free()
	return int32(C.vdo_map_get_int32(m.Ptr, (*C.char)(cName.Ptr), C.gint32(def)))
}

// GetUint32 gets a uint32 value by name from VdoMap.
func (m *VdoMap) GetUint32(name string, def uint32) uint32 {
	cName := clib.NewString(&name)
	defer cName.Free()
	return uint32(C.vdo_map_get_uint32(m.Ptr, (*C.char)(cName.Ptr), C.guint32(def)))
}

// GetInt64 gets a int64 value by name from VdoMap.
func (m *VdoMap) GetInt64(name string, def int64) int64 {
	cName := clib.NewString(&name)
	defer cName.Free()
	return int64(C.vdo_map_get_int64(m.Ptr, (*C.char)(cName.Ptr), C.gint64(def)))
}

// GetUint64 gets a uint64 value by name from VdoMap.
func (m *VdoMap) GetUint64(name string, def uint64) uint64 {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_get_uint64(m.Ptr, (*C.char)(cName.Ptr), C.guint64(def))
	return uint64(C.vdo_map_get_uint64(m.Ptr, (*C.char)(cName.Ptr), C.guint64(def)))
}

// GetDouble gets a double value by name from VdoMap.
func (m *VdoMap) GetDouble(name string, def float64) float64 {
	cName := clib.NewString(&name)
	defer cName.Free()
	return float64(C.vdo_map_get_double(m.Ptr, (*C.char)(cName.Ptr), C.gdouble(def)))
}

// GetString gets a string value by name from VdoMap.
func (m *VdoMap) GetString(name string, def string) string {
	cName := clib.NewString(&name)
	defer cName.Free()
	cValue := clib.NewString(&def)
	defer cValue.Free()
	var size *C.gsize
	return C.GoString(C.vdo_map_get_string(m.Ptr, (*C.char)(cName.Ptr), size, (*C.char)(cValue.Ptr)))
}

// SetByte value in a VdoMap
func (m *VdoMap) SetByte(name string, value byte) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_byte(m.Ptr, (*C.char)(cName.Ptr), C.guchar(value))
}

// Sets bool value in VdoMap
func (m *VdoMap) SetBoolean(name string, value bool) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_boolean(m.Ptr, (*C.char)(cName.Ptr), (C.gboolean)(clib.GoBooleanToC(value)))
}

// SetInt16 sets an int16 value by name in VdoMap.
func (m *VdoMap) SetInt16(name string, value int16) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_int16(m.Ptr, (*C.char)(cName.Ptr), C.gint16(value))
}

// SetUint16 sets a uint16 value by name in VdoMap.
func (m *VdoMap) SetUint16(name string, value uint16) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_uint16(m.Ptr, (*C.char)(cName.Ptr), C.guint16(value))
}

// SetInt32 sets an int32 value by name in VdoMap.
func (m *VdoMap) SetInt32(name string, value int32) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_int32(m.Ptr, (*C.char)(cName.Ptr), C.gint32(value))
}

// SetUint32 sets a uint32 value by name in VdoMap.
func (m *VdoMap) SetUint32(name string, value uint32) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_uint32(m.Ptr, (*C.char)(cName.Ptr), C.guint32(value))
}

// SetInt64 sets an int64 value by name in VdoMap.
func (m *VdoMap) SetInt64(name string, value int64) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_int64(m.Ptr, (*C.char)(cName.Ptr), C.gint64(value))
}

// SetUint64 sets a uint64 value by name in VdoMap.
func (m *VdoMap) SetUint64(name string, value uint64) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_uint64(m.Ptr, (*C.char)(cName.Ptr), C.guint64(value))
}

// SetDouble sets a double value by name in VdoMap.
func (m *VdoMap) SetDouble(name string, value float64) {
	cName := clib.NewString(&name)
	defer cName.Free()
	C.vdo_map_set_double(m.Ptr, (*C.char)(cName.Ptr), C.gdouble(value))
}

// SetString sets a string value by name in VdoMap.
func (m *VdoMap) SetString(name string, value string) {
	cName := clib.NewString(&name)
	defer cName.Free()
	cValue := clib.NewString(&value)
	defer cValue.Free()
	C.vdo_map_set_string(m.Ptr, (*C.char)(cName.Ptr), (*C.char)(cValue.Ptr))
}
