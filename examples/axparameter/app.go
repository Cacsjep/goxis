package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
)

var (
	err        error
	serial_nbr string
	app        *goxis.AcapApplication
)

func main() {
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	// Parameters outside the application's group requires qualification.
	// This could also done via vapix and dbus acap.RetrieveVapixCredentials() and acap.VapixGet()
	if serial_nbr, err = app.ParamHandler.Get("Properties.System.SerialNumber"); err != nil {
		app.Syslog.Error(err.Error())
	} else {
		app.Syslog.Info(fmt.Sprintf("SerialNumber: %s", serial_nbr))
	}

	// Act on changes to IsCustomized as soon as they happen.
	err = app.ParamHandler.RegisterCallback("IsCustomized", func(name, value string, userdata any) {
		app.Syslog.Info(fmt.Sprintf("Param Callback Invoked, Parameter Name: %s, Value: %s, Userdata: %s", name, value, userdata.(string)))
	}, "myuserdata")

	// Signal handler automatically internally created for SIGTERM, SIGINT
	// This blocks now the main thread.
	app.Run()
}
