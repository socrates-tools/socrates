package runstore

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/runstate"
)

func TestAReaderNeverSeesHalfAState(t *testing.T) {
	store := New(t.TempDir())
	run := started(t, store)

	held := map[runstate.State]bool{runstate.Starting: true, runstate.Reviewing: true, runstate.Escalated: true}

	type observation struct {
		state runstate.State
		err   error
	}
	done := make(chan struct{})
	results := make(chan []observation, 1)

	go func() {
		var seen []observation
		for {
			select {
			case <-done:
				results <- seen
				return
			default:
			}
			loaded, err := store.Load(run.Name)
			if err != nil {
				seen = append(seen, observation{err: err})
				continue
			}
			seen = append(seen, observation{state: loaded.State})
		}
	}()

	run.State = runstate.Reviewing
	for range 200 {
		if err := store.Save(run); err != nil {
			close(done)
			<-results
			t.Fatalf("saving the run: %v", err)
		}
	}

	run.State = runstate.Escalated
	run.Verdict = &Verdict{Outcome: runstate.Escalated}
	if err := store.Save(run); err != nil {
		close(done)
		<-results
		t.Fatalf("ending the run: %v", err)
	}
	close(done)

	seen := <-results
	if len(seen) == 0 {
		t.Fatal("the reader observed nothing, so it proved nothing")
	}
	for _, observed := range seen {
		if observed.err != nil {
			t.Errorf("a read did not yield a complete state: %v", observed.err)
			continue
		}
		if !held[observed.state] {
			t.Errorf("a read yielded %q, which the run never held", observed.state)
		}
	}
}

func TestACrashMidWriteLeavesThePreviousState(t *testing.T) {
	store := New(t.TempDir())
	run := started(t, store)

	run.State = runstate.Reviewing
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the run: %v", err)
	}

	directory := store.Path(run.Name)
	crashed, err := os.CreateTemp(directory, recordName+".")
	if err != nil {
		t.Fatalf("writing the temporary file a crash would leave: %v", err)
	}
	if _, err := crashed.WriteString(`{"state": "cle`); err != nil {
		t.Fatalf("writing half a state: %v", err)
	}
	if err := crashed.Close(); err != nil {
		t.Fatalf("closing the temporary file: %v", err)
	}

	loaded, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the run after the crash: %v", err)
	}
	if loaded.State != runstate.Reviewing {
		t.Errorf("the run reads as %q, want %q", loaded.State, runstate.Reviewing)
	}

	if err := os.Remove(crashed.Name()); err != nil {
		t.Fatalf("removing the temporary file: %v", err)
	}
	if err := store.Save(run); err != nil {
		t.Fatalf("saving the run again: %v", err)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}
	if len(entries) != 1 || entries[0].Name() != recordName {
		var names []string
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Errorf("a completed write left %v behind, want only %s", names, recordName)
	}
}

func TestASecondRunInTheSameSecondIsRefused(t *testing.T) {
	store := New(t.TempDir())
	repository := t.TempDir()
	startedAt := time.Date(2026, 8, 19, 4, 15, 30, 0, time.UTC)

	if _, err := store.Create(Start{StartedAt: startedAt, RepositoryRoot: repository}); err != nil {
		t.Fatalf("creating the run: %v", err)
	}
	if _, err := store.Create(Start{StartedAt: startedAt, RepositoryRoot: repository}); err == nil {
		t.Fatal("a second run in the same second was accepted, so one of them overwrote the other")
	}
}

func TestCreateRecordsWhatItWasGiven(t *testing.T) {
	store := New(t.TempDir())
	repository := t.TempDir()
	startedAt := time.Date(2026, 8, 19, 4, 15, 30, 500_000_000, time.UTC)

	run, err := store.Create(Start{
		StartedAt:      startedAt,
		RepositoryRoot: filepath.Join(repository, "elsewhere", ".."),
		Base:           "base-commit",
		Snapshot:       "snapshot-commit",
		Intent:         "what the person asked for",
	})
	if err != nil {
		t.Fatalf("creating the run: %v", err)
	}

	loaded, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the run: %v", err)
	}

	if loaded.Repository != repository {
		t.Errorf("the run names %q as its repository, want %q", loaded.Repository, repository)
	}
	if loaded.Base != "base-commit" || loaded.Snapshot != "snapshot-commit" {
		t.Errorf("the run names %q and %q, want the base and snapshot it was given", loaded.Base, loaded.Snapshot)
	}
	if loaded.Intent != "what the person asked for" {
		t.Errorf("the run states its intent as %q", loaded.Intent)
	}
	if !loaded.StartedAt.Equal(startedAt.Truncate(time.Second)) {
		t.Errorf("the run started at %s, want %s", loaded.StartedAt, startedAt.Truncate(time.Second))
	}
	if loaded.State != runstate.Starting {
		t.Errorf("a new run is %q, want %q", loaded.State, runstate.Starting)
	}
	if loaded.Verdict != nil {
		t.Errorf("a new run already holds the verdict %+v", *loaded.Verdict)
	}
}

func started(t *testing.T, store Store) *Run {
	t.Helper()

	run, err := store.Create(Start{
		StartedAt:      time.Date(2026, 8, 19, 4, 15, 30, 0, time.UTC),
		RepositoryRoot: t.TempDir(),
		Base:           "base-commit",
		Snapshot:       "snapshot-commit",
		Intent:         "what the person asked for",
	})
	if err != nil {
		t.Fatalf("creating the run: %v", err)
	}
	return run
}

func TestARunFromBeforeTheFieldIsReadAsTheOriginalCLI(t *testing.T) {
	before := Run{Name: "20260101T000000Z-00000000", State: runstate.Clean}

	if before.Vendor() != agent.DefaultVendor {
		t.Errorf("a run recorded without a reviewer CLI reads as %q, want %q", before.Vendor(), agent.DefaultVendor)
	}
	if before.Tier() != string(agent.TheDefault().Tier) {
		t.Errorf("a run recorded without a reviewer CLI reads as the tier %q, want %q", before.Tier(), agent.TheDefault().Tier)
	}

	named := Run{ReviewerCLI: "codex", ReviewerTier: string(agent.Parsed)}
	if named.Vendor() != "codex" || named.Tier() != string(agent.Parsed) {
		t.Errorf("a run recorded with a reviewer CLI reads as %q on the %q tier", named.Vendor(), named.Tier())
	}
}

func TestCreateRecordsTheReviewerCLIItWasGivenAndItsTier(t *testing.T) {
	store := New(t.TempDir())

	run, err := store.Create(Start{StartedAt: time.Now(), RepositoryRoot: t.TempDir(), ReviewerCLI: "codex"})
	if err != nil {
		t.Fatalf("creating the run: %v", err)
	}

	loaded, err := store.Load(run.Name)
	if err != nil {
		t.Fatalf("reading the run: %v", err)
	}
	codex, _ := agent.Named("codex")
	if loaded.ReviewerCLI != "codex" || loaded.Tier() != string(codex.Tier) {
		t.Errorf("the run reads as %q on the %q tier, want codex on %q", loaded.ReviewerCLI, loaded.Tier(), codex.Tier)
	}
}
