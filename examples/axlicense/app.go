package main

import "github.com/Cacsjep/goxis/pkg/acapapp"

// This example demonstrates how to check the license status of an ACAP application.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/licensekey
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	major_version := 0
	minor_version := 1
	isValid, err := app.IsLicenseValid(major_version, minor_version)
	if err != nil {
		app.Syslog.Crit(err.Error())
	}

	if isValid {
		app.Syslog.Info("License is valid")
	} else {
		app.Syslog.Warn("Invalid License")
	}
}
