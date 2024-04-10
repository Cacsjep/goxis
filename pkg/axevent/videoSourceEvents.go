package axevent

import "github.com/Cacsjep/goxis/pkg/utils"

func DayNightEventKvs(videoSourceConfigurationToken *int, day *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "DayNightVision", nil, nil, []*KeyValueEntrie{
		{Key: "VideoSourceConfigurationToken", Value: videoSourceConfigurationToken, ValueType: AXValueTypeInt},
		{Key: "day", Value: day, ValueType: AXValueTypeBool},
	})
}

type DayNightEvent struct {
	VideoSourceConfigurationToken int  `eventKey:"VideoSourceConfigurationToken"`
	Day                           bool `eventKey:"day"`
}

func LiveStreamAccessedEventKvs() (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "LiveStreamAccessed", nil, nil, []*KeyValueEntrie{
		{Key: "accessed", ValueType: AXValueTypeBool},
	})
}

type LiveStreamAccessedEvent struct {
	Accessed bool `eventKey:"accessed"`
}

func AutofocusEventKvs(videoSourceConfigurationToken *int) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "Autofocus", nil, nil, []*KeyValueEntrie{
		{Key: "VideoSourceConfigurationToken", Value: videoSourceConfigurationToken, ValueType: AXValueTypeInt},
		{Key: "focus", ValueType: AXValueTypeDouble},
	})
}

type AutofocusEvent struct {
	VideoSourceConfigurationToken int     `eventKey:"VideoSourceConfigurationToken"`
	Focus                         float64 `eventKey:"focus"`
}

func TamperingEventKvs(channel *int, tampering *int) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "Tampering", nil, nil, []*KeyValueEntrie{
		{Key: "channel", Value: channel, ValueType: AXValueTypeInt},
		{Key: "tampering", Value: tampering, ValueType: AXValueTypeInt},
	})
}

type TamperingEvent struct {
	Channel   int `eventKey:"channel"`
	Tampering int `eventKey:"tampering"`
}

func MotionAlarmEventKvs(source *string, state *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "MotionAlarm", nil, nil, []*KeyValueEntrie{
		{Key: "Source", Value: source, ValueType: AXValueTypeString},
		{Key: "State", Value: state, ValueType: AXValueTypeBool},
	})
}

type MotionAlarmEvent struct {
	Source string `eventKey:"Source"`
	State  bool   `eventKey:"State"`
}

func GlobalSceneChangeEventKvs(source *string, state *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("VideoSource", "GlobalSceneChange", utils.StrPtr("ImagingService"), nil, []*KeyValueEntrie{
		{Key: "Source", Value: source, ValueType: AXValueTypeString},
		{Key: "State", Value: state, ValueType: AXValueTypeBool},
	})
}

type GlobalSceneChangeEvent struct {
	Source string `eventKey:"Source"`
	State  bool   `eventKey:"State"`
}
