package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

type HomebrewConfig struct {
	Path            string            `json:"path" yaml:"path" validate:"required"` // homebrew 根目录
	InstallScript   string            `json:"install_script" yaml:"install_script" validate:"required"`
	UninstallScript string            `json:"uninstall_script" yaml:"uninstall_script" validate:"required"`
	PresetIntro     string            `json:"preset_intro" yaml:"preset_intro" validate:"required"`
	Env             HomebrewEnvConfig `json:"env" yaml:"env" validate:"required"`
}

type HomebrewEnvConfig struct {
	BrewGitRemote string `json:"brew_git_remote" yaml:"brew_git_remote" validate:"required"`
	CoreGitRemote string `json:"core_git_remote" yaml:"core_git_remote" validate:"required"`
	APIDomain     string `json:"api_domain" yaml:"api_domain" validate:"required"`
	BottleDomain  string `json:"bottle_domain" yaml:"bottle_domain" validate:"required"`
}

type homebrew struct {
	HomebrewConfig
}

func NewHomebrew(c HomebrewConfig) (*homebrew, error) {
	v := validator.New()
	if err := v.Struct(&c); err != nil {
		return nil, err
	}
	return &homebrew{c}, nil
}

func (h *homebrew) Init() {
	os.Setenv("NONINTERACTIVE", "1")
	os.Setenv("HOMEBREW_BREW_GIT_REMOTE", "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git")
	os.Setenv("HOMEBREW_CORE_GIT_REMOTE", "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git")
	os.Setenv("HOMEBREW_API_DOMAIN", "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api")
	os.Setenv("HOMEBREW_BOTTLE_DOMAIN", "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles")
}

func (h *homebrew) Install() error {
	if _, err := os.Stat(h.Path); err == nil {
		fmt.Println("Already install Homebrew")
		return nil
	}

	// 下载 install 脚本
	outStr, errStr, err := execCmd("curl", "-fsSL", h.InstallScript)
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'curl -fsSL %s': %v, errStr: %s, check %s\n", h.InstallScript, err, errStr, h.PresetIntro)
	}

	// 执行脚本
	errStr, err = execCmdWithStdout("/bin/bash", "-c", outStr)
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call '/bin/bash -c install.sh': %v, errStr: %s, check %s\n", err, errStr, h.PresetIntro)
	}

	return nil
}

func (h *homebrew) Uninstall() error {
	if _, err := os.Stat(h.Path); err != nil {
		fmt.Println("Homebrew is not installed")
		return nil
	}

	// 下载 uninstall 脚本
	outStr, errStr, err := execCmd("curl", "-fsSL", h.UninstallScript)
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'curl -fsSL %s': %v, errStr: %s\n", h.UninstallScript, err, errStr)
	}

	// 执行脚本
	errStr, err = execCmdWithStdout("/bin/bash", "-c", outStr)
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call '/bin/bash -c uninstall.sh': %v, errStr: %s\n", err, errStr)
	}

	return nil
}

func (h *homebrew) InstallGCC() error {
	outStr, errStr, err := execCmd("gcc", "--version")
	if err != nil || len(errStr) > 0 {
		fmt.Printf("failed to call Run(): %v\n", err)
		return fmt.Errorf("failed to call 'brew install gcc': %v, errStr: %s\n", err, errStr)
	}
	if strings.Contains(outStr, "Copyright") {
		fmt.Printf("Already install %s\n", strings.Split(outStr, "\n")[0])
		return nil
	}

	outStr, errStr, err = execCmd("brew", "install", "gcc")
	if err != nil || len(errStr) > 0 {
		fmt.Printf("failed to call Run(): %v\n", err)
		return fmt.Errorf("failed to call 'brew install gcc': %v, errStr: %s\n", err, errStr)
	}

	return nil
}

func (h *homebrew) InstallGo() error {
	outStr, errStr, err := execCmd("go", "version")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'go version': %v, errStr: %s\n", err, errStr)
	}
	if strings.Contains(outStr, "go version") {
		fmt.Printf("Already install %s\n", strings.Split(outStr, " ")[2])
		return nil
	}

	errStr, err = execCmdWithStdout("brew", "install", "go")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'brew install go': %v, errStr: %s\n", err, errStr)
	}

	// 注入 go proxy env
	errStr, err = execCmdWithStdout("go", "env", "-w", "GO111MODULE=on")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'go env -w GO111MODULE=on': %v, errStr: %s\n", err, errStr)
	}
	errStr, err = execCmdWithStdout("go", "env", "-w", "GOPROXY=https://goproxy.cn,direct")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'go env -w GOPROXY=https://goproxy.cn,direct': %v, errStr: %s\n", err, errStr)
	}

	return nil
}

func (h *homebrew) InstallPodman() error {
	outStr, errStr, err := execCmd("podman-compose", "--version")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'podman-compose --version': %v, errStr: %s\n", err, errStr)
	}
	if strings.Contains(outStr, "podman-compose version") {
		fmt.Printf("Already install %s\n", outStr)
		return nil
	}

	errStr, err = execCmdWithStdout("brew", "install", "podman", "podman-compose")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'brew install podman podman-compose': %v, errStr: %s\n", err, errStr)
	}

	return nil
}
