package main

import (
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
)

const (
	Blue  = "\033[34m"
	Reset = "\033[0m"
	Green = "\033[32m"
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

func listEapDirectory() {
	entries, err := os.ReadDir("./eap")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		fmt.Println("EAP:", e.Name())
	}
}

func printCompatibility(buildConfig *BuildConfiguration) {
	fmt.Println("Acap Compatibility:")
	// Maps for SDK to Firmware compatibility
	sdkToFirmware := map[string]string{
		"3.0": "9.70 and later",
		"3.1": "9.80 (LTS) and later",
		"3.2": "10.2 and later",
		"3.3": "10.5 and later",
		"3.4": "10.6 and later",
		"3.5": "10.9 and later",
	}

	// Maps for Native SDK version to Firmware compatibility
	nativeSdkToFirmware := map[string]string{
		"1.0":  "10.7 and later until LTS",
		"1.1":  "10.9 and later until LTS",
		"1.2":  "10.10 and later until LTS",
		"1.3":  "10.12 (LTS)",
		"1.4":  "11.0 and later until LTS",
		"1.5":  "11.1 and later until LTS",
		"1.6":  "11.2 and later until LTS",
		"1.7":  "11.3 and later until LTS",
		"1.8":  "11.4 and later until LTS",
		"1.9":  "11.5 and later until LTS",
		"1.10": "11.6 and later until LTS",
		"1.11": "11.7 and later until LTS",
		"1.12": "11.8 and later until LTS",
		"1.13": "11.9 and later until LTS",
	}

	// Check if it's using the native SDK or standard SDK
	if buildConfig.Sdk == "acap-native-sdk" {
		if firmware, ok := nativeSdkToFirmware[buildConfig.Version]; ok {
			fmt.Printf("     ACAP Native SDK %s%s%s, compatible with AXIS OS version: %s%s%s\n", Blue, buildConfig.Version, Reset, Green, firmware, Reset)
		} else {
			log.Printf("     Unknown ACAP Native SDK version: %s\n", buildConfig.Version)
		}
	} else if buildConfig.Sdk == "acap-sdk" {
		if firmware, ok := sdkToFirmware[buildConfig.Version]; ok {
			fmt.Printf("     ACAP3 SDK %s%s%s, compatible with firmware version: %s%s%s\n", Blue, buildConfig.Version, Reset, Green, firmware, Reset)
		} else {
			log.Printf("     Unknown ACAP3 SDK version: %s\n", buildConfig.Version)
		}
	} else {
		log.Printf("     Unknown SDK configuration: %s\n", buildConfig.Sdk)
	}

	schemaToFirmware := map[string]string{
		"1.0":   "10.7",
		"1.1":   "10.7",
		"1.2":   "10.7",
		"1.3":   "10.9",
		"1.3.1": "11.0",
		"1.4.0": "11.7",
		"1.5.0": "11.8",
		"1.6.0": "11.9",
	}

	if firmware, ok := schemaToFirmware[buildConfig.Manifest.SchemaVersion]; ok {
		fmt.Printf("     Schema %s%s%s is compatible with firmware version: %s%s%s\n", Blue, buildConfig.Manifest.SchemaVersion, Reset, Green, firmware, Reset)
	} else {
		log.Printf("     Unknown Schema version: %s\n", buildConfig.Manifest.SchemaVersion)
	}

	archToChips := map[string][]string{
		"armv7hf": {"ARTPEC-6", "ARTPEC-7", "i.MX 6SoloX", "i.MX 6ULL"},
		"aarch64": {"ARTPEC-8", "CV25", "S5", "S5L"},
	}

	if chips, ok := archToChips[buildConfig.Arch]; ok {
		chipsStr := strings.Join(chips, ", ")

		fmt.Printf("     Supported architecture: %s%s%s with chips: %s%s%s\n",
			Blue, buildConfig.Arch, Reset,
			Green, chipsStr, Reset)
	} else {
		fmt.Println("     Unsupported architecture.")
	}
}
