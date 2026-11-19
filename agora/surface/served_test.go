package surface_test

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/document"
	"github.com/socrates-tools/socrates/agora/surface"
)

func (s *standing) servedTo(t *testing.T, session, token string) (int, []byte) {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, s.at(surface.SessionPrefix+session+surface.DocumentRoute), nil)
	if err != nil {
		t.Fatalf("addressing the surface: %v", err)
	}
	if token != "" {
		request.Header.Set(surface.TokenHeader, token)
	}
	answered, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("reaching the surface: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()

	said, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading what the surface said: %v", err)
	}
	return answered.StatusCode, said
}

func (s *standing) served(t *testing.T) (surface.Served, []byte) {
	t.Helper()

	code, said := s.servedTo(t, s.session, s.token)
	if code != http.StatusOK {
		t.Fatalf("the document was answered with %d: %s", code, said)
	}
	var read surface.Served
	if err := json.Unmarshal(said, &read); err != nil {
		t.Fatalf("reading the document the surface served: %v\n%s", err, said)
	}
	return read, said
}

func TestTheDocumentIsServedAsWrittenBesideWhatOnlyTheServerKnows(t *testing.T) {
	// scenario: A request carrying the session's token is answered with the document and the server part
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	read := document.Document{
		Title:    "two things to settle",
		Preamble: "read these first",
		Decisions: []document.Decision{
			{
				ID:       "c-1",
				Name:     "the lock",
				Severity: "warning",
				Question: "should the cache be read before the lock is taken",
				Mode:     document.One,
				Options: []document.Option{
					{Label: "take the lock first", Meaning: "slower", Command: "patch c-1", Recommended: true},
					{Label: "leave it", Command: "agree c-1"},
				},
				Evidence:     document.Evidence{{Caption: "the read", Code: "cache.Get(key)"}},
				Illustration: document.Illustration{Path: "q3.html"},
			},
			{
				ID:       "c-2",
				Question: "which of these should the write path keep",
				Mode:     document.Many,
				Options:  []document.Option{{Label: "the write", Command: "keep c-2 write"}},
			},
		},
		Source:  []document.Entry{{ID: "c-1", Name: "the lock", Severity: "warning"}},
		Settled: []document.Settled{{ID: "c-0", Question: "which one", Chosen: []string{"this"}, Words: "because"}},
	}
	held := serving(t, oneOtherDecision(), 0)
	opened, err := held.server.Open(read, written(t, read), folder)
	if err != nil {
		t.Fatalf("opening a session: %v", err)
	}
	held.session, held.token = opened.Session, opened.Token

	got, _ := held.served(t)

	want := surface.Served{
		Document: surface.ServedDocument{
			Title:    "two things to settle",
			Preamble: "read these first",
			Decisions: []surface.ServedDecision{
				{
					ID:       "c-1",
					Name:     "the lock",
					Severity: "warning",
					Question: "should the cache be read before the lock is taken",
					Mode:     document.One,
					Options: []surface.ServedOption{
						{Label: "take the lock first", Meaning: "slower", Recommended: true},
						{Label: "leave it"},
					},
					Evidence: []surface.ServedBlock{{Caption: "the read", Code: "cache.Get(key)"}},
				},
				{
					ID:       "c-2",
					Question: "which of these should the write path keep",
					Mode:     document.Many,
					Options:  []surface.ServedOption{{Label: "the write"}},
					Evidence: []surface.ServedBlock{},
				},
			},
			Source:  []surface.ServedEntry{{ID: "c-1", Name: "the lock", Severity: "warning"}},
			Settled: []surface.ServedSettled{{ID: "c-0", Question: "which one", Chosen: []string{"this"}, Words: "because"}},
		},
		Standing: surface.Standing{
			Revision: 1,
			Over:     false,
			Illustrations: map[string]string{
				"c-1": surface.SessionPrefix + held.session + "/illustrations/c-1/q3.html",
			},
			Answered:     []string{"c-1"},
			FilesPerTray: 10,
		},
	}
	if difference := cmp.Diff(want, got); difference != "" {
		t.Errorf("the document served differs from what was written (-want +got):\n%s", difference)
	}
}

func TestTheDocumentServedCarriesNoCommandAndNoLocalPath(t *testing.T) {
	// scenario: No command and no path on the caller's disk reaches the page
	folder := aFolderHolding(t, map[string]string{"q3.html": theDrawing})
	held := servingFrom(t, illustrated("q3.html"), folder)

	_, said := held.served(t)

	for _, withheld := range []string{"patch c-1", "keep c-2 metric", folder, "command"} {
		if strings.Contains(string(said), withheld) {
			t.Errorf("the document served carries %q:\n%s", withheld, said)
		}
	}
}

func TestADocumentIsRefusedWithoutThisSessionsToken(t *testing.T) {
	// scenario: No token and another session's token are refused as every guarded route refuses them
	held := serving(t, twoDecisions(), 0)
	other, err := held.server.Open(oneOtherDecision(), written(t, oneOtherDecision()), "")
	if err != nil {
		t.Fatalf("opening a second session: %v", err)
	}

	for named, token := range map[string]string{"no token": "", "another session's token": other.Token} {
		t.Run(named, func(t *testing.T) {
			code, said := held.servedTo(t, held.session, token)
			if code != http.StatusForbidden {
				t.Errorf("the document was answered with %d", code)
			}
			if strings.Contains(string(said), twoDecisions().Title) {
				t.Errorf("a refused request was told the document:\n%s", said)
			}
		})
	}

	if code, _ := held.servedTo(t, "fedcba9876543210", held.token); code != http.StatusNotFound {
		t.Errorf("the document of a session the surface does not hold was answered with %d", code)
	}
}

func TestTheDocumentAfterARevisionIsTheNewOneAtALaterRevision(t *testing.T) {
	// scenario: A revision moves the revision on and the new document is the one served
	held := serving(t, twoDecisions(), 0)
	before, _ := held.served(t)

	revised := oneOtherDecision()
	if err := held.server.Revise(held.session, revised, written(t, revised), ""); err != nil {
		t.Fatalf("revising the session: %v", err)
	}
	after, _ := held.served(t)

	if after.Document.Title != revised.Title {
		t.Errorf("the document served is titled %q, want %q", after.Document.Title, revised.Title)
	}
	if after.Standing.Revision <= before.Standing.Revision {
		t.Errorf("the revision moved from %d to %d", before.Standing.Revision, after.Standing.Revision)
	}
}

func TestEveryIllustrationAddressServedOpensThatDecisionsCopy(t *testing.T) {
	// scenario: Each address in the server part opens the drawing its decision named
	folder := aFolderHolding(t, map[string]string{
		"q3.html":           theDrawing,
		"other/second.html": `<!doctype html><title>second</title><p id="second">second drawing</p>`,
	})
	read := twoDecisions()
	read.Decisions[0].Illustration = document.Illustration{Path: filepath.Join(folder, "other", "second.html")}
	read.Decisions[1].Illustration = document.Illustration{Path: "q3.html"}
	held := serving(t, oneOtherDecision(), 0)
	opened, err := held.server.Open(read, written(t, read), folder)
	if err != nil {
		t.Fatalf("opening a session on two illustrated decisions: %v", err)
	}
	held.session, held.token = opened.Session, opened.Token

	got, _ := held.served(t)

	if len(got.Standing.Illustrations) != 2 {
		t.Fatalf("the server part holds %d illustration addresses, want 2: %v", len(got.Standing.Illustrations), got.Standing.Illustrations)
	}
	for decision, bodyHolds := range map[string]string{"c-1": "second drawing", "c-2": "$1,240.00"} {
		answered := held.responseTo(t, held.at(got.Standing.Illustrations[decision]))
		if answered.code != http.StatusOK || !strings.Contains(answered.body, bodyHolds) {
			t.Errorf("the address of %s answered %d with %q, want the drawing holding %q", decision, answered.code, answered.body, bodyHolds)
		}
	}
}

func TestAnAddressIsServedForNoDecisionWithoutAnIllustration(t *testing.T) {
	held := serving(t, twoDecisions(), 0)

	got, _ := held.served(t)

	if got.Standing.Illustrations == nil || len(got.Standing.Illustrations) != 0 {
		t.Errorf("a document with no illustration was served the addresses %#v, want an empty map", got.Standing.Illustrations)
	}
	if got.Standing.Answered == nil || len(got.Standing.Answered) != 0 {
		t.Errorf("a document recommending nothing was served the answered %#v, want an empty list", got.Standing.Answered)
	}
}

func TestWhatCrossesCarriesEveryFieldOfTheDocumentButTheCommandAndThePath(t *testing.T) {
	// scenario: A field added to the document in Go reaches the page or is named as withheld
	withheld := map[string]bool{"Option.command": true, "Decision.illustration": true}

	pairs := map[reflect.Type]reflect.Type{
		reflect.TypeFor[document.Document](): reflect.TypeFor[surface.ServedDocument](),
		reflect.TypeFor[document.Decision](): reflect.TypeFor[surface.ServedDecision](),
		reflect.TypeFor[document.Option]():   reflect.TypeFor[surface.ServedOption](),
		reflect.TypeFor[document.Block]():    reflect.TypeFor[surface.ServedBlock](),
		reflect.TypeFor[document.Entry]():    reflect.TypeFor[surface.ServedEntry](),
		reflect.TypeFor[document.Settled]():  reflect.TypeFor[surface.ServedSettled](),
	}
	for written, served := range pairs {
		carried := map[string]bool{}
		for at := range served.NumField() {
			carried[jsonNameOf(served.Field(at))] = true
		}
		for at := range written.NumField() {
			name := jsonNameOf(written.Field(at))
			if !carried[name] && !withheld[written.Name()+"."+name] {
				t.Errorf("%s.%s is in the document and neither crosses in %s nor is named as withheld", written.Name(), name, served.Name())
			}
		}
	}
}

func jsonNameOf(field reflect.StructField) string {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	return name
}
