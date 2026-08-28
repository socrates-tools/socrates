package surface

import (
	"os"
	"path/filepath"

	"github.com/rui-ktei/socrates/internal/refusal"
)

const (
	HomeVariable  = "AGORA_HOME"
	toolName      = "agora"
	homeDirectory = "." + toolName

	sessionFile  = "session.json"
	documentFile = "document.json"
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
