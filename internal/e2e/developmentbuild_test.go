package e2e_test

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/testing/source"
)

const editedStylesheet = "body { background: #123456 }\n"

func theStylesheetServedAt(t *testing.T, address string) string {
	t.Helper()

	parsed, err := url.Parse(address)
	if err != nil {
		t.Fatalf("reading the address %q: %v", address, err)
	}
	answered, err := http.Get(parsed.Scheme + "://" + parsed.Host + "/assets/surface.css")
	if err != nil {
		t.Fatalf("reaching the stylesheet: %v", err)
	}
	defer func() { _ = answered.Body.Close() }()
	if answered.StatusCode != http.StatusOK {
		t.Fatalf("the stylesheet was answered with %d", answered.StatusCode)
	}
	served, err := io.ReadAll(answered.Body)
	if err != nil {
		t.Fatalf("reading the stylesheet: %v", err)
	}
	return string(served)
}

func aBuildFromACopyOfTheTree(t *testing.T, tags ...string) (binary, bundle string) {
	t.Helper()

	root := source.CopyOfTree(t)
	binary = filepath.Join(t.TempDir(), "socrates")
	source.BuildIn(t, root, "./cmd/socrates", binary, tags...)
	return binary, filepath.Join(root, "agora", "assets", "bundle")
}

func TestADevelopmentBuildServesAnEditedStylesheetAfterAReloadWithNoRebuild(t *testing.T) {
	t.Parallel()

	binary, bundle := aBuildFromACopyOfTheTree(t, "agoradev")
	held := standing(t)
	held.binary = binary
	address := held.serve(oneDecisionDocument)

	embedded, err := os.ReadFile(filepath.Join(bundle, "surface.css"))
	if err != nil {
		t.Fatal(err)
	}
	if served := theStylesheetServedAt(t, address); served != string(embedded) {
		t.Fatalf("before any edit the development build served a stylesheet that is not the bundle on disk")
	}

	if err := os.WriteFile(filepath.Join(bundle, "surface.css"), []byte(editedStylesheet), 0o644); err != nil {
		t.Fatal(err)
	}
	if served := theStylesheetServedAt(t, address); served != editedStylesheet {
		t.Fatalf("after the bundle on disk changed the development build served %q rather than the edit", truncated(served))
	}
}

func TestABuildWithoutTheTagServesTheEmbeddedBundleEvenWhenTheBundleOnDiskChanged(t *testing.T) {
	t.Parallel()

	binary, bundle := aBuildFromACopyOfTheTree(t)
	embedded, err := os.ReadFile(filepath.Join(bundle, "surface.css"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bundle, "surface.css"), []byte(editedStylesheet), 0o644); err != nil {
		t.Fatal(err)
	}

	held := standing(t)
	held.binary = binary
	address := held.serve(oneDecisionDocument)

	if served := theStylesheetServedAt(t, address); served != string(embedded) {
		t.Fatalf("a build without the tag served %q rather than the embedded bundle", truncated(served))
	}
}

func truncated(text string) string {
	if len(text) > 80 {
		return strings.TrimSpace(text[:80])
	}
	return text
}
