package docker

import (
	"context"
	"fmt"
	"io"
	"os"

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
	targetArc = os.Getenv("TARGET_ARCH")
	targetOS  = os.Getenv("TARGET_OS")
)

func init() {
	if targetArc == "" {
		logrus.Warnf("'TARGET_ARCH' empty, defaulting to 'amd64'")
		targetArc = "amd64"
	}
	if targetOS == "" {
		logrus.Warnf("'TARGET_OS' empty, defaulting to 'linux'")
		targetOS = "linux"
	}

	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logrus.WithError(err).Fatalf("Error creating docker client")
	}

	Client = &DockerClient{c: c}
}

func (dc *DockerClient) Create(image, dir, cmd, file, id string,
	logger *logrus.Entry) (string, io.ReadCloser, error) {
	ctx := context.Background()
	var timeout = 10

	// TODO: echo is here to flush the logs to stdout.
	// This is a hack and should be removed.
	fullCmd := fmt.Sprintf("%s %s;echo ''", cmd, file)
	logger.WithField("Cmd", fullCmd).WithField("Dir", dir).Debugf("Container create")

	containerBody, err := dc.c.ContainerCreate(ctx, &container.Config{
		OpenStdin:       true,
		Tty:             true,
		AttachStdin:     true,
		AttachStdout:    true,
		Image:           image,
		NetworkDisabled: true,
		StopTimeout:     &timeout,
		Cmd: strslice.StrSlice{
			"/bin/sh", "-c", fullCmd,
		},
		WorkingDir: "/ci",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: dir,
				Target: "/ci",
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
		// TODO: This is RPI thing
		// Variant:      "v8",
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

// Delete kills and removes a container from the host
func (dc *DockerClient) Delete(ctx context.Context, containerID string) error {
	return dc.c.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})
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
