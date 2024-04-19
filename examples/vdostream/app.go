package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

var (
	vdo_format = axvdo.VdoFormatH265
	stream_cfg = axvdo.VideoSteamConfiguration{Format: &vdo_format}
	err        error
)

// This example demonstrates how to use the axvdo to receive video frames from the camera.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/vdostream
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	// FrameProvider for easy go channeld based frame receiving
	if err = app.NewFrameProvider(stream_cfg); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Start the frameprovider
	if err = app.FrameProvider.Start(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// ! VideoFrame has also an Error field for unexpected errors.
	// * Expected errors are detected automatically in the Frameprovider and force a stream restart.
	for {
		select {
		case frame := <-app.FrameProvider.FrameStreamChannel:
			if frame.Error != nil {
				app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}
			app.Syslog.Info(frame.String())

			// Do something with the frame like frame.Data()
		}
	}
}
