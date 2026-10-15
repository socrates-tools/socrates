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

	"github.com/socrates-tools/socrates/internal/process"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/dialogue"
	"github.com/socrates-tools/socrates/review/guidance"
	"github.com/socrates-tools/socrates/review/prompt"
	"github.com/socrates-tools/socrates/review/runstate"
	"github.com/socrates-tools/socrates/review/snapshot"
)

const recordName = "run.json"

type Verdict struct {
	Outcome runstate.State `json:"outcome"`
	Reason  string         `json:"reason,omitempty"`
}

type Reviewed struct {
	Round    int    `json:"round"`
	Base     string `json:"base"`
	Snapshot string `json:"snapshot"`
}

type Run struct {
	Name            string           `json:"name"`
	Repository      string           `json:"repository"`
	Base            string           `json:"base"`
	Snapshot        string           `json:"snapshot"`
	Worktree        string           `json:"worktree,omitempty"`
	Intent          string           `json:"intent,omitempty"`
	NotMyChange     bool             `json:"not_my_change,omitempty"`
	ReviewGuide     guidance.Outcome `json:"review_guide"`
	GuidanceCeiling int64            `json:"guidance_ceiling,omitempty"`
	Diff            prompt.Delivery  `json:"diff,omitempty"`
	StartedAt       time.Time        `json:"started_at"`
	State           runstate.State   `json:"state"`
	Verdict         *Verdict         `json:"verdict,omitempty"`
	ReviewerCLI     string           `json:"reviewer_cli,omitempty"`
	ReviewerTier    string           `json:"reviewer_tier,omitempty"`
	Reviewer        int              `json:"reviewer,omitempty"`
	Agent           int              `json:"agent,omitempty"`
	Reviews         []Reviewed       `json:"reviews,omitempty"`
}

func (r Run) Vendor() string {
	if r.ReviewerCLI == "" {
		return agent.DefaultVendor
	}
	return r.ReviewerCLI
}

func (r Run) Authorship() dialogue.Authorship {
	return dialogue.Authorship(r.NotMyChange)
}

func (r Run) Tier() string {
	if r.ReviewerCLI == "" {
		return string(agent.TheDefault().Tier)
	}
	return r.ReviewerTier
}

func (r *Run) Snapshotted(round int, base, snapshot string) {
	reviewed := Reviewed{Round: round, Base: base, Snapshot: snapshot}

	for position, already := range r.Reviews {
		if already.Round == round {
			r.Reviews[position] = reviewed
			return
		}
	}
	r.Reviews = append(r.Reviews, reviewed)
}

func (r *Run) Reviewed(round int) (Reviewed, bool) {
	for _, already := range r.Reviews {
		if already.Round == round {
			return already, true
		}
	}
	return Reviewed{}, false
}

func (r *Run) Advance(state runstate.State) error {
	if state.Terminal() {
		return fmt.Errorf(terminalIsNotAStep, state)
	}
	return r.transition(state, nil)
}

func (r *Run) End(outcome runstate.State, reason string) error {
	if !outcome.Terminal() {
		return fmt.Errorf(notTerminalIsNotAVerdict, outcome)
	}
	if outcome == runstate.Failed && strings.TrimSpace(reason) == "" {
		return errors.New(aFailureStatesWhy)
	}
	return r.transition(outcome, &Verdict{Outcome: outcome, Reason: reason})
}

func (r *Run) transition(state runstate.State, verdict *Verdict) error {
	if !state.Valid() {
		return fmt.Errorf(notARunState, state)
	}
	if r.State.Terminal() {
		return fmt.Errorf(alreadyEnded, r.Name, r.State, state)
	}

	r.State = state
	r.Verdict = verdict
	return nil
}

type Start struct {
	StartedAt       time.Time
	RepositoryRoot  string
	Base            string
	Snapshot        string
	Intent          string
	NotMyChange     bool
	ReviewGuide     guidance.Outcome
	GuidanceCeiling int64
	ReviewerCLI     string
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
		return nil, fmt.Errorf(reading, s.root, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}

const DefaultRetention = 14 * 24 * time.Hour

func (s Store) Sweep(now time.Time, retention time.Duration) error {
	if retention <= 0 {
		retention = DefaultRetention
	}

	names, err := s.Names()
	if err != nil {
		return err
	}

	for _, name := range names {
		run, err := s.recorded(name)
		if err != nil {
			continue
		}
		if !run.State.Terminal() || now.Sub(run.StartedAt) < retention {
			continue
		}
		if err := s.Discard(name); err != nil {
			return err
		}
	}
	return nil
}

func (s Store) Discard(name string) error {
	directory := s.Path(name)

	if err := snapshot.Unseal(directory); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := os.RemoveAll(directory); err != nil {
		return fmt.Errorf(removing, directory, err)
	}
	return nil
}

func (s Store) Select(identifier, repositoryRoot string) (string, error) {
	if identifier != "" {
		return identifier, nil
	}

	key, err := RepositoryKey(repositoryRoot)
	if err != nil {
		return "", err
	}
	names, err := s.Names()
	if err != nil {
		return "", err
	}

	latest, found := Latest(names, key)
	if !found {
		return "", refusal.New(noRunForThisRepository, repositoryRoot)
	}
	return latest, nil
}

func (s Store) Create(start Start) (*Run, error) {
	key, err := RepositoryKey(start.RepositoryRoot)
	if err != nil {
		return nil, err
	}
	absolute, err := filepath.Abs(start.RepositoryRoot)
	if err != nil {
		return nil, fmt.Errorf(resolving, start.RepositoryRoot, err)
	}

	run := &Run{
		Name:            Name(start.StartedAt, key),
		Repository:      filepath.Clean(absolute),
		Base:            start.Base,
		Snapshot:        start.Snapshot,
		Intent:          start.Intent,
		NotMyChange:     start.NotMyChange,
		ReviewGuide:     start.ReviewGuide,
		GuidanceCeiling: start.GuidanceCeiling,
		ReviewerCLI:     start.ReviewerCLI,
		StartedAt:       start.StartedAt.UTC().Truncate(time.Second),
		State:           runstate.Starting,
	}
	if vendor, known := agent.Named(start.ReviewerCLI); known {
		run.ReviewerTier = string(vendor.Tier)
	}

	if err := os.MkdirAll(s.root, 0o755); err != nil {
		return nil, fmt.Errorf(creating, s.root, err)
	}

	directory := s.Path(run.Name)
	if err := os.Mkdir(directory, 0o755); err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf(alreadyExists, run.Name, run.Repository)
		}
		return nil, fmt.Errorf(creating, directory, err)
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
	if !reviewerGone(run) {
		return run, nil
	}

	settled, err := s.recorded(name)
	if err != nil {
		return nil, err
	}
	if reviewerGone(settled) {
		settled.State = runstate.Failed
		settled.Verdict = &Verdict{
			Outcome: runstate.Failed,
			Reason:  reviewerWentAway,
		}
	}
	return settled, nil
}

func (s Store) recorded(name string) (*Run, error) {
	path := filepath.Join(s.Path(name), recordName)

	content, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, refusal.Wrapping(err, noSuchRun, name, s.root)
	}
	if err != nil {
		return nil, fmt.Errorf(reading, path, err)
	}

	var run Run
	if err := json.Unmarshal(content, &run); err != nil {
		return nil, fmt.Errorf(parsing, path, err)
	}
	return &run, nil
}

func reviewerGone(run *Run) bool {
	if run.State.Terminal() || run.State == runstate.AwaitingResponse || run.Reviewer == 0 {
		return false
	}
	return !process.Alive(run.Reviewer)
}

func (s Store) Save(run *Run) error {
	directory := s.Path(run.Name)
	if _, err := os.Stat(directory); err != nil {
		return fmt.Errorf(noDirectory, run.Name, s.root, err)
	}

	stored, err := s.recorded(run.Name)
	switch {
	case err == nil:
		if stored.State.Terminal() && !sameOutcome(stored, run) {
			return fmt.Errorf(alreadyEndedOnDisk, run.Name, stored.State, run.State)
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
		return fmt.Errorf(encoding, run.Name, err)
	}
	encoded = append(encoded, '\n')

	return replaceFile(filepath.Join(directory, recordName), encoded)
}
