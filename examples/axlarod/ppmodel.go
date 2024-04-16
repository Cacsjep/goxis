package main

import (
	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

func (lea *larodExampleApplication) InitalizePPModel() error {
	if lea.PPModel, err = lea.app.Larod.NewPreProccessModel(
		"cpu-proc",
		axlarod.LarodResolution{Width: lea.streamWidth, Height: lea.streamHeight},
		axlarod.LarodResolution{Width: lea.streamWidth, Height: lea.streamHeight},
		axlarod.PreProccessOutputFormatRgbInterleaved,
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

func (lea *larodExampleApplication) feedPPModel(fdata []byte) error {
	return lea.PPModel.Inputs[0].CopyDataInto(fdata)
}

func (lea *larodExampleApplication) getPPResult() ([]byte, error) {
	return lea.PPModel.Outputs[0].GetData(lea.streamWidth * lea.streamHeight * 3)
}

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
