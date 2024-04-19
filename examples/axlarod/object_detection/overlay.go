package main

import (
	"fmt"
	"strings"

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
	for _, obj := range lea.prediction_result.Detections {
		scaled_box := obj.Box.Scale(renderEvent.Stream.Width, renderEvent.Stream.Height)
		cords := scaled_box.ToCords64()
		renderEvent.CairoCtx.DrawBoundingBox(
			cords.X,
			cords.Y,
			cords.W,
			cords.H,
			axoverlay.ColorMaterialBlue,
			fmt.Sprintf("%s %d%%", strings.ToUpper(coco_labels[int(obj.Class)]), int(obj.Score*100)),
			axoverlay.ColorWite,
			17,
			"sans",
			170,
		)
		lea.app.Syslog.Infof("Render overlay for object: %s, score %d%%", coco_labels[int(obj.Class)], int(obj.Score*100))
	}

}
