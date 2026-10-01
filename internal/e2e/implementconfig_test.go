package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

func implementConfigurationPath(session *harness.Session) string {
	return filepath.Join(session.Root, "implement", "config.yaml")
}

func stateImplementConfiguration(t *testing.T, session *harness.Session, content string) {
	t.Helper()

	path := implementConfigurationPath(session)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func implementSettingsIn(t *testing.T, printed string) map[string][2]string {
	t.Helper()

	reported := map[string][2]string{}
	for line := range strings.SplitSeq(printed, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		reported[fields[0]] = [2]string{fields[1], fields[2]}
	}
	return reported
}

func TestImplementConfigWithNoFileReportsTheThreeSettingsAtTheirDefaults(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config with no file exited %d: %s", code, failures)
	}

	want := map[string][2]string{
		"concurrency":              {"4", string(config.Default)},
		"roles.implementer.model":  {"sonnet", string(config.Default)},
		"roles.implementer.effort": {"high", string(config.Default)},
	}
	if diff := cmp.Diff(want, implementSettingsIn(t, printed)); diff != "" {
		t.Errorf("what socrates implement config reports differs (-want +got):\n%s", diff)
	}
	if first := firstLine(printed); first != implementConfigurationPath(session) {
		t.Errorf("socrates implement config names %q, and the file under the stated home is %s", first, implementConfigurationPath(session))
	}
	if _, err := os.Stat(implementConfigurationPath(session)); !os.IsNotExist(err) {
		t.Errorf("socrates implement config created %s", implementConfigurationPath(session))
	}
}

func TestImplementConfigReportsWhatAPersonWroteAsSet(t *testing.T) {
	session := harness.Start(t)
	stateImplementConfiguration(t, session, "concurrency: 2\nroles:\n  implementer:\n    effort: max\n")

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}

	want := map[string][2]string{
		"concurrency":              {"2", string(config.Set)},
		"roles.implementer.model":  {"sonnet", string(config.Default)},
		"roles.implementer.effort": {"max", string(config.Set)},
	}
	if diff := cmp.Diff(want, implementSettingsIn(t, printed)); diff != "" {
		t.Errorf("what socrates implement config reports differs (-want +got):\n%s", diff)
	}
}

func TestImplementConfigRefusesAValueOutOfRangeAndAFieldItDoesNotHaveNamingTheFile(t *testing.T) {
	for name, content := range map[string]string{
		"a concurrency of zero": "concurrency: 0\n",
		"an unknown effort":     "roles:\n  implementer:\n    effort: extreme\n",
		"an unknown field":      "parallelism: 3\n",
	} {
		t.Run(name, func(t *testing.T) {
			session := harness.Start(t)
			stateImplementConfiguration(t, session, content)

			printed, failures, code := session.Top("implement", "config")
			if code != exit.Unusable {
				t.Fatalf("socrates implement config with %s exited %d, want %d", name, code, exit.Unusable)
			}
			if !strings.Contains(failures, implementConfigurationPath(session)) {
				t.Errorf("the refusal does not name the file: %s", failures)
			}
			if printed != "" {
				t.Errorf("socrates implement config reported settings while the file is refused:\n%s", printed)
			}
		})
	}
}

func TestImplementReadsNothingFromReviewsFile(t *testing.T) {
	session := harness.Start(t)
	stateConfiguration(t, session, "diff_ceiling: 4096\n")

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}
	for _, reported := range implementSettingsIn(t, printed) {
		if reported[1] != string(config.Default) {
			t.Errorf("a setting reads as %s while only review's file states anything:\n%s", reported[1], printed)
		}
	}
}

func TestTheTopLevelHelpListsImplementAndImplementHelpListsWhatItHolds(t *testing.T) {
	session := harness.Start(t)

	top, failures, code := session.Top("--help")
	if code != exit.NothingAsked || !strings.Contains(top, "implement") {
		t.Errorf("socrates --help exited %d and reads %q%s, without listing implement", code, top, failures)
	}

	own, failures, code := session.Top("implement", "--help")
	if code != exit.NothingAsked || !strings.Contains(own, "config") {
		t.Errorf("socrates implement --help exited %d and reads %q%s, without listing config", code, own, failures)
	}
}
