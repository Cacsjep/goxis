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

	// Define the structure and metadata for an MQTT event. This includes specifying the event name,
	// a human-readable nice name, and the key-value structure that the event will contain.
	mqttEvent := &acapapp.CameraPlatformEvent{
		Name:     "MqttEvent",                               // Unique identifier for the event.
		NiceName: utils.NewStringPointer("Mqtt Event Send"), // Human-readable name for the event.
		KvsEntries: []*axevent.KeyValueEntrie{ // Define the keys and their value types that this event will include.
			{Key: "foo", ValueType: axevent.AXValueTypeInt},
			{Key: "bar", ValueType: axevent.AXValueTypeDouble},
			{Key: "baz", ValueType: axevent.AXValueTypeString},
			{Key: "qux", ValueType: axevent.AXValueTypeBool},
		},
		DataMarkers: []*axevent.AxEventKeyValueSetDataMark{ // Mark keys that represent event data.
			{Key: "foo"},
			{Key: "bar"},
			{Key: "baz"},
			{Key: "qux"},
		},
		NiceNames: []*axevent.AxEventKeyValueSetNiceNames{ // Assign 'nice names' for keys, improving readability in UI.
			{Key: "foo", KeyNiceName: utils.NewStringPointer("Foo Value")},
			{Key: "bar", KeyNiceName: utils.NewStringPointer("Bar Value")},
			{Key: "baz", KeyNiceName: utils.NewStringPointer("Baz Value")},
			{Key: "qux", KeyNiceName: utils.NewStringPointer("Qux Value")},
		},
		Stateless: true, // Indicates that the event does not maintain any state between emissions.
	}

	// Register the MQTT event with the application. This step declares the event structure to the system,
	// allowing it to recognize and properly route the event when sent.
	mqtt_event_id, err := app.AddCameraPlatformEvent(mqttEvent)
	if err != nil {
		app.Syslog.Critf("Error adding event declaration: %s", err.Error()) // Log critical error if the event declaration fails.
	}

	// Periodically send the MQTT event every second. This goroutine simulates an ongoing event emission,
	// like sensor readings or status updates.
	go func() {
		for true {
			time.Sleep(1 * time.Second) // Wait for 1 second before sending the next event.

			// Attempt to send a newly created MQTT event with dynamic values.
			if err := app.SendPlatformEvent(mqtt_event_id, func() (*axevent.AXEvent, error) {
				// Dynamically generate event data with random values.
				return mqttEvent.NewEvent(acapapp.KeyValueMap{
					"foo": rand.Int(),          // Random integer value.
					"bar": rand.Float64(),      // Random floating-point value.
					"baz": "i want to be json", // Static string value.
					"qux": rand.Intn(2) == 1,   // Random boolean value (true or false).
				})
			}); err != nil {
				app.Syslog.Errorf("Error sending event: %s", err.Error()) // Log error if event sending fails.
			} else {
				app.Syslog.Infof("Event send") // Log information on successful event send.
			}
		}
	}()

	// Execute the main loop of the application. This call blocks and allows the application to continue running,
	// responding to system events, and handling the periodic event emissions initiated above.
	app.Run()
}
