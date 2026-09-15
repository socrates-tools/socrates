package assets_test

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/agora/assets"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

var stamped = []string{"src", "fonts", "build.mjs", "package.json", "pnpm-lock.yaml"}

func TestTheCommittedBundleIsWhatItsSourceBuilds(t *testing.T) {
	root := source.Root(t)
	client := filepath.Join(root, "agora", "web")

	digest := sha256.New()
	for _, named := range stamped {
		for _, held := range within(t, filepath.Join(client, named)) {
			relative, err := filepath.Rel(client, held)
			if err != nil {
				t.Fatalf("naming %s: %v", held, err)
			}
			digest.Write([]byte(filepath.ToSlash(relative)))
			write(t, digest, held)
		}
	}

	built := hex.EncodeToString(digest.Sum(nil))
	if built != assets.Sum() {
		t.Errorf("the committed bundle was built from a different source than the one in web:\ncommitted %s\nsource    %s\nrun make assets and commit what it writes", assets.Sum(), built)
	}
}

func TestTheBundleHoldsWhatThePageAsksFor(t *testing.T) {
	// scenario: Rendering needs no network
	for _, asked := range []string{"surface.css", "surface.js"} {
		held, err := fs.ReadFile(assets.Files(), asked)
		if err != nil {
			t.Errorf("the bundle holds no %s, and the page asks for it: %v", asked, err)
			continue
		}
		if len(held) == 0 {
			t.Errorf("the bundle holds an empty %s", asked)
		}
	}

	faces, err := fs.Glob(assets.Files(), "fonts/*.woff2")
	if err != nil {
		t.Fatalf("reading what the bundle holds: %v", err)
	}
	if len(faces) == 0 {
		t.Error("the bundle embeds no font, so the page would ask a host for one")
	}
}

var absoluteStyleReference = regexp.MustCompile(`url\(\s*["']?(https?:|//)`)

var absoluteRequest = regexp.MustCompile("(fetch|EventSource|WebSocket|importScripts|open)\\(\\s*[\"'`]\\s*(https?:|//)")

func TestNothingInTheBundleReachesAHost(t *testing.T) {
	// scenario: The page reaches nothing
	style, err := fs.ReadFile(assets.Files(), "surface.css")
	if err != nil {
		t.Fatalf("reading the stylesheet: %v", err)
	}
	for _, reaching := range absoluteStyleReference.FindAllString(string(style), -1) {
		t.Errorf("the stylesheet fetches %q, and the page reaches nothing outside the surface", reaching)
	}

	client, err := fs.ReadFile(assets.Files(), "surface.js")
	if err != nil {
		t.Fatalf("reading the client: %v", err)
	}
	for _, reaching := range absoluteRequest.FindAllString(string(client), -1) {
		t.Errorf("the client requests %q, and the page reaches nothing outside the surface", reaching)
	}
}

func TestNoDiagramEngineIsEmbedded(t *testing.T) {
	// scenario: A diff renders from a patch
	held, err := fs.ReadFile(assets.Files(), "surface.js")
	if err != nil {
		t.Fatalf("reading the client: %v", err)
	}

	for _, engine := range []string{"mermaid", "graphviz", "viz.js", "cytoscape", "dagre", "elkjs", "d3-force"} {
		if strings.Contains(strings.ToLower(string(held)), engine) {
			t.Errorf("the client embeds %q, and an illustration is drawn by the caller rather than laid out by a machine", engine)
		}
	}
}

func within(t *testing.T, root string) []string {
	t.Helper()

	held, err := os.Stat(root)
	if err != nil {
		t.Fatalf("reading %s: %v", root, err)
	}
	if !held.IsDir() {
		return []string{root}
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatalf("reading %s: %v", root, err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	slices.Sort(names)

	var found []string
	for _, name := range names {
		found = append(found, within(t, filepath.Join(root, name))...)
	}
	return found
}

func write(t *testing.T, digest hash.Hash, path string) {
	t.Helper()

	held, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	digest.Write(held)
}
