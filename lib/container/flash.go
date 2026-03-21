package container

import (
	"fmt"
	"os"
	"os/exec"
	"takyon/lib/ui"
)

func FlashContainer(containerName, distro string) error {
	mountpoint := GetImageMount(containerName)

	ui.Step("Bootstrapping container %s (%s)", containerName, distro)

	if !IsMounted(containerName) {
		ui.Error("Constainer %s is not mounted", containerName)
		return fmt.Errorf("aborting operation")
	}

	var cmd *exec.Cmd
	switch distro {
	case "kali":
		cmd = exec.Command(
			"debootstrap",
			"--arch=amd64",
			"--variant=minbase",
			"kali-rolling",
			mountpoint,
			"http://http.kali.org/kali",
		)
	case "debian", "ubuntu":
		cmd = exec.Command("debootstrap", "--arch=amd64", "bookworm", mountpoint, "http://deb.debian.org/debian")
	case "minimal-arch":
		cmd = exec.Command("pacstrap", mountpoint, "base")
	case "arch":
		cmd = exec.Command("pacstrap", mountpoint, "base bash coreutils fastfetch nano git sudo")
	default:
		return fmt.Errorf("unsupported distro: %s", distro)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.Error("Failed to flash container: %v", err)
		return err
	}

	ui.Success("Container %s bootstrapped successfully", containerName)
	return nil
}
