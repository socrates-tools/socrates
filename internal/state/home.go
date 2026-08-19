package state

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	HomeVariable  = "SOCRATES_STATE_HOME"
	xdgVariable   = "XDG_STATE_HOME"
	runsDirectory = "runs"
)

func Home() (string, error) {
	if stated := os.Getenv(HomeVariable); stated != "" {
		return stated, nil
	}
	if shared := os.Getenv(xdgVariable); shared != "" {
		return filepath.Join(shared, "socrates"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("socrates keeps its runs under your home directory, and it could not find one: %w", err)
	}
	return filepath.Join(home, ".local", "state", "socrates"), nil
}

func Runs() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, runsDirectory), nil
}
