package skill

import (
	"fmt"
	"strings"

	"github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/dialogue"
)

const (
	Name         = "socrates-review"
	HandoverName = Name + "-handover"
	Surface      = "agora"
)

var (
	Ceilings         = skill.Ceilings{skill.ClaudeCode: 15708, skill.Codex: 16257}
	HandoverCeilings = skill.Ceilings{skill.ClaudeCode: 17000, skill.Codex: 17549}
)

var Namespaces = map[string]string{
	"review": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document(host skill.Host) string {
	return assembled(skill.FrontMatter(host, Name, description), purpose, host.Guidance(reviewPage, true), rules(), thePage(theAnswerAPageCarries, passingTheAttachments), fmt.Sprintf(theLoop, host.Flag), prohibitions)
}

func thePage(answer, attachments string) string {
	return fmt.Sprintf(pageTemplate, answer, attachments)
}

func assembled(parts ...string) string {
	return skill.Assembled(parts...)
}

func rules() string {
	return stating(dialogue.Rules())
}

func stating(rules []string) string {
	stated := []string{rulesHeading, ""}
	for _, rule := range rules {
		stated = append(stated, fmt.Sprintf(ruleLine, rule))
	}
	return strings.Join(stated, "\n")
}

func For(host skill.Host) []skill.Skill {
	return []skill.Skill{
		{Name: Name, Host: host, Document: Document, Ceiling: Ceilings.For(host)},
		{Name: HandoverName, Host: host, Document: HandoverDocument, Ceiling: HandoverCeilings.For(host)},
	}
}

func All() []skill.Skill {
	var every []skill.Skill
	for _, host := range skill.Hosts() {
		every = append(every, For(host)...)
	}
	return every
}
