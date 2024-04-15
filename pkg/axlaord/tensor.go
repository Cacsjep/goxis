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

const LAROD_TENSOR_MAX_LEN = 12

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

type LarodTensorLayout int

const (
	LarodTensorLayoutInvalid LarodTensorLayout = iota
	LarodTensorLayoutUnspecified
	LarodTensorLayoutNHWC
	LarodTensorLayoutNCHW
	LarodTensorLayout420SP
	LarodTensorLayoutMax
)

type LarodTensor struct {
	ptr     *C.larodTensor
	TmpFile *TmpFile
}

type LarodModelIO struct {
	Inputs        []*LarodTensor
	InputsCount   uint
	InputPitches  *LarodTensorPitches
	Outputs       []*LarodTensor
	OutputsCount  uint
	OutputPitches *LarodTensorPitches
}

type LarodTensorPitches struct {
	Pitches [LAROD_TENSOR_MAX_LEN]uint
	Len     uint
}

func (lmt *LarodModelIO) String() string {
	return fmt.Sprintf("LarodModelIO{InputsCount: %d, InputPitches %d, OutputsCount: %d, OutputPitches: %d}", lmt.InputsCount, lmt.InputPitches.Pitches, lmt.OutputsCount, lmt.OutputPitches.Pitches)
}

func (model *LarodModel) CreateModelInputs() ([]*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelInputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	model.inputTensorPtr = tensors
	length := uint(numTensors)
	result := make([]*LarodTensor, length)
	for i := 0; i < int(length); i++ {
		tensorPtr := (**C.larodTensor)(unsafe.Pointer(uintptr(unsafe.Pointer(tensors)) + uintptr(i)*unsafe.Sizeof(*tensors)))
		result[i] = &LarodTensor{ptr: *tensorPtr}
	}
	return result, length, nil
}

func (model *LarodModel) CreateModelOutputs() ([]*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelOutputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	model.outputTensorPtr = tensors
	length := uint(numTensors)
	result := make([]*LarodTensor, length)
	for i := 0; i < int(length); i++ {
		tensorPtr := (**C.larodTensor)(unsafe.Pointer(uintptr(unsafe.Pointer(tensors)) + uintptr(i)*unsafe.Sizeof(*tensors)))
		result[i] = &LarodTensor{ptr: *tensorPtr}
	}
	return result, length, nil
}

func (model *LarodModel) CreateModelTensors(model_defs *ModelTmpMapDefiniton) (*LarodModelIO, error) {
	inputs, inputsCount, err := model.CreateModelInputs()
	if err != nil {
		return nil, fmt.Errorf("failed to create model inputs: %w", err)
	}

	outputs, outputsCount, err := model.CreateModelOutputs()
	if err != nil {
		return nil, fmt.Errorf("failed to create model outputs: %w", err)
	}

	inputsPitches, err := inputs[0].GetTensorPitches()
	if err != nil {
		return nil, fmt.Errorf("failed to get input tensor pitches: %w", err)
	}

	outputsPitches, err := outputs[0].GetTensorPitches()
	if err != nil {
		return nil, fmt.Errorf("failed to get output tensor pitches: %w", err)
	}

	model.LarodModelIO = &LarodModelIO{
		Inputs:        inputs,
		InputsCount:   inputsCount,
		InputPitches:  inputsPitches,
		Outputs:       outputs,
		OutputsCount:  outputsCount,
		OutputPitches: outputsPitches,
	}

	err = model.MapModelTmpFiles(model_defs)
	if err != nil {
		return nil, fmt.Errorf("failed to map model tmp files: %w", err)
	}

	return model.LarodModelIO, nil
}

// GetTensorPitches retrieves the pitch information of a tensor and converts it to a Go struct.
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

func (tensor *LarodTensor) CopyDataInto(data []byte) error {
	return CopyDataToMappedMemory(tensor.TmpFile.MemoryAddress, data)
}

func (tensor *LarodTensor) GetData(size int) ([]byte, error) {
	return CopyDataFromMappedMemory(tensor.TmpFile.MemoryAddress, size)
}
