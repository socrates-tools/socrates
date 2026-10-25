package runstore

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writtenFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}

func TestAnAttachmentIsCopiedIntoTheRunStore(t *testing.T) {
	round := roundFor(t)
	picture := writtenFile(t, "before.png", "\x89PNG\r\n\x1a\nbytes of a picture")
	trace := writtenFile(t, "trace.txt", "a trace the person captured")

	attached, err := round.Attach([]string{picture, trace})
	if err != nil {
		t.Fatalf("attaching: %v", err)
	}
	if len(attached) != 2 {
		t.Fatalf("%d attachments came back, want 2: %+v", len(attached), attached)
	}

	t.Run("The copy lives inside the round and survives the original", func(t *testing.T) {
		if err := os.Remove(picture); err != nil {
			t.Fatalf("deleting the original: %v", err)
		}
		held, err := os.ReadFile(attached[0].Path)
		if err != nil {
			t.Fatalf("reading the copy: %v", err)
		}
		if string(held) != "\x89PNG\r\n\x1a\nbytes of a picture" {
			t.Errorf("the copy holds %q", held)
		}
		if !filepath.IsAbs(attached[0].Path) || !strings.HasPrefix(attached[0].Path, round.Directory) {
			t.Errorf("%s is not an absolute path inside %s", attached[0].Path, round.Directory)
		}
	})

	t.Run("Each carries its name and type", func(t *testing.T) {
		if attached[0].Name != "before.png" || attached[0].Type != "image/png" {
			t.Errorf("the picture is recorded as %+v", attached[0])
		}
		if attached[1].Name != "trace.txt" || attached[1].Type != "text/plain" {
			t.Errorf("the trace is recorded as %+v", attached[1])
		}
	})

	t.Run("A second attachment of the same name does not overwrite the first", func(t *testing.T) {
		again := writtenFile(t, "trace.txt", "a different trace")
		second, err := round.Attach([]string{again})
		if err != nil {
			t.Fatalf("attaching again: %v", err)
		}
		if second[0].Path == attached[1].Path {
			t.Errorf("both copies are at %s", second[0].Path)
		}
		held, _ := os.ReadFile(attached[1].Path)
		if string(held) != "a trace the person captured" {
			t.Errorf("the first copy now holds %q", held)
		}
	})
}

func TestAnAttachmentThatIsNotAFileIsRefusedAndNothingIsKept(t *testing.T) {
	round := roundFor(t)
	present := writtenFile(t, "present.txt", "here")
	missing := filepath.Join(t.TempDir(), "missing.png")

	for name, sources := range map[string][]string{
		"a path that does not exist": {present, missing},
		"a directory":                {present, t.TempDir()},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := round.Attach(sources)
			if err == nil {
				t.Fatal("it was attached")
			}
			if !strings.Contains(err.Error(), sources[1]) {
				t.Errorf("the refusal does not name %s: %v", sources[1], err)
			}
			if _, err := os.Stat(filepath.Join(round.Directory, "attachments")); err == nil {
				entries, _ := os.ReadDir(filepath.Join(round.Directory, "attachments"))
				if len(entries) != 0 {
					t.Errorf("a refused attach left %d copies behind", len(entries))
				}
			}
		})
	}
}
