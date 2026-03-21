package container

import (
	"fmt"
	"os"
	"os/exec"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"
)

type EnterContainerOption struct {
	Name  string
	User  string
	Shell string
}

func EnterContainer(param EnterContainerOption) error {

	if param.Shell == "" || param.User == "" {
		ui.Warn("Missing required params to enter")
		return ui.AbortErr
	}

	if !cutils.ImageExist(param.Name) {
		ui.Warn("Container %s does not exists", param.Name)
		return ui.AbortErr
	}

	if !cutils.IsMounted(param.Name) {
		err := MountDiskImage(param.Name)
		if err != nil {
			return err
		}
	}

	mount_path := cutils.GetImageMount(param.Name)

	ui.Step("Entering chroot as %s", param.User)

	script := fmt.Sprintf(`
set -e

mount --make-rprivate /

mount -t proc proc %s/proc
mount -t sysfs sys %s/sys
mount --bind /dev %s/dev
mount -t devpts devpts %s/dev/pts

exec chroot --userspec %s %s %s
`,
		mount_path,
		mount_path,
		mount_path,
		mount_path,
		param.User,
		mount_path,
		param.Shell,
	)

	cmd := exec.Command(
		"unshare",
		"--mount",
		"--pid",
		"--fork",
		"--mount-proc",
		"bash",
		"-c",
		script,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		ui.Warn("Failed to enter disk image")
		return err
	}

	return nil
}
