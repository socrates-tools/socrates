package convention

import (
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/google/go-cmp/cmp"
)

func TestAnEmbeddedAssetWithNoEntryFailsTheCheck(t *testing.T) {
	embedded := fstest.MapFS{
		"bundle/surface.js":                   &fstest.MapFile{Data: []byte("(()=>{})();")},
		"bundle/fonts/IBMPlexSerif-400.woff2": &fstest.MapFile{Data: []byte("a face somebody bundled")},
	}
	notice := "bundle/surface.js is the client, and it is socrates' own work.\n"

	unaccounted, err := EmbeddedFilesWithoutANotice(embedded, notice)
	if err != nil {
		t.Fatalf("reading what is embedded: %v", err)
	}

	want := []string{"bundle/fonts/IBMPlexSerif-400.woff2"}
	if diff := cmp.Diff(want, unaccounted); diff != "" {
		t.Errorf("the unaccounted files differ (-want +got):\n%s", diff)
	}
}

func TestABundlerThatDeletesNoticesFailsTheSuite(t *testing.T) {
	notice := "==========================\n" +
		"highlight.js 11.12.0\n" +
		"BSD 3-Clause License.\n" +
		"==========================\n\n" +
		"Copyright (c) 2006, Ivan Sagalaev.\nAll rights reserved.\n"
	bundled := []string{"highlight.js"}

	carried := "(()=>{})();\n/*! Copyright (c) 2006, Ivan Sagalaev. */\n"
	if missing := NoticesTheClientDoesNotCarry(carried, notice, bundled); len(missing) != 0 {
		t.Errorf("the check reported %v against a bundle that carries the notice", missing)
	}

	deleted := "(()=>{})();\n"
	want := []MissingNotice{{Library: "highlight.js", Line: "Copyright (c) 2006, Ivan Sagalaev."}}
	if diff := cmp.Diff(want, NoticesTheClientDoesNotCarry(deleted, notice, bundled)); diff != "" {
		t.Errorf("what the bundle does not carry differs (-want +got):\n%s", diff)
	}
}

func TestAnAssetInTheSourceButNotInTheBinaryNeedsNoEntry(t *testing.T) {
	tree := t.TempDir()
	write(t, tree, filepath.Join("web", "fonts", "IBMPlexSerif-400.woff2"), "a face nobody bundled")

	embedded := fstest.MapFS{
		"bundle/surface.js": &fstest.MapFile{Data: []byte("(()=>{})();")},
	}
	notice := "bundle/surface.js is the client, and it is socrates' own work.\n"

	unaccounted, err := EmbeddedFilesWithoutANotice(embedded, notice)
	if err != nil {
		t.Fatalf("reading what is embedded: %v", err)
	}

	if len(unaccounted) != 0 {
		t.Errorf("the check asked for an entry for %v, and what a binary distributes is what a notice is a condition of", unaccounted)
	}
}
