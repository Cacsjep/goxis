package main

import (
	"math/rand"
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// This example demonstrates how to send a custom event with dynamic values.
// The event is declared with a set of keys and their types, and then sent with
// a set of values that correspond to the keys. The event is sent every second.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/axevent/send_event
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	// Platform event is a struct that contains the event declaration for onvif events for an ACAP
	// application. These events are then listed under the "Application" tab in the camera's web interface.
	// Like Appname: <Event Nice Name> simiarly to Object Analytics or VMD ACAP.
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

	// Add the platform event to app and get the event id.
	// Its declare under the hood the event on the event handler.
	myevent_id, err := app.AddCameraPlatformEvent(myEvent)
	if err != nil {
		app.Syslog.Critf("Error adding event declaration: %s", err.Error())
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)

			// Attempt to send a newly created event with dynamic values.
			// with AcapApplication we try to abstract the low level apis and provide a simple interface to send events,
			// live SendPlatformEvent what accepts the event id and a function that returns a new event that is sent.
			err := app.SendPlatformEvent(myevent_id, func() (*axevent.AXEvent, error) {
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
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axevent needs a running event loop to handle the events callbacks corretly
	app.Run()
}
