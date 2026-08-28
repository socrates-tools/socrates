package surface

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/rui-ktei/socrates/internal/refusal"
	"go.yaml.in/yaml/v3"
)

const configurationFile = "config.yaml"

type Configuration struct {
	Bind      string `yaml:"bind"`
	ReachedAt string `yaml:"reached_at"`
}

func ConfigurationPath() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configurationFile), nil
}

func LoadConfiguration() (Configuration, error) {
	path, err := ConfigurationPath()
	if err != nil {
		return Configuration{}, err
	}
	return LoadConfigurationFrom(path)
}

func LoadConfigurationFrom(path string) (Configuration, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Configuration{}, nil
	}
	if err != nil {
		return Configuration{}, refusal.New(unreadableConfiguration, path, err)
	}

	var stated Configuration
	decoder := yaml.NewDecoder(strings.NewReader(string(content)))
	decoder.KnownFields(true)
	if err := decoder.Decode(&stated); err != nil && !errors.Is(err, io.EOF) {
		return Configuration{}, refusal.New(unparsableConfiguration, path, err)
	}
	return stated, nil
}
