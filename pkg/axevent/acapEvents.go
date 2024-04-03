package axevent

import "github.com/Cacsjep/goxis/pkg/utils"

func CameraApplicationPlatformVmdCamera1ProfileANYEventKvs(active *bool) (*AXEventKeyValueSet, error) {
	return NewTnsAxisEvent("CameraApplicationPlatform", "VMD", utils.NewStringPointer("Camera1ProfileANY"), nil, []*KeyValueEntrie{
		{key: "active", value: active, value_type: AXValueTypeBool},
	})
}

type CameraApplicationPlatformVmdCamera1ProfileANYEvent struct {
	Active bool `eventKey:"active"`
}
