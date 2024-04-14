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

type LarodAccess int

const (
	LarodAccessInvalid LarodAccess = iota
	LarodAccessPrivate
	LarodAccessPublic
)

type LarodModel struct {
	ptr *C.larodModel
}

// LarodLoadModel loads a new model onto a specified device.
func (l *Larod) LoadModel(fd *uintptr, dev *LarodDevice, access LarodAccess, name string, params *LarodMap) (*LarodModel, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cError *C.larodError
	var cParams *C.larodMap
	if params != nil {
		cParams = params.ptr
	} else {
		cParams = nil
	}

	_fd := C.int(-1)
	if fd != nil {
		_fd = C.int(*fd)
	}

	cModel := C.larodLoadModel(l.conn.ptr, _fd, dev.ptr, C.larodAccess(access), cName, cParams, &cError)
	if cModel == nil {
		if cError != nil {
			return nil, newLarodError(cError)
		}
		return nil, fmt.Errorf("larodLoadModel failed without setting an error")
	}
	return &LarodModel{ptr: cModel}, nil
}

func (l *Larod) LoadModelWithDeviceName(fd *uintptr, dev_name string, access LarodAccess, name string, params *LarodMap) (*LarodModel, error) {
	device, err := l.GetDeviceByName(dev_name)
	if err != nil {
		return nil, err
	}
	return l.LoadModel(fd, device, access, name, params)
}
