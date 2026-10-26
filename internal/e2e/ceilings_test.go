package e2e_test

import (
	"strconv"
	"strings"
	"testing"
)

const (
	pageShowsARefusal = `document.querySelector('.state').dataset.tone === 'error'`
	whatThePageSays   = `document.querySelector('.state').textContent`
	mebibyte          = 1 << 20
)

func (r *reader) addFilesOfSize(tray string, size int, names ...string) {
	r.t.Helper()

	r.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(r.t, tray)+`);
  tray.add(`+jsonOfList(r.t, names)+`.map((name) => new File([new Uint8Array(`+strconv.Itoa(size)+`)], name, { type: 'application/octet-stream' })));
})()`, nil)
}

func TestAFileOverTwentyMebibytesRefusesThePressNamingTheCeilingAndKeepsNothing(t *testing.T) {
	// scenario: A file over 20 MiB refuses the submit, the page names the ceiling and nothing is written
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	browser.evaluate(chose("c-1", 0), nil)
	browser.addFilesOfSize(trayOfDecision("c-1"), 21*mebibyte, "huge.bin")
	browser.click(".send")
	browser.waitFor(pageShowsARefusal)

	var said string
	browser.evaluate(whatThePageSays, &said)
	if !strings.Contains(said, "20 MiB") {
		t.Errorf("the page says %q of a file over the ceiling, and does not name 20 MiB", said)
	}
	var receiptHidden bool
	browser.evaluate(`document.querySelector('.receipt').hidden`, &receiptHidden)
	if !receiptHidden {
		t.Error("the page drew a receipt for an answer the surface refused")
	}
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 0 {
		t.Errorf("a refused press left the folders %v", folders)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("a refused press was delivered: %+v", delivered)
	}

	browser.evaluate(`document.querySelector(`+jsonOf(t, trayOfDecision("c-1"))+`).clear()`, nil)
	browser.click(".send")
	browser.waitFor(`!document.querySelector('.receipt').hidden`)
	answered := held.polledAnswer()
	for _, decision := range answered.Decisions {
		if decision.ID == "c-1" {
			if len(decision.Attachments) != 0 {
				t.Errorf("pressing again without the file delivered %+v", decision.Attachments)
			}
			return
		}
	}
	t.Errorf("pressing again without the file delivered %+v, and no answer to c-1", answered)
}

func TestFilesTotallingOverOneHundredMebibytesRefuseThePressNamingTheCeiling(t *testing.T) {
	// scenario: Files totalling over 100 MiB refuse the submit and the page names the ceiling
	t.Parallel()
	held, browser := read(t, twoDecisionDocument)

	browser.evaluate(chose("c-1", 0), nil)
	browser.addFilesOfSize(trayOfDecision("c-1"), 19*mebibyte, "a.bin", "b.bin", "c.bin", "d.bin", "e.bin", "f.bin")
	browser.click(".send")
	browser.waitFor(pageShowsARefusal)

	var said string
	browser.evaluate(whatThePageSays, &said)
	if !strings.Contains(said, "100 MiB") {
		t.Errorf("the page says %q of files over the answer's ceiling, and does not name 100 MiB", said)
	}
	if folders := foldersBesideAnswers(t, held.sessionFolder()); len(folders) != 0 {
		t.Errorf("a refused press left the folders %v", folders)
	}
	if delivered, arrived := held.awaiting(shortWindow); arrived {
		t.Errorf("a refused press was delivered: %+v", delivered)
	}
}
