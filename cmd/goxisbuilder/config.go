package main

import "github.com/Cacsjep/goxis/pkg/axmanifest"

// BuildConfiguration defines the configuration parameters for building
// the EAP application, including details such as architecture, manifest details,
// and flags indicating whether to install the application, start it,
// build examples, or watch logs.
type BuildConfiguration struct {
	Manifest      *axmanifest.ApplicationManifestSchema
	ManifestPath  string
	ImageName     string
	Ip            string
	Pwd           string
	Arch          string
	DoStart       bool
	DoInstall     bool
	BuildExamples bool
	AppDirectory  string
	Sdk           string
	UbunutVersion string
	Version       string
	WithLibav     bool
	GoArch        string
	GoArm         string
	CrossPrefix   string
	LowestSdk     bool
	Watch         bool
}
