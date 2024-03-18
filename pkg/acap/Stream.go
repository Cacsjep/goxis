package acap

import (
	"fmt"
	"time"
)

type StreamRotation int

const (
	StreamRotationNone StreamRotation = 0
	StreamRotation90   StreamRotation = 90
	StreamRotation180  StreamRotation = 180
	StreamRotation270  StreamRotation = 270
)

type ZipGopMode int

const (
	ZipGOPModeFixed   ZipGopMode = 0
	ZipGOPModeDynamic ZipGopMode = 1
)

type ZipFPSMode int

const (
	ZipFPSModeFixed   ZipFPSMode = 0
	ZipFPSModeDynamic ZipFPSMode = 1
)

type ZipSkipMode int

const (
	ZipSkipModeDrop  ZipSkipMode = 0
	ZipSkipModeEmpty ZipSkipMode = 1
)

type VideoSteamConfiguration struct {
	Format              *VdoFormat              //Video stream format as VdoFormat.
	BufferCount         *int                    // How many in-flight buffers are desired.
	BufferStrategy      *VdoBufferStrategy      // Buffering Strategy as VdoBufferStrategy.
	Input               *int                    // Video input, 1 ... inmax. 0 is invalid. No view areas.
	Channel             *int                    // Video channel, 0 is overview, 1, 2, ... are view areas.
	Width               *int                    // Video stream horizontal resolution.
	Height              *int                    // Video stream vertical resolution.
	Framerate           *int                    // Video stream vertical resolution.
	Compression         *int                    // Video stream compression, Axis standard range [0:100]
	Rotation            *StreamRotation         // Video stream rotation, normally [0,90,180,270].
	HorizontalFlip      *bool                   // Video stream horizontal flip (mirroring).
	VerticalFlip        *bool                   // Video stream vertical flip.
	Monochrome          *bool                   // Video stream monochrome encoding.
	DynamicGOP          *bool                   // Enable dynamic gop
	DynamicBitrate      *bool                   // Enable dynamic bitrate
	DynamicFramerate    *bool                   // Enable dynamic framerate
	DynamicCompression  *bool                   // Enable dynamic compression
	Qpi                 *uint32                 // QP value for I-frames.
	Qpp                 *uint32                 // QP value for P-frames.
	Bitrate             *uint32                 // Video stream bitrate (bps)
	RateControlMode     *VdoRateControlMode     // Bitrate control mode.
	RateControlPriority *VdoRateControlPriority // Bitrate control priority.
	GOPLength           *uint32                 // GOP length.
	// H.264 Specific Settings
	H264Profile *VdoH264Profile // H.265 profile as VdoH264Profile
	H265Profile *VdoH265Profile // H.264 profile as VdoH264Profile
	// Zipstream Specific Settings
	ZipStrength     *uint32              // Zipstream strength.
	ZipMaxGOPLength *uint32              // Zipstream maximum GOP length.
	ZipGOPMode      *ZipGopMode          // Zipstream GOP mode [0 = fixed, 1 = dynamic].
	ZipFPSMode      *ZipFPSMode          // Zipstream framerate control mode: [0 = fixed, 1 = dynamic].
	ZipSkipMode     *ZipSkipMode         // Zipstream frame skip mode: [0 = drop, 1 = empty].
	ZipMinFPSNum    *uint32              // Zipstream minimum framerate numerator.
	ZipMinFPSDen    *uint32              // Zipstream minimum framerate denominator.
	ZipProfile      *VdoZipStreamProfile // Zipstream profile.
	// ABR Specific Settings
	// The following ABR specific settings are supported with VDO_RATE_CONTROL_MODE_ABR.
	AbrTarget_bitrate *uint32 // Stream target bitrate (bps)
	AbrRetention_time *uint32 // Retention time in seconds
}

type VideoFrame struct {
	SequenceNbr   uint
	Timestamp     time.Time
	Size          uint
	Data          []byte
	Type          VdoFrameType
	Info          *VdoMap
	Error         error
	ErrorExpected bool
}

func (f *VideoFrame) String() string {
	return fmt.Sprintf("SequenceNbr: %d, Timestamp: %s, Size: %d, Type: %s",
		f.SequenceNbr, f.Timestamp.Format("2006-01-02 15:04:05"), f.Size, f.Type.String())
}

func NewVideoFrame(frame *VdoFrame, data []byte) *VideoFrame {
	return &VideoFrame{
		SequenceNbr: frame.GetSequenceNbr(),
		Timestamp:   time.Unix(frame.GetCustomTimestamp()/1000000, (frame.GetCustomTimestamp()%1000000)*1000),
		Size:        frame.GetSize(),
		Type:        frame.GetFrameType(),
		Data:        data,
	}
}

func NewVideoStreamFromConfig(stream_cfg VideoSteamConfiguration) (*VdoStream, error) {
	vdoMap := VideoStreamConfigToVdoMap(stream_cfg)
	defer vdoMap.Unref()
	return NewStream(vdoMap)
}

// Gets a vdo buffer and frame data and wraps it into a *VideoFrame
// Errors are set in the struct if vdo error expected ErrorExpected is set.
// ErrorExpected should use to restart a stream because it means maintance like WDR change.
func GetVideoFrame(vdo_stream *VdoStream) *VideoFrame {
	vdo_buf, err := vdo_stream.GetBuffer()
	if err != nil {
		if vdoErr, ok := err.(*VdoError); ok && vdoErr.Expected {
			return &VideoFrame{Error: vdoErr, ErrorExpected: vdoErr.Expected}
		}
		// Retry when its not an expected error
		return GetVideoFrame(vdo_stream)
	}

	vdo_frame, err := vdo_buf.GetFrame()
	if err != nil {
		return &VideoFrame{Error: err, ErrorExpected: false}
	}

	buff_data, err := vdo_buf.GetBytes()
	if err != nil {
		return &VideoFrame{Error: err, ErrorExpected: false}
	}

	defer vdo_stream.BufferUnref(vdo_buf)
	return NewVideoFrame(vdo_frame, buff_data)
}

func VideoStreamConfigToVdoMap(cfg VideoSteamConfiguration) *VdoMap {
	m := NewVdoMap()

	// Utilizing a helper function to streamline nil checks and assignments
	setUint32IfNotNil := func(key string, value *uint32) {
		if value != nil {
			m.SetUint32(key, *value)
		}
	}
	setIntIfNotNil := func(key string, value *int) {
		if value != nil {
			m.SetUint32(key, uint32(*value))
		}
	}
	setBoolIfNotNil := func(key string, value *bool) {
		if value != nil {
			m.SetBoolean(key, *value)
		}
	}

	if cfg.Format != nil {
		m.SetUint32("format", uint32(*cfg.Format))
	}
	setIntIfNotNil("buffer.count", cfg.BufferCount)
	if cfg.BufferStrategy != nil {
		m.SetUint32("buffer.strategy", uint32(*cfg.BufferStrategy))
	}
	setIntIfNotNil("input", cfg.Input)
	setIntIfNotNil("channel", cfg.Channel)
	setIntIfNotNil("width", cfg.Width)
	setIntIfNotNil("height", cfg.Height)
	setIntIfNotNil("framerate", cfg.Framerate)
	setIntIfNotNil("compression", cfg.Compression)
	if cfg.Rotation != nil {
		m.SetUint32("rotation", uint32(*cfg.Rotation))
	}
	setBoolIfNotNil("horizontal_flip", cfg.HorizontalFlip)
	setBoolIfNotNil("vertical_flip", cfg.VerticalFlip)
	setBoolIfNotNil("monochrome", cfg.Monochrome)
	setBoolIfNotNil("dynamic.gop", cfg.DynamicGOP)
	setBoolIfNotNil("dynamic.bitrate", cfg.DynamicBitrate)
	setBoolIfNotNil("dynamic.framerate", cfg.DynamicFramerate)
	setBoolIfNotNil("dynamic.compression", cfg.DynamicCompression)
	setUint32IfNotNil("qp.i", cfg.Qpi)
	setUint32IfNotNil("qp.p", cfg.Qpp)
	setUint32IfNotNil("bitrate", cfg.Bitrate)
	if cfg.RateControlMode != nil {
		m.SetUint32("rc.mode", uint32(*cfg.RateControlMode))
	}
	if cfg.RateControlPriority != nil {
		m.SetUint32("rc.prio", uint32(*cfg.RateControlPriority))
	}
	setUint32IfNotNil("gop_length", cfg.GOPLength)
	if cfg.H264Profile != nil {
		m.SetUint32("h264.profile", uint32(*cfg.H264Profile))
	}
	if cfg.H265Profile != nil {
		m.SetUint32("h265.profile", uint32(*cfg.H265Profile))
	}
	setUint32IfNotNil("zip.strength", cfg.ZipStrength)
	setUint32IfNotNil("zip.max_gop_length", cfg.ZipMaxGOPLength)
	if cfg.ZipGOPMode != nil {
		m.SetUint32("zip.gop_mode", uint32(*cfg.ZipGOPMode))
	}
	if cfg.ZipFPSMode != nil {
		m.SetUint32("zip.fps_mode", uint32(*cfg.ZipFPSMode))
	}
	if cfg.ZipSkipMode != nil {
		m.SetUint32("zip.skip_mode", uint32(*cfg.ZipSkipMode))
	}
	setUint32IfNotNil("zip.min_fps_num", cfg.ZipMinFPSNum)
	setUint32IfNotNil("zip.min_fps_den", cfg.ZipMinFPSDen)
	if cfg.ZipProfile != nil {
		m.SetUint32("zip.profile", uint32(*cfg.ZipProfile))
	}
	setUint32IfNotNil("abr.target_bitrate", cfg.AbrTarget_bitrate)
	setUint32IfNotNil("abr.retention_time", cfg.AbrRetention_time)

	return m
}

func RestartStream(stream *VdoStream, stream_cfg VideoSteamConfiguration) (*VdoStream, error) {
	stream.Stop()
	stream.Unref()
	time.Sleep(time.Second * 2)
	return CreateAndStartStream(stream_cfg)
}

func CreateAndStartStream(stream_cfg VideoSteamConfiguration) (*VdoStream, error) {
	stream, err := NewVideoStreamFromConfig(stream_cfg)
	if err != nil {
		return nil, err
	}
	if err = stream.Start(); err != nil {
		return nil, err
	}
	return stream, nil
}
