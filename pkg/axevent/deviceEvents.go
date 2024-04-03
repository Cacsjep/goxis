package axevent

import "github.com/Cacsjep/goxis/pkg/utils"

func SupervisedPortEventKvs(port *int, tampered *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "IO", utils.NewStringPointer("SupervisedPort"), nil, []*KeyValueEntrie{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "tampered", value: tampered, value_type: AXValueTypeBool},
	})
}

type SupervisedPortEvent struct {
	Port     int  `eventKey:"port"`
	Tampered bool `eventKey:"tampered"`
}

func VirtualInputEventKvs(port *int, active *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "IO", utils.NewStringPointer("VirtualInput"), nil, []*KeyValueEntrie{
		{key: "port", value: port, value_type: AXValueTypeInt},
		{key: "active", value: active, value_type: AXValueTypeBool},
	})
}

type VirtualInputEvent struct {
	Port   int  `eventKey:"port"`
	Active bool `eventKey:"active"`
}

func StorageFailureEventKvs(disk_id *string, disruption *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "HardwareFailure", utils.NewStringPointer("StorageFailure"), nil, []*KeyValueEntrie{
		{key: "disk_id", value: disk_id, value_type: AXValueTypeString},
		{key: "disruption", value: disruption, value_type: AXValueTypeBool},
	})
}

type StorageFailureEvent struct {
	DiskId     string `eventKey:"disk_id"`
	Disruption bool   `eventKey:"disruption"`
}

func HeaterStatusEventKvs(heater *int, running *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Heater", utils.NewStringPointer("Status"), nil, []*KeyValueEntrie{
		{key: "heater", value: heater, value_type: AXValueTypeInt},
		{key: "running", value: running, value_type: AXValueTypeBool},
	})
}

type HeaterStatusEvent struct {
	Heater  int  `eventKey:"heater"`
	Running bool `eventKey:"running"`
}

func SystemReadyStatusEventKvs(ready *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Status", utils.NewStringPointer("SystemReady"), nil, []*KeyValueEntrie{
		{key: "ready", value: ready, value_type: AXValueTypeBool},
	})
}

type SystemReadyStatusEvent struct {
	Ready int `eventKey:"ready"`
}

func TriggerRelayEventKvs(relayToken *int, logicalState *string) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Trigger", utils.NewStringPointer("Relay"), nil, []*KeyValueEntrie{
		{key: "RelayToken", value: relayToken, value_type: AXValueTypeInt},
		{key: "LogicalState", value: logicalState, value_type: AXValueTypeString},
	})
}

type TriggerRelayEvent struct {
	RelayToken   int    `eventKey:"RelayToken"`
	LogicalState string `eventKey:"LogicalState"`
}

func DigitalInputEventKvs(inputToken *int, logicalState *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Trigger", utils.NewStringPointer("DigitalInput"), nil, []*KeyValueEntrie{
		{key: "InputToken", value: inputToken, value_type: AXValueTypeInt},
		{key: "LogicalState", value: logicalState, value_type: AXValueTypeBool},
	})
}

type DigitalInputEvent struct {
	InputToken   int  `eventKey:"InputToken"`
	LogicalState bool `eventKey:"LogicalState"`
}
