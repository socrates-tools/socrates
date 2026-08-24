package skill_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rui-ktei/socrates/internal/skill"
)

func TestProseNamingTheProductIsNotAnInvocation(t *testing.T) {
	var commands []string
	for _, named := range skill.Named() {
		if !slices.Contains(commands, named.Command) {
			commands = append(commands, named.Command)
		}
	}
	slices.Sort(commands)

	want := []string{"continue", "findings", "respond", "review", "status", "watch"}
	if diff := cmp.Diff(want, commands); diff != "" {
		t.Errorf("the commands the skill is read to name differ (-want +got):\n%s\na prose sentence about socrates is not an invocation, and an invocation is written as code", diff)
	}
}

func TestNoInlineMarkerSurvivesIntoTheDocument(t *testing.T) {
	document := skill.Document()

	if counted := strings.Count(document, "|"); counted != 0 {
		t.Errorf("the document carries %d inline markers that were never turned into code", counted)
	}
	if counted := strings.Count(document, "`"); counted%2 != 0 {
		t.Errorf("the document holds %d code markers, so one code span is unclosed", counted)
	}
}
