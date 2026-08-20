package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/process"
	"github.com/rui-ktei/socrates/internal/runlog"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/runstore"
	"github.com/rui-ktei/socrates/internal/snapshot"
	"github.com/rui-ktei/socrates/internal/state"
	"github.com/rui-ktei/socrates/internal/testing/gitrepo"
	"github.com/rui-ktei/socrates/internal/testing/simulator"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

const terminalWait = 30 * time.Second

type Session struct {
	t          *testing.T
	Binary     string
	Repository *gitrepo.Repository
	Simulator  *simulator.Simulator
	StateHome  string
	Store      runstore.Store
}

func Start(t *testing.T) *Session {
	t.Helper()

	repository := gitrepo.New(t)
	simulated := simulator.Install(t)

	home := t.TempDir()
	t.Setenv(state.HomeVariable, home)

	session := &Session{
		t:          t,
		Binary:     buildScrat(t),
		Repository: repository,
		Simulator:  simulated,
		StateHome:  home,
		Store:      runstore.New(filepath.Join(home, "runs")),
	}
	t.Cleanup(session.liftEverySeal)
	t.Cleanup(session.stopEveryReviewer)
	return session
}

func buildScrat(t *testing.T) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "scrat")
	source.Build(t, "./cmd/scrat", binary)
	return binary
}

func (s *Session) Scrat(arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.ScratIn(s.Repository.Root, arguments...)
}

func (s *Session) ScratIn(directory string, arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.ScratWith(directory, nil, arguments...)
}

func (s *Session) ScratWith(directory string, environment []string, arguments ...string) (string, string, int) {
	s.t.Helper()

	command := exec.Command(s.Binary, arguments...)
	command.Dir = directory
	if len(environment) > 0 {
		command.Env = append(os.Environ(), environment...)
	}

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures

	code := 0
	if err := command.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			s.t.Fatalf("running scrat %s: %v", strings.Join(arguments, " "), err)
		}
		code = exited.ExitCode()
	}

	return out.String(), failures.String(), code
}

func (s *Session) Review(arguments ...string) string {
	s.t.Helper()

	out, failures, code := s.Scrat(append([]string{"review"}, arguments...)...)
	if code != 0 {
		s.t.Fatalf("scrat review exited %d: %s", code, failures)
	}

	name, _, _ := strings.Cut(strings.TrimSpace(out), "\n")
	if name == "" {
		s.t.Fatal("scrat review printed no run identifier")
	}
	return name
}

func (s *Session) Run(name string) *runstore.Run {
	s.t.Helper()

	run, err := s.Store.Load(name)
	if err != nil {
		s.t.Fatalf("reading the run %s: %v", name, err)
	}
	return run
}

func (s *Session) Recorded(name string) *runstore.Run {
	s.t.Helper()

	content, err := os.ReadFile(filepath.Join(s.Store.Path(name), "run.json"))
	if err != nil {
		s.t.Fatalf("reading run.json for %s: %v", name, err)
	}

	run := &runstore.Run{}
	if err := json.Unmarshal(content, run); err != nil {
		s.t.Fatalf("parsing run.json for %s: %v", name, err)
	}
	return run
}

func (s *Session) WaitForTerminal(name string) *runstore.Run {
	s.t.Helper()

	deadline := time.Now().Add(terminalWait)
	for time.Now().Before(deadline) {
		run := s.Run(name)
		if run.State.Terminal() {
			return run
		}
		time.Sleep(20 * time.Millisecond)
	}

	run := s.Run(name)
	s.t.Fatalf("the run %s is still %s after %s", name, run.State, terminalWait)
	return run
}

func (s *Session) WaitFor(name string, state runstate.State) *runstore.Run {
	s.t.Helper()

	deadline := time.Now().Add(terminalWait)
	for time.Now().Before(deadline) {
		run := s.Run(name)
		if run.State == state {
			return run
		}
		if run.State.Terminal() {
			s.t.Fatalf("the run %s ended as %s, want %s: %s", name, run.State, state, reasonOf(run))
		}
		time.Sleep(20 * time.Millisecond)
	}

	s.t.Fatalf("the run %s never reached %s", name, state)
	return nil
}

func (s *Session) Round(name string, number int) runstore.Round {
	s.t.Helper()
	return s.Store.Round(name, number)
}

func (s *Session) File(path string) string {
	s.t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		s.t.Fatalf("reading %s: %v", path, err)
	}
	return string(content)
}

func (s *Session) Events(name string) []runlog.Event {
	s.t.Helper()

	events, err := runlog.Read(filepath.Join(s.Store.Path(name), runlog.FileName))
	if err != nil {
		s.t.Fatalf("reading the event log for %s: %v", name, err)
	}
	return events
}

func (s *Session) SupervisorLog(name string) string {
	s.t.Helper()

	content, err := os.ReadFile(filepath.Join(s.Store.Path(name), "supervisor.log"))
	if err != nil {
		return ""
	}
	return string(content)
}

func (s *Session) Names() []string {
	s.t.Helper()

	names, err := s.Store.Names()
	if err != nil {
		s.t.Fatalf("listing the runs: %v", err)
	}
	return names
}

func (s *Session) liftEverySeal() {
	_ = snapshot.Unseal(s.StateHome)
}

func (s *Session) stopEveryReviewer() {
	for _, name := range s.Names() {
		run, err := s.Store.Load(name)
		if err != nil {
			continue
		}
		_ = process.KillGroup(run.Agent)
		_ = process.Kill(run.Reviewer)
	}
}

func reasonOf(run *runstore.Run) string {
	if run.Verdict == nil {
		return "no verdict"
	}
	return run.Verdict.Reason
}
