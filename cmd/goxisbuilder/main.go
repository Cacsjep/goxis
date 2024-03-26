package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/Cacsjep/goxis/pkg/manifest"
)

func main() {
	showHelp := flag.Bool("h", false, "Show usage")
	ip := flag.String("ip", "", "IP for camera where eap is installed")
	manifest_path := flag.String("manifest", "manifest.json", "Specify manifest default is manifest.json")
	pwd := flag.String("pwd", "", "Root password for camera where eap is installed")
	arch := flag.String("arch", "aarch64", "ACAP Architecture: aarch64 or armv7hf")
	doStart := flag.Bool("start", false, "Start after install")
	doInstall := flag.Bool("install", false, "Install on camera")
	buildExamples := flag.Bool("build-examples", false, "Build Examples")
	lowesSdkVersion := flag.Bool("lowsdk", false, "Build for Firmware > 10.9 With SDK version 1.1 your manifest should use: 1.3")
	watch := flag.Bool("watch", false, "Watch the package log after build")
	appDirectory := flag.String("appdir", "", "Full path of application directroy to build from")
	withLibav := flag.Bool("libav", false, "Compile libav for binding it with go-astiav")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(1)
	}

	amf, err := manifest.LoadManifest(path.Join(*appDirectory, *manifest_path))
	if err != nil {
		panic(err)
	}

	buildConfig := BuildConfiguration{
		Arch:          *arch,
		Manifest:      amf,
		ManifestPath:  *manifest_path,
		Ip:            *ip,
		Pwd:           *pwd,
		DoStart:       *doStart,
		DoInstall:     *doInstall,
		BuildExamples: *buildExamples,
		LowestSdk:     *lowesSdkVersion,
		WithLibav:     *withLibav,
		Watch:         *watch,
	}
	ctx := context.Background()
	cli, err := newDockerClient()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	if *lowesSdkVersion {
		buildConfig.Sdk = "acap-sdk"
		buildConfig.UbunutVersion = "20.04"
		buildConfig.Version = "3.5"
	} else {
		buildConfig.Sdk = "acap-native-sdk"
		buildConfig.UbunutVersion = "22.04"
		buildConfig.Version = "1.13"
	}

	if *arch == "aarch64" {
		buildConfig.ImageName = "acap:aarch64"
		buildConfig.GoArch = "arm64"
		buildConfig.CrossPrefix = "aarch64-linux-gnu-"
	} else if *arch == "armv7hf" {
		buildConfig.ImageName = "acap:arm"
		buildConfig.GoArch = "arm"
		buildConfig.GoArm = "7"
		buildConfig.CrossPrefix = "arm-linux-gnueabihf-"
	} else {
		panic("Architecture invalid should be either aarch64 or armv7hf")
	}

	if *appDirectory == "" {
		if !*buildExamples {
			panic("appdir must be set")
		}
		for _, e := range examples {
			buildConfig.AppDirectory = fmt.Sprintf("examples/%s", e)
			if err := buildAndRunContainer(ctx, cli, &buildConfig); err != nil {
				log.Fatalf("Error: %s\n", err)
			}

		}
		return
	} else {
		buildConfig.AppDirectory = *appDirectory
		if err := buildAndRunContainer(ctx, cli, &buildConfig); err != nil {
			log.Fatalf("Error: %s\n", err)
		}
	}

	if buildConfig.Watch {
		// Setup a channel to listen for interrupt signal (Ctrl+C)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		url := fmt.Sprintf("http://%s/axis-cgi/admin/systemlog.cgi?appname=%s", *ip, amf.ACAPPackageConf.Setup.AppName)
	Loop:
		for {
			select {
			case <-ticker.C:
				getLog(url, *pwd)
			case <-sigChan:
				fmt.Println("Interrupt received, stopping...")
				break Loop
			}
		}
	}
}
