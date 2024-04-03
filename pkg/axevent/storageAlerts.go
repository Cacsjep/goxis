package axevent

func StorageAlertEventKvs(wear *int, overall_health *int, temperature *int, alert *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("Storage", "Alert", nil, nil, []*KeyValueEntrie{
		{key: "wear", value: wear, value_type: AXValueTypeInt},
		{key: "overall_health", value: overall_health, value_type: AXValueTypeInt},
		{key: "temperature", value: temperature, value_type: AXValueTypeInt},
		{key: "alert", value: alert, value_type: AXValueTypeBool},
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
		{key: "disk_id", value: disk_id, value_type: AXValueTypeString},
		{key: "disruption", value: disruption, value_type: AXValueTypeBool},
	})
}

type StorageDisruptionEvent struct {
	DiskId     string `eventKey:"disk_id"`
	Disruption bool   `eventKey:"disruption"`
}

func StorageRecordingEventKvs(recording *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("Storage", "Recording", nil, nil, []*KeyValueEntrie{
		{key: "recording", value: recording, value_type: AXValueTypeBool},
	})
}

type StorageRecordingEvent struct {
	Recording bool `eventKey:"recording"`
}
