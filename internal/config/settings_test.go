package config_test

import (
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/agent"
	"github.com/socrates-tools/socrates/internal/config"
	"github.com/socrates-tools/socrates/internal/guidance"
	"github.com/socrates-tools/socrates/internal/loopback"
	"github.com/socrates-tools/socrates/internal/prompt"
	"github.com/socrates-tools/socrates/internal/runlog"
	"github.com/socrates-tools/socrates/internal/runstore"
)

var everySetting = []string{
	"reviewer_cli",
	"roles.reviewer.model",
	"roles.reviewer.effort",
	"diff_ceiling",
	"guidance_ceiling",
	"agent_deadline",
	"log_ceiling",
	"run_retention",
	"surface.bind",
	"surface.reached_at",
}

func TestEverySettingIsReportedInTheOrderTheFileIsWritten(t *testing.T) {
	var keys []string
	for _, setting := range (config.Config{}).Settings() {
		keys = append(keys, setting.Key)
	}
	if diff := cmp.Diff(everySetting, keys); diff != "" {
		t.Errorf("what socrates-review reports about itself differs from what the file holds (-want +got):\n%s", diff)
	}
}

func TestWithNoFileAtAllEverySettingIsADefault(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		if setting.Source != config.Default {
			t.Errorf("%s reads as %s with no configuration file at all", setting.Key, setting.Source)
		}
	}
}

func TestOnlyWhatTheFileStatesReadsAsSet(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, `reviewer_cli: codex
run_retention: 72h
`))
	if err != nil {
		t.Fatalf("loading a configuration file naming two settings: %v", err)
	}

	set := map[string]string{}
	for _, setting := range loaded.Settings() {
		if setting.Source == config.Set {
			set[setting.Key] = setting.Value
		}
	}
	if diff := cmp.Diff(map[string]string{"reviewer_cli": "codex", "run_retention": "72h"}, set); diff != "" {
		t.Errorf("what reads as set differs from what the file states (-want +got):\n%s", diff)
	}
}

func TestTheModelAndEffortAreTheOnesTheConfiguredReviewerCLIGets(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "reviewer_cli: codex\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming codex: %v", err)
	}

	reported := valuesOf(loaded.Settings())
	codex, _ := agent.Named("codex")
	if want := codex.Routing[agent.Reviewer]; reported["roles.reviewer.model"] != want.Model || reported["roles.reviewer.effort"] != want.Effort {
		t.Errorf("a machine configured for codex reports the model %q and the effort %q, and codex gets %q and %q",
			reported["roles.reviewer.model"], reported["roles.reviewer.effort"], want.Model, want.Effort)
	}
	if reported["roles.reviewer.model"] != "" {
		t.Errorf("the model reported under codex is %q, and codex picks its own", reported["roles.reviewer.model"])
	}
}

func TestAReviewerCLISocratesReviewDoesNotDriveIsReportedRatherThanCorrected(t *testing.T) {
	unsupported := agent.Unsupported()[0].Name

	loaded, err := config.LoadFrom(write(t, "reviewer_cli: "+unsupported+"\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming %s: %v", unsupported, err)
	}

	reported := valuesOf(loaded.Settings())
	if reported["reviewer_cli"] != unsupported {
		t.Errorf("reviewer_cli reports %q, and the file states %s, which socrates-review refuses to drive when a review starts", reported["reviewer_cli"], unsupported)
	}
	for _, key := range []string{"roles.reviewer.model", "roles.reviewer.effort"} {
		if reported[key] != "" {
			t.Errorf("%s reports %q, and socrates-review routes nothing to a CLI it does not drive", key, reported[key])
		}
	}

	shipped, _ := agent.Named(agent.DefaultVendor)
	if reported["roles.reviewer.model"] == shipped.Routing[agent.Reviewer].Model {
		t.Error("the model reported is the one socrates-review shipped with, and the file names another CLI")
	}
}

func TestASettingNothingHasAValueForIsReportedEmpty(t *testing.T) {
	reported := valuesOf((config.Config{}).Settings())
	if reported["surface.reached_at"] != "" {
		t.Errorf("surface.reached_at reports %q, and nothing states one", reported["surface.reached_at"])
	}
}

func TestEveryReportedDefaultIsTheOneItsOwnerStates(t *testing.T) {
	for _, of := range []config.Config{{}, {ReviewerCLI: "codex"}} {
		drifted := defaultsThatDrifted(of, of.Settings())
		for _, one := range drifted {
			t.Errorf("configured for %s, %s", of.Vendor(), one)
		}
	}
}

func TestADefaultWrittenOutRatherThanReadFromItsOwnerIsCaught(t *testing.T) {
	hardcoded := []config.Setting{{Key: "run_retention", Value: "168h", Source: config.Default}}

	drifted := defaultsThatDrifted(config.Config{}, hardcoded)
	if len(drifted) != 1 {
		t.Fatalf("a run_retention written out as 168h rather than read from runstore reports %v", drifted)
	}
	if !strings.Contains(drifted[0], "run_retention") {
		t.Errorf("the drift reads %q, without naming the setting that drifted", drifted[0])
	}
}

func defaultsThatDrifted(of config.Config, reported []config.Setting) []string {
	vendor, _ := agent.Named(of.Vendor())
	routing := vendor.Routing[agent.Reviewer]

	owners := map[string]func(string) bool{
		"reviewer_cli":          reads(agent.DefaultVendor),
		"roles.reviewer.model":  reads(routing.Model),
		"roles.reviewer.effort": reads(routing.Effort),
		"diff_ceiling":          counts(int64(prompt.DefaultCeiling)),
		"guidance_ceiling":      counts(guidance.DefaultCeiling),
		"agent_deadline":        lasts(agent.DefaultDeadline),
		"log_ceiling":           counts(runlog.DefaultCeiling),
		"run_retention":         lasts(runstore.DefaultRetention),
		"surface.bind":          reads(loopback.DefaultAddress),
		"surface.reached_at":    reads(""),
	}

	var drifted []string
	for _, setting := range reported {
		if setting.Source != config.Default {
			continue
		}
		owned, known := owners[setting.Key]
		if !known {
			drifted = append(drifted, setting.Key+" is reported and no package owns a default for it")
			continue
		}
		if !owned(setting.Value) {
			drifted = append(drifted, setting.Key+" reports "+strconv.Quote(setting.Value)+", and its owning package states something else")
		}
	}
	return drifted
}

func reads(owned string) func(string) bool {
	return func(reported string) bool { return reported == owned }
}

func counts(owned int64) func(string) bool {
	return func(reported string) bool {
		read, err := strconv.ParseInt(reported, 10, 64)
		return err == nil && read == owned
	}
}

func lasts(owned time.Duration) func(string) bool {
	return func(reported string) bool {
		read, err := time.ParseDuration(reported)
		return err == nil && read == owned
	}
}

func valuesOf(settings []config.Setting) map[string]string {
	values := map[string]string{}
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}
	return values
}

func TestEverySettingCarriesTheDefaultItFallsBackTo(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		if setting.Default != setting.Value {
			t.Errorf("%s reports the value %q and the default %q, and nothing is set", setting.Key, setting.Value, setting.Default)
		}
	}
}

func TestASetSettingStillCarriesTheDefaultItWouldFallBackTo(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "run_retention: 72h\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming one setting: %v", err)
	}

	for _, setting := range loaded.Settings() {
		if setting.Key != "run_retention" {
			continue
		}
		if setting.Source != config.Set || setting.Value != "72h" {
			t.Fatalf("run_retention reads %q as %s, and the file states 72h", setting.Value, setting.Source)
		}
		if want := lasting(runstore.DefaultRetention); setting.Default != want {
			t.Errorf("run_retention carries the default %q, and runstore states %q", setting.Default, want)
		}
	}
}

func TestOnlyTheReviewerCLIOffersAFixedSetOfValues(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		offered := len(setting.Choices) > 0
		if offered != (setting.Key == "reviewer_cli") {
			t.Errorf("%s offers %v, and only reviewer_cli has a fixed set of values", setting.Key, setting.Choices)
		}
	}
}

func TestTheReviewerCLIsValuesAreTheOnesSocratesReviewDrives(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		if setting.Key != "reviewer_cli" {
			continue
		}
		if diff := cmp.Diff(agent.SupportedNames(), setting.Choices); diff != "" {
			t.Errorf("reviewer_cli offers something other than the CLIs socrates-review drives (-want +got):\n%s", diff)
		}
		for _, unsupported := range agent.Unsupported() {
			if slices.Contains(setting.Choices, unsupported.Name) {
				t.Errorf("reviewer_cli offers %s, which a review refuses", unsupported.Name)
			}
		}
	}
}

func lasting(span time.Duration) string {
	return strconv.FormatInt(int64(span/time.Hour), 10) + "h"
}
