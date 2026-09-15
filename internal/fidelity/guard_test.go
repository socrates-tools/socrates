package fidelity_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
	"testing"

	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/corpus"
	"github.com/socrates-tools/socrates/internal/fidelity"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var behindATag = []string{"record_test.go", "criteria_test.go", "drift_test.go", "livesmoke_test.go", "forgesmoke_test.go"}

const fidelityPackage = "./internal/fidelity"

var replayingPackages = []string{"./internal/corpus", "./internal/claudesim"}

func TestTheOrdinarySuiteReachesNoNetwork(t *testing.T) {
	listed := listPackage(t, fidelityPackage)

	compiled := slices.Concat(listed.GoFiles, listed.TestGoFiles, listed.XTestGoFiles)
	for _, guarded := range behindATag {
		if slices.Contains(compiled, guarded) {
			t.Errorf("the default build compiles %s, which reaches something outside this machine", guarded)
		}
		if !slices.Contains(listed.IgnoredGoFiles, guarded) {
			t.Errorf("%s is not excluded from the default build: give it a build tag of its own", guarded)
		}
	}
}

func TestEveryFileBehindATagIsThere(t *testing.T) {
	listed := listPackage(t, fidelityPackage)

	for _, guarded := range behindATag {
		if _, err := os.Stat(guarded); err != nil {
			t.Errorf("%s is named as being behind a build tag and is not there: %v", guarded, err)
		}
	}
	if len(listed.IgnoredGoFiles) != len(behindATag) {
		t.Errorf("the default build ignores %v, and the files behind a tag are %v", listed.IgnoredGoFiles, behindATag)
	}
}

func TestReplayingARecordedAnswerNeedsNoBuildTag(t *testing.T) {
	for _, replaying := range replayingPackages {
		if ignored := listPackage(t, replaying).IgnoredGoFiles; len(ignored) > 0 {
			t.Errorf("%s keeps %v out of the default build, and replaying a recorded answer reaches nothing", replaying, ignored)
		}
	}

	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus from the ordinary suite: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("the corpus holds no recorded answer, so nothing can be replayed")
	}
}

type packageFiles struct {
	GoFiles        []string
	TestGoFiles    []string
	XTestGoFiles   []string
	IgnoredGoFiles []string
}

func listPackage(t *testing.T, named string) packageFiles {
	t.Helper()

	command := exec.Command("go", "list", "-json", named)
	command.Dir = source.Root(t)
	command.Env = append(os.Environ(), source.Toolchain()...)

	var out, failures bytes.Buffer
	command.Stdout = &out
	command.Stderr = &failures
	if err := command.Run(); err != nil {
		t.Fatalf("listing the package: %v: %s", err, failures.String())
	}

	var listed packageFiles
	if err := json.Unmarshal(out.Bytes(), &listed); err != nil {
		t.Fatalf("reading what go list said: %v", err)
	}
	return listed
}

func TestTheRecordingsCarryWhatSocratesReadsFromTheOpeningEvent(t *testing.T) {
	// scenario: What the reviewer held is read rather than assumed
	recordings, err := fidelity.Recordings(fidelity.RecordingsFor("claude"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(recordings) == 0 {
		t.Fatal("there are no recordings, so nothing pins the fields socrates depends on")
	}

	for _, recording := range recordings {
		opening, found := systemShape(recording.Shapes)
		if !found {
			t.Errorf("%s records no opening event", recording.Path)
			continue
		}
		if !slices.Contains(opening.Fields, "tools") {
			t.Errorf("%s does not record the session's tool list, and socrates reads it to know what the reviewer held: a recording taken without it would let the field disappear unnoticed", recording.Path)
		}
		if !slices.Contains(opening.Fields, "mcp_servers") {
			t.Errorf("%s does not record the sources of tools the session had connected, and socrates reads them because a source is named before its tools reach the list beside it: a recording taken without it would let the field disappear unnoticed", recording.Path)
		}
	}
}

func TestTheRecordingsCarryWhatSocratesReadsFromTheResultEvent(t *testing.T) {
	// scenario: A non-zero exit after the agent said why
	recordings, err := fidelity.Recordings(fidelity.RecordingsFor("claude"))
	if err != nil {
		t.Fatalf("%v", err)
	}
	if len(recordings) == 0 {
		t.Fatal("there are no recordings, so nothing pins the fields socrates depends on")
	}

	for _, recording := range recordings {
		terminal, found := shapeNamed(recording.Shapes, "result")
		if !found {
			t.Errorf("%s records no result event", recording.Path)
			continue
		}
		if !slices.Contains(terminal.Fields, "is_error") {
			t.Errorf("%s does not record whether the result reports an error, and socrates reads it to know a failed call from an answered one: a recording taken without it would let the field disappear unnoticed", recording.Path)
		}
		if !slices.Contains(terminal.Fields, "result") {
			t.Errorf("%s does not record the account the agent gives of its own failure, and socrates carries it as the reason a call failed: a recording taken without it would leave the reason silently falling back to the subtype", recording.Path)
		}
	}
}

func systemShape(shapes []fidelity.Shape) (fidelity.Shape, bool) {
	return shapeNamed(shapes, "system")
}

func shapeNamed(shapes []fidelity.Shape, kind string) (fidelity.Shape, bool) {
	for _, shape := range shapes {
		if shape.Type == kind {
			return shape, true
		}
	}
	return fidelity.Shape{}, false
}

func TestEverySupportedVendorIsCheckedAgainstItsOwnRecordings(t *testing.T) {
	supported := agent.Supported()
	if len(supported) == 0 {
		t.Fatal("no reviewer CLI is supported, so this check reads nothing")
	}

	for _, vendor := range supported {
		recordings, err := fidelity.Recordings(fidelity.RecordingsFor(vendor.Name))
		if err != nil {
			t.Errorf("%s is supported and has no recordings of its own: %v", vendor.Name, err)
			continue
		}

		for _, request := range fidelity.Requests() {
			if _, found := recordingNamed(recordings, request.Name); !found {
				t.Errorf("%s has no recording named %s, so the drift check would pass without comparing it", vendor.Name, request.Name)
			}
		}
	}
}

func TestWhatTheExtractingTierCostsIsRecorded(t *testing.T) {
	observed, err := fidelity.Observations(fidelity.CriteriaDirectory)
	if err != nil {
		t.Fatalf("reading what was recorded about the reviewer CLIs: %v", err)
	}

	for _, vendor := range agent.Supported() {
		if vendor.Tier != agent.Parsed {
			continue
		}

		found := false
		for _, one := range observed {
			if one.Vendor == vendor.Name && one.Criterion == theExtractingTiersCost {
				found = true
				if one.Observed["attempts"] == "" {
					t.Errorf("%s records no number of attempts, so what the extracting tier costs is asserted rather than measured", vendor.Name)
				}
			}
		}
		if !found {
			t.Errorf("%s is on the extracting tier and nothing records how many attempts a live review took", vendor.Name)
		}
	}
}

func recordingNamed(recordings []fidelity.Recording, name string) (fidelity.Recording, bool) {
	for _, recording := range recordings {
		if recording.Name == name {
			return recording, true
		}
	}
	return fidelity.Recording{}, false
}

const theExtractingTiersCost = "what-the-extracting-tier-costs"
