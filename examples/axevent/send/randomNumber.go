package main

import (
	"math/rand/v2"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Declare an event with two random numbers
func randomNumbersEventKvs() (*axevent.AXEventKeyValueSet, error) {
	kvs, err := axevent.NewTnsAxisEvent("CameraApplicationPlatform", "SendEventExample", utils.NewStringPointer("RandomNumbers"), nil, []*axevent.KeyValueEntrie{
		{Key: "random_int", ValueType: axevent.AXValueTypeInt},
		{Key: "random_float", ValueType: axevent.AXValueTypeDouble},
	})
	if err != nil {
		return nil, err
	}
	kvs.MarkAsSource("random_int", nil)
	kvs.MarkAsUserDefined("random_int", nil, nil)

	kvs.MarkAsSource("random_float", nil)
	kvs.MarkAsUserDefined("random_float", nil, nil)

	// Set a nice name so it looks better in the webinterface
	if err := kvs.AddNiceNames(
		"random_int",
		nil,
		utils.NewStringPointer("Random Integer Number"),
		nil,
	); err != nil {
		return nil, err
	}

	// Set a nice name so it looks better in the webinterface
	if err := kvs.AddNiceNames(
		"random_float",
		nil,
		utils.NewStringPointer("Random Float Number"),
		nil,
	); err != nil {
		return nil, err
	}

	// Used to display a nice name in the event list when adding a rule via webinterface
	// like SendEventExample: RandomNumbers
	if err := kvs.AddNiceNames(
		"topic2",
		utils.NewStringPointer(axevent.OnfivNameSpaceTnsAxis),
		nil,
		utils.NewStringPointer("SendEventExample: RandomNumbers"),
	); err != nil {
		return nil, err
	}

	return kvs, nil
}

// Declare event with randomNumbersEventKvs as stateless
//
// Works like https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax_event_stateless_declaration_example_8c-example.html
func declareRandomNumbersEvent(app *acapapp.AcapApplication) (int, error) {
	random_numbers_kvs, err := randomNumbersEventKvs()
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
	kvs := axevent.NewAXEventKeyValueSet()
	kvs.AddKeyValue("random_int", nil, rand.IntN(100), axevent.AXValueTypeInt)
	kvs.AddKeyValue("random_float", nil, rand.Float64(), axevent.AXValueTypeDouble)
	new_event := axevent.NewAxEvent(kvs, nil)
	defer kvs.Free()
	return new_event
}
