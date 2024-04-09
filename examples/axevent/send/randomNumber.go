package main

import (
	"math/rand/v2"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Declare event with randomNumbersEventKvs as stateless
//
// Works like https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax_event_stateless_declaration_example_8c-example.html
func declareRandomNumbersEvent(app *acapapp.AcapApplication) (int, error) {
	random_numbers_kvs, err := axevent.NewCameraApplicationPlatformEvent(
		app.Manifest.ACAPPackageConf.Setup,
		"RandomNumbers",
		utils.NewStringPointer("Random Numbers"),
		[]*axevent.KeyValueEntrie{
			{Key: "random_int", ValueType: axevent.AXValueTypeInt},
			{Key: "random_float", ValueType: axevent.AXValueTypeDouble},
		},
		nil,
		[]*axevent.AxEventKeyValueSetDataMark{
			{Key: "random_int"},
			{Key: "random_float"},
		},
		nil,
		[]*axevent.AxEventKeyValueSetNiceNames{
			{Key: "random_int", KeyNiceName: utils.NewStringPointer("Random Integer Number")},
			{Key: "random_float", KeyNiceName: utils.NewStringPointer("Random Float Number")},
		},
	)
	if err != nil {
		return 0, err
	}
	stateless := true
	event_id, err := app.EventHandler.Declare(random_numbers_kvs, stateless, func(subscription int, userdata any) {
		app.Syslog.Infof("Random Number declaration complete, subscription: %d", subscription)
	}, nil)
	return event_id, err
}

// Create a new random number event
func newRandomNumberEvent() *axevent.AXEvent {
	return axevent.NewAxEvent(axevent.NewAXEventKeyValueSetFromEntries([]axevent.KeyValueEntrie{
		{Key: "random_int", Value: rand.IntN(100), ValueType: axevent.AXValueTypeInt},
		{Key: "random_float", Value: rand.Float64(), ValueType: axevent.AXValueTypeDouble},
	}), nil)

}
