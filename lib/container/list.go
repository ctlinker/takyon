package container

import (
	"os"
	"strings"
	"takyon/lib/env"
	"takyon/lib/ui"
)

func ListContainers() {
	store := env.ReadEnv().ContainerDirPath

	// check if store exists
	info, err := os.Stat(store)
	if err != nil {
		if os.IsNotExist(err) {
			ui.Warn("No containers found: %s does not exist", store)
			return
		}
		ui.Error("Failed to read container store %s: %v", store, err)
		return
	}

	if !info.IsDir() {
		ui.Error("%s exists but is not a directory", store)
		return
	}

	// read directory contents
	entries, err := os.ReadDir(store)
	if err != nil {
		ui.Error("Failed to list container store %s: %v", store, err)
		return
	}

	if len(entries) == 0 {
		ui.Info("No containers found in %s", store)
		return
	}

	ui.Info("Containers in %s:", store)
	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.TrimSuffix(entry.Name(), ".img")
			state := "idle"
			mount := GetImageMount(entry.Name())

			if IsMounted(name) {
				state = "mounted"
			} else if IsCorrupted(name) {
				state = "error"
			}

			// assuming images are files, not directories
			ui.Step("Container: %-20s State: %-10s Mount: %s", name, state, mount)
		}
	}
}
