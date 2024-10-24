package main

import (
	"fmt"
	"os"
	"testing"
)

func Test_demo1(t *testing.T) {
	fmt.Println(os.LookupEnv("HOMEBREW_PREFIX"))

	h, err := NewHomebrew(HomebrewConfig{
		Path:            "/home/linuxbrew/.linuxbrew",
		InstallScript:   "https://ghp.ci/https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh",
		UninstallScript: "https://ghp.ci/https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh",
		PresetIntro:     "https://docs.brew.sh/Homebrew-on-Linux#requirements",
		Env: HomebrewEnvConfig{
			BrewGitRemote: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git",
			CoreGitRemote: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git",
			APIDomain:     "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api",
			BottleDomain:  "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles",
		},
	})
	if err != nil {
		t.Error(err)
	}
	h.Init()
	if err := h.Install(); err != nil {
		t.Error(err)
	}
	if err := h.InstallGCC(); err != nil {
		t.Error(err)
	}
	if err := h.InstallGo(); err != nil {
		t.Error(err)
	}
	if err := h.InstallPodman(); err != nil {
		t.Error(err)
	}
}
