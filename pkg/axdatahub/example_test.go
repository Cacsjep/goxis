//go:build cgo

package axdatahub_test

import (
	"log"

	"github.com/Cacsjep/goxis/pkg/axdatahub"
)

// Consume AXIS Scene Metadata frames and ended object tracks for video channel
// 1. See the package documentation for the manifest declaration and its AXIS
// OS 13.0.114 caveat.
func ExampleOpen() {
	if err := axdatahub.Load(); err != nil {
		log.Printf("no Device Data Hub on this camera: %v", err)
		return
	}
	sub, err := axdatahub.Open(axdatahub.SubscriptionOptions{
		Name:           "scene consumer",
		Topics:         []string{axdatahub.TopicSceneFrame},
		OptionalTopics: []string{axdatahub.TopicSceneObjectTrack},
		RequireTopics:  true,
		Instances:      []map[string]int64{{axdatahub.InstanceKeyChannel: 1}},
	})
	if err != nil {
		log.Printf("subscribe: %v", err)
		return
	}
	defer sub.Close()
	for {
		select {
		case sample := <-sub.Samples():
			switch sample.Topic {
			case axdatahub.TopicSceneFrame:
				frame, err := axdatahub.DecodeSceneFrame(sample.JSON)
				if err == nil {
					log.Printf("%s: %d detections", frame.Timestamp, len(frame.Detections))
				}
			case axdatahub.TopicSceneObjectTrack:
				track, err := axdatahub.DecodeObjectTrack(sample.JSON)
				if err == nil {
					log.Printf("track %s ended after %.1f s", track.ID, track.Duration)
				}
			}
		case <-sub.Failed():
			log.Printf("connection lost: %v", sub.Err())
			return
		}
	}
}
