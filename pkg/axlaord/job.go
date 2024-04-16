package axlarod

/*
#cgo pkg-config: liblarod
#include "larod.h"
*/
import "C"
import (
	"fmt"
	"time"
)

type JobResult struct {
	ExecutionTime float64
	OutputData    []byte
}

type JobRequest struct {
	ptr *C.larodJobRequest
}

func (model *LarodModel) Execute(conn *LarodConnection) error {
	if model.Job == nil {
		return fmt.Errorf("JobRequest is nil")
	}
	var cError *C.larodError
	result := C.larodRunJob(conn.ptr, model.Job.ptr, &cError)
	if result == C.bool(false) {
		if cError != nil {
			errorMsg := C.GoString(cError.msg)
			C.larodClearError(&cError)
			return fmt.Errorf("larodRunJob failed: %s", errorMsg)
		}
		return fmt.Errorf("larodRunJob failed with unknown error")
	}

	return nil
}

// ExecuteJob orchestrates the execution of a model with data setting and getting
func (l *Larod) ExecuteJob(model *LarodModel, dataSetterFunc func() error, dataGetterFunc func() ([]byte, error)) (*JobResult, error) {
	start := time.Now()
	if err := dataSetterFunc(); err != nil {
		return nil, fmt.Errorf("error setting data: %w", err)
	}

	// Execute the model
	if err := model.Execute(l.conn); err != nil {
		return nil, fmt.Errorf("model execution failed: %w", err)
	}

	// Get data
	outputData, err := dataGetterFunc()
	if err != nil {
		return nil, fmt.Errorf("error getting data: %w", err)
	}
	return &JobResult{OutputData: outputData, ExecutionTime: float64(time.Since(start).Milliseconds())}, nil
}

// Wrap the larodCreateJobRequest C function.
func (model *LarodModel) CreateJobRequest(inputTensors []*LarodTensor, outputTensors []*LarodTensor, params *LarodMap) (*JobRequest, error) {
	numInputs := C.size_t(len(inputTensors))
	numOutputs := C.size_t(len(outputTensors))

	cInputTensors := make([]*C.larodTensor, numInputs)
	for i, tensor := range inputTensors {
		cInputTensors[i] = tensor.ptr
	}

	cOutputTensors := make([]*C.larodTensor, numOutputs)
	for i, tensor := range outputTensors {
		cOutputTensors[i] = tensor.ptr
	}

	var cError *C.larodError

	var cParams *C.larodMap
	if params != nil {
		cParams = params.ptr
	} else {
		cParams = nil
	}

	job := C.larodCreateJobRequest(model.ptr, &cInputTensors[0], numInputs, &cOutputTensors[0], numOutputs, cParams, &cError)

	if job == nil {
		if cError != nil {
			return nil, fmt.Errorf("larodCreateJobRequest failed: %s", C.GoString(cError.msg))
		}
		return nil, fmt.Errorf("larodCreateJobRequest failed with unknown error")
	}

	model.Job = &JobRequest{ptr: job}
	return model.Job, nil
}
