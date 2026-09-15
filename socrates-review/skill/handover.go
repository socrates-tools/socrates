package skill

import "github.com/socrates-tools/socrates/socrates-review/dialogue"

func HandoverDocument() string {
	return assembled(handoverFrontMatter, handoverPurpose, handoverRules(), thePage(theAnswerAHandoverPageCarries), theHandoverLoop, handoverProhibitions)
}

func handoverRules() string {
	return stating(dialogue.HandoverRules())
}
