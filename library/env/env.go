package env

import (
	"os"
	"path/filepath"
)

type VAR struct {
	CONTAINER_STORE_PATH string
	CONTAINER_MOUNT_PATH string
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func ReadEnv() VAR {
	// defaults
	base := filepath.Join("/", "var", "lib", "takyon")
	store := filepath.Join(base, "containers")
	mount := filepath.Join("/", "mnt", "takyon")

	return VAR{
		CONTAINER_STORE_PATH: getEnv("TAKYON_CONTAINER_STORE_PATH", store),
		CONTAINER_MOUNT_PATH: getEnv("TAKYON_CONTAINER_MOUNT_PATH", mount),
	}
}

func SetupEnvDirectories() error {
	env := ReadEnv()

	dirs := []string{
		env.CONTAINER_STORE_PATH,
		env.CONTAINER_MOUNT_PATH,
	}

	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}

	return nil
}
