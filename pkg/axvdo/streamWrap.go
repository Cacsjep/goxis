package axvdo

import (
	"fmt"
	"reflect"
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

func (vsc *VideoSteamConfiguration) RgbFrameSize() int {
	return *vsc.Width * *vsc.Height * 3
}

// StreamStats holds statistical information about a video stream, including metrics such as bitrate, frame rate, resolution, and more.
// It provides a detailed view of the stream's current configuration and performance.
type StreamStats struct {
	// Bitrate of the video stream in bits per second.
	Bitrate uint32
	// BufferType describes the type of buffering strategy used for the stream.
	BufferType string
	// Channel represents the specific channel or stream identifier.
	Channel uint32
	// Format specifies the video format used in the stream.
	Format int16
	// Framerate is the number of frames per second in the video stream.
	Framerate uint32
	// GOPLength represents the length of a Group of Pictures in frames.
	GOPLength uint32
	// H26xIntraRefresh indicates the intra refresh rate for H.264/H.265 video streams.
	H26xIntraRefresh uint32
	// Height of the video stream in pixels.
	Height uint32
	// HorizontalFlip indicates whether the video stream is flipped horizontally.
	HorizontalFlip bool
	// ID is a unique identifier for the stream.
	ID uint32
	// InitialBitrate is the initial bitrate of the video stream in bits per second.
	InitialBitrate uint32
	// InitialQPb is the initial quantization parameter for B-frames.
	InitialQPb uint32
	// InitialQPi is the initial quantization parameter for I-frames.
	InitialQPi uint32
	// InitialQPp is the initial quantization parameter for P-frames.
	InitialQPp uint32
	// Overlays describes any overlays applied to the video stream.
	Overlays string
	// Peers represents the number of peers connected to the video stream.
	Peers uint32
	// QPb is the current quantization parameter for B-frames.
	QPb uint32
	// QPi is the current quantization parameter for I-frames.
	QPi uint32
	// QPp is the current quantization parameter for P-frames.
	QPp uint32
	// Rotation indicates the rotation applied to the video stream, in degrees.
	Rotation uint32
	// Running indicates whether the video stream is currently active.
	Running bool
	// SquarePixel indicates whether square pixels are used in the video stream.
	SquarePixel uint32
	// StatisticsAccumulatedBytes is the total number of bytes processed by the stream.
	StatisticsAccumulatedBytes uint64
	// StatisticsAccumulatedIDRBytes is the total number of IDR frame bytes processed by the stream.
	StatisticsAccumulatedIDRBytes uint64
	// StatisticsBitCount is the total number of bits processed by the stream.
	StatisticsBitCount uint32
	// StatisticsBitrate is the current bitrate of the video stream in bits per second.
	StatisticsBitrate uint32
	// StatisticsDuration is the total duration of video processed by the stream, in milliseconds.
	StatisticsDuration int64
	// StatisticsDynamicFramerate is the current dynamic frame rate of the video stream.
	StatisticsDynamicFramerate uint32
	// StatisticsFailedFrames is the total number of frames that failed to process.
	StatisticsFailedFrames uint32
	// StatisticsFrameCount is the total number of frames processed by the stream.
	StatisticsFrameCount uint32
	// StatisticsFramerate is the current frame rate of the video stream.
	StatisticsFramerate uint32
	// StatisticsIDRFrameCount is the total number of IDR frames processed by the stream.
	StatisticsIDRFrameCount uint32
	// StatisticsLastFrameTS is the timestamp of the last frame processed.
	StatisticsLastFrameTS uint64
	// StatisticsReclaimCount is the number of times stream resources have been reclaimed.
	StatisticsReclaimCount uint32
	// Width of the video stream in pixels.
	Width uint32
	// ZipProfile indicates the compression profile used for the video stream.
	ZipProfile int16
}

// PrintStreamStats prints the fields of the StreamStats.
func (s *StreamStats) PrintStreamStats() {
	val := reflect.Indirect(reflect.ValueOf(s))
	typeOfStats := val.Type()

	fmt.Println("StreamStats:")
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typeOfStats.Field(i).Name
		fieldValue := field.Interface()
		fmt.Printf("  %s: %v\n", fieldName, fieldValue)
	}
}

// VideoFrame represents a single frame of video data, including metadata such as the sequence number, timestamp, and size.
// It also includes information about the type of frame and any errors encountered.
type VideoFrame struct {
	SequenceNbr   uint         // The sequence number of the frame.
	Timestamp     time.Time    // The timestamp when the frame was captured.
	Size          uint         // The size of the frame data in bytes.
	Data          []byte       // The raw data of the video frame
	Type          VdoFrameType // Type describes the frame type (e.g., I-frame, P-frame, B-frame).
	Error         error        // Error contains any error that occurred while processing the frame.
	ErrorExpected bool         // ErrorExpected indicates whether the error was expected on vdo maintance
	HeaderSize    int          // HeaderSize is the size of the frame header
}

// String returns a string representation of the VideoFrame, including sequence number, timestamp, size, and frame type.
func (f *VideoFrame) String() string {
	return fmt.Sprintf("SequenceNbr: %d, Timestamp: %s, Size: %d, Type: %s",
		f.SequenceNbr, f.Timestamp.Format("2006-01-02 15:04:05"), f.Size, f.Type.String())
}

func (f *VideoFrame) HeaderData() []byte {
	return f.Data[:f.HeaderSize]
}

// NewVideoFrame creates a new VideoFrame instance from a VdoFrame and its data.
// This function extracts relevant information from the VdoFrame, including the sequence number, timestamp, size, and frame type, and packages it into a VideoFrame structure.
func NewVideoFrame(frame *VdoFrame, data []byte, header_size int) *VideoFrame {
	return &VideoFrame{
		SequenceNbr: frame.GetSequenceNbr(),
		Timestamp:   time.Unix(frame.GetCustomTimestamp()/1000000, (frame.GetCustomTimestamp()%1000000)*1000),
		Size:        frame.GetSize(),
		Type:        frame.GetFrameType(),
		Data:        data,
		HeaderSize:  header_size,
	}
}

// NewVideoStreamFromConfig creates and initializes a new video stream based on the provided configuration.
// It converts the configuration into a format understood by the underlying video streaming system and starts the stream.
func NewVideoStreamFromConfig(stream_cfg VideoSteamConfiguration) (*VdoStream, error) {
	vdoMap := VideoStreamConfigToVdoMap(stream_cfg)
	defer vdoMap.Unref()
	return NewStream(vdoMap)
}

// GetVideoFrame retrieves a video frame from a video stream.
// If an expected error occurs (e.g., for stream maintenance), the function returns a VideoFrame with the Error and ErrorExpected fields set.
// This function is recursive and will retry fetching a frame if a recoverable error occurs.
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

	header_size := vdo_frame.GetHeaderSize()
	defer vdo_stream.BufferUnref(vdo_buf)
	return NewVideoFrame(vdo_frame, buff_data, header_size)
}

// VideoStreamConfigToVdoMap converts a VideoSteamConfiguration object into a VdoMap.
// This conversion allows the configuration to be applied to the video stream by translating high-level configuration options into the format expected by the streaming system.
func VideoStreamConfigToVdoMap(cfg VideoSteamConfiguration) *VdoMap {
	m := NewVdoMap()

	// Utility functions to streamline conditional setting of configuration parameters
	setUint32IfNotNil := func(key string, value *uint32) {
		if value != nil {
			m.SetUint32(key, *value)
		}
	}
	setIntIfNotNil := func(key string, value *int) {
		if value != nil {
			m.SetInt16(key, int16(*value))
		}
	}
	setBoolIfNotNil := func(key string, value *bool) {
		if value != nil {
			m.SetBoolean(key, *value)
		}
	}

	if cfg.Format != nil {
		m.SetInt16("format", int16(*cfg.Format))
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
		m.SetInt16("rc.mode", int16(*cfg.RateControlMode))
	}
	if cfg.RateControlPriority != nil {
		m.SetInt16("rc.prio", int16(*cfg.RateControlPriority))
	}
	setUint32IfNotNil("gop_length", cfg.GOPLength)
	if cfg.H264Profile != nil {
		m.SetInt16("h264.profile", int16(*cfg.H264Profile))
	}
	if cfg.H265Profile != nil {
		m.SetInt16("h265.profile", int16(*cfg.H265Profile))
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
		m.SetInt16("zip.profile", int16(*cfg.ZipProfile))
	}
	setUint32IfNotNil("abr.target_bitrate", cfg.AbrTarget_bitrate)
	setUint32IfNotNil("abr.retention_time", cfg.AbrRetention_time)

	return m
}

// GetChannel extracts the channel number from a VideoSteamConfiguration if it's set, otherwise returns a default value of 0.
// This method simplifies accessing the channel value, providing a safe way to handle nil pointers within the configuration.
func (vsc *VideoSteamConfiguration) GetChannel() int {
	if vsc.Channel != nil {
		return *vsc.Channel
	}
	return 0
}
