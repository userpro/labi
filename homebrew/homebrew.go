package homebrew

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/go-playground/validator/v10"
)

type HomebrewConfig struct {
	Path            string      `json:"path" yaml:"path" validate:"required"` // homebrew 根目录
	InstallScript   string      `json:"install_script" yaml:"install_script" validate:"required"`
	UninstallScript string      `json:"uninstall_script" yaml:"uninstall_script" validate:"required"`
	PresetIntro     string      `json:"preset_intro" yaml:"preset_intro" validate:"required"`
	Env             HomebrewEnv `json:"env" yaml:"env" validate:"required"`
}

type HomebrewEnv struct {
	BrewGitRemote string `json:"brew_git_remote" yaml:"brew_git_remote" validate:"required"`
	CoreGitRemote string `json:"core_git_remote" yaml:"core_git_remote" validate:"required"`
	APIDomain     string `json:"api_domain" yaml:"api_domain" validate:"required"`
	BottleDomain  string `json:"bottle_domain" yaml:"bottle_domain" validate:"required"`
}

type homebrew struct {
	HomebrewConfig
	podman *podman
}

func NewHomebrew(c HomebrewConfig) (*homebrew, error) {
	v := validator.New()
	if err := v.Struct(&c); err != nil {
		return nil, err
	}
	return &homebrew{HomebrewConfig: c}, nil
}

func (h *homebrew) Initialize() error {
	h.setEnv()
	if err := h.InstallSelf(); err != nil {
		return err
	}

	// if err := h.installGCC(); err != nil {
	// 	return err
	// }

	// if err := h.installGo(); err != nil {
	// 	return err
	// }

	if err := h.installPodman(); err != nil {
		return err
	}

	return nil
}

func (h *homebrew) setEnv() {
	os.Setenv("NONINTERACTIVE", "1")
	os.Setenv("HOMEBREW_BREW_GIT_REMOTE", h.Env.BrewGitRemote)
	os.Setenv("HOMEBREW_CORE_GIT_REMOTE", h.Env.CoreGitRemote)
	os.Setenv("HOMEBREW_API_DOMAIN", h.Env.APIDomain)
	os.Setenv("HOMEBREW_BOTTLE_DOMAIN", h.Env.BottleDomain)
}

type homebrewInstalledPackage struct {
	Path    string
	Name    string
	Version string
}

// Info 检查安装情况 如果已安装返回信息
func (h *homebrew) Info(name string) (*homebrewInstalledPackage, error) {
	outStr, errStr, err := execCmd("brew", "info", name)
	if err != nil || len(errStr) > 0 {
		return nil, fmt.Errorf("failed to call 'brew info %s': %v, errStr: %s\n", name, err, errStr)
	}
	lines := strings.Split(outStr, "\n")
	// 安装情况 Installed / Not installed
	// 如果未安装
	if strings.Compare("Installed", strings.TrimSpace(lines[3])) != 0 {
		return nil, nil
	}
	contents := strings.Fields(strings.TrimSpace(lines[4]))

	return &homebrewInstalledPackage{
		Path:    contents[0],
		Name:    name,
		Version: path.Base(contents[0]),
	}, nil
}

// Config 输出 homebrew 完整配置信息
func (h *homebrew) Config() (map[string]string, error) {
	cfg := map[string]string{}

	outStr, errStr, err := execCmd("brew", "config")
	if err != nil || len(errStr) > 0 {
		return nil, fmt.Errorf("failed to call 'brew config': %v, errStr: %s\n", err, errStr)
	}

	lines := strings.Split(outStr, "\n")
	for _, line := range lines[1:] {
		contents := strings.Split(line, ":")
		if len(contents) <= 1 {
			continue
		}

		cfg[strings.TrimSpace(contents[0])] = strings.TrimSpace(contents[1])
	}

	return cfg, nil
}

// Install 安装特定包
func (h *homebrew) Install(name string) error {
	errStr, err := execCmdWithStdout("brew", "install", name)
	if err != nil || len(errStr) > 0 {
		fmt.Printf("failed to call Run(): %v\n", err)
		return fmt.Errorf("failed to call 'brew install %s': %v, errStr: %s\n", name, err, errStr)
	}

	return nil
}

// Uninstall 卸载特定包
func (h *homebrew) Uninstall(name string) error {
	errStr, err := execCmdWithStdout("brew", "remove", name)
	if err != nil || len(errStr) > 0 {
		fmt.Printf("failed to call Run(): %v\n", err)
		return fmt.Errorf("failed to call 'brew remove %s': %v, errStr: %s\n", name, err, errStr)
	}

	return nil
}

// InstallSelf 安装 homebrew 本身
func (h *homebrew) InstallSelf() error {
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

// UninstallSelf 卸载 homebrew 本身
func (h *homebrew) UninstallSelf() error {
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
