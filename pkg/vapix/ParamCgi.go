package vapix

import (
	"errors"

	"github.com/Cacsjep/goxis/pkg/shared"
)

// Get all params via parameter cgi API
// useInternalVapix -> set to true if it should called via internal vapix call via ACAP
// requires dbus configuration in manifest.
// When not useInternalVapix is used a
// host_location need to be set like http://192.168.0.90
func VapixParamCgiGetAll(creds *shared.VapixCreds, useInternalVapix bool, host_location *string) (map[string]string, error) {
	cgi_path := "/axis-cgi/param.cgi?action=list"
	var url string
	if useInternalVapix {
		url = InternalUrlConstruct(cgi_path)
	} else {
		if host_location == nil {
			return nil, errors.New("host_location is nil")
		}
		hl := *host_location
		url = hl + cgi_path
	}
	paramsResult := VapixGet(creds, url)
	if paramsResult.IsOk {
		defer paramsResult.ResponseReader.Close()
		return ParseKeyValueRequestBody(paramsResult.ResponseReader)
	} else {
		return nil, paramsResult.Error
	}
}
