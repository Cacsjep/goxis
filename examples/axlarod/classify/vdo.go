package main

import "github.com/Cacsjep/goxis/pkg/axvdo"

// InitializeAndStartVdo configures and starts a video stream based on predefined settings.
// It sets the video format to YUV and applies the specified resolution and framerate from the larodExampleApplication struct.
// This function handles the creation and activation of the frame provider which captures video frames.
// Returns an error if there are issues initializing or starting the video frame provider.
func (lea *larodExampleApplication) InitalizeAndStartVdo() error {
	vdo_format := axvdo.VdoFormatYUV
	stream_cfg := axvdo.VideoSteamConfiguration{Format: &vdo_format, Width: &lea.streamWidth, Height: &lea.streamHeight, Framerate: &lea.fps}

	if err = lea.app.NewFrameProvider(stream_cfg); err != nil {
		return err
	}

	if err = lea.app.FrameProvider.Start(); err != nil {
		return err
	}
	return nil
}
