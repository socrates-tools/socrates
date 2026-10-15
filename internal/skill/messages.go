package skill

const (
	noSuchBinary  = "the skill names %s %s, and this repository builds no binary called %s"
	noSuchCommand = "the skill names %s %s, and %s offers no such command"
	noSuchFlag    = "the skill names %s %s %s, and that command offers no such flag"

	skillOnAHost = "%s on %s"

	theSkillIsTooLarge   = "the skill renders %d bytes and its ceiling is %d: every session that reads it reads all of it, so move the reasoning for a rule into docs/adr/ rather than raising this"
	theSkillHasNoCeiling = "%s states no ceiling, so nothing bounds what every session that reads it reads"

	theTargetCannotBeWritten = "the skill could not be written to %s: %s. Name a directory you can write with %s --host <host> --to"

	installedForAPerson = "the %s skill is at %s, and only a person starts it, by typing %s"

	aDirectoryNeedsAHost = "--%s names a directory, and the directory holds one host's rendering. Choose whose with --%s, as %s"
	noSuchHost           = "--%s %q names no host. Use %s"

	hostFlagStated = "--%s %s"
	hostFlagsJoin  = " or "

	hostUsage = "write only the rendering for this host, %s, to its own directory, instead of every host's"

	noHomeDirectory = "the skill installs under a home directory by default, and none could be found: %s. Name a directory to write it to with %s --host and --to"
)
