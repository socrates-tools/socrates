package runstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/socrates-tools/socrates/internal/dialogue"
	"github.com/socrates-tools/socrates/internal/findings"
	"github.com/socrates-tools/socrates/internal/hunks"
)

const (
	roundPrefix   = "r"
	promptFile    = "prompt.txt"
	streamFile    = "stream.jsonl"
	findingsFile  = "findings.json"
	changeFile    = "change.json"
	responsesFile = "responses.jsonl"
)

type Round struct {
	Number    int
	Directory string
}

func (s Store) Round(name string, number int) Round {
	return Round{
		Number:    number,
		Directory: filepath.Join(s.Path(name), roundPrefix+strconv.Itoa(number)),
	}
}

func (s Store) CreateRound(name string, number int) (Round, error) {
	round := s.Round(name, number)
	if err := os.MkdirAll(round.Directory, 0o755); err != nil {
		return Round{}, fmt.Errorf(creating, round.Directory, err)
	}
	return round, nil
}

func (s Store) Rounds(name string) ([]int, error) {
	directory := s.Path(name)

	entries, err := os.ReadDir(directory)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf(reading, directory, err)
	}

	var numbers []int
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		digits, prefixed := strings.CutPrefix(entry.Name(), roundPrefix)
		if !prefixed {
			continue
		}
		number, err := strconv.Atoi(digits)
		if err != nil {
			continue
		}
		numbers = append(numbers, number)
	}

	slices.Sort(numbers)
	return numbers, nil
}

func (s Store) LatestRound(name string) (Round, bool, error) {
	numbers, err := s.Rounds(name)
	if err != nil || len(numbers) == 0 {
		return Round{}, false, err
	}
	return s.Round(name, numbers[len(numbers)-1]), true, nil
}

func (r Round) PromptPath() string {
	return filepath.Join(r.Directory, promptFile)
}

func (r Round) StreamPath() string {
	return filepath.Join(r.Directory, streamFile)
}

func (r Round) FindingsPath() string {
	return filepath.Join(r.Directory, findingsFile)
}

func (r Round) Write(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf(writing, path, err)
	}
	return nil
}

func (r Round) ChangePath() string {
	return filepath.Join(r.Directory, changeFile)
}

func (r Round) WriteChange(change hunks.Change) error {
	encoded, err := json.MarshalIndent(change, "", "  ")
	if err != nil {
		return fmt.Errorf(encodingTheChange, r.ChangePath(), err)
	}
	return r.Write(r.ChangePath(), string(encoded)+"\n")
}

func (r Round) Change() (hunks.Change, bool, error) {
	content, err := os.ReadFile(r.ChangePath())
	if errors.Is(err, fs.ErrNotExist) {
		return hunks.Change{}, false, nil
	}
	if err != nil {
		return hunks.Change{}, false, fmt.Errorf(reading, r.ChangePath(), err)
	}

	var read hunks.Change
	if err := json.Unmarshal(content, &read); err != nil {
		return hunks.Change{}, false, fmt.Errorf(parsing, r.ChangePath(), err)
	}
	return read, true, nil
}

func (r Round) ResponsesPath() string {
	return filepath.Join(r.Directory, responsesFile)
}

func (r Round) Responses() (dialogue.Responses, error) {
	content, err := os.ReadFile(r.ResponsesPath())
	if errors.Is(err, fs.ErrNotExist) {
		return dialogue.Responses{}, nil
	}
	if err != nil {
		return dialogue.Responses{}, fmt.Errorf(reading, r.ResponsesPath(), err)
	}

	recorded, err := dialogue.Parse(content)
	if err != nil {
		return dialogue.Responses{}, fmt.Errorf(parsing, r.ResponsesPath(), err)
	}
	return recorded, nil
}

func (r Round) AppendResponse(response dialogue.Response) error {
	line, err := dialogue.Line(response)
	if err != nil {
		return err
	}
	return r.appendLine(line)
}

func (r Round) appendLine(line []byte) error {
	file, err := os.OpenFile(r.ResponsesPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf(creating, r.ResponsesPath(), err)
	}
	if _, err := file.Write(line); err != nil {
		_ = file.Close()
		return fmt.Errorf(writing, r.ResponsesPath(), err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf(closing, r.ResponsesPath(), err)
	}
	return nil
}

func (r Round) Findings() (findings.Document, bool, error) {
	content, err := os.ReadFile(r.FindingsPath())
	if errors.Is(err, fs.ErrNotExist) {
		return findings.Document{}, false, nil
	}
	if err != nil {
		return findings.Document{}, false, fmt.Errorf(reading, r.FindingsPath(), err)
	}

	document, err := findings.Parse(content)
	if err != nil {
		return findings.Document{}, false, fmt.Errorf(parsing, r.FindingsPath(), err)
	}
	return document, true, nil
}
