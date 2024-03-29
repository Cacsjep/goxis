package main

import (
	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/axvdo"
)

var (
	err error
	app *goxis.AcapApplication
	// FrameProvider for easy go channeld based frame recv
	fp *goxis.FrameProvider
	// The format for the vdo example
	vdo_format = axvdo.VdoFormatH265
	// Stream configuration
	stream_cfg = axvdo.VideoSteamConfiguration{
		Format: &vdo_format,
	}
)

func main() {
	app = goxis.NewAcapApplication()

	// FrameProvider for easy interact with VDO
	// Easy method to creates VDO streams via go struct axvdo.VideoSteamConfiguration
	// It automtically restarts vdo stream on maintance.
	// Provide also access to stream stats via go structs
	if fp, err = app.NewFrameProvider(stream_cfg); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Start the frameprovider
	if err = fp.Start(); err != nil {
		app.Syslog.Crit(err.Error())
	}
	app.AddCloseCleanFunc(fp.Stop)

	// Enter channel based recv for *VideoFrame from FrameStreamChannel
	// *VideoFrame holds also Error that are either expected during maintanance
	// or unexpected errors. Expected errors are detected automatically in the
	// Frameprovider and force a stream restart.
	// All Errors that are recived here are unexpected errors if frame.Error not nil!
	for {
		select {
		case frame := <-fp.FrameStreamChannel:
			if frame.Error != nil {
				app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
				continue
			}
			app.Syslog.Info(frame.String())
		}
	}

}
