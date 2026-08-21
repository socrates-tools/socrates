package skill_test

import (
	"strings"
	"testing"

	"github.com/rui-ktei/socrates/internal/skill"
)

func TestTheThreeProhibitionsArePresent(t *testing.T) {
	document := strings.ToLower(skill.Document())

	for _, prohibition := range []struct {
		named  string
		stated []string
	}{
		{named: "socrates patches nothing", stated: []string{"patches nothing", "your own working tree"}},
		{named: "there is no way to explain the change", stated: []string{"no way to", "the design"}},
		{named: "a rejection is a claim rather than a dismissal", stated: []string{"rejection", "next reviewer", "without a reason is refused"}},
	} {
		t.Run(prohibition.named, func(t *testing.T) {
			for _, stated := range prohibition.stated {
				if !strings.Contains(document, stated) {
					t.Errorf("the skill does not say %q, and this is the prohibition that %s", stated, prohibition.named)
				}
			}
		})
	}
}

func TestNoChannelIsImpliedForTheAuthorsAccount(t *testing.T) {
	permitted := map[string]string{
		"--intent":   "the one line the person stated, which is the only account of the change socrates carries to a first-round reviewer",
		"--finding":  "which finding an answer is about",
		"--fixed":    "that a finding was patched",
		"--rejected": "that a finding is wrong",
		"--note":     "what was done about a patched finding, which reaches the next reviewer as a claim it judges",
		"--reason":   "why a rejected finding is wrong, which reaches the next reviewer as a claim it judges",
		"--to":       "where to install the skill, which reaches no reviewer at all",
	}

	commands := map[string]string{
		"review":    "starts a review, and carries one line of intent and nothing else",
		"status":    "reports on a run",
		"findings":  "reads what a reviewer said",
		"respond":   "records one answer to one finding",
		"continue":  "starts the next round from what was recorded",
		"watch":     "prints a run's events into a person's terminal",
		"cancel":    "stops a run",
		"logs":      "prints what a run recorded",
		"skill":     "installs this skill",
		"supervise": "conducts a review that scrat review started",
	}

	for _, named := range skill.Named() {
		if _, allowed := commands[named.Command]; !allowed {
			t.Errorf("the skill names scrat %s, and nothing says why it is not a way to send the reviewer the session's own account of the change", named.Command)
			continue
		}
		if named.Flag == "" {
			continue
		}
		if _, allowed := permitted[named.Flag]; !allowed {
			t.Errorf("the skill names scrat %s %s, and nothing says why it is not a way to send the reviewer the session's own account of the change", named.Command, named.Flag)
		}
	}
}
