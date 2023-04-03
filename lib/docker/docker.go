package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Squwid/bytegolf/lib/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/sirupsen/logrus"
)

type DockerClient struct {
	c *client.Client
}

var Client *DockerClient

var (
	targetArc     = os.Getenv("TARGET_ARCH")
	targetOS      = os.Getenv("TARGET_OS")
	targetVariant = os.Getenv("TARGET_VARIANT") // v8 for rpi.
)

// Init initializes the docker client.
func Init() {
	if targetArc == "" {
		log.GetLogger().Warnf("'TARGET_ARCH' empty, defaulting to 'amd64'")
		targetArc = "amd64"
	}
	if targetOS == "" {
		log.GetLogger().Warnf("'TARGET_OS' empty, defaulting to 'linux'")
		targetOS = "linux"
	}

	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.GetLogger().WithError(err).Fatalf("Error creating docker client")
	}

	Client = &DockerClient{c: c}
}

func (dc *DockerClient) Create(
	image string,
	absHostCodePath string,
	targetFileName string,
	cmd string,
	id string,
	testInputFile string,
	logger *logrus.Entry) (string, io.ReadCloser, error) {
	ctx := context.Background()

	fullCmd := fmt.Sprintf("%s %s", cmd, targetFileName)
	containerBody, err := dc.c.ContainerCreate(ctx, &container.Config{
		OpenStdin:       true,
		Tty:             false,
		AttachStdin:     true,
		AttachStdout:    true,
		Image:           image,
		NetworkDisabled: true,
		Cmd: strslice.StrSlice{
			"/bin/sh", "-c", fullCmd,
		},
		WorkingDir: "/ci",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   absHostCodePath,
				Target:   "/ci/" + targetFileName,
				ReadOnly: true,
			},
			{
				Type:     mount.TypeBind,
				Source:   "/home/bytegolf-inputs/" + testInputFile,
				Target:   "/ci/input.txt",
				ReadOnly: true,
			},
		},

		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"mode": "blocking",
			},
		},

		AutoRemove: true,
	}, &network.NetworkingConfig{}, &v1.Platform{
		Architecture: targetArc,
		OS:           targetOS,
		Variant:      targetVariant,
	}, id)
	if err != nil {
		return "", nil, err
	}

	return dc.start(ctx, containerBody.ID)
}

// Kill terminates the container process but does not remove the container from the docker host.
func (dc *DockerClient) Kill(ctx context.Context, containerID string) error {
	return dc.c.ContainerKill(ctx, containerID, "KILL")
}

// Stats returns the stats of the container. This must be called after the container is done
// running.
func (dc *DockerClient) Stats(ctx context.Context, containerID string) (types.ContainerStats, error) {
	return dc.c.ContainerStats(ctx, containerID, true)
}

// Delete kills and removes a container from the host
func (dc *DockerClient) Delete(ctx context.Context, containerID string) error {
	return dc.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
}

func (dc *DockerClient) Wait(ctx context.Context, containerID string) (int, error) {
	statusCh, errCh := dc.c.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return -1, err
		}
	case status := <-statusCh:
		return int(status.StatusCode), nil
	}
	return -1, nil
}

func (dc *DockerClient) start(ctx context.Context, containerID string) (string, io.ReadCloser, error) {
	if err := dc.c.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return "", nil, err
	}
	reader, err := dc.c.ContainerLogs(context.Background(), containerID, types.ContainerLogsOptions{
		Timestamps: false,
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return "", nil, err
	}

	return containerID, reader, nil
}
