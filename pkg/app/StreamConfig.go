package app

import "github.com/Cacsjep/goxis/pkg/axvdo"

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
	Format              *axvdo.VdoFormat              //Video stream format as VdoFormat.
	BufferCount         *int                          // How many in-flight buffers are desired.
	BufferStrategy      *axvdo.VdoBufferStrategy      // Buffering Strategy as VdoBufferStrategy.
	Input               *int                          // Video input, 1 ... inmax. 0 is invalid. No view areas.
	Channel             *int                          // Video channel, 0 is overview, 1, 2, ... are view areas.
	Width               *int                          // Video stream horizontal resolution.
	Height              *int                          // Video stream vertical resolution.
	Framerate           *int                          // Video stream vertical resolution.
	Compression         *int                          // Video stream compression, Axis standard range [0:100]
	Rotation            *StreamRotation               // Video stream rotation, normally [0,90,180,270].
	HorizontalFlip      *bool                         // Video stream horizontal flip (mirroring).
	VerticalFlip        *bool                         // Video stream vertical flip.
	Monochrome          *bool                         // Video stream monochrome encoding.
	DynamicGOP          *bool                         // Enable dynamic gop
	DynamicBitrate      *bool                         // Enable dynamic bitrate
	DynamicFramerate    *bool                         // Enable dynamic framerate
	DynamicCompression  *bool                         // Enable dynamic compression
	Qpi                 *uint32                       // QP value for I-frames.
	Qpp                 *uint32                       // QP value for P-frames.
	Bitrate             *uint32                       // Video stream bitrate (bps)
	RateControlMode     *axvdo.VdoRateControlMode     // Bitrate control mode.
	RateControlPriority *axvdo.VdoRateControlPriority // Bitrate control priority.
	GOPLength           *uint32                       // GOP length.
	// H.264 Specific Settings
	H264Profile *axvdo.VdoH264Profile // H.265 profile as VdoH264Profile
	H265Profile *axvdo.VdoH265Profile // H.264 profile as VdoH264Profile
	// Zipstream Specific Settings
	ZipStrength     *uint32                    // Zipstream strength.
	ZipMaxGOPLength *uint32                    // Zipstream maximum GOP length.
	ZipGOPMode      *ZipGopMode                // Zipstream GOP mode [0 = fixed, 1 = dynamic].
	ZipFPSMode      *ZipFPSMode                // Zipstream framerate control mode: [0 = fixed, 1 = dynamic].
	ZipSkipMode     *ZipSkipMode               // Zipstream frame skip mode: [0 = drop, 1 = empty].
	ZipMinFPSNum    *uint32                    // Zipstream minimum framerate numerator.
	ZipMinFPSDen    *uint32                    // Zipstream minimum framerate denominator.
	ZipProfile      *axvdo.VdoZipStreamProfile // Zipstream profile.
	// ABR Specific Settings
	// The following ABR specific settings are supported with VDO_RATE_CONTROL_MODE_ABR.
	AbrTarget_bitrate *uint32 // Stream target bitrate (bps)
	AbrRetention_time *uint32 // Retention time in seconds
}

func VideoStreamConfigToVdoMap(cfg VideoSteamConfiguration) *axvdo.VdoMap {
	m := axvdo.NewVdoMap()

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
