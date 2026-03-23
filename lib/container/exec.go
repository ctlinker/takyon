package container

import (
	"fmt"
	"os"
	"path/filepath"
	"takyon/lib/container/cutils"
	"takyon/lib/ui"
)

type ExecInContainerOption struct {
	Name   string
	User   string
	Script string
}

func ExecInContainer(param ExecInContainerOption) error {
	if param.Script == "" || param.User == "" {
		ui.Warn("Missing required params to enter")
		return ui.AbortErr
	}

	if !cutils.ImageExist(param.Name) {
		ui.Warn("Container %s does not exist", param.Name)
		return ui.AbortErr
	}

	if !cutils.IsMounted(param.Name) {
		if err := MountDiskImage(param.Name); err != nil {
			return err
		}
	}

	mountPath := cutils.GetImageMount(param.Name)

	// --- 🧱 1. Create temp script inside container fs
	scriptPath := filepath.Join(mountPath, "tmp", "takyon_exec.sh")

	scriptContent := fmt.Sprintf(`#!/bin/bash
set -e
%s
`, param.Script)

	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return err
	}

	// ensure cleanup
	defer os.Remove(scriptPath)

	ui.Step("Executing script in container as %s", param.User)

	// --- 🧠 2. Runtime bootstrap script (host-side)
	runtimeScript := fmt.Sprintf(`
set -e

mount --make-rprivate /

mount -t proc proc %[1]s/proc
mount -t sysfs sys %[1]s/sys
mount --bind /dev %[1]s/dev
mount -t devpts devpts %[1]s/dev/pts

exec chroot --userspec %[2]s %[1]s /bin/bash /tmp/takyon_exec.sh
`,
		mountPath,
		param.User,
	)

	cmd := cutils.MkTTYCommand(
		"unshare",
		"--mount",
		"--pid",
		"--fork",
		"--mount-proc",
		"bash",
		"-c",
		runtimeScript,
	)

	if err := cmd.Run(); err != nil {
		ui.Warn("Failed to exec in container")
		return err
	}

	return nil
}
