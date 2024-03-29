package goxis

import (
	"errors"
	"fmt"

	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

var onlyOnce = false

type OverlayProvider struct {
	renderCallback       axoverlay.AxOverlayRenderCallback
	adjustmentCallback   axoverlay.AxOverlayAdjustmentCallback
	streamSelectCallback axoverlay.AxOverlayStreamSelectCallback
	settings             *axoverlay.AxOverlaySettings
	palleteColors        map[int]axoverlay.AxOverlayPaletteColor
	overlays             map[int]*Overlay
}

type Overlay struct {
	overlayId        int
	OverlayData      *axoverlay.AxOverlayOverlayData
	Camera           int
	UseMaxResolution bool
	Userdata         any
}

type StreamSelectEvent struct {
	Camera                  int
	Width, Height, Rotation int
	IsMirrored              bool
	StreamType              axoverlay.AxOverlayStreamType
}

// Creates a default overlay with anchor center and ARGB32 colorspace
func NewAnchorCenterRrgbaOverlay(positonType axoverlay.AxOverlayPositionType, userData any) *Overlay {
	return &Overlay{
		UseMaxResolution: true,
		OverlayData: &axoverlay.AxOverlayOverlayData{
			AnchorPoint:  axoverlay.AxOverlayAnchorCenter,
			PositionType: positonType,
			Colorspace:   axoverlay.AxOverlayColorspaceARGB32,
		},
		Userdata: userData,
	}
}

// Creates a cairo backend overlay Provider, this can only created once !!!
func NewOverlayProvider(renderCallback axoverlay.AxOverlayRenderCallback, adjustmentCallback axoverlay.AxOverlayAdjustmentCallback, streamSelectCallback axoverlay.AxOverlayStreamSelectCallback) (*OverlayProvider, error) {
	var err error

	if !onlyOnce {
		onlyOnce = true
	} else {
		return nil, errors.New("Only one overlay provider could created")
	}

	if !axoverlay.AxOverlayIsBackendSupported(axoverlay.AxOverlayCairoImageBackend) {
		return nil, errors.New("Cairo backend not supported")
	}

	op := &OverlayProvider{
		renderCallback:       renderCallback,
		adjustmentCallback:   adjustmentCallback,
		streamSelectCallback: streamSelectCallback,
		overlays:             make(map[int]*Overlay),
	}

	op.settings = axoverlay.NewAxOverlaySettings(op.renderCallback, op.adjustmentCallback, op.streamSelectCallback, axoverlay.AxOverlayCairoImageBackend)
	if err = axoverlay.AxOverlayInit(op.settings); err != nil {
		return nil, err
	}
	return op, nil
}

func (op *OverlayProvider) AddOverlay(overlay *Overlay) (overlayId int, err error) {
	if overlay.UseMaxResolution {
		if err := overlay.SetMaxResolution(overlay.Camera); err != nil {
			return 0, err
		}
	}
	axoverlay.AxOverlayDataInitalze(overlay.OverlayData)
	if overlay.overlayId, err = axoverlay.AxOverlayCreateOverlay(overlay.OverlayData, overlay.Userdata); err != nil {
		return 0, err
	}
	op.overlays[overlay.overlayId] = overlay
	return overlay.overlayId, nil
}

func (op *OverlayProvider) RemoveOverlay(overlayId int) error {
	if overlay, found := op.overlays[overlayId]; found {
		overlay.Destroy()
		delete(op.overlays, overlayId)
		return nil
	}
	return fmt.Errorf("Overlay with ID: %d not found", overlayId)
}

func (op *OverlayProvider) Redraw() error {
	return axoverlay.AxOverlayRedraw()
}

func (op *OverlayProvider) Cleanup() {
	op.settings.Free()
	for _, overlay := range op.overlays {
		op.RemoveOverlay(overlay.overlayId)
	}
	axoverlay.AxOvlerayDeleteHandle()
	axoverlay.AxOverlayCleanup()
}

func (ov *Overlay) SetMaxResolution(camera int) (err error) {
	if ov.OverlayData.Width, ov.OverlayData.Height, err = axoverlay.AxOverlayGetMaxResolution(1); err != nil {
		return err
	}
	return nil
}

func (ov *Overlay) Destroy() {
	ov.OverlayData.Free()
	axoverlay.AxOverlayDestroyOverlay(ov.overlayId)
}
