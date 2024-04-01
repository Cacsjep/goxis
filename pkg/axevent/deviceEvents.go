package axevent

// Topic: Device Event: SupervisedPort, IsProp: true, Source: {port xsd:int}, Data {tampered xsd:boolean}
func SupervisedPortEvent(port int, tampered *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("Device", "IO", NewStringPointer("SupervisedPort"), nil, []*keyvalues{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "tampered", value: tampered, value_type: AXValueTypeBool},
	})
}

// Topic: Device Event: VirtualInput, IsProp: true, Source: {port xsd:int}, Data {active xsd:boolean}
func VirtualInputEvent(port int, active *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("Device", "IO", NewStringPointer("VirtualInput"), nil, []*keyvalues{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "active", value: active, value_type: AXValueTypeBool},
	})
}
