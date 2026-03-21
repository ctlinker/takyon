package container

import (
	"fmt"
	"os"
	"path/filepath"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"
	"takyon/lib/utils"
)

func MountDiskImage(containerName string) error {
	image := cutils.GetImagePath(containerName)
	mountpoint := cutils.GetImageMount(containerName)

	if !utils.FileExists(image) {
		ui.Error("Container image %s does not exist", containerName)
		return fmt.Errorf("aborting operation")
	}

	if cutils.IsMounted(containerName) {
		ui.Warn("Container %s is already mounted at %s", containerName, mountpoint)
		return fmt.Errorf("aborting operation")
	}

	ui.Step("Creating mountpoint directory: %s", mountpoint)
	if err := os.MkdirAll(mountpoint, 0755); err != nil {
		ui.Error("Failed to create mountpoint: %v", err)
		return err
	}

	ui.Step("Mounting container %s", containerName)
	if err := cutils.Run("mount", "-o", "loop", image, mountpoint); err != nil {
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

	// adding host resolv.conf
	host_resolv := filepath.Join("/", "etc", "resolv.conf")
	img_resolv := filepath.Join(mountpoint, "etc", "resolv.conf")

	if !utils.FileExists(img_resolv) && utils.FileExists(host_resolv) {
		ui.Step("Inheriting host's /etc/resolv.conf file")
		err := cutils.Run("cp", host_resolv, img_resolv)
		if err != nil {
			ui.Warn("Failed to inherit %s, manualy write it to your image or you may not have access to internet", host_resolv)
		}
	}

	ui.Success("Container %s mounted at %s", containerName, mountpoint)
	return nil
}

func UmountDiskImage(containerName string) error {
	mountpoint := cutils.GetImageMount(containerName)

	if !cutils.IsMounted(containerName) {
		ui.Warn("Container %s is not mounted at", containerName)
		return ui.AbortErr
	}

	ui.Step("Umounting mountpoint directory: %s", mountpoint)
	if err := cutils.Run("umount", "-R", mountpoint); err != nil {
		ui.Error("Failed to umount mountpoint: %v", mountpoint)
		return err
	}

	ui.Step("Removing mountpoint directory: %s", mountpoint)
	if rm_err := os.RemoveAll(mountpoint); rm_err != nil {
		ui.Warn("Failed removet mountpoint directory: %v", mountpoint)
		return rm_err
	}

	ui.Success("Container %s umounted successfully", containerName)
	return nil
}
