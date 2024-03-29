# About

Goxis provides golang bindings for AXIS ACAP API's.
The acap package contains low level wrappers around ACAP API's.
Main package contains more high level interface to create ACAP Applications.

> [!NOTE]  
> This is a hobby project and is still in progress !

# Prerequisites

- Docker for building the ACAP applications

## Install
```
go get github.com/Cacsjep/goxis
```

## Goxis AcapApplication
AcapApplication is a high level based abstraction for an ACAP application,
it loads at runtime the ***manifest.json***. 

### Whats the purpose of AcapApplication ?
The `AcapApplication` acts as a foundational abstraction layer designed to streamline the handling of common tasks such as syslog logging, managing the GObject main loop, and the manipulation of parameters and events.

Upon instantiation, `AcapApplication` undertakes several crucial steps::
- **Manifest Parsing:** It reads the manifest file to extract the application's name.
- **Parameter Management:** Initializes an `axparameter` instance, enabling the application to get, set, and remove parameters efficiently.
- **Event Handling:** Sets up an `axevent` handler to facilitate event processing.
- **GMain Loop Preparation:** Configures a GMain loop complete with signal handlers, ensuring robust event management.

`AcapApplication` offers access to a variety of powerful functionalities, such as::
- **FrameProvider:** Facilitates easy interaction with `axvdo`, streamlining video-related operations.
- **StorageProvider:** Offers straightforward access to the camera's storage, enhancing data management capabilities.
- `app.IsLicenseValid(major_version int, minor_version int)`: Verifies the validity of the application's license for the specified version.
- `app.Run()`: Activates the GMain loop within the application, allowing for continuous operation.
- `app.GetSnapshot(video_channel int)`: Captures and retrieves a JPEG snapshot from a given video channel.

Also FrameProvider can use from app to easy access the axvdo,
or using StorageProvider to easy access camera storage.

### Create a new AcapApplication
- Create a new directory like *myawesomeacap*
- Create a .go file

```go
package main

import (
	"github.com/Cacsjep/goxis"
)

func main() {
	app := goxis.NewAcapApplication()
	app.Syslog.Info("Hello from My awesome acap")
}
```

Build it with [Goxisbuilder](#goxisbuilder): 
```
.\goxisbuilder.exe -appdir="myawesomeacap"
```


## Goxisbuilder
Goxisbuilder is a command line tool that can build docker acap applications.
```
go install github.com/Cacsjep/goxis/cmd/goxisbuilder@latest
```

### Application Structure
You can look into any example, a application needs:
 - .go file with a main function
 - LICENSE 
 - manifest.json

### Usage
```
.\goxisbuilder.exe -appdir="example/license"
```

> [!NOTE]  
> appdir is a directory with a main.go and needs to have a ***LICENSE*** and a ***manifest.json***

| Flag            | Description                                                                                           | Default          |
|-----------------|-------------------------------------------------------------------------------------------------------|------------------|
| `-h`            | Displays this help message.                                                                           |                  |
| `-appdir`       | The full path to the application directory from which to build.                                       |  `""`            |
| `-arch`         | The architecture for the ACAP application: 'aarch64' or 'armv7hf'.                                    |  `"aarch64"`     |
| `-build-examples`| Set to true to build example applications.                                                           |  `false`         |
| `-install`      | Set to true to install the application on the camera.                                                 |  `false`         |
| `-libav`        | Set to true to compile libav for binding with go-astiav.                                              |  `false`         |
| `-lowsdk`       | Set to true to build for firmware versions greater than 10.9 with SDK version 1.1. This adjusts the manifest to use version 1.3.|                  |
| `-manifest`     | The path to the manifest file. Defaults to 'manifest.json'.                                           |  `"manifest.json"`|
| `-ip`           | The IP address of the camera where the EAP application is installed.                                  |  `""`            |
| `-pwd`          | The root password for the camera where the EAP application is installed.                              |  `""`            |
| `-start`        | Set to true to start the application after installation.                                              |  `false`         |
| `-watch`        | Set to true to monitor the package log after building.                                                |  `false`         |

## Examples
Just digg into examples to see how you can use goxis.
Currently we had:
  - axevent 	| Demonstrate how to subscribe to an Virutal Input state change
  - axoverlay	| Render rects via axolveray api
  - axparameter   | Demonstrate how to get an parameter and listen to changes
  - license 	| Show how to obtain the license state
  - vdostream 	| High level wrapper demonstration to get video frames (stream)
  - webserver     | Reverse proxy webserver with fiber


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


# Todos:
  - Adding Larod API