package main

import (
	"fmt"
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axoverlay"
)

// This example demonstrate how to use overlay provider to draw overlays on the camera stream.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/axoverlay
// ! Note: Overlay callbacks only invoked when stream is viewed via web ui or rtsp etc..
var (
	app             *acapapp.AcapApplication
	overlayProvider *acapapp.OverlayProvider
	err             error
	overlay_id_rect int
	overlay_id_text int
	counter         int
)

// streamSelectCallback can be used to select which streams to render overlays to.
// Note that YCBCR streams are always skipped since these are used for analytics.
// ! Just for demo demonstration
func streamSelectCallback(streamSelectEvent *axoverlay.OverlayStreamSelectEvent) bool {
	return true
}

// adjustmentCallback is called when an overlay needs adjustments.
// This let developers make adjustments to the size and position of their overlays for each stream.
// This callback function is called prior to rendering every time when an overlay
// is rendered on a stream, which is useful if the resolution has been
// updated or rotation has changed.
func adjustmentCallback(adjustmentEvent *axoverlay.OverlayAdjustmentEvent) {
	app := adjustmentEvent.Userdata.(*acapapp.AcapApplication)
	app.Syslog.Infof("Adjust callback for overlay-%d: %dx%d", adjustmentEvent.OverlayId, adjustmentEvent.OverlayWidth, adjustmentEvent.OverlayHeight)
	app.Syslog.Infof("Adjust callback for stream: %dx%d", adjustmentEvent.Stream.Width, adjustmentEvent.Stream.Height)

	*adjustmentEvent.OverlayWidth = adjustmentEvent.Stream.Width
	*adjustmentEvent.OverlayHeight = adjustmentEvent.Stream.Height
}

// renderCallback is called whenever the system redraws an overlay
// This can happen in two cases, Redraw() is called or a new stream is started.
func renderCallback(renderEvent *axoverlay.OverlayRenderEvent) {
	app := renderEvent.Userdata.(*acapapp.AcapApplication)
	app.Syslog.Infof("Render callback for camera: %d", renderEvent.Stream.Camera)
	app.Syslog.Infof("Render callback for overlay-%d: %dx%d", renderEvent.OverlayId, renderEvent.OverlayWidth, renderEvent.OverlayHeight)
	app.Syslog.Infof("Render callback for stream: %dx%d", renderEvent.Stream.Width, renderEvent.Stream.Height)

	if renderEvent.OverlayId == overlay_id_text {
		renderEvent.CairoCtx.DrawText(fmt.Sprintf("Counter: %d", counter), 10, 10, 32.0, "serif", axoverlay.ColorBlack)
	} else if renderEvent.OverlayId == overlay_id_rect {
		renderEvent.CairoCtx.DrawTransparent(renderEvent.Stream.Width, renderEvent.Stream.Height)
		renderEvent.CairoCtx.DrawRect(0, 0, float64(renderEvent.Stream.Width), float64(renderEvent.Stream.Height/4), axoverlay.ColorMaterialRed, 9.6)
		renderEvent.CairoCtx.DrawRect(0, float64(renderEvent.Stream.Height*3/4), float64(renderEvent.Stream.Width), float64(renderEvent.Stream.Height), axoverlay.ColorMaterialRed, 9.6)
	} else {
		app.Syslog.Warn("Unknown overlay id!")
	}
}

func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app = acapapp.NewAcapApplication()

	// Overlayprovider is an highlevel wrapper around AxOvleray to make life easier
	if overlayProvider, err = acapapp.NewOverlayProvider(renderCallback, adjustmentCallback, streamSelectCallback); err != nil {
		panic(err)
	}
	app.AddCloseCleanFunc(overlayProvider.Cleanup)

	// we pass app as userdata to access syslog from app in callbacks
	if overlay_id_rect, err = overlayProvider.AddOverlay(acapapp.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayCustomNormalized, app)); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// we pass app as userdata to access syslog from app in callbacks
	if overlay_id_text, err = overlayProvider.AddOverlay(acapapp.NewAnchorCenterRrgbaOverlay(axoverlay.AxOverlayTopLeft, app)); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Draw overlays
	if err = overlayProvider.Redraw(); err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Overlay update - increasing counter and call redraw to invoke a new render call
	go func() {
		for true {
			time.Sleep(time.Second * 1)
			counter++
			if err = overlayProvider.Redraw(); err != nil {
				app.Syslog.Crit(err.Error())
			}
		}
	}()

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axoverlay needs a running event loop to handle the overlay callbacks corretly
	app.Run()
}
