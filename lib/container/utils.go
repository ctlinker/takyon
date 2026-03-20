package container

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"slices"
	"takyon/lib/env"
	"takyon/lib/utils"
)

func ImageExist(containerName string) bool {
	img := GetImagePath(containerName)
	return utils.FileExists(img)
}

func GetImageMount(container_name string) string {
	return filepath.Join(env.ReadEnv().ContainerMountPath, container_name)
}

func GetImagePath(container_name string) string {
	return filepath.Join(env.ReadEnv().ContainerDirPath, fmt.Sprintf("%s.img", container_name))
}

func IsMounted(container_name string) bool {
	mount := GetImageMount(container_name)
	return utils.MountExists(mount)
}

func IsSupportedDiskFormat(input string) bool {
	format := []string{"ext4", "ext3", "ext2"}
	return slices.Contains(format, input)
}

func IsCorrupted(name string) bool {
	img := GetImagePath(name)

	cmd := exec.Command("blkid", img)
	if err := cmd.Run(); err != nil {
		return true
	}

	return false
}
