package main

import "github.com/Cacsjep/goxis/pkg/manifest"

type BuildConfiguration struct {
	Manifest      *manifest.ApplicationManifestSchema
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
