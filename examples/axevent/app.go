package main

import (
	"github.com/Cacsjep/goxis"
	"github.com/Cacsjep/goxis/pkg/axevent"
)

var (
	err          error
	subscription int
	app          *goxis.AcapApplication
)

// This example uses axevent library for subscribing to an ONVIF event.
// Tipp: Use Axis Metadata Monitor to see live with events are produced by camera
//
//	https://www.axis.com/developer-community/axis-metadata-monitor
func main() {
	app = goxis.NewAcapApplication()

	/* Initialize an AXEventKeyValueSet that matches Virtual Input 1.
	 *
	 *    tns1:topic0=Device
	 * tnsaxis:topic1=IO
	 * tnsaxis:topic2=VirtualInput
	 *           port=1    		 <-- Subscribe to port number 1
	 *           active=NULL     <-- Subscribe to all states
	 */

	set := axevent.NewAXEventKeyValueSet()
	err = set.AddKeyValue("topic0", &axevent.OnfivNameSpaceTns1, "Device", axevent.AXValueTypeString)
	err = set.AddKeyValue("topic1", &axevent.OnfivNameSpaceTnsAxis, "IO", axevent.AXValueTypeString)
	err = set.AddKeyValue("topic2", &axevent.OnfivNameSpaceTnsAxis, "VirtualInput", axevent.AXValueTypeString)
	err = set.AddKeyValue("port", nil, 1, axevent.AXValueTypeInt)
	err = set.AddKeyValue("active", nil, nil, axevent.AXValueTypeBool)

	// Subscribe to the event.
	// You can test the callback via changing the state of the virtual input via:
	// Activate: 	http://<ip>/axis-cgi/virtualinput/activate.cgi?schemaversion=1&port=1
	// Deactivate:  http://<ip>/axis-cgi/virtualinput/deactivate.cgi?schemaversion=1&port=1
	// A note on callback functions:
	//  	The callback functions registered with the AXEventHandler
	//		will be called from the GMainLoop thread in the default context.
	//		This means that the client may not prevent callback functions from returning,
	//		nor should any lengthy processing be made in the callback functions.
	//		Failure to comply with this convention will prevent the event system from,
	//		or delay it in, sending or delivering any more events to the calling application.
	subscription, err = app.EventHandler.Subscribe(set, func(subscription int, event *axevent.AXEvent, userdata any) {

		// Get the key value set from event
		kvs := event.GetKeyValueSet()

		// Get the port value
		port, err := kvs.GetInteger("port", nil)
		if err != nil {
			app.Syslog.Error("Unable to get port value from event key value set")
			return
		}

		// Get the active value
		active, err := kvs.GetBoolean("active", nil)
		if err != nil {
			app.Syslog.Error("Unable to get active value from event key value set")
			return
		}

		app.Syslog.Infof(
			"VIO Callback, Port: %d, Active: %t, Timestamp: %s, Userdata: %s",
			port,
			active,
			event.GetTimestamp().Format("2006-01-02 15:04:05"),
			userdata,
		)

	}, "my importand user data")

	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Signal handler automatically internally created for SIGTERM, SIGINT
	// This blocks now the main thread.
	app.Run()
}
