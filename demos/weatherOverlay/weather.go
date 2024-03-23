package main

import (
	"encoding/json"

	"github.com/innotechdevops/openmeteo"
)

// Define the struct for the "current_weather_units" object in your JSON
type CurrentWeatherUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature   string `json:"temperature"`
	Windspeed     string `json:"windspeed"`
	Winddirection string `json:"winddirection"`
	IsDay         string `json:"is_day"`
	Weathercode   string `json:"weathercode"`
}

// Define the struct for the "current_weather" object in your JSON
type CurrentWeather struct {
	Time          string  `json:"time"`
	Interval      int     `json:"interval"`
	Temperature   float64 `json:"temperature"`
	Windspeed     float64 `json:"windspeed"`
	Winddirection int     `json:"winddirection"`
	IsDay         int     `json:"is_day"`
	Weathercode   int     `json:"weathercode"`
}

// Define the top-level struct to match your JSON structure
type WeatherData struct {
	Latitude             float64             `json:"latitude"`
	Longitude            float64             `json:"longitude"`
	GenerationTimeMs     float64             `json:"generationtime_ms"`
	UTCOffsetSeconds     int                 `json:"utc_offset_seconds"`
	Timezone             string              `json:"timezone"`
	TimezoneAbbreviation string              `json:"timezone_abbreviation"`
	Elevation            float64             `json:"elevation"`
	CurrentWeatherUnits  CurrentWeatherUnits `json:"current_weather_units"`
	CurrentWeather       CurrentWeather      `json:"current_weather"`
}

func (w *WeatherApp) GetWeather() (*WeatherData, error) {
	param := openmeteo.Parameter{
		Latitude:       openmeteo.Float32(float32(w.Lat)),
		Longitude:      openmeteo.Float32(float32(w.Long)),
		CurrentWeather: openmeteo.Bool(true),
	}
	m := openmeteo.New()
	jsonData, err := m.Execute(param)
	if err != nil {
		return nil, err
	}
	var data WeatherData
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (w *WeatherApp) UpdateWeather() {
	var err error
	if w.LastData, err = w.GetWeather(); err != nil {
		w.AcapApp.Syslog.Errorf("Failed to retrieve weather data: %s", err.Error())
		return
	}
	w.AcapApp.Syslog.Infof("Temperature: %.1f%s", w.LastData.CurrentWeather.Temperature, w.LastData.CurrentWeatherUnits.Temperature)
}
