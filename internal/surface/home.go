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

	recordFile        = "surface.json"
	sessionsDirectory = "sessions"
	pageFile          = "session.json"
	documentFile      = "document.json"
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

func RecordPath() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, recordFile), nil
}

func SessionsPath() (string, error) {
	home, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, sessionsDirectory), nil
}

func SessionPath(session string) (string, error) {
	sessions, err := SessionsPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(sessions, session), nil
}

func PagePath(session string) (string, error) {
	within, err := SessionPath(session)
	if err != nil {
		return "", err
	}
	return filepath.Join(within, pageFile), nil
}

func DocumentPath(session string) (string, error) {
	within, err := SessionPath(session)
	if err != nil {
		return "", err
	}
	return filepath.Join(within, documentFile), nil
}
