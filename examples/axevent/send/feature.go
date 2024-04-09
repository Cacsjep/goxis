package main

import (
	"math/rand"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Declare an event with two random numbers
func featureEventKvs() (*axevent.AXEventKeyValueSet, error) {
	kvs, err := axevent.NewTnsAxisEvent("CameraApplicationPlatform", "SendEventExample", utils.NewStringPointer("Feature"), nil, []*axevent.KeyValueEntrie{
		{Key: "feature", Value: "myfeature", ValueType: axevent.AXValueTypeString},
		{Key: "enabled", Namespace: &axevent.OnfivNameSpaceTnsAxis, Value: true, ValueType: axevent.AXValueTypeBool},
	})
	if err != nil {
		return nil, err
	}

	// ax_event_key_value_set_mark_as_source(set, "feature", NULL, NULL);
	if err := kvs.MarkAsSource("feature", nil); err != nil {
		return nil, err
	}

	// ax_event_key_value_set_mark_as_data(set, "enabled", "tnsaxis", NULL);
	if err := kvs.MarkAsData("enabled", utils.NewStringPointer(axevent.OnfivNameSpaceTnsAxis)); err != nil {
		return nil, err
	}

	// ax_event_key_value_set_mark_as_user_defined(set, "feature", NULL, "tag-on-key-value", NULL);
	if err := kvs.MarkAsUserDefined(
		"feature",
		nil,
		utils.NewStringPointer("tag-on-key-value"),
	); err != nil {
		return nil, err
	}

	// Set a nice name and value so it looks better in the webinterface
	if err := kvs.AddNiceNames(
		"feature",
		nil,
		utils.NewStringPointer("Feature"),
		utils.NewStringPointer("My Feature"),
	); err != nil {
		return nil, err
	}

	// ax_event_key_value_set_mark_as_user_defined(set, "topic1", "tnsaxis", "tag-on-key-value", NULL);
	if err := kvs.MarkAsUserDefined(
		"topic1",
		utils.NewStringPointer(axevent.OnfivNameSpaceTnsAxis),
		utils.NewStringPointer("tag-on-key-value"),
	); err != nil {
		return nil, err
	}

	// ax_event_key_value_set_add_nice_names(set, "enabled", "tnsaxis", "Key nice name", "Value nice name", NULL);
	if err := kvs.AddNiceNames(
		"enabled",
		utils.NewStringPointer(axevent.OnfivNameSpaceTnsAxis),
		utils.NewStringPointer("Enabled"),
		nil,
	); err != nil {
		return nil, err
	}

	// Used to display a nice name in the event list when adding a rule via webinterface
	// like SendEventExample: Feature Event
	if err := kvs.AddNiceNames(
		"topic2",
		utils.NewStringPointer(axevent.OnfivNameSpaceTnsAxis),
		nil,
		utils.NewStringPointer("SendEventExample: Feature Event"),
	); err != nil {
		return nil, err
	}
	return kvs, nil
}

// Declare event with featureEventKvs as statefull
//
// Works like https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax_event_property_state_declaration_example_8c-example.html
func declareFeatureEvent(app *acapapp.AcapApplication) (int, error) {
	feature_kvs, err := featureEventKvs()
	if err != nil {
		return 0, err
	}
	stateless := false /* statefull, Indicated a property state */
	event_id, err := app.EventHandler.Declare(feature_kvs, stateless, func(subscription int, userdata any) {
		app.Syslog.Infof("Feature Event declaration complete, subscription: %d", subscription)
	}, nil)
	return event_id, err
}

// Create a new random number event
func newFeatureEvent() *axevent.AXEvent {
	kvs := axevent.NewAXEventKeyValueSet()
	kvs.AddKeyValue("feature", nil, "myfeature", axevent.AXValueTypeString)
	kvs.AddKeyValue("enabled", &axevent.OnfivNameSpaceTnsAxis, rand.Intn(2) == 1, axevent.AXValueTypeBool)
	new_event := axevent.NewAxEvent(kvs, nil)
	defer kvs.Free()
	return new_event
}
