package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import (
	"fmt"
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
	ptr **C.larodTensor // Pointer to an array of tensor pointers
}

type LarodModelIO struct {
	Inputs        *LarodTensor
	InputsCount   uint
	InputPitches  *LarodTensorPitches
	Outputs       *LarodTensor
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

func (model *LarodModel) CreateModelInputs() (*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelInputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	return &LarodTensor{ptr: tensors}, uint(numTensors), nil
}

func (model *LarodModel) CreateModelOutputs() (*LarodTensor, uint, error) {
	var numTensors C.size_t
	var cError *C.larodError
	tensors := C.larodCreateModelOutputs(model.ptr, &numTensors, &cError)
	if cError != nil {
		return nil, 0, newLarodError(cError)
	}
	return &LarodTensor{ptr: tensors}, uint(numTensors), nil
}

func (model *LarodModel) CreateModelTensors() (*LarodModelIO, error) {
	inputs, inputsCount, err := model.CreateModelInputs()
	if err != nil {
		return nil, fmt.Errorf("failed to create model inputs: %w", err)
	}

	outputs, outputsCount, err := model.CreateModelOutputs()
	if err != nil {
		return nil, fmt.Errorf("failed to create model outputs: %w", err)
	}

	inputsPitches, err := inputs.GetTensorPitches()
	if err != nil {
		return nil, fmt.Errorf("failed to get input tensor pitches: %w", err)
	}

	outputsPitches, err := outputs.GetTensorPitches()
	if err != nil {
		return nil, fmt.Errorf("failed to get output tensor pitches: %w", err)
	}

	return &LarodModelIO{
		Inputs:        inputs,
		InputsCount:   inputsCount,
		InputPitches:  inputsPitches,
		Outputs:       outputs,
		OutputsCount:  outputsCount,
		OutputPitches: outputsPitches,
	}, nil
}

// GetTensorPitches retrieves the pitch information of a tensor and converts it to a Go struct.
func (tensor *LarodTensor) GetTensorPitches() (*LarodTensorPitches, error) {
	var cError *C.larodError
	cPitches := C.larodGetTensorPitches(*tensor.ptr, &cError)
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
