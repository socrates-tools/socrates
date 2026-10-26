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
	"github.com/chromedp/cdproto/input"
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

type drawnTile struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Thumbnail   *string `json:"thumbnail"`
	Badge       *string `json:"badge"`
	Removable   bool    `json:"removable"`
	RemoveShown bool    `json:"removeShown"`
}

type drawnAttach struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type drawnTray struct {
	Tiles   []drawnTile  `json:"tiles"`
	Attach  *drawnAttach `json:"attach"`
	Refusal *string      `json:"refusal"`
	Files   []string     `json:"files"`
	Max     int          `json:"max"`
	Text    string       `json:"text"`
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
    tiles: tray.drawnTiles.map((tile) => ({ name: tile.name, type: tile.type, thumbnail: tile.thumbnail, badge: tile.badge, removable: tile.removeAt !== null, removeShown: tile.removeShown })),
    attach: tray.drawnAttach,
    refusal: tray.drawnRefusal,
    files: tray.files.map((file) => file.name),
    max: tray.max,
    text: tray.drawnText,
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
	if len(empty.Tiles) != 0 || empty.Refusal != nil {
		t.Errorf("the empty state draws %+v, want no tile and no refusal", empty)
	}
	if empty.Attach == nil || empty.Attach.Name != "Attach files" || empty.Attach.Title != "Attach files" {
		t.Errorf("the empty state draws the attach control %+v, want one named and titled Attach files", empty.Attach)
	}
	if strings.Contains(empty.Text, "/") || strings.Contains(empty.Text, "Drop") {
		t.Errorf("the empty state draws the text %q, want no counter and no hint", empty.Text)
	}

	holding := page.trayDrawn(trayOf(trayHoldingKey))
	if len(holding.Tiles) != 2 || holding.Refusal != nil {
		t.Fatalf("the holding state draws %+v, want two tiles and no refusal", holding)
	}
	if holding.Tiles[0].Thumbnail == nil || !strings.HasPrefix(*holding.Tiles[0].Thumbnail, "blob:") || holding.Tiles[0].Badge != nil {
		t.Errorf("the image tile draws %+v, want a thumbnail from a blob: URL and no badge", holding.Tiles[0])
	}
	if holding.Tiles[1].Thumbnail != nil || holding.Tiles[1].Badge == nil || *holding.Tiles[1].Badge == "" {
		t.Errorf("the non-image tile draws %+v, want a type badge and no thumbnail", holding.Tiles[1])
	}

	full := page.trayDrawn(trayOf(trayFullKey))
	if len(full.Tiles) != trayCeiling || full.Refusal != nil {
		t.Errorf("the full state draws %d tiles and the refusal %v, want ten and none", len(full.Tiles), full.Refusal)
	}

	refusing := page.trayDrawn(trayOf(trayRefusingKey))
	if len(refusing.Tiles) != trayCeiling || refusing.Refusal == nil || !strings.Contains(*refusing.Refusal, "10") {
		t.Errorf("the refusing state draws %d tiles and the refusal %v, want ten tiles and a refusal naming 10", len(refusing.Tiles), refusing.Refusal)
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
	if drawn.Tiles[0].Thumbnail == nil || !strings.HasPrefix(*drawn.Tiles[0].Thumbnail, "blob:") {
		t.Errorf("the attached image draws %+v, want a thumbnail from a blob: URL", drawn.Tiles[0])
	}
	if drawn.Tiles[1].Badge == nil || *drawn.Tiles[1].Badge != "TXT" || drawn.Tiles[1].Thumbnail != nil {
		t.Errorf("the attached text file draws %+v, want the badge TXT and no thumbnail", drawn.Tiles[1])
	}
	if got := page.announcedByTray()["change"]; got < 1 {
		t.Errorf("the tray announced change %d times, want at least once", got)
	}
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).thumbnailsLoaded`)
	if got := page.violations(); len(got) != 0 {
		t.Errorf("attaching drew against the page's policy: %v", got)
	}
}

func TestRemovingATileWithItsControlRemovesThatFileFromTheTrayAndAnnouncesChange(t *testing.T) {
	// scenario: Removing a tile removes that file from files and announces change
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayHoldingKey)
	page.listenToTray(tray)
	before := page.trayDrawn(tray)

	x, y := page.pointAt(tray, `tray.drawnTiles[0].removeAt`)
	page.clickAt(x, y)
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).files.length === 1`)

	after := page.trayDrawn(tray)
	if diff := cmp.Diff(before.Files[1:], after.Files); diff != "" {
		t.Errorf("the tray holds after removing the first tile (-want +got):\n%s", diff)
	}
	if len(after.Tiles) != 1 {
		t.Errorf("the tray draws %d tiles, want one", len(after.Tiles))
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
	if len(drawn.Files) != trayCeiling || len(drawn.Tiles) != trayCeiling {
		t.Errorf("the tray holds %d files and draws %d tiles, want ten of each", len(drawn.Files), len(drawn.Tiles))
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
	if len(drawn.Files) != 0 || len(drawn.Tiles) != 0 || drawn.Refusal != nil {
		t.Errorf("the cleared tray draws %+v, want nothing held and no refusal", drawn)
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

func TestADisabledTrayTakesNoDropAndDrawsNoAttachOrRemoveControl(t *testing.T) {
	// scenario: A disabled tray takes no dropped file and draws no attach or remove control
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayDisabledKey)
	page.listenToTray(tray)

	var disabled bool
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).disabled`, &disabled)
	drawnDisabled := page.trayDrawn(tray)
	if !disabled || drawnDisabled.Attach != nil || len(drawnDisabled.Tiles) == 0 {
		t.Fatalf("the disabled state reads disabled %v and draws %+v, want it disabled, holding tiles and drawing no attach control", disabled, drawnDisabled)
	}
	for _, tile := range drawnDisabled.Tiles {
		if tile.Removable {
			t.Errorf("the disabled state draws a remove control on %s", tile.Name)
		}
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

	refused := page.addToTray(tray, "pasted-in.png")
	if diff := cmp.Diff([]string{"pasted-in.png"}, refused); diff != "" {
		t.Errorf("add on a disabled tray returned (-want +got):\n%s", diff)
	}
	drawn := page.trayDrawn(tray)
	if diff := cmp.Diff(before, drawn.Files); diff != "" {
		t.Errorf("add attached a file to a disabled tray (-before +after):\n%s", diff)
	}
	if drawn.Refusal != nil {
		t.Errorf("a disabled tray drew the refusal %q for a file it never took", *drawn.Refusal)
	}
	if announced := page.announcedByTray(); announced["change"] != 0 || announced["refused"] != 0 {
		t.Errorf("a disabled tray announced %v for an add", announced)
	}
}

func (r *reader) hoverAt(x, y float64) {
	r.t.Helper()
	r.do(sending(input.DispatchMouseEvent, input.DispatchMouseEventParams{Type: input.DispatchMouseEventTypeMouseMoved, X: x, Y: y}))
}

func (r *reader) dragFilesOver(selector, kind string) {
	r.t.Helper()

	r.evaluate(`(() => {
  const transfer = new DataTransfer();
  transfer.items.add(`+aFileNamedByType+`);
  document.querySelector(`+jsonOf(r.t, selector)+`).dispatchEvent(new DragEvent(`+jsonOf(r.t, kind)+`, { dataTransfer: transfer, bubbles: true, cancelable: true, composed: true }));
  return true;
})()`, nil)
}

func (r *reader) waitForTheTraysBorder(selector string, accent bool) {
	r.t.Helper()

	comparison := "!=="
	if accent {
		comparison = "==="
	}
	r.waitFor(`getComputedStyle(document.querySelector(` + jsonOf(r.t, selector) + `)).borderTopColor ` + comparison + ` ` + jsonOf(r.t, r.pageColourOf("--accent")))
}

func TestATileIsASquareAndAnImageThumbnailFillsIt(t *testing.T) {
	// scenario: A tile is a square and an image's thumbnail fills it
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayHoldingKey)

	var tiles []struct {
		Width, Height, ThumbnailWidth, ThumbnailHeight *float64
	}
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).drawnTiles.map((tile) => ({ Width: tile.width, Height: tile.height, ThumbnailWidth: tile.thumbnailWidth, ThumbnailHeight: tile.thumbnailHeight }))`, &tiles)

	if len(tiles) != 2 {
		t.Fatalf("the tray draws %d tiles, want 2", len(tiles))
	}
	for at, tile := range tiles {
		if *tile.Width != *tile.Height || *tile.Width < 48 {
			t.Errorf("tile %d is %vx%v, want a square of a size worth looking at", at, *tile.Width, *tile.Height)
		}
	}
	if tiles[0].ThumbnailWidth == nil || *tiles[0].ThumbnailWidth < *tiles[0].Width-2 || *tiles[0].ThumbnailHeight < *tiles[0].Height-2 {
		t.Errorf("the image's thumbnail is %v wide and %v high in a tile %v square, want it to fill the tile", tiles[0].ThumbnailWidth, tiles[0].ThumbnailHeight, *tiles[0].Width)
	}
}

func TestATilesRemoveControlIsHiddenAtRestAndShownOnHoverAndOnFocus(t *testing.T) {
	// scenario: A tile's remove control is hidden at rest and shown on hover and on focus
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayHoldingKey)

	if page.trayDrawn(tray).Tiles[0].RemoveShown {
		t.Fatal("the remove control is shown on a tile nothing points at")
	}

	x, y := page.pointAt(tray, `tray.drawnTiles[0].removeAt`)
	page.hoverAt(x, y)
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).drawnTiles[0].removeShown`)

	page.hoverAt(1, 1)
	page.waitFor(`!document.querySelector(` + jsonOf(t, tray) + `).drawnTiles[0].removeShown`)

	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).shadowRoot.querySelector('li button').focus(); true`, nil)
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).drawnTiles[0].removeShown`)
}

func TestFilesDraggedOverATrayHighlightItsOwnBorderUntilTheyLeaveOrDrop(t *testing.T) {
	// scenario: Files dragged over a tray highlight its own border until they leave or are dropped
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayEmptyKey)

	page.waitForTheTraysBorder(tray, false)

	page.dragFilesOver(tray, "dragenter")
	page.waitFor(`document.querySelector(` + jsonOf(t, tray) + `).draggedOver`)
	page.waitForTheTraysBorder(tray, true)

	page.dragFilesOver(tray, "dragleave")
	page.waitFor(`!document.querySelector(` + jsonOf(t, tray) + `).draggedOver`)
	page.waitForTheTraysBorder(tray, false)

	page.dragFilesOver(tray, "dragenter")
	page.waitForTheTraysBorder(tray, true)
	page.dragFilesOver(tray, "drop")
	page.waitFor(`!document.querySelector(` + jsonOf(t, tray) + `).draggedOver`)
	page.waitForTheTraysBorder(tray, false)
}

func TestADisabledTrayIsNotHighlightedByFilesDraggedOverIt(t *testing.T) {
	// scenario: A disabled tray is not highlighted by files dragged over it
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayDisabledKey)

	page.dragFilesOver(tray, "dragenter")

	var draggedOver bool
	page.evaluate(`document.querySelector(`+jsonOf(t, tray)+`).draggedOver`, &draggedOver)
	if draggedOver {
		t.Error("a disabled tray is marked as dragged over")
	}
	page.waitForTheTraysBorder(tray, false)
}

func TestTheTrayDrawnAroundAFieldIsTheOnlyBorderAndTilesSitAboveTheTextAndThePlusBelowItAtTheLeft(t *testing.T) {
	// scenario: The tray's border is the only border, the tiles sit above the text and the plus below it at the left inside the border
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayHoldingKey)

	var laid struct {
		TrayLeft, TrayTop, TrayRight, TrayBottom float64
		TrayBorder                               string
		FieldBorder                              string
		FieldBackground                          string
		TilesBottom                              float64
		TextTop, TextBottom                      float64
		PlusLeft, PlusTop, PlusRight, PlusBottom float64
		HoldsTheField                            bool
	}
	page.evaluate(`(() => {
  const tray = document.querySelector(`+jsonOf(t, tray)+`);
  tray.scrollIntoView({ block: 'center' });
  const box = tray.getBoundingClientRect();
  const field = tray.querySelector('ui-field');
  const text = field.shadowRoot.querySelector('textarea');
  const textBox = text.getBoundingClientRect();
  const tile = tray.shadowRoot.querySelector('li').getBoundingClientRect();
  const plus = tray.shadowRoot.querySelector('[data-attach]').getBoundingClientRect();
  return {
    TrayLeft: box.left, TrayTop: box.top, TrayRight: box.right, TrayBottom: box.bottom,
    TrayBorder: getComputedStyle(tray).borderTopWidth,
    FieldBorder: getComputedStyle(text).borderTopWidth,
    FieldBackground: getComputedStyle(text).backgroundColor,
    TilesBottom: tile.bottom,
    TextTop: textBox.top, TextBottom: textBox.bottom,
    PlusLeft: plus.left, PlusTop: plus.top, PlusRight: plus.right, PlusBottom: plus.bottom,
    HoldsTheField: field !== null,
  };
})()`, &laid)

	if !laid.HoldsTheField {
		t.Fatal("the tray draws no field inside it")
	}
	if laid.TrayBorder != "1px" || laid.FieldBorder != "0px" {
		t.Errorf("the tray's border is %s and the field's is %s, want 1px and 0px", laid.TrayBorder, laid.FieldBorder)
	}
	if laid.FieldBackground != "rgba(0, 0, 0, 0)" {
		t.Errorf("the field paints %s inside the tray, want nothing of its own", laid.FieldBackground)
	}
	if laid.TilesBottom > laid.TextTop {
		t.Errorf("the tiles end at %v, below where the text starts at %v", laid.TilesBottom, laid.TextTop)
	}
	if laid.PlusTop < laid.TextBottom {
		t.Errorf("the plus starts at %v, above where the text ends at %v", laid.PlusTop, laid.TextBottom)
	}
	if laid.PlusLeft < laid.TrayLeft || laid.PlusBottom > laid.TrayBottom || laid.PlusRight > laid.TrayLeft+(laid.TrayRight-laid.TrayLeft)/2 {
		t.Errorf("the plus sits at %v-%v, want inside the border and in its left half %v-%v", laid.PlusLeft, laid.PlusRight, laid.TrayLeft, laid.TrayRight)
	}
}

func TestFocusingTheFieldInsideATrayRingsTheTraysBorder(t *testing.T) {
	// scenario: Focusing the field inside a tray draws the focus ring on the tray's border
	t.Parallel()
	page := catalogued(t)
	tray := trayOf(trayEmptyKey)

	x, y := page.pointAt(tray, `(() => { const box = tray.querySelector('ui-field').getBoundingClientRect(); return { x: box.left + box.width / 2, y: box.top + box.height / 2 }; })()`)
	page.clickAt(x, y)
	page.waitFor(`getComputedStyle(document.querySelector(` + jsonOf(t, tray) + `)).boxShadow !== 'none'`)
	page.waitForTheTraysBorder(tray, true)
}
