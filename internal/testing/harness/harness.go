package harness

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/socrates-devtool/socrates/internal/process"
	"github.com/socrates-devtool/socrates/internal/runlog"
	"github.com/socrates-devtool/socrates/internal/runstate"
	"github.com/socrates-devtool/socrates/internal/runstore"
	"github.com/socrates-devtool/socrates/internal/snapshot"
	"github.com/socrates-devtool/socrates/internal/state"
	"github.com/socrates-devtool/socrates/internal/testing/gitrepo"
	"github.com/socrates-devtool/socrates/internal/testing/simulator"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

const (
	withoutProcessGroups = "noprocessgroups"

	terminalWait = 30 * time.Second
	reapingWait  = 10 * time.Second
	reapingPoll  = 10 * time.Millisecond
	quietSweeps  = 3
)

type Session struct {
	t          *testing.T
	Binary     string
	Repository *gitrepo.Repository
	Simulator  *simulator.Simulator
	Home       string
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
		Binary:     buildSocrates(t),
		Repository: repository,
		Simulator:  simulated,
		Home:       home,
		Store:      runstore.New(filepath.Join(home, "runs")),
	}
	t.Cleanup(session.liftEverySeal)
	t.Cleanup(session.stopEveryReviewer)
	return session
}

func (s *Session) WhereNothingReaps() *Session {
	s.t.Helper()

	elsewhere := *s
	elsewhere.Binary = buildSocrates(s.t, withoutProcessGroups)
	return &elsewhere
}

func buildSocrates(t *testing.T, tags ...string) string {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "socrates")
	source.BuildWithTags(t, "./cmd/socrates", binary, tags...)
	return binary
}

func (s *Session) Socrates(arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.SocratesIn(s.Repository.Root, arguments...)
}

func (s *Session) SocratesIn(directory string, arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.SocratesWith(directory, nil, arguments...)
}

func (s *Session) SocratesReading(input string, arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.socrates(s.Repository.Root, nil, strings.NewReader(input), arguments...)
}

func (s *Session) SocratesWith(directory string, environment []string, arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.socrates(directory, environment, nil, arguments...)
}

func (s *Session) socrates(directory string, environment []string, input io.Reader, arguments ...string) (string, string, int) {
	s.t.Helper()

	command := exec.Command(s.Binary, arguments...)
	command.Dir = directory
	command.Stdin = input
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
			s.t.Fatalf("running socrates %s: %v", strings.Join(arguments, " "), err)
		}
		code = exited.ExitCode()
	}

	return out.String(), failures.String(), code
}

type Running struct {
	t         *testing.T
	arguments []string
	command   *exec.Cmd
	out       *bytes.Buffer
	failures  *bytes.Buffer
	finished  chan error
}

func (s *Session) Starting(arguments ...string) *Running {
	s.t.Helper()

	command := exec.Command(s.Binary, arguments...)
	command.Dir = s.Repository.Root

	running := &Running{t: s.t, arguments: arguments, command: command, out: &bytes.Buffer{}, failures: &bytes.Buffer{}, finished: make(chan error, 1)}
	command.Stdout = running.out
	command.Stderr = running.failures

	if err := command.Start(); err != nil {
		s.t.Fatalf("starting socrates %s: %v", strings.Join(arguments, " "), err)
	}
	go func() { running.finished <- command.Wait() }()
	s.t.Cleanup(running.stop)
	return running
}

func (r *Running) Wait() (string, string, int) {
	r.t.Helper()

	select {
	case err := <-r.finished:
		code := 0
		if err != nil {
			var exited *exec.ExitError
			if !errors.As(err, &exited) {
				r.t.Fatalf("running socrates %s: %v", strings.Join(r.arguments, " "), err)
			}
			code = exited.ExitCode()
		}
		return r.out.String(), r.failures.String(), code
	case <-time.After(terminalWait):
		r.t.Fatalf("socrates %s is still running after %s", strings.Join(r.arguments, " "), terminalWait)
		return "", "", 0
	}
}

func (r *Running) stop() {
	select {
	case err := <-r.finished:
		r.finished <- err
	default:
		_ = r.command.Process.Kill()
	}
}

func (s *Session) Review(arguments ...string) string {
	s.t.Helper()

	out, failures, code := s.Socrates(append([]string{"review"}, arguments...)...)
	if code != 0 {
		s.t.Fatalf("socrates review exited %d: %s", code, failures)
	}

	name, _, _ := strings.Cut(strings.TrimSpace(out), "\n")
	if name == "" {
		s.t.Fatal("socrates review printed no run identifier")
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
	_ = snapshot.Unseal(s.Home)
}

func (s *Session) stopEveryReviewer() {
	deadline := time.Now().Add(reapingWait)
	quiet := 0

	for quiet < quietSweeps {
		alive := s.killWhatIsStillRunning()
		if alive == 0 {
			quiet++
		} else {
			quiet = 0
		}
		if time.Now().After(deadline) {
			s.t.Errorf("%d of the processes this session started were alive %s after being killed, and whatever they write next lands in a directory the test is deleting", alive, reapingWait)
			return
		}
		time.Sleep(reapingPoll)
	}
}

func (s *Session) killWhatIsStillRunning() int {
	alive := 0

	for _, name := range s.Names() {
		run, err := s.Store.Load(name)
		if err != nil {
			continue
		}
		if process.Alive(run.Agent) {
			alive++
			_ = process.KillGroup(run.Agent)
		}
		if process.Alive(run.Reviewer) {
			alive++
			_ = process.Kill(run.Reviewer)
		}
	}

	for _, pid := range s.Simulator.Processes() {
		if process.Alive(pid) {
			alive++
			_ = process.Kill(pid)
		}
	}

	return alive
}

func reasonOf(run *runstore.Run) string {
	if run.Verdict == nil {
		return "no verdict"
	}
	return run.Verdict.Reason
}
