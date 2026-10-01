package settingsmenu

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	partialSuffix = ".partial"
	directoryMode = 0o755
	fileMode      = 0o644
)

func (d Definition) held(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, refusal.New(cannotBeRead, d.MenuCommand, path, err)
	}
	return content, nil
}

func (d Definition) writeAtomically(path string, content []byte) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, directoryMode); err != nil {
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}

	being, err := os.CreateTemp(directory, filepath.Base(path)+partialSuffix)
	if err != nil {
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	partial := being.Name()
	defer func() { _ = os.Remove(partial) }()

	if _, err := being.Write(content); err != nil {
		_ = being.Close()
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	if err := being.Chmod(fileMode); err != nil {
		_ = being.Close()
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	if err := being.Close(); err != nil {
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	if err := os.Rename(partial, path); err != nil {
		return refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	return nil
}

func (d Definition) withSetting(path string, key []string, value string, removing bool) ([]byte, bool, error) {
	content, err := d.held(path)
	if err != nil {
		return nil, false, err
	}

	document, err := documentOf(content)
	if err != nil {
		return nil, false, refusal.New(cannotBeRead, d.MenuCommand, path, err)
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
		return nil, false, refusal.New(cannotBeWritten, d.MenuCommand, path, err)
	}
	if err := d.Parse(path, written); err != nil {
		return nil, false, err
	}
	return written, true, nil
}

func (d Definition) WriteKey(key []string, value string) error {
	return d.changeKey(key, value, false)
}

func (d Definition) RemoveKey(key []string) error {
	return d.changeKey(key, "", true)
}

func (d Definition) changeKey(key []string, value string, removing bool) error {
	written, changed, err := d.withSetting(d.Path, key, value, removing)
	if err != nil || !changed {
		return err
	}
	return d.writeAtomically(d.Path, written)
}
