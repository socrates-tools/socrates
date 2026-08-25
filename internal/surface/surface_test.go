package surface_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/document"
	"github.com/rui-ktei/socrates/internal/surface"
)

const token = "a-token-the-page-carries"

func twoDecisions() document.Document {
	return document.Document{
		Title: "two things to settle",
		Decisions: []document.Decision{
			{
				ID:       "c-1",
				Question: "should the cache be read before the lock is taken",
				Mode:     document.One,
				Options: []document.Option{
					{Label: "take the lock first", Command: "patch c-1"},
					{Label: "leave it", Command: `agree c-1 --reason "acceptable"`},
				},
			},
			{
				ID:       "c-2",
				Question: "which of these should the write path keep",
				Mode:     document.Many,
				Options: []document.Option{
					{Label: "the write", Command: "keep c-2 write"},
					{Label: "the metric", Command: "keep c-2 metric"},
					{Label: "the log line", Command: "keep c-2 log"},
				},
			},
		},
	}
}

type standing struct {
	server *surface.Server
	home   string
}

func serving(t *testing.T, read document.Document, idle time.Duration) *standing {
	t.Helper()

	home := t.TempDir()
	t.Setenv(surface.HomeVariable, home)

	server, err := surface.Serve(read, surface.DefaultAddress, token, idle)
	if err != nil {
		t.Fatalf("serving: %v", err)
	}
	t.Cleanup(server.Stop)

	return &standing{server: server, home: home}
}

func (s *standing) at(route string) string {
	return "http://" + s.server.Address() + route
}

func (s *standing) send(t *testing.T, sent surface.Sent) int {
	t.Helper()

	written, err := json.Marshal(sent)
	if err != nil {
		t.Fatalf("writing what a person sent: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, s.at("/answer"), bytes.NewReader(written))
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Header.Set(surface.TokenHeader, token)
	request.Header.Set("content-type", "application/json")

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()
	return answered.StatusCode
}

func (s *standing) await(t *testing.T, within time.Duration) (surface.Answered, bool) {
	t.Helper()

	open := surface.Session{Address: s.server.Address(), Token: token}
	read, arrived, err := surface.Reach(open).Await(context.Background(), within)
	if err != nil {
		t.Fatalf("waiting on the surface: %v", err)
	}
	return read, arrived
}

func TestAWildcardBindIsRefusedRatherThanQuietlyNarrowed(t *testing.T) {
	for _, refused := range []string{"0.0.0.0:0", ":8080", "[::]:0", "192.168.1.10:0", "example.test:80"} {
		t.Run(refused, func(t *testing.T) {
			before := open(t)

			server, err := surface.Serve(twoDecisions(), refused, token, 0)
			if err == nil {
				server.Stop()
				t.Fatalf("the surface bound %q, and it binds a loopback address only", refused)
			}
			if !strings.Contains(err.Error(), refused) {
				t.Errorf("refusing %q does not name it: %v", refused, err)
			}
			if after := open(t); after != before {
				t.Errorf("refusing %q left %d listeners open, and it starts no server", refused, after-before)
			}
		})
	}
}

func TestALoopbackBindIsAccepted(t *testing.T) {
	for _, accepted := range []string{"127.0.0.1:0", "localhost:0", "[::1]:0"} {
		t.Run(accepted, func(t *testing.T) {
			if err := surface.CheckLoopback(accepted); err != nil {
				t.Errorf("the surface refuses the loopback address %q: %v", accepted, err)
			}
		})
	}
}

func open(t *testing.T) int {
	t.Helper()

	held, err := os.ReadDir("/proc/self/fd")
	if err != nil {
		return 0
	}
	return len(held)
}

func TestARequestClaimingAnotherHostIsRejected(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	request, err := http.NewRequest(http.MethodGet, held.at("/"), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	request.Host = "somewhere.test"

	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()

	if answered.StatusCode != http.StatusMisdirectedRequest {
		t.Errorf("a request naming another host was answered with %d", answered.StatusCode)
	}
	said, _ := io.ReadAll(answered.Body)
	if !strings.Contains(string(said), "somewhere.test") {
		t.Errorf("rejecting a request naming another host does not name it: %s", said)
	}
}

func TestAnAnswerWithoutTheTokenIsRejected(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	written, err := json.Marshal(surface.Sent{Message: "let me in"})
	if err != nil {
		t.Fatalf("writing what a person sent: %v", err)
	}

	answered, err := http.Post(held.at("/answer"), "application/json", bytes.NewReader(written))
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer answered.Body.Close()

	if answered.StatusCode != http.StatusForbidden {
		t.Errorf("an answer carrying no token was answered with %d", answered.StatusCode)
	}

	read, arrived := held.await(t, 200*time.Millisecond)
	if arrived {
		t.Errorf("an answer carrying no token was delivered: %+v", read)
	}
}

func TestAnAnswerRoundTrips(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	if code := held.send(t, surface.Sent{
		Decisions: []surface.Choice{
			{ID: "c-1", Chosen: []int{0}, Words: "the wait is worth it"},
			{ID: "c-2", Chosen: []int{0, 2}},
		},
		Message: "one more look please",
	}); code != http.StatusNoContent {
		t.Fatalf("sending an answer was answered with %d", code)
	}

	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}

	want := surface.Answered{
		Decisions: []surface.Decided{
			{ID: "c-1", Chosen: []string{"take the lock first"}, Commands: []string{"patch c-1"}, Words: "the wait is worth it"},
			{ID: "c-2", Chosen: []string{"the write", "the log line"}, Commands: []string{"keep c-2 write", "keep c-2 log"}},
		},
		Message: "one more look please",
	}
	if diff := cmp.Diff(want, read); diff != "" {
		t.Errorf("what came back differs (-want +got):\n%s", diff)
	}
	if read.Ended {
		t.Error("an answer sent without ending reports the exchange as over")
	}
}

func TestMultipleChoiceCarriesEveryTickInDocumentOrder(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-2", Chosen: []int{2, 0, 1}}}})

	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}
	want := []string{"the write", "the metric", "the log line"}
	if diff := cmp.Diff(want, read.Decisions[0].Chosen); diff != "" {
		t.Errorf("what a multiple-choice decision carried differs (-want +got):\n%s", diff)
	}
}

func TestTheCommandIsWhatComesBack(t *testing.T) {
	read := twoDecisions()
	held := serving(t, read, 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}}}})

	answered, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}

	stated := read.Decisions[0].Options[1].Command
	if diff := cmp.Diff([]string{stated}, answered.Decisions[0].Commands); diff != "" {
		t.Errorf("the command that came back is not the option's own string (-document +delivered):\n%s", diff)
	}
}

func TestWordsAloneAreACompleteAnswer(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Words: "none of these is the question"}}})

	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("words alone were sent and nothing was delivered")
	}

	decided := read.Decisions[0]
	if len(decided.Chosen) != 0 || len(decided.Commands) != 0 {
		t.Errorf("an answer with no option chosen carries %v", decided)
	}
	if decided.Words != "none of these is the question" {
		t.Errorf("the words a person wrote came back as %q", decided.Words)
	}
}

func TestEndingIsDistinguishable(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}, Ends: true})

	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}
	if !read.Ended {
		t.Error("an answer sent with ending does not report the exchange as over")
	}

	if code := held.send(t, surface.Sent{Message: "one more thing"}); code != http.StatusConflict {
		t.Errorf("a closed exchange accepted a further answer with %d", code)
	}
	if further, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("a closed exchange delivered a further answer: %+v", further)
	}
}

func TestDeliveryHappensOnce(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})

	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}
	if again, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("an answer already delivered was delivered again: %+v", again)
	}
}

func TestAPollThatStartedBeforeThePageWasOpenStillWaits(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	waiting := make(chan surface.Answered, 1)
	go func() {
		read, arrived := held.await(t, 10*time.Second)
		if arrived {
			waiting <- read
		}
		close(waiting)
	}()

	time.Sleep(100 * time.Millisecond)
	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})

	read, arrived := <-waiting
	if !arrived {
		t.Fatal("a wait that started before anything was sent failed rather than waiting")
	}
	if read.Decisions[0].ID != "c-1" {
		t.Errorf("what the wait returned is %+v", read)
	}
}

func TestTheSurfaceRunsNothingAndWritesOnlyItsOwnState(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	before := within(t, held.home)
	held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}}}})
	if _, arrived := held.await(t, 5*time.Second); !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}

	if diff := cmp.Diff(before, within(t, held.home)); diff != "" {
		t.Errorf("an answer arriving wrote outside the surface's own state (-before +after):\n%s", diff)
	}
}

func within(t *testing.T, root string) []string {
	t.Helper()

	var held []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		held = append(held, relative)
		return nil
	})
	if err != nil {
		t.Fatalf("reading what the surface wrote: %v", err)
	}
	slices.Sort(held)
	return held
}
