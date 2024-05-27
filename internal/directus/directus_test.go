package directus_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Piitschy/drctsdm/internal/directus"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	adminEmail    = "ad@min.de"
	adminPassword = "admin"
)

func NewDirectusContainer(t *testing.T, version string) (context.Context, testcontainers.Container, *directus.Directus) {
	if version == "" {
		version = "latest"
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "directus/directus:latest",
			ExposedPorts: []string{"8055/tcp"},
			WaitingFor:   wait.ForLog("Server started at http://0.0.0.0:8055"),
			Env: map[string]string{
				"ADMIN_EMAIL":    adminEmail,
				"ADMIN_PASSWORD": adminPassword,
			},
		},
		Started: true,
	})

	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to start container: %s", err)
	}
	url, err := container.Host(ctx)
	hostPort, err := container.MappedPort(ctx, "8055")
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to get mapped port or url: %s", err)
	}
	time.Sleep(1 * time.Second)
	d, err := directus.NewDirectus(fmt.Sprintf("http://%s:%s", url, hostPort.Port()), "")
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to create Directus instance: %s", err)
	}
	err = d.TestConnection()
	if err != nil {
		container.Terminate(ctx)
		t.Fatalf("Failed to test connection: %s", err)
	}
	return ctx, container, d
}

func TestNewDirectus(t *testing.T) {
	ctx, container, _ := NewDirectusContainer(t, "latest")
	defer container.Terminate(ctx)
}
