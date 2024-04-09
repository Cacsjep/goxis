package main

import (
	"math/rand"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// Declare event with featureEventKvs as statefull
//
// Works like https://axiscommunications.github.io/acap-documentation/3.5/api/axevent/html/ax_event_property_state_declaration_example_8c-example.html
func declareFeatureEvent(app *acapapp.AcapApplication) (int, error) {
	feature_kvs, err := axevent.NewCameraApplicationPlatformEvent(
		app.Manifest.ACAPPackageConf.Setup,
		"Feature",
		nil,
		[]*axevent.KeyValueEntrie{
			{Key: "feature", Value: "myfeature", ValueType: axevent.AXValueTypeString},
			{Key: "enabled", Namespace: &axevent.OnfivNameSpaceTnsAxis, Value: true, ValueType: axevent.AXValueTypeBool},
		},
		[]*axevent.AxEventKeyValueSetSourceMark{
			{Key: "feature"},
		},
		[]*axevent.AxEventKeyValueSetDataMark{
			{Key: "enabled", Namespace: &axevent.OnfivNameSpaceTnsAxis},
		},
		[]*axevent.AxEventKeyValueSetUserDefineMark{
			{Key: "feature", Tag: &axevent.OnfivTagOnKeyValue},
			{Key: "topic1", Namespace: &axevent.OnfivNameSpaceTnsAxis, Tag: &axevent.OnfivTagOnKeyValue},
		},
		[]*axevent.AxEventKeyValueSetNiceNames{
			{Key: "feature", KeyNiceName: utils.NewStringPointer("Feature"), ValueNiceName: utils.NewStringPointer("My Feature")},
			{Key: "enabled", Namespace: &axevent.OnfivNameSpaceTnsAxis, KeyNiceName: utils.NewStringPointer("Enabled")},
		},
	)
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
	kvs := axevent.NewAXEventKeyValueSetFromEntries([]axevent.KeyValueEntrie{
		{Key: "feature", Value: "myfeature", ValueType: axevent.AXValueTypeString},
		{Key: "enabled", Namespace: &axevent.OnfivNameSpaceTnsAxis, Value: rand.Intn(2) == 1, ValueType: axevent.AXValueTypeBool},
	})
	new_event := axevent.NewAxEvent(kvs, nil)
	defer kvs.Free()
	return new_event
}
