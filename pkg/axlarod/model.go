package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
#include <sys/mman.h>

*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

type LarodAccess int

const (
	LarodAccessInvalid LarodAccess = iota
	LarodAccessPrivate
	LarodAccessPublic
)

type LarodModel struct {
	ptr              *C.larodModel
	inputTensorsPtr  **C.larodTensor
	outputTensorsPtr **C.larodTensor
	maps             []*LarodMap
	Fd               uintptr
	Job              *JobRequest
	Name             string
	Inputs           []*LarodTensor
	InputsCount      uint
	InputPitches     *LarodTensorPitches
	Outputs          []*LarodTensor
	OutputsCount     uint
	OutputPitches    *LarodTensorPitches
}

// LarodLoadModel loads a new model onto a specified device.
func (l *Larod) LoadModel(file_path *string, dev *LarodDevice, access LarodAccess, name string, params *LarodMap) (*LarodModel, error) {
	_fd := C.int(-1)

	if file_path != nil {
		model_file, err := os.Open(*file_path)
		if err != nil {
			return nil, err
		}
		_fd = C.int(model_file.Fd())
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cError *C.larodError
	var cParams *C.larodMap
	if params != nil {
		cParams = params.ptr
	} else {
		cParams = nil
	}

	cModel := C.larodLoadModel(l.conn.ptr, _fd, dev.ptr, C.larodAccess(access), cName, cParams, &cError)
	if cModel == nil {
		if cError != nil {
			return nil, newLarodError(cError)
		}
		return nil, fmt.Errorf("larodLoadModel failed without setting an error")
	}

	maps := make([]*LarodMap, 0)
	maps = append(maps, params)

	model := &LarodModel{ptr: cModel, maps: maps, Name: name}
	return model, nil
}

// LoadModelWithDeviceName loads a new model onto a specified device by name.
func (l *Larod) LoadModelWithDeviceName(file_path *string, dev_name string, access LarodAccess, name string, params *LarodMap) (*LarodModel, error) {
	device, err := l.GetDeviceByName(dev_name)
	if err != nil {
		return nil, err
	}
	return l.LoadModel(file_path, device, access, name, params)
}

// Seek the memory mapped file to the beginning for all output tensors.
func (m *LarodModel) RewindAllOutputsMemMapFiles() error {
	for _, tensor_ouput := range m.Outputs {
		if err := tensor_ouput.MemMapFile.Rewind(); err != nil {
			return err
		}
	}
	return nil
}

// LoadModelWithDeviceID loads a new model onto a specified device by ID.
func (m *LarodModel) Destroy() {
	if m == nil {
		return
	}
	if m.ptr == nil {
		return
	}
	C.larodDestroyModel(&m.ptr)
}
