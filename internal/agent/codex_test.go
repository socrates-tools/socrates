package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/findings"
)

func codexUnderTest(t *testing.T) *codexAgent {
	t.Helper()

	vendor, known := Named("codex")
	if !known {
		t.Fatal("codex is not in the statement")
	}
	built, err := newCodex(vendor, Options{Routing: Configuration{Reviewer: {Model: "gpt-5.1-codex", Effort: "high"}}})
	if err != nil {
		t.Fatalf("constructing the codex agent: %v", err)
	}
	return built.(*codexAgent)
}

func TestTheSandboxIsStated(t *testing.T) {
	built := codexUnderTest(t)
	scratch := t.TempDir()

	arguments, _, err := built.arguments(Routing{Model: "gpt-5.1-codex", Effort: "high"}, Request{Scratch: scratch, Schema: json.RawMessage(findings.Schema)})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}

	if arguments[0] != "exec" {
		t.Errorf("the invocation begins %q, want the non-interactive command", arguments[0])
	}
	if !slices.Contains(arguments, "--json") {
		t.Errorf("the invocation does not ask for the event stream: %v", arguments)
	}

	sandbox := slices.Index(arguments, "--sandbox")
	if sandbox < 0 || arguments[sandbox+1] != "read-only" {
		t.Errorf("the invocation does not state the sandbox, and a default can change in a release: %v", arguments)
	}
}

func TestNothingAsksForTheMachinesOwnCustomisations(t *testing.T) {
	built := codexUnderTest(t)

	arguments, _, err := built.arguments(Routing{Model: "gpt-5.1-codex", Effort: "high"}, Request{Scratch: t.TempDir(), Schema: json.RawMessage(findings.Schema)})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}

	if !slices.Contains(arguments, "--ignore-user-config") {
		t.Errorf("the invocation loads the machine's own configuration: %v", arguments)
	}
	if !statesOverride(arguments, "project_doc_max_bytes=0") {
		t.Errorf("the invocation lets the reviewed repository instruct its own reviewer: %v", arguments)
	}
}

func TestTheRoleReachesTheCodexInvocation(t *testing.T) {
	built := codexUnderTest(t)

	arguments, _, err := built.arguments(Routing{Model: "gpt-5.1-codex", Effort: "high"}, Request{Scratch: t.TempDir(), Schema: json.RawMessage(findings.Schema)})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}

	model := slices.Index(arguments, "-m")
	if model < 0 || arguments[model+1] != "gpt-5.1-codex" {
		t.Errorf("the role's model does not reach the invocation: %v", arguments)
	}
	if !statesOverride(arguments, "model_reasoning_effort=high") {
		t.Errorf("the role's effort does not reach the invocation: %v", arguments)
	}
	if arguments[len(arguments)-1] != "-" {
		t.Errorf("the invocation ends %q, and the prompt is written to the process rather than put in its arguments", arguments[len(arguments)-1])
	}
}

func TestAnEffortNobodyStatedIsNotPassed(t *testing.T) {
	// scenario: A role with no effort to state
	built := codexUnderTest(t)

	arguments, _, err := built.arguments(Routing{Model: "gpt-5.1-codex"}, Request{Scratch: t.TempDir(), Schema: json.RawMessage(findings.Schema)})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}

	for _, stated := range arguments {
		if strings.HasPrefix(stated, "model_reasoning_effort=") {
			t.Errorf("the invocation states %q where the role states no effort", stated)
		}
	}
	if model := slices.Index(arguments, "-m"); model < 0 || arguments[model+1] != "gpt-5.1-codex" {
		t.Errorf("the model is not what the role asked for once the effort is absent: %v", arguments)
	}
}

func TestTheInvocationCarriesTheSchemaAsAFileUnderSocratesHome(t *testing.T) {
	built := codexUnderTest(t)
	home := t.TempDir()

	arguments, written, err := built.arguments(Routing{Model: "gpt-5.1-codex"}, Request{Scratch: home, Schema: json.RawMessage(findings.Schema)})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}

	stated := slices.Index(arguments, "--output-schema")
	if stated < 0 {
		t.Fatalf("the invocation carries no schema: %v", arguments)
	}
	path := arguments[stated+1]

	if path != written {
		t.Errorf("the invocation names %s and the schema was written to %s", path, written)
	}
	if filepath.Dir(path) != home {
		t.Errorf("the schema was written to %s, and it belongs in the round's directory under socrates' home", path)
	}

	held, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading the schema the reviewer is held to: %v", err)
	}
	if string(held) != findings.Schema {
		t.Errorf("the file holds a schema other than the one socrates validates against:\n%s", held)
	}
}

func TestNoSchemaMeansNoFile(t *testing.T) {
	built := codexUnderTest(t)

	arguments, written, err := built.arguments(Routing{Model: "gpt-5.1-codex"}, Request{Scratch: t.TempDir()})
	if err != nil {
		t.Fatalf("composing the invocation: %v", err)
	}
	if written != "" {
		t.Errorf("a call supplying no schema wrote one to %s", written)
	}
	if slices.Contains(arguments, "--output-schema") {
		t.Errorf("a call supplying no schema still names one: %v", arguments)
	}
}

func TestTheCodexStreamIsRead(t *testing.T) {
	completed := strings.Join([]string{
		`{"type":"thread.started","thread_id":"01a06ab9-098c-7050-8ce8-2f19fe60653e"}`,
		`{"type":"turn.started"}`,
		`{"type":"item.completed","item":{"id":"item_0","type":"error","message":"Model metadata not found. Defaulting to fallback metadata."}}`,
		`{"type":"item.completed","item":{"id":"item_1","type":"agent_message","text":"` + `{\"findings\":[],\"risk_level\":\"low\",\"risk_rationale\":\"nothing here\"}` + `"}}`,
		`{"type":"turn.completed","usage":{"input_tokens":1200,"cached_input_tokens":800,"cache_write_input_tokens":7,"output_tokens":42,"reasoning_output_tokens":8}}`,
		"",
	}, "\n")

	read, err := parseCodexStream(strings.NewReader(completed))
	if err != nil {
		t.Fatalf("reading a stream codex wrote: %v", err)
	}
	if !read.completed {
		t.Error("a turn that completed did not read as completed")
	}
	if read.threadID != "01a06ab9-098c-7050-8ce8-2f19fe60653e" {
		t.Errorf("the session reads %q", read.threadID)
	}
	if read.text != `{"findings":[],"risk_level":"low","risk_rationale":"nothing here"}` {
		t.Errorf("the answer reads %q", read.text)
	}
	want := Usage{InputTokens: 1200, OutputTokens: 42, CacheReadTokens: 800, CacheCreationTokens: 7}
	if read.usage != want {
		t.Errorf("what it cost reads %+v, want %+v", read.usage, want)
	}
}

func TestACodexTurnThatFailedCarriesItsAccount(t *testing.T) {
	failed := strings.Join([]string{
		`{"type":"thread.started","thread_id":"t1"}`,
		`{"type":"turn.started"}`,
		`{"type":"error","message":"Reconnecting... 1/5 (unexpected status 401 Unauthorized)"}`,
		`{"type":"turn.failed","error":{"message":"Missing environment variable: ` + "`NOTHING_IS_SET_HERE`" + `."}}`,
		"",
	}, "\n")

	read, err := parseCodexStream(strings.NewReader(failed))
	if err != nil {
		t.Fatalf("reading a stream codex wrote: %v", err)
	}
	if read.completed {
		t.Error("a turn that failed read as completed")
	}
	if !read.failed {
		t.Error("a turn that failed did not read as failed")
	}
	if !strings.Contains(read.account, "Missing environment variable") {
		t.Errorf("the account reads %q, without what the terminal event said", read.account)
	}
	if !strings.Contains(read.account, "401 Unauthorized") {
		t.Errorf("the account reads %q, without what it said on the way", read.account)
	}
}

func TestACodexStreamThatSaidNothingTerminalFails(t *testing.T) {
	for _, answered := range []struct {
		named  string
		stream string
	}{
		{named: "a stream with no terminal event", stream: `{"type":"thread.started","thread_id":"t1"}` + "\n" + `{"type":"turn.started"}` + "\n"},
		{named: "silence", stream: ""},
	} {
		t.Run(answered.named, func(t *testing.T) {
			read, err := parseCodexStream(strings.NewReader(answered.stream))
			if err != nil {
				t.Fatalf("reading the stream: %v", err)
			}
			if read.completed {
				t.Errorf("%s read as a turn that completed", answered.named)
			}
		})
	}
}

func TestATruncatedCodexStreamIsRefused(t *testing.T) {
	truncated := `{"type":"thread.started","thread_id":"t1"}` + "\n" + `{"type":"turn.comp`

	if _, err := parseCodexStream(strings.NewReader(truncated)); err == nil {
		t.Error("a stream truncated part-way through a line was read as whole")
	}
}

func TestALineThatIsNotACodexEventIsRefused(t *testing.T) {
	notAnEvent := `{"type":"thread.started","thread_id":"t1"}` + "\n" + "this is not an event at all\n" + `{"type":"turn.completed"}` + "\n"

	if _, err := parseCodexStream(strings.NewReader(notAnEvent)); err == nil {
		t.Error("a stream carrying a line that is not an event was read as whole")
	}
}

func statesOverride(arguments []string, stated string) bool {
	for position, argument := range arguments {
		if argument == "-c" && position+1 < len(arguments) && arguments[position+1] == stated {
			return true
		}
	}
	return false
}
