//go:build unix

package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func OwnGroup(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func Detach(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
}

func KillGroup(pid int) error {
	if pid <= 0 {
		return nil
	}
	if err := syscall.Kill(-pid, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
		return fmt.Errorf("killing the process group %d: %w", pid, err)
	}
	return nil
}

func Kill(pid int) error {
	if pid <= 0 {
		return nil
	}
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil && err != syscall.ESRCH {
		return fmt.Errorf("killing the process %d: %w", pid, err)
	}
	return nil
}

func Alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	if err := syscall.Kill(pid, 0); err != nil {
		return false
	}
	return !reaped(pid)
}

func reaped(pid int) bool {
	stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return false
	}
	closing := strings.LastIndexByte(string(stat), ')')
	if closing < 0 {
		return false
	}
	fields := strings.Fields(string(stat)[closing+1:])
	return len(fields) > 0 && fields[0] == "Z"
}
