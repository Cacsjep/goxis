package main

import (
	"fmt"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axparameter"
)

var (
	err        error
	serial_nbr string
	app        *acapapp.AcapApplication
)

func main() {
	app = acapapp.NewAcapApplication()

	if serial_nbr, err = app.ParamHandler.Get("Properties.System.SerialNumber"); err != nil {
		app.Syslog.Error(err.Error())
	} else {
		app.Syslog.Info(fmt.Sprintf("SerialNumber: %s", serial_nbr))
	}

	// Parameters "IsCustomized" is declared in the manifest.json
	// OnChange register a callback for the given parameter via ax_parameter_register_callback
	// ! Note: As docs suggest: Callback functions should avoid blocking calls, i.e. never call any axparameter method as this will very likely lead to a deadlock.
	// If you want to access params use a gorutine as shown below.
	if err = app.ParamHandler.OnChange("IsCustomized", func(e *axparameter.ParameterChangeEvent) {
		app.Syslog.Info(fmt.Sprintf("(OnChange) Param Changed, Parameter Name: %s, Value: %s", e.Name, e.Value))
		go app.ParamHandler.Get("Properties.System.SerialNumber")
	}); err != nil {
		app.Syslog.Error(err.Error())
	}

	// OnAnyChange register a callback for any parameter change via ax_parameter_register_callback
	// ! Note: Callbacks a registerd via there name so the onChange and OnAnyChange will conflict, use either OnChange or OnAnyChange
	/* if err = app.ParamHandler.OnAnyChange(func(e *axparameter.ParameterChangeEvent) {
		app.Syslog.Info(fmt.Sprintf("(OnAnyChange) Param Changed, Parameter Name: %s, Value: %s", e.Name, e.Value))
	}); err != nil {
		app.Syslog.Error(err.Error())
	} */

	app.Run()
}
