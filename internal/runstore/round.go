package runstore

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	promptFile   = "prompt.txt"
	streamFile   = "stream.jsonl"
	findingsFile = "findings.json"
)

type Round struct {
	Number    int
	Directory string
}

func (s Store) Round(name string, number int) Round {
	return Round{
		Number:    number,
		Directory: filepath.Join(s.Path(name), fmt.Sprintf("r%d", number)),
	}
}

func (s Store) CreateRound(name string, number int) (Round, error) {
	round := s.Round(name, number)
	if err := os.MkdirAll(round.Directory, 0o755); err != nil {
		return Round{}, fmt.Errorf("creating %s: %w", round.Directory, err)
	}
	return round, nil
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
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
