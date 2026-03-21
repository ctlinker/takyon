package container

import (
	"fmt"
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
		if entry.IsDir() {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".img")
		state := "idle"
		mount := GetImageMount(entry.Name())
		size := "unknown"

		if imgInfo, err := entry.Info(); err == nil {
			size = formatSize(imgInfo.Size())
		} else {
			ui.Warn("Failed to retrieve size for %s: %v", name, err)
		}

		if IsMounted(name) {
			state = "mounted"
		} else if IsCorrupted(name) {
			state = "error"
		}

		ui.Step(
			"Container: %-20s State: %-10s Size: %-10s Mount: %s",
			name, state, size, mount,
		)
	}
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
