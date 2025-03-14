package testdb

import (
	"context"
	"fmt"
	"net"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDBContainer struct {
	testcontainers.Container
	URI string
}

type TestLogConsumer struct{}

func (g TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Println(l.LogType, string(l.Content))
}

var logConsumer TestLogConsumer

const (
	db       = "test_db"
	user     = "test_user"
	password = "test_password"
)

func New(ctx context.Context) (*TestDBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:17.2",
		ExposedPorts: []string{"5432/tcp", "8080/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       db,
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForLog(`listening on IPv4 address "0.0.0.0", port 5432`),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("testdb New Generic Container failed: %w", err)
	}

	err = container.StartLogProducer(ctx)
	if err != nil {
		return nil, fmt.Errorf("testdb New StartLogProducer failed: %w", err)
	}
	container.FollowOutput(logConsumer)

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("testdb New MappedPort failed: %w", err)
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("testdb New Host failed: %w", err)
	}

	host := net.JoinHostPort(hostIP, mappedPort.Port())
	uri := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, db)

	return &TestDBContainer{
		Container: container,
		URI:       uri,
	}, nil
}

func (t *TestDBContainer) Close(ctx context.Context) error {
	err := t.Container.Terminate(ctx)
	if err != nil {
		return fmt.Errorf("testdb close failed: %w", err)
	}
	return nil
}
