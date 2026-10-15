package skill

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/state"
	"go.yaml.in/yaml/v3"
)

const RecordFile = "installed.yaml"

type installation struct {
	AgentsDirectory   string            `yaml:"agents_directory,omitempty"`
	AgentsDirectories map[string]string `yaml:"agents_directories,omitempty"`
}

func (i installation) directories() map[string]string {
	held := map[string]string{}
	if i.AgentsDirectory != "" {
		held[skillkit.ClaudeCode.Flag] = i.AgentsDirectory
	}
	for flag, directory := range i.AgentsDirectories {
		held[flag] = directory
	}
	return held
}

func RecordPath() (string, error) {
	home, err := state.Implement.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, RecordFile), nil
}

func Record(host skillkit.Host, agentsDirectory string) error {
	path, err := RecordPath()
	if err != nil {
		return err
	}
	held, err := readRecord(path)
	if err != nil {
		return err
	}
	held[host.Flag] = agentsDirectory

	content, err := yaml.Marshal(installation{AgentsDirectories: held})
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

func Recorded(host skillkit.Host) (string, bool, error) {
	path, err := RecordPath()
	if err != nil {
		return "", false, err
	}
	held, err := readRecord(path)
	if err != nil {
		return "", false, err
	}
	directory, installed := held[host.Flag]
	return directory, installed, nil
}

func readRecord(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, refusal.Wrapping(err, theRecordCannotBeRead, path, err)
	}

	var recorded installation
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&recorded); err != nil {
		return nil, refusal.Wrapping(err, theRecordCannotBeRead, path, err)
	}
	held := recorded.directories()
	if len(held) == 0 {
		err := errors.New(namesNoAgentsDirectory)
		return nil, refusal.Wrapping(err, theRecordCannotBeRead, path, err)
	}
	return held, nil
}
