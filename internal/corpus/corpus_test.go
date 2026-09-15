package corpus_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/corpus"
	"github.com/socrates-tools/socrates/internal/findings"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

// scenario: A live review's answer is kept
func TestAnEntryHoldsTheAnswerAndItsCircumstances(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("the corpus holds no recorded answer, so nothing here proves anything")
	}

	answers := filepath.Join(source.Root(t), corpus.Directory)
	for _, entry := range entries {
		t.Run(entry.Name, func(t *testing.T) {
			if err := entry.Validate(); err != nil {
				t.Errorf("the recorded answer is missing its circumstances: %v", err)
			}

			stored, err := os.ReadFile(filepath.Join(answers, entry.Name, "answer.json"))
			if err != nil {
				t.Fatalf("reading what was stored: %v", err)
			}
			if !bytes.Equal(stored, entry.Document) {
				t.Errorf("the document read back is not the document stored: %d bytes stored, %d read back", len(stored), len(entry.Document))
			}

			if _, err := findings.Parse(entry.Document); err != nil {
				t.Errorf("the recorded answer does not parse as a findings document: %v", err)
			}
		})
	}
}

func TestAnEntryThatNamesNoModelIsRefused(t *testing.T) {
	complete := corpus.Entry{
		Name:       "an-answer",
		Model:      "opus",
		Effort:     "high",
		Round:      1,
		Base:       "a9a3207995f002632584efb33312c14123de4560",
		Snapshot:   "4445282d62930b67e699165cda6a212f8c2e47ab",
		Run:        "20260820T091048Z-6616db11",
		RecordedOn: "2026-08-20",
		Subject:    "a change worth one line",
		Document:   []byte(`{"findings":[],"risk_level":"low","risk_rationale":"nothing"}`),
	}
	if err := complete.Validate(); err != nil {
		t.Fatalf("an entry carrying everything is refused: %v", err)
	}

	for absent, without := range map[string]func(corpus.Entry) corpus.Entry{
		"model":  func(e corpus.Entry) corpus.Entry { e.Model = ""; return e },
		"effort": func(e corpus.Entry) corpus.Entry { e.Effort = ""; return e },
		"round":  func(e corpus.Entry) corpus.Entry { e.Round = 0; return e },
	} {
		t.Run("without the "+absent, func(t *testing.T) {
			err := without(complete).Validate()
			if err == nil {
				t.Fatalf("an entry that names no %s is accepted", absent)
			}
			if !strings.Contains(err.Error(), absent) {
				t.Errorf("the refusal does not name what is absent: %v", err)
			}
		})
	}
}

func TestAnEntryTheCorpusCannotHoldIsRefused(t *testing.T) {
	reviewed := t.TempDir()
	writeModule(t, reviewed, source.ModulePath)

	for named, entry := range map[string]corpus.Entry{
		"a document socrates never routed":    withDocument(`{"findings":"none that matter","risk_level":"low","risk_rationale":"r"}`),
		"a document that is not JSON":         withDocument("I read the change and it looks fine"),
		"a name embedding drops":              withName(".first-look"),
		"a name embedding drops, underscored": withName("_first-look"),
		"a name that is the corpus itself":    withName("."),
		"a name that climbs out":              withName("../../etc"),
	} {
		t.Run(named, func(t *testing.T) {
			into := t.TempDir()
			if _, err := corpus.Record(into, reviewed, entry); err == nil {
				t.Fatalf("the corpus accepted %s", named)
			}

			held, err := corpus.Read(filepath.Join(into, corpus.Directory))
			if err == nil && len(held) > 0 {
				t.Errorf("the refusal left %d answers in the corpus", len(held))
			}
		})
	}
}

func TestEveryRecordedAnswerIsVisibleToTheCorpus(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}

	directories, err := os.ReadDir(filepath.Join(source.Root(t), corpus.Directory))
	if err != nil {
		t.Fatalf("reading the answers directory: %v", err)
	}

	held := 0
	for _, directory := range directories {
		if directory.IsDir() {
			held++
		}
	}
	if held != len(entries) {
		t.Errorf("%d answers are on disk and %d are visible to the corpus", held, len(entries))
	}
}

func TestRecordingFromAnotherRepositoryIsRefused(t *testing.T) {
	foreign := t.TempDir()
	writeModule(t, foreign, "github.com/somebody/else")

	into := t.TempDir()
	_, err := corpus.Record(into, foreign, recordable())
	if err == nil {
		t.Fatal("recording an answer about another repository is accepted")
	}
	if !strings.Contains(err.Error(), foreign) {
		t.Errorf("the refusal does not name the repository: %v", err)
	}

	if _, err := os.Stat(filepath.Join(into, corpus.Directory)); !os.IsNotExist(err) {
		t.Errorf("the refusal wrote into the corpus anyway: %v", err)
	}
}

func TestRecordingFromARepositoryThatDeclaresNoModuleIsRefused(t *testing.T) {
	// scenario: Recording from another repository is refused
	for _, foreign := range []struct {
		named    string
		declared string
		written  bool
	}{
		{named: "a repository that is not Go at all", written: false},
		{named: "a module file that declares nothing", declared: "go 1.25\n", written: true},
	} {
		t.Run(foreign.named, func(t *testing.T) {
			repository := t.TempDir()
			if foreign.written {
				if err := os.WriteFile(filepath.Join(repository, "go.mod"), []byte(foreign.declared), 0o644); err != nil {
					t.Fatalf("writing a module file that declares nothing: %v", err)
				}
			}

			into := t.TempDir()
			_, err := corpus.Record(into, repository, recordable())
			if err == nil {
				t.Fatal("recording an answer about a repository that is not socrates is accepted")
			}
			if !strings.Contains(err.Error(), repository) {
				t.Errorf("the refusal does not name the repository: %v", err)
			}
			for _, misreading := range []string{"no such file", "file does not exist"} {
				if strings.Contains(err.Error(), misreading) {
					t.Errorf("the refusal reads as %q rather than as a repository that is not socrates: %v", misreading, err)
				}
			}
			if !strings.Contains(err.Error(), "socrates") {
				t.Errorf("the refusal does not say what the condition is: %v", err)
			}

			if _, err := os.Stat(filepath.Join(into, corpus.Directory)); !os.IsNotExist(err) {
				t.Errorf("the refusal wrote into the corpus anyway: %v", err)
			}
		})
	}
}

func TestARecordedAnswerIsStoredByteForByte(t *testing.T) {
	reviewed := t.TempDir()
	writeModule(t, reviewed, source.ModulePath)

	into := t.TempDir()
	entry := recordable()
	at, err := corpus.Record(into, reviewed, entry)
	if err != nil {
		t.Fatalf("recording an answer: %v", err)
	}
	if !strings.HasPrefix(at, into) {
		t.Errorf("the answer was recorded outside the root it was given: %s", at)
	}

	recorded, err := corpus.Read(filepath.Join(into, corpus.Directory))
	if err != nil {
		t.Fatalf("reading back what was recorded: %v", err)
	}
	if len(recorded) != 1 {
		t.Fatalf("the corpus holds %d answers, want one", len(recorded))
	}

	if !bytes.Equal(recorded[0].Document, entry.Document) {
		t.Errorf("the document read back is not the document the reviewer returned:\nreturned %q\nread back %q", entry.Document, recorded[0].Document)
	}
	if recorded[0].Name != entry.Name || recorded[0].Model != entry.Model || recorded[0].Round != entry.Round {
		t.Errorf("the circumstances read back differ: %+v", recorded[0])
	}
}

func withDocument(document string) corpus.Entry {
	entry := recordable()
	entry.Document = []byte(document)
	return entry
}

func withName(name string) corpus.Entry {
	entry := recordable()
	entry.Name = name
	return entry
}

func recordable() corpus.Entry {
	return corpus.Entry{
		Name:       "one-live-round",
		Model:      "opus",
		Effort:     "high",
		Round:      1,
		Base:       "a9a3207995f002632584efb33312c14123de4560",
		Snapshot:   "4445282d62930b67e699165cda6a212f8c2e47ab",
		Run:        "20260820T091048Z-6616db11",
		RecordedOn: "2026-08-21",
		Subject:    "a change worth one line",
		Document:   []byte(`{"findings":[{"id":"f1","severity":"info","file":"a.go","line":1,"description":"d","action":"no-op"}],"risk_level":"low","risk_rationale":"r"}`),
	}
}

func writeModule(t *testing.T, directory, module string) {
	t.Helper()

	if err := os.WriteFile(filepath.Join(directory, "go.mod"), []byte("module "+module+"\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatalf("writing go.mod: %v", err)
	}
}

func TestALiveReviewsAnswerIsKeptWithTheVendorBesideIt(t *testing.T) {
	entries, err := corpus.Entries()
	if err != nil {
		t.Fatalf("reading the corpus: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("the corpus holds no recorded answer, so this check reads nothing")
	}

	for _, entry := range entries {
		answered := entry.Reviewer
		if answered == "" {
			answered = agent.DefaultVendor
		}
		if !agent.Supports(answered) {
			t.Errorf("%s was answered by %s, which socrates does not drive", entry.Name, answered)
		}
		if entry.Model == "" {
			t.Errorf("%s names no model beside the reviewer that answered", entry.Name)
		}
		if entry.Round < 1 {
			t.Errorf("%s names the round %d, and rounds are counted from one", entry.Name, entry.Round)
		}
	}
}
