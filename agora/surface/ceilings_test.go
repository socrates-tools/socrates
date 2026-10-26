package surface_test

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
)

const mebibyte = 1 << 20

func aFileOf(size int) carried {
	return carried{name: "big.bin", kind: "application/octet-stream", held: bytes.Repeat([]byte{'x'}, size)}
}

func filesAttachedTo(count int) (surface.Sent, map[string]carried) {
	names := make([]string, count)
	parts := map[string]carried{}
	for at := range names {
		names[at] = fmt.Sprintf("file-%d", at+1)
		parts[names[at]] = carried{name: fmt.Sprintf("trace-%d.txt", at+1), kind: "text/plain", held: theTrace}
	}
	return surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: names}}}, parts
}

func TestTheCeilingsAreTwentyMebibytesAFileOneHundredAnAnswerAndTenAFilesATray(t *testing.T) {
	// scenario: The ceilings the surface states are 20 MiB a file, 100 MiB an answer and ten files a tray
	if surface.FileCeiling != 20*mebibyte || surface.AnswerCeiling != 100*mebibyte || surface.TrayCeiling != 10 {
		t.Errorf("the ceilings are %d, %d and %d, want 20 MiB, 100 MiB and 10", surface.FileCeiling, surface.AnswerCeiling, surface.TrayCeiling)
	}
	if surface.BodyCeiling != mebibyte {
		t.Errorf("the body ceiling of the document and control routes is %d, want 1 MiB", surface.BodyCeiling)
	}
}

func TestAFileOfExactlyTwentyMebibytesIsAcceptedAndOneByteMoreIsRefusedNamingTheCeiling(t *testing.T) {
	// scenario: A file over 20 MiB refuses the submit and the answer names the ceiling
	held := serving(t, twoDecisions(), 0)
	sent := surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1"}}}}

	code, said := held.sendCarryingSaid(t, sent, map[string]carried{"file-1": aFileOf(surface.FileCeiling + 1)})
	if code != http.StatusRequestEntityTooLarge {
		t.Errorf("a file one byte over the ceiling was answered with %d, want %d", code, http.StatusRequestEntityTooLarge)
	}
	if !strings.Contains(said, "20 MiB") {
		t.Errorf("the refusal of a file over the ceiling reads %q, and does not name 20 MiB", said)
	}
	if delivered, arrived := held.await(t, 200*time.Millisecond); arrived {
		t.Errorf("a refused answer was delivered: %+v", delivered)
	}

	read := collectedWith(t, held, sent, map[string]carried{"file-1": aFileOf(surface.FileCeiling)})
	if got := read.Decisions[0].Attachments[0].Bytes; got != surface.FileCeiling {
		t.Errorf("a file of exactly the ceiling was stored as %d bytes, want %d", got, surface.FileCeiling)
	}
}

func TestFilesTotallingOverOneHundredMebibytesRefuseTheSubmitNamingTheCeiling(t *testing.T) {
	// scenario: Files each under 20 MiB totalling over 100 MiB refuse the submit and the answer names the ceiling
	held := serving(t, twoDecisions(), 0)
	names := []string{"file-1", "file-2", "file-3", "file-4", "file-5", "file-6"}
	parts := map[string]carried{}
	for _, name := range names {
		parts[name] = aFileOf(19 * mebibyte)
	}

	code, said := held.sendCarryingSaid(t, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: names}}}, parts)

	if code != http.StatusRequestEntityTooLarge {
		t.Errorf("files totalling 114 MiB were answered with %d, want %d", code, http.StatusRequestEntityTooLarge)
	}
	if !strings.Contains(said, "100 MiB") {
		t.Errorf("the refusal of an answer over the ceiling reads %q, and does not name 100 MiB", said)
	}
	if left := undeliveredFolders(t, filepath.Join(held.home, "sessions", held.session)); len(left) != 0 {
		t.Errorf("a refused answer left the folders %v", left)
	}
}

func TestAnEleventhFileInOneTrayIsRefusedAndATenthIsNot(t *testing.T) {
	// scenario: A request naming an eleventh file in one tray is refused, in a decision and in the message
	held := serving(t, twoDecisions(), 0)

	tenSent, tenParts := filesAttachedTo(surface.TrayCeiling)
	if code := held.sendCarrying(t, tenSent, tenParts); code != http.StatusNoContent {
		t.Fatalf("ten files in one tray were answered with %d", code)
	}

	elevenSent, elevenParts := filesAttachedTo(surface.TrayCeiling + 1)
	code, said := held.sendCarryingSaid(t, elevenSent, elevenParts)
	if code == http.StatusNoContent {
		t.Error("eleven files in one decision's tray were accepted")
	}
	if !strings.Contains(said, "10") {
		t.Errorf("the refusal of an eleventh file reads %q, and does not name the ceiling of 10", said)
	}

	_, messageParts := filesAttachedTo(surface.TrayCeiling + 1)
	inTheMessage := surface.Sent{Message: "many files", Attached: elevenSent.Decisions[0].Attached}
	if code := held.sendCarrying(t, inTheMessage, messageParts); code == http.StatusNoContent {
		t.Error("eleven files in the message's tray were accepted")
	}
}

func TestTenFilesInEachOfTwoTraysAreAcceptedTogether(t *testing.T) {
	// scenario: The ceiling of ten holds a tray, not an answer
	held := serving(t, twoDecisions(), 0)
	parts := map[string]carried{}
	var inDecision, inMessage []string
	for at := 1; at <= 2*surface.TrayCeiling; at++ {
		name := fmt.Sprintf("file-%d", at)
		parts[name] = carried{name: name + ".txt", kind: "text/plain", held: theTrace}
		if at <= surface.TrayCeiling {
			inDecision = append(inDecision, name)
			continue
		}
		inMessage = append(inMessage, name)
	}

	read := collectedWith(t, held, surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: inDecision}}, Attached: inMessage}, parts)

	if len(read.Decisions[0].Attachments) != surface.TrayCeiling || len(read.Attachments) != surface.TrayCeiling {
		t.Errorf("two full trays were collected as %d and %d files", len(read.Decisions[0].Attachments), len(read.Attachments))
	}
}

func TestARefusedSubmitLeavesTheUncollectedAnswerAndTheDiskAsTheyWere(t *testing.T) {
	// scenario: After any refusal the exchange's uncollected answer and the disk are as they were before the submit
	overTheEleventh, elevenParts := filesAttachedTo(surface.TrayCeiling + 1)
	refusals := map[string]struct {
		sent  surface.Sent
		parts map[string]carried
	}{
		"a file over twenty mebibytes": {
			sent:  surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{1}, Attached: []string{"file-1"}}}},
			parts: map[string]carried{"file-1": aFileOf(surface.FileCeiling + 1)},
		},
		"an eleventh file in a tray": {sent: overTheEleventh, parts: elevenParts},
		"a part the answer does not name": {
			sent:  surface.Sent{Message: "other"},
			parts: map[string]carried{"file-1": aFileOf(10)},
		},
		"a name no part carries": {
			sent:  surface.Sent{Message: "other", Attached: []string{"file-1"}},
			parts: nil,
		},
	}
	for name, refusal := range refusals {
		t.Run(name, func(t *testing.T) {
			held := serving(t, twoDecisions(), 0)
			session := filepath.Join(held.home, "sessions", held.session)
			if code := held.sendCarrying(t, answerCarryingPictureAndTrace(), pictureAndTrace()); code != http.StatusNoContent {
				t.Fatalf("sending the first answer was answered with %d", code)
			}
			before := within(t, session)
			foldersBefore := undeliveredFolders(t, session)

			if code := held.sendCarrying(t, refusal.sent, refusal.parts); code == http.StatusNoContent {
				t.Fatalf("%s was accepted", name)
			}

			if diff := cmp.Diff(before, within(t, session)); diff != "" {
				t.Errorf("a refused submit changed the disk (-before +after):\n%s", diff)
			}
			if diff := cmp.Diff(foldersBefore, undeliveredFolders(t, session)); diff != "" {
				t.Errorf("a refused submit changed the folders of the session (-before +after):\n%s", diff)
			}
			read, arrived := held.await(t, 5*time.Second)
			if !arrived || read.Message != answerCarryingPictureAndTrace().Message || len(read.Attachments) != 1 || len(read.Decisions[0].Attachments) != 1 {
				t.Errorf("after a refused submit the uncollected answer reads %+v, want the first answer whole", read)
			}
		})
	}
}

func TestTheAnswerRouteAcceptsWhatTheDocumentRouteRefusesOverOneMebibyte(t *testing.T) {
	// scenario: A body over 1 MiB is accepted at the answer route when it is files and refused at the document route
	held := serving(t, twoDecisions(), 0)
	sent := surface.Sent{Decisions: []surface.Choice{{ID: "c-1", Chosen: []int{0}, Attached: []string{"file-1"}}}}

	read := collectedWith(t, held, sent, map[string]carried{"file-1": aFileOf(2 * mebibyte)})
	if read.Decisions[0].Attachments[0].Bytes != 2*mebibyte {
		t.Errorf("a 2 MiB file was stored as %d bytes", read.Decisions[0].Attachments[0].Bytes)
	}

	code := held.posting(t, held.at(surface.SessionPrefix+held.session+"/document"), held.token, aDocumentLargerThan(surface.BodyCeiling))
	if code != http.StatusRequestEntityTooLarge {
		t.Errorf("a document over 1 MiB was answered with %d, want %d", code, http.StatusRequestEntityTooLarge)
	}
}

func TestTheAnswerPartItselfIsHeldToTheBodyCeilingAndNotTheFileCeiling(t *testing.T) {
	// scenario: The JSON of the answer is refused over 1 MiB however small the files are
	held := serving(t, twoDecisions(), 0)

	code, said := held.sendCarryingSaid(t, surface.Sent{Message: strings.Repeat("x", surface.BodyCeiling+1)}, nil)

	if code != http.StatusRequestEntityTooLarge {
		t.Errorf("an answer whose words alone pass 1 MiB was answered with %d, want %d", code, http.StatusRequestEntityTooLarge)
	}
	if !strings.Contains(said, fmt.Sprint(surface.BodyCeiling)) {
		t.Errorf("the refusal reads %q, and does not name the %d bytes", said, surface.BodyCeiling)
	}
}
