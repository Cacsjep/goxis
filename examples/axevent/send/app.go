package main

import (
	"time"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
)

// Tipp: Use Axis Metadata Monitor to see live which events are produced by camera
// https://www.axis.com/developer-community/axis-metadata-monitor
func main() {
	app := acapapp.NewAcapApplication()

	random_numbers_event_id, err := declareRandomNumbersEvent(app)
	if err != nil {
		app.Syslog.Critf("Error declaring random numbers event: %s", err.Error())
	}
	app.Syslog.Infof("Random numbers event declared with id: %d", random_numbers_event_id)

	feature_event_id, err := declareFeatureEvent(app)
	if err != nil {
		app.Syslog.Critf("Error declaring features event: %s", err.Error())
	}
	app.Syslog.Infof("Features event declared with id: %d", feature_event_id)

	go func() {
		for true {
			time.Sleep(1 * time.Second)
			sendEvent(app, "Random", random_numbers_event_id, newRandomNumberEvent())
			sendEvent(app, "Feature", feature_event_id, newFeatureEvent())
		}
	}()

	app.AddCloseCleanFunc(func() {
		app.EventHandler.Undeclare(random_numbers_event_id)
		app.EventHandler.Undeclare(feature_event_id)
	})
	app.Run()
}

// Send event to the event handler
func sendEvent(app *acapapp.AcapApplication, event_name string, event_id int, event *axevent.AXEvent) {
	if err := app.EventHandler.SendEvent(event_id, event); err != nil {
		app.Syslog.Errorf("Error sending %s event: %s", event_name, err.Error())
	} else {
		app.Syslog.Infof("Event %s send", event_name)
	}
}
