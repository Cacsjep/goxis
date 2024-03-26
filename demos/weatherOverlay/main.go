// This demo application draws a weateroverlay
package main

import (
	"image/color"
	"time"

	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/acap"
)

type WeatherApp struct {
	AcapApp              *goxis.AcapApplication
	OvProvider           *goxis.OverlayProvider
	Lat                  float64
	Long                 float64
	Color                color.RGBA
	CircleColor          color.RGBA
	Position             acap.AxOverlayPositionType
	Size                 float64
	LastData             *WeatherData
	TemperatureOverlayId int
	NordDirection        float64
}

func main() {
	var err error
	w := WeatherApp{
		Position: acap.AxOverlayBottomLeft,
		AcapApp:  goxis.NewAcapApplication(),
	}

	if err = w.LoadParams(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to load parameters: %s", err.Error())
	} else {
		w.AcapApp.Syslog.Infof("Coordinates set to Latitude: %f, Longitude: %f", w.Lat, w.Long)
	}

	// Overlayprovider is an highlevel wrapper around AxOvleray to make life easier
	if w.OvProvider, err = goxis.NewOverlayProvider(renderCallback, adjustmentCallback, nil); err != nil {
		w.AcapApp.Syslog.Critf("Failed to initialize OverlayProvider: %s", err.Error())
	}
	w.AcapApp.AddCloseCleanFunc(w.OvProvider.Cleanup)

	if err = w.RegisterParamsCallbacks(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to set up parameter callbacks: %s", err.Error())
	}

	if w.TemperatureOverlayId, err = w.OvProvider.AddOverlay(goxis.NewAnchorCenterRrgbaOverlay(w.Position, &w)); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to add temperature overlay: %s", err.Error())
	}

	go func() {
		for true {
			w.AcapApp.Syslog.Info("Update overlay")
			w.UpdateWeather()
			w.Redraw()
			time.Sleep(time.Second * 60)
		}
	}()

	// Enter main loop, stops automatically via signals
	w.AcapApp.Run()
}
