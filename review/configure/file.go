package configure

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/config"
)

const (
	partialSuffix = ".partial"
	directoryMode = 0o755
	fileMode      = 0o644
)

func held(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, refusal.New(cannotBeRead, path, err)
	}
	return content, nil
}

func writeAtomically(path string, content []byte) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, directoryMode); err != nil {
		return refusal.New(cannotBeWritten, path, err)
	}

	being, err := os.CreateTemp(directory, filepath.Base(path)+partialSuffix)
	if err != nil {
		return refusal.New(cannotBeWritten, path, err)
	}
	partial := being.Name()
	defer func() { _ = os.Remove(partial) }()

	if _, err := being.Write(content); err != nil {
		_ = being.Close()
		return refusal.New(cannotBeWritten, path, err)
	}
	if err := being.Chmod(fileMode); err != nil {
		_ = being.Close()
		return refusal.New(cannotBeWritten, path, err)
	}
	if err := being.Close(); err != nil {
		return refusal.New(cannotBeWritten, path, err)
	}
	if err := os.Rename(partial, path); err != nil {
		return refusal.New(cannotBeWritten, path, err)
	}
	return nil
}

func withSetting(path, key, value string, removing bool) ([]byte, bool, error) {
	content, err := held(path)
	if err != nil {
		return nil, false, err
	}

	document, err := documentOf(content)
	if err != nil {
		return nil, false, refusal.New(cannotBeRead, path, err)
	}

	if removing {
		if !document.states(key) {
			return nil, false, nil
		}
		document.remove(key)
	} else {
		document.set(key, value)
	}

	written, err := document.render()
	if err != nil {
		return nil, false, refusal.New(cannotBeWritten, path, err)
	}
	if err := config.Parse(path, written); err != nil {
		return nil, false, err
	}
	return written, true, nil
}
