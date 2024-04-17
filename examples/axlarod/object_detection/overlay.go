package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

func (lea *larodExampleApplication) InitOverlay() error {
	if lea.overlayProvider, err = acapapp.NewOverlayProvider(renderCallback, nil, streamSelectCallback); err != nil {
		return err
	}
	lea.app.AddCloseCleanFunc(lea.overlayProvider.Cleanup)
	if _, err = lea.overlayProvider.AddOverlay(acapapp.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, lea)); err != nil {
		return err
	}
	return nil
}

// we want all streams
func streamSelectCallback(streamSelectEvent *axoverlay.OverlayStreamSelectEvent) bool {
	return true
}

func renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	lea := renderEvent.Userdata.(*larodExampleApplication)
	renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)
	for _, obj := range lea.prediction_result.Detections {
		scaled_box := obj.Box.Scale(renderEvent.Stream.Width, renderEvent.Stream.Height)
		cords := scaled_box.ToCords64()
		renderEvent.CairoCtx.DrawText(coco_labels[int(obj.Class)], cords.X+5, cords.Y+5, 24, "sans", axoverlay.ColorMaterialBlue)
		renderEvent.CairoCtx.DrawRect(
			cords.X,
			cords.Y,
			cords.W,
			cords.H,
			axoverlay.ColorMaterialBlue,
			3,
		)
		lea.app.Syslog.Infof("Render overlay for object: %s, score %d%%, cords: %v", coco_labels[int(obj.Class)], int(obj.Score*100), cords)
	}

}
