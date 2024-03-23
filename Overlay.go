package goxis

import (
	"errors"
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acap"
)

var onlyOnce = false

type OverlayProvider struct {
	renderCallback       acap.AxOverlayRenderFunc
	adjustmentCallback   acap.AxOverlayAdjustmentFunc
	streamSelectCallback acap.AxOverlayStreamSelectFunc
	settings             *acap.AxOverlaySettings
	palleteColors        map[int]acap.AxOverlayPaletteColor
	overlays             map[int]*Overlay
}

type Overlay struct {
	overlayId        int
	OverlayData      *acap.AxOverlayOverlayData
	Camera           int
	UseMaxResolution bool
	Userdata         any
}

type StreamSelectEvent struct {
	Camera                  int
	Width, Height, Rotation int
	IsMirrored              bool
	StreamType              acap.AxOverlayStreamType
}

// Creates a cairo backend overlay Provider, this can only created once !!!
func NewOverlayProvider(renderCallback acap.AxOverlayRenderFunc, adjustmentCallback acap.AxOverlayAdjustmentFunc, streamSelectCallback acap.AxOverlayStreamSelectFunc) (*OverlayProvider, error) {
	var err error

	if !onlyOnce {
		onlyOnce = true
	} else {
		return nil, errors.New("Only one overlay provider could created")
	}

	if !acap.AxOverlayIsBackendSupported(acap.AxOverlayCairoImageBackend) {
		return nil, errors.New("Cairo backend not supported")
	}

	op := &OverlayProvider{
		renderCallback:       renderCallback,
		adjustmentCallback:   adjustmentCallback,
		streamSelectCallback: streamSelectCallback,
		overlays:             make(map[int]*Overlay),
	}

	op.settings = acap.NewAxOverlaySettings(op.renderCallback, op.adjustmentCallback, op.streamSelectCallback, acap.AxOverlayCairoImageBackend)
	if err = acap.AxOverlayInit(op.settings); err != nil {
		return nil, err
	}
	/* for color_index, palleteColor := range op.palleteColors {
		if err = acap.AxOverlaySetPaletteColor(color_index, palleteColor); err != nil {
			return nil, err
		}
	} */
	return op, nil
}

func (op *OverlayProvider) AddOverlay(overlay *Overlay) (overlayId int, err error) {
	if overlay.UseMaxResolution {
		if err := overlay.SetMaxResolution(overlay.Camera); err != nil {
			return 0, err
		}
	}
	acap.AxOverlayDataInitalze(overlay.OverlayData)
	if overlay.overlayId, err = acap.AxOverlayCreateOverlay(overlay.OverlayData, overlay.Userdata); err != nil {
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
	return acap.AxOverlayRedraw()
}

func (op *OverlayProvider) Cleanup() {
	op.settings.Free()
	for _, overlay := range op.overlays {
		op.RemoveOverlay(overlay.overlayId)
	}
	acap.AxOvlerayDeleteHandle()
	acap.AxOverlayCleanup()
}

func (ov *Overlay) SetMaxResolution(camera int) (err error) {
	if ov.OverlayData.Width, ov.OverlayData.Height, err = acap.AxOverlayGetMaxResolution(1); err != nil {
		return err
	}
	return nil
}

func (ov *Overlay) Destroy() {
	ov.OverlayData.Free()
	acap.AxOverlayDestroyOverlay(ov.overlayId)
}
