package main

import (
	"github.com/Cacsjep/go-astiav"
	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

func main() {
	app := goxis.NewAcapApplication()
	server_url := "rtmp://fr.castr.io/static/live_2259a220ed2011ee831901652c9e69d2"
	vdo_format := acap.VdoFormatH264
	hprofile := acap.VdoH264ProfileBaseline
	zip_off := acap.VdoZipStreamProfileNone
	doff := false
	bitrate := uint32(5000000)
	brmode := acap.VdoRateControlModeCBR
	stream_cfg := acap.VideoSteamConfiguration{
		Format:          &vdo_format,
		H264Profile:     &hprofile,
		ZipProfile:      &zip_off,
		DynamicGOP:      &doff,
		Bitrate:         &bitrate,
		RateControlMode: &brmode,
	}

	fp, err := app.NewFrameProvider(stream_cfg)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	app.AddCloseCleanFunc(fp.Stop)

	if err := fp.Start(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	stats, err := fp.Stats()
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	stats.StreamStats.PrintStreamStats()

	r, err := NewRtmpStreamer(app, server_url, &RtmpStreamConfig{
		Width:       int(stats.StreamStats.Width),
		Height:      int(stats.StreamStats.Height),
		Fps:         int(stats.StreamStats.Framerate),
		CodecId:     astiav.CodecIDH264,
		Pixelformat: astiav.PixelFormatYuv420P,
	})
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	app.AddCloseCleanFunc(r.ForceStop)
	app.AddCloseCleanFunc(r.Free)

	firstIDRFrame := false

	for frame := range fp.FrameStreamChannel {
		if frame.Error != nil {
			app.Syslog.Errorf("Unexpected Vdo Error: %s", frame.Error.Error())
			continue
		}
		if frame.Type == acap.VdoFrameTypeH264IDR && !firstIDRFrame {
			if err := r.Start(frame.HeaderData()); err != nil {
				app.Syslog.Crit(err.Error())
			}
			firstIDRFrame = true
		}
		if err := r.Write(frame.Data); err != nil {
			app.Syslog.Error(err.Error())
		}
	}

}
