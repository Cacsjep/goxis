package main

import (
	"fmt"
	"time"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

const VDO_CHANNEL = 0

var (
	VDO_FORMAT       = acap.VdoFormatH265
	VDO_H265_PROFILE = acap.VdoH265ProfileMain
	VDO_STREAM_CFG   = acap.VideoSteamConfiguration{
		Format:      &VDO_FORMAT,
		H265Profile: &VDO_H265_PROFILE,
	}
)

func main() {
	var err error
	var fp *goxis.FrameProvider
	var app *goxis.AcapApplication

	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}

	if fp, err = goxis.NewFrameProvider(VDO_STREAM_CFG); err != nil {
		app.Syslog.Error(err.Error())
		panic(err)
	}

	fp.Start()
	defer fp.Stop()

	for {
		select {
		case frame := <-fp.FrameStreamChannel:
			if frame.Error != nil {
				fmt.Println("ERR", frame.Error.Error())
				continue
			}
			fmt.Println(frame.String())

		case <-time.After(10 * time.Second):
			fmt.Println("No frame received in 10 seconds.")
		}
	}
}
