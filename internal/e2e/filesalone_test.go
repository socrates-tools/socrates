package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

const whatTheIndexSaysOfFirstDecision = `({
  count: document.querySelector('.answered-count').textContent,
  ticked: document.querySelector('.entry-link[data-decision="c-1"] .answered').dataset.answered,
  card: document.querySelector('.decision[data-decision="c-1"] .attribution .answered').dataset.answered,
})`

type theIndexSays struct {
	Count  string `json:"count"`
	Ticked string `json:"ticked"`
	Card   string `json:"card"`
}

func TestAFileAloneOnADecisionIsTickedInTheIndexCountedAndSentWithoutAConfirmation(t *testing.T) {
	// scenario: A decision holding only a file is ticked in the index and counted, and a press raises no confirmation naming it
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)
	pictureAt, _ := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	var before theIndexSays
	browser.evaluate(whatTheIndexSaysOfFirstDecision, &before)
	if before.Ticked != "no" || before.Count != "1" {
		t.Fatalf("a decision nobody answered is ticked %q among %q answered", before.Ticked, before.Count)
	}

	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)

	var after theIndexSays
	browser.evaluate(whatTheIndexSaysOfFirstDecision, &after)
	if diff := cmp.Diff(theIndexSays{Count: "2", Ticked: "yes", Card: "yes"}, after); diff != "" {
		t.Errorf("a decision holding only a file is not ticked and counted (-want +got):\n%s", diff)
	}

	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	var opened bool
	browser.evaluate(confirmationIsOpen, &opened)
	if opened {
		t.Error("a press with a file alone on one decision and an answer on the other raised a confirmation")
	}
}

func TestAPressWithAFileAloneOnOneDecisionAndAnotherLeftBlankNamesOnlyTheBlankOneAndPollsTheFileWithNoWords(t *testing.T) {
	// scenario: The confirmation names the blank decision and never the one holding only a file, and poll prints that decision with its attachment and no words
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)
	pictureAt, _ := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.evaluate(chose("c-2", 2), nil)
	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.click(".send")
	browser.waitFor(confirmationIsOpen)

	var named []string
	browser.evaluate(whatTheConfirmationNames, &named)
	if diff := cmp.Diff([]string{"c-2"}, named); diff != "" {
		t.Errorf("the confirmation names other decisions than the blank one (-want +got):\n%s", diff)
	}
	var heading string
	browser.evaluate(whatTheConfirmationIsHeaded, &heading)
	if !strings.Contains(heading, "1 of 2") {
		t.Errorf("the confirmation is headed %q, and it counts the decision holding a file as blank", heading)
	}

	browser.click(".send-anyway")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	answered := held.polledAnswer()
	if len(answered.Decisions) != 1 || answered.Decisions[0].ID != "c-1" {
		t.Fatalf("poll printed %+v, want only the decision holding the file", answered.Decisions)
	}
	decided := answered.Decisions[0]
	if decided.Words != "" || len(decided.Chosen) != 0 {
		t.Errorf("poll printed words %q and chosen %v for a decision holding only a file", decided.Words, decided.Chosen)
	}
	if len(decided.Attachments) != 1 || decided.Attachments[0].Name != somePicture {
		t.Errorf("poll printed the attachments %+v, want the picture alone", decided.Attachments)
	}
}

const whatTheSummaryNames = `Array.from(document.querySelectorAll('.panel .summary-row')).map((row) => ({
  said: row.querySelector('.summary-said').textContent.trim(),
  files: Array.from(row.querySelectorAll('.summary-file')).map((file) => file.textContent.trim()),
}))`

type summaryNaming struct {
	Said  string   `json:"said"`
	Files []string `json:"files"`
}

func TestTheSummaryNamesEachAttachedFileUnderItsDecisionAndUnderTheMessageBeforeAPress(t *testing.T) {
	// scenario: The summary names each attached file under its decision and under the message before a press
	t.Parallel()
	held := standing(t)
	browser := opening(t, held.serve(twoDecisionDocument), chromedp.WindowSize(aWindowWideEnoughForThePanel, aShortWindow))
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.attachByTheButton(chosen, composerTray, traceAt)

	var rows []summaryNaming
	browser.evaluate(whatTheSummaryNames, &rows)
	if len(rows) != 3 {
		t.Fatalf("the summary holds %d rows for two decisions and a message holding a file: %+v", len(rows), rows)
	}
	if diff := cmp.Diff([]string{somePicture}, rows[0].Files); diff != "" {
		t.Errorf("the summary names other files under the first decision (-want +got):\n%s", diff)
	}
	if len(rows[1].Files) != 0 {
		t.Errorf("the summary names %v under a decision that holds none", rows[1].Files)
	}
	if diff := cmp.Diff([]string{someTrace}, rows[2].Files); diff != "" {
		t.Errorf("the summary names other files under the message (-want +got):\n%s", diff)
	}
	if strings.HasSuffix(rows[0].Said, ": nothing") {
		t.Errorf("the summary says %q of a decision holding a file", rows[0].Said)
	}

	browser.evaluate(`document.querySelector(`+jsonOf(t, trayOfDecision("c-1"))+`).clear()`, nil)
	var emptied []summaryNaming
	browser.evaluate(whatTheSummaryNames, &emptied)
	if len(emptied) != 3 || len(emptied[0].Files) != 0 || !strings.HasSuffix(emptied[0].Said, ": nothing") {
		t.Errorf("the summary still names a file removed from its decision: %+v", emptied)
	}
}

func TestTheReceiptNamesEachFileThatWentUnderItsDecisionAndUnderTheMessage(t *testing.T) {
	// scenario: The receipt names each file that went, under its decision and under the message
	t.Parallel()
	_, browser := read(t, twoDecisionDocument)
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	chosen := browser.pickingFiles()

	browser.attachByTheButton(chosen, trayOfDecision("c-1"), pictureAt)
	browser.attachByTheButton(chosen, composerTray, traceAt)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)

	var named struct {
		Decision []string `json:"decision"`
		Message  []string `json:"message"`
	}
	browser.evaluate(`({
  decision: Array.from(document.querySelectorAll('.receipt-decisions .receipt-file')).map((file) => file.textContent.trim()),
  message: Array.from(document.querySelectorAll('.receipt-message-files .receipt-file')).map((file) => file.textContent.trim()),
})`, &named)
	if diff := cmp.Diff([]string{somePicture}, named.Decision); diff != "" {
		t.Errorf("the receipt names other files under the decisions (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{someTrace}, named.Message); diff != "" {
		t.Errorf("the receipt names other files under the message (-want +got):\n%s", diff)
	}
}
