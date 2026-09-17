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

	Ceiling         = 15677
	HandoverCeiling = 16969
)

var Namespaces = map[string]string{
	"review": Name,
	"agora":  Surface,
}

func Binaries() []string {
	return []string{Name, Surface}
}

func Document() string {
	return assembled(frontMatter, purpose, rules(), thePage(theAnswerAPageCarries), theLoop, prohibitions)
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

func All() []skill.Skill {
	return []skill.Skill{
		{Name: Name, Document: Document, Ceiling: Ceiling},
		{Name: HandoverName, Document: HandoverDocument, Ceiling: HandoverCeiling},
	}
}
