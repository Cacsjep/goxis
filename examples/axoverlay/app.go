package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

var (
	err             error
	max_w           int
	max_h           int
	app             *goxis.AcapApplication
	overlayProvider *goxis.OverlayProvider
	overlay_id_rect int
	overlay_id_text int
)

func cleanBackground(ctx *acap.CairoContext, width, height float64) {
	ctx.SetSourceRGBA(goxis.ColorTransparent)
	ctx.SetOperator(acap.OPERATOR_SOURCE)
	ctx.Rectangle(0, 0, width, height)
	ctx.Fill()
}

func streamSelectCallback(camera, width, height, rotation int, isMirrored bool, streamType acap.AxOverlayStreamType) bool {
	fmt.Printf("Stream Select | Camera: %d, Width: %d, Height: %d, Rotation: %d, IsMirrored: %t, StreamType: %v\n",
		camera, width, height, rotation, isMirrored, streamType)
	return false
}

func adjustmentCallback(OverlayId int, stream *acap.AxOverlayStreamData, positionType *acap.AxOverlayPositionType, OverlayX, OverlayY *float32, OverlayWidth, OverlayHeight *int, userData any) {
	fmt.Printf("Adjustment | ID: %d, Stream: %+v, PositionType: %s, OverlayX: %f, OverlayY: %f, OverlayWidth: %d, OverlayHeight: %d\n",
		OverlayId, stream, positionType.String(), *OverlayX, *OverlayY, *OverlayWidth, *OverlayHeight)
}
func renderCallback(cairoCtx *acap.CairoContext, OverlayId int, stream *acap.AxOverlayStreamData, positionType acap.AxOverlayPositionType, OverlayX, OverlayY float32, OverlayWidth, OverlayHeight int, userData any) {
	fmt.Printf("Render | ID: %d, Stream: %+v, PositionType: %s, OverlayX: %f, OverlayY: %f, OverlayWidth: %d, OverlayHeight: %d\n",
		OverlayId, stream, positionType.String(), OverlayX, OverlayY, OverlayWidth, OverlayHeight)

	if OverlayId == overlay_id_text {
		cairoCtx.DrawText("hallo", 10, 10, 32.0, "serif", goxis.ColorMaterialBlack)
	} else {
		cleanBackground(cairoCtx, float64(stream.Width), float64(stream.Height))
		cairoCtx.DrawRect(0, 0, float64(stream.Width), float64(stream.Height/4), goxis.ColorMaterialRed, 9.6)
		cairoCtx.DrawRect(0, float64(stream.Height*3/4), float64(stream.Width), float64(stream.Height), goxis.ColorMaterialRed, 9.6)
	}
}

// This example uses axoverlay example to draw a rectangle
func main() {
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	if overlayProvider, err = goxis.NewOverlayProvider(renderCallback, adjustmentCallback, nil); err != nil {
		panic(err)
	}
	defer overlayProvider.Cleanup()

	if overlay_id_rect, err = overlayProvider.AddOverlay(&goxis.Overlay{
		UseMaxResolution: true,
		OverlayData: &acap.AxOverlayOverlayData{
			AnchorPoint:  acap.AxOverlayAnchorCenter,
			PositionType: acap.AxOverlayCustomNormalized,
			Colorspace:   acap.AxOverlayColorspaceARGB32,
		},
	}); err != nil {
		panic(err)
	}
	fmt.Println("Overlay ID Rectangle", overlay_id_rect)

	if overlay_id_text, err = overlayProvider.AddOverlay(&goxis.Overlay{
		UseMaxResolution: true,
		OverlayData: &acap.AxOverlayOverlayData{
			AnchorPoint:  acap.AxOverlayAnchorCenter,
			PositionType: acap.AxOverlayTopLeft,
			Colorspace:   acap.AxOverlayColorspaceARGB32,
		},
	}); err != nil {
		panic(err)
	}
	fmt.Println("Overlay ID Text", overlay_id_text)

	// Call redraw fo first render
	if err = overlayProvider.Redraw(); err != nil {
		panic(err)
	}

	// Enter main loop
	app.Run()
}
