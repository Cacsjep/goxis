package main

import (
	"archive/tar"
	"context"
	"encoding/json"
	"flag"
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

	"github.com/Cacsjep/goxis/pkg/manifest"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"

	dac "github.com/Snawoot/go-http-digest-auth-client"
)

type BuildConfiguration struct {
	buildArgs map[string]*string
	imageName string
}

type BuildFlags struct {
	Ip           string
	Pwd          string
	Arch         string
	DoStart      bool
	DoInstall    bool
	AppDirectory string
	WithLibav    bool
}

var examples []string = []string{
	"webserver",
	"axparameter",
	"axevent",
	"axstorage",
	"license",
	"vdostream",
}

func boolToStr(b bool) string {
	if b {
		return "YES"
	}

	return "NO"
}

func buildArm64(amf *manifest.ApplicationManifestSchema, bf BuildFlags) *BuildConfiguration {
	buildArgs := map[string]*string{
		"ARCH":         ptr("aarch64"),
		"GO_ARCH":      ptr("arm64"),
		"APP_NAME":     ptr(amf.ACAPPackageConf.Setup.AppName),
		"IP_ADDR":      ptr(bf.Ip),
		"PASSWORD":     ptr(bf.Pwd),
		"START":        ptr(boolToStr(bf.DoStart)),
		"INSTALL":      ptr(boolToStr(bf.DoInstall)),
		"GO_APP":       ptr(bf.AppDirectory),
		"CROSS_FILE":   ptr("cross_aarch64.txt"),
		"CROSS_PREFIX": ptr("aarch64-linux-gnu-"),
		"COMP_LIBAV":   ptr(boolToStr(bf.WithLibav)),
	}
	return &BuildConfiguration{buildArgs: buildArgs, imageName: "acaparm64"}
}

func buildArmv7hf(amf *manifest.ApplicationManifestSchema, bf BuildFlags) *BuildConfiguration {
	buildArgs := map[string]*string{
		"ARCH":         ptr("armv7hf"),
		"GO_ARCH":      ptr("arm"),
		"GO_ARM":       ptr("7"),
		"APP_NAME":     ptr(amf.ACAPPackageConf.Setup.AppName),
		"IP_ADDR":      ptr(bf.Ip),
		"PASSWORD":     ptr(bf.Pwd),
		"START":        ptr(boolToStr(bf.DoStart)),
		"INSTALL":      ptr(boolToStr(bf.DoInstall)),
		"GO_APP":       ptr(bf.AppDirectory),
		"CROSS_FILE":   ptr("cross_armv7hf.txt"),
		"CROSS_PREFIX": ptr("arm-linux-gnueabihf-"),
		"COMP_LIBAV":   ptr(boolToStr(bf.WithLibav)),
	}
	return &BuildConfiguration{buildArgs: buildArgs, imageName: "acaparmv7hf"}
}

func dockerBuild(bc *BuildConfiguration) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("Current Directory:", currentDir)
	buildContext, err := archive.TarWithOptions(currentDir, &archive.TarOptions{})
	if err != nil {
		panic(err)
	}
	defer buildContext.Close()
	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{bc.imageName},
		BuildArgs:  bc.buildArgs,
		Remove:     true,
	}
	buildResponse, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		panic(err)
	}
	defer buildResponse.Body.Close()
	decoder := json.NewDecoder(buildResponse.Body)
	fmt.Println("--- Start Docker Image Build ---")
	for {
		var m map[string]interface{}
		if err := decoder.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if stream, ok := m["stream"]; ok {
			fmt.Print(stream.(string))
		}
	}

	// Create a container from the built image
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: bc.imageName,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	copyFromContainer, _, err := cli.CopyFromContainer(ctx, resp.ID, "/opt/eap")
	if err != nil {
		panic(err)
	}
	defer copyFromContainer.Close()

	os.Mkdir("eap", 664)

	tr := tar.NewReader(copyFromContainer)
	var foundFile bool
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}

		if header.Typeflag == tar.TypeReg {
			outputFile, err := os.Create(header.Name)
			if err != nil {
				continue
			}
			defer outputFile.Close()

			if _, err := io.Copy(outputFile, tr); err != nil {
				continue
			}
			foundFile = true
		}
	}

	if !foundFile {
		println("No file found in the archive")
	}

	// Stop and remove the container after copying the file
	if err := cli.ContainerStop(ctx, resp.ID, container.StopOptions{}); err != nil {
		panic(err)
	}

	if err := cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Complete")
}

func main() {
	showHelp := flag.Bool("h", false, "Show usage")
	ip := flag.String("ip", "", "IP for camera where eap is installed")
	pwd := flag.String("pwd", "", "Root password for camera where eap is installed")
	arch := flag.String("arch", "aarch64", "ACAP Architecture: aarch64 or armv7hf")
	doStart := flag.Bool("start", false, "Start after install")
	doInstall := flag.Bool("install", false, "Install on camera")
	buildExamples := flag.Bool("build-examples", false, "Build Examples")
	getPackageLog := flag.Bool("watch", false, "Watch the package log after build")
	appDirectory := flag.String("appdir", "", "Full path of application directroy to build from")
	withLibav := flag.Bool("libav", false, "Compile libav for binding it with go-astiav")
	flag.Parse()
	if *showHelp {
		flag.Usage()
		os.Exit(1)
	}

	if *appDirectory == "" {
		if !*buildExamples {
			panic("appdir must be set")
		}
		for _, e := range examples {
			doNot := false
			example := fmt.Sprintf("examples/%s", e)
			buildApp(&example, arch, ip, pwd, doInstall, &doNot, &doNot, &doNot)
		}
	} else {
		buildApp(appDirectory, arch, ip, pwd, doInstall, withLibav, doStart, getPackageLog)
	}
}

func buildApp(appDirectory *string, arch *string, ip *string, pwd *string, doInstall *bool, withLibav *bool, doStart *bool, getPackageLog *bool) {
	fmt.Println("### Starting ACAP Builder ###")
	amf, err := manifest.LoadManifest(*appDirectory + "/manifest.json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Build: %s \n", amf.ACAPPackageConf.Setup.AppName)
	fmt.Printf("Architecture: %s \n", *arch)
	fmt.Printf("Version: %s \n", amf.ACAPPackageConf.Setup.Version)
	fmt.Printf("Vendor: %s \n", amf.ACAPPackageConf.Setup.Vendor)
	fmt.Printf("Build Libav: %t \n", *withLibav)

	fmt.Printf("IP: %s \n", *ip)
	fmt.Printf("Pwd: %s \n", *pwd)
	fmt.Printf("Install: %t \n", *doInstall)
	fmt.Printf("Start after Install: %t \n", *doStart)
	fmt.Printf("Watch after Install: %t \n", *doStart)

	if (*doStart || *doInstall || *getPackageLog) && (*ip == "" || *pwd == "") {
		panic("When install/starting/watch is used, you need to provide both ip and pwd")
	}

	bf := BuildFlags{Ip: *ip, Pwd: *pwd, DoStart: *doStart, DoInstall: *doInstall, Arch: *arch, AppDirectory: *appDirectory, WithLibav: *withLibav}

	if *arch == "aarch64" {
		dockerBuild(buildArm64(amf, bf))
	} else if *arch == "armv7hf" {
		dockerBuild(buildArmv7hf(amf, bf))
	} else {
		panic("Architecture invalid should be either aarch64 or armv7hf")
	}

	if *getPackageLog {
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
