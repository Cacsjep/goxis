package main

import "github.com/Cacsjep/goxis/pkg/acapapp"

func main() {
	app := acapapp.NewAcapApplication()
	if isValid, err := app.IsLicenseValid(1, 0); err != nil {
		app.Syslog.Crit(err.Error())
	} else {
		if isValid {
			app.Syslog.Info("License is valid")
		} else {
			app.Syslog.Warn("Invalid License")
		}
	}
}
