package e2e_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/socrates-tools/socrates/implement/config"
	"github.com/socrates-tools/socrates/implement/skill"
	"github.com/socrates-tools/socrates/internal/exit"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/internal/testing/harness"
)

const implementConcurrencyKey = "concurrency"

func openImplementMenu(session *harness.Session) *harness.Terminal {
	menu := session.TopInATerminal("implement", "configure")
	menu.Await("roles.implementer.effort")
	return menu
}

func installedDefinitionIn(t *testing.T, session *harness.Session) string {
	t.Helper()

	directory := t.TempDir()
	if _, failures, code := session.Top("implement", "skill", "--host", skillkit.ClaudeCode.Flag, "--agents-to", directory); code != exit.NothingAsked {
		t.Fatalf("socrates implement skill exited %d: %s", code, failures)
	}
	return filepath.Join(directory, skill.AgentFileName())
}

func repositoryTopLevel(t *testing.T, session *harness.Session) string {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(session.Repository.Root)
	if err != nil {
		t.Fatalf("resolving %s: %v", session.Repository.Root, err)
	}
	return resolved
}

func TestImplementConfigureWritesAConfirmedConcurrencyAndConfigReportsItAsSet(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)

	menu := openImplementMenu(session)
	menu.MoveTo(implementConcurrencyKey)
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "2")
	menu.Await(path)
	menu.Type("q")
	if _, code := menu.Wait(); code != exit.NothingAsked {
		t.Errorf("quitting the menu exited %d, want %d", code, exit.NothingAsked)
	}

	if held := fileAt(t, path); held != "concurrency: 2\n" {
		t.Errorf("the file holds %q, want only the setting that was confirmed", held)
	}
	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}
	if reported := implementSettingsIn(t, printed)[implementConcurrencyKey]; reported != [2]string{"2", string(config.Set)} {
		t.Errorf("socrates implement config reports concurrency as %v after the menu set it", reported)
	}
}

func TestImplementConfigureDrawsTheSameSettingsConfigReports(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	stateImplementConfiguration(t, session, "concurrency: 2\nroles:\n  implementer:\n    effort: max\n")

	printed, failures, code := session.Top("implement", "config")
	if code != exit.NothingAsked {
		t.Fatalf("socrates implement config exited %d: %s", code, failures)
	}
	menu := openImplementMenu(session)
	drawn := menu.AwaitSettling()
	menu.Type("q")
	menu.Wait()

	for key, reported := range implementSettingsIn(t, printed) {
		row, found := rowFor(drawn, key)
		if !found || row.Value != reported[0] || string(row.Source) != reported[1] {
			t.Errorf("the menu says %+v for %s, and socrates implement config says %v:\n%s", row, key, reported, drawn)
		}
	}
}

func TestImplementConfigureOffersTheModelAliasesAndTakesAFullIdentifierToo(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.model")
	menu.Type(harness.Enter)
	drawn := menu.Await("use the default")
	for _, alias := range config.ModelAliases() {
		if !strings.Contains(drawn, alias) {
			t.Errorf("the model list offers no %s:\n%s", alias, drawn)
		}
	}

	menu.Type(strings.Repeat(harness.Down, len(config.ModelAliases())))
	menu.AwaitSettling()
	typeValue(menu, "claude-opus-4-1-20250805")
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	if held := fileAt(t, path); !strings.Contains(held, "model: claude-opus-4-1-20250805") {
		t.Errorf("the file holds %q, and a full identifier was typed", held)
	}
}

func TestImplementConfigureRewritesTheInstalledDefinitionAtOnceForAModelAndAnEffort(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	definition := installedDefinitionIn(t, session)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.model")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "opus")
	menu.Await(path)
	if written := session.File(definition); !strings.Contains(written, "model: opus\n") {
		t.Errorf("the definition does not state the model that was confirmed:\n%s", written)
	}

	menu.MoveTo("roles.implementer.effort")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "max")
	drawn := menu.AwaitSettling()
	menu.Type("q")
	menu.Wait()

	written := session.File(definition)
	for _, want := range []string{"model: opus\n", "effort: max\n"} {
		if !strings.Contains(written, want) {
			t.Errorf("the definition does not state %q:\n%s", want, written)
		}
	}
	if !strings.Contains(drawn, definition) {
		t.Errorf("the menu does not say which definition it rewrote:\n%s", drawn)
	}
	if held := fileAt(t, path); !strings.Contains(held, "model: opus") || !strings.Contains(held, "effort: max") {
		t.Errorf("the file holds %q", held)
	}
}

func TestImplementConfigureWithNoDefinitionInstalledWritesTheSettingAndSaysHowTheDefinitionIsWritten(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.effort")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "low")
	drawn := menu.Await(path)
	menu.Type("q")
	menu.Wait()

	if !strings.Contains(drawn, "socrates implement skill") {
		t.Errorf("the menu does not say which command writes the definition:\n%s", drawn)
	}
	if held := fileAt(t, path); !strings.Contains(held, "effort: low") {
		t.Errorf("the file holds %q, and the setting is written whether or not a definition is installed", held)
	}
	if _, err := os.Stat(filepath.Join(session.Root, "implement", skill.RecordFile)); !os.IsNotExist(err) {
		t.Errorf("the menu recorded an installation: %v", err)
	}
}

func TestImplementConfigureUsingTheDefaultRemovesTheLineAndAnEmptiedBlockGoesWithItsLastSetting(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	stateImplementConfiguration(t, session, "# kept\nconcurrency: 2\nroles:\n  implementer:\n    effort: max\n")

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.effort")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")
	menu.Await(path)
	menu.Type("q")
	menu.Wait()

	if held := fileAt(t, path); held != "# kept\nconcurrency: 2\n" {
		t.Errorf("the file holds %q, want the emptied roles block gone and the rest kept", held)
	}
}

func TestImplementConfigureRefusesAValueTheLoaderRefusesBeforeTheFileIsTouched(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	held := "concurrency: 2\n"
	stateImplementConfiguration(t, session, held)

	menu := openImplementMenu(session)
	menu.MoveTo(implementConcurrencyKey)
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "0")
	menu.Await("positive whole number")
	menu.Type(harness.Escape)
	menu.Await("roles.implementer.effort")
	menu.Type("q")
	menu.Wait()

	if after := fileAt(t, path); after != held {
		t.Errorf("the file holds %q after a refused value, and it held %q", after, held)
	}
}

func TestImplementConfigureOutsideARepositoryOffersNoDefaultBranch(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	menu := session.TopInATerminalIn(t.TempDir(), "implement", "configure")
	drawn := menu.Await("roles.implementer.effort")
	menu.Type("q")
	menu.Wait()

	if strings.Contains(drawn, "default_branch") {
		t.Errorf("the menu offers a default branch where there is no repository:\n%s", drawn)
	}
}

func TestImplementConfigureInsideARepositoryRemembersItsDefaultBranchAndRemovingItLeavesOtherRepositoriesAlone(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	session.Repository.Git("branch", "dev")
	here := repositoryTopLevel(t, session)
	stateImplementConfiguration(t, session, "default_branches:\n  /elsewhere/other.repo: main\n")

	menu := session.TopInATerminal("implement", "configure")
	menu.Await("default_branch")
	menu.MoveTo("default_branch")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "dev")
	drawn := menu.Await(path)
	if row, found := rowFor(drawn, "default_branch"); !found || row.Value != "dev" || string(row.Source) != string(config.Set) {
		t.Errorf("the menu came back reading %+v for default_branch:\n%s", row, drawn)
	}
	held := fileAt(t, path)
	if !strings.Contains(held, here+": dev") || !strings.Contains(held, "/elsewhere/other.repo: main") {
		t.Errorf("the file holds %q after remembering a branch", held)
	}

	menu.MoveTo("default_branch")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")
	menu.Await("now uses the default")
	menu.Type("q")
	menu.Wait()

	after := fileAt(t, path)
	if after != "default_branches:\n  /elsewhere/other.repo: main\n" {
		t.Errorf("the file holds %q after the entry was removed, want only the other repository's", after)
	}
}

func TestImplementConfigureRefusesABranchTheRepositoryDoesNotHoldBeforeTheFileIsTouched(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)

	menu := session.TopInATerminal("implement", "configure")
	menu.Await("default_branch")
	menu.MoveTo("default_branch")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "nonesuch")
	drawn := menu.Await("holds no local or remote-tracking branch")
	menu.Type(harness.Escape)
	menu.Await("roles.implementer.effort")
	menu.Type("q")
	menu.Wait()

	if !strings.Contains(drawn, "nonesuch") {
		t.Errorf("the refusal does not name the branch:\n%s", drawn)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("a refused branch created %s", path)
	}
}

func TestImplementConfigureWithNoTerminalRefusesNamingImplementConfigAndWritesNothing(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)

	printed, failures, code := session.Top("implement", "configure")
	if code != exit.Unusable {
		t.Fatalf("socrates implement configure with no terminal exited %d, want %d:\n%s%s", code, exit.Unusable, printed, failures)
	}
	if !strings.Contains(failures, "socrates implement config") {
		t.Errorf("the refusal does not name the command that reads settings without a terminal: %s", failures)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("a refused socrates implement configure created %s", path)
	}
}

func TestImplementConfigureTakesNoArguments(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	_, failures, code := session.Top("implement", "configure", "extra")
	if code != exit.Unusable || !strings.Contains(failures, "socrates implement configure") {
		t.Errorf("an argument exited %d with %q", code, failures)
	}
}

func installedCodexRoleIn(t *testing.T, session *harness.Session) string {
	t.Helper()

	directory := t.TempDir()
	if _, failures, code := session.Top("implement", "skill", "--host", skillkit.Codex.Flag, "--agents-to", directory); code != exit.NothingAsked {
		t.Fatalf("socrates implement skill --host codex exited %d: %s", code, failures)
	}
	return filepath.Join(directory, skill.CodexRoleFileName())
}

func TestImplementConfigureOffersACodexModelTakingATypedIdentifierOrTheDefaultAndACodexEffortOfferingCodexsEfforts(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.codex.model")
	menu.Type(harness.Enter)
	model := menu.Await("use the default")
	menu.Type(harness.Escape)
	menu.AwaitSettling()
	menu.MoveTo("roles.implementer.codex.effort")
	menu.Type(harness.Enter)
	effort := menu.Await("use the default")
	menu.Type(harness.Escape)
	menu.AwaitSettling()
	menu.Type("q")
	menu.Wait()

	offered := optionsIn(model)
	if len(offered) != 2 || offered[1] != "use the default" {
		t.Errorf("the Codex model row offers %v, want a typed identifier and use the default", offered)
	}
	for _, each := range config.CodexEfforts() {
		if !slices.ContainsFunc(optionsIn(effort), func(row string) bool { return strings.Fields(row)[0] == each }) {
			t.Errorf("the Codex effort row offers no %s:\n%s", each, effort)
		}
	}
}

func TestImplementConfigureRewritesTheInstalledCodexRoleAtOnceAndLeavesClaudeCodesDefinitionAlone(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	definition := installedDefinitionIn(t, session)
	role := installedCodexRoleIn(t, session)
	atInstall := session.File(definition)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.codex.model")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	typeValue(menu, "gpt-5.5")
	menu.Await(path)
	if written := session.File(role); !strings.Contains(written, "model = \"gpt-5.5\"\n") {
		t.Errorf("the role does not state the Codex model that was confirmed:\n%s", written)
	}

	menu.MoveTo("roles.implementer.codex.effort")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "ultra")
	drawn := menu.AwaitSettling()
	menu.Type("q")
	menu.Wait()

	written := session.File(role)
	for _, want := range []string{"model = \"gpt-5.5\"\n", "model_reasoning_effort = \"ultra\"\n"} {
		if !strings.Contains(written, want) {
			t.Errorf("the role does not state %q:\n%s", want, written)
		}
	}
	if !strings.Contains(drawn, role) {
		t.Errorf("the menu does not say which role it rewrote:\n%s", drawn)
	}
	if session.File(definition) != atInstall {
		t.Error("confirming a Codex setting rewrote Claude Code's definition")
	}
}

func TestImplementConfigureUsingTheDefaultForTheCodexModelTakesTheModelLineOutOfTheRole(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	stateImplementConfiguration(t, session, "roles:\n  implementer:\n    codex:\n      model: gpt-5.5\n")
	role := installedCodexRoleIn(t, session)
	if !strings.Contains(session.File(role), "model = \"gpt-5.5\"") {
		t.Fatalf("the role installed from the file does not state its model:\n%s", session.File(role))
	}

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.codex.model")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "use the default")
	menu.Await("now uses the default")
	menu.Type("q")
	menu.Wait()

	if strings.Contains(session.File(role), "\nmodel = ") {
		t.Errorf("the role still carries a model line:\n%s", session.File(role))
	}
	if held := fileAt(t, path); strings.Contains(held, "gpt-5.5") {
		t.Errorf("the file still holds the model: %q", held)
	}
}

func TestImplementConfigureRewritesTheInstalledClaudeCodeDefinitionAndLeavesTheCodexRoleAlone(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	definition := installedDefinitionIn(t, session)
	role := installedCodexRoleIn(t, session)
	atInstall := session.File(role)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.model")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "opus")
	menu.AwaitSettling()
	menu.Type("q")
	menu.Wait()

	if !strings.Contains(session.File(definition), "model: opus\n") {
		t.Errorf("the definition was not rewritten:\n%s", session.File(definition))
	}
	if session.File(role) != atInstall {
		t.Error("confirming a Claude Code setting rewrote the Codex role")
	}
}

func TestImplementConfigureWithNoCodexRoleInstalledWritesTheSettingAndNamesSocratesImplementSkill(t *testing.T) {
	t.Parallel()
	session := harness.Start(t)
	path := implementConfigurationPath(session)
	installedDefinitionIn(t, session)

	menu := openImplementMenu(session)
	menu.MoveTo("roles.implementer.codex.effort")
	menu.Type(harness.Enter)
	menu.Await("use the default")
	chooseOption(t, menu, "low")
	drawn := menu.Await(path)
	menu.Type("q")
	menu.Wait()

	if !strings.Contains(drawn, "socrates implement skill") || !strings.Contains(drawn, "Codex") {
		t.Errorf("the menu does not say which command writes the Codex role:\n%s", drawn)
	}
	if held := fileAt(t, path); !strings.Contains(held, "effort: low") {
		t.Errorf("the file holds %q, and the setting is written whether or not a role is installed", held)
	}
}
