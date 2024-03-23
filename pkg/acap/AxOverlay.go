package acap

/*
#cgo LDFLAGS: -laxoverlay
#cgo pkg-config: gio-2.0 glib-2.0 cairo axoverlay
#include <axoverlay.h>
#include <cairo/cairo.h>

extern gboolean GoAxOverlayStreamSelectFunction(gint camera, gint width, gint height, gint rotation, gboolean is_mirrored, enum axoverlay_stream_type type);
extern void GoAxOverlayAdjustmentFunction(gint id, struct axoverlay_stream_data *stream, enum axoverlay_position_type *postype, gfloat *overlay_x, gfloat *overlay_y, gint *overlay_width, gint *overlay_height, gpointer user_data);
extern void GoAxOverlayRenderFunction(gpointer rendering_context, gint id, struct axoverlay_stream_data *stream, enum axoverlay_position_type postype, gfloat overlay_x, gfloat overlay_y, gint overlay_width, gint overlay_height, gpointer user_data);
*/
import "C"
import (
	"fmt"
	"runtime/cgo"
	"unsafe"
)

// https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/src/api/axoverlay/html/axoverlaypage.html

var streamSelectCallback AxOverlayStreamSelectFunc
var adjustmentCallback AxOverlayAdjustmentFunc
var renderCallback AxOverlayRenderFunc
var overlayUserDataHandle cgo.Handle

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

func (opt AxOverlayPositionType) String() string {
	switch opt {
	case AxOverlayTopLeft:
		return "TopLeft"
	case AxOverlayTopRight:
		return "TopRight"
	case AxOverlayBottomLeft:
		return "BottomLeft"
	case AxOverlayBottomRight:
		return "BottomRight"
	case AxOverlayCustomNormalized:
		return "CustomNormalized"
	case AxOverlayCustomSource:
		return "CustomSource"
	default:
		return "Unkown"
	}
}

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

func (oap AxOverlayAnchorPoint) String() string {
	if oap == AxOverlayAnchorTopLeft {
		return "TopLeft"
	} else {
		return "Center"
	}
}

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
	Ptr        *C.struct_axoverlay_stream_data
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
	ptr           *C.struct_axoverlay_overlay_data
}

func AxOverlayDataInitalze(overlay_data *AxOverlayOverlayData) {
	overlay_data.ptr = (*C.struct_axoverlay_overlay_data)(C.malloc(C.size_t(unsafe.Sizeof(*overlay_data.ptr))))
	if overlay_data.ptr == nil {
		panic("Failed to allocate memory for axoverlay_overlay_data")
	}
	AxOverlayInitOverlayData(overlay_data)
	overlay_data.ptr.postype = C.enum_axoverlay_anchor_point(overlay_data.PositionType)
	overlay_data.ptr.anchor_point = C.enum_axoverlay_position_type(overlay_data.AnchorPoint)
	overlay_data.ptr.colorspace = C.enum_axoverlay_colorspace(overlay_data.Colorspace)
	overlay_data.ptr.x = C.gfloat(overlay_data.X)
	overlay_data.ptr.y = C.gfloat(overlay_data.Y)
	overlay_data.ptr.width = C.gint(overlay_data.Width)
	overlay_data.ptr.height = C.gint(overlay_data.Height)
	overlay_data.ptr.scale_to_stream = goBooleanToC(overlay_data.ScaleToStream)
}

func (s *AxOverlayOverlayData) Free() {
	C.free(unsafe.Pointer(s.ptr))
}

// axoverlayInitOverlayData initializes an axoverlay_overlay_data struct with default values.
func AxOverlayInitOverlayData(data *AxOverlayOverlayData) {
	C.axoverlay_init_overlay_data(data.ptr)
}

// axoverlay_palette_color defines a color in the overlay palette.
type AxOverlayPaletteColor struct {
	R, G, B, A byte
	Pixelate   bool
}

// Callback function types for stream selection, adjustment, and rendering.
type (
	AxOverlayStreamSelectFunc func(camera, width, height, rotation int, isMirrored bool, streamType AxOverlayStreamType) bool
	AxOverlayAdjustmentFunc   func(id int, stream *AxOverlayStreamData, positionType *AxOverlayPositionType, x, y *float32, width, height *int, userData any)
	AxOverlayRenderFunc       func(cairoCtx *CairoContext, id int, stream *AxOverlayStreamData, positionType AxOverlayPositionType, x, y float32, width, height int, userData any)
)

// axoverlay_settings is a struct to hold overlay settings.
type AxOverlaySettings struct {
	backend AxOverlayBackendType
	ptr     *C.struct_axoverlay_settings
}

// AxoverlayInitAxoverlaySettings initializes axoverlay_settings with default values.
func NewAxOverlaySettings(render AxOverlayRenderFunc, adjustment AxOverlayAdjustmentFunc, selectCallback AxOverlayStreamSelectFunc, backend AxOverlayBackendType) *AxOverlaySettings {
	settings := &AxOverlaySettings{}

	settings.ptr = (*C.struct_axoverlay_settings)(C.malloc(C.size_t(unsafe.Sizeof(*settings.ptr))))
	if settings.ptr == nil {
		panic("Failed to allocate memory for axoverlay_settings")
	}

	C.axoverlay_init_axoverlay_settings(settings.ptr)

	if render != nil {
		settings.ptr.render_callback = (C.axoverlay_render_function)(C.GoAxOverlayRenderFunction)
		renderCallback = render
	}

	if adjustment != nil {
		settings.ptr.adjustment_callback = (C.axoverlay_adjustment_function)(C.GoAxOverlayAdjustmentFunction)
		adjustmentCallback = adjustment
	}

	if selectCallback != nil {
		settings.ptr.select_callback = (C.axoverlay_stream_select_function)(C.GoAxOverlayStreamSelectFunction)
		streamSelectCallback = selectCallback
	}
	settings.ptr.backend = C.enum_axoverlay_backend_type(backend)
	return settings
}

func (s *AxOverlaySettings) Free() {
	C.free(unsafe.Pointer(s.ptr))
}

//export GoAxOverlayStreamSelectFunction
func GoAxOverlayStreamSelectFunction(camera C.gint, width C.gint, height C.gint, rotation C.gint, isMirrored C.gboolean, streamType C.enum_axoverlay_stream_type) C.gboolean {
	if streamSelectCallback != nil {
		return goBooleanToC(streamSelectCallback(int(camera), int(width), int(height), int(rotation), ctoGoBoolean(isMirrored), AxOverlayStreamType(streamType)))
	}
	return C.FALSE
}

//export GoAxOverlayAdjustmentFunction
func GoAxOverlayAdjustmentFunction(id C.gint, stream *C.struct_axoverlay_stream_data, postype *C.enum_axoverlay_position_type, overlayX *C.gfloat, overlayY *C.gfloat, overlayWidth *C.gint, overlayHeight *C.gint, userData unsafe.Pointer) {
	if adjustmentCallback != nil {
		var goOverlayX float32 = float32(*overlayX)
		var goOverlayY float32 = float32(*overlayY)
		var goOverlayWidth int = int(*overlayWidth)
		var goOverlayHeight int = int(*overlayHeight)
		goPostype := AxOverlayPositionType(*postype)
		handle := cgo.Handle(userData)
		adjustmentCallback(int(id), newStreamDataFromC(stream), &goPostype, &goOverlayX, &goOverlayY, &goOverlayWidth, &goOverlayHeight, handle.Value())
		*overlayX = C.gfloat(goOverlayX)
		*overlayY = C.gfloat(goOverlayY)
		*overlayWidth = C.gint(goOverlayWidth)
		*overlayHeight = C.gint(goOverlayHeight)
		*postype = C.enum_axoverlay_position_type(goPostype)
	}
}

func newStreamDataFromC(stream *C.struct_axoverlay_stream_data) *AxOverlayStreamData {
	// TODO: Type
	return &AxOverlayStreamData{Ptr: stream, ID: int(stream.id), Camera: int(stream.camera), Width: int(stream.width), Height: int(stream.height), Rotation: int(stream.rotation)}
}

//export GoAxOverlayRenderFunction
func GoAxOverlayRenderFunction(renderingContext C.gpointer, id C.gint, stream *C.struct_axoverlay_stream_data, postype C.enum_axoverlay_position_type, overlayX C.gfloat, overlayY C.gfloat, overlayWidth C.gint, overlayHeight C.gint, userData unsafe.Pointer) {
	if renderCallback != nil {
		handle := cgo.Handle(userData)
		renderCallback(NewCairoCtxFromC(renderingContext), int(id), newStreamDataFromC(stream), AxOverlayPositionType(postype), float32(overlayX), float32(overlayY), int(overlayWidth), int(overlayHeight), handle.Value())
	}
}

// AxoverlayInit initializes the axoverlay system with specified settings.
func AxOverlayInit(settings *AxOverlaySettings) error {
	var gerr *C.GError
	C.axoverlay_init(settings.ptr, &gerr)
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

// axoverlayCreateOverlay creates an overlay with the specified data.
func AxOverlayCreateOverlay(data *AxOverlayOverlayData, user_data any) (int, error) {
	var gerr *C.GError
	overlayUserDataHandle = cgo.NewHandle(user_data)
	id := C.axoverlay_create_overlay(data.ptr, (C.gpointer)(unsafe.Pointer(overlayUserDataHandle)), &gerr)
	err := newGError(gerr)
	if err != nil {
		overlayUserDataHandle.Delete()
	}
	return int(id), err
}

func AxOvlerayDeleteHandle() {
	overlayUserDataHandle.Delete()
}

// axoverlayDestroyOverlay destroys the overlay with the given ID.
func AxOverlayDestroyOverlay(id int) error {
	fmt.Println("Destroy")
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

func AxOverlayGetMaxResolution(camera int) (int, int, error) {
	var w int
	var h int
	var err error
	if w, err = AxOverlayGetMaxResolutionWidth(camera); err != nil {
		return 0, 0, err
	}
	if h, err = AxOverlayGetMaxResolutionHeight(camera); err != nil {
		return 0, 0, err
	}
	return w, h, nil
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
		R:        uint8(color.red),
		G:        uint8(color.green),
		B:        uint8(color.blue),
		A:        uint8(color.alpha),
		Pixelate: ctoGoBoolean(color.pixelate),
	}, nil
}

// axoverlaySetPaletteColor sets a palette color by index.
func AxOverlaySetPaletteColor(index int, color AxOverlayPaletteColor) error {
	var gerr *C.GError
	cColor := C.struct_axoverlay_palette_color{
		red:      C.guchar(color.R),
		green:    C.guchar(color.G),
		blue:     C.guchar(color.B),
		alpha:    C.guchar(color.A),
		pixelate: goBooleanToC(color.Pixelate),
	}
	C.axoverlay_set_palette_color(C.gint(index), &cColor, &gerr)
	return newGError(gerr)
}
