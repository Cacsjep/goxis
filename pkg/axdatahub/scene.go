package axdatahub

// AXIS Scene Metadata on the Device Data Hub. Pure Go (no cgo), so it builds
// and tests anywhere.
//
// The types follow the published data schemas (Frame-v1, ObjectTrack-v1):
// https://developer.axis.com/vapix/device-data-hub/topics/
// They differ from the Message Broker topics (axmdb): a frame has
// "detections" and "track_events" at the top level instead of a "frame"
// object with "observations" and "operations", a detection names its track in
// "object_track_id", and an object track has a "path" and "parts" instead of
// "observations", and its image has a "crop_box".

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// The released AXIS Scene Metadata topics. Each has one instance per video
// channel, keyed by the integer instance key "channel_id".
const (
	TopicSceneFrame          = "com.axis.scene.frame.v1"
	TopicSceneObjectTrack    = "com.axis.scene.object_track.v1"
	TopicSceneObjectSnapshot = "com.axis.scene.object_snapshot.v1"

	// InstanceKeyChannel is the instance key of the scene topics.
	InstanceKeyChannel = "channel_id"
)

// Tracking event types in a frame's track_events (schema definitions Rename
// and TrackEnded).
const (
	TrackEventRename     = "Rename"
	TrackEventTrackEnded = "TrackEnded"
)

// BoundingBox is in normalized image coordinates (0..1, origin top left).
type BoundingBox struct {
	Left   float64 `json:"left"`
	Top    float64 `json:"top"`
	Right  float64 `json:"right"`
	Bottom float64 `json:"bottom"`
}

// IsZero reports whether all four edges are zero (no box).
func (b BoundingBox) IsZero() bool {
	return b.Left == 0 && b.Top == 0 && b.Right == 0 && b.Bottom == 0
}

// ColorClass is one colour candidate and its score.
type ColorClass struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

// ObjectClass is one classification candidate. The class-specific fields are
// set only for the classes that carry them: clothing colours and carries_bag
// for Human, colours for Car, Bus, Truck and Bike, face_visible for Head.
type ObjectClass struct {
	Type                string       `json:"type"`
	Score               float64      `json:"score"`
	UpperClothingColors []ColorClass `json:"upper_clothing_colors,omitempty"`
	LowerClothingColors []ColorClass `json:"lower_clothing_colors,omitempty"`
	Colors              []ColorClass `json:"colors,omitempty"`
	CarriesBag          *bool        `json:"carries_bag,omitempty"`
	FaceVisible         *float64     `json:"face_visible,omitempty"`
}

// ObjectDetection is one object in a frame.
type ObjectDetection struct {
	BoundingBox   BoundingBox  `json:"bounding_box"`
	Class         *ObjectClass `json:"class,omitempty"`
	ObjectTrackID string       `json:"object_track_id,omitempty"`
	ImageID       string       `json:"image_id,omitempty"`
}

// TrackingEvent is a track_events entry. Rename sets FromID and ToID (FromID
// ends and merges into ToID); TrackEnded sets ObjectTrackID.
type TrackingEvent struct {
	Type          string `json:"type"`
	FromID        string `json:"from_id,omitempty"`
	ToID          string `json:"to_id,omitempty"`
	ObjectTrackID string `json:"object_track_id,omitempty"`
}

// SceneFrame is a com.axis.scene.frame.v1 sample: every detection in one video
// frame of one channel.
type SceneFrame struct {
	ChannelID   int               `json:"channel_id"`
	Timestamp   time.Time         `json:"timestamp"`
	Detections  []ObjectDetection `json:"detections"`
	TrackEvents []TrackingEvent   `json:"track_events"`
}

// Snapshot is a JPEG crop of an object (base64 in Data) and where it was cut
// from the frame.
type Snapshot struct {
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      string      `json:"data"`
	CropBox   BoundingBox `json:"crop_box"`
}

// TimeAndPosition is one point of an object track's path.
type TimeAndPosition struct {
	Timestamp   time.Time    `json:"timestamp"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// ObjectTrackReference names a track that was combined into another one.
type ObjectTrackReference struct {
	ObjectTrackID string                 `json:"object_track_id,omitempty"`
	Parts         []ObjectTrackReference `json:"parts,omitempty"`
}

// ObjectTrack is a com.axis.scene.object_track.v1 sample: one object observed
// over time, published when its track ends.
type ObjectTrack struct {
	ChannelID    int                    `json:"channel_id"`
	ID           string                 `json:"id"`
	Classes      []ObjectClass          `json:"classes"`
	Duration     float64                `json:"duration"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Image        *Snapshot              `json:"image,omitempty"`
	Path         []TimeAndPosition      `json:"path"`
	Parts        []ObjectTrackReference `json:"parts"`
	MinimumSpeed *float64               `json:"minimum_speed,omitempty"`
	MaximumSpeed *float64               `json:"maximum_speed,omitempty"`
}

// PartIDs returns the IDs of every track combined into t, depth first, without
// duplicates and without t's own ID.
func (t ObjectTrack) PartIDs() []string {
	var out []string
	seen := map[string]struct{}{t.ID: {}}
	var walk func([]ObjectTrackReference)
	walk = func(refs []ObjectTrackReference) {
		for _, ref := range refs {
			if ref.ObjectTrackID != "" {
				if _, dup := seen[ref.ObjectTrackID]; !dup {
					seen[ref.ObjectTrackID] = struct{}{}
					out = append(out, ref.ObjectTrackID)
				}
			}
			walk(ref.Parts)
		}
	}
	walk(t.Parts)
	return out
}

// ErrNotSceneFrame and ErrNotObjectTrack mark payloads that parse as JSON but
// lack the fields the schema requires.
var (
	ErrNotSceneFrame  = errors.New("axdatahub: not a scene frame (no timestamp)")
	ErrNotObjectTrack = errors.New("axdatahub: not an object track (no id)")
)

// DecodeSceneFrame parses a com.axis.scene.frame.v1 payload.
func DecodeSceneFrame(payload string) (SceneFrame, error) {
	var f SceneFrame
	if err := json.Unmarshal([]byte(payload), &f); err != nil {
		return SceneFrame{}, fmt.Errorf("axdatahub: scene frame: %w", err)
	}
	if f.Timestamp.IsZero() {
		return SceneFrame{}, ErrNotSceneFrame
	}
	return f, nil
}

// DecodeObjectTrack parses a com.axis.scene.object_track.v1 payload.
func DecodeObjectTrack(payload string) (ObjectTrack, error) {
	var t ObjectTrack
	if err := json.Unmarshal([]byte(payload), &t); err != nil {
		return ObjectTrack{}, fmt.Errorf("axdatahub: object track: %w", err)
	}
	if t.ID == "" {
		return ObjectTrack{}, ErrNotObjectTrack
	}
	return t, nil
}
