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

func (dev *LarodDevice) GetName() (string, error) {
	var cError *C.larodError
	cName := C.larodGetDeviceName(dev.ptr, &cError)
	if cName == nil {
		return "", newLarodError(cError)
	}
	return C.GoString(cName), nil
}

func (l *Larod) GetDeviceByName(name string) (*LarodDevice, error) {
	for _, device := range l.Devices {
		if device.Name == name {
			return device, nil
		}
	}
	return nil, fmt.Errorf("device not found")
}

func (l *Larod) ListDevices() ([]*LarodDevice, error) {
	var numDevices C.size_t
	var cError *C.larodError

	// Call the C function
	cDevices := C.larodListDevices(l.conn.ptr, &numDevices, &cError)
	if cDevices == nil {
		return nil, newLarodError(cError)
	}
	// Note: Do not free cDevices. Its lifetime is managed by the connection.

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
