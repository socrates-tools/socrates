package surface

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/socrates-tools/socrates/internal/refusal"
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
	return replace(path, append(written, '\n'), writingTheRecord)
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

func ForgetRecordOf(token string) error {
	if recorded, err := ReadRecord(); err == nil && recorded.Token != token {
		return nil
	}
	return ForgetRecord()
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
	if err := replace(filepath.Join(within, pageFile), append(written, '\n'), writingTheSession); err != nil {
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

func WriteAnswer(session string, answered Answered, staged string) (Answered, error) {
	answers, err := AnswersPath(session)
	if err != nil {
		return Answered{}, err
	}
	if err := os.MkdirAll(answers, 0o700); err != nil {
		return Answered{}, fmt.Errorf(writingTheAnswer, err)
	}

	numbers, err := collectedNumbers(answers)
	if err != nil {
		return Answered{}, fmt.Errorf(writingTheAnswer, err)
	}
	folder, err := AnswerFolderPath(session, nextNumber(numbers))
	if err != nil {
		return Answered{}, err
	}

	kept := answered
	if staged != "" {
		kept = relocated(answered, staged, folder)
	}
	written, err := json.MarshalIndent(kept, "", "  ")
	if err != nil {
		return Answered{}, fmt.Errorf(writingTheAnswer, err)
	}

	if err := settle(staged, folder); err != nil {
		return Answered{}, fmt.Errorf(writingTheAnswer, err)
	}
	if err := replace(filepath.Join(folder, answerFile), append(written, '\n'), writingTheAnswer); err != nil {
		unsettle(staged, folder)
		return Answered{}, err
	}
	return kept, nil
}

func settle(staged, folder string) error {
	if staged == "" {
		return os.Mkdir(folder, 0o700)
	}
	return os.Rename(staged, folder)
}

func unsettle(staged, folder string) {
	if staged == "" || os.Rename(folder, staged) != nil {
		_ = os.RemoveAll(folder)
	}
}

func ReadAnswer(session string) (Answered, bool, error) {
	if filepath.Base(session) != session || !filepath.IsLocal(session) {
		return Answered{}, false, refusal.New(sessionIsNotAName, session)
	}
	answers, err := AnswersPath(session)
	if err != nil {
		return Answered{}, false, err
	}

	numbers, err := collectedNumbers(answers)
	if err != nil {
		return Answered{}, false, fmt.Errorf(readingTheKeptAnswer, err)
	}
	if len(numbers) == 0 {
		return Answered{}, false, nil
	}
	folder, err := AnswerFolderPath(session, slices.Max(numbers))
	if err != nil {
		return Answered{}, false, err
	}

	held, err := os.ReadFile(filepath.Join(folder, answerFile))
	if err != nil {
		return Answered{}, false, fmt.Errorf(readingTheKeptAnswer, err)
	}

	var answered Answered
	if err := json.Unmarshal(held, &answered); err != nil {
		return Answered{}, false, fmt.Errorf(readingTheKeptAnswer, err)
	}
	return answered, true, nil
}

func collectedNumbers(answers string) ([]int, error) {
	entries, err := os.ReadDir(answers)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var numbers []int
	for _, entry := range entries {
		number, err := strconv.Atoi(entry.Name())
		if err == nil && entry.IsDir() && number > 0 {
			numbers = append(numbers, number)
		}
	}
	return numbers, nil
}

func nextNumber(numbers []int) int {
	if len(numbers) == 0 {
		return 1
	}
	return slices.Max(numbers) + 1
}

func replace(path string, written []byte, failure string) error {
	beside, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return fmt.Errorf(failure, err)
	}
	defer func() { _ = os.Remove(beside.Name()) }()

	if _, err := beside.Write(written); err != nil {
		_ = beside.Close()
		return fmt.Errorf(failure, err)
	}
	if err := beside.Close(); err != nil {
		return fmt.Errorf(failure, err)
	}
	if err := os.Rename(beside.Name(), path); err != nil {
		return fmt.Errorf(failure, err)
	}
	return nil
}
