package axevent

// Topic: Device Event: SupervisedPort, IsProp: true, Source: {port xsd:int}, Data {tampered xsd:boolean}
func SupervisedPortEventKvs(port int, tampered *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("Device", "IO", NewStringPointer("SupervisedPort"), nil, []*KeyValueEntrie{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "tampered", value: tampered, value_type: AXValueTypeBool},
	})
}

type SupervisedPortEvent struct {
	Port     int  `eventKey:"port"`
	Tampered bool `eventKey:"tampered"`
}

// Topic: Device Event: VirtualInput, IsProp: true, Source: {port xsd:int}, Data {active xsd:boolean}
func VirtualInputEventKvs(port int, active *bool) (*AXEventKeyValueSet, error) {
	return tnsAxisEvent("Device", "IO", NewStringPointer("VirtualInput"), nil, []*KeyValueEntrie{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "active", value: active, value_type: AXValueTypeBool},
	})
}

type VirtualInputEvent struct {
	Port   int  `eventKey:"port"`
	Active bool `eventKey:"active"`
}
