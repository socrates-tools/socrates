package runstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/runstate"
)

const recordName = "run.json"

type Verdict struct {
	Outcome runstate.State `json:"outcome"`
	Reason  string         `json:"reason,omitempty"`
}

type Run struct {
	Name       string         `json:"name"`
	Repository string         `json:"repository"`
	Base       string         `json:"base"`
	Snapshot   string         `json:"snapshot"`
	Worktree   string         `json:"worktree,omitempty"`
	Intent     string         `json:"intent,omitempty"`
	StartedAt  time.Time      `json:"started_at"`
	State      runstate.State `json:"state"`
	Verdict    *Verdict       `json:"verdict,omitempty"`
	Reviewer   int            `json:"reviewer,omitempty"`
	Agent      int            `json:"agent,omitempty"`
}

func (r *Run) Advance(state runstate.State) error {
	if state.Terminal() {
		return fmt.Errorf("%s is terminal, so it is a verdict rather than a step", state)
	}
	return r.transition(state, nil)
}

func (r *Run) End(outcome runstate.State, reason string) error {
	if !outcome.Terminal() {
		return fmt.Errorf("%s is not terminal, so it is not a verdict", outcome)
	}
	if outcome == runstate.Failed && strings.TrimSpace(reason) == "" {
		return errors.New("a failed run states why, so that nothing reads it as a review that found nothing")
	}
	return r.transition(outcome, &Verdict{Outcome: outcome, Reason: reason})
}

func (r *Run) transition(state runstate.State, verdict *Verdict) error {
	if !state.Valid() {
		return fmt.Errorf("%q is not a run state", state)
	}
	if r.State.Terminal() {
		return fmt.Errorf("run %s already ended as %s, so it cannot become %s", r.Name, r.State, state)
	}

	r.State = state
	r.Verdict = verdict
	return nil
}

type Start struct {
	StartedAt      time.Time
	RepositoryRoot string
	Base           string
	Snapshot       string
	Intent         string
}

type Store struct {
	root string
}

func New(root string) Store {
	return Store{root: root}
}

func (s Store) Root() string {
	return s.root
}

func (s Store) Path(name string) string {
	return filepath.Join(s.root, name)
}

func (s Store) Names() ([]string, error) {
	entries, err := os.ReadDir(s.root)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", s.root, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

func (s Store) Create(start Start) (*Run, error) {
	key, err := RepositoryKey(start.RepositoryRoot)
	if err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(start.RepositoryRoot)
	if err != nil {
		return nil, fmt.Errorf("resolving %s: %w", start.RepositoryRoot, err)
	}

	run := &Run{
		Name:       Name(start.StartedAt, key),
		Repository: filepath.Clean(absolute),
		Base:       start.Base,
		Snapshot:   start.Snapshot,
		Intent:     start.Intent,
		StartedAt:  start.StartedAt.UTC().Truncate(time.Second),
		State:      runstate.Starting,
	}

	if err := os.MkdirAll(s.root, 0o755); err != nil {
		return nil, fmt.Errorf("creating %s: %w", s.root, err)
	}

	directory := s.Path(run.Name)
	if err := os.Mkdir(directory, 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf("run %s already exists, so a second review of %s in the same second cannot start: wait a moment and run it again", run.Name, run.Repository)
		}
		return nil, fmt.Errorf("creating %s: %w", directory, err)
	}

	if err := write(directory, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s Store) Load(name string) (*Run, error) {
	run, err := s.recorded(name)
	if err != nil {
		return nil, err
	}
	observe(run)
	return run, nil
}

func (s Store) recorded(name string) (*Run, error) {
	path := filepath.Join(s.Path(name), recordName)

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var run Run
	if err := json.Unmarshal(content, &run); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return &run, nil
}

func observe(run *Run) {
	if run.State.Terminal() || run.Reviewer == 0 || process.Alive(run.Reviewer) {
		return
	}

	run.State = runstate.Failed
	run.Verdict = &Verdict{
		Outcome: runstate.Failed,
		Reason:  "the reviewer went away before it reported anything, so nothing has been reviewed",
	}
}

func (s Store) Save(run *Run) error {
	directory := s.Path(run.Name)
	if _, err := os.Stat(directory); err != nil {
		return fmt.Errorf("run %s has no directory under %s: %w", run.Name, s.root, err)
	}

	stored, err := s.recorded(run.Name)
	switch {
	case err == nil:
		if stored.State.Terminal() && !sameOutcome(stored, run) {
			return fmt.Errorf("run %s already ended as %s, so it cannot be written as %s", run.Name, stored.State, run.State)
		}
	case !errors.Is(err, fs.ErrNotExist):
		return err
	}

	return write(directory, run)
}

func sameOutcome(stored, incoming *Run) bool {
	if stored.State != incoming.State {
		return false
	}
	if (stored.Verdict == nil) != (incoming.Verdict == nil) {
		return false
	}
	if stored.Verdict == nil {
		return true
	}
	return *stored.Verdict == *incoming.Verdict
}

func write(directory string, run *Run) error {
	encoded, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding run %s: %w", run.Name, err)
	}
	encoded = append(encoded, '\n')

	temporary, err := os.CreateTemp(directory, recordName+".")
	if err != nil {
		return fmt.Errorf("creating a temporary file in %s: %w", directory, err)
	}
	defer os.Remove(temporary.Name())

	if _, err := temporary.Write(encoded); err != nil {
		temporary.Close()
		return fmt.Errorf("writing %s: %w", temporary.Name(), err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("flushing %s: %w", temporary.Name(), err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("closing %s: %w", temporary.Name(), err)
	}

	target := filepath.Join(directory, recordName)
	if err := os.Rename(temporary.Name(), target); err != nil {
		return fmt.Errorf("renaming onto %s: %w", target, err)
	}
	return nil
}
