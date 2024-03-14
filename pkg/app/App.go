// App is an Highlevel abstraction for an ACAP Application
// Automatic declarations and loadings
//     Manifest at runtime to get APPName and APP ID for example
//     Syslog
//	   AXParameter Instance
//	   AXEventHandler Instance
//	   GMainLoop Instance

package app

import (
	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/axlicense"
	"github.com/Cacsjep/goxis/pkg/axparam"
	"github.com/Cacsjep/goxis/pkg/glib"
	"github.com/Cacsjep/goxis/pkg/manifest"
	"github.com/Cacsjep/goxis/pkg/syslog"
)

type AcapApplication struct {
	Manifest     *manifest.ApplicationManifestSchema
	Syslog       *syslog.Syslog
	ParamHandler *axparam.AXParameter
	EventHandler *axevent.AXEventHandler
	Mainloop     *glib.GMainLoop
}

func NewAcapApplication() (*AcapApplication, error) {
	m, err := manifest.LoadManifest("manifest.json")
	if err != nil {
		return nil, err
	}

	p, err := axparam.AXParameterNew(m.ACAPPackageConf.Setup.AppName)
	if err != nil {
		return nil, err
	}

	app := AcapApplication{
		Manifest:     m,
		Syslog:       syslog.NewSyslog(m.ACAPPackageConf.Setup.AppName, syslog.LOG_PID|syslog.LOG_CONS, syslog.LOG_USER),
		ParamHandler: p,
		EventHandler: axevent.NewEventHandler(),
		Mainloop:     glib.NewMainLoop(),
	}

	return &app, nil
}

func (a *AcapApplication) IsLicenseValid(major_version int, minor_version int) bool {
	return axlicense.LicensekeyVerify(
		a.Manifest.ACAPPackageConf.Setup.AppName,
		a.Manifest.ACAPPackageConf.Setup.AppID,
		major_version,
		minor_version,
	)
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
