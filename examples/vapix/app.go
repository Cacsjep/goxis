package main

import (
	"io"

	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/dbus"
	"github.com/Cacsjep/goxis/pkg/vapix"
)

// Activate virtual input port 1 via vapix get method
func activateVirtualInputPort(username, password string, port int) ([]byte, error) {
	r := vapix.VapixGet(username, password, vapix.InternalVapixUrlPathJoin("/axis-cgi/virtualinput/activate.cgi?schemaversion=1&port=1"))
	if r.IsOk {
		defer r.ResponseReader.Close()
		if rawxml, err := io.ReadAll(r.ResponseReader); err != nil {
			return nil, err
		} else {
			return rawxml, nil
		}
	}
	return nil, r.Error
}

// List all parameters via vapix get method
func listAllParameters(username, password string) (map[string]string, error) {
	r := vapix.VapixGet(username, password, vapix.InternalVapixUrlPathJoin("/axis-cgi/param.cgi?action=list"))
	if r.IsOk {
		defer r.ResponseReader.Close()
		if params, err := vapix.ParseKeyValueRequestBody(r.ResponseReader); err != nil {
			return nil, err
		} else {
			return params, nil
		}
	}
	return nil, r.Error

}

func main() {
	app := acapapp.NewAcapApplication()

	// Retrieve VAPIX credentials
	username, password, err := dbus.RetrieveVapixCredentials("root")
	if err != nil {
		app.Syslog.Critf("Failed to retrieve VAPIX credentials: %s", err.Error())
	}

	// Activate virtual input port 1
	app.Syslog.Infof("User: %s, Password: %s", username, password)
	rawxml, err := activateVirtualInputPort(username, password, 1)
	if err != nil {
		app.Syslog.Critf("Failed to activate virtual input port: %s", err.Error())
	}
	app.Syslog.Infof("VirtualInputPort XML Response: %v", string(rawxml))

	// List all parameters
	params, err := listAllParameters(username, password)
	if err != nil {
		app.Syslog.Critf("Failed to list all parameters: %s", err.Error())
	}
	app.Syslog.Infof("Parameters: %v", params)
}
