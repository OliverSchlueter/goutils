package containers

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"log/slog"
)

func StartNATS(ctx context.Context) (string, error) {
	cReq := testcontainers.ContainerRequest{
		Image:        "nats",
		ExposedPorts: []string{"4222/tcp"},
	}
	gReq := testcontainers.GenericContainerRequest{
		ContainerRequest: cReq,
		Started:          true,
		Reuse:            false,
	}

	var err error
	nc, err := testcontainers.GenericContainer(ctx, gReq)
	if err != nil {
		return "", fmt.Errorf("could not start nats container: %w", err)
	}

	port, err := nc.MappedPort(ctx, "4222")
	if err != nil {
		return "", fmt.Errorf("could not get port: %w", err)
	}

	slog.Info("Started NATS test container on port: " + port.Port())

	return port.Port(), nil
}
