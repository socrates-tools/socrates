package skill

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"go.yaml.in/yaml/v3"
)

const RecordFile = "installed.yaml"

type installation struct {
	AgentsDirectory string `yaml:"agents_directory"`
}

func RecordPath() (string, error) {
	home, err := state.Implement.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, RecordFile), nil
}

func Record(agentsDirectory string) error {
	path, err := RecordPath()
	if err != nil {
		return err
	}
	content, err := yaml.Marshal(installation{AgentsDirectory: agentsDirectory})
	if err != nil {
		return refusal.Wrapping(err, theRecordCannotBeWritten, path, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return refusal.Wrapping(err, theRecordCannotBeWritten, path, err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return refusal.Wrapping(err, theRecordCannotBeWritten, path, err)
	}
	return nil
}

func Recorded() (string, bool, error) {
	path, err := RecordPath()
	if err != nil {
		return "", false, err
	}
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, refusal.Wrapping(err, theRecordCannotBeRead, path, err)
	}

	var recorded installation
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&recorded); err != nil || recorded.AgentsDirectory == "" {
		if err == nil {
			err = errors.New("it names no agents directory")
		}
		return "", false, refusal.Wrapping(err, theRecordCannotBeRead, path, err)
	}
	return recorded.AgentsDirectory, true, nil
}
