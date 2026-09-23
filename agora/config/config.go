package config

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Surface Surface `yaml:"surface"`
}

type Surface struct {
	Bind      string `yaml:"bind"`
	ReachedAt string `yaml:"reached_at"`
}

func Path() (string, error) {
	return state.Agora.Configuration()
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
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	decoder.KnownFields(true)
	if err := decoder.Decode(&stated); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, refusal.New(unparsableFile, path, err)
	}
	return stated, nil
}
