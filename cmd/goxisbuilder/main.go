package main

import (
	"context"
	"flag"
	"os"
	"path"

	"github.com/Cacsjep/goxis/pkg/manifest"
)

func main() {
	showHelp := flag.Bool("h", false, "Displays this help message.")
	ip := flag.String("ip", "", "The IP address of the camera where the EAP application is installed.")
	manifest_path := flag.String("manifest", "manifest.json", "The path to the manifest file. Defaults to 'manifest.json'.")
	pwd := flag.String("pwd", "", "The root password for the camera where the EAP application is installed.")
	arch := flag.String("arch", "aarch64", "The architecture for the ACAP application: 'aarch64' or 'armv7hf'.")
	doStart := flag.Bool("start", false, "Set to true to start the application after installation.")
	doInstall := flag.Bool("install", false, "Set to true to install the application on the camera.")
	buildExamples := flag.Bool("build-examples", false, "Set to true to build example applications.")
	lowestSdkVersion := flag.Bool("lowsdk", false, "Set to true to build for firmware versions greater than 10.9 with SDK version 1.1. This adjusts the manifest to use version 1.3.")
	watch := flag.Bool("watch", false, "Set to true to monitor the package log after building.")
	appDirectory := flag.String("appdir", "", "The full path to the application directory from which to build.")
	withLibav := flag.Bool("libav", false, "Set to true to compile libav for binding with go-astiav.")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(1)
	}

	// Load and validate the application manifest file specified by the user.
	amf, err := manifest.LoadManifest(path.Join(*appDirectory, *manifest_path))
	if err != nil {
		panic(err)
	}

	// Setup the build configuration based on the parsed flags and loaded manifest.
	buildConfig := BuildConfiguration{
		AppDirectory:  *appDirectory,
		Arch:          *arch,
		Manifest:      amf,
		ManifestPath:  *manifest_path,
		Ip:            *ip,
		Pwd:           *pwd,
		DoStart:       *doStart,
		DoInstall:     *doInstall,
		BuildExamples: *buildExamples,
		LowestSdk:     *lowestSdkVersion,
		WithLibav:     *withLibav,
		Watch:         *watch,
	}
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		handleError("Failed create new docker client", err)
	}

	configureSdk(*lowestSdkVersion, &buildConfig)
	configureArchitecture(*arch, &buildConfig)
	buildApplication(ctx, cli, &buildConfig)

	if buildConfig.Watch {
		watchPackageLog(&buildConfig)
	}
}
