package main

import (
	"math/rand"
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

func main() {
	// Initialize a new ACAP application instance.
	app := acapapp.NewAcapApplication()

	myEvent := &acapapp.CameraPlatformEvent{
		Name:     "myevent",                // Unique identifier for the event.
		NiceName: utils.StrPtr("My Event"), // Human-readable name for the event.
		Entries: []*acapapp.EventEntry{
			{Key: "foo", ValueType: axevent.AXValueTypeInt, IsData: utils.BoolPtr(true), KeyNiceName: utils.StrPtr("Foo Value")},
			{Key: "bar", ValueType: axevent.AXValueTypeDouble, IsData: utils.BoolPtr(true), KeyNiceName: utils.StrPtr("Bar Value")},
			{Key: "baz", ValueType: axevent.AXValueTypeString, IsData: utils.BoolPtr(true), KeyNiceName: utils.StrPtr("Baz Value")},
			{Key: "qux", ValueType: axevent.AXValueTypeBool, IsData: utils.BoolPtr(true), KeyNiceName: utils.StrPtr("Qux Value")},
		},
		Stateless: true,
	}

	mqtt_event_id, err := app.AddCameraPlatformEvent(myEvent)
	if err != nil {
		app.Syslog.Critf("Error adding event declaration: %s", err.Error())
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)

			// Attempt to send a newly created event with dynamic values.
			err := app.SendPlatformEvent(mqtt_event_id, func() (*axevent.AXEvent, error) {
				return myEvent.NewEvent(acapapp.KeyValueMap{
					"foo": rand.Int(),        // Random integer value.
					"bar": rand.Float64(),    // Random floating-point value.
					"baz": "oh yeah",         // Static string value.
					"qux": rand.Intn(2) == 1, // Random boolean value (true or false).
				})
			})

			if err != nil {
				app.Syslog.Errorf("Error sending event: %s", err.Error())
			} else {
				app.Syslog.Infof("Event send")
			}
		}
	}()

	// Run gmain loop with signal handler attached.
	app.Run()
}
