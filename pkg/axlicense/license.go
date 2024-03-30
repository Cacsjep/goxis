/*
Package axlicense provides a Go wrapper around the Axis License Key library, facilitating license verification and management for ACAP applications. This package simplifies the process of validating application licenses, retrieving license expiration dates, and understanding license states, leveraging the license key system implemented in Axis devices.

The Axis License Key system is designed to ensure that ACAP applications running on Axis devices are properly licensed. It verifies licenses based on application-specific details such as name, ID, and version numbers, provided during the application's deployment. This system helps protect against unauthorized distribution and use of applications.

Requirements:

This package requires linking against both the static and dynamic parts of the License Key library to ensure secure license verification. The appropriate linker flags (`-Wl,-Bstatic -llicensekey_stat -Wl,-Bdynamic -llicensekey -ldl`) must be specified in the Go package for correct compilation and linkage.

Anti-Debugging Measures:

To safeguard against debugging attempts aimed at circumventing license checks, it is recommended to incorporate anti-debugging techniques such as `ptrace` within the application's initialization code. This package, however, focuses on the licensing aspect and does not implement such measures directly.

Usage Notes:

- The application name passed to the license verification functions must match the `APPNAME` specified in the application's `package.conf`.
- The application ID is assigned by Axis and must match the `APPID` also specified in `package.conf`.
- Major and minor version numbers should correspond to those declared in the application's configuration, aligning with `APPMAJORVERSION` and `APPMINORVERSION`.

This wrapper aims to streamline the process of working with Axis License Keys for Go developers, abstracting the complexities of direct C library interactions and providing a more Go-idiomatic approach to license management for ACAP applications.
*/
package axlicense

/*
#cgo LDFLAGS: -Wl,-Bstatic -llicensekey_stat -Wl,-Bdynamic -llicensekey -ldl
#include <licensekey.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/licensekey/html/licensekey_8h.html
// License requires to enable them in mainfest, and an APPID registerd at AXIS to create Licenses for the application
//
//	"copyProtection": {
//		"method": "axis"
//	}
//
// Portal Link: https://www.axis.com/partner_pages/compatible_applications/
func LicensekeyVerify(app_name string, app_id int, major_version int, minor_version int) (valid bool) {
	cAppName := C.CString(app_name)
	defer C.free(unsafe.Pointer(cAppName))

	state := C.licensekey_verify(cAppName, C.int(app_id), C.int(major_version), C.int(minor_version))
	if int(state) == 1 {
		return true
	}
	return false
}

// TODO: Bring this to work
func LicensekeyGetExpDate(app_name string) (time.Time, error) {
	cAppName := C.CString(app_name)
	defer C.free(unsafe.Pointer(cAppName))

	str_date := C.licensekey_get_exp_date(cAppName, nil)
	if str_date == nil {
		return time.Now(), errors.New("the expiration date couldn't be read")
	}
	go_str_date := C.GoString(str_date)
	if go_str_date == "0" {
		return time.Now(), errors.New("the expiration date is '0'")
	}
	date, err := time.Parse("2006-01-02", C.GoString(str_date))
	if err != nil {
		return time.Now(), err
	}
	return date, err
}
