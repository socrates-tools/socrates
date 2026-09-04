package prompt

import (
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates-review/internal/agent"
	"github.com/socrates-tools/socrates-review/internal/findings"
	"github.com/socrates-tools/socrates-review/internal/testquality"
)

func smallChange() Request {
	return Request{
		Intent:     "make the parser reject an empty document",
		Repository: "/home/somebody/work/example",
		Base:       "1111111111111111111111111111111111111111",
		Snapshot:   "2222222222222222222222222222222222222222",
		Changed:    []string{"internal/parse/parse.go", "internal/parse/parse_test.go"},
		Diff:       "diff --git a/internal/parse/parse.go b/internal/parse/parse.go\n+\tif len(document) == 0 {\n+\t\treturn errors.New(\"an empty document\")\n+\t}\n",
	}
}

func TestOrderIsStableForTheSameRequest(t *testing.T) {
	request := smallChange()
	first := Build(request).Text

	for attempt := range 5 {
		again := Build(request).Text
		if again != first {
			t.Fatalf("build %d differs from the first:\n%s", attempt+2, cmp.Diff(first, again))
		}
	}
}

func TestSectionsAreDeliveredInOneFixedOrder(t *testing.T) {
	want := []string{
		TaskHeading,
		ContextHeading,
		IntentHeading,
		ChangedHeading,
		ChangeHeading,
		ReviewingHeading,
		TestQualityHeading,
		ReportingHeading,
	}

	got := Headings(Sections(Build(smallChange()).Text))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("the delivered sections differ (-want +got):\n%s", diff)
	}
}

func TestAnEmptySectionIsAbsent(t *testing.T) {
	request := smallChange()
	request.Intent = "   \n\t "

	delivered := Build(request).Text
	sections := Sections(delivered)

	if _, present := BodyOf(sections, IntentHeading); present {
		t.Errorf("an unstated intent still delivered the %q heading:\n%s", IntentHeading, delivered)
	}
	if strings.Contains(delivered, IntentHeading) {
		t.Errorf("the intent heading survives somewhere in the prompt:\n%s", delivered)
	}

	for _, part := range sections {
		if strings.TrimSpace(part.Body) == "" {
			t.Errorf("the section %q was delivered with nothing under it", part.Heading)
		}
	}
}

func TestIntentReachesTheReviewerAsAClaimAndNotAsFact(t *testing.T) {
	request := smallChange()
	sections := Sections(Build(request).Text)

	stated, present := BodyOf(sections, IntentHeading)
	if !present {
		t.Fatal("a stated intent delivered no intent section")
	}
	if !strings.Contains(stated, request.Intent) {
		t.Errorf("the intent section does not carry what the person said:\n%s", stated)
	}
	if !strings.Contains(stated, "say") {
		t.Errorf("the intent section does not mark the intent as something claimed:\n%s", stated)
	}
}

func TestNoIntentIsStated(t *testing.T) {
	request := smallChange()
	request.Intent = ""

	sections := Sections(Build(request).Text)
	if _, present := BodyOf(sections, IntentHeading); present {
		t.Error("a review with no stated intent still delivered an intent section")
	}

	for _, heading := range []string{TaskHeading, ContextHeading, ChangedHeading, ChangeHeading, ReviewingHeading, TestQualityHeading, ReportingHeading} {
		if _, present := BodyOf(sections, heading); !present {
			t.Errorf("the missing intent took %q with it", heading)
		}
	}
}

func TestASmallChangeCarriesItsDiff(t *testing.T) {
	request := smallChange()
	delivered := Build(request)

	if delivered.Diff != Inlined {
		t.Errorf("a %d byte diff under the %d byte ceiling was delivered as %q", len(request.Diff), DefaultCeiling, delivered.Diff)
	}

	change, present := BodyOf(Sections(delivered.Text), ChangeHeading)
	if !present {
		t.Fatal("the prompt has no section for the change itself")
	}
	if !strings.Contains(change, request.Diff) {
		t.Errorf("the change section does not carry the diff:\n%s", change)
	}
}

func TestALargeChangeIsRead(t *testing.T) {
	request := smallChange()
	request.Ceiling = 64
	request.Diff = strings.Repeat("+a line of a generated file\n", 100)

	delivered := Build(request)
	if delivered.Diff != Read {
		t.Errorf("a %d byte diff over the %d byte ceiling was delivered as %q", len(request.Diff), request.Ceiling, delivered.Diff)
	}

	change, present := BodyOf(Sections(delivered.Text), ChangeHeading)
	if !present {
		t.Fatal("the prompt has no section for the change itself")
	}
	if strings.Contains(change, request.Diff) {
		t.Errorf("the change section carries the diff that is over the ceiling:\n%s", change)
	}
	for _, named := range []string{request.Base, request.Snapshot, strconv.Itoa(request.Ceiling)} {
		if !strings.Contains(change, named) {
			t.Errorf("the change section does not name %s:\n%s", named, change)
		}
	}

	listed, present := BodyOf(Sections(delivered.Text), ChangedHeading)
	if !present {
		t.Fatal("the prompt has no section for the files the change touches")
	}
	for _, path := range request.Changed {
		if !strings.Contains(listed, path) {
			t.Errorf("the file list does not carry %s:\n%s", path, listed)
		}
	}
}

func TestTheCeilingIsOneSetting(t *testing.T) {
	request := smallChange()
	request.Diff = strings.Repeat("+a line\n", 1000)

	if delivered := Build(request); delivered.Diff != Inlined {
		t.Errorf("with no ceiling stated, a %d byte diff was delivered as %q, and the default is %d bytes",
			len(request.Diff), delivered.Diff, DefaultCeiling)
	}

	request.Ceiling = len(request.Diff) - 1
	if delivered := Build(request); delivered.Diff != Read {
		t.Errorf("with the ceiling one byte under the diff, it was delivered as %q", delivered.Diff)
	}

	request.Ceiling = len(request.Diff)
	if delivered := Build(request); delivered.Diff != Inlined {
		t.Errorf("with the ceiling exactly the size of the diff, it was delivered as %q", delivered.Diff)
	}
}

func TestTheActionDefinitionsAreTheOnesTheClassifierImplements(t *testing.T) {
	reported, present := BodyOf(Sections(Build(smallChange()).Text), ReportingHeading)
	if !present {
		t.Fatal("the prompt says nothing about what to report")
	}

	defined := quotedValuesIn(reported)

	for _, action := range findings.Actions() {
		if !slices.Contains(defined, string(action)) {
			t.Errorf("the prompt does not define %q, which the classifier accepts", action)
		}
	}
	for _, value := range defined {
		if slices.Contains(findings.Actions(), findings.Action(value)) {
			continue
		}
		if slices.Contains(findings.Severities(), findings.Severity(value)) {
			continue
		}
		if slices.Contains(findings.Risks(), findings.Risk(value)) {
			continue
		}
		t.Errorf("the prompt defines %q, which is not a severity, an action or a risk level the code has", value)
	}

	unrecognised := findings.Finding{
		ID:          "f1",
		Severity:    findings.Error,
		File:        "internal/parse/parse.go",
		Line:        1,
		Description: "an action nobody stated",
		Action:      "whatever-the-reviewer-felt-like",
	}
	stated := defaultActionIn(t, reported)
	if applied := unrecognised.EffectiveAction(); applied != stated {
		t.Errorf("the prompt says the default action is %q, and the classifier applies %q", stated, applied)
	}
}

func defaultActionIn(t *testing.T, reported string) findings.Action {
	t.Helper()

	const preamble = "When in doubt the action is "
	_, after, found := strings.Cut(reported, preamble)
	if !found {
		t.Fatalf("the prompt does not say what the action is when in doubt:\n%s", reported)
	}

	quoted := quotedValuesIn(after)
	if len(quoted) == 0 {
		t.Fatalf("the prompt names no action after %q:\n%s", preamble, after)
	}
	return findings.Action(quoted[0])
}

func quotedValuesIn(text string) []string {
	var quoted []string
	remaining := text

	for {
		_, after, found := strings.Cut(remaining, `"`)
		if !found {
			return quoted
		}
		value, rest, closed := strings.Cut(after, `"`)
		if !closed {
			return quoted
		}
		if value != "" && !strings.Contains(value, " ") {
			quoted = append(quoted, value)
		}
		remaining = rest
	}
}

func TestEverySeverityAndRiskLevelTheCodeHasIsDefinedForTheReviewer(t *testing.T) {
	reported, present := BodyOf(Sections(Build(smallChange()).Text), ReportingHeading)
	if !present {
		t.Fatal("the prompt says nothing about what to report")
	}

	for _, severity := range findings.Severities() {
		if !strings.Contains(reported, strconv.Quote(string(severity))) {
			t.Errorf("the prompt does not define the severity %q", severity)
		}
	}
	for _, risk := range findings.Risks() {
		if !strings.Contains(reported, strconv.Quote(string(risk))) {
			t.Errorf("the prompt does not define the risk level %q", risk)
		}
	}

	if len(severityDefinitions) != len(findings.Severities()) {
		t.Errorf("the prompt defines %d severities and the code has %d", len(severityDefinitions), len(findings.Severities()))
	}
	if len(actionDefinitions) != len(findings.Actions()) {
		t.Errorf("the prompt defines %d actions and the code has %d", len(actionDefinitions), len(findings.Actions()))
	}
	if len(riskDefinitions) != len(findings.Risks()) {
		t.Errorf("the prompt defines %d risk levels and the code has %d", len(riskDefinitions), len(findings.Risks()))
	}
}

func TestTheSeverityThatBlocksIsTheSeverityTheCodeBlocksOn(t *testing.T) {
	reported, present := BodyOf(Sections(Build(smallChange()).Text), ReportingHeading)
	if !present {
		t.Fatal("the prompt says nothing about what to report")
	}

	for _, severity := range findings.Severities() {
		line := lineNaming(reported, strconv.Quote(string(severity)))
		if line == "" {
			t.Fatalf("the prompt does not define the severity %q", severity)
		}
		if strings.Contains(line, "holds the change") != severity.Blocks() {
			t.Errorf("the prompt says of %q: %q, and the code blocks on it: %t", severity, line, severity.Blocks())
		}
	}
}

func lineNaming(text, needle string) string {
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, needle) {
			return line
		}
	}
	return ""
}

func TestOneBlockManyRoles(t *testing.T) {
	roles := agent.Roles()
	if len(roles) != 1 || roles[0] != agent.Reviewer {
		t.Fatalf("socrates now has the roles %v, and every role that can write, repair or review tests belongs in this test", roles)
	}

	body, present := BodyOf(Sections(Build(smallChange()).Text), TestQualityHeading)
	if !present {
		t.Fatalf("the prompt for %s carries no test-quality section", roles[0])
	}
	shared := strings.TrimSpace(testquality.Rule)
	if !strings.HasPrefix(body, shared) {
		t.Errorf("the test-quality block delivered to %s is not the shared block (-shared +delivered):\n%s",
			roles[0], cmp.Diff(shared, body))
	}
}

func TestSectionsRoundTripThroughTheDeliveredPrompt(t *testing.T) {
	delivered := Build(smallChange()).Text

	var rebuilt strings.Builder
	for position, part := range Sections(delivered) {
		if position > 0 {
			rebuilt.WriteString("\n")
		}
		rebuilt.WriteString(headingMarker + part.Heading + "\n\n" + part.Body + "\n")
	}

	if rebuilt.String() != delivered {
		t.Errorf("reading the sections back does not reproduce the prompt (-delivered +read):\n%s",
			cmp.Diff(delivered, rebuilt.String()))
	}
}

func TestTheParsedTierIsToldHowToDeliverTheDocument(t *testing.T) {
	asked := smallChange()
	asked.DocumentEndsTheMessage = true

	extracting := Build(asked)
	constrained := Build(smallChange())

	instruction := findings.DeliveryInstruction()
	if !strings.Contains(extracting.Text, instruction) {
		t.Errorf("the delivered prompt for a reviewer whose document is extracted does not say how to deliver it:\n%s", extracting.Text)
	}
	if strings.Contains(constrained.Text, instruction) {
		t.Error("the delivered prompt for a reviewer whose vendor constrains the schema asks it to end its message with the document as well")
	}

	reporting, found := bodyOf(Sections(extracting.Text), ReportingHeading)
	if !found {
		t.Fatalf("the delivered prompt carries no %q section", ReportingHeading)
	}
	if !strings.Contains(reporting, instruction) {
		t.Errorf("the delivery instruction is not in the section that says what to report:\n%s", reporting)
	}
}

func bodyOf(sections []Section, heading string) (string, bool) {
	for _, part := range sections {
		if part.Heading == heading {
			return part.Body, true
		}
	}
	return "", false
}
