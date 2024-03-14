// Package glib provides Go bindings for GLib, a low-level core library that forms the basis
// for projects such as GTK+ and GNOME. It provides data structure handling for C, portability
// wrappers, and interfaces for such runtime functionality as an event loop, threads, dynamic
// loading, and an object system.
//
// This particular implementation focuses on the GMainLoop functionality, enabling the creation,
// control, and cleanup of main event loops within Go applications using GLib.
//
// The use of this package requires the GLib library to be installed on your system.
package glib

/*
#cgo pkg-config: glib-2.0
#include <glib.h>
#include <stdio.h>
*/
import "C"
import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// GMainLoop is a representation of GLib's GMainLoop structure, encapsulating the functionality
// of a main event loop. A GMainLoop is a loop that processes events for a context.
type GMainLoop struct {
	Ptr *C.GMainLoop
}

// NewMainLoop creates and returns a new GMainLoop. This function creates a new event loop
// object with a reference count of 1. The loop is set to not automatically quit upon the
// termination of its last source, allowing for manual control.
//
// https://docs.gtk.org/glib/ctor.MainLoop.new.html
//
// Returns:
//   - *GMainLoop: A pointer to the newly created GMainLoop.
func NewMainLoop() *GMainLoop {
	return &GMainLoop{Ptr: C.g_main_loop_new(nil, C.FALSE)}
}

// Run begins the execution of the GMainLoop. This function locks the current OS thread,
// runs the event loop, and processes events for the context until Quit is called. It's important
// to manage concurrency and ensure that Start is called from the appropriate execution context,
// as it will lock the calling thread.
//
// Run also register a signal Handler for SIGTERM and SIGINT

// https://docs.gtk.org/glib/method.MainLoop.run.html
func (g *GMainLoop) Run() {
	SignalHandler(g.Quit)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	C.g_main_loop_run(g.Ptr)
}

// Quit stops the GMainLoop, causing Start to return. It is safe to call this function from any
// thread. The function decreases the reference count of the GMainLoop object and frees it if
// this was the last reference. After calling Quit, the GMainLoop should not be used anymore.
//
// https://docs.gtk.org/glib/method.MainLoop.quit.html
func (g *GMainLoop) Quit() {
	C.g_main_loop_quit(g.Ptr)
	C.g_main_loop_unref(g.Ptr)
}

func SignalHandler(handler func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigs
		handler()
		os.Exit(0)
	}()
}
