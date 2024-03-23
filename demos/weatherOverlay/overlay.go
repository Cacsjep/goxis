package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acap"
)

func adjustmentCallback(adjustmentEvent *acap.OverlayAdjustmentEvent) {
	*adjustmentEvent.OverlayWidth = adjustmentEvent.Stream.Width
	*adjustmentEvent.OverlayHeight = adjustmentEvent.Stream.Height
}

func renderCallback(renderEvent *acap.OverlayRenderEvent) {
	wapp := renderEvent.Userdata.(*WeatherApp)

	if renderEvent.OverlayId == wapp.TemperatureOverlayId {
		temperatureText := fmt.Sprintf("%.1f%s", wapp.LastData.CurrentWeather.Temperature, wapp.LastData.CurrentWeatherUnits.Temperature)
		renderEvent.CairoCtx.DrawText(temperatureText, 10, 10, 32.0, "serif", wapp.Color)
	}
}

func (w *WeatherApp) Redraw() {
	if err := w.OvProvider.Redraw(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to redraw overlays: %s", err.Error())
	}
}
