package surface

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	HomeVariable = "AGORA_STATE_HOME"
	xdgVariable  = "XDG_STATE_HOME"
	toolName     = "agora"

	sessionFile  = "session.json"
	documentFile = "document.json"
)

func Home() (string, error) {
	if stated := os.Getenv(HomeVariable); stated != "" {
		return stated, nil
	}
	if shared := os.Getenv(xdgVariable); shared != "" {
		return filepath.Join(shared, toolName), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf(noHomeDirectory, err)
	}
	return filepath.Join(home, ".local", "state", toolName), nil
}

func SessionPath() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, sessionFile), nil
}

func DocumentPath() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, documentFile), nil
}
