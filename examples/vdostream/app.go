package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

var (
	err error
	app *goxis.AcapApplication
	// FrameProvider for easy go channeld based frame recv
	fp *goxis.FrameProvider
	// The format for the vdo example
	vdo_format = acap.VdoFormatH265
	// Stream configuration
	stream_cfg = acap.VideoSteamConfiguration{
		Format: &vdo_format,
	}
)

func main() {
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	// FrameProvider for easy interact with VDO
	// Easy method to creates VDO streams via go struct acap.VideoSteamConfiguration
	// It automtically restarts vdo stream on maintance.
	// Provide also access to stream stats via go structs
	if fp, err = app.NewFrameProvider(stream_cfg); err != nil {
		app.Syslog.Error(err.Error())
		panic(err)
	}

	// Start the frameprovider
	if err = fp.Start(); err != nil {
		app.Syslog.Error(err.Error())
		panic(err)
	}
	defer fp.Stop()

	// Enter channel based recv for *VideoFrame from FrameStreamChannel
	// *VideoFrame holds also Error that are either expected during maintanance
	// or unexpected errors. Expected errors are detected automatically in the
	// Frameprovider and force a stream restart.
	// All Errors that are recived here are unexpected errors if frame.Error not nil!
	for {
		select {
		case frame := <-fp.FrameStreamChannel:
			if frame.Error != nil {
				app.Syslog.Error(fmt.Sprintf("Unexpected Vdo Error: %s", frame.Error.Error()))
				continue
			}
			app.Syslog.Info(frame.String())
		}
	}

}
