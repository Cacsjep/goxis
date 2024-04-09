package axevent

import "github.com/Cacsjep/goxis/pkg/utils"

func SupervisedPortEventKvs(port *int, tampered *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "IO", utils.NewStringPointer("SupervisedPort"), nil, []*KeyValueEntrie{
		{Key: "port", Value: port, ValueType: AXValueTypeInt},
		{Key: "tampered", Value: tampered, ValueType: AXValueTypeBool},
	})
}

type SupervisedPortEvent struct {
	Port     int  `eventKey:"port"`
	Tampered bool `eventKey:"tampered"`
}

func VirtualInputEventKvs(port *int, active *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "IO", utils.NewStringPointer("VirtualInput"), nil, []*KeyValueEntrie{
		{Key: "port", Value: port, ValueType: AXValueTypeInt},
		{Key: "active", Value: active, ValueType: AXValueTypeBool},
	})
}

type VirtualInputEvent struct {
	Port   int  `eventKey:"port"`
	Active bool `eventKey:"active"`
}

func StorageFailureEventKvs(disk_id *string, disruption *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "HardwareFailure", utils.NewStringPointer("StorageFailure"), nil, []*KeyValueEntrie{
		{Key: "disk_id", Value: disk_id, ValueType: AXValueTypeString},
		{Key: "disruption", Value: disruption, ValueType: AXValueTypeBool},
	})
}

type StorageFailureEvent struct {
	DiskId     string `eventKey:"disk_id"`
	Disruption bool   `eventKey:"disruption"`
}

func HeaterStatusEventKvs(heater *int, running *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Heater", utils.NewStringPointer("Status"), nil, []*KeyValueEntrie{
		{Key: "heater", Value: heater, ValueType: AXValueTypeInt},
		{Key: "running", Value: running, ValueType: AXValueTypeBool},
	})
}

type HeaterStatusEvent struct {
	Heater  int  `eventKey:"heater"`
	Running bool `eventKey:"running"`
}

func SystemReadyStatusEventKvs(ready *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Status", utils.NewStringPointer("SystemReady"), nil, []*KeyValueEntrie{
		{Key: "ready", Value: ready, ValueType: AXValueTypeBool},
	})
}

type SystemReadyStatusEvent struct {
	Ready int `eventKey:"ready"`
}

func TriggerRelayEventKvs(relayToken *int, logicalState *string) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Trigger", utils.NewStringPointer("Relay"), nil, []*KeyValueEntrie{
		{Key: "RelayToken", Value: relayToken, ValueType: AXValueTypeInt},
		{Key: "LogicalState", Value: logicalState, ValueType: AXValueTypeString},
	})
}

type TriggerRelayEvent struct {
	RelayToken   int    `eventKey:"RelayToken"`
	LogicalState string `eventKey:"LogicalState"`
}

func DigitalInputEventKvs(inputToken *int, logicalState *bool) (*AXEventKeyValueSet, error) {
	return NewTns1AxisEvent("Device", "Trigger", utils.NewStringPointer("DigitalInput"), nil, []*KeyValueEntrie{
		{Key: "InputToken", Value: inputToken, ValueType: AXValueTypeInt},
		{Key: "LogicalState", Value: logicalState, ValueType: AXValueTypeBool},
	})
}

type DigitalInputEvent struct {
	InputToken   int  `eventKey:"InputToken"`
	LogicalState bool `eventKey:"LogicalState"`
}
