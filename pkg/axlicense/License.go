package axlicense

/*
#cgo LDFLAGS: -Wl,-Bstatic -llicensekey_stat -Wl,-Bdynamic -llicensekey -ldl
#include <licensekey.h>
*/
import "C"
import (
	"errors"
	"time"

	"github.com/Cacsjep/goxis/pkg/clib"
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
	cAppName := clib.NewString(&app_name)
	defer cAppName.Free()
	state := C.licensekey_verify((*C.char)(cAppName.Ptr), C.int(app_id), C.int(major_version), C.int(minor_version))
	if int(state) == 1 {
		return true
	}
	return false
}

// TODO: Bring this to work
func LicensekeyGetExpDate(app_name string) (time.Time, error) {
	cAppName := clib.NewString(&app_name)
	defer cAppName.Free()
	str_date := C.licensekey_get_exp_date((*C.char)(cAppName.Ptr), nil)
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
