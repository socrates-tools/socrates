package e2e_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/prompt"
	"github.com/rui-ktei/socrates/internal/testing/harness"
	"github.com/rui-ktei/socrates/internal/testquality"
)

func TestThePromptIsOnDisk(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n\nfunc Added() int {\n\treturn 1\n}\n")

	name := session.Review("--intent", "add a function")
	session.WaitForTerminal(name)

	kept := session.File(session.Round(name, 1).PromptPath())
	delivered := session.Simulator.OnlyInvocation().Prompt

	if kept != delivered {
		t.Errorf("what the round kept is not what the reviewer was given (-kept +delivered):\n%s", cmp.Diff(kept, delivered))
	}
	if kept == "" {
		t.Error("the round kept an empty prompt")
	}
}

func TestEveryRuleReachesTheReviewer(t *testing.T) {
	// scenario: Assertions are made on the delivered prompt
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n\nfunc Added() int {\n\treturn 1\n}\n")

	name := session.Review("--intent", "add a function")
	session.WaitForTerminal(name)

	delivered := session.Simulator.OnlyInvocation().Prompt

	rules := []struct {
		rule   string
		needle string
	}{
		{rule: "counterexample tracing", needle: "construct at least one concrete input or state and trace it through the code"},
		{rule: "a full pass", needle: "do not stop at the first valid finding"},
		{rule: "the anti-astronaut clamp", needle: "Do not infer a systemic flaw from code shape, duplication or architectural preference alone"},
		{rule: "a concrete path for a design finding", needle: "concrete reachable path"},
		{rule: "no test running", needle: "Do not run tests"},
		{rule: "no deterministic-tool findings", needle: "Do not report styling, formatting, linting, compilation or type-checking issues"},
		{rule: "anchoring", needle: "one-indexed line"},
		{rule: "emptiness", needle: "return an empty findings array"},
		{rule: "the action definitions", needle: "When in doubt the action is"},
		{rule: "the test-quality rule", needle: strings.TrimSpace(testquality.Rule)},
	}

	for _, required := range rules {
		if !strings.Contains(delivered, required.needle) {
			t.Errorf("the delivered prompt does not carry %s", required.rule)
		}
	}
}

func TestIntentReachesTheReviewer(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	const intent = "stop the parser accepting an empty document"
	name := session.Review("--intent", intent)
	session.WaitForTerminal(name)

	sections := prompt.Sections(session.Simulator.OnlyInvocation().Prompt)
	stated, present := prompt.BodyOf(sections, prompt.IntentHeading)
	if !present {
		t.Fatalf("the delivered prompt has no %q section", prompt.IntentHeading)
	}
	if !strings.Contains(stated, intent) {
		t.Errorf("the intent section does not carry what the person stated:\n%s", stated)
	}
}

func TestNoIntentIsStated(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n")

	name := session.Review()
	session.WaitForTerminal(name)

	delivered := session.Simulator.OnlyInvocation().Prompt
	if _, present := prompt.BodyOf(prompt.Sections(delivered), prompt.IntentHeading); present {
		t.Errorf("a review with no stated intent delivered an intent section anyway:\n%s", delivered)
	}
	if strings.Contains(delivered, prompt.IntentHeading) {
		t.Errorf("the intent heading reached the reviewer with nothing under it:\n%s", delivered)
	}
}

func TestTheAuthorsRationaleHasNoChannel(t *testing.T) {
	session := harness.Start(t)

	documented, _, code := session.Socrates("review", "--help")
	if code != 0 {
		t.Fatalf("socrates review --help exited %d", code)
	}

	taking := flagsTakingAValue(documented)
	if diff := cmp.Diff([]string{"--intent"}, taking); diff != "" {
		t.Errorf("the free text a caller may supply differs (-want +got):\n%s\n%s", diff, documented)
	}

	for _, invented := range []string{"--rationale", "--summary", "--notes", "--author-notes", "--session", "--reasoning"} {
		_, failures, code := session.Socrates("review", invented, "I refactored the parser and it all works")
		if code == 0 {
			t.Errorf("socrates review %s was accepted", invented)
		}
		if !strings.Contains(failures, invented) {
			t.Errorf("refusing %s does not name it: %s", invented, failures)
		}
	}

	const positional = "I refactored the parser and it all works"
	_, failures, code := session.Socrates("review", positional)
	if code == 0 {
		t.Error("socrates review accepted free text as a positional argument")
	}
	if !strings.Contains(failures, positional) {
		t.Errorf("refusing a positional argument does not name what was refused: %s", failures)
	}
}

func flagsTakingAValue(documented string) []string {
	var taking []string

	for line := range strings.SplitSeq(documented, "\n") {
		fields := strings.Fields(line)
		for position, field := range fields {
			if !strings.HasPrefix(field, "--") || position+1 >= len(fields) {
				continue
			}
			switch fields[position+1] {
			case "string", "int", "duration", "stringArray", "stringSlice":
				taking = append(taking, field)
			}
		}
	}

	return taking
}

func TestASmallChangeCarriesItsDiff(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n\nfunc Added() int {\n\treturn 1\n}\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.Diff != prompt.Inlined {
		t.Errorf("the run records the diff as %q, want %q", run.Diff, prompt.Inlined)
	}

	delivered := session.Simulator.OnlyInvocation().Prompt
	for _, carried := range []string{"diff --git a/internal/added.go b/internal/added.go", "+func Added() int {", "+\treturn 1"} {
		if !strings.Contains(delivered, carried) {
			t.Errorf("the delivered prompt does not carry %q:\n%s", carried, delivered)
		}
	}
}

func TestALargeChangeIsRead(t *testing.T) {
	session := harness.Start(t)

	var generated strings.Builder
	for line := range 6000 {
		generated.WriteString("var Generated" + strconv.Itoa(line) + " = \"a line nobody wrote by hand\"\n")
	}
	session.Repository.Write("internal/generated/generated.go", "package generated\n\n"+generated.String())

	if generated.Len() <= prompt.DefaultCeiling {
		t.Fatalf("the generated change is %d bytes, which is under the %d byte ceiling", generated.Len(), prompt.DefaultCeiling)
	}

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.Diff != prompt.Read {
		t.Errorf("the run records the diff as %q, want %q", run.Diff, prompt.Read)
	}

	delivered := session.Simulator.OnlyInvocation().Prompt
	if strings.Contains(delivered, "var Generated5999") {
		t.Error("the delivered prompt carries a diff that is over the ceiling")
	}
	if len(delivered) > prompt.DefaultCeiling {
		t.Errorf("the delivered prompt is %d bytes, over the %d byte ceiling it was meant to respect", len(delivered), prompt.DefaultCeiling)
	}
	for _, named := range []string{"internal/generated/generated.go", run.Base, run.Snapshot} {
		if !strings.Contains(delivered, named) {
			t.Errorf("the delivered prompt does not name %s:\n%s", named, delivered)
		}
	}
}

func TestTheCeilingIsOneSetting(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("socrates.yaml", "inlined_diff_ceiling: 1\n")
	session.Repository.Write(".socrates.yaml", "inlined_diff_ceiling: 1\n")
	session.Repository.Write(".socrates/config.yaml", "inlined_diff_ceiling: 1\n")
	session.Repository.Write("internal/added.go", "package internal\n\nfunc Added() int {\n\treturn 1\n}\n")

	name := session.Review()
	run := session.WaitForTerminal(name)

	if run.Diff != prompt.Inlined {
		t.Errorf("a repository lowered the ceiling to one byte and the diff was delivered as %q", run.Diff)
	}
	if !strings.Contains(session.Simulator.OnlyInvocation().Prompt, "+func Added() int {") {
		t.Error("a repository lowered the ceiling and the diff stopped reaching the reviewer")
	}
}

func TestOrderIsStableAcrossTwoReviews(t *testing.T) {
	session := harness.Start(t)
	session.Repository.Write("internal/added.go", "package internal\n\nfunc Added() int {\n\treturn 1\n}\n")

	first := session.Review("--intent", "add a function")
	session.WaitForTerminal(first)

	waitForTheNextSecond()

	second := session.Review("--intent", "add a function")
	session.WaitForTerminal(second)

	invocations := session.Simulator.Invocations()
	if len(invocations) != 2 {
		t.Fatalf("the reviewer was invoked %d times, want twice", len(invocations))
	}
	if invocations[0].Prompt != invocations[1].Prompt {
		t.Errorf("two reviews of the same change delivered different prompts (-first +second):\n%s",
			cmp.Diff(invocations[0].Prompt, invocations[1].Prompt))
	}

	if session.File(session.Round(first, 1).PromptPath()) != session.File(session.Round(second, 1).PromptPath()) {
		t.Error("the two rounds kept different prompts")
	}
}

func waitForTheNextSecond() {
	time.Sleep(time.Until(time.Now().Truncate(time.Second).Add(time.Second + 20*time.Millisecond)))
}
