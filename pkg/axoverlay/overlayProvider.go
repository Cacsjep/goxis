package axoverlay

import (
	"errors"
	"fmt"
)

var onlyOnce = false

type OverlayProvider struct {
	renderCallback       AxOverlayRenderCallback
	adjustmentCallback   AxOverlayAdjustmentCallback
	streamSelectCallback AxOverlayStreamSelectCallback
	settings             *AxOverlaySettings
	palleteColors        map[int]AxOverlayPaletteColor
	overlays             map[int]*Overlay
}

type Overlay struct {
	overlayId        int
	OverlayData      *AxOverlayOverlayData
	Camera           int
	UseMaxResolution bool
	Userdata         any
}

type StreamSelectEvent struct {
	Camera                  int
	Width, Height, Rotation int
	IsMirrored              bool
	StreamType              AxOverlayStreamType
}

// Creates a default overlay with anchor center and ARGB32 colorspace
func NewAnchorCenterRrgbaOverlay(positonType AxOverlayPositionType, userData any) *Overlay {
	return &Overlay{
		UseMaxResolution: true,
		OverlayData: &AxOverlayOverlayData{
			AnchorPoint:  AxOverlayAnchorCenter,
			PositionType: positonType,
			Colorspace:   AxOverlayColorspaceARGB32,
		},
		Userdata: userData,
	}
}

// Creates a cairo backend overlay Provider, this can only created once !!!
func NewOverlayProvider(renderCallback AxOverlayRenderCallback, adjustmentCallback AxOverlayAdjustmentCallback, streamSelectCallback AxOverlayStreamSelectCallback) (*OverlayProvider, error) {
	var err error

	if !onlyOnce {
		onlyOnce = true
	} else {
		return nil, errors.New("Only one overlay provider could created")
	}

	if !AxOverlayIsBackendSupported(AxOverlayCairoImageBackend) {
		return nil, errors.New("Cairo backend not supported")
	}

	op := &OverlayProvider{
		renderCallback:       renderCallback,
		adjustmentCallback:   adjustmentCallback,
		streamSelectCallback: streamSelectCallback,
		overlays:             make(map[int]*Overlay),
	}

	op.settings = NewAxOverlaySettings(op.renderCallback, op.adjustmentCallback, op.streamSelectCallback, AxOverlayCairoImageBackend)
	if err = AxOverlayInit(op.settings); err != nil {
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
	AxOverlayDataInitalze(overlay.OverlayData)
	if overlay.overlayId, err = AxOverlayCreateOverlay(overlay.OverlayData, overlay.Userdata); err != nil {
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
	return AxOverlayRedraw()
}

func (op *OverlayProvider) Cleanup() {
	op.settings.Free()
	for _, overlay := range op.overlays {
		op.RemoveOverlay(overlay.overlayId)
	}
	AxOvlerayDeleteHandle()
	AxOverlayCleanup()
}

func (ov *Overlay) SetMaxResolution(camera int) (err error) {
	if ov.OverlayData.Width, ov.OverlayData.Height, err = AxOverlayGetMaxResolution(1); err != nil {
		return err
	}
	return nil
}

func (ov *Overlay) Destroy() {
	ov.OverlayData.Free()
	AxOverlayDestroyOverlay(ov.overlayId)
}
