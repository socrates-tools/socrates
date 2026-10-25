package e2e_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/agora/surface"
	"github.com/socrates-tools/socrates/internal/exit"
)

const composerTray = ".composer .tray"

func trayOfDecision(identifier string) string {
	return `.decision[data-decision="` + identifier + `"] .tray`
}

type picker struct {
	mu    sync.Mutex
	paths []string
}

func (p *picker) offer(paths ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paths = paths
}

func (p *picker) offered() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.paths
}

func (r *reader) pickingFiles() *picker {
	r.t.Helper()

	chosen := &picker{}
	opened := chromedp.Events(r.context, page.FileChooserOpened)
	go func() {
		for held, err := range opened {
			if err != nil {
				return
			}
			_ = chromedp.Do(r.context, sending(dom.SetFileInputFiles, dom.SetFileInputFilesParams{Files: chosen.offered(), BackendNodeID: held.BackendNodeID}))
		}
	}()
	r.do(sending(page.SetInterceptFileChooserDialog, page.SetInterceptFileChooserDialogParams{Enabled: true}))
	return chosen
}

func (r *reader) attachByTheButton(chosen *picker, tray string, paths ...string) {
	r.t.Helper()

	var before int
	r.evaluate(`document.querySelector(`+jsonOf(r.t, tray)+`).files.length`, &before)
	chosen.offer(paths...)
	x, y := r.pointAt(tray, `tray.attachAt`)
	r.clickAt(x, y)
	r.waitFor(fmt.Sprintf(`document.querySelector(%s).files.length === %d`, jsonOf(r.t, tray), before+len(paths)))
}

func (s *surfaced) polledAnswer() surface.Answered {
	s.t.Helper()

	printed, failures, code := s.poll()
	if code != exit.NothingAsked {
		s.t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	var answered surface.Answered
	if err := json.Unmarshal([]byte(printed), &answered); err != nil {
		s.t.Fatalf("reading what agora poll printed: %v\n%s", err, printed)
	}
	return answered
}

func (s *surfaced) sessionFolder() string {
	return filepath.Join(s.surfaceHome(), "sessions", s.opened)
}

func foldersBesideAnswers(t *testing.T, session string) []string {
	t.Helper()

	entries, err := os.ReadDir(session)
	if err != nil {
		t.Fatalf("reading the session: %v", err)
	}
	var folders []string
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "answers" {
			folders = append(folders, entry.Name())
		}
	}
	return folders
}

func TestAFileAttachedByTheButtonToADecisionAndOneToTheMessageComeBackFromPollWithAbsolutePathsHoldingTheirBytes(t *testing.T) {
	// scenario: A file attached to a decision and one to the message, sent with a press, appear in what poll prints with absolute paths whose bytes match
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.evaluate(chose("c-1", 0), nil)
	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.attachByTheButton(chosen, composerTray, traceAt)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	answered := held.polledAnswer()

	var decided []surface.Attachment
	for _, decision := range answered.Decisions {
		if decision.ID == "c-1" {
			decided = decision.Attachments
		}
	}
	pictureBytes, err := os.ReadFile(pictureAt)
	if err != nil {
		t.Fatal(err)
	}
	traceBytes, err := os.ReadFile(traceAt)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []struct {
		got   []surface.Attachment
		name  string
		kind  string
		bytes []byte
	}{
		{decided, somePicture, "image/png", pictureBytes},
		{answered.Attachments, someTrace, "text/plain", traceBytes},
	} {
		if len(expected.got) != 1 {
			t.Fatalf("poll printed %d attachments for %s, want one: %+v", len(expected.got), expected.name, answered)
		}
		attachment := expected.got[0]
		if attachment.Name != expected.name || attachment.Type != expected.kind || attachment.Bytes != int64(len(expected.bytes)) {
			t.Errorf("poll printed the attachment %+v, want %s of type %s and %d bytes", attachment, expected.name, expected.kind, len(expected.bytes))
		}
		if !filepath.IsAbs(attachment.Path) {
			t.Errorf("poll printed the path %q for %s, and it is not absolute", attachment.Path, expected.name)
		}
		onDisk, err := os.ReadFile(attachment.Path)
		if err != nil {
			t.Fatalf("reading %s at the path poll printed: %v", expected.name, err)
		}
		if diff := cmp.Diff(string(expected.bytes), string(onDisk)); diff != "" {
			t.Errorf("the bytes at %s are not what was attached (-attached +stored):\n%s", attachment.Path, diff)
		}
	}

	if got := browser.violations(); len(got) != 0 {
		t.Errorf("attaching and sending drew against the page's policy: %v", got)
	}

	printedAgain, failures, code := held.poll("--last")
	if code != exit.NothingAsked {
		t.Fatalf("agora poll --last exited %d: %s", code, failures)
	}
	var again surface.Answered
	if err := json.Unmarshal([]byte(printedAgain), &again); err != nil {
		t.Fatalf("reading what agora poll --last printed: %v", err)
	}
	if diff := cmp.Diff(answered, again); diff != "" {
		t.Errorf("poll --last printed other paths than the poll that collected them (-collected +last):\n%s", diff)
	}
}

func TestAnAnswerSentWithNoFileIsPolledWithNoAttachmentsFieldAtAll(t *testing.T) {
	// scenario: A poll of an answer with no attachments prints exactly what it prints today
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	browser.evaluate(chose("c-1", 0), nil)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	printed, failures, code := held.poll()
	if code != exit.NothingAsked {
		t.Fatalf("agora poll exited %d: %s", code, failures)
	}
	if strings.Contains(printed, "attachments") {
		t.Errorf("an answer sent with no file printed an attachments field:\n%s", printed)
	}
}

func TestAmendingBeforeThePollSendsEveryFileAgainAndTheReplacedAnswersFolderIsGoneFromDisk(t *testing.T) {
	// scenario: Amending before the poll sends every file again, and the replaced answer's folder is gone from disk
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.evaluate(chose("c-1", 0), nil)
	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 1 {
		t.Fatalf("the uncollected answer's files are in %v, want one folder beside answers", folders)
	}

	browser.click(".amend")
	browser.waitFor(`document.querySelector('.receipt').hidden`)
	browser.attachByTheButton(chosen, composerTray, traceAt)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	answered := held.polledAnswer()

	var decided []surface.Attachment
	for _, decision := range answered.Decisions {
		if decision.ID == "c-1" {
			decided = decision.Attachments
		}
	}
	if len(decided) != 1 || decided[0].Name != somePicture || len(answered.Attachments) != 1 || answered.Attachments[0].Name != someTrace {
		t.Errorf("the amended answer carries %+v, want the picture sent again on c-1 and the trace on the message", answered)
	}
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 0 {
		t.Errorf("the replaced answer's folder is still on disk: %v", folders)
	}
	var stored int
	if err := filepath.WalkDir(filepath.Join(held.sessionFolder(), "answers"), func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			stored++
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if stored != 3 {
		t.Errorf("the session keeps %d files for the one collected answer, want answer.json and the two attachments", stored)
	}
}

func TestANewDocumentEmptiesEveryTrayAndLeavesNoUncollectedAnswersFolderOnDisk(t *testing.T) {
	// scenario: A revision empties every tray, and no uncollected answer's folder is left on disk
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.evaluate(chose("c-1", 0), nil)
	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)
	browser.click(".amend")
	browser.waitFor(`document.querySelector('.receipt').hidden`)
	browser.attachByTheButton(chosen, composerTray, traceAt)
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 1 {
		t.Fatalf("the uncollected answer's files are in %v, want one folder beside answers", folders)
	}

	held.serveAgain(reusingIdentifierReply)
	browser.waitFor(`document.querySelector('.heading h1').textContent === 'an unrelated question reusing the identifier'`)

	var heldByEachTray []int
	browser.evaluate(`Array.from(document.querySelectorAll('.tray')).map((tray) => tray.files.length)`, &heldByEachTray)
	if diff := cmp.Diff([]int{0, 0}, heldByEachTray); diff != "" {
		t.Errorf("the trays of the revised page hold files (-want +got):\n%s", diff)
	}
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 0 {
		t.Errorf("an uncollected answer's folder survived a revision: %v", folders)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("an answer left uncollected across a revision was delivered: %+v", delivered)
	}
}

func TestAClosedExchangeDrawsNoTray(t *testing.T) {
	// scenario: A closed exchange draws no tray
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	var drawnWhileOpen []bool
	browser.evaluate(drawnTrays, &drawnWhileOpen)
	if diff := cmp.Diff([]bool{true, true, true}, drawnWhileOpen); diff != "" {
		t.Fatalf("an open exchange draws its trays (-want +got):\n%s", diff)
	}

	browser.evaluate(chose("c-1", 0), nil)
	browser.click(".finish")
	browser.waitFor(`document.querySelector('.send').disabled === true`)

	var drawnAfterFinishing []bool
	browser.evaluate(drawnTrays, &drawnAfterFinishing)
	if diff := cmp.Diff([]bool{false, false, false}, drawnAfterFinishing); diff != "" {
		t.Errorf("a page whose exchange was finished draws trays (-want +got):\n%s", diff)
	}

	held.end(held.opened)
	address := held.record().Address
	reopened := opening(t, "http://"+address+surface.SessionPrefix+held.opened)
	var drawnWhenLoadedOver []bool
	reopened.evaluate(drawnTrays, &drawnWhenLoadedOver)
	if diff := cmp.Diff([]bool{false, false, false}, drawnWhenLoadedOver); diff != "" {
		t.Errorf("a page loaded after its exchange was over draws trays (-want +got):\n%s", diff)
	}
}

const drawnTrays = `Array.from(document.querySelectorAll('.tray')).map((tray) => tray.getBoundingClientRect().width > 0 && getComputedStyle(tray).display !== 'none')`

func TestAStatedNameHoldingAPathSeparatorOrDotsIsStoredInsideTheAnswersFolder(t *testing.T) {
	// scenario: A stated name holding a path separator or .. is stored inside the answer's folder
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	browser.evaluate(chose("c-1", 0), nil)
	browser.addToTray(trayOfDecision("c-1"), "../../escape.png", `..\..\escape.txt`)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	answered := held.polledAnswer()

	answers := filepath.Join(held.sessionFolder(), "answers", "1")
	var stored []surface.Attachment
	for _, decision := range answered.Decisions {
		if decision.ID == "c-1" {
			stored = decision.Attachments
		}
	}
	if len(stored) != 2 {
		t.Fatalf("poll printed %d attachments, want the two that were attached: %+v", len(stored), answered)
	}
	for _, attachment := range stored {
		relative, err := filepath.Rel(answers, attachment.Path)
		if err != nil || !filepath.IsLocal(relative) {
			t.Errorf("%q is stored at %s, outside the answer's folder %s", attachment.Name, attachment.Path, answers)
		}
		if _, err := os.Stat(attachment.Path); err != nil {
			t.Errorf("%q is not at the path poll printed: %v", attachment.Name, err)
		}
	}
	for _, escaped := range []string{filepath.Join(held.surfaceHome(), "escape.png"), filepath.Join(filepath.Dir(held.surfaceHome()), "escape.png"), filepath.Join(held.sessionFolder(), "escape.png")} {
		if _, err := os.Stat(escaped); err == nil {
			t.Errorf("a stated name wrote %s", escaped)
		}
	}
}
