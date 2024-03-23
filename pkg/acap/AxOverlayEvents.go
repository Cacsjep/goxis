package acap

/*
#cgo LDFLAGS: -laxoverlay
#cgo pkg-config: gio-2.0 glib-2.0 cairo axoverlay
#include <axoverlay.h>
#include <cairo/cairo.h>
*/
import "C"
import (
	"fmt"
)

// OverlayRenderEvent wraps callback data for axoverlay_render_function.
type OverlayRenderEvent struct {
	CairoCtx                    *CairoContext         // A pointer to the Cairo rendering context.
	OverlayId                   int                   // The ID of the overlay.
	Stream                      *AxOverlayStreamData  // Information about the stream being rendered to.
	PositionType                AxOverlayPositionType // The position type of the overlay.
	OverlayX, OverlayY          float32               // The x and y coordinates of the overlay, can be adjusted.
	OverlayWidth, OverlayHeight int                   // The width and height of the overlay, can be adjusted.
	Userdata                    any                   // Optional user data associated with this overlay.
}

// String method for OverlayRenderEvent
func (e OverlayRenderEvent) String() string {
	return fmt.Sprintf("OverlayRenderEvent{CairoCtx: %v, OverlayId: %d, Stream: %v, PositionType: %v, OverlayX: %f, OverlayY: %f, OverlayWidth: %d, OverlayHeight: %d, Userdata: %v}",
		e.CairoCtx, e.OverlayId, e.Stream, e.PositionType, e.OverlayX, e.OverlayY, e.OverlayWidth, e.OverlayHeight, e.Userdata)
}

// OverlayAdjustmentEvent wraps callback data for axoverlay_adjustment_function.
type OverlayAdjustmentEvent struct {
	OverlayId                   int                    // The ID of the overlay.
	Stream                      *AxOverlayStreamData   // The stream that the overlay is displayed on.
	PositionType                *AxOverlayPositionType // The position type of the overlay.
	OverlayX, OverlayY          *float32               // Pointers to the x and y coordinates of the overlay, can be adjusted.
	OverlayWidth, OverlayHeight *int                   // Pointers to the width and height of the overlay, can be adjusted.
	Userdata                    any                    // Optional user data associated with this overlay.
}

// String method for OverlayAdjustmentEvent
func (e OverlayAdjustmentEvent) String() string {
	return fmt.Sprintf("OverlayAdjustmentEvent{OverlayId: %d, Stream: %v, PositionType: %v, OverlayX: %v, OverlayY: %v, OverlayWidth: %v, OverlayHeight: %v, Userdata: %v}",
		e.OverlayId, e.Stream, e.PositionType, e.OverlayX, e.OverlayY, e.OverlayWidth, e.OverlayHeight, e.Userdata)
}

// OverlayStreamSelectEvent wraps callback data for axoverlay_stream_select_function.
type OverlayStreamSelectEvent struct {
	Camera                  int                 // The ID of the camera for the stream.
	Width, Height, Rotation int                 // The width, height, and rotation of the stream.
	IsMirrored              bool                // TRUE if mirroring is enabled for the stream.
	StreamType              AxOverlayStreamType // The type of the stream.
}

// String method for OverlayStreamSelectEvent
func (e OverlayStreamSelectEvent) String() string {
	return fmt.Sprintf("OverlayStreamSelectEvent{Camera: %d, Width: %d, Height: %d, Rotation: %d, IsMirrored: %t, StreamType: %v}",
		e.Camera, e.Width, e.Height, e.Rotation, e.IsMirrored, e.StreamType)
}
