package main

import (
	"github.com/Cacsjep/goxis/pkg/axlarod"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

// InitializePPModel initializes a preprocessing model tailored for video processing.
// It sets the model to operate in the specified RGB mode and resolution based on the current application settings.
// An error is returned if the model fails to initialize or if any issues occur during setup.
func (lea *larodExampleApplication) InitalizePPModel(rgbMode axlarod.PreProccessOutputFormat) error {
	cropMap, err := axlarod.CreateCropMap(lea.cocoInputWidth, lea.cocoInputHeight, lea.streamWidth, lea.streamHeight)
	if err != nil {
		return err
	}
	if lea.PPModel, err = lea.app.Larod.NewPreProccessModel(
		"cpu-proc",
		axlarod.LarodResolution{Width: lea.streamWidth, Height: lea.streamHeight},
		axlarod.LarodResolution{Width: lea.cocoInputWidth, Height: lea.cocoInputHeight},
		rgbMode,
		cropMap,
	); err != nil {
		return err
	}
	lea.app.AddModelCleaner(lea.PPModel)
	return err
}

// PreProcess handles the preprocessing of video frames using the initialized preprocessing model.
// It manages the flow of data into the model and retrieves the processed output.
// Returns a JobResult containing the processed data or an error if preprocessing fails.
func (lea *larodExampleApplication) PreProcess(frame *axvdo.VideoFrame) (*axlarod.JobResult, error) {
	var result *axlarod.JobResult
	if result, err = lea.app.Larod.ExecuteJob(lea.PPModel, func() error {
		return lea.PPModel.Inputs[0].CopyDataInto(frame.Data)
	}, func() (any, error) {
		return lea.PPModel.Outputs[0].GetData(lea.streamWidth * lea.streamHeight * 3)
	}); err != nil {
		return nil, err
	}
	return result, nil
}
