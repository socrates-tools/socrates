package convention

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/internal/testing/source"
)

func TestARefusalNamesTheNextAction(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, actionless := range wording.Actionless {
		t.Errorf("%s refuses without naming what the person can do about it", actionless)
	}
}

func TestNoEmDashReachesAPerson(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, dashed := range wording.EmDashes {
		t.Errorf("%s carries an em dash", dashed)
	}
}

func TestEveryUserFacingStringIsInItsPackagesRegistry(t *testing.T) {
	wording, err := UserFacingStrings(source.Root(t))
	if err != nil {
		t.Fatalf("reading what socrates says: %v", err)
	}

	for _, inline := range wording.Inline {
		t.Errorf("%s states a message where it is used: move it into the package's %s", inline, RegistryFile)
	}
}

func TestAMessageStatedAtTheCallSiteIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() error {
	return fmt.Errorf("the thing could not be done: %w", nil)
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "thing.go") + ":6"}
	if diff := cmp.Diff(want, wording.Inline); diff != "" {
		t.Errorf("the reported messages differ (-want +got):\n%s", diff)
	}
}

func TestAMessageInTheRegistryIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() error {
	return fmt.Errorf(couldNotBeDone, nil)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const couldNotBeDone = "the thing could not be done: %w"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a message held in the registry was reported: %v", wording.Inline)
	}
}

func TestAFormatWithNoProseIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "fmt"

func Do() string {
	return fmt.Sprintf("%s/%d", "a", 1)
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a format string carrying no prose was reported: %v", wording.Inline)
	}
}

func TestARefusalWithNoNextActionIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "github.com/socrates-tools/socrates/internal/refusal"

func Do() error {
	return refusal.New(somethingWentWrong)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const somethingWentWrong = "the thing is in a state socrates cannot use"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "thing.go") + ":6"}
	if diff := cmp.Diff(want, wording.Actionless); diff != "" {
		t.Errorf("the reported refusals differ (-want +got):\n%s", diff)
	}
}

func TestARefusalThatNamesTheNextActionIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

import "github.com/socrates-tools/socrates/internal/refusal"

func Do() error {
	return refusal.New(somethingWentWrong)
}
`)
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const somethingWentWrong = "the thing is in a state socrates cannot use: fix the thing and run the command again"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Actionless) != 0 {
		t.Errorf("a refusal naming what to do next was reported: %v", wording.Actionless)
	}
}

func TestAnEmDashIsReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "messages.go"), `package thing

const dashed = "socrates refuses — and here is why"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "thing", "messages.go") + ":3"}
	if diff := cmp.Diff(want, wording.EmDashes); diff != "" {
		t.Errorf("the reported em dashes differ (-want +got):\n%s", diff)
	}
}

func TestALabelAMenuDrawsIsReportedWhereItIsStated(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("review", "configure", "view.go"), `package configure

func lines() []string {
	return []string{"What do you want to configure?"}
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("review", "configure", "view.go") + ":4"}
	if diff := cmp.Diff(want, wording.Inline); diff != "" {
		t.Errorf("the reported labels differ (-want +got):\n%s", diff)
	}
}

func TestALabelAMenuDrawsFromItsRegistryIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("review", "configure", "view.go"), `package configure

func lines() []string {
	return []string{whatToConfigure}
}
`)
	write(t, root, filepath.Join("review", "configure", "messages.go"), `package configure

const whatToConfigure = "What do you want to configure?"
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a label held in the registry was reported: %v", wording.Inline)
	}
}

func TestALabelAPackageOutsideTheMenuStatesInPassingIsNotReported(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "thing", "thing.go"), `package thing

func Names() []string {
	return []string{"a name with spaces in it"}
}
`)

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a literal in a package that builds no view was reported: %v", wording.Inline)
	}
}

func TestALabelAMenuDrawsFromARawStringIsReportedToo(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("review", "configure", "view.go"), "package configure\n\nfunc lines() []string {\n\treturn []string{`What do you want to configure?`}\n}\n")

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("review", "configure", "view.go") + ":4"}
	if diff := cmp.Diff(want, wording.Inline); diff != "" {
		t.Errorf("the reported labels differ (-want +got):\n%s", diff)
	}
}

func TestThePatternAMenuPackageMatchesOnIsNotALabel(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("review", "configure", "document.go"), "package configure\n\nimport \"regexp\"\n\nvar blankLine = regexp.MustCompile(`(?m)^[ \\t]*\\r?\\n`)\n")

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}
	if len(wording.Inline) != 0 {
		t.Errorf("a pattern the package matches on was reported as a label: %v", wording.Inline)
	}
}

func TestALabelTheSharedMenuDrawsFromARawStringIsReportedToo(t *testing.T) {
	root := t.TempDir()
	write(t, root, filepath.Join("internal", "settingsmenu", "view.go"), "package settingsmenu\n\nfunc lines() []string {\n\treturn []string{`What do you want to configure?`}\n}\n")

	wording, err := UserFacingStrings(root)
	if err != nil {
		t.Fatalf("reading the fabricated package: %v", err)
	}

	want := []string{filepath.Join("internal", "settingsmenu", "view.go") + ":4"}
	if diff := cmp.Diff(want, wording.Inline); diff != "" {
		t.Errorf("the reported labels differ (-want +got):\n%s", diff)
	}
}

const clientModule = "held/messages.ts"

func clientWordingOf(t *testing.T, content string) (Wording, error) {
	t.Helper()

	root := t.TempDir()
	write(t, root, clientModule, content)
	return ClientWording(root, []string{clientModule})
}

func TestAnEmDashInTheClientsMessagesIsReportedWithTheModuleAndTheLine(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {}\nexport const said = {\n  greeting: 'Hello — there',\n}\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if diff := cmp.Diff([]string{clientModule + ":3"}, wording.EmDashes); diff != "" {
		t.Errorf("the reported em dashes differ (-want +got):\n%s", diff)
	}
}

func TestAnEmDashWrittenAsAnEscapeInTheClientsMessagesIsReported(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {}\nexport const said = \"Hello \\u2014 there\"\nexport const also = 'a \\u{2014} b'\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	want := []string{clientModule + ":2", clientModule + ":3"}
	if diff := cmp.Diff(want, wording.EmDashes); diff != "" {
		t.Errorf("the reported em dashes differ (-want +got):\n%s", diff)
	}
}

func TestAnEmDashInATemplateOfTheClientsMessagesIsReported(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {}\nexport const said = (n: number) => `${n} sent — waiting`\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if diff := cmp.Diff([]string{clientModule + ":2"}, wording.EmDashes); diff != "" {
		t.Errorf("the reported em dashes differ (-want +got):\n%s", diff)
	}
}

func TestAClientsMessageWithNoEmDashIsNotReported(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {}\nexport const said = 'Sent - waiting'\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if len(wording.EmDashes) != 0 {
		t.Errorf("a message with no em dash was reported: %v", wording.EmDashes)
	}
}

func TestAClientRefusalWithNoNextActionIsReported(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {\n  gone: 'The page lost its connection.',\n  sent: \"Sending failed.\",\n  quoted: `Nothing was sent.`,\n}\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	want := []string{clientModule + ":2", clientModule + ":3", clientModule + ":4"}
	if diff := cmp.Diff(want, wording.Actionless); diff != "" {
		t.Errorf("the reported refusals differ (-want +got):\n%s", diff)
	}
}

func TestAClientRefusalThatNamesTheNextActionIsNotReported(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {\n  gone: 'The page lost its connection. Check it and send again.',\n  'quoted-key': 'Check the page.',\n}\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if len(wording.Actionless) != 0 {
		t.Errorf("a refusal naming what to do next was reported: %v", wording.Actionless)
	}
}

func TestAClientRefusalWrittenAsATemplateIsJudgedAsOneMessage(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {\n  gone: (n: number) => `${n} files were refused.`,\n  kept: (n: number) => `${n} files were refused. Remove them and send again.`,\n}\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if diff := cmp.Diff([]string{clientModule + ":2"}, wording.Actionless); diff != "" {
		t.Errorf("the reported refusals differ (-want +got):\n%s", diff)
	}
}

func TestAClientMessageOutsideTheRefusalsNeedsNoNextAction(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {\n  gone: 'Check the connection.',\n}\nexport const said = {\n  open: 'Open',\n  copied: 'Copied',\n}\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if len(wording.Actionless) != 0 {
		t.Errorf("a message that is not a refusal was reported: %v", wording.Actionless)
	}
}

func TestARefusalsGroupWithNestedBracesIsReadToItsEnd(t *testing.T) {
	wording, err := clientWordingOf(t, "export const refusals = {\n  gone: { text: 'Check the connection.' },\n  sent: 'Sending failed.',\n} as const\nexport const said = 'Sending failed.'\n")
	if err != nil {
		t.Fatalf("reading the fabricated module: %v", err)
	}

	if diff := cmp.Diff([]string{clientModule + ":3"}, wording.Actionless); diff != "" {
		t.Errorf("the reported refusals differ (-want +got):\n%s", diff)
	}
}

func TestAClientsMessagesModuleThatIsGoneIsAnError(t *testing.T) {
	_, err := ClientWording(t.TempDir(), []string{clientModule})
	if err == nil {
		t.Fatal("a module that does not exist was read as saying nothing")
	}
}

func TestAClientsMessagesModuleWithNoRefusalsGroupIsAnError(t *testing.T) {
	_, err := clientWordingOf(t, "export const said = 'Open'\n")
	if err == nil {
		t.Fatal("a module declaring no refusals group was read as holding no refusals")
	}
}

func TestAClientsMessagesModuleWithAnUnfinishedStringIsAnError(t *testing.T) {
	_, err := clientWordingOf(t, "export const refusals = {}\nexport const said = 'Open\n")
	if err == nil {
		t.Fatal("a module with an unterminated string was read")
	}
}

func TestTheClientsMessagesCarryNoEmDash(t *testing.T) {
	wording, err := ClientWording(source.Root(t), ClientMessages)
	if err != nil {
		t.Fatalf("reading what the page says: %v", err)
	}

	for _, dashed := range wording.EmDashes {
		t.Errorf("%s carries an em dash", dashed)
	}
}

func TestAClientRefusalNamesTheNextAction(t *testing.T) {
	wording, err := ClientWording(source.Root(t), ClientMessages)
	if err != nil {
		t.Fatalf("reading what the page says: %v", err)
	}

	for _, actionless := range wording.Actionless {
		t.Errorf("%s refuses without naming what the person can do about it", actionless)
	}
}
