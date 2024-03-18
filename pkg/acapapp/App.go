// App is an Highlevel abstraction for an ACAP Application
// Automatic declarations and loadings
//     Manifest at runtime to get APPName and APP ID for example
//     Syslog
//	   AXParameter Instance
//	   AXEventHandler Instance
//	   GMainLoop Instance

package acapapp

import (
	"strconv"

	"github.com/Cacsjep/goxis/pkg/acap"
	"github.com/Cacsjep/goxis/pkg/manifest"
)

type AcapApplication struct {
	Manifest     *manifest.ApplicationManifestSchema
	Syslog       *acap.Syslog
	ParamHandler *acap.AXParameter
	EventHandler *acap.AXEventHandler
	Mainloop     *acap.GMainLoop
}

func NewAcapApplication() (*AcapApplication, error) {
	m, err := manifest.LoadManifest("manifest.json")
	if err != nil {
		return nil, err
	}

	pApp, err := acap.AXParameterNew(&m.ACAPPackageConf.Setup.AppName)
	if err != nil {
		return nil, err
	}

	app := AcapApplication{
		Manifest:     m,
		Syslog:       acap.NewSyslog(m.ACAPPackageConf.Setup.AppName, acap.LOG_PID|acap.LOG_CONS, acap.LOG_USER),
		ParamHandler: pApp,
		EventHandler: acap.NewEventHandler(),
		Mainloop:     acap.NewMainLoop(),
	}

	return &app, nil
}

func (a *AcapApplication) IsLicenseValid(major_version int, minor_version int) (bool, error) {
	appId, err := strconv.Atoi(a.Manifest.ACAPPackageConf.Setup.AppID)
	if err != nil {
		return false, err
	}
	return acap.LicensekeyVerify(
		a.Manifest.ACAPPackageConf.Setup.AppName,
		appId,
		major_version,
		minor_version,
	), nil
}

func (a *AcapApplication) Start() {
	a.Mainloop.Run()
}

func (a *AcapApplication) Stop() {
	a.Mainloop.Quit()
	a.ParamHandler.Free()
	a.EventHandler.Free()
	a.Syslog.Close()
}

// Get an jpeg snapshot trough vdo for given vdo channel
func (a *AcapApplication) GetSnapsot(video_channel int) ([]byte, error) {
	settings := acap.NewVdoMap()
	settings.SetUint32("channel", uint32(video_channel))
	settings.SetUint32("format", uint32(acap.VdoFormatJPEG))
	defer settings.Unref()
	snapshotBuffer, err := acap.Snapshot(settings)
	if err != nil {
		return nil, err
	}
	defer snapshotBuffer.Unref()
	return snapshotBuffer.GetBytes()
}
