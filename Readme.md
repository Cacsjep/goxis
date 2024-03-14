
# Nativ SDK API's

## Axparameter

https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/#parameter-api

Its possible to get any parameter (Camera, ACAP) via this API, but you can only add or set parameters,
for your ACAP App.

#### Get Parameter
```
if mac, err := app.ParamGet("Properties.System.SerialNumber"); err != nil {
    app.Syslog.Error(err.Error())
} else {
    app.Syslog.Info(fmt.Sprintf("Mac-Address: %s", mac))
}
```

#### Add Parameter
```
if err := app.ParamAdd("IsCustomized", "yes", "bool:no,yes"); err != nil {
    app.Syslog.Error(err.Error())
}
```
#### Parameter Change callback
```
myparam := "IsCustomized"
if err := app.ParamAddCallback(myparam, func(name, value string, app *goxis.AcapApplication) {
	app.Syslog.Info(fmt.Sprintf("Callback invoked for param: %s, new-value: %s", name, value))
}); err != nil {
	app.Syslog.Error(err.Error())
} else {
	app.Syslog.Info(fmt.Sprintf("Callback registerd for param: %s", myparam))
}
```


# Event Channels
To enable retrival of event via an channel enable them via the AppConfiguration,
for example to get IO events enable RegisterIOChannel.

### IO Example
```
package main

import (
	"fmt"

	"github.com/Cacsjep/goxis"
)

func main() {
	app, err := goxis.NewApp(&goxis.AppConfiguration{RegisterIOChannel: true})
	if err != nil {
		panic(err)
	}
	defer app.Shutdown()

	if err = app.Load(); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case value, ok := <-app.IoEventChannel:
				if !ok {
					fmt.Println("IoEventChannel closed. Exiting goroutine.")
					return
				}
				fmt.Printf("IO Event | %s:%d=%t \n", value.PortTypeName, value.Port, value.State)
			}
		}
	}()

	app.Run()
}
```
# CGI API's

## Parameters
https://www.axis.com/vapix-library/subjects/t10175981/section/t10036014/display

#### Get a map of all existing parameters
```
paramMap, err := app.VapixParamCgiGetAll()
if err != nil {
    app.Syslog.Error(err.Error())
} else {
    app.Syslog.Info(fmt.Sprintf("I0.Appearance.Resolution: %s", paramMap["root.Image.I0.Appearance.Resolution"]))
}
```
#### Update multiple params
```
toUpdate := []*goxis.Param{
    {Key: "root.Audio.A0.Name", Value: "acs"},
}
err = app.VapixParamCgiUpdate(toUpdate)
if err != nil {
    app.Syslog.Error(err.Error())
} else {
    app.Syslog.Info("root.Audio.A0.Name update successfully")
}
```


# VAPIX API's

## Virtual Input
https://www.axis.com/vapix-library/subjects/t10175981/section/t10074527/display?section=t10074527-t10036012

#### Activate
```
if state_changed, err := app.VapixVirtualInputChange(goxis.VIO_Activate, 2); err != nil {
    app.Syslog.Error(err.Error())
} else {
    app.Syslog.Info(fmt.Sprintf("VIO 2 ON State Changed: %t", state_changed))
}
```
    

#### Deactivate
```
if state_changed, err := app.VapixVirtualInputChange(goxis.VIO_Deactivate, 2); err != nil {
    app.Syslog.Error(err.Error())
} else {
    app.Syslog.Info(fmt.Sprintf("VIO 2 OFF State Changed: %t", state_changed))
}
```
    

# Webserver with Reverse Proxy
https://github.com/AxisCommunications/acap-native-sdk-examples/tree/main/web-server

To use a webserver like fiber we use the reverse proxy support for ACAP,
therefrore we need to declare reverseProxy item to Manifest.

Currently there is now way to declare a settingsPage directly to the reverseProxy path,
so we need a redirect.html that redirects to the correct path.
Thats needed because we want to server our own html files and dont want that Camera Apache 
server did this for us.

### Redirect HTML
- Create a folder called "html" in app directory and place the redirect.html in it.
We replace the current location href with the apiPath that is configured in the Manifest.
```
<!DOCTYPE html>
<html lang="en">
<script>
    window.onload = function() {
        window.location.href = window.location.href.replace("redirect.html", "goxis"); // apiPath
    }
</script>
</html>
```

### Mainfest
Adding a reverse proxy configuration like this example shown
```
"configuration": {
    "settingPage": "redirect.html",
    "reverseProxy": [
        {
            "apiPath": "goxis",
            "target": "http://localhost:2001",
            "access": "admin"
        }
    ]
}
```

### Fiber example
An Example for an fiber webserver with embeding static files and index.html

- Create a folder called "static" in app directory and place your static files in there.
- Create and index.html file in the app directory
```
package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/Cacsjep/goxis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed index.html
var f embed.FS

//go:embed static/*
var embedDirStatic embed.FS

func main() {
	app, err := goxis.NewApp(&goxis.AppConfiguration{})
	if err != nil {
		panic(err)
	}
	defer app.Shutdown()
	baseUrl := app.ReversProxyUrlFirstEntry()
	fapp := fiber.New()
	fapp.Use(baseUrl, filesystem.New(filesystem.Config{
		Root: http.FS(f),
	}))
	fapp.Use(baseUrl+"/static", filesystem.New(filesystem.Config{
		Root:       http.FS(embedDirStatic),
		PathPrefix: "static",
		Browse:     true,
	}))
	go fapp.Listen("127.0.0.1:2001")
	fmt.Println("Start Gmain loop")
	app.Run()
}
```

# ACAP API Docs
https://axiscommunications.github.io/acap-documentation/docs/acap-sdk-version-3/api/

# ACAP Native SDK hardware compatibility

**Last Modified: 08.03.2024**

| Chip          | Architecture |
|---------------|--------------|
| ARTPEC-6      | armv7hf      |
| ARTPEC-7      | armv7hf      |
| ARTPEC-8      | aarch64      |
| CV25          | aarch64      |
| i.MX 6SoloX   | armv7hf      |
| i.MX 6ULL     | armv7hf      |
| S5            | aarch64      |
| S5L           | aarch64      |

# Manifest schema version mapping

**Mapping table for schema, firmware and SDK version.**

It’s recommended to use the latest manifest version available for the minimum firmware version targeted.

| Schema | Firmware | SDK  | Description |
|--------|----------|------|-------------|
| 1.0    | 10.7     | 1.0  | Initial basic version |
| 1.1    | 10.7     | 1.0  | Additional fields, mainly for technical reasons |
| 1.2    | 10.7     | 1.0  | Enables uninstall functionality which is required by e.g. docker-compose-acap |
| 1.3    | 10.9     | 1.1  | Architecture will be automatically generated and added to manifest at packaging step |
| 1.3.1  | 11.0     | 1.4  | Bugfixes; Allow = in runOptions and maxLength of appName should be 26 |
| 1.4.0  | 11.7     | 1.11 | Allow new characters ( ) , . ! ? & ' for vendor field |
| 1.5.0  | 11.8     | 1.12 | - Add support for reverse proxy configuration.<br>- Add access policy for ACAP application web content.<br>- Allow - character in secondary groups of linux resources.<br>- Allow strings in requiredMethods and conditionalMethods under dbus to end with .* to match all methods of a D-Bus interface. |
| 1.6.0  | 11.9     | 1.13 | - Add support for characters $ and \ in apiPath of the reverse proxy configuration.<br>- Add optional field $schema that can point out a manifest schema to use for manifest validation and auto-completion.<br>- Allow strings in requiredMethods and conditionalMethods under dbus to contain -. |

# SDK for software compatibility

**Choose the appropriate SDK version based on what firmware version you want supporting your ACAP application.**

| SDK version | Compatible with firmware version |
|-------------|----------------------------------|
| SDK 3.0     | 9.70 and later                   |
| SDK 3.1     | 9.80 (LTS) and later             |
| SDK 3.2     | 10.2 and later                   |
| SDK 3.3     | 10.5 and later                   |
| SDK 3.4     | 10.6 and later                   |
| SDK 3.5     | 10.9 and later                   |

# ACAP Native SDK

**ACAP Release vs ACAP Native SDK Image version vs Compatible with AXIS OS version**

| ACAP Release | ACAP Native SDK Image version | Compatible with AXIS OS version |
|--------------|-------------------------------|---------------------------------|
| 4.0          | 1.0                           | 10.7 and later until LTS        |
| 4.1          | 1.1                           | 10.9 and later until LTS        |
| 4.2          | 1.2                           | 10.10 and later until LTS       |
| 4.3          | 1.3                           | 10.12 (LTS)                     |
| 4.4          | 1.4                           | 11.0 and later until LTS        |
| 4.5          | 1.5                           | 11.1 and later until LTS        |
| 4.6          | 1.6                           | 11.2 and later until LTS        |
| 4.7          | 1.7                           | 11.3 and later until LTS        |
| 4.8          | 1.8                           | 11.4 and later until LTS        |
| 4.9          | 1.9                           | 11.5 and later until LTS        |
| 4.10         | 1.10                          | 11.6 and later until LTS        |
| 4.11         | 1.11                          | 11.7 and later until LTS        |
| 4.12         | 1.12                          | 11.8 and later until LTS        |
| 4.13         | 1.13                          | 11.9 and later until LTS        |

