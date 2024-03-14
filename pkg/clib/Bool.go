package clib

/*
#include <glib.h> // Make sure to include glib.h for gboolean definition
#include <stdio.h>
*/
import "C"

// Bool wraps a GLib gboolean type, offering a Go representation to facilitate
// the use of boolean values in GLib functions and callbacks.
type Bool struct {
	Ptr C.gboolean // Ptr holds the GLib gboolean value.
}

// NewBool initializes and returns a new Bool instance. This function is useful
// for creating a Bool to interact with GLib functions or when starting with a
// new GLib boolean value.
func NewBool() *Bool {
	return &Bool{}
}

// ToGolang converts the GLib gboolean to a native Go bool. This method is crucial
// for translating boolean values from C to Go, enabling Go programs to accurately
// work with the boolean values returned by or used in GLib functions.
func (cs *Bool) ToGolang() bool {
	return CtoGoBoolean(cs.Ptr)
}

// CtoGoBoolean takes a C.gboolean and converts it to a Go bool. Since GLib
// gboolean is typically defined as an integer type where FALSE is 0 and TRUE
// is anything non-zero, this function ensures accurate representation of these
// values in Go's boolean type.
//
// Parameters:
//  - cBool: The C.gboolean value to convert.
//
// Returns:
//  - bool: The converted Go bool value, true if cBool is not FALSE, false otherwise.
func CtoGoBoolean(cBool C.gboolean) bool {
	return cBool != C.FALSE
}

// GoBooleanToC takes a Go bool and converts it to a GLib gboolean. This is
// the inverse operation of CtoGoBoolean, enabling Go boolean values to be
// accurately represented and used in GLib functions that expect gboolean values.
//
// Parameters:
//  - b: The Go bool value to convert.
//
// Returns:
//  - C.gboolean: The converted GLib gboolean value, 1 (TRUE) if b is true, 0 (FALSE) otherwise.
func GoBooleanToC(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}
