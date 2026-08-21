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

	for _, named := range skill.Named() {
		if named.Flag == "" {
			continue
		}
		if _, allowed := permitted[named.Flag]; !allowed {
			t.Errorf("the skill names scrat %s %s, and nothing says why it is not a way to send the reviewer the session's own account of the change", named.Command, named.Flag)
		}
	}

	for _, word := range strings.Fields(skill.Document()) {
		named := strings.Trim(word, "`.,:;\"'()[]")
		if strings.Contains(named, "/") || namesAFile(named) {
			t.Errorf("the skill names %q, and a file is somewhere a session could write its own account of the change for a reviewer to read", named)
		}
	}
}

func namesAFile(word string) bool {
	for _, extension := range []string{".md", ".txt", ".json", ".yaml", ".yml", ".log"} {
		if strings.HasSuffix(word, extension) {
			return true
		}
	}
	return false
}
