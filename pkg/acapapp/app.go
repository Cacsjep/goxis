// Package acapapp provides a high-level abstraction for an Axis Communications Application Platform (ACAP) application.
// It encapsulates the application's manifest, system logging, parameter handling, event handling, and the main event loop
// to facilitate easy development of ACAP applications. This includes automatic loading of the application's manifest,
// initialization of syslog for logging, handling of application parameters, event handling, and the GMainLoop for the main event loop.
package acapapp

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/Cacsjep/goxis/pkg/axevent"
	"github.com/Cacsjep/goxis/pkg/axlicense"
	"github.com/Cacsjep/goxis/pkg/axmanifest"
	"github.com/Cacsjep/goxis/pkg/axparameter"
	"github.com/Cacsjep/goxis/pkg/axsyslog"
	"github.com/Cacsjep/goxis/pkg/axvdo"
	"github.com/Cacsjep/goxis/pkg/glib"
)

// AcapApplication provides a high-level abstraction for an Axis Communications Application Platform (ACAP) application.
// It encapsulates the application's manifest, system logging, parameter handling, event handling, and the main event loop
// to facilitate easy development of ACAP applications. This includes automatic loading of the application's manifest,
// initialization of syslog for logging, handling of application parameters, event handling, and the GMainLoop for the main event loop.
type AcapApplication struct {
	Manifest        *axmanifest.ApplicationManifestSchema
	Syslog          *axsyslog.Syslog
	ParamHandler    *axparameter.AXParameter
	EventHandler    *axevent.AXEventHandler
	Mainloop        *glib.GMainLoop
	OnCloseCleaners []func()
}

// NewAcapApplication initializes a new AcapApplication instance, loading the application's manifest,
// setting up the syslog, parameter handler, event handler, and main loop. It returns an initialized AcapApplication instance.
//
// ! Note: Since this is the entry point, it panic in case of an error,
// this could happen if manifest could not loaded or parameter instance could not be created
func NewAcapApplication() *AcapApplication {
	m, err := axmanifest.LoadManifest("manifest.json")
	if err != nil {
		panic(err)
	}

	pApp, err := axparameter.AXParameterNew(m.ACAPPackageConf.Setup.AppName)
	if err != nil {
		panic(err)
	}

	app := AcapApplication{
		Manifest:        m,
		Syslog:          axsyslog.NewSyslog(m.ACAPPackageConf.Setup.AppName, axsyslog.LOG_PID|axsyslog.LOG_CONS, axsyslog.LOG_USER),
		ParamHandler:    pApp,
		EventHandler:    axevent.NewEventHandler(),
		Mainloop:        glib.NewMainLoop(),
		OnCloseCleaners: []func(){},
	}

	showHelp := flag.Bool("h", false, "Displays this help message.")
	consoleLog := flag.Bool("consoleLog", false, "Enable console logging")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(1)
	}

	if *consoleLog {
		app.Syslog.EnableConsole()
	}

	return &app
}

// IsLicenseValid checks the validity of the application's license for a given major and minor version.
// It returns true if the license is valid, or false along with an error if the check fails.
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

// Start initiates the main event loop of the application, beginning its execution.
func (a *AcapApplication) Run() {
	SignalHandler(a.Close)
	a.Mainloop.Run()
}

// Add close or clean functions to app so in case of signals these are correct handled
func (a *AcapApplication) AddCloseCleanFunc(f func()) {
	a.OnCloseCleaners = append(a.OnCloseCleaners, f)
}

// Close terminates the application's main event loop and releases resources associated with the syslog, parameter handler,
// event handler, and main loop. This should be called to cleanly shut down the application.
func (a *AcapApplication) Close() {
	a.Syslog.Info("Stop Application")
	for _, f := range a.OnCloseCleaners {
		f()
	}
	a.Mainloop.Quit()     // Terminate the main loop.
	a.ParamHandler.Free() // Release the parameter handler.
	a.EventHandler.Free() // Release the event handler.
	a.Syslog.Close()      // Close the syslog.
}

// GetSnapshot captures a JPEG snapshot from the specified video channel and returns it as a byte slice.
// It sets up the required settings for capturing the snapshot, captures it, and then returns the snapshot data or an error if the capture fails.
func (a *AcapApplication) GetSnapshot(video_channel int) ([]byte, error) {
	settings := axvdo.NewVdoMap()                             // Create a new settings map for the snapshot.
	settings.SetUint32("channel", uint32(video_channel))      // Set the video channel.
	settings.SetUint32("format", uint32(axvdo.VdoFormatJPEG)) // Set the snapshot format to JPEG.
	defer settings.Unref()                                    // Ensure settings are unreferenced after use.

	snapshotBuffer, err := axvdo.Snapshot(settings) // Capture the snapshot.
	if err != nil {
		return nil, err
	}
	defer snapshotBuffer.Unref() // Ensure the snapshot buffer is unreferenced after use.

	return snapshotBuffer.GetBytes() // Return the snapshot data.
}

// AcapWebBaseUri returns the base path for an webserver that is used with reverse proxy
// reverse proxy uri for acap are:
// /local/<appname>/<apipath>
func (a *AcapApplication) AcapWebBaseUri() (string, error) {
	if len(a.Manifest.ACAPPackageConf.Configuration.ReverseProxy) == 0 {
		return "", errors.New("No reverse proxy configuration set in manifest")
	}
	pkgcfg := a.Manifest.ACAPPackageConf
	return fmt.Sprintf("/local/%s/%s", pkgcfg.Setup.AppName, pkgcfg.Configuration.ReverseProxy[0].ApiPath), nil
}
