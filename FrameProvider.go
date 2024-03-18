package goxis

import (
	"github.com/Cacsjep/goxis/pkg/acap"
)

type FrameProviderState int

const (
	FrameProviderStateError FrameProviderState = iota
	FrameProviderStateRunning
	FrameProviderStateStopped
	FrameProviderStateStarted
	FrameProviderStateRestarting
	FrameProviderStateInit
	MaxRestartRetries int = 4
)

type FrameProvider struct {
	Config             acap.VideoSteamConfiguration
	stream             *acap.VdoStream
	state              FrameProviderState
	running            bool
	FrameStreamChannel chan *acap.VideoFrame
	restartRetries     int
}

func (fp *FrameProvider) NewFrameProvider(config acap.VideoSteamConfiguration) (*FrameProvider, error) {
	stream, err := fp.createStream()
	if err != nil {
		return nil, err
	}
	return &FrameProvider{
		stream:             stream,
		Config:             config,
		state:              FrameProviderStateInit,
		FrameStreamChannel: make(chan *acap.VideoFrame, 1),
		running:            false,
	}, nil
}

func (fp *FrameProvider) createStream() (*acap.VdoStream, error) {
	return acap.NewVideoStreamFromConfig(fp.Config)
}

func (fp *FrameProvider) Start() error {
	if err := fp.stream.Start(); err != nil {
		return err
	}

	fp.state = FrameProviderStateStarted
	fp.running = true

	go func() {
		for fp.running {
			video_frame := acap.GetVideoFrame(fp.stream)
			if video_frame.Error != nil {
				if video_frame.ErrorExpected {
					if err := fp.Restart(); err != nil {
						if fp.restartRetries == MaxRestartRetries {
							fp.state = FrameProviderStateError
							break
						}
						fp.restartRetries += 1
					}
					continue
				}
			}
			fp.restartRetries = 0
			fp.FrameStreamChannel <- video_frame
		}
	}()
	return nil
}

func (fp *FrameProvider) Restart() error {
	var err error
	fp.state = FrameProviderStateRestarting
	fp.Stop()

	if fp.stream, err = fp.createStream(); err != nil {
		return err
	}
	return fp.Start()
}

func (fp *FrameProvider) Stop() {
	fp.running = false
	fp.state = FrameProviderStateStopped
	fp.stream.Stop()
	fp.stream.Unref()
}

func (fp *FrameProvider) State() FrameProviderState {
	return fp.state
}
