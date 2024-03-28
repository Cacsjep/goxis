package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Cacsjep/go-astiav"
	"github.com/Cacsjep/goxis"
)

type RtmpStreamer struct {
	AcapApp             *goxis.AcapApplication
	CurrentPkt          *astiav.Packet
	OutputFormatContext *astiav.FormatContext
	OutputStream        *astiav.Stream
	RtmpServerUrl       string
	FrameCounter        int64
	ptsIncrement        int64
	StreamConfig        *RtmpStreamConfig
	ioContext           *astiav.IOContext
}

type RtmpStreamConfig struct {
	Fps         int
	CodecId     astiav.CodecID
	Width       int
	Height      int
	Pixelformat astiav.PixelFormat
}

func NewRtmpStreamer(acap_app *goxis.AcapApplication, server_uri string, stream_cfg *RtmpStreamConfig) (*RtmpStreamer, error) {
	var err error

	r := RtmpStreamer{
		AcapApp:       acap_app,
		RtmpServerUrl: server_uri,
		StreamConfig:  stream_cfg,
	}

	astiav.SetLogLevel(astiav.LogLevel(56))
	astiav.SetLogCallback(func(c astiav.Classer, l astiav.LogLevel, fmt, msg string) {
		var cs string
		if c != nil {
			if cl := c.Class(); cl != nil {
				cs = " - class: " + cl.String()
			}
		}
		r.AcapApp.Syslog.Infof("ffmpeg: %s%s - level: %d\n", strings.TrimSpace(msg), cs, l)
	})

	r.CurrentPkt = astiav.AllocPacket()

	r.OutputFormatContext, err = astiav.AllocOutputFormatContext(nil, "flv", r.RtmpServerUrl)
	if err != nil {
		return nil, fmt.Errorf("allocating output format context failed: %s", err.Error())
	}
	if r.OutputFormatContext == nil {
		return nil, errors.New("output format context is nil")
	}
	//r.OutputFormatContext.Flags().Add(astiav.FormatContextFlag(astiav.CodecContextFlagGlobalHeader))
	r.OutputStream = r.OutputFormatContext.NewStream(nil)
	if r.OutputStream == nil {
		return nil, errors.New("output stream is nil")
	}
	var codecContext *astiav.CodecContext
	if codecContext = astiav.AllocCodecContext(astiav.FindDecoder(stream_cfg.CodecId)); codecContext == nil {
		return nil, errors.New("main: codec context is nil")
	}

	codecContext.Flags().Add(astiav.CodecContextFlagGlobalHeader)
	codecContext.SetTimeBase(astiav.NewRational(1, 90000))
	codecContext.SetHeight(stream_cfg.Height)
	codecContext.SetWidth(stream_cfg.Width)
	codecContext.SetPixelFormat(stream_cfg.Pixelformat)
	codecContext.SetChannels(0)

	if err = r.OutputStream.CodecParameters().FromCodecContext(codecContext); err != nil {
		return nil, fmt.Errorf("setting codec parameters failed: %w", err)
	}

	r.OutputStream.CodecParameters().SetCodecType(astiav.MediaTypeVideo)
	r.OutputStream.CodecParameters().SetCodecTag(0)
	return &r, nil
}

func (r *RtmpStreamer) openContext() error {
	if !r.OutputFormatContext.OutputFormat().Flags().Has(astiav.IOFormatFlagNofile) {
		var err error
		r.ioContext, err = astiav.OpenIOContext(r.RtmpServerUrl, astiav.NewIOContextFlags(astiav.IOContextFlagWrite))
		if err != nil {
			return fmt.Errorf("opening io context failed: %w", err)
		}
		r.OutputFormatContext.SetPb(r.ioContext)
	}
	return nil
}

func (r *RtmpStreamer) Start(extraData []byte) error {
	var err error

	if err = r.OutputStream.CodecParameters().SetExtraData(extraData); err != nil {
		return err
	}

	if err = r.openContext(); err != nil {
		return err
	}
	if err = r.OutputFormatContext.WriteHeader(nil); err != nil {
		return fmt.Errorf("writing header failed: %w", err)
	}
	r.ptsIncrement = 90000 / int64(r.StreamConfig.Fps)
	r.FrameCounter = 0
	return nil
}

func (r *RtmpStreamer) Write(video_data []byte) error {
	var err error
	r.CurrentPkt.FromData(video_data)
	r.CurrentPkt.SetPts(r.FrameCounter * r.ptsIncrement)
	r.CurrentPkt.SetDts(r.CurrentPkt.Pts())
	r.CurrentPkt.SetDuration(r.ptsIncrement)
	if err = r.OutputFormatContext.WriteInterleavedFrame(r.CurrentPkt); err != nil {
		return fmt.Errorf("writing frame failed: %w", err)
	}
	r.FrameCounter++
	return nil
}

func (r *RtmpStreamer) Stop() error {
	var err error
	if err = r.OutputFormatContext.WriteTrailer(); err != nil {
		return fmt.Errorf("writing trailer failed: %w", err)
	}
	r.FrameCounter = 0
	r.ioContext.Closep()
	return nil
}

func (r *RtmpStreamer) ForceStop() {
	r.OutputFormatContext.WriteTrailer()
	r.ioContext.Closep()
}

func (r *RtmpStreamer) Free() {
	r.CurrentPkt.Free()
	r.OutputFormatContext.Free()
}
