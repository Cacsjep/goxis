package app

import "github.com/Cacsjep/goxis/pkg/axvdo"

func (a *AcapApplication) NewVideoStream(stream_cfg VideoSteamConfiguration) (*axvdo.VdoStream, error) {
	vdoMap := VideoStreamConfigToVdoMap(stream_cfg)
	defer vdoMap.Unref()
	return axvdo.NewStream(vdoMap)
}
