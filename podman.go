package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
)

type podman struct {
	conn context.Context
}

func (p *podman) Init() error {
	podmanPath, ok := os.LookupEnv("XDG_RUNTIME_DIR")
	if !ok {
		return errors.New("not found XDG_RUNTIME_DIR env, set first please")
	}

	// 链接podman
	podmanSocketPath := "unix://" + podmanPath + "/podman/podman.sock"
	conn, err := bindings.NewConnection(context.Background(), podmanSocketPath)
	if err != nil {
		return fmt.Errorf("connect to %s err %v", podmanSocketPath, err)
	}
	p.conn = conn

	return nil
}

func (p *podman) ListImages() ([]*types.ImageSummary, error) {
	return images.List(p.conn, nil)
}

func (p *podman) ComposeBuild(path string) error {
	return nil
}
