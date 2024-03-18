package acap

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

func parseCredentials(credentialsString string) (*VapixCreds, error) {
	parts := strings.SplitN(credentialsString, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("error parsing credential string '%s'", credentialsString)
	}
	return &VapixCreds{Username: parts[0], Password: parts[1]}, nil
}

func RetrieveVapixCredentials(username string) (*VapixCreds, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("error connecting to D-Bus: %s", err)
	}
	defer conn.Close()

	busName := "com.axis.HTTPConf1"
	objectPath := dbus.ObjectPath("/com/axis/HTTPConf1/VAPIXServiceAccounts1")
	interfaceName := "com.axis.HTTPConf1.VAPIXServiceAccounts1"
	methodName := "GetCredentials"

	var credentialsString string
	obj := conn.Object(busName, objectPath)
	call := obj.Call(interfaceName+"."+methodName, 0, username)
	if call.Err != nil {
		return nil, fmt.Errorf("error invoking D-Bus method: %s", call.Err)
	}

	call.Store(&credentialsString)
	credentials, err := parseCredentials(credentialsString)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}
