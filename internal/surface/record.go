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

const (
	tokenBytes   = 24
	sessionBytes = 8
)

type Record struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Process int    `json:"process"`
}

type Page struct {
	Session string `json:"session"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

func MintToken() (string, error) {
	held := make([]byte, tokenBytes)
	if _, err := rand.Read(held); err != nil {
		return "", fmt.Errorf(mintingTheToken, err)
	}
	return hex.EncodeToString(held), nil
}

func MintSession() (string, error) {
	held := make([]byte, sessionBytes)
	if _, err := rand.Read(held); err != nil {
		return "", fmt.Errorf(mintingTheSession, err)
	}
	return hex.EncodeToString(held), nil
}

func WriteRecord(open Record) error {
	path, err := RecordPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf(writingTheRecord, err)
	}

	written, err := json.MarshalIndent(open, "", "  ")
	if err != nil {
		return fmt.Errorf(writingTheRecord, err)
	}
	return replace(path, append(written, '\n'))
}

func ReadRecord() (Record, error) {
	path, err := RecordPath()
	if err != nil {
		return Record{}, err
	}

	held, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Record{}, errors.New(noSurfaceIsOpen)
	}
	if err != nil {
		return Record{}, fmt.Errorf(readingTheRecord, err)
	}

	var open Record
	if err := json.Unmarshal(held, &open); err != nil {
		return Record{}, fmt.Errorf(readingTheRecord, err)
	}
	if open.Address == "" || open.Token == "" {
		return Record{}, errors.New(noSurfaceIsOpen)
	}
	return open, nil
}

func ForgetRecord() error {
	path, err := RecordPath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(forgettingTheRecord, err)
	}
	return nil
}

func WritePage(open Page, raw []byte) error {
	within, err := SessionPath(open.Session)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(within, 0o700); err != nil {
		return fmt.Errorf(writingTheSession, err)
	}

	written, err := json.MarshalIndent(open, "", "  ")
	if err != nil {
		return fmt.Errorf(writingTheSession, err)
	}
	if err := replace(filepath.Join(within, pageFile), append(written, '\n')); err != nil {
		return err
	}
	return WriteDocument(open.Session, raw)
}

func WriteDocument(session string, raw []byte) error {
	path, err := DocumentPath(session)
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

func ForgetPage(session string) error {
	within, err := SessionPath(session)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(within); err != nil {
		return fmt.Errorf(forgettingTheSession, err)
	}
	return nil
}

func replace(path string, written []byte) error {
	beside, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return fmt.Errorf(writingTheRecord, err)
	}
	defer func() { _ = os.Remove(beside.Name()) }()

	if _, err := beside.Write(written); err != nil {
		_ = beside.Close()
		return fmt.Errorf(writingTheRecord, err)
	}
	if err := beside.Close(); err != nil {
		return fmt.Errorf(writingTheRecord, err)
	}
	if err := os.Rename(beside.Name(), path); err != nil {
		return fmt.Errorf(writingTheRecord, err)
	}
	return nil
}
