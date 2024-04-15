package main

import "github.com/Cacsjep/goxis/pkg/axvdo"

func (lea *larodExampleApplication) InitalizeAndStartVdo() error {
	vdo_format := axvdo.VdoFormatYUV
	stream_cfg := axvdo.VideoSteamConfiguration{Format: &vdo_format, Width: &lea.streamWidth, Height: &lea.streamHeight, Framerate: &lea.fps}

	if lea.fp, err = lea.app.NewFrameProvider(stream_cfg); err != nil {
		return err
	}

	if err = lea.fp.Start(); err != nil {
		return err
	}

	lea.app.AddCloseCleanFunc(func() {
		lea.fp.Stop()
	})
	return nil
}
