package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Tipp: Use Axis Metadata Monitor to see live with events are produced by camera
// https://www.axis.com/developer-community/axis-metadata-monitor
func main() {
	app := acapapp.NewAcapApplication()

	// VirtualInputEvent is a helper function to create a AXEventKeyValueSet for a VirtualInput event.
	vio_event, err := axevent.VirtualInputEventKvs(utils.NewIntPointer(1), nil)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// DayNightEventKvs is a helper function to create a AXEventKeyValueSet for a DayNight event.
	dn_event, err := axevent.DayNightEventKvs(utils.NewIntPointer(1), nil)
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
		var vi axevent.VirtualInputEvent
		if err := axevent.UnmarshalEvent(e, &vi); err != nil {
			app.Syslog.Error(err.Error())
			return
		}
		app.Syslog.Infof("VirtualInput Port: %d, Active: %t", vi.Port, vi.Active)
	})

	dn_subscription_id, err := app.EventHandler.OnEvent(dn_event, func(e *axevent.Event) {
		var dn axevent.DayNightEvent
		if err := axevent.UnmarshalEvent(e, &dn); err != nil {
			app.Syslog.Error(err.Error())
			return
		}
		app.Syslog.Infof("DayNight, VideoSource: %d, Day: %t", dn.VideoSourceConfigurationToken, dn.Day)
	})

	app.Syslog.Infof("VirtualInput SubsId: %d, DayNight SubsId: %d", vio_subscription_id, dn_subscription_id)

	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Signal handler automatically internally created for SIGTERM, SIGINT
	app.Run()
}
