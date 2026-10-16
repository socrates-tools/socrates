package e2e_test

import (
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"

	designskill "github.com/socrates-tools/socrates/design/skill"
	skillkit "github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/skill"
)

var illustrationTheSkillShows = regexp.MustCompile(`"illustration": (\{"path": "[^"]+"\})`)

func TestTheIllustrationEachSkillShowsIsServedWithoutRefusal(t *testing.T) {
	// scenario: An agent following either skill writes an illustration agora serves
	for name, document := range map[string]string{
		"design": designskill.Document(skillkit.ClaudeCode),
		"review": skill.Document(skillkit.ClaudeCode),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			shown := illustrationTheSkillShows.FindStringSubmatch(document)
			if shown == nil {
				t.Fatalf("the %s skill shows no illustration field", name)
			}

			held := standing(t)
			aDrawing(t, held, "q3.html", "<!doctype html><p>drawn</p>")

			address := held.serve(strings.Replace(oneDecisionDocument, `"mode": "single-choice",`, `"illustration": `+shown[1]+`, "mode": "single-choice",`, 1))

			answered, err := http.Get(address + "/illustrations/c-1/q3.html")
			if err != nil {
				t.Fatalf("reaching the surface: %v", err)
			}
			defer func() { _ = answered.Body.Close() }()
			served, _ := io.ReadAll(answered.Body)
			if answered.StatusCode != http.StatusOK || string(served) != `<!doctype html><p>drawn</p><script src="/assets/height.js"></script>` {
				t.Errorf("the illustration the %s skill shows was answered with %d: %s", name, answered.StatusCode, served)
			}
		})
	}
}
