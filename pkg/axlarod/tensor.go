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

const LAROD_TENSOR_MAX_LEN = 12 // Maximum length of the tensor.

// LarodTensorDataType defines different data types that a tensor can represent.
type LarodTensorDataType int

const (
	LarodTensorDataTypeInvalid LarodTensorDataType = iota
	LarodTensorDataTypeUnspecified
	LarodTensorDataTypeBool
	LarodTensorDataTypeUint8
	LarodTensorDataTypeInt8
	LarodTensorDataTypeUint16
	LarodTensorDataTypeInt16
	LarodTensorDataTypeUint32
	LarodTensorDataTypeInt32
	LarodTensorDataTypeUint64
	LarodTensorDataTypeInt64
	LarodTensorDataTypeFloat16
	LarodTensorDataTypeFloat32
	LarodTensorDataTypeFloat64
	LarodTensorDataTypeMax
)

// LarodTensorLayout represents the memory layout of a tensor.
type LarodTensorLayout int

const (
	LarodTensorLayoutInvalid LarodTensorLayout = iota
	LarodTensorLayoutUnspecified
	LarodTensorLayoutNHWC  // NHWC layout: channels last
	LarodTensorLayoutNCHW  // NCHW layout: channels first
	LarodTensorLayout420SP // 420SP layout: semi-planar YCbCr format
	LarodTensorLayoutMax
)

// LarodTensor encapsulates a tensor structure with a pointer to its C representation.
type LarodTensor struct {
	ptr        *C.larodTensor
	MemMapFile *MemMapFile
}

// LarodTensorPitches represents the memory layout pitches of a tensor.
type LarodTensorPitches struct {
	Pitches [LAROD_TENSOR_MAX_LEN]uint
	Len     uint
}

// CreateModelInputs allocates and returns pointers to input tensors for a model,
// along with the count of these tensors.
func (model *LarodModel) CreateModelInputs() ([]*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelInputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	model.inputTensorsPtr = tensors
	length := uint(numTensors)
	result := make([]*LarodTensor, length)
	for i := 0; i < int(length); i++ {
		tensorPtr := (**C.larodTensor)(unsafe.Pointer(uintptr(unsafe.Pointer(tensors)) + uintptr(i)*unsafe.Sizeof(*tensors)))
		result[i] = &LarodTensor{ptr: *tensorPtr}
	}
	return result, length, nil
}

// CreateModelOutputs allocates and returns pointers to output tensors for a model,
// along with the count of these tensors.
func (model *LarodModel) CreateModelOutputs() ([]*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelOutputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	model.outputTensorsPtr = tensors
	length := uint(numTensors)
	result := make([]*LarodTensor, length)
	for i := 0; i < int(length); i++ {
		tensorPtr := (**C.larodTensor)(unsafe.Pointer(uintptr(unsafe.Pointer(tensors)) + uintptr(i)*unsafe.Sizeof(*tensors)))
		result[i] = &LarodTensor{ptr: *tensorPtr}
	}
	return result, length, nil
}

// CreateModelTensors initializes and configures model tensors according to provided memory configurations.
func (model *LarodModel) CreateModelTensors(model_defs *MemMapConfiguration) error {
	inputs, inputsCount, err := model.CreateModelInputs()
	if err != nil {
		return fmt.Errorf("failed to create model inputs: %w", err)
	}

	outputs, outputsCount, err := model.CreateModelOutputs()
	if err != nil {
		return fmt.Errorf("failed to create model outputs: %w", err)
	}

	inputsPitches, err := inputs[0].GetTensorPitches()
	if err != nil {
		return fmt.Errorf("failed to get input tensor pitches: %w", err)
	}

	outputsPitches, err := outputs[0].GetTensorPitches()
	if err != nil {
		return fmt.Errorf("failed to get output tensor pitches: %w", err)
	}

	model.Inputs = inputs
	model.InputsCount = inputsCount
	model.InputPitches = inputsPitches
	model.Outputs = outputs
	model.OutputsCount = outputsCount
	model.OutputPitches = outputsPitches

	err = model.MapModelTmpFiles(model_defs)
	if err != nil {
		return fmt.Errorf("failed to map model tmp files: %w", err)
	}

	return nil
}

// GetTensorPitches retrieves the pitch information of a tensor and converts it to a Go struct.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#aaade104831363ee594d9b1bb49583fe5
func (tensor *LarodTensor) GetTensorPitches() (*LarodTensorPitches, error) {
	var cError *C.larodError
	cPitches := C.larodGetTensorPitches(tensor.ptr, &cError)
	if cPitches == nil {
		if cError != nil {
			return nil, newLarodError(cError)
		}
		return nil, fmt.Errorf("failed to get tensor pitches without a specific error")
	}

	// Convert C array and len into Go equivalent
	goPitches := &LarodTensorPitches{
		Len: uint(cPitches.len),
	}
	for i := 0; i < int(goPitches.Len); i++ {
		goPitches.Pitches[i] = uint(cPitches.pitches[i])
	}

	return goPitches, nil
}

// SetTensorFd sets the file descriptor of a tensor.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#a4423cb2a9da02f3a21766a2713edbaa5
func (tensor *LarodTensor) SetTensorFd(fd uintptr) error {
	var cError *C.larodError
	result := C.larodSetTensorFd(tensor.ptr, C.int(fd), &cError)
	if result == C.bool(false) {
		if cError != nil {
			return newLarodError(cError)
		}
		return fmt.Errorf("failed to set tensor file descriptor without a specific error")
	}
	return nil
}

// Get data type of a tensor.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#a07bd3d2345f650c3b46704284c728a3b
func (tensor *LarodTensor) GetDataType() (LarodTensorDataType, error) {
	var cError *C.larodError
	result := LarodTensorDataType(C.larodGetTensorDataType(tensor.ptr, &cError)) // Convert C type to Go type here

	if result == LarodTensorDataTypeInvalid {
		if cError != nil {
			return LarodTensorDataTypeInvalid, newLarodError(cError)
		}
		return LarodTensorDataTypeInvalid, fmt.Errorf("failed to get tensor data type without a specific error")
	}

	return result, nil
}

// CopyDataInto copies data into the memory mapped file associated with a tensor.
func (tensor *LarodTensor) CopyDataInto(data []byte) error {
	return CopyDataToMappedMemory(tensor.MemMapFile.MemoryAddress, data)
}

// GetData retrieves data from the memory mapped file associated with a tensor.
func (tensor *LarodTensor) GetData(size int) ([]byte, error) {
	return CopyDataFromMappedMemory(tensor.MemMapFile.MemoryAddress, size)
}
