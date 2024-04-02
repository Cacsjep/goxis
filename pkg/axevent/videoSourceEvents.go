package axevent

//Topic: VideoSource Event: DayNightVision, IsProp: true, Source: {VideoSourceConfigurationToken xsd:int}, Data {day xsd:boolean}
func DayNightEventKvs(videoSourceConfigurationToken *int, day *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "DayNightVision", nil, nil, []*KeyValueEntrie{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "day", value: day, value_type: AXValueTypeBool},
	})
}

type DayNightEvent struct {
	VideoSourceConfigurationToken int  `eventKey:"VideoSourceConfigurationToken"`
	Day                           bool `eventKey:"day"`
}

//Topic: VideoSource Event: LiveStreamAccessed, IsProp: true, Source: { }, Data {accessed xsd:boolean}
func LiveStreamAccessedEventKvs() (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "LiveStreamAccessed", nil, nil, []*KeyValueEntrie{
		{key: "accessed", value_type: AXValueTypeBool},
	})
}

type LiveStreamAccessedEvent struct {
	Accessed bool `eventKey:"accessed"`
}

//Topic: VideoSource Event: Autofocus, IsProp: true, Source: {VideoSourceConfigurationToken xsd:int}, Data {focus xsd:double}
func AutofocusEventKvs(videoSourceConfigurationToken *int) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "Autofocus", nil, nil, []*KeyValueEntrie{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "focus", value_type: AXValueTypeDouble},
	})
}

type AutofocusEvent struct {
	VideoSourceConfigurationToken int     `eventKey:"VideoSourceConfigurationToken"`
	Focus                         float64 `eventKey:"focus"`
}

// Topic: VideoSource Event: Tampering, IsProp: false, Source: {channel xsd:int}, Data {tampering xsd:int}
func TamperingEventKvs(channel *int) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "Tampering", nil, nil, []*KeyValueEntrie{
		{key: "channel", value: channel, value_type: AXValueTypeInt},
		{key: "tampering", value_type: AXValueTypeInt},
	})
}

type TamperingEvent struct {
	Channel   int `eventKey:"channel"`
	Tampering int `eventKey:"tampering"`
}

// Topic: VideoSource Event: MotionAlarm, IsProp: true, Source: {Source tt:ReferenceToken}, Data {State xsd:boolean}
func MotionAlarmEventKvs(source *string) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "MotionAlarm", nil, nil, []*KeyValueEntrie{
		{key: "Source", value: source, value_type: AXValueTypeString},
		{key: "State", value_type: AXValueTypeBool},
	})
}

type MotionAlarmEvent struct {
	Source string `eventKey:"Source"`
	State  bool   `eventKey:"State"`
}
