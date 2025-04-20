package containers

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"log/slog"
)

var natsContainer testcontainers.Container

func StartNATS(ctx context.Context) (string, error) {
	cReq := testcontainers.ContainerRequest{
		Image:        "nats",
		ExposedPorts: []string{"4222/tcp"},
		HostConfigModifier: func(cfg *container.HostConfig) {
			cfg.PortBindings = nat.PortMap{
				"4222/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "4222",
					},
				},
			}
		},
	}
	gReq := testcontainers.GenericContainerRequest{
		ContainerRequest: cReq,
		Started:          true,
		Reuse:            false,
	}

	var err error
	natsContainer, err = testcontainers.GenericContainer(ctx, gReq)
	if err != nil {
		return "", fmt.Errorf("could not start nats container: %w", err)
	}

	port, err := natsContainer.MappedPort(ctx, "4222")
	if err != nil {
		return "", fmt.Errorf("could not get port: %w", err)
	}

	slog.Info("Started NATS test container on port: " + port.Port())

	return port.Port(), nil
}

func StopNATS(ctx context.Context) error {
	err := natsContainer.Terminate(ctx)
	if err != nil {
		return fmt.Errorf("could not stop nats container: %w", err)
	}
	slog.Info("Stopped NATS test container")

	return nil
}
