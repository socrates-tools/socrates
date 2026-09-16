package e2e_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/internal/exit"
	"github.com/socrates-tools/socrates/internal/testing/harness"
	"github.com/socrates-tools/socrates/review/agent"
	"github.com/socrates-tools/socrates/review/config"
)

func TestTheMenuListsEverySettingWithTheValueInForce(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	drawn := menu.Await(everySettingReported[len(everySettingReported)-1])

	for _, key := range everySettingReported {
		if !strings.Contains(drawn, key) {
			t.Errorf("the menu names nothing for %s:\n%s", key, drawn)
		}
	}
	if !strings.Contains(drawn, string(config.Default)) {
		t.Errorf("the menu says of no setting that its value is a default:\n%s", drawn)
	}

	menu.Type("q")
	_, code := menu.Wait()
	if code != exit.NothingAsked {
		t.Errorf("quitting the menu exited %d, want %d", code, exit.NothingAsked)
	}
}

func TestTheCursorMovesUpAndDownTheSettings(t *testing.T) {
	session := harness.Start(t)
	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[len(everySettingReported)-1])

	first := settingUnderTheCursor(t, menu.AwaitSettling())
	if first != everySettingReported[0] {
		t.Errorf("the menu opens on %s, and the first setting is %s", first, everySettingReported[0])
	}

	menu.Type(harness.Down)
	if moved := settingUnderTheCursor(t, menu.AwaitSettling()); moved != everySettingReported[1] {
		t.Errorf("pressing down moved the cursor to %s, want %s", moved, everySettingReported[1])
	}

	menu.Type(harness.Up)
	if back := settingUnderTheCursor(t, menu.AwaitSettling()); back != first {
		t.Errorf("pressing up moved the cursor to %s, want %s", back, first)
	}

	menu.Type("q")
	menu.Wait()
}

func TestOpeningAndQuittingTheMenuWritesNoFile(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.Type(harness.Down)
	menu.Type(harness.Up)
	menu.Type("q")
	menu.Wait()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("opening the menu and quitting created %s", path)
	}
}

func TestTheMenuAndTheReportAgreeAboutWhatIsInForce(t *testing.T) {
	session := harness.Start(t)
	stateConfiguration(t, session, "reviewer_cli: codex\ndiff_ceiling: 4096\n")

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates-review config exited %d: %s", code, failures)
	}

	menu := session.SocratesInATerminal("configure")
	drawn := menu.Await(everySettingReported[len(everySettingReported)-1])
	menu.Type("q")
	menu.Wait()

	for key, reported := range settingsIn(t, printed) {
		row, found := rowFor(drawn, key)
		if !found {
			t.Errorf("the menu lists no row for %s:\n%s", key, drawn)
			continue
		}
		if row.Value != reported.Value || row.Source != reported.Source {
			t.Errorf("the menu says %s is %q as %s, and socrates-review config says %q as %s", key, row.Value, row.Source, reported.Value, reported.Source)
		}
	}
}

func TestChoosingAReviewerCLIWritesAFileHoldingOnlyThatSetting(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("clearing the configuration the harness wrote: %v", err)
	}

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")

	chooseOption(t, menu, "codex")
	menu.Await("codex")
	menu.Type("q")
	menu.Wait()

	if held := fileAt(t, path); held != "reviewer_cli: codex\n" {
		t.Errorf("the file holds %q, want only the setting that was chosen", held)
	}
}

func TestTheListOffersOnlyTheReviewerCLIsSocratesReviewDrives(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	drawn := menu.Await("use the default")

	for _, supported := range agent.SupportedNames() {
		if !strings.Contains(drawn, supported) {
			t.Errorf("the list offers no %s, and socrates-review drives it:\n%s", supported, drawn)
		}
	}
	for _, unsupported := range agent.Unsupported() {
		if strings.Contains(drawn, unsupported.Name) {
			t.Errorf("the list offers %s, and a review refuses it:\n%s", unsupported.Name, drawn)
		}
	}
	if !strings.Contains(drawn, "in force") {
		t.Errorf("nothing in the list is marked as the one in force:\n%s", drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[1])
	menu.Type("q")
	menu.Wait()
}

func TestTheConfirmationNamesTheSettingTheValueAndThePath(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "codex")

	drawn := menu.Await(path)
	for _, named := range []string{"reviewer_cli", "codex", path} {
		if !strings.Contains(drawn, named) {
			t.Errorf("the confirmation names no %s:\n%s", named, drawn)
		}
	}
	if !strings.Contains(drawn, "next round of a run already going") {
		t.Errorf("the confirmation does not say what the change reaches:\n%s", drawn)
	}

	if row, found := rowFor(drawn, "reviewer_cli"); !found || row.Value != "codex" || row.Source != config.Set {
		t.Errorf("the menu came back reading %+v, want codex as set:\n%s", row, drawn)
	}

	menu.Type("q")
	menu.Wait()
}

func TestACommentAboveAKeySurvivesAnUnrelatedSetting(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "# the ceiling I chose, and why\ndiff_ceiling: 4096\n\n# I keep runs for three days\nrun_retention: 72h\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "codex")
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	after := fileAt(t, path)
	if after != held+"reviewer_cli: codex\n" {
		t.Errorf("the file now holds:\n%s\nand it held:\n%s", after, held)
	}
}

func TestTwoSettingsInOneSittingLeaveBothAndQuittingAfterOneLeavesThatOne(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("clearing the configuration the harness wrote: %v", err)
	}

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])

	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "codex")
	menu.Await(path)

	if held := fileAt(t, path); !strings.Contains(held, "reviewer_cli: codex") {
		t.Fatalf("the first setting was not written before the second was chosen: %q", held)
	}

	menu.MoveTo("run_retention")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "72h")
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	held := fileAt(t, path)
	if !strings.Contains(held, "reviewer_cli: codex") || !strings.Contains(held, "run_retention: 72h") {
		t.Errorf("the file holds %q, and two settings were chosen in one sitting", held)
	}
}

func TestAValueIsWrittenEvenWhereTheNamedReviewerCLIIsNotInstalled(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)

	elsewhere := ""
	for _, supported := range agent.SupportedNames() {
		if supported != session.Vendor.Name {
			elsewhere = supported
		}
	}
	if elsewhere == "" {
		t.Fatal("only one reviewer CLI is supported, so this test has nothing to choose that is not installed")
	}

	nothingInstalled := t.TempDir()
	if found := lookIn(t, nothingInstalled, agent.SupportedNames()); len(found) != 0 {
		t.Fatalf("%v are reachable from a directory holding nothing, so this test proves nothing", found)
	}

	menu := session.SocratesInATerminalWith([]string{"PATH=" + nothingInstalled}, "configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, elsewhere)
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	if held := fileAt(t, path); !strings.Contains(held, "reviewer_cli: "+elsewhere) {
		t.Errorf("the file holds %q, and installing a reviewer CLI after configuring it is a legitimate order to do things in", held)
	}
}

func TestASettingWithNoFixedValuesOpensAFieldHoldingTheValueInForce(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("agent_deadline")
	menu.Type(harness.Enter)
	drawn := menu.Await("use the default")

	if !strings.Contains(drawn, "30m") {
		t.Errorf("the field does not hold the value in force:\n%s", drawn)
	}
	if strings.Contains(drawn, "choose") {
		t.Errorf("a setting with no fixed set of values offers a list to choose from:\n%s", drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()
}

func TestEverySettingWithNoFixedValuesOpensAFieldHoldingTheValueInForce(t *testing.T) {
	session := harness.Start(t)

	printed, failures, code := session.Socrates("config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates-review config exited %d: %s", code, failures)
	}
	reported := settingsIn(t, printed)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])

	typed := 0
	for _, key := range everySettingReported {
		menu.MoveTo(key)
		menu.Type(harness.Enter)
		drawn := menu.Await("use the default")

		if offersAList(drawn) {
			menu.Type(harness.Escape)
			menu.Await(everySettingReported[0])
			continue
		}
		typed++
		if field := fieldIn(t, drawn); field != reported[key].Value {
			t.Errorf("the field for %s holds %q, and the value in force is %q", key, field, reported[key].Value)
		}

		menu.Type(harness.Escape)
		menu.Await(everySettingReported[0])
	}

	menu.Type("q")
	menu.Wait()

	if typed != len(everySettingReported)-1 {
		t.Errorf("%d settings opened a field, and all but reviewer_cli should have", typed)
	}
}

func TestAnEmptyFieldIsRefusedRatherThanTakenForTheDefault(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "run_retention: 72h\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("run_retention")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	menu.Type(strings.Repeat("\x7f", 20))
	menu.AwaitSettling()
	menu.Type(harness.Enter)

	drawn := menu.Await("needs a value")
	if !strings.Contains(drawn, "use the default") {
		t.Errorf("what was said does not point at the way back to the default:\n%s", drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()

	if after := fileAt(t, path); after != held {
		t.Errorf("an emptied field left the file reading %q, and it was written as %q", after, held)
	}
}

func TestQuittingTheMenuIsAKeyChosenForIt(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.Type(harness.Escape)
	menu.AwaitSettling()

	if drawn := menu.Drawn(); !strings.Contains(drawn, everySettingReported[0]) {
		t.Errorf("escape left the menu, and the key that quits is q:\n%s", drawn)
	}

	menu.Type("q")
	if _, code := menu.Wait(); code != exit.NothingAsked {
		t.Errorf("quitting with q exited %d, want %d", code, exit.NothingAsked)
	}
}

func TestASettingWithNoValueInForceOpensAnEmptyField(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("surface.reached_at")
	menu.Type(harness.Enter)
	drawn := menu.Await("use the default")

	if field := fieldIn(t, drawn); field != "" {
		t.Errorf("the field holds %q, and nothing states a value for surface.reached_at:\n%s", field, drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()
}

func TestADurationASizeAndAnAddressAreEachTypedAndWritten(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)

	typed := []struct {
		key   string
		value string
	}{
		{key: "agent_deadline", value: "45m"},
		{key: "diff_ceiling", value: "8192"},
		{key: "surface.bind", value: "127.0.0.1:18019"},
		{key: "roles.reviewer.model", value: "sonnet"},
		{key: "roles.reviewer.effort", value: "low"},
		{key: "guidance_ceiling", value: "8192"},
		{key: "log_ceiling", value: "2097152"},
		{key: "run_retention", value: "72h"},
		{key: "surface.reached_at", value: "http://localhost:18019/"},
	}

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	for _, one := range typed {
		menu.MoveTo(one.key)
		menu.Type(harness.Enter)
		menu.Await("use the default")
		typeValue(menu, one.value)
		menu.Await(one.key + " set to")
	}
	menu.Type("q")
	menu.Wait()

	loaded, err := config.LoadFrom(path)
	if err != nil {
		t.Fatalf("the menu wrote a file the loader refuses: %v\n%s", err, fileAt(t, path))
	}
	reported := map[string]config.Setting{}
	for _, setting := range loaded.Settings() {
		reported[setting.Key] = setting
	}
	for _, one := range typed {
		if reported[one.key].Value != one.value || reported[one.key].Source != config.Set {
			t.Errorf("%s reads %q as %s after %q was typed", one.key, reported[one.key].Value, reported[one.key].Source, one.value)
		}
	}
}

func TestNoListOfModelNamesIsOfferedAnywhere(t *testing.T) {
	session := harness.Start(t)
	shipped, _ := agent.Named(agent.DefaultVendor)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("roles.reviewer.model")
	menu.Type(harness.Enter)
	drawn := menu.Await("use the default")

	offered := optionsIn(drawn)
	if len(offered) != 2 {
		t.Errorf("the model offers %v, and a shipped list of model names goes stale:\n%s", offered, drawn)
	}
	if !strings.Contains(drawn, shipped.Routing[agent.Reviewer].Model) {
		t.Errorf("the field does not hold the model in force:\n%s", drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()
}

func TestAValueTheLoaderRefusesIsRefusedBeforeTheFileIsTouched(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "diff_ceiling: 4096\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("diff_ceiling")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "-1")

	drawn := menu.Await("negative")
	if !strings.Contains(drawn, "diff_ceiling") {
		t.Errorf("the refusal does not name the setting:\n%s", drawn)
	}
	if after := fileAt(t, path); after != held {
		t.Errorf("a value the loader refuses changed the file to %q", after)
	}
	if field := fieldIn(t, drawn); field != "-1" {
		t.Errorf("the person was dropped out of the field, which now reads %q:\n%s", field, drawn)
	}

	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()

	if after := fileAt(t, path); after != held {
		t.Errorf("the file reads %q after the refusal, and it was written as %q", after, held)
	}
}

func TestEscapingOutOfAFieldWritesNothing(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "diff_ceiling: 4096\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("diff_ceiling")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	menu.Type("9")
	menu.AwaitSettling()
	menu.Type(harness.Escape)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()

	if after := fileAt(t, path); after != held {
		t.Errorf("escaping out of a field left %q, and the file was written as %q", after, held)
	}
}

func TestEverySettingOffersTheDefaultAndNamesIt(t *testing.T) {
	session := harness.Start(t)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])

	for _, key := range everySettingReported {
		menu.MoveTo(key)
		menu.Type(harness.Enter)
		drawn := menu.Await("use the default")

		if !slices.ContainsFunc(optionsIn(drawn), func(offered string) bool { return strings.HasPrefix(offered, "use the default") }) {
			t.Errorf("%s offers no way back to the default:\n%s", key, drawn)
		}
		if named := defaultNamedIn(drawn); key != "surface.reached_at" && named == "" {
			t.Errorf("%s offers the default without naming what it would fall back to:\n%s", key, drawn)
		}

		menu.Type(harness.Escape)
		menu.Await(everySettingReported[0])
	}

	menu.Type("q")
	menu.Wait()
}

func TestUsingTheDefaultLeavesAFileThatNoLongerNamesTheSetting(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	stateConfiguration(t, session, "# the ceiling I chose\ndiff_ceiling: 4096\nrun_retention: 72h\n")

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("run_retention")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")

	drawn := menu.Await("now uses the default")
	menu.Type("q")
	menu.Wait()

	after := fileAt(t, path)
	if strings.Contains(after, "run_retention") {
		t.Errorf("the file still names run_retention:\n%s", after)
	}
	if !strings.Contains(after, "# the ceiling I chose") || !strings.Contains(after, "diff_ceiling: 4096") {
		t.Errorf("the removal took something else with it:\n%s", after)
	}
	if strings.Contains(after, "336h") {
		t.Errorf("the current default was written into the file rather than left out:\n%s", after)
	}
	if row, found := rowFor(drawn, "run_retention"); !found || row.Source != config.Default {
		t.Errorf("the menu came back reading %+v for run_retention:\n%s", row, drawn)
	}
}

func TestUsingTheDefaultOnTheOnlySettingLeavesAFileNamingNothing(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	stateConfiguration(t, session, "roles:\n  reviewer:\n    model: sonnet\n")

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("roles.reviewer.model")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	after := fileAt(t, path)
	if strings.TrimSpace(after) != "" {
		t.Errorf("the file holds %q, and nothing is set any more", after)
	}
	for _, key := range everySettingReported {
		if strings.Contains(after, key) {
			t.Errorf("the file names %s, and nothing was set:\n%s", key, after)
		}
	}
}

func TestUsingTheDefaultOnASettingAlreadyDefaultSaysSoRatherThanFailing(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "diff_ceiling: 4096\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("run_retention")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")

	drawn := menu.Await("already the default")
	if !strings.Contains(drawn, "run_retention") {
		t.Errorf("what was said names no setting:\n%s", drawn)
	}
	menu.Type("q")
	if _, code := menu.Wait(); code != exit.NothingAsked {
		t.Errorf("defaulting a setting that was already a default exited %d", code)
	}
	if after := fileAt(t, path); after != held {
		t.Errorf("the file reads %q, and it was written as %q", after, held)
	}
}

func TestAWriteThatCannotHappenIsSaidRatherThanSwallowed(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	if err := os.Chmod(session.Home, 0o555); err != nil {
		t.Fatalf("taking write permission off %s: %v", session.Home, err)
	}
	t.Cleanup(func() { _ = os.Chmod(session.Home, 0o755) })

	menu := session.SocratesInATerminal("configure")
	menu.Await(everySettingReported[0])
	menu.MoveTo("reviewer_cli")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "codex")

	drawn, code := menu.Wait()
	if code != exit.Unusable {
		t.Errorf("a write that could not happen exited %d, want %d:\n%s", code, exit.Unusable, drawn)
	}
	if !strings.Contains(drawn, path) {
		t.Errorf("nothing said names the file that could not be written:\n%s", drawn)
	}
	if strings.Contains(drawn, "set to codex") {
		t.Errorf("socrates-review said the setting was written and it was not:\n%s", drawn)
	}
}

func TestWithNoTerminalConfigureRefusesAndWritesNothing(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		t.Fatalf("clearing the configuration the harness wrote: %v", err)
	}

	printed, failures, code := session.Socrates("configure")
	if code != exit.Unusable {
		t.Fatalf("socrates-review configure with no terminal exited %d, want %d:\n%s%s", code, exit.Unusable, printed, failures)
	}
	if !strings.Contains(failures, "socrates-review config") {
		t.Errorf("the refusal does not name the command that reads settings without a terminal: %s", failures)
	}
	for _, library := range []string{"bubbletea", "charm", "tea", "ultraviolet"} {
		if strings.Contains(strings.ToLower(failures), library) {
			t.Errorf("the refusal names %s, which is a dependency the person never chose: %s", library, failures)
		}
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("a refused socrates-review configure created %s", path)
	}
}

func TestInputRedirectedFromAFileStillRefuses(t *testing.T) {
	session := harness.Start(t)

	_, failures, code := session.SocratesReading("q\n", "configure")
	if code != exit.Unusable {
		t.Fatalf("socrates-review configure with input redirected exited %d, want %d", code, exit.Unusable)
	}
	if !strings.Contains(failures, "socrates-review config") {
		t.Errorf("the refusal does not name the command that reads settings without a terminal: %s", failures)
	}
}

func TestAFileThatWillNotParseOpensAndNamesTheLineThatIsWrong(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	held := "diff_ceiling: 4096\nreviewer_temperature: 0.5\nrun_retention: 72h\n"
	stateConfiguration(t, session, held)

	menu := session.SocratesInATerminal("configure")
	drawn := menu.Await("reviewer_temperature")

	if !strings.Contains(drawn, "line 2") {
		t.Errorf("the line named is not the one that is wrong:\n%s", drawn)
	}
	if !strings.Contains(drawn, "reset reviewer_temperature") {
		t.Errorf("the only offer is not resetting the setting to its default:\n%s", drawn)
	}

	menu.Type("q")
	menu.Wait()
	if after := fileAt(t, path); after != held {
		t.Errorf("declining left the file reading %q, and it was written as %q", after, held)
	}
}

func TestResettingTheUnreadableSettingLeavesAFileTheLoaderAccepts(t *testing.T) {
	session := harness.Start(t)
	path := configurationPath(t, session)
	stateConfiguration(t, session, "# the ceiling I chose\ndiff_ceiling: 4096\n\nreviewer_temperature: 0.5\nrun_retention: 72h\n")

	menu := session.SocratesInATerminal("configure")
	menu.Await("reviewer_temperature")
	menu.Type(harness.Enter)
	menu.Await(everySettingReported[0])
	menu.Type("q")
	menu.Wait()

	after := fileAt(t, path)
	if strings.Contains(after, "reviewer_temperature") {
		t.Errorf("the setting socrates-review could not read is still there:\n%s", after)
	}
	for _, kept := range []string{"# the ceiling I chose", "diff_ceiling: 4096", "run_retention: 72h"} {
		if !strings.Contains(after, kept) {
			t.Errorf("%q was dropped along with the setting that was reset:\n%s", kept, after)
		}
	}
	if _, err := config.LoadFrom(path); err != nil {
		t.Errorf("the repaired file is still one the loader refuses: %v\n%s", err, after)
	}
}

func TestTheExceptionIsOneCommandWide(t *testing.T) {
	session := harness.Start(t)
	stateConfiguration(t, session, "reviewer_temperature: 0.5\n")

	for _, command := range [][]string{{"config"}, {"begin"}} {
		_, failures, code := session.Socrates(command...)
		if code != exit.Unusable {
			t.Errorf("socrates-review %s with a file socrates-review cannot read exited %d, want %d", command[0], code, exit.Unusable)
		}
		if !strings.Contains(failures, configurationPath(t, session)) {
			t.Errorf("socrates-review %s does not refuse by naming the file: %s", command[0], failures)
		}
	}
}

func offersAList(drawn string) bool {
	return strings.Contains(drawn, "enter to choose")
}

func lookIn(t *testing.T, directory string, named []string) []string {
	t.Helper()

	var found []string
	for _, one := range named {
		if _, err := os.Stat(filepath.Join(directory, one)); err == nil {
			found = append(found, one)
		}
	}
	return found
}

func chooseOption(t *testing.T, menu *harness.Terminal, label string) {
	t.Helper()

	for range 20 {
		drawn := menu.AwaitSettling()
		if strings.HasPrefix(harness.Under(drawn), label) {
			menu.Type(harness.Enter)
			return
		}
		menu.Type(harness.Down)
	}
	t.Fatalf("the cursor never reached %q, and the screen reads:\n%s", label, menu.Drawn())
}

func typeValue(menu *harness.Terminal, value string) {
	menu.Type(strings.Repeat("\x7f", 40))
	menu.AwaitSettling()
	menu.Type(value)
	menu.AwaitSettling()
	menu.Type(harness.Enter)
}

func optionsIn(drawn string) []string {
	var offered []string
	for line := range strings.SplitSeq(drawn, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "(") {
			continue
		}
		if after, marked := strings.CutPrefix(line, "> "); marked {
			offered = append(offered, strings.TrimSpace(after))
			continue
		}
		if strings.HasPrefix(line, "  ") {
			offered = append(offered, trimmed)
		}
	}
	return offered
}

func fieldIn(t *testing.T, drawn string) string {
	t.Helper()

	under := harness.Under(drawn)
	if under == "" && !strings.Contains(drawn, "use the default") {
		t.Fatalf("nothing on the screen is under the cursor:\n%s", drawn)
	}
	return strings.TrimSuffix(under, "_")
}

func defaultNamedIn(drawn string) string {
	for _, offered := range optionsIn(drawn) {
		if inside, is := strings.CutPrefix(offered, "use the default ("); is {
			return strings.TrimSuffix(inside, ")")
		}
	}
	return ""
}

func settingUnderTheCursor(t *testing.T, drawn string) string {
	t.Helper()

	under := strings.Fields(harness.Under(drawn))
	if len(under) == 0 {
		t.Fatalf("no setting is under the cursor:\n%s", drawn)
	}
	return under[0]
}

func rowFor(drawn, key string) (config.Setting, bool) {
	for line := range strings.SplitSeq(drawn, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == ">" {
			fields = fields[1:]
		}
		if fields[0] != key || len(fields) < 2 {
			continue
		}
		source := config.Source(fields[len(fields)-1])
		if source != config.Set && source != config.Default {
			continue
		}
		return config.Setting{Key: key, Value: strings.Join(fields[1:len(fields)-1], " "), Source: source}, true
	}
	return config.Setting{}, false
}

func fileAt(t *testing.T, path string) string {
	t.Helper()

	held, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ""
	}
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(held)
}
