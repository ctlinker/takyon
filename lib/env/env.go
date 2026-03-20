package env

import (
	"os"
	"path/filepath"
	"takyon/lib/ui"
)

type VAR struct {
	ContainerDirPath   string
	ContainerMountPath string
	SharedDir          string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func ReadEnv() VAR {
	// defaults
	base := "/var/lib/takyon"
	store := filepath.Join(base, ".takyon")
	images := filepath.Join(store, "images")
	mount := "/mnt/takyon"

	return VAR{
		ContainerDirPath:   getEnv("TAKYON_CONTAINER_DIR_PATH", images),
		ContainerMountPath: getEnv("TAKYON_CONTAINER_MOUNT_PATH", mount),
		SharedDir:          getEnv("TAKYON_SHARED_DIR", base),
	}
}

func SetupDir() error {
	env := ReadEnv()

	dirs := []string{
		env.ContainerDirPath,
		env.ContainerMountPath,
		env.SharedDir,
	}

	for _, d := range dirs {
		ui.Step("Ensuring directory exists: %s", d)
		if err := os.MkdirAll(d, 0755); err != nil {
			ui.Error("Failed to create directory %s: %v", d, err)
			return err
		}
	}

	ui.Success("All Takyon directories are ready")
	return nil
}
