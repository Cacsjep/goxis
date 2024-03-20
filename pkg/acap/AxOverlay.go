package acap

/*
#cgo pkg-config: gio-2.0 glib-2.0 cairo axoverlay
#include <axoverlay.h>
#include <cairo/cairo.h>
*/
import "C"

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axoverlay/html/axoverlaypage.html

// axoverlay_error_code mirrors the error codes from the C API.
type AxOverlayErrorCode int

const (
	AxOverlayErrorInvalidValue       AxOverlayErrorCode = 1000
	AxOverlayErrorInternal           AxOverlayErrorCode = 2000
	AxOverlayErrorUnexpected         AxOverlayErrorCode = 3000
	AxOverlayErrorGeneric            AxOverlayErrorCode = 4000
	AxOverlayErrorInvalidArgument    AxOverlayErrorCode = 5000
	AxOverlayErrorServiceUnavailable AxOverlayErrorCode = 6000
	AxOverlayErrorBackend            AxOverlayErrorCode = 7000
)

// axoverlay_colorspace defines color space types similar to the C enumeration.
type AxOverlayColorspace int

const (
	AxOverlayColorspaceARGB32 AxOverlayColorspace = iota
	AxOverlayColorspace4BitPalette
	AxOverlayColorspace1BitPalette
	AxOverlayColorspaceUndefined
)

// axoverlay_position_type defines overlay position types.
type AxOverlayPositionType int

const (
	AxOverlayTopLeft AxOverlayPositionType = iota
	AxOverlayTopRight
	AxOverlayBottomLeft
	AxOverlayBottomRight
	AxOverlayCustomNormalized
	AxOverlayCustomSource
)

// axoverlay_backend_type defines backend types for the overlay system.
type AxOverlayBackendType int

const (
	AxOverlayCairoImageBackend AxOverlayBackendType = 1
	AxOverlayOpenGLESBackend
	AxOverlayOpenBackend
)

// axoverlay_anchor_point defines anchor points for overlays.
type AxOverlayAnchorPoint int

const (
	AxOverlayAnchorTopLeft AxOverlayAnchorPoint = iota
	AxOverlayAnchorCenter
)

// axoverlay_stream_type defines stream types for overlay input.
type AxOverlayStreamType int

const (
	AxOverlayStreamJPEG AxOverlayStreamType = iota
	AxOverlayStreamH264
	AxOverlayStreamH265
	AxOverlayStreamYCbCr
	AxOverlayStreamVOUT
	AxOverlayStreamOther
)

// axoverlay_stream_data is a struct to hold stream data.
type AxOverlayStreamData struct {
	ID         int
	Camera     int
	Width      int
	Height     int
	Rotation   int
	IsMirrored bool
	Type       AxOverlayStreamType
}

// axoverlay_overlay_data is a struct to hold overlay configuration.
type AxOverlayOverlayData struct {
	AnchorPoint   AxOverlayAnchorPoint
	PositionType  AxOverlayPositionType
	X, Y          float32
	Width, Height int
	ZPriority     int
	Colorspace    AxOverlayColorspace
	ScaleToStream bool
	Ptr           *C.struct_axoverlay_overlay_data
}

// axoverlay_palette_color defines a color in the overlay palette.
type AxOverlayPaletteColor struct {
	Red, Green, Blue, Alpha byte
	Pixelate                bool
}

// Callback function types for stream selection, adjustment, and rendering.
type (
	AxOverlayStreamSelectFunc func(camera, width, height, rotation int, isMirrored bool, streamType AxOverlayStreamType) bool
	AxOverlayAdjustmentFunc   func(id int, stream *AxOverlayStreamData, positionType *AxOverlayPositionType, x, y *float32, width, height *int, userData interface{})
	AxOverlayRenderFunc       func(renderingContext interface{}, id int, stream *AxOverlayStreamData, positionType AxOverlayPositionType, x, y float32, width, height int, userData interface{})
)

// axoverlay_settings is a struct to hold overlay settings.
type AxOverlaySettings struct {
	RenderCallback     AxOverlayRenderFunc
	AdjustmentCallback AxOverlayAdjustmentFunc
	SelectCallback     AxOverlayStreamSelectFunc
	Backend            AxOverlayBackendType
	Ptr                *C.struct_axoverlay_settings
}

// AxoverlayInitAxoverlaySettings initializes axoverlay_settings with default values.
func AxOverlayInitAxoverlaySettings(settings *AxOverlaySettings) {
	C.axoverlay_init_axoverlay_settings(settings.Ptr)
}

// AxoverlayInit initializes the axoverlay system with specified settings.
func AxOverlayInit(settings *AxOverlaySettings) error {
	var gerr *C.GError
	C.axoverlay_init(settings.Ptr, &gerr)
	return newGError(gerr)
}

// axoverlayCleanup frees up allocated resources.
func AxOverlayCleanup() {
	C.axoverlay_cleanup()
}

// axoverlayReloadStreams reloads all stream information.
func AxOverlayReloadStreams() error {
	var gerr *C.GError
	C.axoverlay_reload_streams(&gerr)
	return newGError(gerr)
}

// axoverlayRedraw signals the system that a redraw should be done.
func AxOverlayRedraw() error {
	var gerr *C.GError
	C.axoverlay_redraw(&gerr)
	return newGError(gerr)
}

// axoverlayInitOverlayData initializes an axoverlay_overlay_data struct with default values.
func AxOverlayInitOverlayData(data *AxOverlayOverlayData) {
	C.axoverlay_init_overlay_data(data.Ptr)
}

// axoverlayCreateOverlay creates an overlay with the specified data.
func AxOverlayCreateOverlay(data *AxOverlayOverlayData) (int, error) {
	var gerr *C.GError
	id := C.axoverlay_create_overlay(data.Ptr, nil, &gerr)
	return int(id), newGError(gerr)
}

// axoverlayDestroyOverlay destroys the overlay with the given ID.
func AxOverlayDestroyOverlay(id int) error {
	var gerr *C.GError
	C.axoverlay_destroy_overlay(C.gint(id), &gerr)
	return newGError(gerr)
}

// axoverlaySetOverlayPosition updates the position of an existing overlay.
func AxOverlaySetOverlayPosition(id int, positionType AxOverlayPositionType, x, y float32) error {
	var gerr *C.GError
	C.axoverlay_set_overlay_position(C.gint(id), C.enum_axoverlay_position_type(positionType), C.gfloat(x), C.gfloat(y), &gerr)
	return newGError(gerr)
}

// axoverlaySetOverlaySize updates the size of an existing overlay.
func AxOverlaySetOverlaySize(id, width, height int) error {
	var gerr *C.GError
	C.axoverlay_set_overlay_size(C.gint(id), C.gint(width), C.gint(height), &gerr)
	return newGError(gerr)
}

// axoverlayGetMaxResolutionWidth reads the maximum resolution width for a camera.
func AxOverlayGetMaxResolutionWidth(camera int) (int, error) {
	var gerr *C.GError
	width := C.axoverlay_get_max_resolution_width(C.gint(camera), &gerr)
	return int(width), newGError(gerr)
}

// axoverlayGetMaxResolutionHeight reads the maximum resolution height for a camera.
func AxOverlayGetMaxResolutionHeight(camera int) (int, error) {
	var gerr *C.GError
	height := C.axoverlay_get_max_resolution_height(C.gint(camera), &gerr)
	return int(height), newGError(gerr)
}

// axoverlayIsBackendSupported checks if a specified backend is supported.
func AxOverlayIsBackendSupported(backend AxOverlayBackendType) bool {
	return ctoGoBoolean(C.axoverlay_is_backend_supported(C.enum_axoverlay_backend_type(backend)))
}

// axoverlayGetNumberOfPaletteColors returns the number of palette colors.
func AxOverlayGetNumberOfPaletteColors() (int, error) {
	var gerr *C.GError
	count := C.axoverlay_get_number_of_palette_colors(&gerr)
	return int(count), newGError(gerr)
}

// axoverlayGetPaletteColor retrieves a palette color by index.
func AxOverlayGetPaletteColor(index int) (AxOverlayPaletteColor, error) {
	var gerr *C.GError
	var color C.struct_axoverlay_palette_color
	C.axoverlay_get_palette_color(C.gint(index), &color, &gerr)
	err := newGError(gerr)
	if err != nil {
		return AxOverlayPaletteColor{}, err
	}
	return AxOverlayPaletteColor{
		Red:      uint8(color.red),
		Green:    uint8(color.green),
		Blue:     uint8(color.blue),
		Alpha:    uint8(color.alpha),
		Pixelate: ctoGoBoolean(color.pixelate),
	}, nil
}

// axoverlaySetPaletteColor sets a palette color by index.
func AxOverlaySetPaletteColor(index int, color AxOverlayPaletteColor) error {
	var gerr *C.GError
	cColor := C.struct_axoverlay_palette_color{
		red:      C.guchar(color.Red),
		green:    C.guchar(color.Green),
		blue:     C.guchar(color.Blue),
		alpha:    C.guchar(color.Alpha),
		pixelate: goBooleanToC(color.Pixelate),
	}
	C.axoverlay_set_palette_color(C.gint(index), &cColor, &gerr)
	return newGError(gerr)
}
