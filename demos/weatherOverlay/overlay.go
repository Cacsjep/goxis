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

		overlayText := fmt.Sprintf(
			"%s %.1f%s  %.1f%s",
			WindDirectionToGermanHimmelsrichtung(float64(wapp.LastData.CurrentWeather.Winddirection)),
			wapp.LastData.CurrentWeather.Windspeed,
			wapp.LastData.CurrentWeatherUnits.Windspeed,
			wapp.LastData.CurrentWeather.Temperature,
			wapp.LastData.CurrentWeatherUnits.Temperature,
		)
		wapp.AcapApp.Syslog.Infof("Apply changes to overlay... Size: %f, Pos: %s, Color: %v, Text: %s", wapp.Size, wapp.Position.String(), wapp.Color, overlayText)

		renderEvent.CairoCtx.SelectFontFace("serif", 0, 0)
		renderEvent.CairoCtx.SetFontSize(wapp.Size)
		textExtents := renderEvent.CairoCtx.TextExtents(overlayText)
		var x, y float64
		var x2 float64
		padding := wapp.Size * 0.33
		arrowSize := wapp.Size
		switch wapp.Position {
		case acap.AxOverlayTopLeft:
			x = padding + arrowSize
			y = padding + textExtents.Height
			x2 = x + padding
		case acap.AxOverlayTopRight:
			x = float64(renderEvent.Stream.Width) - textExtents.Width - padding
			y = padding + textExtents.Height
			x2 = x
		case acap.AxOverlayBottomLeft:
			x = padding + arrowSize
			y = float64(renderEvent.Stream.Height) - padding
			x2 = x + padding
		case acap.AxOverlayBottomRight:
			x = float64(renderEvent.Stream.Width) - textExtents.Width - padding
			y = float64(renderEvent.Stream.Height) - padding
			x2 = x
		}

		renderEvent.CairoCtx.DrawArrow(x, y, wapp.Size, float64(wapp.LastData.CurrentWeather.Winddirection), wapp.Color)
		renderEvent.CairoCtx.SetSourceRGBA(wapp.Color)
		renderEvent.CairoCtx.MoveTo(x2, y)
		renderEvent.CairoCtx.ShowText(overlayText)
	}
}

func WindDirectionToGermanHimmelsrichtung(degrees float64) string {
	directions := []string{"N", "NNO", "NO", "ONO", "O", "OSO", "SO", "SSO", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW", "N"} // Add "N" again to complete the loop
	index := int((degrees + 11.25) / 22.5)                                                                                          // Divide by 22.5 degrees for each direction, add 11.25 for offset
	return directions[index%16]                                                                                                     // Use modulo to wrap around the array if necessary
}

func (w *WeatherApp) Redraw() {
	if err := w.OvProvider.Redraw(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to redraw overlays: %s", err.Error())
	}
}

func (w *WeatherApp) UpdatePosition(value string) {
	switch value {
	case "tr":
		w.Position = acap.AxOverlayTopRight
	case "br":
		w.Position = acap.AxOverlayBottomRight
	case "bl":
		w.Position = acap.AxOverlayBottomLeft
	case "tl":
		w.Position = acap.AxOverlayTopLeft
	}
}

func (w *WeatherApp) UpdateSize(value string) {
	switch value {
	case "small":
		w.Size = 32.0
	case "medium":
		w.Size = 32.0 + 21.33
	case "large":
		w.Size = 32.0 + 2*21.33
	case "xlarge":
		w.Size = 96.0
	}
}
