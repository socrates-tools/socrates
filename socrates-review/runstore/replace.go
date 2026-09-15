package runstore

import (
	"fmt"
	"os"
	"path/filepath"
)

func replaceFile(path string, content []byte) error {
	directory := filepath.Dir(path)

	temporary, err := os.CreateTemp(directory, filepath.Base(path)+".")
	if err != nil {
		return fmt.Errorf(creatingATemporaryFile, directory, err)
	}
	defer func() { _ = os.Remove(temporary.Name()) }()

	if _, err := temporary.Write(content); err != nil {
		_ = temporary.Close()
		return fmt.Errorf(writing, temporary.Name(), err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		return fmt.Errorf(flushing, temporary.Name(), err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf(closing, temporary.Name(), err)
	}
	if err := os.Rename(temporary.Name(), path); err != nil {
		return fmt.Errorf(renaming, path, err)
	}
	return nil
}
