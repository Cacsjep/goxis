// Package axmanifest provides a Go implementation of the Application Manifest Schema for Axis Camera Application Platform (ACAP) applications.
package axmanifest

import (
	"encoding/json"
	"io"
	"os"
)

// ApplicationManifestSchema is an ACAP application manifest, schema 1.x or 2.x.
// Schema 2 (ACAP Native SDK 12.10+, required on AXIS OS 13) adds
// setup.vendorId and setup.compatibleOsVersions, makes setup.architecture
// required and drops setup.embeddedSdkVersion.
type ApplicationManifestSchema struct {
	SchemaVersion   string          `json:"schemaVersion" validate:"required,regexp=^[12]\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)$"`
	Resources       Resources       `json:"resources,omitempty"`
	ACAPPackageConf ACAPPackageConf `json:"acapPackageConf" validate:"required"`
}

// Resources represents host resources the application requires access to.
type Resources struct {
	DBus  DBusResources  `json:"dbus,omitempty"`
	Linux LinuxResources `json:"linux,omitempty"`
	// DeviceDataHub grants access to the Device Data Hub (schema 1.11.0 or
	// 2.2.0 and later).
	DeviceDataHub *DeviceDataHubResource `json:"deviceDataHub,omitempty"`
	// DeepLearningProcessor declares use of the DLPU (schema 2.x). On AXIS
	// OS 13 an application that does not declare it gets no DLPU.
	DeepLearningProcessor *DeepLearningProcessorResource `json:"deepLearningProcessor,omitempty"`
}

// DeviceDataHubResource enables the Device Data Hub for the application.
// Subscribing to public topics (such as AXIS Scene Metadata) needs only
// Enabled. AccessControlList entries grant rights on the application's own
// topics; they are kept as raw JSON (see the manifest schema in the SDK).
type DeviceDataHubResource struct {
	Enabled           bool              `json:"enabled"`
	AccessControlList []json.RawMessage `json:"accessControlList,omitempty"`
}

// DeepLearningProcessorResource declares the deep learning processor (DLPU).
type DeepLearningProcessorResource struct {
	Enabled  bool `json:"enabled"`
	Required bool `json:"required,omitempty"`
}

// DBusResources details D-Bus resources on the host system that the application requires or desires access to.
type DBusResources struct {
	RequiredMethods    []string `json:"requiredMethods,omitempty"`
	ConditionalMethods []string `json:"conditionalMethods,omitempty"`
}

// LinuxResources specifies Linux resources on the host system that the application requires or desires access to.
type LinuxResources struct {
	User LinuxUser `json:"user,omitempty"`
}

// LinuxUser describes a dynamic user the application shall run as.
type LinuxUser struct {
	Groups []string `json:"groups,omitempty"`
}

// ACAPPackageConf represents an ACAP package configuration object.
type ACAPPackageConf struct {
	Setup          Setup          `json:"setup" validate:"required"`
	Installation   Installation   `json:"installation,omitempty"`
	Uninstallation Uninstallation `json:"uninstallation,omitempty"`
	Configuration  Configuration  `json:"configuration,omitempty"`
	CopyProtection CopyProtection `json:"copyProtection"`
}

// Setup includes ACAP application identification and information settings.
type Setup struct {
	AppName            string `json:"appName" validate:"required,alphanum,max=26"`
	AppID              string `json:"appId" validate:"omitempty,numeric"`
	Architecture       string `json:"architecture" validate:"omitempty,oneof=all aarch64 armv7hf"`
	EmbeddedSdkVersion string `json:"embeddedSdkVersion,omitempty" validate:"omitempty,semver"` // schema 1.x only
	FriendlyName       string `json:"friendlyName" validate:"omitempty"`
	User               User   `json:"user" validate:"required,dive"`
	RunMode            string `json:"runMode" validate:"required,oneof=respawn once never"`
	RunOptions         string `json:"runOptions" validate:"omitempty,regexp=^[\\w /\\.\\=-]+$"`
	Vendor             string `json:"vendor" validate:"required,regexp=^[\\w /\\-\\(\\),\\.\\!\\?&']+$"`
	VendorUrl          string `json:"vendorUrl" validate:"omitempty,url"`
	Version            string `json:"version" validate:"required,semver"`
	// VendorID is the vendor id the AXIS ACAP Service Portal assigns
	// (schema 2.x, required).
	VendorID string `json:"vendorId,omitempty" validate:"omitempty,hexadecimal,len=10"`
	// CompatibleOsVersions lists the AXIS OS version ranges the application
	// supports (schema 2.x, required). From AXIS OS 12.10 the device refuses
	// to install an application outside every range.
	CompatibleOsVersions []OsVersionRange `json:"compatibleOsVersions,omitempty"`
}

// OsVersionRange is an AXIS OS version range, such as {Min: "12.10", Max: "13"}.
// A missing patch or minor version includes all of them.
type OsVersionRange struct {
	Min string `json:"min,omitempty"`
	Max string `json:"max,omitempty"`
}

// User defines static user and group information.
type User struct {
	Username string `json:"username" validate:"required,lowercase,alphanum,startswithalpha,max=32"`
	Group    string `json:"group" validate:"required,lowercase,alphanum,startswithalpha,max=32"`
}

// Installation holds ACAP application installation settings.
type Installation struct {
	PostInstallScript string `json:"postInstallScript" validate:"omitempty,regexp=^[\\w\\.-]+$"`
}

// Uninstallation encompasses ACAP application uninstallation settings.
type Uninstallation struct {
	PreUninstallScript string `json:"preUninstallScript" validate:"omitempty,regexp=^[\\w\\.-]+$"`
}

// Configuration pertains to ACAP application interaction setup.
type Configuration struct {
	SettingPage  string             `json:"settingPage" validate:"omitempty,regexp=^[\\w\\.-]+$"`
	HttpConfig   []HttpConfigItem   `json:"httpConfig,omitempty"`
	ParamConfig  []ParamConfigItem  `json:"paramConfig,omitempty"`
	ReverseProxy []ReverseProxyItem `json:"reverseProxy,omitempty"`
}

// HttpConfigItem describes a web server configuration object, which can be either a CGI configuration or web content configuration.
type HttpConfigItem struct {
	Type   string `json:"type" validate:"required,oneof=transferCgi fastCgi directory"`
	Name   string `json:"name" validate:"omitempty,regexp=^[\\w/\\.-]+$"`
	Access string `json:"access" validate:"omitempty,oneof=admin operator viewer"`
}

// ParamConfigItem defines a parameter configuration object for interaction with the Parameter API or the VAPIX API.
type ParamConfigItem struct {
	Name    string `json:"name" validate:"required"`
	Default string `json:"default" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

// ReverseProxyItem describes a reverse proxy configuration object, which can be for a unix domain socket or a TCP connection.
type ReverseProxyItem struct {
	ApiPath string `json:"apiPath" validate:"required,regexp=^[\\w\\-.+*?$\\(){}\\[\\]\\\\][\\w\\-.+*?$\\(){}\\[\\]\\\\/]+$"`
	ApiType string `json:"apiType" validate:"omitempty,oneof=http ws fcgi"`
	Target  string `json:"target" validate:"required"`
	Access  string `json:"access" validate:"required,oneof=admin operator viewer anonymous"`
}

// CopyProtection outlines the ACAP application's copy protection utilization.
type CopyProtection struct {
	Method string `json:"method" validate:"required,oneof=none axis custom"`
}

func LoadManifest(manifest_path string) (*ApplicationManifestSchema, error) {
	var err error
	var m_file *os.File
	var m_content []byte
	ams := ApplicationManifestSchema{}

	if m_file, err = os.Open(manifest_path); err != nil {
		return nil, err
	}
	defer m_file.Close()

	if m_content, err = io.ReadAll(m_file); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(m_content, &ams); err != nil {
		return nil, err
	}

	return &ams, nil
}
