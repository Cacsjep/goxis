package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
)

var (
	err        error
	serial_nbr string
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
	if app, err = goxis.NewAcapApplication(); err != nil {
		panic(err)
	}
	defer app.Close()

	// Use of ParamHandler
	// ParamHandler has all axparameter functions wrapped for golang.
	// Add, Get, Set, List or callback registering.

	// Parameters outside the application's group requires qualification.
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
	// This blocks not the main thread
	app.Mainloop.Run()

	app.Syslog.Info("Application was stopped")
}
