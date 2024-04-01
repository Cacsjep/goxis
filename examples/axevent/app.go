package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
)

var (
	err              error
	vio_subscription int
	dn_subscription  int
	app              *acapapp.AcapApplication
)

// This example uses axevent library for subscribing to an ONVIF event.
// Tipp: Use Axis Metadata Monitor to see live with events are produced by camera
//
//	https://www.axis.com/developer-community/axis-metadata-monitor
func main() {
	app = acapapp.NewAcapApplication()

	/* Initialize an AXEventKeyValueSet that matches Virtual Input 1.
	 *
	 *    tns1:topic0=Device
	 * tnsaxis:topic1=IO
	 * tnsaxis:topic2=VirtualInput
	 *           port=1    		 <-- Subscribe to port number 1
	 *           active=NULL     <-- Subscribe to all states
	 */
	vio_event, err := axevent.VirtualInputEvent(1, nil)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// VirtualInputEvent is a helper function to create a AXEventKeyValueSet for a VirtualInput event.
	/* func VirtualInputEvent(port int) (*AXEventKeyValueSet, error) {
		vio_event := NewAXEventKeyValueSet()
		if err := vio_event.AddKeyValue("topic0", &OnfivNameSpaceTns1, "Device", AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic0: %w", err)
		}
		if err := vio_event.AddKeyValue("topic1", &OnfivNameSpaceTnsAxis, "IO", AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic1: %w", err)
		}
		if err := vio_event.AddKeyValue("topic2", &OnfivNameSpaceTnsAxis, "VirtualInput", AXValueTypeString); err != nil {
			return nil, fmt.Errorf("failed to add key-value for topic2: %w", err)
		}
		if err := vio_event.AddKeyValue("port", nil, port, AXValueTypeInt); err != nil {
			return nil, fmt.Errorf("failed to add key-value for port: %w", err)
		}
		if err := vio_event.AddKeyValue("active", nil, true, AXValueTypeBool); err != nil { // Assuming active is always true, as nil was passed for value before
			return nil, fmt.Errorf("failed to add key-value for active: %w", err)
		}
		return vio_event, nil
	} */

	// OnEvent create a subscription callback for the given event key value set.
	// You can test via changing the state of the virtual input via:
	// Activate: 	http://<ip>/axis-cgi/virtualinput/activate.cgi?schemaversion=1&port=1
	// Deactivate:  http://<ip>/axis-cgi/virtualinput/deactivate.cgi?schemaversion=1&port=1
	// A note on callback functions:
	//  	The callback functions registered with the AXEventHandler
	//		will be called from the GMainLoop thread in the default context.
	//		This means that the client may not prevent callback functions from returning,
	//		nor should any lengthy processing be made in the callback functions.
	//		Failure to comply with this convention will prevent the event system from,
	//		or delay it in, sending or delivering any more events to the calling application.
	//		For this reason, it is recommended to use a gorutine for any processing that may take time.
	vio_subscription, err = app.EventHandler.OnEvent(vio_event, func(e *axevent.Event) {

		// Get the port value
		port, err := e.Kvs.GetInteger("port", nil)
		if err != nil {
			app.Syslog.Error("Unable to get port value from event key value set")
			return
		}

		// Get the active value
		active, err := e.Kvs.GetBoolean("active", nil)
		if err != nil {
			app.Syslog.Error("Unable to get active value from event key value set")
			return
		}

		app.Syslog.Infof(
			"VIO Callback, Port: %d, Active: %t, Timestamp: %s",
			port,
			active,
			e.Timestamp.Format("2006-01-02 15:04:05"),
		)
	})
	source := 1 // Starts with 1
	dn_event, err := axevent.DayNightEvent(&source, nil)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	dn_subscription, err = app.EventHandler.OnEvent(dn_event, func(e *axevent.Event) {
		day, err := e.Kvs.GetBoolean("day", nil)
		if err != nil {
			app.Syslog.Error("Unable to get day value from event key value set")
			return
		}
		app.Syslog.Infof(
			"DN Callback, Day: %t, Timestamp: %s",
			day,
			e.Timestamp.Format("2006-01-02 15:04:05"),
		)
	})

	app.Syslog.Infof("VIO Subscription ID: %d", vio_subscription)
	app.Syslog.Infof("DN Subscription ID: %d", dn_subscription)

	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	// Signal handler automatically internally created for SIGTERM, SIGINT
	// This blocks now the main thread.
	app.Run()
}
