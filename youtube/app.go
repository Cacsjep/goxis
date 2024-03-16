package main

import (
	"fmt"
	"os"

	"github.com/Cacsjep/goxis/pkg/app"
	"github.com/Cacsjep/goxis/pkg/axvdo"
	"github.com/tinyzimmer/go-gst/gst"
)

func main() {
	gst.Init(&os.Args)
	pipeline, err := gst.NewPipelineFromString("videotestsrc ! videoconvert ! autovideosink")
	if err != nil {
		panic(err)
	}
	fmt.Println(pipeline)
	acapApp, err := app.NewAcapApplication()
	if err != nil {
		panic(err)
	}

	reso, err := acapApp.GetVdoChannelMaxResolution(0)
	if err != nil {
		panic(err)
	}

	var stream *axvdo.VdoStream
	format := axvdo.VdoFormatH265
	if stream, err = acapApp.NewVideoStream(app.VideoSteamConfiguration{
		Format: &format,
		Width:  &reso.Width,
		Height: &reso.Height,
	}); err != nil {
		panic(err)
	}

	if err = stream.Start(); err != nil {
		panic(err)
	}

	for true {
		buf, vdo_err := stream.GetBuffer()
		if err != nil {
			panic(fmt.Sprintf("Vdo Error: %s, Vdo Error Code: %d, Vdo Error Excepted: %t", vdo_err.Err.Error(), vdo_err.Code, vdo_err.Expected))
		}

		frame, err := buf.GetFrame()
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}

		data, err := buf.GetBytesUnsafe()
		if err != nil {
			panic(err)
		}
		fmt.Println("Got jpeg frame with size:", len(data), "TS:", frame.GetTimestamp())
		stream.BufferUnref(buf)
	}
}
