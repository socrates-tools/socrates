package outcomes_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/corpus"
	"github.com/rui-ktei/socrates/internal/corpus/outcomes"
	"github.com/rui-ktei/socrates/internal/runstate"
	"github.com/rui-ktei/socrates/internal/testing/source"
)

const (
	escalating = `{"findings":[{"id":"f1","severity":"warning","file":"a.go","line":1,"description":"d","action":"ask-user"}],"risk_level":"medium","risk_rationale":"r"}`
	parking    = `{"findings":[{"id":"f1","severity":"warning","file":"a.go","line":1,"description":"d","action":"auto-fix"}],"risk_level":"medium","risk_rationale":"r"}`
	nothing    = `{"findings":[],"risk_level":"low","risk_rationale":"r"}`
)

func TestTheReportCountsWhatItClassified(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}

	report, err := outcomes.Of(entries)
	if err != nil {
		t.Fatalf("classifying the corpus: %v", err)
	}
	if report.Entries != len(entries) {
		t.Errorf("the report accounts for %d entries, and the corpus holds %d", report.Entries, len(entries))
	}

	reachable := outcomes.Reachable()
	if len(reachable) < 2 {
		t.Fatalf("the routing rule reaches %v, so this report says nothing", reachable)
	}

	classified := 0
	for _, reach := range report.Reaches {
		if !slices.Contains(reachable, reach.State) {
			t.Errorf("the report holds the outcome %s, which the routing rule does not reach", reach.State)
		}
		classified += len(reach.Recorded)
	}
	if classified != report.Entries {
		t.Errorf("the report classified %d answers out of %d recorded", classified, report.Entries)
	}
	if len(report.Reaches) != len(reachable) {
		t.Errorf("the report holds %d outcomes, and the routing rule reaches %d", len(report.Reaches), len(reachable))
	}

	asserting, err := outcomes.Asserting(source.Root(t))
	if err != nil {
		t.Fatalf("counting what the suite asserts: %v", err)
	}
	rendered := report.Render(asserting)
	for _, state := range reachable {
		if !strings.Contains(rendered, string(state)) {
			t.Errorf("the report does not name the outcome %s:\n%s", state, rendered)
		}
	}
	if !strings.Contains(rendered, string(runstate.Escalated)) {
		t.Errorf("the report does not name the outcome every recorded answer reached:\n%s", rendered)
	}

	named := 0
	for _, state := range reachable {
		named += len(asserting[state])
	}
	if named == 0 {
		t.Error("no test in the suite is counted as asserting any outcome, so the report's second number says nothing")
	}
}

func TestAnOutcomeNoReviewerHasReachedIsNamed(t *testing.T) {
	report, err := outcomes.Of([]corpus.Entry{
		{Name: "one", Round: 1, Document: []byte(escalating)},
		{Name: "two", Round: 1, Document: []byte(nothing)},
	})
	if err != nil {
		t.Fatalf("a corpus that reaches no parked round is refused: %v", err)
	}

	rendered := report.Render(map[runstate.State][]string{runstate.AwaitingResponse: {"TestSomethingParks"}})
	parked := lineNaming(rendered, string(runstate.AwaitingResponse))
	if parked == "" {
		t.Fatalf("the report does not name the parked outcome at all:\n%s", rendered)
	}
	if !strings.Contains(rendered, "reaches "+string(runstate.AwaitingResponse)) {
		t.Errorf("the report does not say that only hand-written scenarios reach the parked outcome:\n%s", rendered)
	}
	if !strings.Contains(parked, "0") {
		t.Errorf("the parked outcome does not carry its count of recorded answers: %q", parked)
	}
}

// scenario: A routing change that moves the corpus is visible
func TestTheRecordedDistributionMatchesTheCorpus(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}
	report, err := outcomes.Of(entries)
	if err != nil {
		t.Fatalf("classifying the corpus: %v", err)
	}

	path := filepath.Join(source.Root(t), outcomes.DistributionPath)
	recorded, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}

	if report.Distribution() != string(recorded) {
		t.Errorf("what the corpus reaches is no longer what %s records:\nrecorded:\n%s\nreached:\n%s", outcomes.DistributionPath, recorded, report.Distribution())
	}

	t.Run("A removed entry is visible", func(t *testing.T) {
		if len(entries) < 2 {
			t.Skip("the corpus holds one answer, so removing one leaves nothing to compare")
		}

		removed := entries[len(entries)-1]
		fewer, err := outcomes.Of(entries[:len(entries)-1])
		if err != nil {
			t.Fatalf("classifying a corpus with one answer removed: %v", err)
		}
		if fewer.Distribution() == string(recorded) {
			t.Fatalf("removing %s left the distribution unchanged", removed.Name)
		}
		if strings.Contains(fewer.Distribution(), removed.Name) {
			t.Errorf("the distribution still names %s after it was removed", removed.Name)
		}
	})
}

func TestALiveTaskIsSatisfiedByRecordingWhateverCameBack(t *testing.T) {
	reviewed := t.TempDir()
	if err := os.WriteFile(filepath.Join(reviewed, "go.mod"), []byte("module "+corpus.ModulePath+"\n"), 0o644); err != nil {
		t.Fatalf("writing go.mod: %v", err)
	}

	into := t.TempDir()
	unexpected := corpus.Entry{
		Name:       "a-round-that-parked",
		Model:      "opus",
		Effort:     "high",
		Round:      1,
		Base:       "a9a3207995f002632584efb33312c14123de4560",
		Snapshot:   "4445282d62930b67e699165cda6a212f8c2e47ab",
		Run:        "20260821T010203Z-6616db11",
		RecordedOn: "2026-08-21",
		Subject:    "whatever the reviewer answered",
		Document:   []byte(parking),
	}
	if _, err := corpus.Record(into, reviewed, unexpected); err != nil {
		t.Fatalf("recording an answer nobody expected: %v", err)
	}

	entries, err := corpus.Read(filepath.Join(into, corpus.Directory))
	if err != nil {
		t.Fatalf("reading back what was recorded: %v", err)
	}
	report, err := outcomes.Of(entries)
	if err != nil {
		t.Fatalf("classifying what was recorded: %v", err)
	}

	rendered := report.Render(nil)
	parked := lineNaming(rendered, string(runstate.AwaitingResponse))
	if !strings.Contains(parked, "1") {
		t.Errorf("the report does not count the answer that parked a round: %q\n%s", parked, rendered)
	}
	if strings.Contains(rendered, "reaches "+string(runstate.AwaitingResponse)) {
		t.Errorf("the report still says the parked outcome was only reached by hand:\n%s", rendered)
	}
	if !strings.Contains(rendered, unexpected.Name) {
		t.Errorf("the report does not name the recorded answer:\n%s", rendered)
	}
}

func lineNaming(rendered, state string) string {
	for line := range strings.SplitSeq(rendered, "\n") {
		if strings.HasPrefix(line, state) {
			return line
		}
	}
	return ""
}
