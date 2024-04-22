package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

// Initialize the overlay provider
func (lea *larodExampleApplication) InitOverlay() error {
	if lea.overlayProvider, err = acapapp.NewOverlayProvider(renderCallback, nil, nil); err != nil {
		return err
	}
	lea.app.AddCloseCleanFunc(lea.overlayProvider.Cleanup)
	if _, err = lea.overlayProvider.AddOverlay(acapapp.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, lea)); err != nil {
		return err
	}
	return nil
}

// renderCallback is used to draw bounding boxes from the detections via axoverlay
func renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	lea := renderEvent.Userdata.(*larodExampleApplication)
	renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)
	for _, obj := range lea.detections {
		if obj.Confidence > lea.threshold {
			scaled_box := obj.Box.Scale(renderEvent.Stream.Width, renderEvent.Stream.Height)
			cords := scaled_box.ToCords64()
			renderEvent.CairoCtx.DrawBoundingBox(
				cords.X,
				cords.Y,
				cords.W,
				cords.H,
				axoverlay.ColorMaterialBlue,
				fmt.Sprintf("%s %d%%", yolo_labels[obj.BestClassIdx], int(obj.Confidence*100)),
				axoverlay.ColorWite,
				17,
				"sans",
				170,
			)
		}
	}
}
