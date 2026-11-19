//go:build unix

package runstore

import (
	"fmt"
	"os"
	"syscall"
)

func holdingTheRecord(directory string, do func() error) error {
	held, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf(locking, directory, err)
	}
	defer func() { _ = held.Close() }()

	if err := syscall.Flock(int(held.Fd()), syscall.LOCK_EX); err != nil {
		return fmt.Errorf(locking, directory, err)
	}
	defer func() { _ = syscall.Flock(int(held.Fd()), syscall.LOCK_UN) }()

	return do()
}
