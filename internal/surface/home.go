package surface

import (
	"os"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
)

const (
	RetiredHomeVariable = "AGORA_HOME"
	ReachedAtField      = "surface.reached_at"

	recordFile        = "surface.json"
	sessionsDirectory = "sessions"
	pageFile          = "session.json"
	documentFile      = "document.json"
)

func Home() (string, error) {
	return state.Surface()
}

func CheckNoHomeOfItsOwn() error {
	if os.Getenv(RetiredHomeVariable) == "" {
		return nil
	}
	return refusal.New(surfaceHasNoHomeOfItsOwn, RetiredHomeVariable, state.HomeVariable, RetiredHomeVariable)
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
