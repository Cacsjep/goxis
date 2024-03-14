package manifest

import (
	"encoding/json"
	"io"
	"os"
)

type ApplicationManifestSchema struct {
	SchemaVersion   string          `json:"schemaVersion" validate:"required,regexp=^1\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)$"`
	Resources       Resources       `json:"resources,omitempty"`
	ACAPPackageConf ACAPPackageConf `json:"acapPackageConf" validate:"required"`
}

// Resources represents host resources the application requires access to.
type Resources struct {
	DBus  DBusResources  `json:"dbus,omitempty"`
	Linux LinuxResources `json:"linux,omitempty"`
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
	EmbeddedSdkVersion string `json:"embeddedSdkVersion" validate:"omitempty,semver"`
	FriendlyName       string `json:"friendlyName" validate:"omitempty"`
	User               User   `json:"user" validate:"required,dive"`
	RunMode            string `json:"runMode" validate:"required,oneof=respawn once never"`
	RunOptions         string `json:"runOptions" validate:"omitempty,regexp=^[\\w /\\.\\=-]+$"`
	Vendor             string `json:"vendor" validate:"required,regexp=^[\\w /\\-\\(\\),\\.\\!\\?&']+$"`
	VendorUrl          string `json:"vendorUrl" validate:"omitempty,url"`
	Version            string `json:"version" validate:"required,semver"`
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
