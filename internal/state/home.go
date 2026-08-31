package state

import (
	"os"
	"path/filepath"

	"github.com/socrates-devtool/socrates/internal/refusal"
)

const (
	HomeVariable      = "SOCRATES_HOME"
	homeDirectory     = ".socrates"
	runsDirectory     = "runs"
	surfaceDirectory  = "surface"
	ConfigurationFile = "config.yaml"
)

func Home() (string, error) {
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
	return filepath.Join(own, homeDirectory), nil
}

func Runs() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, runsDirectory), nil
}

func Surface() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, surfaceDirectory), nil
}

func Configuration() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigurationFile), nil
}
