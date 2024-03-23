package main

import (
	"fmt"
	"strconv"

	"github.com/Cacsjep/goxis/pkg/acap"
)

func (w *WeatherApp) LoadParams() error {
	var err error
	var lat_str string
	var long_str string
	var color_str string
	var pos_str string
	var size_str string
	if lat_str, err = w.AcapApp.ParamHandler.Get("Lat"); err != nil {
		return err
	}
	if err = w.UpdateCoords("root.Weatheroverlay.Lat", lat_str); err != nil {
		return err
	}
	if long_str, err = w.AcapApp.ParamHandler.Get("Long"); err != nil {
		return err
	}
	if err = w.UpdateCoords("root.Weatheroverlay.Long", long_str); err != nil {
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
		w.Lat, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return
		}
	} else {
		w.Long, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return
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
