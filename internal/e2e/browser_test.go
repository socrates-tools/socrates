package e2e_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

const (
	browserVariable = "AGORA_BROWSER"
	browserDeadline = 90 * time.Second
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

	profile, err := os.MkdirTemp("/tmp", "ag")
	if err != nil {
		t.Fatalf("making somewhere short for the browser to keep its profile: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(profile) })

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

	requests := &recorded{}
	chromedp.ListenTarget(bounded, func(event any) {
		if sent, is := event.(*network.EventRequestWillBeSent); is {
			requests.note(sent.Request.URL)
		}
	})

	held := &reader{t: t, context: bounded, requests: requests}
	held.do(
		network.Enable(),
		chromedp.Navigate("http://"+address+"/"),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.WaitReady(".send", chromedp.ByQuery),
	)
	held.settled()
	return held
}

func (r *reader) do(actions ...chromedp.Action) {
	r.t.Helper()

	if err := chromedp.Run(r.context, actions...); err != nil {
		r.t.Fatalf("driving the page: %v", err)
	}
}

func (r *reader) settled() {
	r.t.Helper()

	var bound bool
	r.do(chromedp.Poll("!!document.querySelector('.commit') && !!window.getComputedStyle(document.body).backgroundColor", &bound))
}

func (r *reader) click(selector string) {
	r.t.Helper()
	r.do(chromedp.Click(selector, chromedp.ByQuery, chromedp.NodeVisible))
}

func (r *reader) write(selector, words string) {
	r.t.Helper()
	r.do(chromedp.SendKeys(selector, words, chromedp.ByQuery))
}

func (r *reader) evaluate(expression string, into any) {
	r.t.Helper()
	r.do(chromedp.Evaluate(expression, into))
}

func (r *reader) waitFor(expression string) {
	r.t.Helper()

	var held bool
	if err := chromedp.Run(r.context, chromedp.Poll(expression, &held)); err != nil {
		r.t.Fatalf("waiting for %s: %v", expression, err)
	}
}

func (r *reader) emulate(preference string) {
	r.t.Helper()
	r.do(emulation.SetEmulatedMedia().WithFeatures([]*emulation.MediaFeature{
		{Name: "prefers-color-scheme", Value: preference},
	}))
}

func (r *reader) narrow(width, height int64) {
	r.t.Helper()
	r.do(chromedp.EmulateViewport(width, height))
}

func chose(decision string, option int) string {
	return fmt.Sprintf(`document.querySelector('.decision[data-decision="%s"] .choice[value="%d"]').click()`, decision, option)
}

func stateOf(decision string) string {
	return fmt.Sprintf(`document.querySelector('.decision[data-decision="%s"]').dataset.state`, decision)
}

func within(decision, selector string) string {
	return fmt.Sprintf(`.decision[data-decision="%s"] %s`, decision, selector)
}

func nothingLeftTheMachine(t *testing.T, address string, reaching []string) {
	t.Helper()

	for _, asked := range reaching {
		if !strings.HasPrefix(asked, "http://"+address+"/") && !strings.HasPrefix(asked, "data:") {
			t.Errorf("the page requested %q, and every request it issues resolves to the surface itself", asked)
		}
	}
	if len(reaching) < 3 {
		t.Errorf("the page issued %d requests, and a check that saw almost nothing would pass forever: %v", len(reaching), reaching)
	}
}
