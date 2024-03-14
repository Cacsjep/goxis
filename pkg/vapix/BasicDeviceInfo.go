package vapix

import (
	"errors"

	"github.com/Cacsjep/goxis/pkg/shared"
)

func VapixBasicDeviceInfo(creds *shared.VapixCreds) (map[string]interface{}, error) {
	basicDeviceInfoMethod := NewVapixBaseMethodCall("getAllProperties")
	r := VapixPost(creds, InternalUrlConstruct("/axis-cgi/basicdeviceinfo.cgi"), basicDeviceInfoMethod)
	if r.IsOk {
		defer r.ResponseReader.Close()
		if response, err := JsonResponseParser(r.ResponseReader); err != nil {
			return nil, err
		} else {
			if response.Data.PropertiesList == nil {
				return nil, errors.New("PropertiesList are nil")
			}
			return response.Data.PropertiesList, nil
		}
	}
	return nil, r.Error
}
