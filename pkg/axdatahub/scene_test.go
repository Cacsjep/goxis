package axdatahub

import (
	"errors"
	"testing"
	"time"
)

// A frame shaped by the published Frame-v1 schema: detections and track_events
// at the top level, not inside a "frame" object as on the Message Broker.
const frameSample = `{
  "channel_id": 1,
  "timestamp": "2026-11-07T10:17:45.524804Z",
  "detections": [
    {
      "bounding_box": {"left": 0.1, "top": 0.2, "right": 0.3, "bottom": 0.6},
      "class": {"type": "Human", "score": 0.91,
                "upper_clothing_colors": [{"name": "red", "score": 0.8}],
                "lower_clothing_colors": [{"name": "blue", "score": 0.7}]},
      "object_track_id": "7ef28e6e-7cae-4186-b630-23819d0bfd98"
    },
    {
      "bounding_box": {"left": 0.5, "top": 0.5, "right": 0.9, "bottom": 0.8},
      "class": {"type": "Car", "score": 0.6, "colors": [{"name": "gray", "score": 0.4}, {"name": "white", "score": 0.5}]},
      "object_track_id": "c1"
    },
    {
      "bounding_box": {"left": 0.4, "top": 0.4, "right": 0.45, "bottom": 0.5},
      "class": null,
      "object_track_id": null
    }
  ],
  "track_events": [
    {"type": "Rename", "from_id": "old", "to_id": "c1"},
    {"type": "TrackEnded", "object_track_id": "gone"}
  ]
}`

func TestDecodeSceneFrame(t *testing.T) {
	f, err := DecodeSceneFrame(frameSample)
	if err != nil {
		t.Fatal(err)
	}
	if f.ChannelID != 1 {
		t.Fatalf("channel %d", f.ChannelID)
	}
	if want := time.Date(2026, 11, 7, 10, 17, 45, 524804000, time.UTC); !f.Timestamp.Equal(want) {
		t.Fatalf("timestamp %v, want %v", f.Timestamp, want)
	}
	if len(f.Detections) != 3 {
		t.Fatalf("%d detections, want 3", len(f.Detections))
	}
	human := f.Detections[0]
	if human.ObjectTrackID != "7ef28e6e-7cae-4186-b630-23819d0bfd98" || human.Class == nil || human.Class.Type != "Human" || human.Class.Score != 0.91 {
		t.Fatalf("human detection %+v", human)
	}
	if human.BoundingBox != (BoundingBox{Left: 0.1, Top: 0.2, Right: 0.3, Bottom: 0.6}) {
		t.Fatalf("human box %+v", human.BoundingBox)
	}
	if len(human.Class.UpperClothingColors) != 1 || human.Class.UpperClothingColors[0].Name != "red" ||
		len(human.Class.LowerClothingColors) != 1 || human.Class.LowerClothingColors[0].Name != "blue" {
		t.Fatalf("clothing colours %+v", human.Class)
	}
	if car := f.Detections[1]; car.Class == nil || len(car.Class.Colors) != 2 {
		t.Fatalf("car colours %+v", car.Class)
	}
	if blank := f.Detections[2]; blank.Class != nil || blank.ObjectTrackID != "" {
		t.Fatalf("null class and track id should decode empty: %+v", blank)
	}
	if len(f.TrackEvents) != 2 ||
		f.TrackEvents[0] != (TrackingEvent{Type: TrackEventRename, FromID: "old", ToID: "c1"}) ||
		f.TrackEvents[1] != (TrackingEvent{Type: TrackEventTrackEnded, ObjectTrackID: "gone"}) {
		t.Fatalf("track events %+v", f.TrackEvents)
	}
}

func TestDecodeSceneFrameEmptyFrame(t *testing.T) {
	f, err := DecodeSceneFrame(`{"channel_id":2,"timestamp":"2026-11-07T10:00:00Z","detections":null,"track_events":null}`)
	if err != nil {
		t.Fatal(err)
	}
	if f.ChannelID != 2 || len(f.Detections) != 0 || len(f.TrackEvents) != 0 {
		t.Fatalf("%+v", f)
	}
}

func TestDecodeSceneFrameRejectsNonFrames(t *testing.T) {
	// A Message Broker frame has its timestamp inside "frame": not a DDH frame.
	if _, err := DecodeSceneFrame(`{"frame":{"timestamp":"2026-11-07T10:00:00Z","observations":[]}}`); !errors.Is(err, ErrNotSceneFrame) {
		t.Fatalf("err %v, want ErrNotSceneFrame", err)
	}
	if _, err := DecodeSceneFrame(`not json`); err == nil {
		t.Fatal("garbage decoded")
	}
}

// The example payload from AxisCommunications/acap-native-sdk-examples
// device-data-hub/consume-scene-metadata (shortened path, image data).
const objectTrackSample = `{
  "channel_id": 1,
  "classes": [
    {"colors": [{"name": "blue", "score": 0.7316}, {"name": "gray", "score": 0.1604}], "score": 0.4848, "type": "Truck"},
    {"colors": [{"name": "blue", "score": 0.7316}], "score": 0.4812, "type": "Car"},
    {"score": 0.0376, "type": "VehicleOther"}
  ],
  "duration": 2.499,
  "end_time": "2026-03-20T10:17:48.024761Z",
  "id": "25436c51-6a75-4f23-99dc-7ddd638e4e61",
  "image": {
    "crop_box": {"bottom": 0.7349, "left": 0.6769, "right": 0.7983, "top": 0.6139},
    "data": "/9j/4AAQSkZJRgABAQAAAQABAAD/2Q=="
  },
  "parts": [{"object_track_id": "7ef28e6e-7cae-4186-b630-23819d0bfd98"}],
  "path": [
    {"bounding_box": {"bottom": 0.786, "left": 0.0199, "right": 0.0937, "top": 0.7161}, "timestamp": "2026-03-20T10:17:45.524804Z"},
    {"bounding_box": null, "timestamp": "2026-03-20T10:17:46.524804Z"},
    {"bounding_box": {"bottom": 0.6446, "left": 0.9833, "right": 0.999, "top": 0.5598}, "timestamp": "2026-03-20T10:17:48.024761Z"}
  ],
  "start_time": "2026-03-20T10:17:45.524804Z"
}`

func TestDecodeObjectTrack(t *testing.T) {
	tr, err := DecodeObjectTrack(objectTrackSample)
	if err != nil {
		t.Fatal(err)
	}
	if tr.ID != "25436c51-6a75-4f23-99dc-7ddd638e4e61" || tr.ChannelID != 1 || tr.Duration != 2.499 {
		t.Fatalf("%+v", tr)
	}
	if len(tr.Classes) != 3 || tr.Classes[0].Type != "Truck" || len(tr.Classes[0].Colors) != 2 {
		t.Fatalf("classes %+v", tr.Classes)
	}
	if tr.Image == nil || tr.Image.Data == "" || tr.Image.CropBox.Top != 0.6139 {
		t.Fatalf("image %+v", tr.Image)
	}
	if len(tr.Path) != 3 || tr.Path[0].BoundingBox == nil || tr.Path[1].BoundingBox != nil || tr.Path[2].BoundingBox.Left != 0.9833 {
		t.Fatalf("path %+v", tr.Path)
	}
	if tr.StartTime.IsZero() || tr.EndTime.Sub(tr.StartTime) <= 0 {
		t.Fatalf("times %v %v", tr.StartTime, tr.EndTime)
	}
	if got := tr.PartIDs(); len(got) != 1 || got[0] != "7ef28e6e-7cae-4186-b630-23819d0bfd98" {
		t.Fatalf("parts %v", got)
	}
}

func TestObjectTrackPartIDsAreFlatAndUnique(t *testing.T) {
	tr := ObjectTrack{ID: "self", Parts: []ObjectTrackReference{
		{ObjectTrackID: "a", Parts: []ObjectTrackReference{{ObjectTrackID: "b"}, {ObjectTrackID: "self"}}},
		{ObjectTrackID: "b"},
		{Parts: []ObjectTrackReference{{ObjectTrackID: "c"}}},
	}}
	got := tr.PartIDs()
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestDecodeObjectTrackRejectsNonTracks(t *testing.T) {
	if _, err := DecodeObjectTrack(`{"channel_id":1}`); !errors.Is(err, ErrNotObjectTrack) {
		t.Fatalf("err %v, want ErrNotObjectTrack", err)
	}
}
