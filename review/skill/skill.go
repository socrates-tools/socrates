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

	Ceiling         = 15708
	HandoverCeiling = 17000
)

var Namespaces = map[string]string{
	"review": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document(host skill.Host) string {
	return assembled(skill.FrontMatter(host, Name, description), purpose, rules(), thePage(theAnswerAPageCarries), theLoop, prohibitions)
}

func thePage(answer string) string {
	return fmt.Sprintf(pageTemplate, answer)
}

func assembled(parts ...string) string {
	return strings.ReplaceAll(strings.Join(parts, "\n\n")+"\n", skill.InlineCodeMarker, skill.CodeMarker)
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
		{Name: Name, Host: host, Document: Document, Ceiling: Ceiling},
		{Name: HandoverName, Host: host, Document: HandoverDocument, Ceiling: HandoverCeiling},
	}
}

func All() []skill.Skill {
	var every []skill.Skill
	for _, host := range skill.Hosts() {
		every = append(every, For(host)...)
	}
	return every
}
