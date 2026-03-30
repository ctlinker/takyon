package utils

import (
	"os"
	"path/filepath"
	"syscall"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func MountExists(path string) bool {
	var stat, parent syscall.Stat_t

	// stat current path
	if err := syscall.Stat(path, &stat); err != nil {
		return false
	}

	// stat parent path
	parentPath := filepath.Join(path, "..")
	if err := syscall.Stat(parentPath, &parent); err != nil {
		return false
	}

	// if device differs → it's a mount point
	return stat.Dev != parent.Dev
}
