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
	app = goxis.NewAcapApplication()
	if isValid, err = app.IsLicenseValid(1, 0); err != nil {
		app.Syslog.Crit(err.Error())
	}
	if isValid {
		app.Syslog.Info("License is valid")
	} else {
		app.Syslog.Warn("Invalid License")
	}
}
