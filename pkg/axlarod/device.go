package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type LarodDevice struct {
	ptr  *C.larodDevice
	Name string
}

// GetName retrieves the name of a LarodDevice.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#ab0279d8c0983a66e6c4a6476f2be48cd
func (dev *LarodDevice) GetName() (string, error) {
	var cError *C.larodError
	cName := C.larodGetDeviceName(dev.ptr, &cError)
	if cName == nil {
		return "", newLarodError(cError)
	}
	return C.GoString(cName), nil
}

// GetDeviceByName searches for a device by name among the cached devices in a Larod instance.
func (l *Larod) GetDeviceByName(name string) (*LarodDevice, error) {
	for _, device := range l.Devices {
		if device.Name == name {
			return device, nil
		}
	}
	return nil, fmt.Errorf("device not found")
}

// ListDevices queries and lists all devices managed by a Larod instance's connection.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#a23e4b862c843907b8c25164df07b22ed
func (l *Larod) ListDevices() ([]*LarodDevice, error) {
	var numDevices C.size_t
	var cError *C.larodError

	cDevices := C.larodListDevices(l.conn.ptr, &numDevices, &cError)
	if cDevices == nil {
		return nil, newLarodError(cError)
	}

	length := int(numDevices)
	devices := make([]*LarodDevice, length)

	maxSize := 1 << 20
	if length > maxSize {
		return nil, fmt.Errorf("device count exceeds maximum safe array size")
	}
	tmpSlice := (*[1 << 20]*C.larodDevice)(unsafe.Pointer(cDevices))[:length:length]
	for i, cDev := range tmpSlice {
		devices[i] = &LarodDevice{ptr: cDev}
		name, err := devices[i].GetName()
		if err != nil {
			return nil, fmt.Errorf("failed to get device name: %w", err)
		}
		devices[i].Name = name
	}
	l.Devices = devices
	return devices, nil
}
