package container

import (
	"fmt"
	"os"
	"os/exec"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"
)

type ResizeDiskImageOption struct {
	Name string
	Size int
	Unit string
}

func ResizeDiskImage(param ResizeDiskImageOption) error {

	if !cutils.ImageExist(param.Name) {
		ui.Warn("No container of this name could be found")
		return ui.AbortErr
	}

	if cutils.IsMounted(param.Name) {
		ui.Warn("Container %s is currently mounted", param.Name)
		return ui.AbortErr
	}

	if !cutils.IsValidSizeUnit(param.Unit) {
		ui.Warn("Invalid image size format: %s", param.Unit)
		return ui.AbortErr
	}

	image_path := cutils.GetImagePath(param.Name)
	info, _ := os.Stat(image_path)

	new_size := fmt.Sprintf("%d%s", param.Size, param.Unit)

	cur_bytes := info.Size()
	new_bytes, _ := cutils.AnySizeTo(new_size, "B")

	if cur_bytes <= new_bytes {
		ui.Step("Preparing to grow disk image %s", param.Name)
		growDiskImage(growDiskImageOption{
			Image: image_path,
			Size:  new_size,
		})
	} else {
		ui.Warn("shrinking an imag is yet to impl")
		return ui.AbortErr
	}

	ui.Success("Successfully resized disk image %s to %s size", param.Name, new_size)
	return nil
}

type growDiskImageOption struct {
	Image string
	Size  string
}

func growDiskImage(param growDiskImageOption) error {

	ui.Step("Re Truncating disk image")
	tr_err := cutils.Run("truncate", "-s", param.Size, param.Image)
	if tr_err != nil {
		ui.Warn("Failed to re truncate image %s to %s", param.Image, param.Size)
	}

	e2_cmd := cutils.MkTTYCommand("e2fsck", "-f", "-p", param.Image)

	ui.Step("Checking disk image filesystem health")
	if err := e2_cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()

			// Acceptable cases
			if code == 1 || code == 2 {
				ui.Warn("Filesystem corrected during check (code %d)", code)
			} else {
				return fmt.Errorf("e2fsck failed with code %d", code)
			}
		} else {
			return err
		}
	}

	ui.Step("Expending disk image filesystem")
	re2fs_err := cutils.Run("resize2fs", param.Image)
	if re2fs_err != nil {
		ui.Warn("Failed to extend disk filesystem; Are you root ?")
		return re2fs_err
	}

	return nil
}
