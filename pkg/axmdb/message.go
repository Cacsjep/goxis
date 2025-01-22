package axmdb

import (
	"encoding/json"
	"fmt"
	"time"
)

// MessageType defines the interface for all message types.
type MessageType interface {
	TransformMessage(jsonString string) (MessageType, error) // Converts a JSON string into a specific type
}

// SceneDescription implements MessageType for `analytics_scene_description`.
type SceneDescription struct {
	Frame Frame `json:"frame"`
}

func (sd *SceneDescription) String() string {
	return fmt.Sprintf("Frame: %v, Observations: %d, Operations: %d", sd.Frame.Timestamp, len(sd.Frame.Observations), len(sd.Frame.Operations))
}

// TransformMessage parses a JSON string into a SceneDescription.
func (sd SceneDescription) TransformMessage(jsonString string) (MessageType, error) {
	var parsed SceneDescription
	err := json.Unmarshal([]byte(jsonString), &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// ConsolidatedTrack implements MessageType for `consolidated_track`.
type ConsolidatedTrack struct {
	Classes      []Class       `json:"classes"`
	Duration     float64       `json:"duration"`
	EndTime      time.Time     `json:"end_time"`
	ID           string        `json:"id"`
	Image        Image         `json:"image"`
	Observations []Observation `json:"observations"`
	StartTime    time.Time     `json:"start_time"`
}

// TransformMessage parses a JSON string into a ConsolidatedTrack.
func (ct ConsolidatedTrack) TransformMessage(jsonString string) (MessageType, error) {
	var parsed ConsolidatedTrack
	err := json.Unmarshal([]byte(jsonString), &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

// Frame represents the frame structure for `analytics_scene_description`.
type Frame struct {
	Timestamp    time.Time     `json:"timestamp"`
	Observations []Observation `json:"observations"`
	Operations   []Operation   `json:"operations"`
}

// Observation represents observations in both formats.
type Observation struct {
	BoundingBox Box       `json:"bounding_box"`
	Timestamp   time.Time `json:"timestamp,omitempty"` // For consolidated tracks
	TrackID     string    `json:"track_id,omitempty"`  // For scene description
	Class       *Class    `json:"class,omitempty"`     // Optional class information
	Image       *Image    `json:"image,omitempty"`     // Optional image information
}

// Operation represents operations in `analytics_scene_description`.
type Operation struct {
	Type string `json:"type"`
	ID   string `json:"id,omitempty"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// Class represents class details in observations and tracks.
type Class struct {
	Type                string      `json:"type"`
	Score               float64     `json:"score"`
	UpperClothingColors []ColorInfo `json:"upper_clothing_colors,omitempty"`
	LowerClothingColors []ColorInfo `json:"lower_clothing_colors,omitempty"`
	Colors              []ColorInfo `json:"colors,omitempty"` // For consolidated tracks
}

// ColorInfo represents clothing color data.
type ColorInfo struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

// Image represents image data with bounding box.
type Image struct {
	BoundingBox Box       `json:"bounding_box"`
	Data        string    `json:"data"`
	Timestamp   time.Time `json:"timestamp,omitempty"` // Optional for consolidated tracks
}

// Box represents bounding box coordinates.
type Box struct {
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Top    float64 `json:"top"`
}
