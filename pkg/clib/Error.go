package clib

/*
#cgo pkg-config: glib-2.0
#include <stdlib.h>
#include <glib.h>
*/
import "C"
import (
	"errors"
)

// Error wraps a GLib GError to manage its lifecycle and provide Go-style error handling.
type Error struct {
	Ptr *C.GError
}

// NewError creates a new, empty Error instance.
// This is useful when a C function requires a GError** argument to report errors.
func NewError() *Error {
	return &Error{}
}

// IsNil checks if the underlying GError is nil, indicating no error.
func (e *Error) IsNil() bool {
	return e.Ptr == nil
}

// Code returns the error code contained in the GError.
// It returns 0 if the GError is nil.
func (e *Error) Code() int {
	if !e.IsNil() {
		return int(e.Ptr.code)
	}
	return 0
}

// Message returns the error message contained in the GError.
// It returns an empty string if the GError is nil.
func (e *Error) Message() string {
	if !e.IsNil() {
		return C.GoString(e.Ptr.message)
	}
	return ""
}

// IsError checks if the GError contains an error, converts it to a Go error, and frees the GError.
// It returns nil if there is no error.
func (e *Error) IsError() error {
	if !e.IsNil() {
		defer e.Free()
		return errors.New(e.Message())
	}
	return nil
}

func (e *Error) IsErrorReturnCode() (error, int) {
	if !e.IsNil() {
		defer e.Free()
		return errors.New(e.Message()), e.Code()
	}
	return nil, 0
}

// IsErrorOrNotSuccess checks if the GError contains an error, converts it to a Go error, and frees the GError.
// It returns nil if there is no error. also success value is check to have a easy return method for c funcs from axis
func (e *Error) IsErrorOrNotSuccess(success int, not_success_msg string) error {
	if !e.IsNil() {
		defer e.Free()
		return errors.New(e.Message())
	}
	if success == 0 {
		return errors.New(not_success_msg)
	}
	return nil
}

// Free releases the memory allocated for the GError.
// It is safe to call multiple times, but the Error should not be used after being freed.
func (e *Error) Free() {
	if !e.IsNil() {
		C.g_error_free(e.Ptr)
		e.Ptr = nil
	}
}
