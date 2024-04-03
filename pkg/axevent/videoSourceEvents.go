package axevent

import "github.com/Cacsjep/goxis/pkg/utils"

func DayNightEventKvs(videoSourceConfigurationToken *int, day *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "DayNightVision", nil, nil, []*KeyValueEntrie{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "day", value: day, value_type: AXValueTypeBool},
	})
}

type DayNightEvent struct {
	VideoSourceConfigurationToken int  `eventKey:"VideoSourceConfigurationToken"`
	Day                           bool `eventKey:"day"`
}

func LiveStreamAccessedEventKvs() (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "LiveStreamAccessed", nil, nil, []*KeyValueEntrie{
		{key: "accessed", value_type: AXValueTypeBool},
	})
}

type LiveStreamAccessedEvent struct {
	Accessed bool `eventKey:"accessed"`
}

func AutofocusEventKvs(videoSourceConfigurationToken *int) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "Autofocus", nil, nil, []*KeyValueEntrie{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "focus", value_type: AXValueTypeDouble},
	})
}

type AutofocusEvent struct {
	VideoSourceConfigurationToken int     `eventKey:"VideoSourceConfigurationToken"`
	Focus                         float64 `eventKey:"focus"`
}

func TamperingEventKvs(channel *int, tampering *int) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "Tampering", nil, nil, []*KeyValueEntrie{
		{key: "channel", value: channel, value_type: AXValueTypeInt},
		{key: "tampering", value: tampering, value_type: AXValueTypeInt},
	})
}

type TamperingEvent struct {
	Channel   int `eventKey:"channel"`
	Tampering int `eventKey:"tampering"`
}

func MotionAlarmEventKvs(source *string, state *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "MotionAlarm", nil, nil, []*KeyValueEntrie{
		{key: "Source", value: source, value_type: AXValueTypeString},
		{key: "State", value: state, value_type: AXValueTypeBool},
	})
}

type MotionAlarmEvent struct {
	Source string `eventKey:"Source"`
	State  bool   `eventKey:"State"`
}

func GlobalSceneChangeEventKvs(source *string, state *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "GlobalSceneChange", utils.NewStringPointer("ImagingService"), nil, []*KeyValueEntrie{
		{key: "Source", value: source, value_type: AXValueTypeString},
		{key: "State", value: state, value_type: AXValueTypeBool},
	})
}

type GlobalSceneChangeEvent struct {
	Source string `eventKey:"Source"`
	State  bool   `eventKey:"State"`
}
