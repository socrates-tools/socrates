package snapshot

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	everyWriteBit = 0o222
	ownerWriteBit = 0o200
	probeFile     = ".socrates-probe"
)

var errSealNotEnforced = errors.New(aFilesystemThatEnforcesNothing)

func SealNotEnforced(err error) bool {
	return errors.Is(err, errSealNotEnforced)
}

func Seal(path string) error {
	return adjustPermissions(path, func(mode fs.FileMode) fs.FileMode { return mode &^ everyWriteBit })
}

func Unseal(path string) error {
	return adjustPermissions(path, func(mode fs.FileMode) fs.FileMode { return mode | ownerWriteBit })
}

func (w Worktree) Seal() error {
	return Seal(w.Path)
}

func (w Worktree) Unseal() error {
	return Unseal(w.Path)
}

func (w Worktree) ProveTheSeal() error {
	probe := filepath.Join(w.Path, probeFile)

	written, err := os.OpenFile(probe, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return nil
		}
		return fmt.Errorf(probingTheSeal, probe, err)
	}

	refused := refusal.Wrapping(errSealNotEnforced, theSealDoesNotHold, w.Path)
	if err := written.Close(); err != nil {
		return errors.Join(refused, fmt.Errorf(probingTheSeal, probe, err))
	}
	if err := os.Remove(probe); err != nil {
		return errors.Join(refused, fmt.Errorf(removing, probe, err))
	}
	return refused
}

func adjustPermissions(root string, wanted func(fs.FileMode) fs.FileMode) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && !entry.Type().IsRegular() {
			return nil
		}

		known, err := entry.Info()
		if err != nil {
			return err
		}

		permissions := known.Mode().Perm()
		if wanted(permissions) == permissions {
			return nil
		}
		if err := os.Chmod(path, wanted(permissions)); err != nil {
			return fmt.Errorf(changingThePermissionsOf, path, err)
		}
		return nil
	})
}
