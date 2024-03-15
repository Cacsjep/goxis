// App is an Highlevel abstraction for an ACAP Application
// Automatic declarations and loadings
//     Manifest at runtime to get APPName and APP ID for example
//     Syslog
//	   AXParameter Instance
//	   AXEventHandler Instance
//	   GMainLoop Instance

package app

import (
	"strconv"

	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/axlicense"
	"github.com/Cacsjep/goxis/pkg/axparam"
	"github.com/Cacsjep/goxis/pkg/glib"
	"github.com/Cacsjep/goxis/pkg/manifest"
	"github.com/Cacsjep/goxis/pkg/syslog"
)

type AcapApplication struct {
	Manifest           *manifest.ApplicationManifestSchema
	Syslog             *syslog.Syslog
	AppParamHandler    *axparam.AXParameter
	GlobalParamHandler *axparam.AXParameter
	EventHandler       *axevent.AXEventHandler
	Mainloop           *glib.GMainLoop
}

func NewAcapApplication() (*AcapApplication, error) {
	m, err := manifest.LoadManifest("manifest.json")
	if err != nil {
		return nil, err
	}

	pApp, err := axparam.AXParameterNew(&m.ACAPPackageConf.Setup.AppName)
	if err != nil {
		return nil, err
	}

	pGlobal, err := axparam.AXParameterNew(nil)
	if err != nil {
		return nil, err
	}

	app := AcapApplication{
		Manifest:           m,
		Syslog:             syslog.NewSyslog(m.ACAPPackageConf.Setup.AppName, syslog.LOG_PID|syslog.LOG_CONS, syslog.LOG_USER),
		AppParamHandler:    pApp,
		GlobalParamHandler: pGlobal,
		EventHandler:       axevent.NewEventHandler(),
		Mainloop:           glib.NewMainLoop(),
	}

	return &app, nil
}

func (a *AcapApplication) IsLicenseValid(major_version int, minor_version int) (bool, error) {
	appId, err := strconv.Atoi(a.Manifest.ACAPPackageConf.Setup.AppID)
	if err != nil {
		return false, err
	}
	return axlicense.LicensekeyVerify(
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
	a.AppParamHandler.Free()
	a.GlobalParamHandler.Free()
	a.EventHandler.Free()
	a.Syslog.Close()
}
