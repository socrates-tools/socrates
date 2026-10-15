package e2e_test

import (
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	implementskill "github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/convention"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func installedKey(host skillkit.Host, name string) string {
	return filepath.ToSlash(filepath.Join(host.SkillsDirectory, name, skillkit.FileName))
}

func installedAgentKey(host skillkit.Host) string {
	return filepath.ToSlash(filepath.Join(host.AgentsDirectory, implementskill.AgentFileName()))
}

func installedPolicyKey(host skillkit.Host, name string) string {
	return filepath.ToSlash(filepath.Join(host.SkillsDirectory, name, host.PolicyFile))
}

func hostRoots() []string {
	var roots []string
	for _, host := range skillkit.Hosts() {
		for _, directory := range []string{host.SkillsDirectory, host.AgentsDirectory} {
			if root := filepath.Dir(directory); !slices.Contains(roots, root) {
				roots = append(roots, root)
			}
		}
	}
	return roots
}

func pathWithout(t *testing.T, session *harness.Session, binary string) string {
	t.Helper()

	var path string
	for _, entry := range session.Environment() {
		if value, found := strings.CutPrefix(entry, "PATH="); found {
			path = value
		}
	}

	var kept []string
	for _, directory := range filepath.SplitList(path) {
		if _, err := os.Stat(filepath.Join(directory, binary)); err != nil {
			kept = append(kept, directory)
		}
	}
	return "PATH=" + strings.Join(kept, string(os.PathListSeparator))
}

func writtenIn(t *testing.T, root string) map[string]string {
	t.Helper()

	written := map[string]string{}
	walked := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		written[filepath.ToSlash(relative)] = string(content)
		return nil
	})
	if walked != nil && !os.IsNotExist(walked) {
		t.Fatalf("reading what was written under %s: %v", root, walked)
	}
	return written
}

func writtenUnder(t *testing.T, home string) map[string]string {
	t.Helper()

	written := map[string]string{}
	for _, relativeRoot := range hostRoots() {
		for path, content := range writtenIn(t, filepath.Join(home, relativeRoot)) {
			written[filepath.ToSlash(filepath.Join(relativeRoot, path))] = content
		}
	}
	return written
}

func installedOneNamespaceAtATime(t *testing.T) (map[string]string, string) {
	t.Helper()
	session := harness.Start(t)

	var printed strings.Builder
	for _, tool := range convention.Tools {
		help, _, code := session.Top(tool.Namespace, "skill", "--help")
		if code != 0 || !strings.Contains(help, "socrates "+tool.Namespace+" skill [flags]") {
			continue
		}
		out, failures, code := session.Top(tool.Namespace, "skill")
		if code != 0 {
			t.Fatalf("socrates %s skill exited %d: %s\n%s", tool.Namespace, code, failures, out)
		}
		printed.WriteString(out)
	}
	return writtenUnder(t, session.UserHome), strings.ReplaceAll(printed.String(), session.UserHome, "<home>")
}

func linesOf(printed string) []string {
	return slices.Sorted(slices.Values(strings.Split(strings.TrimSpace(printed), "\n")))
}

func TestOneCommandInstallsEverySkillEveryNamespaceOffers(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("skill")
	if code != 0 {
		t.Fatalf("socrates skill exited %d: %s\n%s", code, failures, printed)
	}

	wantWritten, wantPrinted := installedOneNamespaceAtATime(t)
	if len(wantWritten) == 0 {
		t.Fatal("no namespace installed anything, so this check would pass forever")
	}

	written := writtenUnder(t, session.UserHome)
	if diff := cmp.Diff(slices.Sorted(maps.Keys(wantWritten)), slices.Sorted(maps.Keys(written))); diff != "" {
		t.Errorf("socrates skill did not write what every namespace's own skill writes (-want +got):\n%s", diff)
	}
	for path, content := range wantWritten {
		if written[path] != content {
			t.Errorf("socrates skill wrote %s other than the namespace's own skill does", path)
		}
	}
	if diff := cmp.Diff(linesOf(wantPrinted), linesOf(strings.ReplaceAll(printed, session.UserHome, "<home>"))); diff != "" {
		t.Errorf("socrates skill does not say what every namespace's own skill says (-want +got):\n%s", diff)
	}

	for _, one := range allSkills() {
		if _, held := written[installedKey(one.Host, one.Name)]; !held {
			t.Errorf("socrates skill left out %s", one.Name)
		}
	}
	if _, held := written[installedAgentKey(skillkit.ClaudeCode)]; !held {
		t.Errorf("socrates skill left out the implementer's agent definition")
	}
}

func TestOneCommandWritesEveryHostsRenderingWhateverIsOnPath(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	withoutCodex := pathWithout(t, session, "codex")

	printed, failures, code := session.TopWith([]string{withoutCodex}, "skill")
	if code != 0 {
		t.Fatalf("socrates skill exited %d: %s\n%s", code, failures, printed)
	}

	written := writtenUnder(t, session.UserHome)
	for _, one := range allSkills() {
		if written[installedKey(one.Host, one.Name)] != one.Rendered() {
			t.Errorf("socrates skill with no codex on PATH did not write %s for %s as the package renders it", one.Name, one.Host.Name)
		}
		if one.Host.PolicyFile != "" && written[installedPolicyKey(one.Host, one.Name)] != skillkit.PolicyDocument {
			t.Errorf("socrates skill did not write %s beside %s for %s", one.Host.PolicyFile, one.Name, one.Host.Name)
		}
	}
}

func TestOneCommandRefusesAnArgument(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	printed, failures, code := session.Top("skill", "review")
	if code == 0 {
		t.Fatalf("socrates skill review exited 0:\n%s", printed)
	}
	if !strings.Contains(failures, `"review"`) {
		t.Errorf("socrates skill review does not name the argument it refused:\n%s", failures)
	}
	if written := writtenUnder(t, session.UserHome); len(written) != 0 {
		t.Errorf("socrates skill review wrote %v though it was refused", slices.Sorted(maps.Keys(written)))
	}
}
