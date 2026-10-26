package skill

import (
	"fmt"

	"github.com/socrates-tools/socrates/internal/skill"
	"github.com/socrates-tools/socrates/review/dialogue"
)

func HandoverDocument(host skill.Host) string {
	return assembled(skill.FrontMatter(host, HandoverName, handoverDescription), handoverPurpose, host.Guidance(handoverPage, true), handoverRules(), thePage(theAnswerAHandoverPageCarries, ""), fmt.Sprintf(theHandoverLoop, host.Flag), handoverProhibitions)
}

func handoverRules() string {
	return stating(dialogue.HandoverRules())
}
