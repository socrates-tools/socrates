package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
	"go.yaml.in/yaml/v3"
)

const (
	DefaultConcurrency = 4
	DefaultModel       = "sonnet"
	DefaultEffort      = "high"
)

var efforts = []string{"low", "medium", "high", "xhigh", "max"}

func Efforts() []string {
	return slices.Clone(efforts)
}

type Config struct {
	Concurrency *Count `yaml:"concurrency"`
	Roles       Roles  `yaml:"roles"`
}

type Count int

func (c *Count) UnmarshalYAML(node *yaml.Node) error {
	whole, err := strconv.Atoi(node.Value)
	if err != nil {
		return fmt.Errorf(notAWholeNumber, node.Line, node.Value)
	}
	*c = Count(whole)
	return nil
}

type Roles struct {
	Implementer Role `yaml:"implementer"`
}

type Role struct {
	Model  string `yaml:"model"`
	Effort string `yaml:"effort"`
}

func Path() (string, error) {
	return state.Implement.Configuration()
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

func (c Config) validate() error {
	if c.Concurrency != nil && *c.Concurrency < 1 {
		return fmt.Errorf(concurrencyNotPositive, *c.Concurrency)
	}
	if strings.ContainsAny(c.Roles.Implementer.Model, " \t\r\n") {
		return fmt.Errorf(modelHasWhitespace, c.Roles.Implementer.Model)
	}
	if effort := c.Roles.Implementer.Effort; effort != "" && !slices.Contains(efforts, effort) {
		return fmt.Errorf(unknownEffort, effort, efforts)
	}
	return nil
}
