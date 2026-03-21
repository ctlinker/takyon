package cutils

import (
	"os"
	"os/exec"
)

func Run(name string, arg ...string) error {
	return exec.Command(name, arg...).Run()
}

func RunOnTTY(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func MkCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

func MkTTYCommand(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd
}
