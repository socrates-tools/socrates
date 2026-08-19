//go:build !unix

package agent

import (
	"errors"
	"os/exec"
)

func ownProcessGroup(command *exec.Cmd) {
}

func killProcessGroup(pid int) error {
	return errors.New("socrates reaps an agent by process group, which this platform does not have")
}
