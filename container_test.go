package wgiface_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
)

func Test(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		HostConfigModifier: func(config *container.HostConfig) {
			config.Privileged = true
		},
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "./",
			Dockerfile: "container_test.dockerfile",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}
	defer container.Terminate(ctx)

	exitCode, reader, err := container.Exec(ctx, []string{"go", "test"})
	if err != nil {
		panic(err)
	}
	if exitCode != 0 {
		out := strings.Builder{}
		buf := make([]byte, 1024)
		for {
			n, err := reader.Read(buf)
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				panic(err)
			}
			out.Write(buf[:n])
		}
		fmt.Println(out.String())
		t.Errorf("exit code: %d", exitCode)
	}
}
