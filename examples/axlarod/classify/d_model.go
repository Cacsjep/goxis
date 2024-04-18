package main

import (
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

	if lea.DetectionModel, err = lea.app.Larod.NewInferModel(modelFilePath, chipString, model_defs, nil); err != nil {
		return err
	}

	lea.app.AddModelCleaner(lea.DetectionModel)
	return nil
}

type InferenceResult struct {
	Person int
	Car    int
}

// getDResult retrieves the detection results from the model's output tensors.
// Returns the raw byte output and an error if retrieval fails.
func (lea *larodExampleApplication) getDResult() (*InferenceResult, error) {
	person, err := lea.DetectionModel.Outputs[0].GetDataAsInt()
	if err != nil {
		return nil, err
	}
	car, err := lea.DetectionModel.Outputs[1].GetDataAsInt()
	if err != nil {
		return nil, err
	}
	return &InferenceResult{Person: person, Car: car}, nil
}

// Inference executes the model and retrieves the processed results.
// It ensures the model's file pointers are correctly positioned before execution.
// Returns a JobResult containing the inference results or an error if the inference process fails.
func (lea *larodExampleApplication) Inference() (*axlarod.JobResult, error) {

	// Rewind the file position before each job.
	if err = lea.DetectionModel.RewindAllOutputsMemMapFiles(); err != nil {
		return nil, err
	}

	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.DetectionModel, func() error {
		return nil // is feeded via memmap
	}, func() (any, error) {
		return lea.getDResult()
	}); err != nil {
		return nil, err
	}

	return result, nil
}

// InferenceOutputRead converts raw model output data into structured prediction results.
//
// https://github.com/AxisCommunications/acap-native-sdk-examples/blob/7bff215e7673e4a72630bb89f04c2f7b64cf319c/vdo-larod/app/vdo_larod.c#L486C33-L486C50
func (lea *larodExampleApplication) InferenceOutputRead(result *InferenceResult) (*PredictionResult, error) {
	return &PredictionResult{Persons: float32(result.Person) / 255.0 * 100, Car: float32(result.Car) / 255.0 * 100}, nil
}
