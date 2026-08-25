package surface

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const tokenBytes = 24

type Session struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Process int    `json:"process"`
}

func MintToken() (string, error) {
	held := make([]byte, tokenBytes)
	if _, err := rand.Read(held); err != nil {
		return "", fmt.Errorf(mintingTheToken, err)
	}
	return hex.EncodeToString(held), nil
}

func WriteSession(open Session) error {
	path, err := SessionPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf(writingTheSession, err)
	}

	written, err := json.MarshalIndent(open, "", "  ")
	if err != nil {
		return fmt.Errorf(writingTheSession, err)
	}

	beside, err := os.CreateTemp(filepath.Dir(path), sessionFile+".*")
	if err != nil {
		return fmt.Errorf(writingTheSession, err)
	}
	defer os.Remove(beside.Name())

	if _, err := beside.Write(append(written, '\n')); err != nil {
		beside.Close()
		return fmt.Errorf(writingTheSession, err)
	}
	if err := beside.Close(); err != nil {
		return fmt.Errorf(writingTheSession, err)
	}
	if err := os.Rename(beside.Name(), path); err != nil {
		return fmt.Errorf(writingTheSession, err)
	}
	return nil
}

func ReadSession() (Session, error) {
	path, err := SessionPath()
	if err != nil {
		return Session{}, err
	}

	held, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Session{}, errors.New(noSurfaceIsOpen)
	}
	if err != nil {
		return Session{}, fmt.Errorf(readingTheSession, err)
	}

	var open Session
	if err := json.Unmarshal(held, &open); err != nil {
		return Session{}, fmt.Errorf(readingTheSession, err)
	}
	if open.Address == "" || open.Token == "" {
		return Session{}, errors.New(noSurfaceIsOpen)
	}
	return open, nil
}

func ForgetSession() error {
	path, err := SessionPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(forgettingTheSession, err)
	}
	return nil
}

func WriteDocument(raw []byte) error {
	path, err := DocumentPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf(writingTheDocument, err)
	}
	if err := os.WriteFile(path, raw, 0o600); err != nil {
		return fmt.Errorf(writingTheDocument, err)
	}
	return nil
}
