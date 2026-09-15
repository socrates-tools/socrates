package skill_test

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/socrates-tools/socrates/socrates-review/skill"
)

func TestProseNamingTheProductIsNotAnInvocation(t *testing.T) {
	var commands []string
	for _, named := range skill.Named() {
		attributed := named.Binary + " " + named.Command
		if !slices.Contains(commands, attributed) {
			commands = append(commands, attributed)
		}
	}
	slices.Sort(commands)

	want := []string{
		"agora poll",
		"agora serve",
		"socrates-review continue",
		"socrates-review findings",
		"socrates-review respond",
		"socrates-review review",
		"socrates-review status",
		"socrates-review watch",
	}
	if diff := cmp.Diff(want, commands); diff != "" {
		t.Errorf("the commands the skill is read to name differ (-want +got):\n%s\na prose sentence about socrates is not an invocation, and an invocation is written as code", diff)
	}
}
