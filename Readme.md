# About

Goxis is a comprehensive library designed to facilitate the development of ACAP (Axis Camera Application Platform) applications for AXIS cameras. By providing a rich set of features and abstractions, Goxis aims to simplify the integration of camera functionalities into your applications.

This library includes user-friendly wrappers for most APIs provided by the [native SDK](https://axiscommunications.github.io/acap-documentation/docs/api/native-sdk-api.html).

## C-API Wrappers

Goxis organizes its functionality around several key areas of the native SDK, each encapsulated within its own package for ease of use:

- `axevent` - Event API for managing camera events.
- `axoverlay` - Overlay API, including support for Cairo, for drawing over video feeds.
- `axparameter` - Parameter API for managing camera parameters.
- `axstorage` - Edge Storage API for accessing and managing on-camera storage.
- `axvdo` - Video Capture API for handling video streams.
- `axlicense` - License Key API for managing application licenses.
- `axsyslog` - Syslog for logging and diagnostics.

## Additional Packages

Beyond the core API wrappers, Goxis also provides several additional packages designed to speed up the development process and enhance application capabilities:

- `acapapp` - Offers a high-level abstraction for quick and efficient ACAP application development.
- `dbus` - Provides helpers for interacting with the D-Bus interface, including retrieving VAPIX credentials.
- `glib` - Includes helpers for working with GLib, such as managing the main event loop.
- `vapix` - Facilitates the use of the VAPIX API for interacting with camera functionalities.
- `axmanifest` - Aids in loading and parsing manifest files, simplifying application configuration and setup.

# Prerequisites
- Docker for building the ACAP applications
- [goxisbuilder](https://github.com/Cacsjep/goxisbuilder)

## Install
```
go get github.com/Cacsjep/goxis
```

### Whats the purpose of goxis AcapApplication?
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

### Create a new goxis AcapApplication

Creating a new ACAP application with Goxis follows a similar project structure to the standard approach recommended in the [AXIS Documentation](https://axiscommunications.github.io/acap-documentation/docs/develop/application-project-structure.html). 
However, when leveraging the [goxisbuilder](https://github.com/Cacsjep/goxisbuilder) tool for streamlined builds, your application needs to be organized within a specific directory structure.

Each example follow this rule, just check out the examples.

Here's how you can set up your project for success:

- Start by creating a new directory for your project, e.g., `myawesomeacap`. This directory will serve as the container for your application's components.
  - Inside your project directory, create a Go source file (`*.go`). Any copy the code below into it.
  - Add a `LICENSE` file to clearly state the licensing terms under which your application is distributed.
  - Include a `manifest.json` that should be correctly configured. [AXIS Documentation](https://axiscommunications.github.io/acap-documentation/docs/develop/application-project-structure.html#create-a-manifest-file-from-scratch)

```go
package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
)

func main() {
	app := acapapp.NewAcapApplication()
	app.Syslog.Info("Hello from My awesome acap")
}
```

Build it with [Goxisbuilder](https://github.com/Cacsjep/goxisbuilder): 
```
.\goxisbuilder.exe -appdir="myawesomeacap"
```

Checkout the examples to see more about *AcapApplication*.

## Examples

Examples are really close to existing C examples of the [AXIS Native SDK repo](https://github.com/AxisCommunications/acap-native-sdk-examples).

| Example         | Description |
|-----------------|--------------|
| `axevent`	      | Demonstrate how to subscribe to an Virutal Input state change         |
| `axoverlay`	| Render rects via axolveray api                                        |
| `axparameter`   | Demonstrate how to get an parameter and listen to changes             |
| `axstorage`     | Interact with axstorage api                                           |
| `license` 	| Show how to obtain the license state                                  |
| `vdostream` 	| High level wrapper demonstration to get video frames (stream)         |
| `webserver`     | Reverse proxy webserver with fiber                                    |


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
  - Rewrite test package 
  - Adding Larod API