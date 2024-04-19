package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axparameter"
)

// This example demonstrate how to use paramHandler to easy interact with parameters.
//
// Orginal C Example: https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/axparameter
func main() {

	// Initialize a new ACAP application instance.
	// AcapApplication initializes the ACAP application with there name, eventloop, and syslog etc..
	app := acapapp.NewAcapApplication()

	if serial_nbr, err := app.ParamHandler.Get("Properties.System.SerialNumber"); err != nil {
		app.Syslog.Error(err.Error())
	} else {
		app.Syslog.Info(fmt.Sprintf("SerialNumber: %s", serial_nbr))
	}

	// Parameters "IsCustomized" is declared in the manifest.json
	// OnChange register a callback for the given parameter via ax_parameter_register_callback under the hood.
	// ! Note: As docs suggest: Callback functions should avoid blocking calls, i.e. never call any axparameter method as this will very likely lead to a deadlock,
	// * if you want to access params use a gorutine as shown below.
	if err := app.ParamHandler.OnChange("IsCustomized", func(e *axparameter.ParameterChangeEvent) {
		app.Syslog.Info(fmt.Sprintf("(OnChange) Param Changed, Parameter Name: %s, Value: %s", e.Name, e.Value))

		// * using a goroutine to access the parameter is the recommended way
		go app.ParamHandler.Get("Properties.System.SerialNumber")

	}); err != nil {
		app.Syslog.Error(err.Error())
	}

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axparameter needs a running event loop to handle the parameter callbacks corretly
	app.Run()
}
