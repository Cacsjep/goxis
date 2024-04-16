package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/axlarod"
)

// PredictionResult holds the probabilities of detecting specific objects (e.g., persons, cars)
type PredictionResult struct {
	Persons float32
	Car     float32
}

// InitializeDetectionModel configures a detection model with the given model file and hardware chip.
// It sets up memory-mapped file configurations for input and output tensors.
// Returns an error if model initialization fails.
func (lea *larodExampleApplication) InitalizeDetectionModel(modelFilePath string, chipString string) error {
	model_defs := axlarod.MemMapConfiguration{
		InputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: lea.PPModel.Outputs[0].MemMapFile, // Using of ppmodel output as input for detection model
		},
		OutputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: {Size: 1}, // Output Tensor 1, since is not for ambarella-cvflow chips we could use Size: 1
			1: {Size: 1}, // Output Tensor 2, since is not for ambarella-cvflow chips we could use Size: 1
		},
	}

	if lea.DetectionModel, err = lea.app.Larod.NewInferModel(modelFilePath, chipString, model_defs); err != nil {
		return err
	}

	lea.app.AddCloseCleanFunc(func() {
		err := lea.app.Larod.DestroyModel(lea.DetectionModel)
		if err != nil {
			lea.app.Syslog.Errorf("Failed to destroy DetectionModel: %s", err.Error())
		}
	})

	return nil
}

// getDResult retrieves the detection results from the model's output tensors.
// Returns the raw byte output and an error if retrieval fails.
func (lea *larodExampleApplication) getDResult() ([]byte, error) {
	persons, err := lea.DetectionModel.Outputs[0].GetData(1)
	if err != nil {
		return nil, err
	}
	car, err := lea.DetectionModel.Outputs[1].GetData(1)
	if err != nil {
		return nil, err
	}
	output := make([]byte, 2)
	copy(output[0:1], persons)
	copy(output[1:2], car)
	return output, nil
}

// Inference executes the model and retrieves the processed results.
// It ensures the model's file pointers are correctly positioned before execution.
// Returns a JobResult containing the inference results or an error if the inference process fails.
func (lea *larodExampleApplication) Inference() (*axlarod.JobResult, error) {

	// Rewind the file position before each job.
	if err = lea.DetectionModel.Outputs[0].MemMapFile.Rewind(); err != nil {
		return nil, err
	}

	// Rewind the file position before each job.
	if err = lea.DetectionModel.Outputs[1].MemMapFile.Rewind(); err != nil {
		return nil, err
	}

	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.DetectionModel, func() error {
		return nil // is feeded via memmap
	}, func() ([]byte, error) {
		return lea.getDResult()
	}); err != nil {
		return nil, err
	}

	return result, nil
}

// InferenceOutputRead converts raw model output data into structured prediction results.
// Returns a PredictionResult or an error if data conversion fails.
//
// https://github.com/AxisCommunications/acap-native-sdk-examples/blob/7bff215e7673e4a72630bb89f04c2f7b64cf319c/vdo-larod/app/vdo_larod.c#L486C33-L486C50
func (lea *larodExampleApplication) InferenceOutputRead(result []byte) (*PredictionResult, error) {
	if len(result) < 2 { // Check that we have enough bytes to avoid panics
		return nil, fmt.Errorf("not enough data in result")
	}
	person := float32(result[0]) / 255.0 * 100
	car := float32(result[1]) / 255.0 * 100
	return &PredictionResult{Persons: person, Car: car}, nil
}
