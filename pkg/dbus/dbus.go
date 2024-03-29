package dbus

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

func parseCredentials(credentialsString string) (username string, password string, err error) {
	parts := strings.SplitN(credentialsString, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("error parsing credential string '%s'", credentialsString)
	}
	return parts[0], parts[1], nil
}

func RetrieveVapixCredentials(user string) (username string, password string, err error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return "", "", fmt.Errorf("error connecting to D-Bus: %s", err)
	}
	defer conn.Close()

	busName := "com.axis.HTTPConf1"
	objectPath := dbus.ObjectPath("/com/axis/HTTPConf1/VAPIXServiceAccounts1")
	interfaceName := "com.axis.HTTPConf1.VAPIXServiceAccounts1"
	methodName := "GetCredentials"

	var credentialsString string
	obj := conn.Object(busName, objectPath)
	call := obj.Call(interfaceName+"."+methodName, 0, user)
	if call.Err != nil {
		return "", "", fmt.Errorf("error invoking D-Bus method: %s", call.Err)
	}
	call.Store(&credentialsString)
	username, password, err = parseCredentials(credentialsString)
	if err != nil {
		return "", "", err
	}
	return username, password, nil
}
