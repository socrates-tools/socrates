package config_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/internal/refusal"
	"github.com/socrates-tools/socrates/internal/state"
)

func TestNoConfigurationFile(t *testing.T) {
	directory := t.TempDir()
	path := filepath.Join(directory, "socrates", "config.yaml")

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("loading an absent configuration file: %v", err)
	}
	if diff := cmp.Diff(config.Config{}, loaded); diff != "" {
		t.Errorf("an absent configuration file is not the defaults (-want +got):\n%s", diff)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("reading %s: %v", directory, err)
	}
	if len(entries) != 0 {
		t.Errorf("loading an absent configuration file created %v", entries)
	}
}

func TestAnEmptyFileIsTheDefaults(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, ""))
	if err != nil {
		t.Fatalf("loading an empty configuration file: %v", err)
	}
	if diff := cmp.Diff(config.Config{}, loaded); diff != "" {
		t.Errorf("an empty configuration file is not the defaults (-want +got):\n%s", diff)
	}
}

func TestEverySettingIsRead(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, `concurrency: 2
roles:
  implementer:
    model: claude-opus-4-1
    effort: max
`))
	if err != nil {
		t.Fatalf("loading a configuration file naming every setting: %v", err)
	}

	reported := valuesOf(loaded.Settings())
	want := map[string]string{
		"concurrency":              "2",
		"roles.implementer.model":  "claude-opus-4-1",
		"roles.implementer.effort": "max",
	}
	if diff := cmp.Diff(want, reported); diff != "" {
		t.Errorf("what the file states is not what is in force (-want +got):\n%s", diff)
	}
}

func TestAValueOutOfRangeIsRefusedNamingTheFile(t *testing.T) {
	for name, content := range map[string]string{
		"a concurrency of zero":      "concurrency: 0\n",
		"a negative concurrency":     "concurrency: -3\n",
		"a fractional concurrency":   "concurrency: 2.5\n",
		"an effort that is unknown":  "roles:\n  implementer:\n    effort: extreme\n",
		"a model with a space in it": "roles:\n  implementer:\n    model: \"claude opus\"\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := write(t, content)

			_, err := config.LoadFrom(path)
			if !refusal.Is(err) {
				t.Fatalf("loading %s returned %v, want a refusal", name, err)
			}
			if !strings.Contains(err.Error(), path) {
				t.Errorf("the refusal reads %q, without naming %s", err, path)
			}
		})
	}
}

func TestAFieldTheFileDoesNotHaveIsRefusedNamingTheFile(t *testing.T) {
	for name, content := range map[string]string{
		"a top-level field":   "parallelism: 3\n",
		"a role that is not":  "roles:\n  reviewer:\n    model: sonnet\n",
		"a field of the role": "roles:\n  implementer:\n    temperature: 1\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := write(t, content)

			_, err := config.LoadFrom(path)
			if !refusal.Is(err) {
				t.Fatalf("loading %s returned %v, want a refusal", name, err)
			}
			if !strings.Contains(err.Error(), path) {
				t.Errorf("the refusal reads %q, without naming %s", err, path)
			}
		})
	}
}

func TestAFileThatWillNotParseIsRefusedNamingTheFile(t *testing.T) {
	path := write(t, "roles: [not a mapping\n")

	_, err := config.LoadFrom(path)
	if !refusal.Is(err) || !strings.Contains(err.Error(), path) {
		t.Errorf("loading a file that does not parse returned %v, want a refusal naming %s", err, path)
	}
}

func TestTheConfigurationSitsInImplementsOwnDirectoryUnderTheStatedHome(t *testing.T) {
	root := t.TempDir()
	t.Setenv(state.HomeVariable, root)

	path, err := config.Path()
	if err != nil {
		t.Fatalf("resolving the configuration file: %v", err)
	}
	if want := filepath.Join(root, "implement", "config.yaml"); path != want {
		t.Errorf("the configuration file is %s, want %s", path, want)
	}
}

func TestEverySettingIsReportedInTheOrderTheFileIsWritten(t *testing.T) {
	var keys []string
	for _, setting := range (config.Config{}).Settings() {
		keys = append(keys, setting.Key)
	}
	want := []string{"concurrency", "roles.implementer.model", "roles.implementer.effort"}
	if diff := cmp.Diff(want, keys); diff != "" {
		t.Errorf("what socrates implement reports differs from what the file holds (-want +got):\n%s", diff)
	}
}

func TestWithNoFileAtAllEverySettingIsADefault(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		if setting.Source != config.Default {
			t.Errorf("%s reads as %s with no configuration file at all", setting.Key, setting.Source)
		}
		if setting.Default != setting.Value {
			t.Errorf("%s reports the value %q and the default %q, and nothing is set", setting.Key, setting.Value, setting.Default)
		}
	}
}

func TestOnlyWhatTheFileStatesReadsAsSet(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "concurrency: 7\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming one setting: %v", err)
	}

	set := map[string]string{}
	for _, setting := range loaded.Settings() {
		if setting.Source == config.Set {
			set[setting.Key] = setting.Value
		}
	}
	if diff := cmp.Diff(map[string]string{"concurrency": "7"}, set); diff != "" {
		t.Errorf("what reads as set differs from what the file states (-want +got):\n%s", diff)
	}
}

func TestASetSettingStillCarriesTheDefaultItWouldFallBackTo(t *testing.T) {
	loaded, err := config.LoadFrom(write(t, "concurrency: 7\n"))
	if err != nil {
		t.Fatalf("loading a configuration file naming one setting: %v", err)
	}
	for _, setting := range loaded.Settings() {
		if setting.Key == "concurrency" && setting.Default != strconv.Itoa(config.DefaultConcurrency) {
			t.Errorf("concurrency carries the default %q, and the package states %d", setting.Default, config.DefaultConcurrency)
		}
	}
}

func TestEveryReportedDefaultIsTheOneItsOwnerStates(t *testing.T) {
	drifted := defaultsThatDrifted((config.Config{}).Settings())
	for _, one := range drifted {
		t.Error(one)
	}
}

func TestADefaultWrittenOutRatherThanReadFromItsOwnerIsCaught(t *testing.T) {
	hardcoded := []config.Setting{{Key: "concurrency", Value: "5", Source: config.Default}}

	drifted := defaultsThatDrifted(hardcoded)
	if len(drifted) != 1 || !strings.Contains(drifted[0], "concurrency") {
		t.Errorf("a concurrency written out as 5 rather than read from its owner reports %v", drifted)
	}
}

func TestTheDefaultsAreTheOnesTheSpecStates(t *testing.T) {
	if config.DefaultConcurrency != 4 || config.DefaultModel != "sonnet" || config.DefaultEffort != "high" {
		t.Errorf("the defaults are %d, %s and %s, and the spec states 4, sonnet and high", config.DefaultConcurrency, config.DefaultModel, config.DefaultEffort)
	}
}

func TestOnlyTheEffortOffersAFixedSetOfValues(t *testing.T) {
	for _, setting := range (config.Config{}).Settings() {
		offered := len(setting.Choices) > 0
		if offered != (setting.Key == "roles.implementer.effort") {
			t.Errorf("%s offers %v, and only the effort has a fixed set of values", setting.Key, setting.Choices)
		}
	}
	want := []string{"low", "medium", "high", "xhigh", "max"}
	if diff := cmp.Diff(want, config.Efforts()); diff != "" {
		t.Errorf("the efforts offered differ from the ones the spec states (-want +got):\n%s", diff)
	}
}

func defaultsThatDrifted(reported []config.Setting) []string {
	owners := map[string]string{
		"concurrency":              strconv.Itoa(config.DefaultConcurrency),
		"roles.implementer.model":  config.DefaultModel,
		"roles.implementer.effort": config.DefaultEffort,
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
		if setting.Value != owned {
			drifted = append(drifted, setting.Key+" reports "+strconv.Quote(setting.Value)+", and its owning package states "+strconv.Quote(owned))
		}
	}
	return drifted
}

func valuesOf(settings []config.Setting) map[string]string {
	values := map[string]string{}
	for _, setting := range settings {
		values[setting.Key] = setting.Value
	}
	return values
}

func write(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
	return path
}
