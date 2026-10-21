package assets_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/socrates-tools/socrates/agora/assets"
	"github.com/socrates-tools/socrates/internal/convention"
	"github.com/socrates-tools/socrates/internal/testing/source"
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
	bundled, err := convention.LibrariesTheClientBundles(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}
	if len(bundled) == 0 {
		t.Fatalf("the client bundles nothing, so this test proves nothing")
	}

	for _, missing := range convention.NoticesTheBundleDoesNotCarry(assets.Files(), theNotice(t), bundled) {
		switch {
		case missing.Line == "":
			t.Errorf("%s states no copyright line for %s, so there is nothing for the bundle to carry", convention.NoticeFile, missing.Library)
		case missing.File == "":
			t.Errorf("%s names no file of the bundle that %s is bundled into, so there is nowhere to look for %q", convention.NoticeFile, missing.Library, missing.Line)
		default:
			t.Errorf("the %s that ships carries no %q, and %s is bundled into it", missing.File, missing.Line, missing.Library)
		}
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

func TestTheLandingPagesStylesheetCarriesTheNoticeOfWhatItBundles(t *testing.T) {
	bundled, err := convention.LibrariesTheClientBundles(source.Root(t))
	if err != nil {
		t.Fatalf("reading what the client bundles: %v", err)
	}

	site := os.DirFS(filepath.Join(source.Root(t), "site"))
	for _, missing := range convention.NoticesTheLandingPageDoesNotCarry(site, theNotice(t), bundled) {
		if missing.Library == "" {
			t.Errorf("%s names no library bundled into site/%s, and it ships beside the landing page", convention.NoticeFile, missing.File)
			continue
		}
		t.Errorf("the site/%s that ships carries no %q, and %s is bundled into it", missing.File, missing.Line, missing.Library)
	}
}
