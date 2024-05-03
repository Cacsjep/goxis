package axvdo

/*
#cgo pkg-config: vdostream
#include "vdo-types.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/vdostream/html/vdo-types_8h_source.html

// VdoWdrMode corresponds to the C enum VdoWdrMode.
type VdoWdrMode int

const (
	VdoWdrModeNone   VdoWdrMode = C.VDO_WDR_MODE_NONE
	VdoWdrModeLinear VdoWdrMode = C.VDO_WDR_MODE_LINEAR
	VdoWdrMode2X     VdoWdrMode = C.VDO_WDR_MODE_2X
	VdoWdrMode3X     VdoWdrMode = C.VDO_WDR_MODE_3X
	VdoWdrMode4X     VdoWdrMode = C.VDO_WDR_MODE_4X
	VdoWdrModeSensor VdoWdrMode = C.VDO_WDR_MODE_SENSOR
)

// VdoFormat corresponds to the C enum VdoFormat.
type VdoFormat int

const (
	VdoFormatNone      VdoFormat = C.VDO_FORMAT_NONE
	VdoFormatH264      VdoFormat = C.VDO_FORMAT_H264
	VdoFormatH265      VdoFormat = C.VDO_FORMAT_H265
	VdoFormatJPEG      VdoFormat = C.VDO_FORMAT_JPEG
	VdoFormatYUV       VdoFormat = C.VDO_FORMAT_YUV
	VdoFormatBayer     VdoFormat = C.VDO_FORMAT_BAYER
	VdoFormatIVS       VdoFormat = C.VDO_FORMAT_IVS
	VdoFormatRAW       VdoFormat = C.VDO_FORMAT_RAW
	VdoFormatRGBA      VdoFormat = C.VDO_FORMAT_RGBA
	VdoFormatRGB       VdoFormat = C.VDO_FORMAT_RGB
	VdoFormatPlanarRGB VdoFormat = C.VDO_FORMAT_PLANAR_RGB
)

// VdoFormatIsEncoded checks if the given format is one of the encoded formats (H264, H265, or JPEG).
func VdoFormatIsEncoded(format VdoFormat) bool {
	return format == VdoFormatH264 || format == VdoFormatH265 || format == VdoFormatJPEG
}

// VdoFormatIsMotionEncoded checks if the given format is one of the motion encoded formats.
// It is considered motion encoded if it is encoded but not JPEG.
func VdoFormatIsMotionEncoded(format VdoFormat) bool {
	return VdoFormatIsEncoded(format) && format != VdoFormatJPEG
}

// VdoH264Profile corresponds to the C enum VdoH264Profile.
type VdoH264Profile int

const (
	VdoH264ProfileNone     VdoH264Profile = C.VDO_H264_PROFILE_NONE
	VdoH264ProfileBaseline VdoH264Profile = C.VDO_H264_PROFILE_BASELINE
	VdoH264ProfileMain     VdoH264Profile = C.VDO_H264_PROFILE_MAIN
	VdoH264ProfileHigh     VdoH264Profile = C.VDO_H264_PROFILE_HIGH
)

// VdoH265Profile corresponds to the C enum VdoH265Profile.
type VdoH265Profile int

const (
	VdoH265ProfileNone   VdoH265Profile = C.VDO_H265_PROFILE_NONE
	VdoH265ProfileMain   VdoH265Profile = C.VDO_H265_PROFILE_MAIN
	VdoH265ProfileMain10 VdoH265Profile = C.VDO_H265_PROFILE_MAIN_10
)

// VdoRateControlMode corresponds to the C enum VdoRateControlMode.
type VdoRateControlMode int

const (
	VdoRateControlModeNone VdoRateControlMode = C.VDO_RATE_CONTROL_MODE_NONE
	VdoRateControlModeCBR  VdoRateControlMode = C.VDO_RATE_CONTROL_MODE_CBR
	VdoRateControlModeVBR  VdoRateControlMode = C.VDO_RATE_CONTROL_MODE_VBR
	VdoRateControlModeMBR  VdoRateControlMode = C.VDO_RATE_CONTROL_MODE_MBR
	VdoRateControlModeABR  VdoRateControlMode = C.VDO_RATE_CONTROL_MODE_ABR
)

// VdoRateControlPriority corresponds to the C enum VdoRateControlPriority.
type VdoRateControlPriority int

const (
	VdoRateControlPriorityNone          VdoRateControlPriority = C.VDO_RATE_CONTROL_PRIORITY_NONE
	VdoRateControlPriorityFramerate     VdoRateControlPriority = C.VDO_RATE_CONTROL_PRIORITY_FRAMERATE
	VdoRateControlPriorityQuality       VdoRateControlPriority = C.VDO_RATE_CONTROL_PRIORITY_QUALITY
	VdoRateControlPriorityFullFramerate VdoRateControlPriority = C.VDO_RATE_CONTROL_PRIORITY_FULL_FRAMERATE
)

// VdoRateControlPriority corresponds to the C enum VdoRateControlPriority.
type VdoFrameType int

const (
	VdoFrameTypeNone    VdoFrameType = C.VDO_FRAME_TYPE_NONE
	VdoFrameTypeH264SPS VdoFrameType = C.VDO_FRAME_TYPE_H264_SPS
	VdoFrameTypeH264PPS VdoFrameType = C.VDO_FRAME_TYPE_H264_PPS
	VdoFrameTypeH264SEI VdoFrameType = C.VDO_FRAME_TYPE_H264_SEI
	VdoFrameTypeH264IDR VdoFrameType = C.VDO_FRAME_TYPE_H264_IDR
	VdoFrameTypeH264I   VdoFrameType = C.VDO_FRAME_TYPE_H264_I
	VdoFrameTypeH264P   VdoFrameType = C.VDO_FRAME_TYPE_H264_P
	VdoFrameTypeH264B   VdoFrameType = C.VDO_FRAME_TYPE_H264_B
	VdoFrameTypeH265SPS VdoFrameType = C.VDO_FRAME_TYPE_H265_SPS
	VdoFrameTypeH265PPS VdoFrameType = C.VDO_FRAME_TYPE_H265_PPS
	VdoFrameTypeH265VPS VdoFrameType = C.VDO_FRAME_TYPE_H265_VPS
	VdoFrameTypeH265SEI VdoFrameType = C.VDO_FRAME_TYPE_H265_SEI
	VdoFrameTypeH265IDR VdoFrameType = C.VDO_FRAME_TYPE_H265_IDR
	VdoFrameTypeH265I   VdoFrameType = C.VDO_FRAME_TYPE_H265_I
	VdoFrameTypeH265P   VdoFrameType = C.VDO_FRAME_TYPE_H265_P
	VdoFrameTypeH265B   VdoFrameType = C.VDO_FRAME_TYPE_H265_B
	VdoFrameTypeJPEG    VdoFrameType = C.VDO_FRAME_TYPE_JPEG
	VdoFrameTypeYUV     VdoFrameType = C.VDO_FRAME_TYPE_YUV
	VdoFrameTypeRAW     VdoFrameType = C.VDO_FRAME_TYPE_RAW
	VdoFrameTypeRGBA    VdoFrameType = C.VDO_FRAME_TYPE_RGBA
	VdoFrameTypeRGB     VdoFrameType = C.VDO_FRAME_TYPE_RGB
	VdoFrameTypePlanar  VdoFrameType = C.VDO_FRAME_TYPE_PLANAR_RGB
)

func (v VdoFrameType) String() string {
	switch v {
	case VdoFrameTypeNone:
		return "None"
	case VdoFrameTypeH264SPS:
		return "H264SPS"
	case VdoFrameTypeH264PPS:
		return "H264PPS"
	case VdoFrameTypeH264SEI:
		return "H264SEI"
	case VdoFrameTypeH264IDR:
		return "H264IDR"
	case VdoFrameTypeH264I:
		return "H264I"
	case VdoFrameTypeH264P:
		return "H264P"
	case VdoFrameTypeH264B:
		return "H264B"
	case VdoFrameTypeH265SPS:
		return "H265SPS"
	case VdoFrameTypeH265PPS:
		return "H265PPS"
	case VdoFrameTypeH265VPS:
		return "H265VPS"
	case VdoFrameTypeH265SEI:
		return "H265SEI"
	case VdoFrameTypeH265IDR:
		return "H265IDR"
	case VdoFrameTypeH265I:
		return "H265I"
	case VdoFrameTypeH265P:
		return "H265P"
	case VdoFrameTypeH265B:
		return "H265B"
	case VdoFrameTypeJPEG:
		return "JPEG"
	case VdoFrameTypeYUV:
		return "YUV"
	case VdoFrameTypeRAW:
		return "RAW"
	case VdoFrameTypeRGBA:
		return "RGBA"
	case VdoFrameTypeRGB:
		return "RGB"
	case VdoFrameTypePlanarRGB:
		return "PlanarRGB"
	default:
		return fmt.Sprintf("Unknown(%d)", v)
	}
}

type VdoZipStreamProfile int

const (
	VdoZipStreamProfileNone    VdoZipStreamProfile = -1
	VdoZipStreamProfileClassic VdoZipStreamProfile = 0
	VdoZipStreamProfileStorage VdoZipStreamProfile = 1
	VdoZipStreamProfileLive    VdoZipStreamProfile = 2
)

type VdoChunkType uint

const (
	VdoChunkNone  VdoChunkType = 0
	VdoChunkError VdoChunkType = 1 << 31
)

// VdoChunk represents a chunk of data with a specified type.
type VdoChunk struct {
	Data unsafe.Pointer // Pointer to data
	Size uintptr        // Size of the data in bytes
	Type VdoChunkType   // Type of the chunk, based on the VdoChunkType enum
}

// VdoFrameIsEncoded checks if the given frame type is one of the encoded types.
func VdoFrameIsEncoded(frameType VdoFrameType) bool {
	return frameType >= VdoFrameTypeH264SPS && frameType <= VdoFrameTypeJPEG
}

// VdoFrameIsOfFormat returns the VdoFormat corresponding to the given VdoFrameType.
func VdoFrameIsOfFormat(frameType VdoFrameType) VdoFormat {
	switch frameType {
	case VdoFrameTypeH264SPS, VdoFrameTypeH264PPS, VdoFrameTypeH264SEI,
		VdoFrameTypeH264IDR, VdoFrameTypeH264I, VdoFrameTypeH264P, VdoFrameTypeH264B:
		return VdoFormatH264

	case VdoFrameTypeH265SPS, VdoFrameTypeH265PPS, VdoFrameTypeH265VPS,
		VdoFrameTypeH265SEI, VdoFrameTypeH265IDR, VdoFrameTypeH265I, VdoFrameTypeH265P, VdoFrameTypeH265B:
		return VdoFormatH265

	case VdoFrameTypeJPEG:
		return VdoFormatJPEG

	case VdoFrameTypeYUV:
		return VdoFormatYUV

	case VdoFrameTypeRAW:
		return VdoFormatRAW

	case VdoFrameTypeRGBA:
		return VdoFormatRGBA

	case VdoFrameTypeRGB:
		return VdoFormatRGB

	case VdoFrameTypePlanarRGB:
		return VdoFormatPlanarRGB

	default:
		return VdoFormatNone
	}
}

type VdoOverlayAlign int

const (
	VdoOverlayAlignNone   VdoOverlayAlign = -1
	VdoOverlayAlignTop    VdoOverlayAlign = 0
	VdoOverlayAlignBottom VdoOverlayAlign = 1
)

type VdoOverlayColor uint16

const (
	VdoOverlayColorTransparent VdoOverlayColor = 0x0000
	VdoOverlayColorBlack       VdoOverlayColor = 0xF000
	VdoOverlayColorWhite       VdoOverlayColor = 0xFFFF
)

type VdoOverlayTextSize int

const (
	VdoOverlayTextSizeSmall  VdoOverlayTextSize = 16
	VdoOverlayTextSizeMedium VdoOverlayTextSize = 32
	VdoOverlayTextSizeLarge  VdoOverlayTextSize = 48
)

type VdoStreamTimestamp uint

const (
	VdoTimestampNone                  VdoStreamTimestamp = 0
	VdoTimestampUTC                   VdoStreamTimestamp = 1
	VdoTimestampZipstream             VdoStreamTimestamp = 2
	VdoTimestampDiff                  VdoStreamTimestamp = 4
	VdoTimestampMonoCapture           VdoStreamTimestamp = 8
	VdoTimestampMonoServer            VdoStreamTimestamp = 16
	VdoTimestampMonoClient            VdoStreamTimestamp = 32
	VdoTimestampMonoClientServerDiff  VdoStreamTimestamp = VdoTimestampDiff | VdoTimestampMonoClient | VdoTimestampMonoServer
	VdoTimestampMonoClientCaptureDiff VdoStreamTimestamp = VdoTimestampDiff | VdoTimestampMonoClient | VdoTimestampMonoCapture
)

type VdoIntent uint

const (
	VdoIntentNone     VdoIntent = 0
	VdoIntentControl  VdoIntent = 1
	VdoIntentMonitor  VdoIntent = 2
	VdoIntentConsume  VdoIntent = 4
	VdoIntentProduce  VdoIntent = 8
	VdoIntentDefault  VdoIntent = VdoIntentConsume | VdoIntentControl
	VdoIntentEventFD  VdoIntent = 16
	VdoIntentUniverse VdoIntent = ^VdoIntent(0)
)

type VdoStreamEvent uint

const (
	VdoStreamEventNone      VdoStreamEvent = 0x00
	VdoStreamEventStarted   VdoStreamEvent = 0x01
	VdoStreamEventStopped   VdoStreamEvent = 0x02
	VdoStreamEventResource  VdoStreamEvent = 0x10
	VdoStreamEventQuotaSoft VdoStreamEvent = 0x11
	VdoStreamEventQuotaHard VdoStreamEvent = 0x12
	VdoStreamEventZipstream VdoStreamEvent = 0x20
	VdoStreamEventInvalid   VdoStreamEvent = ^VdoStreamEvent(0)
)

type VdoBufferAccess uint

const (
	VdoBufferAccessNone  VdoBufferAccess = 0
	VdoBufferAccessCPURd VdoBufferAccess = 1 << 0
	VdoBufferAccessDEVRd VdoBufferAccess = 1 << 1
	VdoBufferAccessAnyRd VdoBufferAccess = VdoBufferAccessCPURd | VdoBufferAccessDEVRd
	VdoBufferAccessCPUWr VdoBufferAccess = 1 << 8
	VdoBufferAccessDEVWr VdoBufferAccess = 1 << 9
	VdoBufferAccessAnyWr VdoBufferAccess = VdoBufferAccessCPUWr | VdoBufferAccessDEVWr
	VdoBufferAccessCPURW VdoBufferAccess = VdoBufferAccessCPURd | VdoBufferAccessCPUWr
	VdoBufferAccessDEVRW VdoBufferAccess = VdoBufferAccessDEVRd | VdoBufferAccessDEVWr
	VdoBufferAccessAnyRW VdoBufferAccess = VdoBufferAccessCPURW | VdoBufferAccessDEVRW
)

type VdoBufferStrategy int

const (
	VdoBufferStrategyNone     VdoBufferStrategy = 0
	VdoBufferStrategyInput    VdoBufferStrategy = 1
	VdoBufferStrategyExternal VdoBufferStrategy = 2
	VdoBufferStrategyExplicit VdoBufferStrategy = 3
	VdoBufferStrategyInfinite VdoBufferStrategy = 4
)

// VdoMemChunk represents a memory chunk with its data pointer and size.
type VdoMemChunk struct {
	Data     unsafe.Pointer // Pointer to the data
	DataSize uintptr        // Size of the data in bytes
}

// VdoResolution represents the resolution of a video or an image.
type VdoResolution struct {
	Width  int
	Height int
}

func (r VdoResolution) RgbSize() int {
	return r.Width * r.Height * 3
}

// VdoRect represents a rectangular area with width, height, and position (x, y).
type VdoRect struct {
	Width  uint // Width of the rectangle
	Height uint // Height of the rectangle
	X      uint // X coordinate of the rectangle's origin
	Y      uint // Y coordinate of the rectangle's origin
}
