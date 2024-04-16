package main

import (
	"github.com/Cacsjep/goxis/pkg/axlarod"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

// InitializePPModel initializes a preprocessing model tailored for video processing.
// It sets the model to operate in the specified RGB mode and resolution based on the current application settings.
// An error is returned if the model fails to initialize or if any issues occur during setup.
func (lea *larodExampleApplication) InitalizePPModel(rgbMode axlarod.PreProccessOutputFormat) error {
	if lea.PPModel, err = lea.app.Larod.NewPreProccessModel(
		"cpu-proc",
		axlarod.LarodResolution{Width: lea.streamWidth, Height: lea.streamHeight},
		axlarod.LarodResolution{Width: lea.streamWidth, Height: lea.streamHeight},
		rgbMode,
	); err != nil {
		return err
	}

	lea.app.AddCloseCleanFunc(func() {
		err := lea.app.Larod.DestroyModel(lea.PPModel)
		if err != nil {
			lea.app.Syslog.Errorf("Failed to destroy PPModel: %s", err.Error())
		}
	})
	return err
}

// feedPPModel takes video frame data as bytes and copies it into the preprocessing model's input tensor.
// Returns an error if the data copying process fails.
func (lea *larodExampleApplication) feedPPModel(fdata []byte) error {
	return lea.PPModel.Inputs[0].CopyDataInto(fdata)
}

// getPPResult retrieves the output data from the preprocessing model after processing a video frame.
// It returns the output data as bytes or an error if the data retrieval fails.
func (lea *larodExampleApplication) getPPResult() ([]byte, error) {
	return lea.PPModel.Outputs[0].GetData(lea.streamWidth * lea.streamHeight * 3)
}

// PreProcess handles the preprocessing of video frames using the initialized preprocessing model.
// It manages the flow of data into the model and retrieves the processed output.
// Returns a JobResult containing the processed data or an error if preprocessing fails.
func (lea *larodExampleApplication) PreProcess(frame *axvdo.VideoFrame) (*axlarod.JobResult, error) {
	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.PPModel, func() error {
		return lea.feedPPModel(frame.Data)
	}, func() ([]byte, error) {
		return lea.getPPResult()
	}); err != nil {
		return nil, err
	}
	return result, nil
}
