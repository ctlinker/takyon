package container

import (
	"fmt"
	"os/exec"
	"takyon/lib/ui"
)

type CreateDiskImageOption struct {
	Name   string
	Format string
	Size   int
}

func CreateDiskImage(param CreateDiskImageOption) error {
	if ImageExist(param.Name) {
		return fmt.Errorf("container %s already exists", param.Name)
	}

	if !IsSupportedDiskFormat(param.Format) {
		return fmt.Errorf("unsupported disk format: %s", param.Format)
	}

	imgPath := GetImagePath(param.Name)

	ui.Step("Creating disk image: %s (%dMB, %s)", param.Name, param.Size, param.Format)

	// truncate file
	truncErr := exec.Command("truncate", "-s", fmt.Sprintf("%dM", param.Size), imgPath).Run()
	if truncErr != nil {
		ui.Error("Failed to allocate disk image %s: %v", imgPath, truncErr)
		return truncErr
	}

	ui.Step("Formatting disk image with %s", param.Format)

	// mkfs
	mkfsCmd := fmt.Sprintf("mkfs.%s", param.Format)
	mkfsErr := exec.Command(mkfsCmd, imgPath).Run()
	if mkfsErr != nil {
		ui.Error("Failed to format disk image %s: %v", imgPath, mkfsErr)
		return mkfsErr
	}

	ui.Success("Disk image %s created successfully", param.Name)
	return nil
}
