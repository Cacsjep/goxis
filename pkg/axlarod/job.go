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

// Execute runs a job using a LarodModel on a given LarodConnection.
// It submits the job described by model.Job for execution.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#a7c492dcfb18e0a32407dc6078bd50dbe
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

// ExecuteJob orchestrates the execution of a model with data setting and retrieving processes.
// It measures the execution time and returns a JobResult with the execution time and the output data.
// If any step (setting data, executing the model, or retrieving data) fails, it returns an error detailing the failure.
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

// CreateJobRequest initializes a job request for a model using specified input and output tensors and optional parameters.
//
// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/larod/html/larod_8h.html#af25e0293b11b12fb8e296fb20c828159
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

// DestroyJobRequest cleans up resources associated with a JobRequest.
func (job *JobRequest) Destroy() {
	C.larodDestroyJobRequest(&job.ptr)
}
