package agent

import (
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates-review/internal/runlog"
)

func opencodeUnderTest(t *testing.T, recorder runlog.Recorder) *opencodeAgent {
	t.Helper()

	vendor, known := Named("opencode")
	if !known {
		t.Fatal("opencode is not in the statement")
	}
	built, err := newOpencode(vendor, Options{
		Recorder: recorder,
		Routing:  Configuration{Reviewer: {Model: "anthropic/claude-sonnet-4-5", Effort: "high"}},
	})
	if err != nil {
		t.Fatalf("constructing the opencode agent: %v", err)
	}
	return built.(*opencodeAgent)
}

func TestTheOpencodeInvocationStatesTheDirectoryItIsToRead(t *testing.T) {
	// scenario: The repository does not instruct its own reviewer
	built := opencodeUnderTest(t, nil)
	worktree := t.TempDir()

	arguments := built.arguments(Routing{Model: "anthropic/claude-sonnet-4-5", Effort: "high"}, Request{Dir: worktree, Prompt: "review this"})

	if arguments[0] != "run" {
		t.Errorf("the invocation begins %q", arguments[0])
	}
	format := slices.Index(arguments, "--format")
	if format < 0 || arguments[format+1] != "json" {
		t.Errorf("the invocation does not ask for the event stream: %v", arguments)
	}

	directory := slices.Index(arguments, "--dir")
	if directory < 0 || arguments[directory+1] != worktree {
		t.Errorf("the invocation does not state the snapshot worktree, and opencode does not take the repository from the process's working directory: %v", arguments)
	}

	model := slices.Index(arguments, "--model")
	if model < 0 || arguments[model+1] != "anthropic/claude-sonnet-4-5" {
		t.Errorf("the role's model does not reach the invocation: %v", arguments)
	}
	variant := slices.Index(arguments, "--variant")
	if variant < 0 || arguments[variant+1] != "high" {
		t.Errorf("the role's effort does not reach the invocation: %v", arguments)
	}

	for _, stated := range arguments {
		if strings.Contains(stated, "review this") {
			t.Errorf("the invocation carries the prompt in its arguments, and a prompt holding a diff does not fit there: %v", arguments)
		}
	}
}

func TestARoleWithNoEffortToState(t *testing.T) {
	built := opencodeUnderTest(t, nil)

	arguments := built.arguments(Routing{Model: "anthropic/claude-sonnet-4-5"}, Request{Dir: t.TempDir(), Prompt: "review this"})

	if slices.Contains(arguments, "--variant") {
		t.Errorf("the invocation states an effort where the role states none: %v", arguments)
	}
	model := slices.Index(arguments, "--model")
	if model < 0 || arguments[model+1] != "anthropic/claude-sonnet-4-5" {
		t.Errorf("the model is not what the role asked for once the effort is absent: %v", arguments)
	}
}

func TestTheOpencodeStreamIsRead(t *testing.T) {
	completed := strings.Join([]string{
		`{"type":"step_start","timestamp":1,"sessionID":"ses_1","part":{"type":"step-start"}}`,
		`{"type":"text","timestamp":2,"sessionID":"ses_1","part":{"type":"text","text":"I read the change.\n\n"}}`,
		"{\"type\":\"text\",\"timestamp\":3,\"sessionID\":\"ses_1\",\"part\":{\"type\":\"text\",\"text\":\"```json\\n{}\\n```\"}}",
		`{"type":"step_finish","timestamp":4,"sessionID":"ses_1","part":{"type":"step-finish","reason":"stop","tokens":{"total":19182,"input":17366,"output":24,"reasoning":0,"cache":{"write":3,"read":1792}},"cost":0}}`,
		"",
	}, "\n")

	read, err := parseOpencodeStream(strings.NewReader(completed))
	if err != nil {
		t.Fatalf("reading a stream opencode wrote: %v", err)
	}
	if !read.completed {
		t.Error("a run that finished did not read as completed")
	}
	if read.sessionID != "ses_1" {
		t.Errorf("the session reads %q", read.sessionID)
	}
	if !strings.Contains(read.text, "```json") {
		t.Errorf("the final message reads %q, and the parts are joined in the order they arrived", read.text)
	}
	want := Usage{InputTokens: 17366, OutputTokens: 24, CacheReadTokens: 1792, CacheCreationTokens: 3}
	if read.usage != want {
		t.Errorf("what it cost reads %+v, want %+v", read.usage, want)
	}
}

func TestAnOpencodeRunThatFailedCarriesItsAccountAndNoTerminalEvent(t *testing.T) {
	failed := `{"type":"error","timestamp":1,"sessionID":"ses_1","error":{"name":"UnknownError","data":{"message":"Unexpected server error. Check server logs for details.","ref":"err_1"}}}` + "\n"

	read, err := parseOpencodeStream(strings.NewReader(failed))
	if err != nil {
		t.Fatalf("reading a stream opencode wrote: %v", err)
	}
	if read.completed {
		t.Error("a run that failed read as completed, and a failing run emits no step_finish")
	}
	if !strings.Contains(read.account, "UnknownError") || !strings.Contains(read.account, "Unexpected server error") {
		t.Errorf("the account reads %q, without what the error event said", read.account)
	}
}

func TestASessionThatCannotSayWhatItHeld(t *testing.T) {
	for _, vendor := range Supported() {
		if vendor.ReportsItsToolInventory {
			continue
		}

		kept := &keptEvents{}
		recordSessionHeldNothingStated(vendor, kept, Request{Role: Reviewer}, true)

		if !kept.holds(EventToolsUnreported) {
			t.Errorf("a round under %s completed and the run does not record that the session did not say what it held", vendor.Name)
		}
		for _, event := range kept.events {
			if event.Kind == EventUnexpectedTools {
				t.Errorf("a round under %s reads as a session holding only what socrates asked for, and %s reports no inventory", vendor.Name, vendor.Name)
			}
		}
	}
}

func TestASessionThatSaidNothingRecordsNothingWhereItDidNotComplete(t *testing.T) {
	for _, vendor := range Supported() {
		if vendor.ReportsItsToolInventory {
			continue
		}

		kept := &keptEvents{}
		recordSessionHeldNothingStated(vendor, kept, Request{Role: Reviewer}, false)

		if len(kept.events) != 0 {
			t.Errorf("a round under %s that did not complete recorded %v about what the session held", vendor.Name, kept.events)
		}
	}
}

type keptEvents struct {
	events []runlog.Event
}

func (k *keptEvents) Record(event runlog.Event) error {
	k.events = append(k.events, event)
	return nil
}

func (k *keptEvents) holds(kind string) bool {
	for _, event := range k.events {
		if event.Kind == kind {
			return true
		}
	}
	return false
}
