package main

import (
	"github.com/Cacsjep/goxis"
)

var (
	err     error
	isValid bool
	// High level acap wrapper
	app *goxis.AcapApplication
)

func main() {
	// Creates a new ACAP Application based on manifest
	// Internal it creates:
	//    Syslog instance
	//    AxParamter instance
	//    AxEventHandler instance
	//    Gmainloop instance
	// In this example we just need the IsLicenseValid from the AcapApplication
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	if isValid, err = app.IsLicenseValid(1, 0); err != nil {
		panic(err)
	}

	if isValid {
		app.Syslog.Info("License is valid")
	} else {
		app.Syslog.Warn("Invalid License")
	}
}
