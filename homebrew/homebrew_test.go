package homebrew

import (
	"os"
	"testing"
)

func Test_demo1(t *testing.T) {
	t.Log(os.LookupEnv("HOMEBREW_PREFIX"))

	h, err := NewHomebrew(HomebrewConfig{
		// Path:            "/home/linuxbrew/.linuxbrew",
		Path:            "/opt/homebrew",
		InstallScript:   "https://ghp.ci/https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh",
		UninstallScript: "https://ghp.ci/https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh",
		PresetIntro:     "https://docs.brew.sh/Homebrew-on-Linux#requirements",
		Env: HomebrewEnv{
			BrewGitRemote: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/brew.git",
			CoreGitRemote: "https://mirrors.tuna.tsinghua.edu.cn/git/homebrew/homebrew-core.git",
			APIDomain:     "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles/api",
			BottleDomain:  "https://mirrors.tuna.tsinghua.edu.cn/homebrew-bottles",
		},
	})

	if err != nil {
		t.Error(err)
	}
	if err := h.Initialize(); err != nil {
		t.Error(err)
	}
	t.Log(h.Config())

	if err := h.NewPodman(PodmanConfig{
		DockerRegistry: "mirror.ccs.tencentyun.com",
	}); err != nil {
		t.Error(err)
	}
	t.Log(h.CheckServiceStatus("podman", "started"))
}

func Test_demo2(t *testing.T) {
	t.Log(os.UserHomeDir())
	t.Log(os.UserCacheDir())
	t.Log(os.UserConfigDir())
}
