package assets_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-devtool/socrates/internal/assets"
	"github.com/socrates-devtool/socrates/internal/convention"
	"github.com/socrates-devtool/socrates/internal/testing/source"
)

func TestEveryEmbeddedFileIsAccountedForInTheNotice(t *testing.T) {
	unaccounted, err := convention.EmbeddedFilesWithoutANotice(assets.Embedded(), theNotice(t))
	if err != nil {
		t.Fatalf("reading what is embedded: %v", err)
	}

	for _, path := range unaccounted {
		t.Errorf("%s accounts for no %s, and that file ships inside the binary", convention.NoticeFile, path)
	}
}

func TestTheShippedBundleCarriesTheNoticeOfWhatItBundles(t *testing.T) {
	client, err := fs.ReadFile(assets.Files(), "surface.js")
	if err != nil {
		t.Fatalf("reading the embedded client: %v", err)
	}

	bundled, err := convention.LibrariesTheClientBundles(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}
	if len(bundled) == 0 {
		t.Fatalf("the client bundles nothing, so this test proves nothing")
	}

	for _, missing := range convention.NoticesTheClientDoesNotCarry(string(client), theNotice(t), bundled) {
		if missing.Line == "" {
			t.Errorf("%s states no copyright line for %s, so there is nothing for the bundle to carry", convention.NoticeFile, missing.Library)
			continue
		}
		t.Errorf("the bundle that ships carries no %q, and %s is bundled into it", missing.Line, missing.Library)
	}
}

func theNotice(t *testing.T) string {
	t.Helper()

	held, err := os.ReadFile(filepath.Join(source.Root(t), convention.NoticeFile))
	if err != nil {
		t.Fatalf("the root carries no %s: %v", convention.NoticeFile, err)
	}
	return string(held)
}
