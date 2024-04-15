// Package acapapp provides a high-level abstraction for an Axis Communications Application Platform (ACAP) application.
// It encapsulates the application's manifest, system logging, parameter handling, event handling, and the main event loop
// to facilitate easy development of ACAP applications. This includes automatic loading of the application's manifest,
// initialization of syslog for logging, handling of application parameters, event handling, and the GMainLoop for the main event loop.
package acapapp

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/Cacsjep/goxis/pkg/axevent"
	axlarod "github.com/Cacsjep/goxis/pkg/axlaord"
	"github.com/Cacsjep/goxis/pkg/axlicense"
	"github.com/Cacsjep/goxis/pkg/axmanifest"
	"github.com/Cacsjep/goxis/pkg/axparameter"
	"github.com/Cacsjep/goxis/pkg/axsyslog"
	"github.com/Cacsjep/goxis/pkg/axvdo"
	"github.com/Cacsjep/goxis/pkg/glib"
	"github.com/Cacsjep/goxis/pkg/utils"
)

// AcapApplication provides a high-level abstraction for an Axis Communications Application Platform (ACAP) application.
// It encapsulates the application's manifest, system logging, parameter handling, event handling, and the main event loop
// to facilitate easy development of ACAP applications. This includes automatic loading of the application's manifest,
// initialization of syslog for logging, handling of application parameters, event handling, and the GMainLoop for the main event loop.
type AcapApplication struct {
	Manifest            *axmanifest.ApplicationManifestSchema
	Syslog              *axsyslog.Syslog
	ParamHandler        *axparameter.AXParameter
	EventHandler        *axevent.AXEventHandler
	Mainloop            *glib.GMainLoop
	OnCloseCleaners     []func()
	eventDeclarationIds []int
	Larod               *axlarod.Larod
}

// NewAcapApplication initializes a new AcapApplication instance, loading the application's manifest,
// setting up the syslog, parameter handler, event handler, and main loop. It returns an initialized AcapApplication instance.
//
// ! Note: Since this is the entry point, it panic in case of an error,
// this could happen if manifest could not loaded or parameter instance could not be created
func NewAcapApplication() *AcapApplication {
	m, err := axmanifest.LoadManifest("manifest.json")
	if err != nil {
		panic(err)
	}

	pApp, err := axparameter.AXParameterNew(m.ACAPPackageConf.Setup.AppName)
	if err != nil {
		panic(err)
	}

	app := AcapApplication{
		Manifest:        m,
		Syslog:          axsyslog.NewSyslog(m.ACAPPackageConf.Setup.AppName, axsyslog.LOG_PID|axsyslog.LOG_CONS, axsyslog.LOG_USER),
		ParamHandler:    pApp,
		EventHandler:    axevent.NewEventHandler(),
		Mainloop:        glib.NewMainLoop(),
		OnCloseCleaners: []func(){},
	}

	showHelp := flag.Bool("h", false, "Displays this help message.")
	consoleLog := flag.Bool("consoleLog", false, "Enable console logging")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(1)
	}

	if *consoleLog {
		app.Syslog.EnableConsole()
	}

	return &app
}

func (a *AcapApplication) InitalizeLarod() error {
	a.Larod = axlarod.NewLarod()
	if err := a.Larod.Initalize(); err != nil {
		return err
	}
	return nil
}

// IsLicenseValid checks the validity of the application's license for a given major and minor version.
// It returns true if the license is valid, or false along with an error if the check fails.
func (a *AcapApplication) IsLicenseValid(major_version int, minor_version int) (bool, error) {
	appId, err := strconv.Atoi(a.Manifest.ACAPPackageConf.Setup.AppID)
	if err != nil {
		return false, err
	}
	return axlicense.LicensekeyVerify(
		a.Manifest.ACAPPackageConf.Setup.AppName,
		appId,
		major_version,
		minor_version,
	), nil
}

// Start initiates the main event loop of the application, beginning its execution.
func (a *AcapApplication) Run() {
	SignalHandler(a.Close)
	a.Mainloop.Run()
}

// Add close or clean functions to app so in case of signals these are correct handled
func (a *AcapApplication) AddCloseCleanFunc(f func()) {
	a.OnCloseCleaners = append(a.OnCloseCleaners, f)
}

// Close terminates the application's main event loop and releases resources associated with the syslog, parameter handler,
// event handler, and main loop. This should be called to cleanly shut down the application.
func (a *AcapApplication) Close() {
	a.Syslog.Info("Stop Application")
	for _, f := range a.OnCloseCleaners {
		f()
	}
	for _, declaration_id := range a.eventDeclarationIds {
		a.EventHandler.Undeclare(declaration_id)
	}
	if a.Larod != nil {
		a.Larod.Disconnect()
	}
	a.Mainloop.Quit()     // Terminate the main loop.
	a.ParamHandler.Free() // Release the parameter handler.
	a.EventHandler.Free() // Release the event handler.
	a.Syslog.Close()      // Close the syslog.
}

// GetSnapshot captures a JPEG snapshot from the specified video channel and returns it as a byte slice.
// It sets up the required settings for capturing the snapshot, captures it, and then returns the snapshot data or an error if the capture fails.
func (a *AcapApplication) GetSnapshot(video_channel int) ([]byte, error) {
	settings := axvdo.NewVdoMap()                             // Create a new settings map for the snapshot.
	settings.SetUint32("channel", uint32(video_channel))      // Set the video channel.
	settings.SetUint32("format", uint32(axvdo.VdoFormatJPEG)) // Set the snapshot format to JPEG.
	defer settings.Unref()                                    // Ensure settings are unreferenced after use.

	snapshotBuffer, err := axvdo.Snapshot(settings) // Capture the snapshot.
	if err != nil {
		return nil, err
	}
	defer snapshotBuffer.Unref() // Ensure the snapshot buffer is unreferenced after use.

	return snapshotBuffer.GetBytes() // Return the snapshot data.
}

// AcapWebBaseUri returns the base path for an webserver that is used with reverse proxy
// reverse proxy uri for acap are:
// /local/<appname>/<apipath>
func (a *AcapApplication) AcapWebBaseUri() (string, error) {
	if len(a.Manifest.ACAPPackageConf.Configuration.ReverseProxy) == 0 {
		return "", errors.New("No reverse proxy configuration set in manifest")
	}
	pkgcfg := a.Manifest.ACAPPackageConf
	return fmt.Sprintf("/local/%s/%s", pkgcfg.Setup.AppName, pkgcfg.Configuration.ReverseProxy[0].ApiPath), nil
}

// AddCameraPlatformEvent adds the event to the application.
// AcapApplication undaclare the added events on closing or exit signals.
func (a *AcapApplication) AddCameraPlatformEvent(cpe *CameraPlatformEvent) (int, error) {
	event, err := NewCameraApplicationPlatformEvent(
		a.Manifest.ACAPPackageConf.Setup,
		cpe.Name,
		cpe.NiceName,
		cpe.Entries,
	)
	if err != nil {
		return 0, err
	}

	declarationID, err := a.EventHandler.Declare(event, cpe.Stateless, func(subscription int, userdata any) {}, nil)
	if err != nil {
		return 0, err
	}
	a.eventDeclarationIds = append(a.eventDeclarationIds, declarationID)
	return declarationID, nil
}

// SendPlatformEvent sends a platform event with the specified event ID and event creation function.
func (a *AcapApplication) SendPlatformEvent(eventID int, createEventFunc func() (*axevent.AXEvent, error)) error {
	event, err := createEventFunc()
	if err != nil {
		return err
	}
	return a.EventHandler.SendEvent(eventID, event)
}

// CameraPlatformEvent represents an event declaration for the Camera Application Platform.
type CameraPlatformEvent struct {
	// Name is a unique identifier for the event.
	Name string
	// NiceName is an optional human-readable name for the event.
	NiceName *string
	// KeyValueExtendedEntries is a slice of EventEntry structures, each representing a key-value pair with additional metadata.
	Entries []*EventEntry
	// Stateless indicates whether the event is stateless or not.
	// A stateless event, as defined in the AXEvent library, does not maintain any persistent state
	// and exists only at the moment it is sent. Such events are akin to signals or pulses, representing
	// instantaneous occurrences without an ongoing state. Examples of stateless events include momentary
	// interactions or occurrences, such as "a point crossed a line".
	//
	// In the context of event declarations, stateless events require the specification of all constituent
	// keys within the event. For keys with fixed values, these values must be explicitly provided. Conversely,
	// for keys that can assume variable values, their values should be set to nil in the declaration to
	// indicate their dynamic nature. This distinction is crucial for ensuring the correct interpretation
	// and handling of the event data.
	//
	// Setting this field to true characterizes the event as stateless, implying that it is transient and
	// does not relate to any persistent state. This is contrasted with stateful events, where the event's
	// state persists over time and may represent ongoing conditions or statuses, such as the active or
	// inactive state of an I/O port. In the declaration of stateful events, all keys must also be specified,
	// with variable values initialized to their starting state.
	Stateless bool
}

// EventEntry represents a extended version of an key value pair for a ax_event_key_value_set with additional metadata.
type EventEntry struct {
	// Key: The key for the underlaying ax_event_key_value_set
	Key string
	// Namespace: The namespace of the key or nil.
	Namespace *string
	// Value: The data associated with the key. This can be any type of data relevant to the event.
	Value interface{}
	// IsSource: An optional flag that, when set, marks the key as a source. Sources are identifiers
	// used to distinguish between multiple instances of the same event declaration. For example, if a device
	// has multiple I/O ports, the key representing which port the event is for could be marked as a source.
	// It's important to note that while multiple keys can be marked as source, only events with zero or one
	// source keys are eligible for triggering actions.
	IsSource *bool
	// IsData: An optional flag that, when set, marks the key as data. Data keys represent the state or value
	// of what the event is about, such as the high or low state of an I/O port. Similar to IsSource, although
	// it's possible to mark more than one key as data, only events with exactly one data key can be used to
	// trigger actions.
	IsData *bool
	// UserDefined: An optional field that allows users to attach a custom tag to the key-value pair. This tag
	// can be used for additional identification or categorization beyond what is provided by the key and namespace.
	UserDefined *string
	// KeyNiceName: An optional, human-readable name for the key. This is useful for providing clear, understandable
	// labels for keys when displayed to end-users, enhancing the usability and accessibility of event data.
	KeyNiceName *string
	// ValueNiceName: Similar to KeyNiceName, this is an optional, human-readable name for the value. It serves
	// the same purpose of enhancing clarity and understanding for end-users.
	ValueNiceName *string
	// ValueType: Specifies the type of the value, using the AXEventValueType enumeration. This helps ensure
	// consistent interpretation and handling of the value data across different parts of the system.
	ValueType axevent.AXEventValueType
}

type KeyValueMap map[string]interface{}

// NewEvent creates a new AXEvent based on predefined keys and dynamic values.
// It returns the new AXEvent or an error if the values provided do not match the expected types.
// The valuesMap parameter should contain the values for each key in the KeyValueSet.
// The keys in the valuesMap should match the keys in the KeyValueSet of the CameraPlatformEvent.
func (cpe *CameraPlatformEvent) NewEvent(valuesMap KeyValueMap) (*axevent.AXEvent, error) {
	var kvsEntries []axevent.KeyValueEntrie

	for _, entry := range cpe.Entries {
		value, exists := valuesMap[entry.Key]
		if !exists {
			return nil, fmt.Errorf("no value provided for key: %s", entry.Key)
		}

		switch entry.ValueType {
		case axevent.AXValueTypeInt:
			if intValue, ok := value.(int); ok {
				kvsEntries = append(kvsEntries, axevent.KeyValueEntrie{Key: entry.Key, Value: intValue, ValueType: entry.ValueType})
			} else {
				return nil, fmt.Errorf("type mismatch for key %s: expected int", entry.Key)
			}
		case axevent.AXValueTypeDouble:
			if floatValue, ok := value.(float64); ok {
				kvsEntries = append(kvsEntries, axevent.KeyValueEntrie{Key: entry.Key, Value: floatValue, ValueType: entry.ValueType})
			} else {
				return nil, fmt.Errorf("type mismatch for key %s: expected float64", entry.Key)
			}
		case axevent.AXValueTypeString:
			if stringValue, ok := value.(string); ok {
				kvsEntries = append(kvsEntries, axevent.KeyValueEntrie{Key: entry.Key, Value: stringValue, ValueType: entry.ValueType})
			} else {
				return nil, fmt.Errorf("type mismatch for key %s: expected string", entry.Key)
			}
		case axevent.AXValueTypeBool:
			if boolValue, ok := value.(bool); ok {
				kvsEntries = append(kvsEntries, axevent.KeyValueEntrie{Key: entry.Key, Value: boolValue, ValueType: entry.ValueType})
			} else {
				return nil, fmt.Errorf("type mismatch for key %s: expected bool", entry.Key)
			}
		default:
			return nil, fmt.Errorf("unsupported type for key %s", entry.Key)
		}
	}

	return axevent.NewAxEvent(axevent.NewAXEventKeyValueSetFromEntries(kvsEntries), nil), nil
}

// NewCameraApplicationPlatformEvent creates a new AXEventKeyValueSet instance for representing a Camera Application Platform event.
// This function encapsulates the process of initializing an event with specific application setup details, event identifiers,
// key-value pairs for event data, and various types of markers (source, data, user-defined) to provide additional context
// or categorization for the event data. Additionally, it facilitates assigning 'nice names' to event key-value pairs for
// enhanced readability.
//
// Parameters:
//   - app_setup: An axmanifest.Setup structure containing the application setup details. It includes information such as
//     application name and friendly name, which are used to contextualize the event within a specific application platform.
//   - event_name: A string representing the unique identifier of the event.
//   - event_nice_name: An optional pointer to a string that provides a human-readable name for the event. If provided, it
//     overrides the default event name in the context where 'nice names' are used.
//   - event_entries: A slice of pointers to EventEntry structures, each representing a key-value pair with additional metadata
//
// Returns:
//   - A pointer to an AXEventKeyValueSet instance.
//   - An error, which will be non-nil if any part of the event creation process fails.
//
// The function utilizes the NewTnsAxisEvent helper function to initialize the AXEventKeyValueSet, specifying a structured
// set of topics ('topic0' to 'topic3').
// Specifically, the topics are assigned as follows:
//   - 'topic0' is set to "CameraApplicationPlatform", identifying the event as part of the Camera Application Platform.
//     This serves as the primary categorization layer, indicating the event's general domain.
//   - 'topic1' is derived from the `app_setup.AppName`, tying the event to a specific application by its name. This
//     further refines the event's context within the platform, associating it with a particular application's events.
//   - 'topic2' is optionally set to a user-provided string via `event_name` or `event_nice_name`, if provided. This allows
//     for a more descriptive labeling of the event, enhancing the readability and interpretability of the event data.
//     If `event_nice_name` is not null, it prefixes the nice name with the app's friendly name for clearer identification.
//     If both are null, `topic2` effectively utilizes the raw `event_name` for technical identification.
//   - 'topic3' is intentionally left as nil/null.
func NewCameraApplicationPlatformEvent(app_setup axmanifest.Setup, event_name string, event_nice_name *string, event_entries []*EventEntry) (*axevent.AXEventKeyValueSet, error) {

	var kvs_entries []*axevent.KeyValueEntrie
	for _, entry := range event_entries {
		kvs_entries = append(kvs_entries, &axevent.KeyValueEntrie{
			Key:       entry.Key,
			Namespace: entry.Namespace,
			Value:     entry.Value,
			ValueType: entry.ValueType,
		})
	}

	kvs, err := axevent.NewTnsAxisEvent(
		"CameraApplicationPlatform",
		app_setup.AppName,
		utils.StrPtr(event_name),
		nil,
		kvs_entries,
	)

	if err != nil {
		return nil, err
	}

	for _, entry := range event_entries {
		if entry.IsData != nil && *entry.IsData {
			if err := kvs.MarkAsData(entry.Key, entry.Namespace); err != nil {
				return nil, err
			}
		}
		if entry.IsSource != nil && *entry.IsSource {
			if err := kvs.MarkAsSource(entry.Key, entry.Namespace); err != nil {
				return nil, err
			}
		}
		if entry.UserDefined != nil {
			if err := kvs.MarkAsUserDefined(entry.Key, entry.Namespace, entry.UserDefined); err != nil {
				return nil, err
			}
		}
		if entry.KeyNiceName != nil || entry.ValueNiceName != nil {
			if err := kvs.AddNiceNames(entry.Key, entry.Namespace, entry.KeyNiceName, entry.ValueNiceName); err != nil {
				return nil, err
			}
		}
	}

	var nice_name string
	if event_nice_name != nil {
		nice_name = fmt.Sprintf("%s: %s", app_setup.FriendlyName, *event_nice_name)
	} else {
		nice_name = fmt.Sprintf("%s: %s", app_setup.FriendlyName, event_name)
	}

	if err := kvs.AddNiceNames("topic2", &axevent.OnfivNameSpaceTnsAxis, nil, utils.StrPtr(nice_name)); err != nil {
		return nil, err
	}

	return kvs, nil
}

// OnEvent creates a subscription callback for the given event key value set.
func (a *AcapApplication) OnEvent(kvs *axevent.AXEventKeyValueSet, callback func(*axevent.Event)) (subscription int, err error) {
	return a.EventHandler.OnEvent(kvs, callback)
}

// UnmarshalEvent unmarshals the given event into the provided struct.
func UnmarshalEvent(e *axevent.Event, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("value must be a pointer to a struct")
	}

	for i := 0; i < val.Elem().NumField(); i++ {
		field := val.Elem().Field(i)
		if !field.CanSet() {
			continue
		}
		fieldType := val.Elem().Type().Field(i)

		key := fieldType.Tag.Get("eventKey")
		if key == "" {
			key = strings.ToLower(fieldType.Name)
		}

		switch field.Kind() {
		case reflect.Int:
			if intValue, err := e.Kvs.GetInteger(key, nil); err == nil {
				field.SetInt(int64(intValue))
			} else {
				return fmt.Errorf("error getting integer for key %s: %v", key, err)
			}
		case reflect.Float64:
			if fValue, err := e.Kvs.GetDouble(key, nil); err == nil {
				field.SetFloat(fValue)
			} else {
				return fmt.Errorf("error getting double for key %s: %v", key, err)
			}
		case reflect.String:
			if sValue, err := e.Kvs.GetString(key, nil); err == nil {
				field.SetString(sValue)
			} else {
				return fmt.Errorf("error getting string for key %s: %v", key, err)
			}
		case reflect.Bool:
			if boolValue, err := e.Kvs.GetBoolean(key, nil); err == nil {
				field.SetBool(boolValue)
			} else {
				return fmt.Errorf("error getting boolean for key %s: %v", key, err)
			}
		}
	}

	return nil
}
