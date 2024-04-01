package axevent

//Topic: VideoSource Event: DayNightVision, IsProp: true, Source: {VideoSourceConfigurationToken xsd:int}, Data {day xsd:boolean}
func DayNightEvent(videoSourceConfigurationToken *int, day *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "DayNightVision", nil, nil, []*keyvalues{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "day", value: day, value_type: AXValueTypeBool},
	})
}

//Topic: VideoSource Event: LiveStreamAccessed, IsProp: true, Source: { }, Data {accessed xsd:boolean}
func LiveStreamAccessedEvent() (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "LiveStreamAccessed", nil, nil, []*keyvalues{
		{key: "accessed", value_type: AXValueTypeBool},
	})
}

//Topic: VideoSource Event: Autofocus, IsProp: true, Source: {VideoSourceConfigurationToken xsd:int}, Data {focus xsd:double}
func AutofocusEvent(videoSourceConfigurationToken *int) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "Autofocus", nil, nil, []*keyvalues{
		{key: "VideoSourceConfigurationToken", value: videoSourceConfigurationToken, value_type: AXValueTypeInt},
		{key: "focus", value_type: AXValueTypeDouble},
	})
}

// Topic: VideoSource Event: Tampering, IsProp: false, Source: {channel xsd:int}, Data {tampering xsd:int}
func TamperingEvent(channel *int) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "Tampering", nil, nil, []*keyvalues{
		{key: "channel", value: channel, value_type: AXValueTypeInt},
		{key: "tampering", value_type: AXValueTypeInt},
	})
}

// Topic: VideoSource Event: MotionAlarm, IsProp: true, Source: {Source tt:ReferenceToken}, Data {State xsd:boolean}
func MotionAlarmEvent(source *string) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("VideoSource", "MotionAlarm", nil, nil, []*keyvalues{
		{key: "Source", value: source, value_type: AXValueTypeString},
		{key: "State", value_type: AXValueTypeBool},
	})
}
