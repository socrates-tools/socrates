package state

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
)

const (
	HomeVariable      = "SOCRATES_HOME"
	brandDirectory    = ".socrates"
	runsDirectory     = "runs"
	surfaceDirectory  = "surface"
	ConfigurationFile = "config.yaml"
)

type Tool string

const (
	Review    Tool = "review"
	Agora     Tool = "agora"
	Implement Tool = "implement"
)

func Root() (string, error) {
	if stated := os.Getenv(HomeVariable); stated != "" {
		if !filepath.IsAbs(stated) {
			return "", refusal.New(statedHomeIsRelative, HomeVariable, stated)
		}
		return stated, nil
	}

	own, err := os.UserHomeDir()
	if err != nil {
		return "", refusal.Wrapping(err, noHomeDirectory, err, HomeVariable)
	}
	return filepath.Join(own, brandDirectory), nil
}

func (t Tool) HomeUnder(root string) string {
	return filepath.Join(root, string(t))
}

func (t Tool) Home() (string, error) {
	root, err := Root()
	if err != nil {
		return "", err
	}
	return t.HomeUnder(root), nil
}

func (t Tool) Configuration() (string, error) {
	home, err := t.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigurationFile), nil
}

func Runs() (string, error) {
	home, err := Review.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, runsDirectory), nil
}

func Surface() (string, error) {
	home, err := Review.Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, surfaceDirectory), nil
}
