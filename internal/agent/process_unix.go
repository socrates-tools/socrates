//go:build unix

package agent

import (
	"fmt"
	"os/exec"
	"syscall"
)

func ownProcessGroup(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessGroup(pid int) error {
	if pid <= 0 {
		return nil
	}
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
		return fmt.Errorf("killing the process group %d: %w", pid, err)
	}
	return nil
}
