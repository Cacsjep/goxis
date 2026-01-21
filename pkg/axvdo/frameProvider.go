package axvdo

import (
	"sync"
	"sync/atomic"
	"time"
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
	Config             VideoSteamConfiguration // Configuration for the video stream.
	Stream             *VdoStream              // Internal video stream reference.
	state              FrameProviderState      // Current state of the frame provider.
	running            int32                   // Flag indicating whether the frame provider is actively running (atomic).
	wg                 sync.WaitGroup          // WaitGroup to track the frame fetching goroutine.
	FrameStreamChannel chan *VideoFrame        // Channel for delivering video frames to consumers.
	restartRetries     int                     // Counter for the number of restart attempts.
	outReso            *VdoResolution
	frameProccessor    func([]byte) []byte
	RestartCallback    func() // Callback function to be called when the frame provider is restarted
}

// FrameProviderStats provides statistical information about the operation of a FrameProvider.
type FrameProviderStats struct {
	InternalChannelBufferLen int         // The current length of the frame stream channel buffer.
	RestartRetries           int         // The number of restart attempts made since the last successful start.
	StreamStats              StreamStats // Statistics gathered from the video stream.
}

// NewFrameProvider initializes a new FrameProvider with the given configuration and application context.
// It prepares the frame provider for operation but does not start streaming frames until Start is called.
func NewFrameProvider(config VideoSteamConfiguration) (*FrameProvider, error) {
	fp := &FrameProvider{
		Config:             config,
		state:              FrameProviderStateInit,
		FrameStreamChannel: make(chan *VideoFrame, 1),
		// running defaults to 0 (stopped)
	}
	stream, err := fp.createStream()
	if err != nil {
		return nil, err
	}
	fp.Stream = stream
	return fp, nil
}

// createStream initializes the video stream based on the FrameProvider's configuration.
// This is a helper method used internally by the FrameProvider.
func (fp *FrameProvider) createStream() (*VdoStream, error) {
	return NewVideoStreamFromConfig(fp.Config)
}

// Start begins the frame streaming process, marking the FrameProvider as running and initiating the frame fetching loop.
// If an error occurs while starting the stream, it returns the error without altering the provider's state.
// Handles automatic restart in case of an expected Vdo error
func (fp *FrameProvider) Start() error {
	if err := fp.Stream.Start(); err != nil {
		return err
	}

	atomic.StoreInt32(&fp.running, 1)
	fp.state = FrameProviderStateStarted

	fp.wg.Add(1)
	go func() {
		defer fp.wg.Done()
		for atomic.LoadInt32(&fp.running) == 1 {
			video_frame := GetVideoFrame(fp.Stream)
			if video_frame.Error != nil {
				// Check if we were stopped - exit immediately
				if atomic.LoadInt32(&fp.running) == 0 {
					return
				}
				if video_frame.ErrorExpected {
					if fp.state == FrameProviderStateStopped {
						return
					}
					if err := fp.Restart(); err != nil {
						if fp.restartRetries >= MaxRestartRetries {
							fp.state = FrameProviderStateError
							break
						}
						fp.restartRetries++
					} else {
						atomic.StoreInt32(&fp.running, 1)
						fp.state = FrameProviderStateStarted
					}
					continue
				}
				continue
			}
			fp.restartRetries = 0
			fp.FrameStreamChannel <- video_frame
		}
	}()
	return nil
}

// Stop halts the frame streaming process, changing the state of the FrameProvider to stopped and cleaning up resources.
// This function waits for the frame fetching goroutine to exit before releasing stream resources.
func (fp *FrameProvider) Stop() {
	atomic.StoreInt32(&fp.running, 0)
	fp.state = FrameProviderStateStopped
	fp.Stream.Stop()
	fp.wg.Wait() // Wait for goroutine to exit before releasing resources
	fp.Stream.Unref()
}

// Restart attempts to restart the video stream, first stopping the current stream and then re-initializing and starting a new stream.
// It applies a delay before attempting the restart to give the system time to release resources.
// Note: This is called from within the goroutine, so we don't use Stop() which would deadlock waiting for ourselves.
func (fp *FrameProvider) Restart() error {
	if fp.state == FrameProviderStateStopped {
		return nil
	}
	time.Sleep(time.Second * 2)
	var err error
	fp.state = FrameProviderStateRestarting
	// Don't call Stop() here - we're inside the goroutine and Stop() waits for the goroutine.
	// Just stop and unref the stream directly since we are the only user at this point.
	fp.Stream.Stop()
	fp.Stream.Unref()
	if fp.Stream, err = fp.createStream(); err != nil {
		return err
	}

	defer func() {
		if fp.RestartCallback != nil {
			fp.RestartCallback()
		}
	}()

	return fp.Stream.Start()
}

// State returns the current state of the FrameProvider, providing insight into whether it's running, stopped, or in an error state.
func (fp *FrameProvider) State() FrameProviderState {
	return fp.state
}

// IsRunning checks if the FrameProvider is currently active and streaming frames.
func (fp *FrameProvider) IsRunning() bool {
	return atomic.LoadInt32(&fp.running) == 1
}

// Stats gathers and returns statistical information about the frame provider's operation, including internal buffer lengths and stream statistics.
func (fp *FrameProvider) Stats() (*FrameProviderStats, error) {
	m, err := fp.Stream.GetInfo()
	if err != nil {
		return nil, err
	}
	stats := StreamStats{
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
		ZipProfile:                    m.GetInt16("zip.profile", 0),
	}

	return &FrameProviderStats{
		StreamStats:              stats,
		RestartRetries:           fp.restartRetries,
		InternalChannelBufferLen: len(fp.FrameStreamChannel),
	}, nil
}
