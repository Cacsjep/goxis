package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/Cacsjep/goxis/pkg/acap"
)

func (w *WeatherApp) LoadParams() error {
	var err error
	var color_str string
	var circle_color_str string
	var pos_str string
	var size_str string
	if w.Lat, err = w.AcapApp.ParamHandler.GetAsFloat("Lat"); err != nil {
		return err
	}
	if w.Long, err = w.AcapApp.ParamHandler.GetAsFloat("Long"); err != nil {
		return err
	}
	if w.NordDirection, err = w.AcapApp.ParamHandler.GetAsFloat("NordDirection"); err != nil {
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
	if circle_color_str, err = w.AcapApp.ParamHandler.Get("CircleColor"); err != nil {
		return err
	}
	w.Color = w.GetColor(color_str)
	w.CircleColor = w.GetColor(circle_color_str)
	w.UpdatePosition(pos_str)
	w.UpdateSize(size_str)
	return nil
}

func (w *WeatherApp) GetColor(value string) color.RGBA {
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
		w.Size = 130
	case "medium":
		w.Size = 150
	case "large":
		w.Size = 170
	case "xlarge":
		w.Size = 180
	}
}

func (w *WeatherApp) RegisterParamsCallbacks() error {
	callbacks := map[string]func(name, value string) error{
		"Lat": func(name, value string) error {
			var err error
			if w.Lat, err = strconv.ParseFloat(value, 64); err != nil {
				return err
			}
			return nil
		},
		"Long": func(name, value string) error {
			var err error
			if w.Long, err = strconv.ParseFloat(value, 64); err != nil {
				return err
			}
			return nil
		},
		"Color": func(_, value string) error {
			w.Color = w.GetColor(value)
			return nil
		},
		"CircleColor": func(_, value string) error {
			w.CircleColor = w.GetColor(value)
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
		"NordDirection": func(name, value string) error {
			var err error
			if w.NordDirection, err = strconv.ParseFloat(value, 64); err != nil {
				return err
			}
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
