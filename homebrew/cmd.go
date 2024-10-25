package homebrew

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
)

const (
	RuntimeOSMac = iota
	RuntimeOSLinux
	RuntimeOSOther
)

func checkOS() int {
	switch os := runtime.GOOS; os {
	case "darwin":
		return RuntimeOSMac
	case "linux":
		return RuntimeOSLinux
	default:
		return RuntimeOSOther
	}
}

func execCmd(name string, arg ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return outStr, errStr, err
	}
	stdout.Reset()
	stderr.Reset()

	return outStr, errStr, nil
}

func execCmdWithStdout(name string, arg ...string) (string, error) {
	var stderr bytes.Buffer

	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout // 标准输出
	cmd.Stderr = &stderr   // 标准错误
	err := cmd.Run()
	errStr := string(stderr.Bytes())
	if err != nil {
		return errStr, err
	}
	stderr.Reset()

	return errStr, nil
}
