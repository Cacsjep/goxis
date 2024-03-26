package main

import (
	"archive/tar"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// newDockerClient initializes a new Docker client
func newDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return cli, nil
}

// buildAndRunContainer builds a Docker image and runs a container from it
func buildAndRunContainer(ctx context.Context, cli *client.Client, bc *BuildConfiguration) error {
	fmt.Printf("Building image '%s'...\n", bc.ImageName)

	// Build Docker image
	if err := dockerBuild(ctx, cli, bc); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	// Create and start container
	containerID, err := createContainer(ctx, cli, bc.ImageName)
	if err != nil {
		return fmt.Errorf("create container failed: %w", err)
	}

	fmt.Printf("Container '%s' created and started successfully.\n", bc.ImageName)

	if err := copyFromContainer(ctx, cli, containerID); err != nil {
		return fmt.Errorf("copy eap failed: %w", err)
	}

	fmt.Printf("Container data (eap)'%s' copied successfully.\n", bc.ImageName)

	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		panic(err)
	}

	if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{}); err != nil {
		panic(err)
	}

	return nil
}

// dockerBuild performs the Docker image build operation and processes the output
func dockerBuild(ctx context.Context, cli *client.Client, bc *BuildConfiguration) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to create current dir: %w", err)
	}
	buildContext, err := archive.TarWithOptions(currentDir, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("failed to create build context: %w", err)
	}
	defer buildContext.Close()

	options := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{bc.ImageName},
		BuildArgs: map[string]*string{
			"ARCH":           ptr(bc.Arch),
			"SDK":            ptr(bc.Sdk),
			"UBUNTU_VERSION": ptr(bc.UbunutVersion),
			"VERSION":        ptr(bc.Version),
			"GO_ARCH":        ptr(bc.GoArch),
			"GO_ARM":         ptr(bc.GoArm),
			"APP_NAME":       ptr(bc.Manifest.ACAPPackageConf.Setup.AppName),
			"APP_MANIFEST":   ptr(bc.ManifestPath),
			"IP_ADDR":        ptr(bc.Ip),
			"PASSWORD":       ptr(bc.Pwd),
			"START":          ptr(boolToStr(bc.DoStart)),
			"INSTALL":        ptr(boolToStr(bc.DoInstall)),
			"GO_APP":         ptr(bc.AppDirectory),
			"CROSS_PREFIX":   ptr(bc.CrossPrefix),
			"COMP_LIBAV":     ptr(boolToStr(bc.WithLibav)),
		},
		Remove: true,
	}

	buildResponse, err := cli.ImageBuild(ctx, buildContext, options)
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
			panic(err) // Consider handling this error more gracefully
		}

		if stream, ok := m["stream"]; ok {
			fmt.Print(stream.(string))
		} else if errMsg, ok := m["error"]; ok {
			return fmt.Errorf("build error: %s", errMsg)
		} else if status, ok := m["status"]; ok {
			fmt.Println("Status:", status)
		} else if progress, ok := m["progress"]; ok {
			fmt.Println("Progress:", progress)
		} else if aux, ok := m["aux"]; ok {
			fmt.Println("Aux:", aux)
		} else {
			fmt.Println("UNHANDLED MESSAGE", m)
		}
	}

	return nil
}

// createContainer creates and starts a Docker container from an image
func createContainer(ctx context.Context, cli *client.Client, imageName string) (string, error) {
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("container creation failed: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("container start failed: %w", err)
	}

	return resp.ID, nil
}

// copyFromContainer copy our build result
func copyFromContainer(ctx context.Context, cli *client.Client, id string) error {
	copyFromContainer, _, err := cli.CopyFromContainer(ctx, id, "/opt/eap")
	if err != nil {
		return err
	}
	defer copyFromContainer.Close()

	os.Mkdir("eap", 0664)

	tr := tar.NewReader(copyFromContainer)
	var foundFile bool
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
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
		return errors.New("no file found in the archive")
	}

	return nil
}
