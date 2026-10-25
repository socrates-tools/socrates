package e2e_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/google/go-cmp/cmp"
)

const (
	trayEmptyKey     = "empty"
	trayHoldingKey   = "holding-an-image-and-a-text-file"
	trayFullKey      = "full"
	trayRefusingKey  = "refusing"
	trayDisabledKey  = "disabled"
	trayCeiling      = 10
	somePicture      = "before.png"
	someTrace        = "trace.txt"
	traceWords       = "a trace somebody captured\n"
	aFileNamedByType = `new File(['x'], 'n.png', { type: 'image/png' })`
)

type drawnChip struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Thumbnail *string `json:"thumbnail"`
	Badge     *string `json:"badge"`
}

type drawnTray struct {
	Chips   []drawnChip `json:"chips"`
	Counter string      `json:"counter"`
	Refusal *string     `json:"refusal"`
	Files   []string    `json:"files"`
	Max     int         `json:"max"`
}

func trayOf(key string) string {
	return stateOf("ui-file-tray", key) + " ui-file-tray"
}

func (r *reader) trayDrawn(selector string) drawnTray {
	r.t.Helper()

	var drawn drawnTray
	r.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(r.t, selector)+`);
  return {
    chips: tray.drawnChips.map((chip) => ({ name: chip.name, type: chip.type, thumbnail: chip.thumbnail, badge: chip.badge })),
    counter: tray.drawnCounter,
    refusal: tray.drawnRefusal,
    files: tray.files.map((file) => file.name),
    max: tray.max,
  };
})()`, &drawn)
	return drawn
}

func (r *reader) listenToTray(selector string) {
	r.t.Helper()

	r.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(r.t, selector)+`);
  window.__announced = { change: 0, refused: 0 };
  tray.addEventListener('change', () => { window.__announced.change += 1; });
  tray.addEventListener('refused', () => { window.__announced.refused += 1; });
  return true;
})()`, nil)
}

func (r *reader) announcedByTray() map[string]int {
	r.t.Helper()

	var announced map[string]int
	r.evaluate(`window.__announced`, &announced)
	return announced
}

func (r *reader) addToTray(selector string, names ...string) []string {
	r.t.Helper()

	var refused []string
	r.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(r.t, selector)+`);
  const files = `+jsonOfList(r.t, names)+`.map((name) => new File(['held'], name, { type: name.endsWith('.png') ? 'image/png' : 'text/plain' }));
  return tray.add(files).map((file) => file.name);
})()`, &refused)
	return refused
}

func (r *reader) pointAt(selector, expression string) (float64, float64) {
	r.t.Helper()

	var point struct{ X, Y float64 }
	r.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(r.t, selector)+`);
  tray.scrollIntoView({ block: 'center' });
  const at = `+expression+`;
  return { X: at.x, Y: at.y };
})()`, &point)
	return point.X, point.Y
}

func (r *reader) choosingTheseFiles(paths ...string) {
	r.t.Helper()

	opened := chromedp.Events(r.context, page.FileChooserOpened)
	go func() {
		for held, err := range opened {
			if err != nil {
				return
			}
			_ = chromedp.Do(r.context, sending(dom.SetFileInputFiles, dom.SetFileInputFilesParams{Files: paths, BackendNodeID: held.BackendNodeID}))
		}
	}()
	r.do(sending(page.SetInterceptFileChooserDialog, page.SetInterceptFileChooserDialogParams{Enabled: true}))
}

func aPictureAndATraceOnDisk(t *testing.T) (string, string) {
	t.Helper()

	folder := t.TempDir()
	picture := image.NewRGBA(image.Rect(0, 0, 4, 4))
	for x := range 4 {
		for y := range 4 {
			picture.Set(x, y, color.RGBA{R: 200, G: 40, B: 40, A: 255})
		}
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, picture); err != nil {
		t.Fatal(err)
	}
	pictureAt := filepath.Join(folder, somePicture)
	traceAt := filepath.Join(folder, someTrace)
	if err := os.WriteFile(pictureAt, encoded.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(traceAt, []byte(traceWords), 0o600); err != nil {
		t.Fatal(err)
	}
	return pictureAt, traceAt
}

func TestTheCatalogDrawsTheFileTrayEmptyHoldingAnImageAndANonImageFullAtItsCeilingAndRefusing(t *testing.T) {
	// scenario: The catalog draws the file tray empty, holding an image and a non-image, full at its ceiling and refusing
	t.Parallel()
	page := catalogued(t)

	empty := page.trayDrawn(trayOf(trayEmptyKey))
	if len(empty.Chips) != 0 || empty.Counter != "0 / 10" || empty.Refusal != nil {
		t.Errorf("the empty state draws %+v, want no chip, the counter 0 / 10 and no refusal", empty)
	}

	holding := page.trayDrawn(trayOf(trayHoldingKey))
	if len(holding.Chips) != 2 || holding.Counter != "2 / 10" || holding.Refusal != nil {
		t.Fatalf("the holding state draws %+v, want two chips, the counter 2 / 10 and no refusal", holding)
	}
	if holding.Chips[0].Thumbnail == nil || !strings.HasPrefix(*holding.Chips[0].Thumbnail, "blob:") || holding.Chips[0].Badge != nil {
		t.Errorf("the image chip draws %+v, want a thumbnail from a blob: URL and no badge", holding.Chips[0])
	}
	if holding.Chips[1].Thumbnail != nil || holding.Chips[1].Badge == nil || *holding.Chips[1].Badge == "" {
		t.Errorf("the non-image chip draws %+v, want a type badge and no thumbnail", holding.Chips[1])
	}

	full := page.trayDrawn(trayOf(trayFullKey))
	if len(full.Chips) != trayCeiling || full.Counter != "10 / 10" || full.Refusal != nil {
		t.Errorf("the full state draws %d chips, the counter %q and the refusal %v, want ten, 10 / 10 and none", len(full.Chips), full.Counter, full.Refusal)
	}

	refusing := page.trayDrawn(trayOf(trayRefusingKey))
	if len(refusing.Chips) != trayCeiling || refusing.Refusal == nil || !strings.Contains(*refusing.Refusal, "10") {
		t.Errorf("the refusing state draws %d chips and the refusal %v, want ten chips and a refusal naming 10", len(refusing.Chips), refusing.Refusal)
	}
	if got := page.violations(); len(got) != 0 {
		t.Errorf("the catalog's trays drew against the page's policy: %v", got)
	}
}

func TestAttachingAnImageByTheButtonShowsItsThumbnailAndAttachingATextFileShowsATypeBadge(t *testing.T) {
	// scenario: Attaching an image shows its thumbnail, and attaching a text file shows a type badge
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayEmptyKey)
	pictureAt, traceAt := aPictureAndATraceOnDisk(t)
	page.listenToTray(tray)
	page.choosingTheseFiles(pictureAt, traceAt)

	x, y := page.pointAt(tray, `tray.attachAt`)
	page.clickAt(x, y)
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).files.length === 2`)

	drawn := page.trayDrawn(tray)
	if diff := cmp.Diff([]string{somePicture, someTrace}, drawn.Files); diff != "" {
		t.Errorf("the tray holds the files (-want +got):\n%s", diff)
	}
	if drawn.Chips[0].Thumbnail == nil || !strings.HasPrefix(*drawn.Chips[0].Thumbnail, "blob:") {
		t.Errorf("the attached image draws %+v, want a thumbnail from a blob: URL", drawn.Chips[0])
	}
	if drawn.Chips[1].Badge == nil || *drawn.Chips[1].Badge != "TXT" || drawn.Chips[1].Thumbnail != nil {
		t.Errorf("the attached text file draws %+v, want the badge TXT and no thumbnail", drawn.Chips[1])
	}
	if drawn.Counter != "2 / 10" {
		t.Errorf("the tray counts %q, want 2 / 10", drawn.Counter)
	}
	if got := page.announcedByTray()["change"]; got < 1 {
		t.Errorf("the tray announced change %d times, want at least once", got)
	}
	var loaded bool
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).thumbnailsLoaded`, &loaded)
	if !loaded {
		t.Error("the thumbnail did not load from its blob: URL")
	}
	if got := page.violations(); len(got) != 0 {
		t.Errorf("attaching drew against the page's policy: %v", got)
	}
}

func TestRemovingAChipWithItsControlRemovesThatFileFromTheTrayAndAnnouncesChange(t *testing.T) {
	// scenario: Removing a chip removes that file from files and announces change
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayHoldingKey)
	page.listenToTray(tray)
	before := page.trayDrawn(tray)

	x, y := page.pointAt(tray, `tray.drawnChips[0].removeAt`)
	page.clickAt(x, y)
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).files.length === 1`)

	after := page.trayDrawn(tray)
	if diff := cmp.Diff(before.Files[1:], after.Files); diff != "" {
		t.Errorf("the tray holds after removing the first chip (-want +got):\n%s", diff)
	}
	if len(after.Chips) != 1 || after.Counter != "1 / 10" {
		t.Errorf("the tray draws %d chips and the counter %q, want one and 1 / 10", len(after.Chips), after.Counter)
	}
	if got := page.announcedByTray()["change"]; got != 1 {
		t.Errorf("the tray announced change %d times for one removal, want 1", got)
	}
}

func TestAddingAnElevenFileToATrayOfTenLeavesItHoldingTenAnnouncesRefusedAndShowsTheRefusalWhereItWasAdded(t *testing.T) {
	// scenario: Adding an eleventh file to a tray of ten leaves it holding ten, announces refused, and shows the refusal naming the ceiling where it was added
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayFullKey)
	page.listenToTray(tray)

	refused := page.addToTray(tray, "eleventh.txt")

	if diff := cmp.Diff([]string{"eleventh.txt"}, refused); diff != "" {
		t.Errorf("add returned the refused files (-want +got):\n%s", diff)
	}
	drawn := page.trayDrawn(tray)
	if len(drawn.Files) != trayCeiling || len(drawn.Chips) != trayCeiling {
		t.Errorf("the tray holds %d files and draws %d chips, want ten of each", len(drawn.Files), len(drawn.Chips))
	}
	announced := page.announcedByTray()
	if announced["refused"] != 1 || announced["change"] != 0 {
		t.Errorf("the tray announced %v, want refused once and no change", announced)
	}
	if drawn.Refusal == nil || !strings.Contains(*drawn.Refusal, "10") {
		t.Errorf("the tray shows the refusal %v, want one naming 10", drawn.Refusal)
	}
	var inside bool
	page.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(t, tray)+`);
  const box = tray.getBoundingClientRect();
  const at = tray.refusalAt;
  return at.x >= box.left && at.x <= box.right && at.y >= box.top && at.y <= box.bottom;
})()`, &inside)
	if !inside {
		t.Error("the refusal is drawn outside the tray it refused for")
	}
}

func TestAddingMoreFilesThanTheTrayHasRoomForTakesWhatFitsAndRefusesTheRest(t *testing.T) {
	// scenario: An empty tray given twelve files takes ten and refuses two
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayEmptyKey)
	names := make([]string, 0, 12)
	for at := range 12 {
		names = append(names, string(rune('a'+at))+".txt")
	}

	refused := page.addToTray(tray, names...)

	if diff := cmp.Diff(names[trayCeiling:], refused); diff != "" {
		t.Errorf("add refused (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(names[:trayCeiling], page.trayDrawn(tray).Files); diff != "" {
		t.Errorf("the tray holds (-want +got):\n%s", diff)
	}
}

func TestClearingATrayEmptiesItAndAnnouncesChangeAndTakesTheRefusalAway(t *testing.T) {
	// scenario: Clearing a refusing tray empties it, announces change and draws no refusal
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayRefusingKey)
	page.listenToTray(tray)

	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).clear(); true`, nil)

	drawn := page.trayDrawn(tray)
	if len(drawn.Files) != 0 || len(drawn.Chips) != 0 || drawn.Refusal != nil || drawn.Counter != "0 / 10" {
		t.Errorf("the cleared tray draws %+v, want nothing held, no refusal and 0 / 10", drawn)
	}
	if got := page.announcedByTray()["change"]; got != 1 {
		t.Errorf("the tray announced change %d times for one clear, want 1", got)
	}
}

func TestDroppingFilesOnTheTrayAttachesThem(t *testing.T) {
	// scenario: Files dropped on the tray are attached to it
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayEmptyKey)
	page.listenToTray(tray)

	page.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(t, tray)+`);
  const transfer = new DataTransfer();
  transfer.items.add(`+aFileNamedByType+`);
  tray.dispatchEvent(new DragEvent('drop', { dataTransfer: transfer, bubbles: true, cancelable: true, composed: true }));
  return true;
})()`, nil)

	if diff := cmp.Diff([]string{"n.png"}, page.trayDrawn(tray).Files); diff != "" {
		t.Errorf("the tray holds after a drop (-want +got):\n%s", diff)
	}
	if got := page.announcedByTray()["change"]; got != 1 {
		t.Errorf("the tray announced change %d times for one drop, want 1", got)
	}
}

func TestADisabledTrayTakesNoDropAndItsControlsCannotBePressed(t *testing.T) {
	// scenario: A disabled tray takes no dropped file and leaves its attach and remove controls unpressable
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayDisabledKey)
	page.listenToTray(tray)

	var disabled, controlsDisabled bool
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).disabled`, &disabled)
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).drawnControlsDisabled`, &controlsDisabled)
	if !disabled || !controlsDisabled {
		t.Fatalf("the disabled state reads disabled %v with its controls disabled %v, want both", disabled, controlsDisabled)
	}

	before := page.trayDrawn(tray).Files
	page.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(t, tray)+`);
  const transfer = new DataTransfer();
  transfer.items.add(`+aFileNamedByType+`);
  tray.dispatchEvent(new DragEvent('drop', { dataTransfer: transfer, bubbles: true, cancelable: true, composed: true }));
  return true;
})()`, nil)

	if diff := cmp.Diff(before, page.trayDrawn(tray).Files); diff != "" {
		t.Errorf("a disabled tray took a dropped file (-before +after):\n%s", diff)
	}
	if got := page.announcedByTray()["change"]; got != 0 {
		t.Errorf("a disabled tray announced change %d times for a drop", got)
	}
}
