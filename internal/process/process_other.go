//go:build !unix

package process

import (
	"errors"
	"os"
	"os/exec"
)

var unsupported = errors.New("socrates reaps by process group, which this platform does not have")

func OwnGroup(command *exec.Cmd) {
}

func Detach(command *exec.Cmd) {
}

func KillGroup(pid int) error {
	return unsupported
}

func Kill(pid int) error {
	found, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return found.Kill()
}

func Alive(pid int) bool {
	found, err := os.FindProcess(pid)
	return err == nil && found != nil
}
