package homebrew

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/go-playground/validator/v10"
)

type PodmanConfig struct {
	DockerRegistry string `json:"docker_registry" yaml:"docker_registry" validate:"required"`
}

func (h *homebrew) NewPodman(p PodmanConfig) error {
	v := validator.New()
	if err := v.Struct(&p); err != nil {
		return err
	}

	h.podman = NewPodman(p, h)

	if err := h.podman.Fix(); err != nil {
		return err
	}

	if err := h.podman.Optimze(); err != nil {
		return err
	}

	if err := h.podman.Init(); err != nil {
		return err
	}

	return nil
}

type podman struct {
	PodmanConfig
	homebrew *homebrew
	conn     context.Context
	bgOnce   sync.Once
}

func NewPodman(p PodmanConfig, homebrew *homebrew) *podman {
	return &podman{PodmanConfig: p, homebrew: homebrew}
}

func (p *podman) Check() bool {
	pack, err := p.homebrew.Info("podman-compose")
	if err != nil {
		return false
	}
	if pack != nil {
		fmt.Printf("Already install %s:%s\n", pack.Name, pack.Version)
		return true
	}

	return false
}

// 修复 brew service 文件的错误
func (p *podman) Fix() error {
	pack, err := p.homebrew.Info("podman-compose")
	if err != nil {
		return err
	}
	if pack == nil {
		return errors.New("podman not installed")
	}

	podmanServicePath := path.Join(pack.Path, "homebrew.podman.service")
	content, err := os.ReadFile(podmanServicePath)
	if err != nil {
		return fmt.Errorf("read file %s err %v", podmanServicePath, err)
	}
	content = bytes.Replace(content, []byte("--time\\=0"), []byte("--time=0"), 1)
	if err := os.WriteFile(podmanServicePath, content, 0666); err != nil {
		return fmt.Errorf("write file %s err %v", podmanServicePath, err)
	}

	return nil
}

// 替换docker代理为腾讯云代理
func (p *podman) Optimze() error {
	podmanRegistriesPath := path.Join(p.homebrew.Path, "etc/containers/registries.conf")
	f, err := os.OpenFile(podmanRegistriesPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open file %s err %v", podmanRegistriesPath, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("readfull file %s err %v", podmanRegistriesPath, err)
	}

	if bytes.Contains(content, []byte(p.DockerRegistry)) {
		fmt.Printf("already use docker registry %s\n", p.DockerRegistry)
		return nil
	}

	addDockerRegistry := `
[[registry]]
prefix = "docker.io"
location = "%s"
insecure = true
`
	if _, err := f.WriteString(fmt.Sprintf(addDockerRegistry, p.DockerRegistry)); err != nil {
		return fmt.Errorf("write content to file %s err %v", podmanRegistriesPath, err)
	}

	return nil
}

func (p *podman) background() {
	//
}

func (p *podman) Init() error {
	// 启动 podman service
	if err := p.homebrew.StartService("podman"); err != nil {
		return err
	}

	// 获取 podman.sock 路径
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

func (p *podman) ComposeStart(path string) error {
	errStr, err := execCmdWithStdout("podman-compose", "up", "-d")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'podman-compose up -d': %v, errStr: %s\n", err, errStr)
	}

	// todo
	return nil
}

func (p *podman) ComposeStop(path string) error {
	errStr, err := execCmdWithStdout("podman-compose", "down")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'podman-compose down': %v, errStr: %s\n", err, errStr)
	}

	// todo
	return nil
}
