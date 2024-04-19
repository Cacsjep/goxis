package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// This example demonstrates how to subscribe to a VirtualInput event.
// axevent holds some predefinied events like VirtualInputEvent, just create your own Events
// when u need a specific event. Look how VirtualInputEventKvs is build in axevent package.
// using UnmarshalEvent you can convert the event like json.Unmarshal to a struct.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/axevent/subscribe_to_event
//
// Tipp: Use Axis Metadata Monitor to see live which events are produced by camera
// https://www.axis.com/developer-community/axis-metadata-monitor
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	// VirtualInputEventKvs is a helper function to create a AXEventKeyValueSet for a VirtualInput event.
	// You can build your own AXEventKeyValueSet with the NewTns1AxisEvent or NewTnsAxisEvent.
	vio_event, err := axevent.VirtualInputEventKvs(utils.IntPtr(1), nil) // We pass nil because we want to listen to all input states. If you want to listen to a specific state, you can pass utils.BoolPtr(true) or utils.BoolPtr(false)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// OnEvent create a subscription callback for the given event key value set.
	// You can test via changing the state of the virtual input via:
	// Activate: 	http://<ip>/axis-cgi/virtualinput/activate.cgi?schemaversion=1&port=1
	// Deactivate:  http://<ip>/axis-cgi/virtualinput/deactivate.cgi?schemaversion=1&port=1
	// A note on callback functions:
	//  	Any call to axparam in the callback should again should be done via a goroutine.
	//  	Otherwise, the callback will block the event handler.
	vio_subscription_id, err := app.OnEvent(vio_event, func(e *axevent.Event) {
		// You can also build your own events =)
		var vi axevent.VirtualInputEvent
		// You could aslo read manually from the event kvs like
		// e.Kvs.GetInteger("port", nil) or
		// e.Kvs.GetBoolean("active", nil)
		if err := acapapp.UnmarshalEvent(e, &vi); err != nil {
			app.Syslog.Error(err.Error())
			return
		}
		app.Syslog.Infof("VirtualInput Port: %d, Active: %t", vi.Port, vi.Active)
	})
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	app.Syslog.Infof("VirtualInput subscription id: %d", vio_subscription_id)

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axevent needs a running event loop to handle the events callbacks corretly
	app.Run()
}
