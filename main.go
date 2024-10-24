package main

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/specgen"
)

func main() {
	podmanPath, ok := os.LookupEnv("XDG_RUNTIME_DIR")
	if !ok {
		fmt.Println("not found XDG_RUNTIME_DIR env, set first please")
		os.Exit(1)
	}

	// 链接podman
	conn, err := bindings.NewConnection(context.Background(), "unix://"+podmanPath+"/podman/podman.sock")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 列出镜像
	imgs, err := images.List(conn, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, c := range imgs {
		fmt.Printf("[%d] image %s\n", i, c.Names)
	}

	imgName := "quay.io/libpod/alpine_nginx"
	// 查找镜像
	img, err := images.GetImage(conn, imgName, nil)
	if err != nil {
		fmt.Println(err)
		// 下载镜像
		imgIDs, err := images.Pull(conn, imgName, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(imgIDs)
	}
	fmt.Printf("GetImage ID %s, name %v\n", img.ID, img.NamesHistory)

	// 查找容器
	containerName := "foobar"
	inspectData, err := containers.Inspect(conn, containerName, new(containers.InspectOptions).WithSize(true))
	if err != nil {
		fmt.Println(err)
		// 使用镜像创建容器
		s := specgen.NewSpecGenerator(imgName, false)
		s.Name = containerName
		s.OCIRuntime = "crun"
		_, err := containers.CreateWithSpec(conn, s, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Container created.")

		inspectData, err = containers.Inspect(conn, containerName, new(containers.InspectOptions).WithSize(true))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	containerID := inspectData.ID
	fmt.Printf("inspect container id: %s\n", containerID)

	if !inspectData.State.Running {
		// 启动容器
		if err := containers.Start(conn, containerID, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 等待容器启动完成
		_, err = containers.Wait(conn, containerID, new(containers.WaitOptions).WithCondition([]define.ContainerStatus{define.ContainerStateRunning}))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Container %s started.\n", containerID)
	}

	// 停止容器
	err = containers.Stop(conn, containerID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 等待容器停止完成
	_, err = containers.Wait(conn, containerID, new(containers.WaitOptions).WithCondition([]define.ContainerStatus{define.ContainerStateStopped}))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Container %s stopped.\n", containerID)

	// 删除容器
	_, err = containers.Remove(conn, containerID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Container %s removed.\n", containerID)

	// 列出容器
	conts, err := containers.List(conn, new(containers.ListOptions).WithAll(true))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, c := range conts {
		fmt.Printf("[%d] container %v\n", i, c.Names)
	}
}
