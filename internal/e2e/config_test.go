package e2e_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"maps"
	"slices"

	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/config"
)

var everySettingReported = []string{
	"reviewer_cli",
	"roles.reviewer.model",
	"roles.reviewer.effort",
	"diff_ceiling",
	"guidance_ceiling",
	"agent_deadline",
	"log_ceiling",
	"run_retention",
}

func TestWithNoConfigurationFileEverySettingReportsItsDefault(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config with no configuration file exited %d: %s", code, failures)
	}

	reported := settingsIn(t, printed)
	for _, key := range everySettingReported {
		setting, present := reported[key]
		if !present {
			t.Errorf("socrates review config reports nothing for %s:\n%s", key, printed)
			continue
		}
		if setting.Source != config.Default {
			t.Errorf("%s reads as %s with no configuration file at all", key, setting.Source)
		}
	}
	if _, err := os.Stat(configurationPath(t, session)); !os.IsNotExist(err) {
		t.Errorf("socrates review config created %s", configurationPath(t, session))
	}
}

func TestOnlyTheSettingsTheFileStatesReportAsSet(t *testing.T) {
	session := harness.Start(t)
	theFileStates := map[string]string{"diff_ceiling": "4096", "run_retention": "72h"}
	stateConfiguration(t, session, "diff_ceiling: 4096\nrun_retention: 72h\n")

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config exited %d: %s", code, failures)
	}

	reported := settingsIn(t, printed)
	for key, setting := range reported {
		value, statedInTheFile := theFileStates[key]
		want := config.Default
		if statedInTheFile {
			want = config.Set
		}
		if setting.Source != want {
			t.Errorf("%s reads as %s, and the file states %v", key, setting.Source, slices.Sorted(maps.Keys(theFileStates)))
		}
		if statedInTheFile && setting.Value != value {
			t.Errorf("%s reports %q, and the file states %q", key, setting.Value, value)
		}
	}
	if len(reported) != len(everySettingReported) {
		t.Errorf("socrates review config reports %d of the %d settings:\n%s", len(reported), len(everySettingReported), printed)
	}
}

func TestTheModelAndEffortReportedAreTheConfiguredReviewerCLIs(t *testing.T) {
	session := harness.Start(t)
	stateConfiguration(t, session, "reviewer_cli: codex\n")

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config exited %d: %s", code, failures)
	}

	reported := settingsIn(t, printed)
	codex, _ := agent.Named("codex")
	shipped, _ := agent.Named(agent.DefaultVendor)

	if model := reported["roles.reviewer.model"].Value; model != codex.Routing[agent.Reviewer].Model {
		t.Errorf("a machine configured for codex reports the model %q, and codex gets %q", model, codex.Routing[agent.Reviewer].Model)
	}
	if model := reported["roles.reviewer.model"].Value; model == shipped.Routing[agent.Reviewer].Model {
		t.Errorf("the model reported is %q, which is what %s gets, and the file names codex", model, agent.DefaultVendor)
	}
	if effort := reported["roles.reviewer.effort"].Value; effort != codex.Routing[agent.Reviewer].Effort {
		t.Errorf("a machine configured for codex reports the effort %q, and codex gets %q", effort, codex.Routing[agent.Reviewer].Effort)
	}
}

func TestAReviewerCLISocratesReviewWillNotDriveStillReports(t *testing.T) {
	session := harness.Start(t)
	unsupported := agent.Unsupported()[0].Name
	stateConfiguration(t, session, "reviewer_cli: "+unsupported+"\n")

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config with %s configured exited %d: %s", unsupported, code, failures)
	}
	if reported := settingsIn(t, printed)["reviewer_cli"]; reported.Value != unsupported || reported.Source != config.Set {
		t.Errorf("reviewer_cli reports %q as %s, and the file states %s", reported.Value, reported.Source, unsupported)
	}

	_, refused, beginCode := session.Socrates("begin")
	if beginCode != exit.Unusable {
		t.Fatalf("socrates review begin with %s configured exited %d, want %d", unsupported, beginCode, exit.Unusable)
	}
	if !strings.Contains(refused, unsupported) {
		t.Errorf("socrates review begin does not refuse %s by name: %s", unsupported, refused)
	}
}

func TestTheFileNamedIsTheOneUnderTheStatedHome(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config exited %d: %s", code, failures)
	}
	if first := firstLine(printed); first != configurationPath(t, session) {
		t.Errorf("socrates review config names %q, and the file under the stated home is %s", first, configurationPath(t, session))
	}
}

func TestAConfigurationFileThatWillNotParseIsRefusedByConfigToo(t *testing.T) {
	session := harness.Start(t)
	stateConfiguration(t, session, "roles: [not a mapping\n")

	printed, failures, code := session.Socrates("config")
	if code != exit.Unusable {
		t.Fatalf("socrates review config with a file that does not parse exited %d, want %d", code, exit.Unusable)
	}
	if !strings.Contains(failures, configurationPath(t, session)) {
		t.Errorf("the refusal does not name the file: %s", failures)
	}
	if strings.Contains(printed, string(config.Default)) {
		t.Errorf("socrates review config reported settings while the file will not parse:\n%s", printed)
	}
}

func TestConfigReportsNothingOfAgorasSettings(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates review config exited %d: %s", code, failures)
	}
	for key := range settingsIn(t, printed) {
		if strings.HasPrefix(key, "surface.") {
			t.Errorf("socrates review config reports %s, which socrates agora reads from a file of its own:\n%s", key, printed)
		}
	}
}

func TestASurfaceBlockLeftInReviewsFileIsRefusedByNamingAgorasFile(t *testing.T) {
	session := harness.Start(t)
	session.Simulator.Scenario(cleanScenario)
	session.Repository.Write("internal/added.go", "package internal\n")
	held := "diff_ceiling: 4096\nsurface:\n  bind: 127.0.0.1:18019\n"
	stateConfiguration(t, session, held)
	agoras := filepath.Join(session.Root, "agora", "config.yaml")

	for _, command := range []string{"config", "begin"} {
		printed, failures, code := session.Socrates(command)
		if code != exit.Unusable {
			t.Errorf("socrates review %s with a surface block in its file exited %d, want %d: %s%s", command, code, exit.Unusable, printed, failures)
		}
		for _, named := range []string{configurationPath(t, session), "surface", agoras} {
			if !strings.Contains(failures, named) {
				t.Errorf("socrates review %s refuses without naming %s: %s", command, named, failures)
			}
		}
	}

	if names := session.Names(); len(names) != 0 {
		t.Errorf("a run was begun while the file carried a surface block: %v", names)
	}
	if invoked := session.Simulator.Invocations(); len(invoked) != 0 {
		t.Errorf("a reviewer was spawned while the file carried a surface block: %v", invoked)
	}
	if after := fileAt(t, configurationPath(t, session)); after != held {
		t.Errorf("refusing the file rewrote it to %q, and it was written as %q", after, held)
	}
	if _, err := os.Stat(agoras); err == nil {
		t.Errorf("refusing the file wrote %s, and nothing moves a person's settings for them", agoras)
	}
}

func stateConfiguration(t *testing.T, session *harness.Session, content string) {
	t.Helper()

	path := configurationPath(t, session)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating the configuration directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}

func settingsIn(t *testing.T, printed string) map[string]config.Setting {
	t.Helper()

	reported := map[string]config.Setting{}
	for line := range strings.SplitSeq(printed, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		source := config.Source(fields[len(fields)-1])
		if source != config.Set && source != config.Default {
			continue
		}
		reported[fields[0]] = config.Setting{Key: fields[0], Value: strings.Join(fields[1:len(fields)-1], " "), Source: source}
	}
	if len(reported) == 0 {
		t.Fatalf("socrates review config reported no setting at all:\n%s", printed)
	}
	return reported
}
