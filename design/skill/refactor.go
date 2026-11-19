package skill

import (
	"fmt"
	"path/filepath"

	"github.com/socrates-tools/socrates/internal/skill"
)

const RefactorName = "socrates-refactor"

func RefactorDocument(host skill.Host) string {
	return assembled(skill.FrontMatter(host, RefactorName, refactorDescription), refactorPurpose, theComposedSkillsSection, installingWhatIsMissing(host),
		theRunningSection, refactorEscalation(host), refactorInstallingElsewhere(host))
}

func installingWhatIsMissing(host skill.Host) string {
	return fmt.Sprintf(installingWhatIsMissingTemplate, filepath.ToSlash(host.SkillsDirectory), host.SkillsLocation(), host.SkillsCLIAgent, host.SkillsCLIAgent)
}

func refactorEscalation(host skill.Host) string {
	if host != skill.Codex {
		return ""
	}
	return fmt.Sprintf(refactorEscalationTemplate, host.Name)
}

func refactorInstallingElsewhere(host skill.Host) string {
	return fmt.Sprintf(refactorInstallingElsewhereTemplate, host.SkillsLocation(), host.Flag)
}
