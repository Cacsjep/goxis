package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Tipp: Use Axis Metadata Monitor to see live which events are produced by camera
// https://www.axis.com/developer-community/axis-metadata-monitor
func main() {
	app := acapapp.NewAcapApplication()

	// VirtualInputEventKvs is a helper function to create a AXEventKeyValueSet for a VirtualInput event.
	// You can build your own AXEventKeyValueSet with the NewTns1AxisEvent or NewTnsAxisEvent.
	vio_event, err := axevent.VirtualInputEventKvs(utils.NewIntPointer(1), nil) // We pass nil because we want to listen to all input states. If you want to listen to a specific state, you can pass utils.NewBoolPointer(true) or utils.NewBoolPointer(false)
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
	vio_subscription_id, err := app.EventHandler.OnEvent(vio_event, func(e *axevent.Event) {

		// You can also build your own events =)
		var vi axevent.VirtualInputEvent
		// You could aslo read manually from the event kvs like
		// e.Kvs.GetInteger("port", nil) or
		// e.Kvs.GetBoolean("active", nil)
		if err := axevent.UnmarshalEvent(e, &vi); err != nil {
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
	app.Run()
}
