package surface_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
)

var (
	thePicture = []byte("\x89PNG pretend these bytes are a screenshot of the cache page")
	theTrace   = []byte("goroutine 7 [running]:\nmain.readBeforeLock()\n")
)

func pictureAndTrace() map[string]carried {
	return map[string]carried{
		"file-1": {name: "before.png", kind: "image/png", held: thePicture},
		"file-2": {name: "trace.txt", kind: "text/plain", held: theTrace},
	}
}

func answerCarryingPictureAndTrace() surface.Sent {
	return surface.Sent{
		Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1"}}},
		Message:   "the trace is from the failing run",
		Attached:  []string{"file-2"},
	}
}

func collectedWith(t *testing.T, held *standing, sent surface.Sent, parts map[string]carried) surface.Answered {
	t.Helper()

	if code := held.sendCarrying(t, sent, parts); code != http.StatusNoContent {
		t.Fatalf("sending an answer carrying files was answered with %d", code)
	}
	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer carrying files was sent and nothing was delivered")
	}
	return read
}

func TestAnAnswerCarryingFilesIsCollectedWithOneAbsolutePathPerFileAndTheBytesThatWereSent(t *testing.T) {
	// scenario: A file attached to a decision and one to the message appear in what the poll prints with absolute paths whose bytes match
	held := serving(t, twoDecisions(), 0)

	read := collectedWith(t, held, answerCarryingPictureAndTrace(), pictureAndTrace())

	decided := read.Decisions[0].Attachments
	if len(decided) != 1 || len(read.Attachments) != 1 {
		t.Fatalf("the decision carries %d attachments and the message %d, want one each: %+v", len(decided), len(read.Attachments), read)
	}
	for _, expected := range []struct {
		got   surface.Attachment
		name  string
		kind  string
		bytes []byte
	}{
		{decided[0], "before.png", "image/png", thePicture},
		{read.Attachments[0], "trace.txt", "text/plain", theTrace},
	} {
		if expected.got.Name != expected.name || expected.got.Type != expected.kind || expected.got.Bytes != int64(len(expected.bytes)) {
			t.Errorf("the attachment reads %+v, want name %q, type %q and %d bytes", expected.got, expected.name, expected.kind, len(expected.bytes))
		}
		if !filepath.IsAbs(expected.got.Path) {
			t.Errorf("the path %q of %s is not absolute", expected.got.Path, expected.name)
		}
		onDisk, err := os.ReadFile(expected.got.Path)
		if err != nil {
			t.Fatalf("reading %s at the path the answer gave: %v", expected.name, err)
		}
		if !bytes.Equal(onDisk, expected.bytes) {
			t.Errorf("the bytes at %s are not the bytes that were sent", expected.got.Path)
		}
	}

	answers := filepath.Join(held.home, "sessions", held.session, "answers", "1")
	for _, path := range []string{decided[0].Path, read.Attachments[0].Path} {
		if !strings.HasPrefix(path, answers+string(filepath.Separator)) {
			t.Errorf("the file at %s is not under the folder of the answer, %s", path, answers)
		}
	}
	if want := filepath.Join(answers, "c-1", "1-before.png"); decided[0].Path != want {
		t.Errorf("the decision's file is stored at %s, want %s", decided[0].Path, want)
	}
	if want := filepath.Join(answers, "message", "1-trace.txt"); read.Attachments[0].Path != want {
		t.Errorf("the message's file is stored at %s, want %s", read.Attachments[0].Path, want)
	}
}

func TestAnAnswerCarryingFilesIsKeptWholeAndPollLastReadsItsPaths(t *testing.T) {
	// scenario: The last collected answer, files and all, is read back without the surface
	held := serving(t, twoDecisions(), 0)

	collected := collectedWith(t, held, answerCarryingPictureAndTrace(), pictureAndTrace())

	kept, found, err := surface.ReadAnswer(held.session)
	if err != nil || !found {
		t.Fatalf("reading the last collected answer: found %v, %v", found, err)
	}
	if diff := cmp.Diff(collected, kept); diff != "" {
		t.Errorf("the answer the session kept differs from the one handed back (-collected +kept):\n%s", diff)
	}
}

func TestAnAnswerWithNoAttachmentsIsPrintedExactlyAsItWasBeforeAttachmentsExisted(t *testing.T) {
	// scenario: A poll of an answer with no attachments prints exactly what it prints today
	held := serving(t, twoDecisions(), 0)

	if code := held.send(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Words: "fine"}}, Message: "ok"}); code != http.StatusNoContent {
		t.Fatalf("sending an answer was answered with %d", code)
	}
	read, arrived := held.await(t, 5*time.Second)
	if !arrived {
		t.Fatal("an answer was sent and nothing was delivered")
	}

	printed, err := json.Marshal(read)
	if err != nil {
		t.Fatal(err)
	}
	const before = `{"decisions":[{"id":"c-1","chosen":["take the lock first"],"commands":["patch c-1"],"words":"fine"}],"message":"ok","ended":false}`
	if string(printed) != before {
		t.Errorf("an answer with no attachments prints\n%s\nwant\n%s", printed, before)
	}
}

func TestAStatedNameHoldingASeparatorOrDotsIsStoredInsideTheAnswersFolder(t *testing.T) {
	// scenario: A stated name holding a path separator or .. is stored inside the answer's folder
	for _, hostile := range []string{"../../escape.txt", `..\..\escape.txt`, "/etc/passwd", "..", ".", "a/../../b", ".hidden", strings.Repeat("long", 200) + ".txt"} {
		t.Run(hostile, func(t *testing.T) {
			held := serving(t, twoDecisions(), 0)
			parts := map[string]carried{"file-1": {name: hostile, kind: "text/plain", held: theTrace}}

			read := collectedWith(t, held, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1"}}}}, parts)

			answers := filepath.Join(held.home, "sessions", held.session, "answers", "1")
			stored := read.Decisions[0].Attachments[0].Path
			relative, err := filepath.Rel(answers, stored)
			if err != nil || !filepath.IsLocal(relative) {
				t.Fatalf("the file stored for %q is at %s, outside the answer's folder %s", hostile, stored, answers)
			}
			if filepath.Dir(relative) != "c-1" {
				t.Errorf("the file stored for %q is at %s, and it belongs directly in the decision's folder", hostile, relative)
			}
			if onDisk, err := os.ReadFile(stored); err != nil || !bytes.Equal(onDisk, theTrace) {
				t.Errorf("the file stored for %q does not hold what was sent: %v", hostile, err)
			}
			if _, err := os.Stat(filepath.Join(held.home, "escape.txt")); err == nil {
				t.Errorf("%q wrote a file beside the surface's home", hostile)
			}
		})
	}
}

func TestAnAnswerAmendedBeforeTheCollectionLeavesOnlyTheLastAnswersFilesOnDisk(t *testing.T) {
	// scenario: Amending before the poll sends every file again, and the replaced answer's folder is gone from disk
	held := serving(t, twoDecisions(), 0)
	session := filepath.Join(held.home, "sessions", held.session)

	if code := held.sendCarrying(t, answerCarryingPictureAndTrace(), pictureAndTrace()); code != http.StatusNoContent {
		t.Fatalf("sending the first answer was answered with %d", code)
	}
	firstFiles := within(t, session)

	again := map[string]carried{"file-1": {name: "after.png", kind: "image/png", held: []byte("second")}}
	read := collectedWith(t, held, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}, Attached: []string{"file-1"}}}, Message: "amended"}, again)

	if read.Message != "amended" || len(read.Attachments) != 0 || len(read.Decisions[0].Attachments) != 1 {
		t.Fatalf("collecting after an amendment returned %+v, want the amendment alone", read)
	}
	for _, left := range within(t, session) {
		if strings.HasSuffix(left, "before.png") || strings.HasSuffix(left, "trace.txt") {
			t.Errorf("the replaced answer's file %s is still on disk, and the files first sent were %v", left, firstFiles)
		}
	}
	entries, err := os.ReadDir(session)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "answers" {
			t.Errorf("the session holds a folder %q beside answers after the only answer was collected", entry.Name())
		}
	}
}

func TestANewDocumentRemovesEveryUncollectedAnswersFolder(t *testing.T) {
	// scenario: A revision empties every tray and no uncollected answer's folder is left on disk
	held := serving(t, twoDecisions(), 0)
	session := filepath.Join(held.home, "sessions", held.session)

	if code := held.sendCarrying(t, answerCarryingPictureAndTrace(), pictureAndTrace()); code != http.StatusNoContent {
		t.Fatalf("sending an answer was answered with %d", code)
	}
	if len(within(t, session)) == 0 {
		t.Fatal("the answer's files are nowhere on disk before the revision")
	}
	before := within(t, session)

	revised := oneOtherDecision()
	if err := held.server.Revise(held.session, revised, written(t, revised), ""); err != nil {
		t.Fatalf("revising the session: %v", err)
	}

	for _, left := range within(t, session) {
		for _, gone := range before {
			if left == gone && (strings.HasSuffix(left, ".png") || strings.HasSuffix(left, ".txt")) {
				t.Errorf("the uncollected answer's file %s survived a revision", left)
			}
		}
	}
	if delivered, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("an answer left uncollected across a revision was delivered: %+v", delivered)
	}
}

func TestAnAnswerNamingAFileNoPartCarriesIsRefusedAndKeepsNothing(t *testing.T) {
	// scenario: A name the answer states that no part carries refuses the submit with nothing kept
	held := serving(t, twoDecisions(), 0)
	session := filepath.Join(held.home, "sessions", held.session)
	before := within(t, session)

	sent := surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1", "file-9"}}}}
	if code := held.sendCarrying(t, sent, map[string]carried{"file-1": {name: "a.txt", kind: "text/plain", held: theTrace}}); code != http.StatusConflict {
		t.Errorf("an answer naming a part nobody sent was answered with %d, want %d", code, http.StatusConflict)
	}

	if diff := cmp.Diff(before, within(t, session)); diff != "" {
		t.Errorf("a refused answer left files on disk (-before +after):\n%s", diff)
	}
	if left := undeliveredFolders(t, session); len(left) != 0 {
		t.Errorf("a refused answer left the folders %v", left)
	}
	if delivered, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("a refused answer was delivered: %+v", delivered)
	}
}

func TestAPartTheAnswerDoesNotNameIsRefusedAndKeepsNothing(t *testing.T) {
	// scenario: A part the answer does not name refuses the submit with nothing kept
	held := serving(t, twoDecisions(), 0)
	session := filepath.Join(held.home, "sessions", held.session)

	sent := surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1"}}}}
	parts := map[string]carried{
		"file-1": {name: "a.txt", kind: "text/plain", held: theTrace},
		"file-2": {name: "b.txt", kind: "text/plain", held: theTrace},
	}
	if code := held.sendCarrying(t, sent, parts); code != http.StatusConflict {
		t.Errorf("an answer carrying a part it does not name was answered with %d, want %d", code, http.StatusConflict)
	}

	if left := undeliveredFolders(t, session); len(left) != 0 {
		t.Errorf("a refused answer left the folders %v", left)
	}
	if delivered, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("a refused answer was delivered: %+v", delivered)
	}
}

func TestAnAnswerRefusedForItsDecisionKeepsNoFileItCarried(t *testing.T) {
	// scenario: An answer refused whole leaves no file of its own on disk
	held := serving(t, twoDecisions(), 0)
	session := filepath.Join(held.home, "sessions", held.session)

	sent := surface.Sent{Decisions: []surface.Choice{{ID: "nobody-asked", Words: "x", Attached: []string{"file-1"}}}}
	if code := held.sendCarrying(t, sent, map[string]carried{"file-1": {name: "a.txt", kind: "text/plain", held: theTrace}}); code != http.StatusConflict {
		t.Errorf("an answer naming no decision of the page was answered with %d, want %d", code, http.StatusConflict)
	}

	if left := undeliveredFolders(t, session); len(left) != 0 {
		t.Errorf("a refused answer left the folders %v", left)
	}
}

func TestEveryCollectedAnswerKeepsItsFilesUnderItsOwnNumber(t *testing.T) {
	// scenario: A path a poll printed stays readable after later rounds are collected
	held := serving(t, twoDecisions(), 0)

	first := collectedWith(t, held, answerCarryingPictureAndTrace(), pictureAndTrace())
	again := map[string]carried{"file-1": {name: "round-two.png", kind: "image/png", held: []byte("round two")}}
	second := collectedWith(t, held, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}, Attached: []string{"file-1"}}}}, again)

	for _, path := range []string{first.Decisions[0].Attachments[0].Path, first.Attachments[0].Path, second.Decisions[0].Attachments[0].Path} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("a path a poll printed is no longer readable: %v", err)
		}
	}
	if !strings.Contains(second.Decisions[0].Attachments[0].Path, filepath.Join("answers", "2")) {
		t.Errorf("the second collected answer's file is at %s, want it under answers/2", second.Decisions[0].Attachments[0].Path)
	}
	kept, _, err := surface.ReadAnswer(held.session)
	if err != nil || kept.Decisions[0].Attachments[0].Name != "round-two.png" {
		t.Errorf("the last collected answer reads %+v, %v, want the second round's", kept, err)
	}
}

func TestNoRouteOfTheSurfaceReturnsTheBytesOfAnAttachment(t *testing.T) {
	// scenario: No route of the surface returns the bytes of an attachment
	held := serving(t, twoDecisions(), 0)
	read := collectedWith(t, held, answerCarryingPictureAndTrace(), pictureAndTrace())

	session := surface.SessionPrefix + held.session
	routes := []string{
		session + "/answers/1/c-1/1-before.png",
		session + "/answers/1/message/1-trace.txt",
		session + "/attachments/1/c-1/1-before.png",
		session + "/c-1/1-before.png",
		session + "/file-1",
		session + "/illustrations/c-1/1-before.png",
		session,
		"/assets/1-before.png",
		"/attachments/" + held.session + "/1-before.png",
		strings.TrimPrefix(read.Decisions[0].Attachments[0].Path, held.home),
	}
	for _, route := range routes {
		code, said := held.fetch(t, held.at(route))
		if strings.Contains(said, string(thePicture)) || strings.Contains(said, string(theTrace)) {
			t.Errorf("GET %s returned %d and the bytes of an attachment", route, code)
		}
		if code == http.StatusOK && route != session {
			t.Errorf("GET %s was answered with %d, and no route of the surface answers there", route, code)
		}
	}
}

func TestAnAnswerThatIsNotFormDataIsRefused(t *testing.T) {
	// scenario: A body that is not multipart is refused with nothing kept
	held := serving(t, twoDecisions(), 0)

	if code := held.posting(t, held.under("/answer"), held.token, []byte(`{"decisions":[],"message":"x","ends":false}`)); code != http.StatusBadRequest {
		t.Errorf("a JSON body to the answer route was answered with %d, want %d", code, http.StatusBadRequest)
	}
	if delivered, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("a refused answer was delivered: %+v", delivered)
	}
}

func undeliveredFolders(t *testing.T, session string) []string {
	t.Helper()

	entries, err := os.ReadDir(session)
	if err != nil {
		t.Fatalf("reading the session: %v", err)
	}
	var left []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "answers" && entry.Name() != "illustrations" {
			left = append(left, entry.Name())
		}
	}
	return left
}
