package main

import (
	"github.com/Cacsjep/goxis"
)

var (
	err     error
	isValid bool
	app     *goxis.AcapApplication
)

func main() {
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
