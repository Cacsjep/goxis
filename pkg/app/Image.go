package app

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/axvdo"
)

// Returns a JPEG Snapshot with full resolution for give video channel
// Video channel, 0 is overview, 1, 2, ... are view areas.
func (a *AcapApplication) GetSnapsot(video_channel int) ([]byte, error) {
	settings := axvdo.NewVdoMap()
	settings.SetUint32("channel", uint32(video_channel))
	settings.SetUint32("format", uint32(axvdo.VdoFormatJPEG))
	defer settings.Unref()
	snapshotBuffer, err := axvdo.Snapshot(settings)
	if err != nil {
		return nil, fmt.Errorf("Vdo Error: %s, Vdo Error Code: %d, Vdo Error Excepted: %t", err.Err.Error(), err.Code, err.Expected)
	}
	defer snapshotBuffer.Unref()
	return snapshotBuffer.GetBytes()
}

// Returns a list of resolutions for the given video channel
// Video channel, 0 is overview, 1, 2, ... are view areas.
func (a *AcapApplication) GetVdoChannelResolutions(video_channel int) ([]axvdo.VdoResolution, error) {
	s, err := axvdo.VdoChannelGet(uint(video_channel))
	if err != nil {
		return nil, err
	}
	return s.GetResolutions(nil)
}

// Returns the higest resolution for a video channel
// Video channel, 0 is overview, 1, 2, ... are view areas.
func (a *AcapApplication) GetVdoChannelMaxResolution(video_channel int) (*axvdo.VdoResolution, error) {
	resolutions, err := a.GetVdoChannelResolutions(video_channel)
	if err != nil {
		return nil, err
	}

	var highest axvdo.VdoResolution
	maxPixels := 0

	for _, res := range resolutions {
		pixels := res.Width * res.Height
		if pixels > maxPixels {
			highest = res
			maxPixels = pixels
		}
	}

	return &highest, nil
}
