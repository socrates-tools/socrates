package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"github.com/socrates-tools/socrates/review/agent"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	ReviewerCLI     string              `yaml:"reviewer_cli"`
	Roles           agent.Configuration `yaml:"roles"`
	DiffCeiling     int                 `yaml:"diff_ceiling"`
	GuidanceCeiling int                 `yaml:"guidance_ceiling"`
	AgentDeadline   time.Duration       `yaml:"agent_deadline"`
	LogCeiling      int64               `yaml:"log_ceiling"`
	RunRetention    time.Duration       `yaml:"run_retention"`
	Surface         Surface             `yaml:"surface"`
}

type Surface struct {
	Bind      string `yaml:"bind"`
	ReachedAt string `yaml:"reached_at"`
}

func Path() (string, error) {
	return state.Configuration()
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

	return read(path, content)
}

func Parse(path string, content []byte) error {
	_, err := read(path, content)
	return err
}

func read(path string, content []byte) (Config, error) {
	var stated Config
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&stated); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, refusal.New(unparsableFile, path, err)
	}
	if err := stated.validate(); err != nil {
		return Config{}, refusal.New(statesSomethingUnusable, path, err)
	}
	return stated, nil
}

func (c Config) Vendor() string {
	if c.ReviewerCLI == "" {
		return agent.DefaultVendor
	}
	return c.ReviewerCLI
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
		{name: "guidance_ceiling", value: int64(c.GuidanceCeiling)},
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
