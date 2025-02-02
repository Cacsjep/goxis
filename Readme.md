# About

Goxis is a library designed to facilitate the development of ACAP (Axis Camera Application Platform) applications for AXIS cameras.

This library includes user-friendly wrappers for most APIs provided by the [native SDK](https://axiscommunications.github.io/acap-documentation/docs/api/native-sdk-api.html).

[![Discord](https://img.shields.io/badge/Discord-Join%20us-blue?style=for-the-badge&logo=discord)](https://discord.gg/jrE98E6Qe9)

[![GoDocs](https://img.shields.io/badge/go-documentation-blue)](https://pkg.go.dev/github.com/Cacsjep/goxis)



## C-API Wrappers

Goxis organizes its functionality around several key areas of the native SDK, each encapsulated within its own package for ease of use:

[![GoDocs](https://img.shields.io/badge/go%20pkg-documentation-purple)](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg)

- [axevent](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axevent) - Event API for managing camera events.
- [axoverlay](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axoverlay) - Overlay API, including support for Cairo, for drawing over video feeds.
- [axparameter](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axparameter) - Parameter API for managing camera parameters.
- [axstorage](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axstorage) - Edge Storage API for accessing and managing on-camera storage.
- [axvdo](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axvdo) - Video Capture API for handling video streams.
- [axlarod](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axlarod) - Larod API
- [axlicense](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axlicense) - License Key API for managing application licenses.
- [axsyslog](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axsyslog) - Syslog for logging and diagnostics.

## Additional Packages

Beyond the core API wrappers, Goxis also provides several additional packages designed to speed up the development process and enhance application capabilities:

- [acapapp](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/acapapp) - Offers a high-level abstraction for quick and efficient ACAP application development.
- [axmanifest](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/axmanifest) - Aids in loading and parsing manifest files, simplifying application configuration and setup.
- [dbus](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/dbus) - Provides helpers for interacting with the D-Bus interface, including retrieving VAPIX credentials.
- [vapix](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/vapix) - Facilitates the use of the VAPIX API for interacting with camera functionalities.
- [glib](https://pkg.go.dev/github.com/Cacsjep/goxis/pkg/glib) - Includes helpers for working with GLib, such as managing the main event loop.

# Prerequisites
- Docker for building the ACAP applications
- [goxisbuilder](https://github.com/Cacsjep/goxisbuilder)

## Quickstart
```
go install github.com/Cacsjep/goxisbuilder@latest
goxisbuilder.exe -newapp
```

## Examples
Located at this repo [Examples](https://github.com/Cacsjep/goxis_examples)

## Module Installation
```
go get -u github.com/Cacsjep/goxis@latest
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
- **OverlayProvide:** Facilitates easy interaction with `axoverlay` related operations.
- **StorageProvider:** Offers straightforward access to the camera's storage, enhancing data management capabilities.
- `app.IsLicenseValid(major_version int, minor_version int)`: Verifies the validity of the application's license for the specified version.
- `app.Run()`: Activates the GMain loop within the application, allowing for continuous operation.
- `app.GetSnapshot(video_channel int)`: Captures and retrieves a JPEG snapshot from a given video channel.

and more ....

### Create a new goxis application

Just use [goxisbuilder](https://github.com/Cacsjep/goxisbuilder) tool for streamlined builds.

```
.\goxisbuilder.exe -newapp
```

## Events

Most events are already declared in `axevent`. If you miss something, you can manually craft it, create a PR, or just ask! 😊
You can use the AXIS `get_eventlist.py` script (Native SDK Repo - AXIS) to get the XML list to see events that your device supports.

> [!TIP]
> Use AXIS Meta Data Monitor to validate and determine if the desired event is triggered.

#### Namespace Considerations

It is crucial to set the correct namespace for each topic. Incorrect namespaces can cause issues in `axevent` C API mapping. Check carefully if multiple namespaces are involved, and ensure the correct setup.

##### Single Namespace Example: `tnsaxis:CameraApplicationPlatform/ObjectAnalytics/xinternal_data`

When there is only one namespace for all entries, ensure consistency across all topics. In this example, all entries use `OnfivNameSpaceTnsAxis`:
  ```go
    NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "CameraApplicationPlatform"),
    NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ObjectAnalytics"),
    NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "xinternal_data"),
```

##### Multiple Namespace Example: `tns1:Device/tnsaxis:IO/VirtualPort`

When there are multiple namespaces, ensure subsequent entries after a path change use the correct namespace. For instance, `VirtualPort` uses the `tnsaxis` namespace:

 ```go
    NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTns1, "Device"),
    NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "IO"),
    NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "VirtualPort"),
```

##### When creating new events, follow these patterns for namespace and topic consistency
```go
// <tnsaxis:CameraApplicationPlatform>
// 	<ObjectAnalytics>
// 		<xinternal_data wstop:topic="true">
// 			<tt:MessageDescription IsProperty="false">
// 				<tt:Data>
// 					<tt:SimpleItemDescription Name="svgframe" Type="xsd:string"></tt:SimpleItemDescription>
// 				</tt:Data>
// 			</tt:MessageDescription>
// 		</xinternal_data>
// 	</ObjectAnalytics>
// </tnsaxis:CameraApplicationPlatform>
func CameraApplicationPlatformXInternalDataEventKvs(svgFrame *string) *AXEventKeyValueSet {
	return NewAXEventKeyValueSetFromEntries([]KeyValueEntrie{
		NewTopicKeyValueEntrie("topic0", &OnfivNameSpaceTnsAxis, "CameraApplicationPlatform"),
		NewTopicKeyValueEntrie("topic1", &OnfivNameSpaceTnsAxis, "ObjectAnalytics"),
		NewTopicKeyValueEntrie("topic2", &OnfivNameSpaceTnsAxis, "xinternal_data"),
		NewStringKeyValueEntrie("svgframe", svgFrame),
	})
}

type CameraApplicationPlatformXInternalDataEvent struct {
	SvgFrame string `eventKey:"svgframe"`
}
```
