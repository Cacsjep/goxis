package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acap"
)

func (w *WeatherApp) LoadParams() error {
	var err error
	var color_str string
	var pos_str string
	var size_str string

	if w.Lat, err = w.AcapApp.ParamHandler.GetAsFloat("Lat"); err != nil {
		return err
	}

	if w.Long, err = w.AcapApp.ParamHandler.GetAsFloat("Long"); err != nil {
		return err
	}

	if color_str, err = w.AcapApp.ParamHandler.Get("Color"); err != nil {
		return err
	}
	if pos_str, err = w.AcapApp.ParamHandler.Get("Position"); err != nil {
		return err
	}
	if size_str, err = w.AcapApp.ParamHandler.Get("Size"); err != nil {
		return err
	}
	w.UpdateColor(color_str)
	w.UpdatePosition(pos_str)
	w.UpdateSize(size_str)
	return nil
}

func (w *WeatherApp) UpdateCoords(parameterName string, value string) (err error) {
	if parameterName == "root.Weatheroverlay.Lat" {
		if w.Lat, err = w.AcapApp.ParamHandler.GetAsFloat("Lat"); err != nil {
			return err
		}
	} else {
		if w.Long, err = w.AcapApp.ParamHandler.GetAsFloat("Long"); err != nil {
			return err
		}
	}
	return nil
}

func (w *WeatherApp) UpdateColor(value string) {
	switch value {
	case "Transparent":
		w.Color = acap.ColorTransparent
	case "Black":
		w.Color = acap.ColorBlack
	case "White":
		w.Color = acap.ColorWite
	case "Red":
		w.Color = acap.ColorMaterialRed
	case "Green":
		w.Color = acap.ColorMaterialGreen
	case "Blue":
		w.Color = acap.ColorMaterialBlue
	case "Indigo":
		w.Color = acap.ColorMaterialIndigo
	case "Pink":
		w.Color = acap.ColorMaterialPink
	case "Lime":
		w.Color = acap.ColorMaterialLime
	case "DeepPurple":
		w.Color = acap.ColorMaterialDeepPurple
	case "Amber":
		w.Color = acap.ColorMaterialAmber
	case "Teal":
		w.Color = acap.ColorMaterialTeal
	case "Cyan":
		w.Color = acap.ColorMaterialCyan
	case "LightGreen":
		w.Color = acap.ColorMaterialLightGreen
	case "DeepOrange":
		w.Color = acap.ColorMaterialDeepOrange
	case "Brown":
		w.Color = acap.ColorMaterialBrown
	case "Grey":
		w.Color = acap.ColorMaterialGrey
	default:
		w.Color = acap.ColorBlack
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

func (w *WeatherApp) RegisterParamsCallbacks() error {
	callbacks := map[string]func(name, value string) error{
		"Lat": func(name, value string) error {
			return w.UpdateCoords(name, value)
		},
		"Long": func(name, value string) error {
			return w.UpdateCoords(name, value)
		},
		"Color": func(_, value string) error {
			w.UpdateColor(value)
			return nil
		},
		"Position": func(name, value string) error {
			w.UpdatePosition(value)
			return nil
		},
		"Size": func(name, value string) error {
			w.UpdateSize(value)
			return nil
		},
	}

	for paramName, updateFunc := range callbacks {
		callbackFunc := func(name, value string, userdata any) {
			if err := updateFunc(name, value); err != nil {
				w.AcapApp.Syslog.Errorf("Failed to update %s: %s", name, err.Error())
			} else {
				w.AcapApp.Syslog.Infof("%s updated to: %s", name, value)
				w.UpdateWeather()
				w.Redraw()
			}
		}

		if err := w.AcapApp.ParamHandler.RegisterCallback(paramName, callbackFunc, nil); err != nil {
			return fmt.Errorf("error registering %s callback: %w", paramName, err)
		}
	}

	return nil
}
