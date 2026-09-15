package configure

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/socrates-tools/socrates/internal/config"
)

func TestAFileTheMenuNeverTouchedComesBackByteForByte(t *testing.T) {
	held := `# my notes about this file

reviewer_cli: claude

# the ceiling I chose
diff_ceiling: 4096

roles:
  reviewer:
    model: opus

surface:
  bind: 127.0.0.1:18018
`
	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a file a person wrote: %v", err)
	}

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing it back: %v", err)
	}
	if string(written) != held {
		t.Errorf("a file nothing changed came back as:\n%s\nand it was written as:\n%s", written, held)
	}
}

func TestSettingOneKeyLeavesEveryCommentAndBlankLineWhereItWas(t *testing.T) {
	held := `# my notes about this file

reviewer_cli: claude

# the ceiling I chose
diff_ceiling: 4096
`
	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a file a person wrote: %v", err)
	}
	document.set("reviewer_cli", "codex")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing it back: %v", err)
	}

	want := strings.Replace(held, "reviewer_cli: claude", "reviewer_cli: codex", 1)
	if string(written) != want {
		t.Errorf("setting one key produced:\n%s\nwant:\n%s", written, want)
	}
}

func TestSettingAKeyNoFileNamesWritesOnlyThatKey(t *testing.T) {
	document, err := documentOf(nil)
	if err != nil {
		t.Fatalf("reading no file at all: %v", err)
	}
	document.set("reviewer_cli", "codex")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing the first setting a person chose: %v", err)
	}
	if want := "reviewer_cli: codex\n"; string(written) != want {
		t.Errorf("the file holds %q, want %q", written, want)
	}
}

func TestSettingANestedKeyWritesTheBlocksThatHoldIt(t *testing.T) {
	document, err := documentOf(nil)
	if err != nil {
		t.Fatalf("reading no file at all: %v", err)
	}
	document.set("roles.reviewer.model", "sonnet")
	document.set("surface.bind", "127.0.0.1:18019")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing two nested settings: %v", err)
	}
	want := "roles:\n  reviewer:\n    model: sonnet\nsurface:\n  bind: 127.0.0.1:18019\n"
	if string(written) != want {
		t.Errorf("the file holds:\n%s\nwant:\n%s", written, want)
	}
}

func TestRemovingTheOnlySettingLeavesAFileNamingNothing(t *testing.T) {
	document, err := documentOf([]byte("reviewer_cli: codex\n"))
	if err != nil {
		t.Fatalf("reading a file naming one setting: %v", err)
	}
	document.remove("reviewer_cli")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing back a file naming nothing: %v", err)
	}
	if strings.TrimSpace(string(written)) != "" {
		t.Errorf("removing the only setting left %q, and a file naming nothing is empty", written)
	}
}

func TestRemovingTheLastSettingOfABlockTakesTheBlockWithIt(t *testing.T) {
	held := "diff_ceiling: 4096\nroles:\n  reviewer:\n    model: sonnet\n"

	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a file naming a nested setting: %v", err)
	}
	document.remove("roles.reviewer.model")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing back: %v", err)
	}
	if want := "diff_ceiling: 4096\n"; string(written) != want {
		t.Errorf("removing the last setting of a block left:\n%s\nwant:\n%s", written, want)
	}
}

func TestRemovingOneSettingOfABlockLeavesTheRestOfTheBlock(t *testing.T) {
	held := "roles:\n  reviewer:\n    model: sonnet\n    effort: low\n"

	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a file naming two nested settings: %v", err)
	}
	document.remove("roles.reviewer.model")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing back: %v", err)
	}
	if want := "roles:\n  reviewer:\n    effort: low\n"; string(written) != want {
		t.Errorf("removing one nested setting left:\n%s\nwant:\n%s", written, want)
	}
}

func TestRemovingASettingLeavesTheCommentsAroundItAlone(t *testing.T) {
	held := "# what I set and why\ndiff_ceiling: 4096\n\n# the retention I chose\nrun_retention: 72h\n"

	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a commented file: %v", err)
	}
	document.remove("diff_ceiling")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing back: %v", err)
	}
	if !strings.Contains(string(written), "# the retention I chose") {
		t.Errorf("removing one setting took another setting's comment with it:\n%s", written)
	}
	if strings.Contains(string(written), "diff_ceiling") {
		t.Errorf("the setting removed is still named:\n%s", written)
	}
}

func TestRemovingASettingNoFileNamesChangesNothing(t *testing.T) {
	held := "diff_ceiling: 4096\n"

	document, err := documentOf([]byte(held))
	if err != nil {
		t.Fatalf("reading a file naming one setting: %v", err)
	}
	if document.states("run_retention") {
		t.Fatal("a file naming one setting reports a second one as stated")
	}
	document.remove("run_retention")

	written, err := document.render()
	if err != nil {
		t.Fatalf("writing back: %v", err)
	}
	if string(written) != held {
		t.Errorf("removing a setting the file never named produced:\n%s\nwant:\n%s", written, held)
	}
}

func TestAWriteReplacesTheFileRatherThanTruncatingItInPlace(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "config.yaml")
	if err := os.WriteFile(path, []byte("diff_ceiling: 4096\n"), 0o644); err != nil {
		t.Fatalf("writing the file a person already had: %v", err)
	}
	before := inodeOf(t, path)

	if err := writeAtomically(path, []byte("diff_ceiling: 8192\n")); err != nil {
		t.Fatalf("writing the setting a person confirmed: %v", err)
	}

	if after := inodeOf(t, path); after == before {
		t.Errorf("the file at %s is the one that was there, so the write went through the old file rather than through a rename", path)
	}
	if held, err := os.ReadFile(path); err != nil || string(held) != "diff_ceiling: 8192\n" {
		t.Errorf("the file holds %q after the write: %v", held, err)
	}

	left, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		t.Fatalf("reading the directory the file lives in: %v", err)
	}
	if len(left) != 1 {
		t.Errorf("the directory holds %v, and a finished write leaves the file alone in it", left)
	}
}

func TestAWriteCreatesTheDirectoryTheFileLivesIn(t *testing.T) {
	path := filepath.Join(t.TempDir(), "socrates", "review", "config.yaml")

	if err := writeAtomically(path, []byte("diff_ceiling: 8192\n")); err != nil {
		t.Fatalf("writing the first setting on a machine with no home yet: %v", err)
	}
	if held, err := os.ReadFile(path); err != nil || string(held) != "diff_ceiling: 8192\n" {
		t.Errorf("the file holds %q: %v", held, err)
	}
}

func inodeOf(t *testing.T, path string) uint64 {
	t.Helper()

	held, err := os.Stat(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	described, isUnix := held.Sys().(*syscall.Stat_t)
	if !isUnix {
		t.Skipf("%s reports no inode on this platform", path)
	}
	return described.Ino
}

func TestTheSettingNamedIsTheNestedOneThatIsWrong(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	held := "diff_ceiling: 4096\nroles:\n  reviewer:\n    model: sonnet\n    temperature: 0.5\n"

	found := unreadableIn(path, []byte(held))
	if len(found) != 1 {
		t.Fatalf("a file with one unreadable nested setting reports %d of them: %+v", len(found), found)
	}
	if found[0].key != "roles.reviewer.temperature" {
		t.Errorf("the setting named is %q, and the one that is wrong is roles.reviewer.temperature", found[0].key)
	}
	if found[0].line != 5 {
		t.Errorf("the line named is %d, and the one that is wrong is 5", found[0].line)
	}

	written, err := repaired([]byte(held), found[0])
	if err != nil {
		t.Fatalf("resetting it: %v", err)
	}
	if err := config.Parse(path, written); err != nil {
		t.Errorf("the repaired file is one the loader refuses: %v\n%s", err, written)
	}
	for _, kept := range []string{"diff_ceiling: 4096", "model: sonnet"} {
		if !strings.Contains(string(written), kept) {
			t.Errorf("%q was dropped along with the setting that was reset:\n%s", kept, written)
		}
	}
	if strings.Contains(string(written), "temperature") {
		t.Errorf("the setting that was reset is still named:\n%s", written)
	}
}

func TestResettingTheLastSettingOfANestedBlockTakesTheBlockWithIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	held := "diff_ceiling: 4096\nroles:\n  reviewer:\n    temperature: 0.5\n"

	found := unreadableIn(path, []byte(held))
	if len(found) != 1 || found[0].key != "roles.reviewer.temperature" {
		t.Fatalf("the unreadable settings read %+v", found)
	}

	written, err := repaired([]byte(held), found[0])
	if err != nil {
		t.Fatalf("resetting it: %v", err)
	}
	if want := "diff_ceiling: 4096\n"; string(written) != want {
		t.Errorf("the repaired file holds:\n%s\nwant:\n%s", written, want)
	}
}

func TestAFileThatIsNotYAMLAtAllNamesTheLineThatOpenedIt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	held := "roles: [not a mapping\ndiff_ceiling: 4096\n"

	found := unreadableIn(path, []byte(held))
	if len(found) == 0 {
		t.Fatalf("a file that is not YAML at all reports no unreadable setting")
	}
	if found[0].line != 1 {
		t.Errorf("the line named is %d, and the line that opened the flow is 1", found[0].line)
	}

	written, err := repaired([]byte(held), found[0])
	if err != nil {
		t.Fatalf("resetting it: %v", err)
	}
	if err := config.Parse(path, written); err != nil {
		t.Errorf("the repaired file is one the loader refuses: %v\n%s", err, written)
	}
	if !strings.Contains(string(written), "diff_ceiling: 4096") {
		t.Errorf("the setting that was readable was dropped:\n%s", written)
	}
}

func TestAFileNoSingleSettingAccountsForIsNotGuessedAt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")

	if found := unreadableIn(path, []byte("diff_ceiling: 4096\ndiff_ceiling: 8192\n")); len(found) != 0 {
		t.Errorf("a file whose settings each read on their own reports %+v, and socrates-review guesses at nothing", found)
	}
}
