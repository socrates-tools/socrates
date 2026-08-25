package config

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rui-ktei/socrates/internal/agent"
	"github.com/rui-ktei/socrates/internal/refusal"
	"go.yaml.in/yaml/v3"
)

const (
	directoryName = "socrates"
	fileName      = "config.yaml"
)

type Config struct {
	Roles         agent.Configuration `yaml:"roles"`
	DiffCeiling   int                 `yaml:"diff_ceiling"`
	AgentDeadline time.Duration       `yaml:"agent_deadline"`
	LogCeiling    int64               `yaml:"log_ceiling"`
	RunRetention  time.Duration       `yaml:"run_retention"`
}

func Path() (string, error) {
	directory, err := configurationDirectory()
	if err != nil {
		return "", refusal.New(noConfigurationDirectory, err)
	}
	return filepath.Join(directory, directoryName, fileName), nil
}

func configurationDirectory() (string, error) {
	if stated := os.Getenv("XDG_CONFIG_HOME"); filepath.IsAbs(stated) {
		return stated, nil
	}
	return os.UserConfigDir()
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}
	return LoadFrom(path)
}

func LoadFrom(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, refusal.New(unreadableFile, path, err)
	}

	var stated Config
	decoder := yaml.NewDecoder(strings.NewReader(string(content)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&stated); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, refusal.New(unparsableFile, path, err)
	}
	if err := stated.validate(); err != nil {
		return Config{}, refusal.New(statesSomethingUnusable, path, err)
	}
	return stated, nil
}

func (c Config) validate() error {
	for role := range c.Roles {
		if !role.Known() {
			return fmt.Errorf(unknownRole, role, agent.Roles())
		}
	}
	for _, measure := range []struct {
		name  string
		value int64
	}{
		{name: "diff_ceiling", value: int64(c.DiffCeiling)},
		{name: "agent_deadline", value: int64(c.AgentDeadline)},
		{name: "log_ceiling", value: c.LogCeiling},
		{name: "run_retention", value: int64(c.RunRetention)},
	} {
		if measure.value < 0 {
			return fmt.Errorf(negativeMeasure, measure.name)
		}
	}
	return nil
}
