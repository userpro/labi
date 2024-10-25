package homebrew

import (
	"fmt"
)

func (h *homebrew) installGCC() error {
	// 检查是否已经安装
	p, err := h.Info("gcc")
	if err != nil {
		return err
	}
	if p != nil {
		fmt.Printf("Already install %s:%s\n", p.Name, p.Version)
		return nil
	}

	// 安装
	if err := h.Install("gcc"); err != nil {
		return err
	}

	return nil
}

func (h *homebrew) installGo() error {
	// 检查是否已经安装
	p, err := h.Info("go")
	if err != nil {
		return err
	}
	if p != nil {
		fmt.Printf("Already install %s:%s\n", p.Name, p.Version)
		return nil
	}

	// 安装
	if err := h.Install("go"); err != nil {
		return err
	}

	// 注入 go proxy env
	errStr, err := execCmdWithStdout("go", "env", "-w", "GO111MODULE=on")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'go env -w GO111MODULE=on': %v, errStr: %s\n", err, errStr)
	}
	errStr, err = execCmdWithStdout("go", "env", "-w", "GOPROXY=https://goproxy.cn,direct")
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'go env -w GOPROXY=https://goproxy.cn,direct': %v, errStr: %s\n", err, errStr)
	}

	return nil
}

func (h *homebrew) installPodman() error {
	// 检查是否已经安装
	p, err := h.Info("podman-compose")
	if err != nil {
		return err
	}
	if p != nil {
		fmt.Printf("Already install %s:%s\n", p.Name, p.Version)
		return nil
	}

	// 安装
	if err := h.Install("podman"); err != nil {
		return err
	}

	if err := h.Install("podman-compose"); err != nil {
		return err
	}

	return nil
}
