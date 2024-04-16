package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
)

func (lea *larodExampleApplication) InitalizeDetectionModel(modelFilePath string, chipString string) error {
	model_defs := axlarod.MemMapConfiguration{
		InputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: lea.PPModel.Outputs[0].MemMapFile, // Input Tensor 0
		},
		OutputTmpMapFiles: map[int]*axlarod.MemMapFile{
			0: {Size: 4}, // Output Tensor 1
			1: {Size: 4}, // Output Tensor 2
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

func (lea *larodExampleApplication) feedDModel(fdata []byte) error {
	return lea.DetectionModel.Inputs[0].CopyDataInto(fdata)
}

func (lea *larodExampleApplication) getDResult() ([]byte, error) {
	persons, err := lea.DetectionModel.Outputs[0].GetData(4)
	if err != nil {
		return nil, err
	}
	car, err := lea.DetectionModel.Outputs[1].GetData(4)
	if err != nil {
		return nil, err
	}
	output := make([]byte, 8)
	copy(output[0:4], persons)
	copy(output[4:8], car)
	return output, nil
}

func (lea *larodExampleApplication) Inference() (*axlarod.JobResult, error) {
	// Since larodOutputAddr points to the beginning of the fd we should
	// rewind the file position before each job.
	_, err = lea.DetectionModel.Outputs[0].MemMapFile.File.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	_, err = lea.DetectionModel.Outputs[1].MemMapFile.File.Seek(0, 0)
	if err != nil {
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

type PredictionResult struct {
	Persons float32
	Car     float32
}

func (lea *larodExampleApplication) PredictionResultConverter(result []byte) (*PredictionResult, error) {
	if len(result) < 8 {
		return nil, fmt.Errorf("result slice too short, expected at least 8 bytes, got %d", len(result))
	}

	personReader := bytes.NewReader(result[0:4])
	carReader := bytes.NewReader(result[4:8])

	var person, car float32

	// Read the data into the float32 variables
	if err := binary.Read(personReader, binary.LittleEndian, &person); err != nil {
		return nil, fmt.Errorf("failed to read person data: %v", err)
	}
	if err := binary.Read(carReader, binary.LittleEndian, &car); err != nil {
		return nil, fmt.Errorf("failed to read car data: %v", err)
	}

	return &PredictionResult{Persons: person * 100, Car: car * 100}, nil
}
