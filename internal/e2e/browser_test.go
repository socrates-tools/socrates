package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/log"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

const (
	browserVariable    = "AGORA_BROWSER"
	browserDeadline    = 90 * time.Second
	blockedByThePolicy = "has been blocked"
)

var browserCandidates = []string{
	"google-chrome",
	"google-chrome-stable",
	"chromium",
	"chromium-browser",
}

var installedBrowsers = []string{
	"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	"/Applications/Chromium.app/Contents/MacOS/Chromium",
}

var foundBrowser = locateBrowser()

func locateBrowser() string {
	if stated := os.Getenv(browserVariable); stated != "" {
		return stated
	}
	for _, named := range browserCandidates {
		if found, err := exec.LookPath(named); err == nil {
			return found
		}
	}
	for _, path := range installedBrowsers {
		if executable(path) {
			return path
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	for _, pattern := range []string{
		filepath.Join(home, ".cache", "ms-playwright", "chromium-*", "chrome-linux64", "chrome"),
		filepath.Join(home, ".cache", "ms-playwright", "chromium-*", "chrome-linux", "chrome"),
		filepath.Join(home, "Library", "Caches", "ms-playwright", "chromium-*", "chrome-mac*", "Google Chrome for Testing.app", "Contents", "MacOS", "Google Chrome for Testing"),
		filepath.Join(home, "Library", "Caches", "ms-playwright", "chromium-*", "chrome-mac*", "Chromium.app", "Contents", "MacOS", "Chromium"),
		filepath.Join(home, "Applications", "Google Chrome.app", "Contents", "MacOS", "Google Chrome"),
		filepath.Join(home, "Applications", "Chromium.app", "Contents", "MacOS", "Chromium"),
	} {
		matched, _ := filepath.Glob(pattern)
		for _, path := range slices.Backward(matched) {
			if executable(path) {
				return path
			}
		}
	}
	return ""
}

func executable(path string) bool {
	stated, err := os.Stat(path)
	return err == nil && !stated.IsDir() && stated.Mode().Perm()&0o111 != 0
}

func browserPath(t *testing.T) string {
	t.Helper()

	if foundBrowser == "" {
		t.Fatalf("no browser was found, and the page is driven by a real one: install Chrome or set %s to one", browserVariable)
	}
	return foundBrowser
}

type reader struct {
	t        *testing.T
	context  context.Context
	requests *recorded
	refused  *recorded
	ready    string
	painted  string
}

type recorded struct {
	mu       sync.Mutex
	reaching []string
}

func (r *recorded) note(url string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reaching = append(r.reaching, url)
}

func (r *recorded) all() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.reaching...)
}

func opening(t *testing.T, address string, options ...chromedp.ExecAllocatorOption) *reader {
	t.Helper()
	return openingFor(t, address, ".send", ".decision", options...)
}

func openingNothingToAnswer(t *testing.T, address string, options ...chromedp.ExecAllocatorOption) *reader {
	t.Helper()
	return openingFor(t, address, ".nothing-said", ".nothing-said", options...)
}

func openingFor(t *testing.T, address, ready, painted string, options ...chromedp.ExecAllocatorOption) *reader {
	t.Helper()

	profile, err := os.MkdirTemp("/tmp", "ag")
	if err != nil {
		t.Fatalf("making somewhere short for the browser to keep its profile: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(profile) })

	allocated := append([]chromedp.ExecAllocatorOption{
		chromedp.ExecPath(browserPath(t)),
		chromedp.UserDataDir(profile),
		chromedp.Env("TMPDIR=" + filepath.Dir(profile)),
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.DisableGPU,
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("use-mock-keychain", true),
		chromedp.Flag("password-store", "basic"),
	}, options...)

	allocator, releaseAllocator := chromedp.NewExecAllocator(context.Background(), allocated...)
	t.Cleanup(releaseAllocator)

	driving, releaseBrowser := chromedp.NewContext(allocator)
	t.Cleanup(releaseBrowser)

	bounded, cancel := context.WithTimeout(driving, browserDeadline)
	t.Cleanup(cancel)

	requests, refused := &recorded{}, &recorded{}
	requested := chromedp.Events(bounded, network.RequestWillBeSent)
	logged := chromedp.Events(bounded, log.EntryAdded)
	go func() {
		for held, err := range requested {
			if err != nil {
				return
			}
			requests.note(held.Request.URL)
		}
	}()
	go func() {
		for held, err := range logged {
			if err != nil {
				return
			}
			if held.Entry.Source == log.EntrySourceSecurity {
				refused.note(held.Entry.Text)
			}
		}
	}()

	held := &reader{t: t, context: bounded, requests: requests, refused: refused, ready: ready, painted: painted}
	held.do(
		sending(network.Enable, network.EnableParams{}),
		sending(log.Enable, cdp.Empty{}),
		chromedp.Navigate(address),
		chromedp.WaitReady(chromedp.CSS("body")),
		chromedp.WaitReady(chromedp.CSS(ready)),
	)
	held.settled()
	return held
}

func sending[P, R any](command cdp.Command[P, R], params P) chromedp.Action[chromedp.Void] {
	return chromedp.Func(func(ctx context.Context, target *chromedp.Target) error {
		_, err := cdp.Call(ctx, target, command, params)
		return err
	})
}

func (r *reader) do(actions ...chromedp.Action[chromedp.Void]) {
	r.t.Helper()

	if err := chromedp.Do(r.context, actions...); err != nil {
		r.t.Fatalf("driving the page: %v", err)
	}
}

func (r *reader) settled() {
	r.t.Helper()
	r.waitFor("!!document.querySelector('" + r.painted + "') && !!window.getComputedStyle(document.body).backgroundColor")
}

func (r *reader) reload() {
	r.t.Helper()

	r.do(
		chromedp.Reload(),
		chromedp.WaitReady(chromedp.CSS("body")),
		chromedp.WaitReady(chromedp.CSS(r.ready)),
	)
	r.settled()
}

func (r *reader) click(selector string) {
	r.t.Helper()
	r.do(chromedp.Click(chromedp.CSS(selector), chromedp.NodeVisible))
}

func (r *reader) clickAt(x, y float64) {
	r.t.Helper()
	r.do(chromedp.MouseClickXY(x, y))
}

func (r *reader) press(key string) {
	r.t.Helper()
	r.do(chromedp.KeyEvent(key))
}

func (r *reader) write(selector, words string) {
	r.t.Helper()
	r.do(chromedp.SendKeys(chromedp.CSS(selector), words))
}

func (r *reader) evaluate(expression string, into any, options ...chromedp.EvaluateOption) {
	r.t.Helper()

	if into == nil {
		if _, err := chromedp.Run(r.context, chromedp.Evaluate[chromedp.Void](expression, options...)); err != nil {
			r.t.Fatalf("driving the page: %v", err)
		}
		return
	}
	returned, err := chromedp.Run(r.context, chromedp.Evaluate[[]byte](expression, options...))
	if err != nil {
		r.t.Fatalf("driving the page: %v", err)
	}
	if len(returned) == 0 {
		returned = []byte("null")
	}
	if err := json.Unmarshal(returned, into); err != nil {
		r.t.Fatalf("reading what the page returned for %s: %v", expression, err)
	}
}

func awaitingThePromise(asked *runtime.EvaluateParams) {
	asked.AwaitPromise = new(true)
}

func (r *reader) waitFor(expression string) {
	r.t.Helper()

	if _, err := chromedp.Run(r.context, chromedp.Poll[bool](expression)); err != nil {
		r.t.Fatalf("waiting for %s: %v", expression, err)
	}
}

func (r *reader) screenshot() []byte {
	r.t.Helper()

	shot, err := chromedp.Run(r.context, chromedp.CaptureScreenshot())
	if err != nil {
		r.t.Fatalf("capturing what the page painted: %v", err)
	}
	return shot
}

func (r *reader) violations() []string {
	r.t.Helper()

	var held []string
	for _, said := range r.refused.all() {
		if strings.Contains(said, blockedByThePolicy) {
			held = append(held, said)
		}
	}
	return held
}

func (r *reader) emulate(preference string) {
	r.t.Helper()
	r.do(sending(emulation.SetEmulatedMedia, emulation.SetEmulatedMediaParams{Features: []*emulation.MediaFeature{
		{Name: "prefers-color-scheme", Value: preference},
	}}))
}

func (r *reader) narrow(width, height int64) {
	r.t.Helper()
	r.do(chromedp.EmulateViewport(width, height))
	r.waitFor(fmt.Sprintf(`window.innerWidth === %d && window.innerHeight === %d`, width, height))
	r.evaluate(`new Promise((settled) => requestAnimationFrame(() => requestAnimationFrame(settled)))`, nil, awaitingThePromise)
}

func chose(decision string, option int) string {
	return fmt.Sprintf(`document.querySelector('.decision[data-decision="%s"] .choice[value="%d"]').click()`, decision, option)
}

func declined(decision string) string {
	return fmt.Sprintf(`document.querySelector('.decision[data-decision="%s"] .decline').click()`, decision)
}

func within(decision, selector string) string {
	return fmt.Sprintf(`.decision[data-decision="%s"] %s`, decision, selector)
}

func reachedThePage(decision string) string {
	return fmt.Sprintf(`!!document.querySelector('.decision[data-decision="%s"]')`, decision)
}

func nothingLeftTheMachine(t *testing.T, address string, reaching []string) {
	t.Helper()

	origin := originOf(t, address)
	for _, asked := range reaching {
		if !strings.HasPrefix(asked, origin+"/") && !strings.HasPrefix(asked, "data:") {
			t.Errorf("the page requested %q, and every request it issues resolves to the surface itself", asked)
		}
	}
	if len(reaching) < 3 {
		t.Errorf("the page issued %d requests, and a check that saw almost nothing would pass forever: %v", len(reaching), reaching)
	}
}

func originOf(t *testing.T, address string) string {
	t.Helper()

	read, err := url.Parse(address)
	if err != nil {
		t.Fatalf("reading the address the page was served at: %v", err)
	}
	return read.Scheme + "://" + read.Host
}
