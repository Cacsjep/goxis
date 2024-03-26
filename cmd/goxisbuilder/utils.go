package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	dac "github.com/Snawoot/go-http-digest-auth-client"
	"github.com/docker/docker/client"
)

func boolToStr(b bool) string {
	if b {
		return "YES"
	}
	return "NO"
}

func ptr(s string) *string {
	return &s
}

func getLog(url string, pwd string) {
	client := &http.Client{
		Transport: dac.NewDigestTransport("root", pwd, http.DefaultTransport),
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	displayLastLines(string(body), 70)
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func displayLastLines(logContent string, nline int) {
	fmt.Println("Update..")
	time.Sleep(time.Millisecond * 500)
	clearScreen()
	lines := strings.Split(logContent, "\n")
	startLine := 0
	if len(lines) > nline {
		startLine = len(lines) - nline
	}
	for _, line := range lines[startLine:] {
		fmt.Println(line)
	}
}

// handleError logs an error message and exits the program with a status code.
func handleError(message string, err error) {
	log.Printf("Error: %s: %v\n", message, err)
	os.Exit(1) // Exit with a status code indicating failure.
}

// configureArchitecture sets up the build configuration based on the architecture.
func configureArchitecture(arch string, buildConfig *BuildConfiguration) {
	switch arch {
	case "aarch64":
		buildConfig.ImageName = "acap:aarch64"
		buildConfig.GoArch = "arm64"
		buildConfig.CrossPrefix = "aarch64-linux-gnu-"
	case "armv7hf":
		buildConfig.ImageName = "acap:arm"
		buildConfig.GoArch = "arm"
		buildConfig.GoArm = "7"
		buildConfig.CrossPrefix = "arm-linux-gnueabihf-"
	default:
		handleError("Architecture invalid", fmt.Errorf("should be either aarch64 or armv7hf, got %s", arch))
	}
}

// configureSdk sets up the build configuration based on the lowest Sdk flag.
func configureSdk(lowestSdkVersion bool, buildConfig *BuildConfiguration) {
	if lowestSdkVersion {
		buildConfig.Sdk = "acap-sdk"
		buildConfig.UbunutVersion = "20.04"
		buildConfig.Version = "3.5"
	} else {
		buildConfig.Sdk = "acap-native-sdk"
		buildConfig.UbunutVersion = "22.04"
		buildConfig.Version = "1.13"
	}
}

// buildApplication handles the building of the application or examples based on the provided configuration.
func buildApplication(ctx context.Context, cli *client.Client, buildConfig *BuildConfiguration) {
	if buildConfig.AppDirectory == "" {
		if !buildConfig.BuildExamples {
			handleError("Application directory must be set or build examples flag enabled", nil)
		}
		for _, e := range examples {
			buildConfig.AppDirectory = fmt.Sprintf("examples/%s", e)
			if err := buildAndRunContainer(ctx, cli, buildConfig); err != nil {
				handleError("Build fails", err)
			}
		}
		return
	} else {
		// Application build logic here...
		if err := buildAndRunContainer(ctx, cli, buildConfig); err != nil {
			handleError("Failed to build and run container", err)
		}
	}
}

func watchPackageLog(buildConfig *BuildConfiguration) {
	// Setup a channel to listen for interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	url := fmt.Sprintf("http://%s/axis-cgi/admin/systemlog.cgi?appname=%s", buildConfig.Ip, buildConfig.Manifest.ACAPPackageConf.Setup.AppName)
Loop:
	for {
		select {
		case <-ticker.C:
			getLog(url, buildConfig.Pwd)
		case <-sigChan:
			fmt.Println("Interrupt received, stopping...")
			break Loop
		}
	}
}
