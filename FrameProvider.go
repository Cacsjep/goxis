package goxis

import (
	"fmt"
	"time"

	"github.com/Cacsjep/goxis/pkg/acap"
)

// FrameProviderState defines the possible states of a FrameProvider.
type FrameProviderState int

const (
	// FrameProviderStateError indicates an error state where the frame provider cannot recover without intervention.
	FrameProviderStateError FrameProviderState = iota
	// FrameProviderStateStopped indicates the frame provider is stopped and not currently providing frames.
	FrameProviderStateStopped
	// FrameProviderStateStarted indicates the frame provider is actively providing frames.
	FrameProviderStateStarted
	// FrameProviderStateRestarting indicates the frame provider is in the process of restarting.
	FrameProviderStateRestarting
	// FrameProviderStateInit indicates the frame provider is initialized but not yet started.
	FrameProviderStateInit
	// MaxRestartRetries defines the maximum number of restart attempts for the frame provider before entering an error state.
	MaxRestartRetries int = 4
)

// FrameProvider encapsulates the management of video frame streaming, including starting, stopping, and restarting the stream.
type FrameProvider struct {
	Config             acap.VideoSteamConfiguration // Configuration for the video stream.
	stream             *acap.VdoStream              // Internal video stream reference.
	state              FrameProviderState           // Current state of the frame provider.
	running            bool                         // Flag indicating whether the frame provider is actively running.
	FrameStreamChannel chan *acap.VideoFrame        // Channel for delivering video frames to consumers.
	restartRetries     int                          // Counter for the number of restart attempts.
	app                *AcapApplication             // Reference to the application managing this frame provider.
}

// FrameProviderStats provides statistical information about the operation of a FrameProvider.
type FrameProviderStats struct {
	InternalChannelBufferLen int              // The current length of the frame stream channel buffer.
	RestartRetries           int              // The number of restart attempts made since the last successful start.
	StreamStats              acap.StreamStats // Statistics gathered from the video stream.
}

// NewFrameProvider initializes a new FrameProvider with the given configuration and application context.
// It prepares the frame provider for operation but does not start streaming frames until Start is called.
func (a *AcapApplication) NewFrameProvider(config acap.VideoSteamConfiguration) (*FrameProvider, error) {
	fp := &FrameProvider{
		Config:             config,
		state:              FrameProviderStateInit,
		FrameStreamChannel: make(chan *acap.VideoFrame, 30),
		running:            false,
		app:                a,
	}
	stream, err := fp.createStream()
	if err != nil {
		return nil, err
	}
	fp.stream = stream
	return fp, nil
}

// createStream initializes the video stream based on the FrameProvider's configuration.
// This is a helper method used internally by the FrameProvider.
func (fp *FrameProvider) createStream() (*acap.VdoStream, error) {
	return acap.NewVideoStreamFromConfig(fp.Config)
}

// Start begins the frame streaming process, marking the FrameProvider as running and initiating the frame fetching loop.
// If an error occurs while starting the stream, it returns the error without altering the provider's state.
// Handles automatic restart in case of an expected Vdo error
func (fp *FrameProvider) Start() error {
	if err := fp.stream.Start(); err != nil {
		return err
	}

	fp.running = true
	fp.state = FrameProviderStateStarted
	fp.app.Syslog.Info(fmt.Sprintf("VDO Channel(%d): Stream is started", fp.Config.GetChannel()))

	go func() {
		for fp.running {
			video_frame := acap.GetVideoFrame(fp.stream)
			if video_frame.Error != nil {
				if video_frame.ErrorExpected {
					fp.app.Syslog.Warn(fmt.Sprintf("VDO Channel(%d): Restarting stream because vdo is in maintanance mode %s", fp.Config.GetChannel(), video_frame.Error.Error()))
					if err := fp.Restart(); err != nil {
						fp.app.Syslog.Warn(fmt.Sprintf("VDO Channel(%d): Unable to restart stream, try again...: %s", fp.Config.GetChannel(), err.Error()))
						if fp.restartRetries >= MaxRestartRetries {
							fp.state = FrameProviderStateError
							fp.app.Syslog.Error(fmt.Sprintf("VDO Channel(%d): Max retries for stream restart reached, stream is stopped", fp.Config.GetChannel()))
							break
						}
						fp.restartRetries++
					} else {
						fp.app.Syslog.Info(fmt.Sprintf("VDO Channel(%d): Successfully restart stream", fp.Config.GetChannel()))
						fp.running = true
						fp.state = FrameProviderStateStarted
					}
					continue
				}
				fp.app.Syslog.Error(fmt.Sprintf("VDO Channel(%d): Vdo returns an error when getting buffer/frame data %s", fp.Config.Channel, video_frame.Error.Error()))
				continue
			}
			fp.restartRetries = 0
			fp.FrameStreamChannel <- video_frame
		}
	}()
	return nil
}

// Stop halts the frame streaming process, changing the state of the FrameProvider to stopped and cleaning up resources.
func (fp *FrameProvider) Stop() {
	fp.running = false
	fp.state = FrameProviderStateStopped
	fp.stream.Stop()
	fp.stream.Unref()
	fp.app.Syslog.Info(fmt.Sprintf("VDO Channel(%d): Stream is stopped", fp.Config.GetChannel()))
}

// Restart attempts to restart the video stream, first stopping the current stream and then re-initializing and starting a new stream.
// It applies a delay before attempting the restart to give the system time to release resources.
func (fp *FrameProvider) Restart() error {
	time.Sleep(time.Second * 2)
	fp.app.Syslog.Info(fmt.Sprintf("VDO Channel(%d): Try to restart stream", fp.Config.GetChannel()))
	var err error
	fp.state = FrameProviderStateRestarting
	fp.Stop()
	if fp.stream, err = fp.createStream(); err != nil {
		return err
	}
	return fp.stream.Start()
}

// State returns the current state of the FrameProvider, providing insight into whether it's running, stopped, or in an error state.
func (fp *FrameProvider) State() FrameProviderState {
	return fp.state
}

// IsRunning checks if the FrameProvider is currently active and streaming frames.
func (fp *FrameProvider) IsRunning() bool {
	return fp.running
}

// Stats gathers and returns statistical information about the frame provider's operation, including internal buffer lengths and stream statistics.
func (fp *FrameProvider) Stats() (*FrameProviderStats, error) {
	m, err := fp.stream.GetInfo()
	if err != nil {
		return nil, err
	}
	stats := acap.StreamStats{
		Bitrate:                       m.GetUint32("bitrate", 0),
		BufferType:                    m.GetString("buffer.type", ""),
		Channel:                       m.GetUint32("channel", 0),
		Format:                        m.GetUint32("format", 0),
		Framerate:                     m.GetUint32("framerate", 0),
		GOPLength:                     m.GetUint32("gop_length", 0),
		H26xIntraRefresh:              m.GetUint32("h26x.intra_refresh", 0),
		Height:                        m.GetUint32("height", 0),
		HorizontalFlip:                m.GetBoolean("horizontal_flip", false),
		ID:                            m.GetUint32("id", 0),
		InitialBitrate:                m.GetUint32("initial.bitrate", 0),
		InitialQPb:                    m.GetUint32("initial.qp.b", 0),
		InitialQPi:                    m.GetUint32("initial.qp.i", 0),
		InitialQPp:                    m.GetUint32("initial.qp.p", 0),
		Overlays:                      m.GetString("overlays", ""),
		Peers:                         m.GetUint32("peers", 0),
		QPb:                           m.GetUint32("qp.b", 0),
		QPi:                           m.GetUint32("qp.i", 0),
		QPp:                           m.GetUint32("qp.p", 0),
		Rotation:                      m.GetUint32("rotation", 0),
		Running:                       m.GetBoolean("running", false),
		SquarePixel:                   m.GetUint32("squarepixel", 0),
		StatisticsAccumulatedBytes:    m.GetUint64("statistics.accumulated_bytes", 0),
		StatisticsAccumulatedIDRBytes: m.GetUint64("statistics.accumulated_idrbytes", 0),
		StatisticsBitCount:            m.GetUint32("statistics.bit_count", 0),
		StatisticsBitrate:             m.GetUint32("statistics.bitrate", 0),
		StatisticsDuration:            m.GetInt64("statistics.duration", 0),
		StatisticsDynamicFramerate:    m.GetUint32("statistics.dynamic_framerate", 0),
		StatisticsFailedFrames:        m.GetUint32("statistics.failed_frames", 0),
		StatisticsFrameCount:          m.GetUint32("statistics.frame_count", 0),
		StatisticsFramerate:           m.GetUint32("statistics.framerate", 0),
		StatisticsIDRFrameCount:       m.GetUint32("statistics.idrframe_count", 0),
		StatisticsLastFrameTS:         m.GetUint64("statistics.last_frame_ts", 0),
		StatisticsReclaimCount:        m.GetUint32("statistics.reclaim_count", 0),
		Width:                         m.GetUint32("width", 0),
		ZipProfile:                    m.GetUint32("zip.profile", 0),
	}

	return &FrameProviderStats{
		StreamStats:              stats,
		RestartRetries:           fp.restartRetries,
		InternalChannelBufferLen: len(fp.FrameStreamChannel),
	}, nil
}
