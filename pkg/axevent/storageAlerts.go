package axevent

func StorageAlertEventKvs(wear *int, overall_health *int, temperature *int, alert *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("Storage", "Alert", nil, nil, []*KeyValueEntrie{
		{Key: "wear", Value: wear, ValueType: AXValueTypeInt},
		{Key: "overall_health", Value: overall_health, ValueType: AXValueTypeInt},
		{Key: "temperature", Value: temperature, ValueType: AXValueTypeInt},
		{Key: "alert", Value: alert, ValueType: AXValueTypeBool},
	})
}

type StorageAlertEvent struct {
	Wear          int  `eventKey:"wear"`
	OverallHealth int  `eventKey:"overall_health"`
	temperature   int  `eventKey:"temperature"`
	Tampered      bool `eventKey:"alert"`
}

func StorageDisruptionEventKvs(disk_id *string, disruption *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("Storage", "Disruption", nil, nil, []*KeyValueEntrie{
		{Key: "disk_id", Value: disk_id, ValueType: AXValueTypeString},
		{Key: "disruption", Value: disruption, ValueType: AXValueTypeBool},
	})
}

type StorageDisruptionEvent struct {
	DiskId     string `eventKey:"disk_id"`
	Disruption bool   `eventKey:"disruption"`
}

func StorageRecordingEventKvs(recording *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("Storage", "Recording", nil, nil, []*KeyValueEntrie{
		{Key: "recording", Value: recording, ValueType: AXValueTypeBool},
	})
}

type StorageRecordingEvent struct {
	Recording bool `eventKey:"recording"`
}
