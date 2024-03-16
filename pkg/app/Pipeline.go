package app

import (
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
	"github.com/tinyzimmer/go-gst/gst/video"
)

func NewGstPipeline() (*gst.Pipeline, error) {
	return gst.NewPipeline("")
}

func NewRtmpPipeline(rtmp_server string, width int, height int) (*gst.Pipeline, error) {
	pipeline, err := NewGstPipeline()
	if err != nil {
		return nil, err
	}
	elems, err := gst.NewElementMany("appsrc", "videoconvert", "autovideosink")
	if err != nil {
		return nil, err
	}
	pipeline.AddMany(elems...)
	gst.ElementLinkMany(elems...)

	// Get the app sourrce from the first element returned
	src := app.SrcFromElement(elems[0])

	// Specify the format we want to provide as application into the pipeline
	// by creating a video info with the given format and creating caps from it for the appsrc element.
	videoInfo := video.NewInfo().
		WithFormat(video.FormatNV12, uint(width), uint(height)).
		WithFPS(gst.Fraction(2, 1))

	src.SetCaps(videoInfo.ToCaps())
	src.SetProperty("format", gst.FormatTime)
	return nil, nil
}
