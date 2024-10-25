package homebrew

import (
	"fmt"
	"strings"
)

type HomebrewServices struct {
	Name   string
	Status string
	User   string
	File   string
}

// StartService homebrew 启动指定 service
func (h *homebrew) StartService(name string) error {
	// 检查是否已经启动
	if h.CheckServiceStatus(name, "started") {
		return nil
	}

	errStr, err := execCmdWithStdout("brew", "services", "start", name)
	if err != nil || len(errStr) > 0 {
		return fmt.Errorf("failed to call 'brew services start %s': %v, errStr: %s\n", name, err, errStr)
	}

	// 检查 brew services podman 是否已启动
	retryCnt := 3
	for !h.CheckServiceStatus(name, "started") && retryCnt > 0 {
		// 尝试启动 podman
		_, err = execCmdWithStdout("brew", "services", "start", name)
		retryCnt--
	}
	if !h.CheckServiceStatus(name, "started") {
		return fmt.Errorf("already try %d times start %s failed, err %v", retryCnt, name, err)
	}

	return nil
}

// ListServices 列出 homebrew 已运行的 service 列表
func (h *homebrew) ListServices() []HomebrewServices {
	services := []HomebrewServices{}

	outStr, errStr, err := execCmd("brew", "services", "list")
	if err != nil || len(errStr) > 0 {
		return nil
	}
	idx := strings.Index(outStr, "Name")
	if idx < 0 {
		return nil
	}

	outStr = outStr[idx:]
	lines := strings.Split(outStr, "\n")
	for _, line := range lines[1:] {
		contents := strings.Fields(strings.TrimSpace(line))
		if len(contents) <= 1 {
			continue
		}
		s := HomebrewServices{
			Name:   contents[0],
			Status: contents[1],
		}
		if len(contents) == 3 { // 可能没有 User 字段，但是一定会有 File
			s.File = contents[2]
		}
		if len(contents) == 4 {
			s.User = contents[2]
			s.File = contents[3]
		}
		services = append(services, s)
	}

	// fmt.Printf("list services %v\n", services)
	return services
}

// CheckServiceStatus 检查指定 service 运行情况
func (h *homebrew) CheckServiceStatus(name string, status string) bool {
	services := h.ListServices()
	for _, service := range services {
		if service.Name == name && service.Status == status {
			return true
		}
	}
	return false
}
