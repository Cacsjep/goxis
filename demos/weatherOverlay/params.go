package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/Cacsjep/goxis/pkg/acap"
)

func (w *WeatherApp) LoadParams() error {
	var err error
	var lat_str string
	var long_str string
	var color_str string
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
	w.Color = w.UpdateColor(color_str)
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

func (w *WeatherApp) UpdateColor(value string) color.RGBA {
	switch value {
	case "Transparent":
		return acap.ColorTransparent
	case "Black":
		return acap.ColorBlack
	case "White":
		return acap.ColorWite
	case "Red":
		return acap.ColorMaterialRed
	case "Green":
		return acap.ColorMaterialGreen
	case "Blue":
		return acap.ColorMaterialBlue
	case "Indigo":
		return acap.ColorMaterialIndigo
	case "Pink":
		return acap.ColorMaterialPink
	case "Lime":
		return acap.ColorMaterialLime
	case "DeepPurple":
		return acap.ColorMaterialDeepPurple
	case "Amber":
		return acap.ColorMaterialAmber
	case "Teal":
		return acap.ColorMaterialTeal
	case "Cyan":
		return acap.ColorMaterialCyan
	case "LightGreen":
		return acap.ColorMaterialLightGreen
	case "DeepOrange":
		return acap.ColorMaterialDeepOrange
	case "Brown":
		return acap.ColorMaterialBrown
	case "Grey":
		return acap.ColorMaterialGrey
	default:
		return acap.ColorBlack
	}
}

func (w *WeatherApp) RegisterParamsCallbacks() error {
	if err := w.AcapApp.ParamHandler.RegisterCallback("Lat", func(name, value string, userdata any) {
		if err := w.UpdateCoords(name, value); err != nil {
			w.AcapApp.Syslog.Errorf("Failed to update latitude: %s", err.Error())
		} else {
			w.AcapApp.Syslog.Infof("Latitude updated to: %f", w.Lat)
			w.UpdateWeather()
			w.Redraw()
		}
	}, nil); err != nil {
		return fmt.Errorf("error registering latitude callback: %w", err)
	}

	if err := w.AcapApp.ParamHandler.RegisterCallback("Long", func(name, value string, userdata any) {
		if err := w.UpdateCoords(name, value); err != nil {
			w.AcapApp.Syslog.Errorf("Failed to update longitude: %s", err.Error())
		} else {
			w.AcapApp.Syslog.Infof("Longitude updated to: %f", w.Long)
			w.UpdateWeather()
			w.Redraw()
		}
	}, nil); err != nil {
		return fmt.Errorf("error registering longitude callback: %w", err)
	}

	if err := w.AcapApp.ParamHandler.RegisterCallback("Color", func(name, value string, userdata any) {
		w.Color = w.UpdateColor(value)
		w.AcapApp.Syslog.Infof("Color updated to: %s", value)
		w.UpdateWeather()
		w.Redraw()
	}, nil); err != nil {
		return fmt.Errorf("error registering longitude callback: %w", err)
	}

	return nil
}
