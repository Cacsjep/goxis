package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acap"
	"github.com/Cacsjep/goxis/pkg/acapapp"
)

const VDO_CHANNEL = 0

var (
	err              error
	app              *acapapp.AcapApplication
	stream           *acap.VdoStream
	VDO_FORMAT       = acap.VdoFormatH265
	VDO_H265_PROFILE = acap.VdoH265ProfileMain
	VDO_STREAM_CFG   = acap.VideoSteamConfiguration{
		Format:      &VDO_FORMAT,
		H265Profile: &VDO_H265_PROFILE,
	}
)

func main() {
	if app, err = acapapp.NewAcapApplication(); err != nil {
		panic(err)
	}

	if stream, err = acap.CreateAndStartStream(VDO_STREAM_CFG); err != nil {
		app.Syslog.Error(err.Error())
		panic(err)
	}

	for {
		video_frame := acap.GetVideoFrame(stream)
		if video_frame.Error != nil {
			if video_frame.ErrorExpected {
				app.Syslog.Warn(fmt.Sprintf("Vdo stream error, attempting to restart stream..., %s", video_frame.Error.Error()))
				stream, err = acap.RestartStream(stream, VDO_STREAM_CFG)
				if err != nil {
					app.Syslog.Warn(fmt.Sprintf("Unable to restart vdo stream, %s", err.Error()))
				}
				continue
			}
			app.Syslog.Error(err.Error())
			panic(err)
		}
		app.Syslog.Info(video_frame.String())
	}
}
