package container

import (
	"os"
	"path/filepath"
	"takyon/library/env"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Registry map[string]RegistryEntry

type RegistryEntry struct {
	Name      string `yaml:"name"`
	CreatedAt string `yaml:"created_at"`
}

func (r Registry) GetById(id string) (*RegistryEntry, bool) {
	e, ok := r[id]
	return &e, ok
}

func (r Registry) GetByName(name string) (*RegistryEntry, bool) {
	for _, e := range r {
		if e.Name == name {
			return &e, true
		}
	}
	return nil, false
}

func LoadContainerRegistry() (Registry, error) {
	path := getRegistryPath()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var reg Registry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}

	return reg, nil
}

func SaveContainerRegistry(reg Registry) error {
	path := getRegistryPath()

	data, err := yaml.Marshal(reg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func GenRegistryEntryId() uuid.UUID {
	return uuid.New()
}

func getRegistryPath() string {
	store := env.ReadEnv().CONTAINER_STORE_PATH
	path := filepath.Join(store, env.CONTAINER_STORE_REGISTRY_FILE_NAME)
	return path
}
