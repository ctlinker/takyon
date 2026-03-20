package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"takyon/lib/ui"
	"takyon/lib/utils"
)

func MountDiskImage(containerName string) error {
	image := GetImagePath(containerName)
	mountpoint := GetImageMount(containerName)

	if !utils.FileExists(image) {
		ui.Error("Container image %s does not exist", containerName)
		return fmt.Errorf("aborting operation")
	}

	if IsMounted(containerName) {
		ui.Warn("Container %s is already mounted at %s", containerName, mountpoint)
		return fmt.Errorf("aborting operation")
	}

	ui.Step("Creating mountpoint directory: %s", mountpoint)
	if err := os.MkdirAll(mountpoint, 0755); err != nil {
		ui.Error("Failed to create mountpoint: %v", err)
		return err
	}

	ui.Step("Mounting container %s", containerName)
	if err := exec.Command("mount", "-o", "loop", image, mountpoint).Run(); err != nil {
		ui.Error("Failed to mount container: %v", err)
		return err
	}

	// create necessary subdirectories
	for _, dir := range []string{"dev", "proc", "sys", "etc"} {
		path := filepath.Join(mountpoint, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			ui.Error("Failed to create directory %s: %v", path, err)
			return err
		}
		ui.Step("Created directory %s", path)
	}

	ui.Success("Container %s mounted at %s", containerName, mountpoint)
	return nil
}
