package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

var (
	fp         *acapapp.FrameProvider
	vdo_format = axvdo.VdoFormatH265
	stream_cfg = axvdo.VideoSteamConfiguration{Format: &vdo_format}
	err        error
)

func main() {
	app := acapapp.NewAcapApplication()

	// FrameProvider for easy go channeld based frame receiving
	if fp, err = app.NewFrameProvider(stream_cfg); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Start the frameprovider
	if err = fp.Start(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.AddCloseCleanFunc(fp.Stop)

	// ! VideoFrame has also an Error field for unexpected errors.
	// * Expected errors are detected automatically in the Frameprovider and force a stream restart.
	for {
		select {
		case frame := <-fp.FrameStreamChannel:
			if frame.Error != nil {
				app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}
			app.Syslog.Info(frame.String())

			// Do something with the frame like frame.Data()
		}
	}
}
