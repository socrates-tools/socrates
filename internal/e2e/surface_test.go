package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rui-ktei/socrates/internal/surface"
	"github.com/rui-ktei/socrates/internal/testing/gitrepo"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

const oneDecisionDocument = `{
  "title": "one thing to settle",
  "preamble": "what follows is what was seen",
  "decisions": [
    {
      "id": "c-1",
      "question": "should the cache be read before the lock is taken",
      "meaning": "two callers arriving together read the same stale entry",
      "mode": "single-choice",
      "options": [
        {"label": "take the lock first", "meaning": "correctness first", "command": "patch c-1"},
        {"label": "leave it", "meaning": "acceptable here", "command": "agree c-1"}
      ]
    }
  ]
}`

type surfaced struct {
	t      *testing.T
	binary string
	home   string
	where  string
	bind   string
	opened string
}

func standing(t *testing.T) *surfaced {
	t.Helper()

	gitrepo.Hermetic(t)
	return servingSurface(t)
}

func servingSurface(t *testing.T) *surfaced {
	t.Helper()

	binary := filepath.Join(t.TempDir(), "agora")
	source.Build(t, "./cmd/agora", binary)

	home := t.TempDir()
	held := &surfaced{t: t, binary: binary, home: home, where: t.TempDir(), bind: aPortTheOperatingSystemChooses}
	t.Cleanup(func() { held.run("stop") })
	return held
}

const aPortTheOperatingSystemChooses = "127.0.0.1:0"

func (s *surfaced) bindsWhatIsConfigured() {
	s.bind = ""
}

func (s *surfaced) run(arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.reading("", arguments...)
}

func (s *surfaced) reading(input string, arguments ...string) (string, string, int) {
	s.t.Helper()

	command := exec.Command(s.binary, arguments...)
	command.Dir = s.where
	command.Env = append(os.Environ(), surface.HomeVariable+"="+s.home)
	command.Env = append(command.Env, source.Toolchain()...)
	command.Stdin = strings.NewReader(input)

	var printed, failures bytes.Buffer
	command.Stdout = &printed
	command.Stderr = &failures

	code := 0
	if err := command.Run(); err != nil {
		var exited *exec.ExitError
		if !errors.As(err, &exited) {
			s.t.Fatalf("running agora %v: %v", arguments, err)
		}
		code = exited.ExitCode()
	}
	return printed.String(), failures.String(), code
}

func (s *surfaced) serve(document string) string {
	s.t.Helper()

	printed, failures, code := s.reading(document, s.serving()...)
	if code != 0 {
		s.t.Fatalf("agora serve exited %d: %s", code, failures)
	}
	address := strings.TrimSpace(printed)
	s.opened = theSessionIn(s.t, address)
	return address
}

func (s *surfaced) serveAgain(document string) string {
	s.t.Helper()
	return s.serveOnto(s.opened, document)
}

func (s *surfaced) serveOnto(session, document string) string {
	s.t.Helper()

	printed, failures, code := s.reading(document, s.serving("--session", session)...)
	if code != 0 {
		s.t.Fatalf("agora serve --session exited %d: %s", code, failures)
	}
	return strings.TrimSpace(printed)
}

func theSessionIn(t *testing.T, address string) string {
	t.Helper()

	_, session, found := strings.Cut(address, surface.SessionPrefix)
	if !found || session == "" {
		t.Fatalf("agora serve printed %q, and the address it prints holds the session the page belongs to", address)
	}
	return session
}

func (s *surfaced) serving(arguments ...string) []string {
	asked := append([]string{"serve", "--document", "-"}, arguments...)
	if s.bind != "" {
		asked = append(asked, "--address", s.bind)
	}
	return asked
}

func (s *surfaced) record() surface.Record {
	s.t.Helper()

	held, err := os.ReadFile(filepath.Join(s.home, "surface.json"))
	if err != nil {
		s.t.Fatalf("reading what the surface wrote about itself: %v", err)
	}
	var open surface.Record
	if err := json.Unmarshal(held, &open); err != nil {
		s.t.Fatalf("reading back what the surface wrote about itself: %v", err)
	}
	return open
}

func (s *surfaced) poll(arguments ...string) (string, string, int) {
	s.t.Helper()
	return s.run(append([]string{"poll", "--session", s.opened}, arguments...)...)
}

func (s *surfaced) end(session string) {
	s.t.Helper()

	if _, failures, code := s.run("end", "--session", session); code != 0 {
		s.t.Fatalf("agora end exited %d: %s", code, failures)
	}
}

func (s *surfaced) answer(sent surface.Sent) {
	s.t.Helper()
	s.answerOn(s.opened, sent)
}

func (s *surfaced) answerOn(session string, sent surface.Sent) {
	s.t.Helper()

	open := s.record()
	written, err := json.Marshal(sent)
	if err != nil {
		s.t.Fatalf("writing what a person sent: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, "http://"+open.Address+surface.SessionPrefix+session+"/answer", bytes.NewReader(written))
	if err != nil {
		s.t.Fatalf("addressing the surface: %v", err)
	}
	request.Header.Set(surface.TokenHeader, open.Token)
	request.Header.Set("content-type", "application/json")

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		s.t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()
	if answered.StatusCode != http.StatusNoContent {
		said, _ := io.ReadAll(answered.Body)
		s.t.Fatalf("sending an answer was refused with %d: %s", answered.StatusCode, said)
	}
}

func TestTheSurfaceOffersNothingBeyondTheFive(t *testing.T) {
	held := standing(t)

	printed, failures, code := held.run("--help")
	if code != 0 {
		t.Fatalf("agora --help exited %d: %s", code, failures)
	}

	offered := commandsIn(printed)
	offered = slices.DeleteFunc(offered, func(named string) bool { return named == "help" })
	slices.Sort(offered)

	want := []string{"end", "poll", "serve", "stop", "version"}
	if !slices.Equal(want, offered) {
		t.Errorf("the surface offers %v beside the help the command library brings, want %v", offered, want)
	}

	for _, absent := range []string{"publish", "export", "share", "annotate", "open", "screenshot"} {
		if strings.Contains(strings.ToLower(printed), absent) {
			t.Errorf("the surface names %q, and it answers a decision and sends", absent)
		}
	}
}

func TestSocratesKeepsItsTwoNamesAndGainsNoThird(t *testing.T) {
	gitrepo.Hermetic(t)

	directory := t.TempDir()
	socrates := filepath.Join(directory, "socrates")
	short := filepath.Join(directory, shortName())
	agora := filepath.Join(directory, "agora")

	source.Build(t, "./cmd/socrates", socrates)
	source.Build(t, "./cmd/agora", agora)
	if err := os.Symlink(socrates, short); err != nil {
		t.Fatalf("linking the second name beside the first: %v", err)
	}

	underOneName := said(t, socrates, "version")
	underTheOther := said(t, short, "version")
	if underOneName != underTheOther {
		t.Errorf("the two names are not one program:\n%s\n%s", underOneName, underTheOther)
	}
	if !strings.Contains(underTheOther, "socrates") {
		t.Errorf("invoked under its second name the program does not name socrates: %s", underTheOther)
	}

	separate := said(t, agora, "version")
	if strings.Contains(separate, "socrates") {
		t.Errorf("the surface reports itself as socrates: %s", separate)
	}
	if !strings.Contains(separate, "agora") {
		t.Errorf("the surface does not name itself: %s", separate)
	}
}

func shortName() string {
	return string([]rune{'s', 'c', 'r', 'a', 't'})
}

func said(t *testing.T, binary string, arguments ...string) string {
	t.Helper()

	command := exec.Command(binary, arguments...)
	command.Env = append(os.Environ(), source.Toolchain()...)

	var printed bytes.Buffer
	command.Stdout = &printed
	if err := command.Run(); err != nil {
		t.Fatalf("running %s %v: %v", binary, arguments, err)
	}
	return strings.TrimSpace(printed.String())
}

func TestServeReturnsBeforeThePersonHasReadAnything(t *testing.T) {
	held := standing(t)

	started := time.Now()
	address := held.serve(oneDecisionDocument)
	took := time.Since(started)

	if took > 20*time.Second {
		t.Errorf("agora serve took %s to return, and it returns as soon as the page is reachable", took)
	}
	if err := surface.CheckLoopback(originOf(t, address)[len("http://"):]); err != nil {
		t.Fatalf("agora serve printed %q: %v", address, err)
	}

	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching the address agora serve printed: %v", err)
	}
	defer answered.Body.Close()

	page, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the page: %v", err)
	}
	if answered.StatusCode != http.StatusOK {
		t.Fatalf("the page at the printed address answered with %d", answered.StatusCode)
	}
	if !strings.Contains(string(page), "should the cache be read before the lock is taken") {
		t.Errorf("the page at the printed address is not the document that was served:\n%s", page)
	}
}

func TestAWildcardBindIsRefusedByTheBinary(t *testing.T) {
	// scenario: A wildcard bind is refused rather than quietly narrowed
	held := standing(t)

	printed, failures, code := held.reading(oneDecisionDocument, "serve", "--document", "-", "--address", "0.0.0.0:0")
	if code == 0 {
		t.Fatalf("agora serve bound every interface and printed %q", printed)
	}
	if !strings.Contains(failures, "0.0.0.0:0") {
		t.Errorf("refusing a wildcard bind does not name it: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.home, "surface.json")); err == nil {
		t.Error("a refused bind left a session behind, and it starts no server")
	}
}

func TestAFragmentOutsideTheAllowlistStillServes(t *testing.T) {
	// scenario: A fragment outside the allowlist is stripped
	held := standing(t)

	var written map[string]any
	if err := json.Unmarshal([]byte(oneDecisionDocument), &written); err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	decisions := written["decisions"].([]any)
	decisions[0].(map[string]any)["fragment"] = map[string]any{
		"markup": `<figure><script>fetch("https://elsewhere")</script><p onclick="steal()">drawn</p><form><input name="who"></form></figure>`,
	}

	raw, err := json.Marshal(written)
	if err != nil {
		t.Fatalf("writing the document: %v", err)
	}

	address := held.serve(string(raw))

	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()

	page, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the page: %v", err)
	}
	if answered.StatusCode != http.StatusOK {
		t.Fatalf("a document carrying a fragment was served with %d", answered.StatusCode)
	}
	for _, gone := range []string{"<script>", "onclick", "elsewhere", "steal(", `name="who"`} {
		if strings.Contains(string(page), gone) {
			t.Errorf("the page carries %q from the fragment it was given", gone)
		}
	}
	if !strings.Contains(string(page), "drawn") {
		t.Error("neutralising the fragment took the illustration with it")
	}
}

const twoRecommendedDocument = `{
  "title": "one thing to settle",
  "decisions": [
    {
      "id": "c-1",
      "question": "should the cache be read before the lock is taken",
      "mode": "single-choice",
      "options": [
        {"label": "take the lock first", "command": "patch c-1", "recommended": true},
        {"label": "leave it", "command": "agree c-1", "recommended": true}
      ]
    }
  ]
}`

func TestServeRefusesADocumentMarkingTwoOptions(t *testing.T) {
	// scenario: A document marking two options is refused before serving
	held := standing(t)

	printed, failures, code := held.reading(twoRecommendedDocument, "serve", "--document", "-")
	if code == 0 {
		t.Fatalf("a document marking two options was served at %q", printed)
	}
	if !strings.Contains(failures, "c-1") {
		t.Errorf("refusing a document marking two options does not name the decision: %s", failures)
	}
	if _, err := os.Stat(filepath.Join(held.home, "surface.json")); err == nil {
		t.Error("a refused document left a session behind, and it starts no server")
	}
}

func TestADocumentWrittenBeforeTheMarkStillServes(t *testing.T) {
	// scenario: A document that marks nothing is unaffected
	held := standing(t)

	address := held.serve(oneDecisionDocument)

	answered, err := http.Get(address)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()

	page, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the page: %v", err)
	}
	if answered.StatusCode != http.StatusOK {
		t.Fatalf("a document that marks nothing was served with %d", answered.StatusCode)
	}
	if strings.Contains(string(page), "option-mark") || strings.Contains(string(page), "data-recommended") {
		t.Errorf("the page marks an option, and the document marked none:\n%s", page)
	}

	first := strings.Index(string(page), "take the lock first")
	second := strings.Index(string(page), "leave it")
	if first < 0 || second < 0 || first > second {
		t.Errorf("the options are not read in the order the document listed them: %d then %d", first, second)
	}
}

func TestTheRepositoryIsUntouched(t *testing.T) {
	held := standing(t)

	repository := gitrepo.New(t)
	held.where = repository.Root
	before := repository.Status()
	beforeTree := workingTree(t, repository.Root)

	_, first, _, second := held.twoPages()
	for _, session := range []string{first, second} {
		held.answerOn(session, surface.Sent{
			Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "go on"}},
			Message:   "one more look",
		})

		printed, failures, code := held.run("poll", "--session", session)
		if code != 0 {
			t.Fatalf("agora poll exited %d: %s", code, failures)
		}
		if !strings.Contains(printed, "patch c-1") {
			t.Errorf("what agora poll printed does not carry the command the option carries:\n%s", printed)
		}
		held.end(session)
	}

	if _, failures, code := held.run("stop"); code != 0 {
		t.Fatalf("agora stop exited %d: %s", code, failures)
	}

	if after := repository.Status(); after != before {
		t.Errorf("git status changed while two surfaces were served, answered and ended:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if after := workingTree(t, repository.Root); !equalTrees(beforeTree, after) {
		t.Errorf("the working tree changed while two surfaces were served, answered and ended:\nbefore: %v\nafter: %v", beforeTree, after)
	}
}

func equalTrees(before, after map[string]string) bool {
	if len(before) != len(after) {
		return false
	}
	for path, content := range before {
		if after[path] != content {
			return false
		}
	}
	return true
}

const surfaceCeiling = 32 << 20

func TestTheSurfaceWeighsWhatSomebodyHasSeen(t *testing.T) {
	// scenario: A diff renders from a patch
	gitrepo.Hermetic(t)

	binary := filepath.Join(t.TempDir(), "agora")
	source.Build(t, "./cmd/agora", binary)

	held, err := os.Stat(binary)
	if err != nil {
		t.Fatalf("measuring the surface: %v", err)
	}

	t.Logf("the surface weighs %d bytes", held.Size())
	if held.Size() > surfaceCeiling {
		t.Errorf("the surface weighs %d bytes, and the weight somebody looked at and accepted is %d", held.Size(), surfaceCeiling)
	}
}

func (s *surfaced) awaiting(within time.Duration) (surface.Answered, bool) {
	s.t.Helper()
	return s.awaitingOn(s.opened, within)
}

func (s *surfaced) awaitingOn(session string, within time.Duration) (surface.Answered, bool) {
	s.t.Helper()

	read, arrived, err := surface.Reach(s.record()).Await(context.Background(), session, within)
	if err != nil {
		s.t.Fatalf("waiting on the surface: %v", err)
	}
	return read, arrived
}
