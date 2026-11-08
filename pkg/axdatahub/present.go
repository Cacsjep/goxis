package axdatahub

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrUnavailable means the camera has no Device Data Hub client library
// (AXIS OS 12.9 and earlier). It never changes at runtime.
var ErrUnavailable = errors.New("the Device Data Hub client (libdevice-data-hub-client-c) is not available on this device")

// libraryName is the client library the package loads.
const libraryName = "libdevice-data-hub-client-c.so.0"

// libraryDirs are the directories Present searches, after LD_LIBRARY_PATH.
var libraryDirs = []string{"/usr/lib", "/lib", "/usr/lib64", "/lib64", "/usr/local/lib"}

// Present reports, without loading it, whether the client library is on this
// device: nil when the file exists, else an error wrapping ErrUnavailable.
// Loading pulls in the library's C++ dependencies, so an application that only
// keeps the hub as a fallback can check with Present and call Load when it
// needs the hub. Load remains the authoritative check.
func Present() error {
	var dirs []string
	for _, dir := range strings.Split(os.Getenv("LD_LIBRARY_PATH"), ":") {
		if dir != "" {
			dirs = append(dirs, dir)
		}
	}
	dirs = append(dirs, libraryDirs...)
	for _, dir := range dirs {
		if _, err := os.Stat(filepath.Join(dir, libraryName)); err == nil {
			return nil
		}
	}
	return fmt.Errorf("%w: %s not found in %s", ErrUnavailable, libraryName, strings.Join(dirs, ":"))
}
